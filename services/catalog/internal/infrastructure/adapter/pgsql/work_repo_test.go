package pgsql

import (
	"context"
	"errors"
	"reflect"
	"sort"
	"testing"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/codejsha/bookstore-microservices/catalog/generated/infrastructure/port/dao"
	"github.com/codejsha/bookstore-microservices/catalog/generated/infrastructure/port/entity"
	"github.com/codejsha/bookstore-microservices/catalog/internal/application/port/repo"
	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/model/command"
	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/model/option"
)

func newWorkRepo(db *gorm.DB) *workRepository {
	return &workRepository{q: dao.Use(db), db: db}
}

func workUids(results []*repo.WorkResult) []string {
	uids := make([]string, len(results))
	for i, r := range results {
		uids[i] = r.Uid
	}
	return uids
}

func sorted(values []string) []string {
	out := append([]string(nil), values...)
	sort.Strings(out)
	return out
}

func mappedAuthorIds(t *testing.T, db *gorm.DB, workID int64) []int64 {
	t.Helper()
	var ids []int64
	if err := db.Table("work_author_mapping").Where("work_id = ?", workID).Order("author_id").Pluck("author_id", &ids).Error; err != nil {
		t.Fatalf("load author mappings: %v", err)
	}
	return ids
}

func mappedSubjectNames(t *testing.T, db *gorm.DB, workID int64) []string {
	t.Helper()
	var names []string
	err := db.Table("work_subject_mapping").
		Joins("JOIN subject ON subject.id = work_subject_mapping.subject_id").
		Where("work_subject_mapping.work_id = ?", workID).
		Order("subject.name").
		Pluck("subject.name", &names).Error
	if err != nil {
		t.Fatalf("load subject mappings: %v", err)
	}
	return names
}

func TestWorkRepositoryFindAll_noFilter_hydratesAuthorsAndSubjectsInIdOrder(t *testing.T) {
	db := newFullCatalogDB(t)
	hobbit := seedWork(t, db, "w-1", "The Hobbit")
	dune := seedWork(t, db, "w-2", "Dune")
	tolkien := seedAuthor(t, db, "a-1", "Tolkien")
	herbert := seedAuthor(t, db, "a-2", "Herbert")
	fantasy := seedSubject(t, db, "s-1", "Fantasy")
	linkAuthor(t, db, hobbit.Id, tolkien.Id)
	linkAuthor(t, db, dune.Id, herbert.Id)
	linkSubject(t, db, hobbit.Id, fantasy.Id)

	total, got, err := newWorkRepo(db).FindAll(context.Background(), option.NewWorkQueryOption())
	if err != nil {
		t.Fatalf("FindAll: %v", err)
	}

	if total != 2 || !reflect.DeepEqual(workUids(got), []string{"w-1", "w-2"}) {
		t.Fatalf("total = %d, uids = %v, want 2 and [w-1 w-2]", total, workUids(got))
	}
	if !reflect.DeepEqual(got[0].AuthorUids, []string{"a-1"}) || !reflect.DeepEqual(got[0].AuthorNames, []string{"Tolkien"}) {
		t.Errorf("hobbit authors = %v/%v", got[0].AuthorUids, got[0].AuthorNames)
	}
	if !reflect.DeepEqual(got[0].SubjectUids, []string{"s-1"}) || !reflect.DeepEqual(got[0].SubjectNames, []string{"Fantasy"}) {
		t.Errorf("hobbit subjects = %v/%v", got[0].SubjectUids, got[0].SubjectNames)
	}
	if !reflect.DeepEqual(got[1].AuthorUids, []string{"a-2"}) || len(got[1].SubjectUids) != 0 {
		t.Errorf("dune authors = %v, subjects = %v", got[1].AuthorUids, got[1].SubjectUids)
	}
}

func TestWorkRepositoryFindAll_authorFilterGiven_keepsOnlyThatAuthorsWorks(t *testing.T) {
	db := newFullCatalogDB(t)
	hobbit := seedWork(t, db, "w-1", "The Hobbit")
	dune := seedWork(t, db, "w-2", "Dune")
	silmarillion := seedWork(t, db, "w-3", "The Silmarillion")
	tolkien := seedAuthor(t, db, "a-1", "Tolkien")
	herbert := seedAuthor(t, db, "a-2", "Herbert")
	linkAuthor(t, db, hobbit.Id, tolkien.Id)
	linkAuthor(t, db, silmarillion.Id, tolkien.Id)
	linkAuthor(t, db, dune.Id, herbert.Id)

	total, got, err := newWorkRepo(db).FindAll(context.Background(), option.NewWorkQueryOption(
		option.WorkQueryOption{}.WithAuthorUid(strPtr("a-1")),
	))
	if err != nil {
		t.Fatalf("FindAll: %v", err)
	}

	if total != 2 || !reflect.DeepEqual(workUids(got), []string{"w-1", "w-3"}) {
		t.Errorf("total = %d, uids = %v, want 2 and [w-1 w-3]", total, workUids(got))
	}
}

func TestWorkRepositoryFindAll_authorFilterOnSoftDeletedAuthor_matchesNothing(t *testing.T) {
	db := newFullCatalogDB(t)
	hobbit := seedWork(t, db, "w-1", "The Hobbit")
	tolkien := seedAuthor(t, db, "a-1", "Tolkien")
	linkAuthor(t, db, hobbit.Id, tolkien.Id)
	softDelete(t, db, "author", tolkien.Id)

	total, got, err := newWorkRepo(db).FindAll(context.Background(), option.NewWorkQueryOption(
		option.WorkQueryOption{}.WithAuthorUid(strPtr("a-1")),
	))
	if err != nil {
		t.Fatalf("FindAll: %v", err)
	}

	if total != 0 || got == nil || len(got) != 0 {
		t.Errorf("total = %d, results = %v, want 0 and an empty slice", total, got)
	}
}

func TestWorkRepositoryFindAll_subjectFilterGiven_keepsOnlyThatSubjectsWorks(t *testing.T) {
	db := newFullCatalogDB(t)
	hobbit := seedWork(t, db, "w-1", "The Hobbit")
	dune := seedWork(t, db, "w-2", "Dune")
	fantasy := seedSubject(t, db, "s-1", "Fantasy")
	scifi := seedSubject(t, db, "s-2", "Science Fiction")
	linkSubject(t, db, hobbit.Id, fantasy.Id)
	linkSubject(t, db, dune.Id, scifi.Id)

	total, got, err := newWorkRepo(db).FindAll(context.Background(), option.NewWorkQueryOption(
		option.WorkQueryOption{}.WithSubjectUid(strPtr("s-2")),
	))
	if err != nil {
		t.Fatalf("FindAll: %v", err)
	}

	if total != 1 || !reflect.DeepEqual(workUids(got), []string{"w-2"}) {
		t.Errorf("total = %d, uids = %v, want 1 and [w-2]", total, workUids(got))
	}
}

func TestWorkRepositoryFindAll_olKeyFilterGiven_matchesExactKeyOnly(t *testing.T) {
	db := newFullCatalogDB(t)
	hobbit := seedWork(t, db, "w-1", "The Hobbit")
	dune := seedWork(t, db, "w-2", "Dune")
	setColumn(t, db, "work", hobbit.Id, "ol_key", "OL1W")
	setColumn(t, db, "work", dune.Id, "ol_key", "OL12W")

	total, got, err := newWorkRepo(db).FindAll(context.Background(), option.NewWorkQueryOption(
		option.WorkQueryOption{}.WithOlKey(strPtr("OL1W")),
	))
	if err != nil {
		t.Fatalf("FindAll: %v", err)
	}

	if total != 1 || !reflect.DeepEqual(workUids(got), []string{"w-1"}) {
		t.Errorf("total = %d, uids = %v, want 1 and [w-1]", total, workUids(got))
	}
	if got[0].OlKey == nil || *got[0].OlKey != "OL1W" {
		t.Errorf("OlKey = %v, want OL1W", got[0].OlKey)
	}
}

func TestWorkRepositoryFindAll_blankFilters_ignoresThem(t *testing.T) {
	db := newFullCatalogDB(t)
	seedWork(t, db, "w-1", "The Hobbit")
	seedWork(t, db, "w-2", "Dune")

	total, got, err := newWorkRepo(db).FindAll(context.Background(), option.NewWorkQueryOption(
		option.WorkQueryOption{}.WithTitle(strPtr("")),
		option.WorkQueryOption{}.WithAuthorUid(strPtr("")),
		option.WorkQueryOption{}.WithSubjectUid(strPtr("")),
		option.WorkQueryOption{}.WithOlKey(strPtr("")),
	))
	if err != nil {
		t.Fatalf("FindAll: %v", err)
	}

	if total != 2 || len(got) != 2 {
		t.Errorf("total = %d, len = %d, want 2 and 2", total, len(got))
	}
}

func TestWorkRepositoryFindAll_titleDescSortSecondPage_pagesAfterSortingAndKeepsTotal(t *testing.T) {
	db := newFullCatalogDB(t)
	seedWork(t, db, "w-1", "Bravo")
	seedWork(t, db, "w-2", "Delta")
	seedWork(t, db, "w-3", "Alpha")
	seedWork(t, db, "w-4", "Charlie")
	seedWork(t, db, "w-5", "Echo")
	r := newWorkRepo(db)

	total, first, err := r.FindAll(context.Background(), option.NewWorkQueryOption(
		option.WorkQueryOption{}.WithPage(pageOf(2, 0, "title,desc")),
	))
	if err != nil {
		t.Fatalf("FindAll page 0: %v", err)
	}
	_, second, err := r.FindAll(context.Background(), option.NewWorkQueryOption(
		option.WorkQueryOption{}.WithPage(pageOf(2, 1, "title,desc")),
	))
	if err != nil {
		t.Fatalf("FindAll page 1: %v", err)
	}
	_, last, err := r.FindAll(context.Background(), option.NewWorkQueryOption(
		option.WorkQueryOption{}.WithPage(pageOf(2, 2, "title,desc")),
	))
	if err != nil {
		t.Fatalf("FindAll page 2: %v", err)
	}

	if total != 5 {
		t.Errorf("total = %d, want 5", total)
	}
	if !reflect.DeepEqual(workUids(first), []string{"w-5", "w-2"}) {
		t.Errorf("page 0 = %v, want [w-5 w-2]", workUids(first))
	}
	if !reflect.DeepEqual(workUids(second), []string{"w-4", "w-1"}) {
		t.Errorf("page 1 = %v, want [w-4 w-1]", workUids(second))
	}
	if !reflect.DeepEqual(workUids(last), []string{"w-3"}) {
		t.Errorf("page 2 = %v, want [w-3]", workUids(last))
	}
}

func TestWorkRepositoryFindAll_equalSortKeys_breaksTiesById(t *testing.T) {
	db := newFullCatalogDB(t)
	seedWork(t, db, "w-1", "Same")
	seedWork(t, db, "w-2", "Same")
	seedWork(t, db, "w-3", "Same")

	_, got, err := newWorkRepo(db).FindAll(context.Background(), option.NewWorkQueryOption(
		option.WorkQueryOption{}.WithPage(pageOf(10, 0, "title,desc")),
	))
	if err != nil {
		t.Fatalf("FindAll: %v", err)
	}

	if !reflect.DeepEqual(workUids(got), []string{"w-1", "w-2", "w-3"}) {
		t.Errorf("uids = %v, want [w-1 w-2 w-3]", workUids(got))
	}
}

func TestWorkRepositoryFindAll_unknownSortField_fallsBackToIdOrder(t *testing.T) {
	db := newFullCatalogDB(t)
	seedWork(t, db, "w-1", "Zulu")
	seedWork(t, db, "w-2", "Alpha")

	_, got, err := newWorkRepo(db).FindAll(context.Background(), option.NewWorkQueryOption(
		option.WorkQueryOption{}.WithPage(pageOf(10, 0, "id;drop table work,desc")),
	))
	if err != nil {
		t.Fatalf("FindAll: %v", err)
	}

	if !reflect.DeepEqual(workUids(got), []string{"w-1", "w-2"}) {
		t.Errorf("uids = %v, want [w-1 w-2]", workUids(got))
	}
}

func TestWorkRepositoryFindAll_softDeletedWork_excludesItFromRowsAndTotal(t *testing.T) {
	db := newFullCatalogDB(t)
	seedWork(t, db, "w-1", "The Hobbit")
	gone := seedWork(t, db, "w-2", "Dune")
	softDelete(t, db, "work", gone.Id)

	total, got, err := newWorkRepo(db).FindAll(context.Background(), option.NewWorkQueryOption())
	if err != nil {
		t.Fatalf("FindAll: %v", err)
	}

	if total != 1 || !reflect.DeepEqual(workUids(got), []string{"w-1"}) {
		t.Errorf("total = %d, uids = %v, want 1 and [w-1]", total, workUids(got))
	}
}

func TestWorkRepositoryFindAll_noRows_yieldsEmptySliceAndZeroTotal(t *testing.T) {
	db := newFullCatalogDB(t)

	total, got, err := newWorkRepo(db).FindAll(context.Background(), option.NewWorkQueryOption())
	if err != nil {
		t.Fatalf("FindAll: %v", err)
	}

	if total != 0 || got == nil || len(got) != 0 {
		t.Errorf("total = %d, results = %v, want 0 and an empty slice", total, got)
	}
}

func TestWorkRepositoryFindByUid_softDeletedAuthorLinked_omitsThatAuthor(t *testing.T) {
	db := newFullCatalogDB(t)
	work := seedWork(t, db, "w-1", "Good Omens")
	pratchett := seedAuthor(t, db, "a-1", "Pratchett")
	gaiman := seedAuthor(t, db, "a-2", "Gaiman")
	humor := seedSubject(t, db, "s-1", "Humor")
	retired := seedSubject(t, db, "s-2", "Retired")
	linkAuthor(t, db, work.Id, pratchett.Id)
	linkAuthor(t, db, work.Id, gaiman.Id)
	linkSubject(t, db, work.Id, humor.Id)
	linkSubject(t, db, work.Id, retired.Id)
	softDelete(t, db, "author", gaiman.Id)
	softDelete(t, db, "subject", retired.Id)

	got, err := newWorkRepo(db).FindByUid(context.Background(), "w-1")
	if err != nil {
		t.Fatalf("FindByUid: %v", err)
	}

	if got.Id != work.Id || got.Title != "Good Omens" || !got.CreatedAt.Equal(seededAt) {
		t.Errorf("work = %+v", got)
	}
	if !reflect.DeepEqual(got.AuthorUids, []string{"a-1"}) || !reflect.DeepEqual(got.AuthorNames, []string{"Pratchett"}) {
		t.Errorf("authors = %v/%v, want only Pratchett", got.AuthorUids, got.AuthorNames)
	}
	if !reflect.DeepEqual(got.SubjectUids, []string{"s-1"}) || !reflect.DeepEqual(got.SubjectNames, []string{"Humor"}) {
		t.Errorf("subjects = %v/%v, want only Humor", got.SubjectUids, got.SubjectNames)
	}
}

func TestWorkRepositoryFindByUid_unknownUid_failsWithRecordNotFound(t *testing.T) {
	db := newFullCatalogDB(t)
	seedWork(t, db, "w-1", "The Hobbit")

	got, err := newWorkRepo(db).FindByUid(context.Background(), "w-missing")

	if !errors.Is(err, gorm.ErrRecordNotFound) || got != nil {
		t.Errorf("got = %v, err = %v, want nil and gorm.ErrRecordNotFound", got, err)
	}
}

func TestWorkRepositoryFindByUid_softDeletedWork_failsWithRecordNotFound(t *testing.T) {
	db := newFullCatalogDB(t)
	work := seedWork(t, db, "w-1", "The Hobbit")
	softDelete(t, db, "work", work.Id)

	_, err := newWorkRepo(db).FindByUid(context.Background(), "w-1")

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Errorf("err = %v, want gorm.ErrRecordNotFound", err)
	}
}

func TestWorkRepositoryFindOne_existingId_loadsWorkWithRelations(t *testing.T) {
	db := newFullCatalogDB(t)
	seedWork(t, db, "w-1", "Dune")
	work := seedWork(t, db, "w-2", "The Hobbit")
	tolkien := seedAuthor(t, db, "a-1", "Tolkien")
	linkAuthor(t, db, work.Id, tolkien.Id)
	setColumn(t, db, "work", work.Id, "cover_uids", `["c-1","c-2"]`)

	got, err := newWorkRepo(db).FindOne(context.Background(), work.Id)
	if err != nil {
		t.Fatalf("FindOne: %v", err)
	}

	if got.Uid != "w-2" || !reflect.DeepEqual(got.CoverUids, []string{"c-1", "c-2"}) {
		t.Errorf("work = %+v", got)
	}
	if !reflect.DeepEqual(got.AuthorUids, []string{"a-1"}) || len(got.SubjectUids) != 0 {
		t.Errorf("authors = %v, subjects = %v", got.AuthorUids, got.SubjectUids)
	}
}

func TestWorkRepositoryFindOne_unknownOrSoftDeletedId_failsWithRecordNotFound(t *testing.T) {
	db := newFullCatalogDB(t)
	gone := seedWork(t, db, "w-1", "The Hobbit")
	softDelete(t, db, "work", gone.Id)
	r := newWorkRepo(db)

	for _, id := range []int64{gone.Id, 999} {
		if _, err := r.FindOne(context.Background(), id); !errors.Is(err, gorm.ErrRecordNotFound) {
			t.Errorf("FindOne(%d) err = %v, want gorm.ErrRecordNotFound", id, err)
		}
	}
}

func TestWorkRepositoryCreate_authorsAndSubjectsGiven_persistsWorkAndReusesExistingSubjects(t *testing.T) {
	db := newFullCatalogDB(t)
	tolkien := seedAuthor(t, db, "a-1", "Tolkien")
	lewis := seedAuthor(t, db, "a-2", "Lewis")
	fantasy := seedSubject(t, db, "s-1", "Fantasy")
	r := newWorkRepo(db)

	id, err := r.Create(context.Background(), command.WorkCreateCommand{
		Title:            "Collected Letters",
		Description:      strPtr("letters"),
		CoverUids:        []string{"c-1"},
		FirstPublishDate: strPtr("1981"),
		OlKey:            strPtr("OL9W"),
		AuthorUids:       []string{"a-1", "a-2"},
		SubjectNames:     []string{"Fantasy", "Letters"},
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	row := loadRow[entity.WorkEntity](t, db, id)
	if parsed, err := uuid.Parse(row.Uid); err != nil || parsed.Version() != 7 {
		t.Errorf("uid = %q, want a UUIDv7", row.Uid)
	}
	if row.Title != "Collected Letters" || row.Version != 1 || row.CreatedAt.IsZero() {
		t.Errorf("row = %+v", row)
	}
	if row.CoverUid == nil || *row.CoverUid != `["c-1"]` {
		t.Errorf("cover_uids = %v, want [\"c-1\"]", row.CoverUid)
	}
	if got := mappedAuthorIds(t, db, id); !reflect.DeepEqual(got, []int64{tolkien.Id, lewis.Id}) {
		t.Errorf("author mappings = %v, want [%d %d]", got, tolkien.Id, lewis.Id)
	}
	if got := mappedSubjectNames(t, db, id); !reflect.DeepEqual(got, []string{"Fantasy", "Letters"}) {
		t.Errorf("subject mappings = %v, want [Fantasy Letters]", got)
	}
	if n := countRows(t, db, "subject", "name = ?", "Fantasy"); n != 1 {
		t.Errorf("Fantasy subject rows = %d, want the seeded one reused", n)
	}
	if n := countRows(t, db, "work_subject_mapping", "work_id = ? AND subject_id = ?", id, fantasy.Id); n != 1 {
		t.Errorf("mapping to the seeded Fantasy subject = %d, want 1", n)
	}

	got, err := r.FindOne(context.Background(), id)
	if err != nil {
		t.Fatalf("FindOne: %v", err)
	}
	if got.Description == nil || *got.Description != "letters" || got.FirstPublishDate == nil || *got.FirstPublishDate != "1981" ||
		got.OlKey == nil || *got.OlKey != "OL9W" {
		t.Errorf("work = %+v", got)
	}
	if !reflect.DeepEqual(sorted(got.AuthorUids), []string{"a-1", "a-2"}) || !reflect.DeepEqual(sorted(got.SubjectNames), []string{"Fantasy", "Letters"}) {
		t.Errorf("authors = %v, subjects = %v", got.AuthorUids, got.SubjectNames)
	}
}

func TestWorkRepositoryCreate_unknownAuthorUid_rollsBackWorkAndMappings(t *testing.T) {
	db := newFullCatalogDB(t)
	seedAuthor(t, db, "a-1", "Tolkien")

	id, err := newWorkRepo(db).Create(context.Background(), command.WorkCreateCommand{
		Title:        "Orphan",
		AuthorUids:   []string{"a-1", "a-missing"},
		SubjectNames: []string{"Fantasy"},
	})

	if !errors.Is(err, gorm.ErrRecordNotFound) || id != 0 {
		t.Fatalf("id = %d, err = %v, want 0 and gorm.ErrRecordNotFound", id, err)
	}
	if n := countRows(t, db, "work", "1 = 1"); n != 0 {
		t.Errorf("work rows = %d, want the insert rolled back", n)
	}
	if n := countRows(t, db, "work_author_mapping", "1 = 1"); n != 0 {
		t.Errorf("author mappings = %d, want none", n)
	}
	if n := countRows(t, db, "subject", "1 = 1"); n != 0 {
		t.Errorf("subject rows = %d, want none", n)
	}
}

func TestWorkRepositoryCreate_softDeletedAuthorUid_failsWithRecordNotFound(t *testing.T) {
	db := newFullCatalogDB(t)
	gone := seedAuthor(t, db, "a-1", "Tolkien")
	softDelete(t, db, "author", gone.Id)

	_, err := newWorkRepo(db).Create(context.Background(), command.WorkCreateCommand{
		Title:      "Orphan",
		AuthorUids: []string{"a-1"},
	})

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Errorf("err = %v, want gorm.ErrRecordNotFound", err)
	}
}

func TestWorkRepositoryUpdate_scalarFieldsGiven_persistsEachField(t *testing.T) {
	db := newFullCatalogDB(t)
	work := seedWork(t, db, "w-1", "before")

	err := newWorkRepo(db).Update(context.Background(), work.Id, command.WorkUpdateCommand{
		Description:      strPtr("after"),
		CoverUids:        []string{"c-9"},
		FirstPublishDate: strPtr("1937"),
		OlKey:            strPtr("OL7W"),
	})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}

	got := loadRow[entity.WorkEntity](t, db, work.Id)
	if got.Title != "before" || got.Description == nil || *got.Description != "after" ||
		got.CoverUid == nil || *got.CoverUid != `["c-9"]` ||
		got.FirstPublishDate == nil || *got.FirstPublishDate != "1937" ||
		got.OlKey == nil || *got.OlKey != "OL7W" || got.UpdatedAt == nil {
		t.Errorf("row = %+v", got)
	}
}

func TestWorkRepositoryUpdate_authorUidsGiven_replacesAuthorMappings(t *testing.T) {
	db := newFullCatalogDB(t)
	work := seedWork(t, db, "w-1", "Good Omens")
	other := seedWork(t, db, "w-2", "Dune")
	pratchett := seedAuthor(t, db, "a-1", "Pratchett")
	gaiman := seedAuthor(t, db, "a-2", "Gaiman")
	herbert := seedAuthor(t, db, "a-3", "Herbert")
	linkAuthor(t, db, work.Id, pratchett.Id)
	linkAuthor(t, db, other.Id, herbert.Id)

	err := newWorkRepo(db).Update(context.Background(), work.Id, command.WorkUpdateCommand{AuthorUids: []string{"a-2", "a-3"}})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}

	if got := mappedAuthorIds(t, db, work.Id); !reflect.DeepEqual(got, []int64{gaiman.Id, herbert.Id}) {
		t.Errorf("author mappings = %v, want [%d %d]", got, gaiman.Id, herbert.Id)
	}
	if got := mappedAuthorIds(t, db, other.Id); !reflect.DeepEqual(got, []int64{herbert.Id}) {
		t.Errorf("other work's mappings = %v, want them untouched", got)
	}
	if got := loadRow[entity.WorkEntity](t, db, work.Id); got.UpdatedAt == nil {
		t.Error("updated_at not touched by a relation-only update")
	}
}

func TestWorkRepositoryUpdate_emptyAuthorUidsGiven_clearsAuthorMappings(t *testing.T) {
	db := newFullCatalogDB(t)
	work := seedWork(t, db, "w-1", "Good Omens")
	pratchett := seedAuthor(t, db, "a-1", "Pratchett")
	linkAuthor(t, db, work.Id, pratchett.Id)

	err := newWorkRepo(db).Update(context.Background(), work.Id, command.WorkUpdateCommand{AuthorUids: []string{}})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}

	if got := mappedAuthorIds(t, db, work.Id); len(got) != 0 {
		t.Errorf("author mappings = %v, want none", got)
	}
}

func TestWorkRepositoryUpdate_subjectNamesGiven_replacesMappingsAndCreatesMissingSubjects(t *testing.T) {
	db := newFullCatalogDB(t)
	work := seedWork(t, db, "w-1", "Good Omens")
	humor := seedSubject(t, db, "s-1", "Humor")
	seedSubject(t, db, "s-2", "Fantasy")
	linkSubject(t, db, work.Id, humor.Id)

	err := newWorkRepo(db).Update(context.Background(), work.Id, command.WorkUpdateCommand{SubjectNames: []string{"Fantasy", "Apocalypse"}})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}

	if got := mappedSubjectNames(t, db, work.Id); !reflect.DeepEqual(got, []string{"Apocalypse", "Fantasy"}) {
		t.Errorf("subject mappings = %v, want [Apocalypse Fantasy]", got)
	}
	if n := countRows(t, db, "subject", "1 = 1"); n != 3 {
		t.Errorf("subject rows = %d, want 3 (two seeded, one created)", n)
	}
	if n := countRows(t, db, "subject", "name = ?", "Humor"); n != 1 {
		t.Errorf("Humor subject rows = %d, want the unlinked subject kept", n)
	}
}

func TestWorkRepositoryUpdate_unknownAuthorUid_rollsBackFieldAndMappingChanges(t *testing.T) {
	db := newFullCatalogDB(t)
	work := seedWork(t, db, "w-1", "before")
	pratchett := seedAuthor(t, db, "a-1", "Pratchett")
	linkAuthor(t, db, work.Id, pratchett.Id)

	err := newWorkRepo(db).Update(context.Background(), work.Id, command.WorkUpdateCommand{
		Title:      strPtr("after"),
		AuthorUids: []string{"a-missing"},
	})

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("err = %v, want gorm.ErrRecordNotFound", err)
	}
	if got := loadRow[entity.WorkEntity](t, db, work.Id); got.Title != "before" {
		t.Errorf("title = %q, want the rename rolled back", got.Title)
	}
	if got := mappedAuthorIds(t, db, work.Id); !reflect.DeepEqual(got, []int64{pratchett.Id}) {
		t.Errorf("author mappings = %v, want the original mapping restored", got)
	}
}

func TestWorkRepository_relationTablesUnavailable_propagatesQueryError(t *testing.T) {
	cases := []struct {
		name string
		ddl  []string
	}{
		{"authorMappingMissing", []string{workTableDDL}},
		{"subjectMappingMissing", []string{workTableDDL, authorTableDDL, workAuthorMappingTableDDL}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			db := newCatalogDB(t, c.ddl...)
			work := seedWork(t, db, "w-1", "The Hobbit")
			r := newWorkRepo(db)

			if _, _, err := r.FindAll(context.Background(), option.NewWorkQueryOption()); err == nil || errors.Is(err, gorm.ErrRecordNotFound) {
				t.Errorf("FindAll err = %v, want the query failure", err)
			}
			if _, err := r.FindOne(context.Background(), work.Id); err == nil || errors.Is(err, gorm.ErrRecordNotFound) {
				t.Errorf("FindOne err = %v, want the query failure", err)
			}
		})
	}
}

func TestWorkRepository_workTableUnavailable_propagatesQueryError(t *testing.T) {
	db := newCatalogDB(t)
	r := newWorkRepo(db)

	if _, _, err := r.FindAll(context.Background(), option.NewWorkQueryOption()); err == nil {
		t.Error("FindAll err = nil, want the query failure")
	}
	if id, err := r.Create(context.Background(), command.WorkCreateCommand{Title: "x"}); err == nil || id != 0 {
		t.Errorf("Create id = %d, err = %v, want 0 and the insert failure", id, err)
	}
	if err := r.Update(context.Background(), 1, command.WorkUpdateCommand{Title: strPtr("x")}); err == nil {
		t.Error("Update err = nil, want the update failure")
	}
}

func TestWorkRepositoryFindAll_titleFilterGiven_matchesSubstringIgnoringCase(t *testing.T) {
	db := newFullCatalogDB(t)
	seedWork(t, db, "w-1", "The Hobbit")
	seedWork(t, db, "w-2", "Dune")

	total, got, err := newWorkRepo(db).FindAll(context.Background(), option.NewWorkQueryOption(
		option.WorkQueryOption{}.WithTitle(strPtr("hOBb")),
	))
	if err != nil {
		t.Fatalf("FindAll: %v", err)
	}

	if total != 1 || !reflect.DeepEqual(workUids(got), []string{"w-1"}) {
		t.Errorf("total = %d, uids = %v, want 1 and [w-1]", total, workUids(got))
	}
}
