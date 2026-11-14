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

func newAuthorRepo(db *gorm.DB) *authorRepository {
	return &authorRepository{q: dao.Use(db)}
}

func authorUids(results []*repo.AuthorResult) []string {
	uids := make([]string, len(results))
	for i, r := range results {
		uids[i] = r.Uid
	}
	return uids
}

func TestAuthorRepositoryFindAll_workUidFilterGiven_keepsOnlyThatWorksAuthors(t *testing.T) {
	db := newFullCatalogDB(t)
	omens := seedWork(t, db, "w-1", "Good Omens")
	dune := seedWork(t, db, "w-2", "Dune")
	pratchett := seedAuthor(t, db, "a-1", "Pratchett")
	herbert := seedAuthor(t, db, "a-2", "Herbert")
	gaiman := seedAuthor(t, db, "a-3", "Gaiman")
	linkAuthor(t, db, omens.Id, pratchett.Id)
	linkAuthor(t, db, omens.Id, gaiman.Id)
	linkAuthor(t, db, dune.Id, herbert.Id)

	total, got, err := newAuthorRepo(db).FindAll(context.Background(), option.NewAuthorQueryOption(
		option.AuthorQueryOption{}.WithWorkUid(strPtr("w-1")),
	))
	if err != nil {
		t.Fatalf("FindAll: %v", err)
	}

	if total != 2 || !reflect.DeepEqual(authorUids(got), []string{"a-1", "a-3"}) {
		t.Errorf("total = %d, uids = %v, want 2 and [a-1 a-3]", total, authorUids(got))
	}
}

func TestAuthorRepositoryFindAll_workUidFilterOnSoftDeletedWork_matchesNothing(t *testing.T) {
	db := newFullCatalogDB(t)
	withdrawn := seedWork(t, db, "w-1", "Withdrawn")
	pratchett := seedAuthor(t, db, "a-1", "Pratchett")
	linkAuthor(t, db, withdrawn.Id, pratchett.Id)
	softDelete(t, db, "work", withdrawn.Id)

	total, got, err := newAuthorRepo(db).FindAll(context.Background(), option.NewAuthorQueryOption(
		option.AuthorQueryOption{}.WithWorkUid(strPtr("w-1")),
	))
	if err != nil {
		t.Fatalf("FindAll: %v", err)
	}

	if total != 0 || len(got) != 0 {
		t.Errorf("total = %d, uids = %v, want 0 and none", total, authorUids(got))
	}
}

func TestAuthorRepositoryFindAll_olKeyFilterGiven_matchesExactKeyOnly(t *testing.T) {
	db := newFullCatalogDB(t)
	pratchett := seedAuthor(t, db, "a-1", "Pratchett")
	gaiman := seedAuthor(t, db, "a-2", "Gaiman")
	setColumn(t, db, "author", pratchett.Id, "ol_key", "OL1A")
	setColumn(t, db, "author", gaiman.Id, "ol_key", "OL10A")

	total, got, err := newAuthorRepo(db).FindAll(context.Background(), option.NewAuthorQueryOption(
		option.AuthorQueryOption{}.WithOlKey(strPtr("OL1A")),
	))
	if err != nil {
		t.Fatalf("FindAll: %v", err)
	}

	if total != 1 || !reflect.DeepEqual(authorUids(got), []string{"a-1"}) {
		t.Errorf("total = %d, uids = %v, want 1 and [a-1]", total, authorUids(got))
	}
}

func TestAuthorRepositoryFindAll_nameDescSortWithPaging_pagesAfterSortingAndKeepsTotal(t *testing.T) {
	db := newFullCatalogDB(t)
	seedAuthor(t, db, "a-1", "Herbert")
	seedAuthor(t, db, "a-2", "Tolkien")
	seedAuthor(t, db, "a-3", "Asimov")
	gone := seedAuthor(t, db, "a-4", "Zelazny")
	softDelete(t, db, "author", gone.Id)
	r := newAuthorRepo(db)

	total, first, err := r.FindAll(context.Background(), option.NewAuthorQueryOption(
		option.AuthorQueryOption{}.WithName(strPtr("")),
		option.AuthorQueryOption{}.WithPage(pageOf(2, 0, "name,desc")),
	))
	if err != nil {
		t.Fatalf("FindAll page 0: %v", err)
	}
	_, second, err := r.FindAll(context.Background(), option.NewAuthorQueryOption(
		option.AuthorQueryOption{}.WithPage(pageOf(2, 1, "name,desc")),
	))
	if err != nil {
		t.Fatalf("FindAll page 1: %v", err)
	}

	if total != 3 {
		t.Errorf("total = %d, want 3 (soft-deleted author excluded)", total)
	}
	if !reflect.DeepEqual(authorUids(first), []string{"a-2", "a-1"}) {
		t.Errorf("page 0 = %v, want [a-2 a-1]", authorUids(first))
	}
	if !reflect.DeepEqual(authorUids(second), []string{"a-3"}) {
		t.Errorf("page 1 = %v, want [a-3]", authorUids(second))
	}
}

func TestAuthorRepositoryCreate_everyFieldGiven_persistsAndReadsBackByUid(t *testing.T) {
	db := newFullCatalogDB(t)
	r := newAuthorRepo(db)

	id, err := r.Create(context.Background(), command.AuthorCreateCommand{
		Name:           "Isaac Asimov",
		Bio:            strPtr("biochemist"),
		BirthDate:      strPtr("1920-01-02"),
		DeathDate:      strPtr("1992-04-06"),
		PhotoUids:      []string{"ph-1"},
		AlternateNames: []string{"Paul French", "George E. Dale"},
		OlKey:          strPtr("OL34221A"),
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	row := loadRow[entity.AuthorEntity](t, db, id)
	if parsed, err := uuid.Parse(row.Uid); err != nil || parsed.Version() != 7 {
		t.Fatalf("uid = %q, want a UUIDv7", row.Uid)
	}
	if row.Version != 1 || row.CreatedAt.IsZero() {
		t.Errorf("row = %+v", row)
	}

	got, err := r.FindByUid(context.Background(), row.Uid)
	if err != nil {
		t.Fatalf("FindByUid: %v", err)
	}

	want := &repo.AuthorResult{
		Id:             id,
		Uid:            row.Uid,
		Name:           "Isaac Asimov",
		Bio:            strPtr("biochemist"),
		BirthDate:      strPtr("1920-01-02"),
		DeathDate:      strPtr("1992-04-06"),
		PhotoUids:      []string{"ph-1"},
		AlternateNames: []string{"Paul French", "George E. Dale"},
		OlKey:          strPtr("OL34221A"),
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("author = %+v, want %+v", got, want)
	}
}

func TestAuthorRepositoryFindOne_existingId_loadsThatAuthor(t *testing.T) {
	db := newFullCatalogDB(t)
	seedAuthor(t, db, "a-1", "Herbert")
	tolkien := seedAuthor(t, db, "a-2", "Tolkien")

	got, err := newAuthorRepo(db).FindOne(context.Background(), tolkien.Id)
	if err != nil {
		t.Fatalf("FindOne: %v", err)
	}

	if got.Id != tolkien.Id || got.Uid != "a-2" || got.Name != "Tolkien" {
		t.Errorf("author = %+v", got)
	}
}

func TestAuthorRepository_unknownOrSoftDeletedAuthor_failsWithRecordNotFound(t *testing.T) {
	db := newFullCatalogDB(t)
	gone := seedAuthor(t, db, "a-1", "Tolkien")
	softDelete(t, db, "author", gone.Id)
	r := newAuthorRepo(db)

	for _, id := range []int64{gone.Id, 999} {
		if _, err := r.FindOne(context.Background(), id); !errors.Is(err, gorm.ErrRecordNotFound) {
			t.Errorf("FindOne(%d) err = %v, want gorm.ErrRecordNotFound", id, err)
		}
	}
	for _, uid := range []string{"a-1", "a-missing"} {
		if _, err := r.FindByUid(context.Background(), uid); !errors.Is(err, gorm.ErrRecordNotFound) {
			t.Errorf("FindByUid(%q) err = %v, want gorm.ErrRecordNotFound", uid, err)
		}
	}
}

func TestAuthorRepositoryUpdate_everyFieldGiven_persistsEachField(t *testing.T) {
	db := newFullCatalogDB(t)
	author := seedAuthor(t, db, "a-1", "before")
	other := seedAuthor(t, db, "a-2", "bystander")

	err := newAuthorRepo(db).Update(context.Background(), author.Id, command.AuthorUpdateCommand{
		Bio:            strPtr("novelist"),
		BirthDate:      strPtr("1892-01-03"),
		DeathDate:      strPtr("1973-09-02"),
		PhotoUids:      []string{"ph-9"},
		AlternateNames: []string{"J.R.R. Tolkien"},
		OlKey:          strPtr("OL26320A"),
	})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}

	got := loadRow[entity.AuthorEntity](t, db, author.Id)
	if got.Name != "before" || got.Bio == nil || *got.Bio != "novelist" ||
		got.BirthDate == nil || *got.BirthDate != "1892-01-03" || got.DeathDate == nil || *got.DeathDate != "1973-09-02" {
		t.Errorf("row = %+v", got)
	}
	if got.PhotoUid == nil || *got.PhotoUid != `["ph-9"]` || got.AlternateName == nil || *got.AlternateName != `["J.R.R. Tolkien"]` {
		t.Errorf("photo_uids = %v, alternate_names = %v", got.PhotoUid, got.AlternateName)
	}
	if got.OlKey == nil || *got.OlKey != "OL26320A" || got.UpdatedAt == nil {
		t.Errorf("ol_key = %v, updated_at = %v", got.OlKey, got.UpdatedAt)
	}
	if bystander := loadRow[entity.AuthorEntity](t, db, other.Id); bystander.Bio != nil || bystander.OlKey != nil {
		t.Errorf("other author changed: %+v", bystander)
	}
}

func TestAuthorRepository_tablesUnavailable_propagatesQueryError(t *testing.T) {
	db := newCatalogDB(t)
	r := newAuthorRepo(db)

	if _, _, err := r.FindAll(context.Background(), option.NewAuthorQueryOption()); err == nil {
		t.Error("FindAll err = nil, want the query failure")
	}
	if id, err := r.Create(context.Background(), command.AuthorCreateCommand{Name: "x"}); err == nil || id != 0 {
		t.Errorf("Create id = %d, err = %v, want 0 and the insert failure", id, err)
	}
}

func TestAuthorRepositoryFindAll_nameFilterGiven_matchesSubstringIgnoringCase(t *testing.T) {
	db := newFullCatalogDB(t)
	seedAuthor(t, db, "a-1", "Pratchett")
	seedAuthor(t, db, "a-2", "Gaiman")

	total, got, err := newAuthorRepo(db).FindAll(context.Background(), option.NewAuthorQueryOption(
		option.AuthorQueryOption{}.WithName(strPtr("RATCH")),
	))
	if err != nil {
		t.Fatalf("FindAll: %v", err)
	}

	if total != 1 || !reflect.DeepEqual(authorUids(got), []string{"a-1"}) {
		t.Errorf("total = %d, uids = %v, want 1 and [a-1]", total, authorUids(got))
	}
}
