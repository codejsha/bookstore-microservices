package pgsql

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/codejsha/bookstore-microservices/catalog/generated/infrastructure/port/dao"
	"github.com/codejsha/bookstore-microservices/catalog/generated/infrastructure/port/entity"
	"github.com/codejsha/bookstore-microservices/catalog/internal/application/port/repo"
	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/model/command"
	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/model/option"
)

func newEditionRepo(db *gorm.DB) *editionRepository {
	return &editionRepository{q: dao.Use(db)}
}

func editionUids(results []*repo.EditionResult) []string {
	uids := make([]string, len(results))
	for i, r := range results {
		uids[i] = r.Uid
	}
	return uids
}

func TestEditionRepositoryFindAll_noFilter_joinsWorkTitleAndOptionalPublisher(t *testing.T) {
	db := newFullCatalogDB(t)
	hobbit := seedWork(t, db, "w-1", "The Hobbit")
	allen := seedPublisher(t, db, "p-1", "Allen & Unwin")
	seedEdition(t, db, "e-1", "First edition", hobbit, allen)
	seedEdition(t, db, "e-2", "Pocket edition", hobbit, nil)

	total, got, err := newEditionRepo(db).FindAll(context.Background(), option.NewEditionQueryOption())
	if err != nil {
		t.Fatalf("FindAll: %v", err)
	}

	if total != 2 || !reflect.DeepEqual(editionUids(got), []string{"e-1", "e-2"}) {
		t.Fatalf("total = %d, uids = %v, want 2 and [e-1 e-2]", total, editionUids(got))
	}
	first := got[0]
	if first.WorkUid != "w-1" || first.WorkTitle != "The Hobbit" || !first.CreatedAt.Equal(seededAt) {
		t.Errorf("first = %+v", first)
	}
	if first.PublisherUid == nil || *first.PublisherUid != "p-1" || first.PublisherName == nil || *first.PublisherName != "Allen & Unwin" {
		t.Errorf("first publisher = %v/%v, want p-1/Allen & Unwin", first.PublisherUid, first.PublisherName)
	}
	if got[1].PublisherUid != nil || got[1].PublisherName != nil {
		t.Errorf("second publisher = %v/%v, want nil/nil", got[1].PublisherUid, got[1].PublisherName)
	}
}

func TestEditionRepositoryFindAll_workUidFilterGiven_keepsOnlyThatWorksEditions(t *testing.T) {
	db := newFullCatalogDB(t)
	hobbit := seedWork(t, db, "w-1", "The Hobbit")
	dune := seedWork(t, db, "w-2", "Dune")
	seedEdition(t, db, "e-1", "Hobbit hardcover", hobbit, nil)
	seedEdition(t, db, "e-2", "Dune hardcover", dune, nil)
	seedEdition(t, db, "e-3", "Hobbit paperback", hobbit, nil)

	total, got, err := newEditionRepo(db).FindAll(context.Background(), option.NewEditionQueryOption(
		option.EditionQueryOption{}.WithWorkUid(strPtr("w-1")),
	))
	if err != nil {
		t.Fatalf("FindAll: %v", err)
	}

	if total != 2 || !reflect.DeepEqual(editionUids(got), []string{"e-1", "e-3"}) {
		t.Errorf("total = %d, uids = %v, want 2 and [e-1 e-3]", total, editionUids(got))
	}
}

func TestEditionRepositoryFindAll_publisherUidFilterGiven_keepsOnlyThatPublishersEditions(t *testing.T) {
	db := newFullCatalogDB(t)
	hobbit := seedWork(t, db, "w-1", "The Hobbit")
	allen := seedPublisher(t, db, "p-1", "Allen & Unwin")
	harper := seedPublisher(t, db, "p-2", "HarperCollins")
	seedEdition(t, db, "e-1", "First edition", hobbit, allen)
	seedEdition(t, db, "e-2", "Reissue", hobbit, harper)
	seedEdition(t, db, "e-3", "Unattributed", hobbit, nil)

	total, got, err := newEditionRepo(db).FindAll(context.Background(), option.NewEditionQueryOption(
		option.EditionQueryOption{}.WithPublisherUid(strPtr("p-2")),
	))
	if err != nil {
		t.Fatalf("FindAll: %v", err)
	}

	if total != 1 || !reflect.DeepEqual(editionUids(got), []string{"e-2"}) {
		t.Errorf("total = %d, uids = %v, want 1 and [e-2]", total, editionUids(got))
	}
}

func TestEditionRepositoryFindAll_olKeyFilterGiven_matchesExactKeyOnly(t *testing.T) {
	db := newFullCatalogDB(t)
	hobbit := seedWork(t, db, "w-1", "The Hobbit")
	first := seedEdition(t, db, "e-1", "First edition", hobbit, nil)
	second := seedEdition(t, db, "e-2", "Second edition", hobbit, nil)
	setColumn(t, db, "edition", first.Id, "ol_key", "OL1M")
	setColumn(t, db, "edition", second.Id, "ol_key", "OL11M")

	total, got, err := newEditionRepo(db).FindAll(context.Background(), option.NewEditionQueryOption(
		option.EditionQueryOption{}.WithOlKey(strPtr("OL11M")),
	))
	if err != nil {
		t.Fatalf("FindAll: %v", err)
	}

	if total != 1 || !reflect.DeepEqual(editionUids(got), []string{"e-2"}) {
		t.Errorf("total = %d, uids = %v, want 1 and [e-2]", total, editionUids(got))
	}
}

func TestEditionRepositoryFindAll_blankFilters_ignoresThem(t *testing.T) {
	db := newFullCatalogDB(t)
	hobbit := seedWork(t, db, "w-1", "The Hobbit")
	seedEdition(t, db, "e-1", "First edition", hobbit, nil)

	total, got, err := newEditionRepo(db).FindAll(context.Background(), option.NewEditionQueryOption(
		option.EditionQueryOption{}.WithTitle(strPtr("")),
		option.EditionQueryOption{}.WithWorkUid(strPtr("")),
		option.EditionQueryOption{}.WithPublisherUid(strPtr("")),
		option.EditionQueryOption{}.WithIsbn(strPtr("")),
		option.EditionQueryOption{}.WithOlKey(strPtr("")),
		option.EditionQueryOption{}.WithLanguage(strPtr("")),
	))
	if err != nil {
		t.Fatalf("FindAll: %v", err)
	}

	if total != 1 || len(got) != 1 {
		t.Errorf("total = %d, len = %d, want 1 and 1", total, len(got))
	}
}

func TestEditionRepositoryFindAll_titleSortWithPaging_pagesAfterSortingAndKeepsTotal(t *testing.T) {
	db := newFullCatalogDB(t)
	hobbit := seedWork(t, db, "w-1", "The Hobbit")
	seedEdition(t, db, "e-1", "Charlie", hobbit, nil)
	seedEdition(t, db, "e-2", "Alpha", hobbit, nil)
	seedEdition(t, db, "e-3", "Bravo", hobbit, nil)
	r := newEditionRepo(db)

	total, first, err := r.FindAll(context.Background(), option.NewEditionQueryOption(
		option.EditionQueryOption{}.WithPage(pageOf(2, 0, "title,asc")),
	))
	if err != nil {
		t.Fatalf("FindAll page 0: %v", err)
	}
	_, second, err := r.FindAll(context.Background(), option.NewEditionQueryOption(
		option.EditionQueryOption{}.WithPage(pageOf(2, 1, "title,asc")),
	))
	if err != nil {
		t.Fatalf("FindAll page 1: %v", err)
	}

	if total != 3 {
		t.Errorf("total = %d, want 3", total)
	}
	if !reflect.DeepEqual(editionUids(first), []string{"e-2", "e-3"}) {
		t.Errorf("page 0 = %v, want [e-2 e-3]", editionUids(first))
	}
	if !reflect.DeepEqual(editionUids(second), []string{"e-1"}) {
		t.Errorf("page 1 = %v, want [e-1]", editionUids(second))
	}
}

func TestEditionRepositoryFindAll_softDeletedEditionOrWork_excludesThem(t *testing.T) {
	db := newFullCatalogDB(t)
	hobbit := seedWork(t, db, "w-1", "The Hobbit")
	withdrawn := seedWork(t, db, "w-2", "Withdrawn")
	seedEdition(t, db, "e-1", "Kept", hobbit, nil)
	gone := seedEdition(t, db, "e-2", "Deleted edition", hobbit, nil)
	seedEdition(t, db, "e-3", "Edition of deleted work", withdrawn, nil)
	softDelete(t, db, "edition", gone.Id)
	softDelete(t, db, "work", withdrawn.Id)

	total, got, err := newEditionRepo(db).FindAll(context.Background(), option.NewEditionQueryOption())
	if err != nil {
		t.Fatalf("FindAll: %v", err)
	}

	if total != 1 || !reflect.DeepEqual(editionUids(got), []string{"e-1"}) {
		t.Errorf("total = %d, uids = %v, want 1 and [e-1]", total, editionUids(got))
	}
}

func TestEditionRepositoryFindAll_softDeletedPublisher_keepsEditionWithoutPublisher(t *testing.T) {
	db := newFullCatalogDB(t)
	hobbit := seedWork(t, db, "w-1", "The Hobbit")
	defunct := seedPublisher(t, db, "p-1", "Defunct Press")
	seedEdition(t, db, "e-1", "First edition", hobbit, defunct)
	softDelete(t, db, "publisher", defunct.Id)

	total, got, err := newEditionRepo(db).FindAll(context.Background(), option.NewEditionQueryOption())
	if err != nil {
		t.Fatalf("FindAll: %v", err)
	}

	if total != 1 || len(got) != 1 {
		t.Fatalf("total = %d, len = %d, want 1 and 1", total, len(got))
	}
	if got[0].PublisherUid != nil || got[0].PublisherName != nil {
		t.Errorf("publisher = %v/%v, want nil/nil", got[0].PublisherUid, got[0].PublisherName)
	}
}

func TestEditionRepositoryFindByUid_existingUid_mapsEveryColumn(t *testing.T) {
	db := newFullCatalogDB(t)
	hobbit := seedWork(t, db, "w-1", "The Hobbit")
	allen := seedPublisher(t, db, "p-1", "Allen & Unwin")
	pages := int32(310)
	r := newEditionRepo(db)
	id, err := r.Create(context.Background(), command.EditionCreateCommand{
		Title:          "First edition",
		Isbn10:         strPtr("0261102214"),
		Isbn13:         strPtr("9780261102217"),
		NumberOfPages:  &pages,
		PublishDate:    strPtr("1937-09-21"),
		CoverUids:      []string{"c-1"},
		Languages:      []string{"eng", "kor"},
		PhysicalFormat: strPtr("Hardcover"),
		Description:    strPtr("first printing"),
		OlKey:          strPtr("OL1M"),
	}, repo.EditionWorkRef{Id: hobbit.Id, Uid: hobbit.Uid}, &repo.EditionPublisherRef{Id: allen.Id, Uid: allen.Uid})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	row := loadRow[entity.EditionEntity](t, db, id)
	if parsed, err := uuid.Parse(row.Uid); err != nil || parsed.Version() != 7 {
		t.Fatalf("uid = %q, want a UUIDv7", row.Uid)
	}

	got, err := r.FindByUid(context.Background(), row.Uid)
	if err != nil {
		t.Fatalf("FindByUid: %v", err)
	}

	want := &repo.EditionResult{
		Id:             id,
		Uid:            row.Uid,
		Title:          "First edition",
		Isbn10:         strPtr("0261102214"),
		Isbn13:         strPtr("9780261102217"),
		NumberOfPages:  &pages,
		PublishDate:    strPtr("1937-09-21"),
		CoverUids:      []string{"c-1"},
		Languages:      []string{"eng", "kor"},
		PhysicalFormat: strPtr("Hardcover"),
		Description:    strPtr("first printing"),
		WorkUid:        "w-1",
		WorkTitle:      "The Hobbit",
		PublisherUid:   strPtr("p-1"),
		PublisherName:  strPtr("Allen & Unwin"),
		OlKey:          strPtr("OL1M"),
		CreatedAt:      got.CreatedAt,
		UpdatedAt:      got.UpdatedAt,
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("edition = %+v, want %+v", got, want)
	}
	if got.CreatedAt.IsZero() {
		t.Error("created_at not populated")
	}
}

func TestEditionRepositoryFindByUid_unknownOrSoftDeletedUid_failsWithRecordNotFound(t *testing.T) {
	db := newFullCatalogDB(t)
	hobbit := seedWork(t, db, "w-1", "The Hobbit")
	gone := seedEdition(t, db, "e-1", "Deleted edition", hobbit, nil)
	softDelete(t, db, "edition", gone.Id)
	r := newEditionRepo(db)

	for _, uid := range []string{"e-1", "e-missing"} {
		got, err := r.FindByUid(context.Background(), uid)
		if !errors.Is(err, gorm.ErrRecordNotFound) || got != nil {
			t.Errorf("FindByUid(%q) = %v, %v, want nil and gorm.ErrRecordNotFound", uid, got, err)
		}
	}
}

func TestEditionRepositoryFindOne_existingId_loadsThatEdition(t *testing.T) {
	db := newFullCatalogDB(t)
	hobbit := seedWork(t, db, "w-1", "The Hobbit")
	seedEdition(t, db, "e-1", "First edition", hobbit, nil)
	second := seedEdition(t, db, "e-2", "Second edition", hobbit, nil)

	got, err := newEditionRepo(db).FindOne(context.Background(), second.Id)
	if err != nil {
		t.Fatalf("FindOne: %v", err)
	}

	if got.Id != second.Id || got.Uid != "e-2" || got.Title != "Second edition" || got.WorkTitle != "The Hobbit" {
		t.Errorf("edition = %+v", got)
	}
}

func TestEditionRepositoryFindOne_workSoftDeleted_failsWithRecordNotFound(t *testing.T) {
	db := newFullCatalogDB(t)
	withdrawn := seedWork(t, db, "w-1", "Withdrawn")
	edition := seedEdition(t, db, "e-1", "Edition of deleted work", withdrawn, nil)
	softDelete(t, db, "work", withdrawn.Id)

	_, err := newEditionRepo(db).FindOne(context.Background(), edition.Id)

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Errorf("err = %v, want gorm.ErrRecordNotFound", err)
	}
}

func TestEditionRepositoryUpdate_everyFieldGiven_persistsEachField(t *testing.T) {
	db := newFullCatalogDB(t)
	hobbit := seedWork(t, db, "w-1", "The Hobbit")
	edition := seedEdition(t, db, "e-1", "before", hobbit, nil)
	pages := int32(288)

	err := newEditionRepo(db).Update(context.Background(), edition.Id, command.EditionUpdateCommand{
		Isbn10:         strPtr("0261102214"),
		Isbn13:         strPtr("9780261102217"),
		NumberOfPages:  &pages,
		PublishDate:    strPtr("1951"),
		CoverUids:      []string{"c-2"},
		Languages:      []string{"eng"},
		PhysicalFormat: strPtr("Paperback"),
		Description:    strPtr("second edition"),
		OlKey:          strPtr("OL2M"),
	}, nil, nil)
	if err != nil {
		t.Fatalf("Update: %v", err)
	}

	got := loadRow[entity.EditionEntity](t, db, edition.Id)
	if got.Title != "before" || got.WorkUid != "w-1" {
		t.Errorf("untouched columns changed: %+v", got)
	}
	if got.Isbn10 == nil || *got.Isbn10 != "0261102214" || got.Isbn13 == nil || *got.Isbn13 != "9780261102217" ||
		got.NumberOfPage == nil || *got.NumberOfPage != 288 || got.PublishDate == nil || *got.PublishDate != "1951" {
		t.Errorf("identifiers = %+v", got)
	}
	if got.CoverUid == nil || *got.CoverUid != `["c-2"]` || got.Language == nil || *got.Language != `["eng"]` {
		t.Errorf("cover_uids = %v, languages = %v", got.CoverUid, got.Language)
	}
	if got.PhysicalFormat == nil || *got.PhysicalFormat != "Paperback" || got.Description == nil || *got.Description != "second edition" ||
		got.OlKey == nil || *got.OlKey != "OL2M" || got.UpdatedAt == nil {
		t.Errorf("descriptive columns = %+v", got)
	}
}

func TestEditionRepositoryUpdate_noFieldGiven_leavesRowUntouched(t *testing.T) {
	db, rec := newRecordedDB(t, editionTableDDL)
	edition := seedRow(t, db, &entity.EditionEntity{Uid: "e-1", Title: "before", WorkId: 1, WorkUid: "w-1", CreatedAt: seededAt, Version: 1})

	if err := newEditionRepo(db).Update(context.Background(), edition.Id, command.EditionUpdateCommand{}, nil, nil); err != nil {
		t.Fatalf("Update: %v", err)
	}

	if updates := rec.updatesOf("edition"); len(updates) != 0 {
		t.Errorf("UPDATE statements = %q, want none", updates)
	}
}

func TestEditionRepository_tablesUnavailable_propagatesQueryError(t *testing.T) {
	db := newCatalogDB(t)
	r := newEditionRepo(db)

	if _, _, err := r.FindAll(context.Background(), option.NewEditionQueryOption()); err == nil {
		t.Error("FindAll err = nil, want the query failure")
	}
	if _, err := r.FindOne(context.Background(), 1); err == nil || errors.Is(err, gorm.ErrRecordNotFound) {
		t.Errorf("FindOne err = %v, want the query failure, not a not-found", err)
	}
	id, err := r.Create(context.Background(), command.EditionCreateCommand{Title: "x"}, repo.EditionWorkRef{Id: 1, Uid: "w-1"}, nil)
	if err == nil || id != 0 {
		t.Errorf("Create id = %d, err = %v, want 0 and the insert failure", id, err)
	}
}

func TestEditionRepositoryFindAll_titleFilterGiven_matchesSubstringIgnoringCase(t *testing.T) {
	db := newFullCatalogDB(t)
	hobbit := seedWork(t, db, "w-1", "The Hobbit")
	seedEdition(t, db, "e-1", "Illustrated edition", hobbit, nil)
	seedEdition(t, db, "e-2", "Pocket edition", hobbit, nil)

	total, got, err := newEditionRepo(db).FindAll(context.Background(), option.NewEditionQueryOption(
		option.EditionQueryOption{}.WithTitle(strPtr("POCKET")),
	))
	if err != nil {
		t.Fatalf("FindAll: %v", err)
	}

	if total != 1 || !reflect.DeepEqual(editionUids(got), []string{"e-2"}) {
		t.Errorf("total = %d, uids = %v, want 1 and [e-2]", total, editionUids(got))
	}
}

func TestEditionRepositoryFindAll_isbnFilterGiven_matchesEitherIsbnColumn(t *testing.T) {
	db := newFullCatalogDB(t)
	hobbit := seedWork(t, db, "w-1", "The Hobbit")
	first := seedEdition(t, db, "e-1", "First edition", hobbit, nil)
	second := seedEdition(t, db, "e-2", "Second edition", hobbit, nil)
	third := seedEdition(t, db, "e-3", "Third edition", hobbit, nil)
	setColumn(t, db, "edition", first.Id, "isbn10", "026110221X")
	setColumn(t, db, "edition", second.Id, "isbn13", "9780261102217")
	setColumn(t, db, "edition", third.Id, "isbn13", "9780345339683")

	total, got, err := newEditionRepo(db).FindAll(context.Background(), option.NewEditionQueryOption(
		option.EditionQueryOption{}.WithIsbn(strPtr("26110221")),
	))
	if err != nil {
		t.Fatalf("FindAll: %v", err)
	}

	if total != 2 || !reflect.DeepEqual(editionUids(got), []string{"e-1", "e-2"}) {
		t.Errorf("total = %d, uids = %v, want 2 and [e-1 e-2]", total, editionUids(got))
	}
}

func TestEditionRepositoryFindAll_languageFilterGiven_sendsJsonbContainmentPredicate(t *testing.T) {
	db, rec := newRecordedDB(t, catalogSchemaDDL...)

	_, _, _ = newEditionRepo(db).FindAll(context.Background(), option.NewEditionQueryOption(
		option.EditionQueryOption{}.WithLanguage(strPtr("eng")),
	))

	want := `edition.languages @> "[""eng""]"::jsonb`
	for _, statement := range rec.statements {
		if strings.Contains(statement, want) {
			return
		}
	}
	t.Errorf("statements = %q, want one containing %s", rec.statements, want)
}
