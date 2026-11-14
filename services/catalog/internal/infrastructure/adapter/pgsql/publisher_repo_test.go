package pgsql

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/codejsha/bookstore-microservices/catalog/generated/infrastructure/port/dao"
	"github.com/codejsha/bookstore-microservices/catalog/generated/infrastructure/port/entity"
	"github.com/codejsha/bookstore-microservices/catalog/internal/application/port/repo"
	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/model/command"
	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/model/option"
)

func newPublisherRepo(db *gorm.DB) *publisherRepository {
	return &publisherRepository{q: dao.Use(db)}
}

func publisherUids(results []*repo.PublisherResult) []string {
	uids := make([]string, len(results))
	for i, r := range results {
		uids[i] = r.Uid
	}
	return uids
}

func TestPublisherRepositoryFindAll_olKeyFilterGiven_matchesExactKeyOnly(t *testing.T) {
	db := newFullCatalogDB(t)
	allen := seedPublisher(t, db, "p-1", "Allen & Unwin")
	harper := seedPublisher(t, db, "p-2", "HarperCollins")
	setColumn(t, db, "publisher", allen.Id, "ol_key", "OL1P")
	setColumn(t, db, "publisher", harper.Id, "ol_key", "OL10P")

	total, got, err := newPublisherRepo(db).FindAll(context.Background(), option.NewPublisherQueryOption(
		option.PublisherQueryOption{}.WithOlKey(strPtr("OL10P")),
	))
	if err != nil {
		t.Fatalf("FindAll: %v", err)
	}

	if total != 1 || !reflect.DeepEqual(publisherUids(got), []string{"p-2"}) {
		t.Errorf("total = %d, uids = %v, want 1 and [p-2]", total, publisherUids(got))
	}
}

func TestPublisherRepositoryFindAll_nameSortWithPaging_pagesAfterSortingAndKeepsTotal(t *testing.T) {
	db := newFullCatalogDB(t)
	seedPublisher(t, db, "p-1", "Penguin")
	seedPublisher(t, db, "p-2", "Ace")
	seedPublisher(t, db, "p-3", "Gollancz")
	gone := seedPublisher(t, db, "p-4", "Aardvark")
	softDelete(t, db, "publisher", gone.Id)
	r := newPublisherRepo(db)

	total, first, err := r.FindAll(context.Background(), option.NewPublisherQueryOption(
		option.PublisherQueryOption{}.WithName(strPtr("")),
		option.PublisherQueryOption{}.WithOlKey(strPtr("")),
		option.PublisherQueryOption{}.WithPage(pageOf(2, 0, "name")),
	))
	if err != nil {
		t.Fatalf("FindAll page 0: %v", err)
	}
	_, second, err := r.FindAll(context.Background(), option.NewPublisherQueryOption(
		option.PublisherQueryOption{}.WithPage(pageOf(2, 1, "name")),
	))
	if err != nil {
		t.Fatalf("FindAll page 1: %v", err)
	}

	if total != 3 {
		t.Errorf("total = %d, want 3 (soft-deleted publisher excluded)", total)
	}
	if !reflect.DeepEqual(publisherUids(first), []string{"p-2", "p-3"}) {
		t.Errorf("page 0 = %v, want [p-2 p-3]", publisherUids(first))
	}
	if !reflect.DeepEqual(publisherUids(second), []string{"p-1"}) {
		t.Errorf("page 1 = %v, want [p-1]", publisherUids(second))
	}
}

func TestPublisherRepositoryFindAll_noPagingGiven_capsPageAtDefaultSize(t *testing.T) {
	db := newFullCatalogDB(t)
	for i := 0; i < 25; i++ {
		seedPublisher(t, db, uuid.NewString(), "Publisher")
	}

	total, got, err := newPublisherRepo(db).FindAll(context.Background(), option.NewPublisherQueryOption())
	if err != nil {
		t.Fatalf("FindAll: %v", err)
	}

	if total != 25 || len(got) != 20 {
		t.Errorf("total = %d, len = %d, want 25 and the default page of 20", total, len(got))
	}
}

func TestPublisherRepositoryCreate_everyFieldGiven_persistsAndReadsBack(t *testing.T) {
	db := newFullCatalogDB(t)
	r := newPublisherRepo(db)

	id, err := r.Create(context.Background(), command.PublisherCreateCommand{
		Name:    "Allen & Unwin",
		Address: strPtr("London"),
		OlKey:   strPtr("OL1P"),
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	row := loadRow[entity.PublisherEntity](t, db, id)
	if parsed, err := uuid.Parse(row.Uid); err != nil || parsed.Version() != 7 {
		t.Fatalf("uid = %q, want a UUIDv7", row.Uid)
	}
	if row.Version != 1 || row.CreatedAt.IsZero() {
		t.Errorf("row = %+v", row)
	}

	want := &repo.PublisherResult{Id: id, Uid: row.Uid, Name: "Allen & Unwin", Address: strPtr("London"), OlKey: strPtr("OL1P")}
	byUid, err := r.FindByUid(context.Background(), row.Uid)
	if err != nil {
		t.Fatalf("FindByUid: %v", err)
	}
	byId, err := r.FindOne(context.Background(), id)
	if err != nil {
		t.Fatalf("FindOne: %v", err)
	}
	if !reflect.DeepEqual(byUid, want) || !reflect.DeepEqual(byId, want) {
		t.Errorf("byUid = %+v, byId = %+v, want %+v", byUid, byId, want)
	}
}

func TestPublisherRepository_unknownOrSoftDeletedPublisher_failsWithRecordNotFound(t *testing.T) {
	db := newFullCatalogDB(t)
	gone := seedPublisher(t, db, "p-1", "Defunct Press")
	softDelete(t, db, "publisher", gone.Id)
	r := newPublisherRepo(db)

	for _, id := range []int64{gone.Id, 999} {
		if _, err := r.FindOne(context.Background(), id); !errors.Is(err, gorm.ErrRecordNotFound) {
			t.Errorf("FindOne(%d) err = %v, want gorm.ErrRecordNotFound", id, err)
		}
	}
	for _, uid := range []string{"p-1", "p-missing"} {
		if _, err := r.FindByUid(context.Background(), uid); !errors.Is(err, gorm.ErrRecordNotFound) {
			t.Errorf("FindByUid(%q) err = %v, want gorm.ErrRecordNotFound", uid, err)
		}
	}
}

func TestPublisherRepositoryUpdate_addressAndOlKeyGiven_persistsThemAndKeepsName(t *testing.T) {
	db := newFullCatalogDB(t)
	publisher := seedPublisher(t, db, "p-1", "Allen & Unwin")

	err := newPublisherRepo(db).Update(context.Background(), publisher.Id, command.PublisherUpdateCommand{
		Address: strPtr("Sydney"),
		OlKey:   strPtr("OL2P"),
	})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}

	got := loadRow[entity.PublisherEntity](t, db, publisher.Id)
	if got.Name != "Allen & Unwin" || got.Address == nil || *got.Address != "Sydney" || got.OlKey == nil || *got.OlKey != "OL2P" || got.UpdatedAt == nil {
		t.Errorf("row = %+v", got)
	}
}

func TestPublisherRepositoryUpdate_noFieldGiven_leavesRowUntouched(t *testing.T) {
	db, rec := newRecordedDB(t, publisherTableDDL)
	publisher := seedRow(t, db, &entity.PublisherEntity{Uid: "p-1", Name: "before", CreatedAt: seededAt, Version: 1})

	if err := newPublisherRepo(db).Update(context.Background(), publisher.Id, command.PublisherUpdateCommand{}); err != nil {
		t.Fatalf("Update: %v", err)
	}

	if updates := rec.updatesOf("publisher"); len(updates) != 0 {
		t.Errorf("UPDATE statements = %q, want none", updates)
	}
}

func TestPublisherRepository_tablesUnavailable_propagatesQueryError(t *testing.T) {
	db := newCatalogDB(t)
	r := newPublisherRepo(db)

	if _, _, err := r.FindAll(context.Background(), option.NewPublisherQueryOption()); err == nil {
		t.Error("FindAll err = nil, want the query failure")
	}
	if id, err := r.Create(context.Background(), command.PublisherCreateCommand{Name: "x"}); err == nil || id != 0 {
		t.Errorf("Create id = %d, err = %v, want 0 and the insert failure", id, err)
	}
}

func TestPublisherRepositoryFindAll_nameFilterGiven_matchesSubstringIgnoringCase(t *testing.T) {
	db := newFullCatalogDB(t)
	seedPublisher(t, db, "p-1", "Allen & Unwin")
	seedPublisher(t, db, "p-2", "HarperCollins")

	total, got, err := newPublisherRepo(db).FindAll(context.Background(), option.NewPublisherQueryOption(
		option.PublisherQueryOption{}.WithName(strPtr("perCOLL")),
	))
	if err != nil {
		t.Fatalf("FindAll: %v", err)
	}

	if total != 1 || !reflect.DeepEqual(publisherUids(got), []string{"p-2"}) {
		t.Errorf("total = %d, uids = %v, want 1 and [p-2]", total, publisherUids(got))
	}
}
