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

func newSubjectRepo(db *gorm.DB) *subjectRepository {
	return &subjectRepository{q: dao.Use(db)}
}

func subjectUids(results []*repo.SubjectResult) []string {
	uids := make([]string, len(results))
	for i, r := range results {
		uids[i] = r.Uid
	}
	return uids
}

func TestSubjectRepositoryFindAll_nameDescSortWithPaging_pagesAfterSortingAndKeepsTotal(t *testing.T) {
	db := newFullCatalogDB(t)
	seedSubject(t, db, "s-1", "Fantasy")
	seedSubject(t, db, "s-2", "Science Fiction")
	seedSubject(t, db, "s-3", "Adventure")
	gone := seedSubject(t, db, "s-4", "Western")
	softDelete(t, db, "subject", gone.Id)
	r := newSubjectRepo(db)

	total, first, err := r.FindAll(context.Background(), option.NewSubjectQueryOption(
		option.SubjectQueryOption{}.WithName(strPtr("")),
		option.SubjectQueryOption{}.WithPage(pageOf(2, 0, "name,desc")),
	))
	if err != nil {
		t.Fatalf("FindAll page 0: %v", err)
	}
	_, second, err := r.FindAll(context.Background(), option.NewSubjectQueryOption(
		option.SubjectQueryOption{}.WithPage(pageOf(2, 1, "name,desc")),
	))
	if err != nil {
		t.Fatalf("FindAll page 1: %v", err)
	}

	if total != 3 {
		t.Errorf("total = %d, want 3 (soft-deleted subject excluded)", total)
	}
	if !reflect.DeepEqual(subjectUids(first), []string{"s-2", "s-1"}) {
		t.Errorf("page 0 = %v, want [s-2 s-1]", subjectUids(first))
	}
	if !reflect.DeepEqual(subjectUids(second), []string{"s-3"}) {
		t.Errorf("page 1 = %v, want [s-3]", subjectUids(second))
	}
}

func TestSubjectRepositoryFindAll_noSortGiven_ordersById(t *testing.T) {
	db := newFullCatalogDB(t)
	seedSubject(t, db, "s-1", "Zoology")
	seedSubject(t, db, "s-2", "Art")

	total, got, err := newSubjectRepo(db).FindAll(context.Background(), option.NewSubjectQueryOption())
	if err != nil {
		t.Fatalf("FindAll: %v", err)
	}

	if total != 2 || !reflect.DeepEqual(subjectUids(got), []string{"s-1", "s-2"}) {
		t.Errorf("total = %d, uids = %v, want 2 and [s-1 s-2]", total, subjectUids(got))
	}
}

func TestSubjectRepositoryCreate_nameGiven_persistsAndReadsBack(t *testing.T) {
	db := newFullCatalogDB(t)
	r := newSubjectRepo(db)

	id, err := r.Create(context.Background(), command.SubjectCreateCommand{Name: "Fantasy"})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	row := loadRow[entity.SubjectEntity](t, db, id)
	if parsed, err := uuid.Parse(row.Uid); err != nil || parsed.Version() != 7 {
		t.Fatalf("uid = %q, want a UUIDv7", row.Uid)
	}
	if row.Version != 1 || row.CreatedAt.IsZero() {
		t.Errorf("row = %+v", row)
	}

	want := &repo.SubjectResult{Id: id, Uid: row.Uid, Name: "Fantasy"}
	byId, err := r.FindOne(context.Background(), id)
	if err != nil {
		t.Fatalf("FindOne: %v", err)
	}
	byUid, err := r.FindByUid(context.Background(), row.Uid)
	if err != nil {
		t.Fatalf("FindByUid: %v", err)
	}
	byName, err := r.FindByName(context.Background(), "Fantasy")
	if err != nil {
		t.Fatalf("FindByName: %v", err)
	}
	if !reflect.DeepEqual(byId, want) || !reflect.DeepEqual(byUid, want) || !reflect.DeepEqual(byName, want) {
		t.Errorf("byId = %+v, byUid = %+v, byName = %+v, want %+v", byId, byUid, byName, want)
	}
}

func TestSubjectRepository_unknownOrSoftDeletedSubject_failsWithRecordNotFound(t *testing.T) {
	db := newFullCatalogDB(t)
	gone := seedSubject(t, db, "s-1", "Western")
	softDelete(t, db, "subject", gone.Id)
	r := newSubjectRepo(db)

	for _, id := range []int64{gone.Id, 999} {
		if _, err := r.FindOne(context.Background(), id); !errors.Is(err, gorm.ErrRecordNotFound) {
			t.Errorf("FindOne(%d) err = %v, want gorm.ErrRecordNotFound", id, err)
		}
	}
	for _, uid := range []string{"s-1", "s-missing"} {
		if _, err := r.FindByUid(context.Background(), uid); !errors.Is(err, gorm.ErrRecordNotFound) {
			t.Errorf("FindByUid(%q) err = %v, want gorm.ErrRecordNotFound", uid, err)
		}
	}
}

func TestSubjectRepositoryFindByName_unknownName_yieldsNilWithoutError(t *testing.T) {
	db := newFullCatalogDB(t)
	seedSubject(t, db, "s-1", "Fantasy")

	got, err := newSubjectRepo(db).FindByName(context.Background(), "Fantas")

	if err != nil || got != nil {
		t.Errorf("got = %v, err = %v, want nil and no error", got, err)
	}
}

func TestSubjectRepositoryFindByName_queryFails_propagatesErrorInsteadOfNil(t *testing.T) {
	db := newCatalogDB(t)

	got, err := newSubjectRepo(db).FindByName(context.Background(), "Fantasy")

	if err == nil || got != nil {
		t.Errorf("got = %v, err = %v, want nil and the query failure", got, err)
	}
}

func TestSubjectRepositoryFindOrCreateByName_existingName_reusesThatSubject(t *testing.T) {
	db := newFullCatalogDB(t)
	fantasy := seedSubject(t, db, "s-1", "Fantasy")

	id, err := newSubjectRepo(db).FindOrCreateByName(context.Background(), "Fantasy")
	if err != nil {
		t.Fatalf("FindOrCreateByName: %v", err)
	}

	if id != fantasy.Id {
		t.Errorf("id = %d, want the seeded %d", id, fantasy.Id)
	}
	if n := countRows(t, db, "subject", "1 = 1"); n != 1 {
		t.Errorf("subject rows = %d, want 1", n)
	}
}

func TestSubjectRepositoryFindOrCreateByName_unknownName_createsSubject(t *testing.T) {
	db := newFullCatalogDB(t)
	fantasy := seedSubject(t, db, "s-1", "Fantasy")

	id, err := newSubjectRepo(db).FindOrCreateByName(context.Background(), "Horror")
	if err != nil {
		t.Fatalf("FindOrCreateByName: %v", err)
	}

	if id == 0 || id == fantasy.Id {
		t.Fatalf("id = %d, want a new subject id", id)
	}
	row := loadRow[entity.SubjectEntity](t, db, id)
	if parsed, err := uuid.Parse(row.Uid); err != nil || parsed.Version() != 7 {
		t.Errorf("uid = %q, want a UUIDv7", row.Uid)
	}
	if row.Name != "Horror" || row.Version != 1 || row.CreatedAt.IsZero() {
		t.Errorf("row = %+v", row)
	}
}

func TestSubjectRepositoryFindOrCreateByName_queryFails_propagatesErrorWithoutCreating(t *testing.T) {
	db := newCatalogDB(t)

	id, err := newSubjectRepo(db).FindOrCreateByName(context.Background(), "Horror")

	if err == nil || id != 0 {
		t.Errorf("id = %d, err = %v, want 0 and the query failure", id, err)
	}
}

func TestSubjectRepositoryUpdate_noFieldGiven_leavesRowUntouched(t *testing.T) {
	db, rec := newRecordedDB(t, subjectTableDDL)
	subject := seedRow(t, db, &entity.SubjectEntity{Uid: "s-1", Name: "before", CreatedAt: seededAt, Version: 1})

	if err := newSubjectRepo(db).Update(context.Background(), subject.Id, command.SubjectUpdateCommand{}); err != nil {
		t.Fatalf("Update: %v", err)
	}

	if updates := rec.updatesOf("subject"); len(updates) != 0 {
		t.Errorf("UPDATE statements = %q, want none", updates)
	}
	if got := loadRow[entity.SubjectEntity](t, db, subject.Id); got.Name != "before" {
		t.Errorf("name = %q, want before", got.Name)
	}
}

func TestSubjectRepository_tablesUnavailable_propagatesQueryError(t *testing.T) {
	db := newCatalogDB(t)
	r := newSubjectRepo(db)

	if _, _, err := r.FindAll(context.Background(), option.NewSubjectQueryOption()); err == nil {
		t.Error("FindAll err = nil, want the query failure")
	}
	if id, err := r.Create(context.Background(), command.SubjectCreateCommand{Name: "x"}); err == nil || id != 0 {
		t.Errorf("Create id = %d, err = %v, want 0 and the insert failure", id, err)
	}
}

func TestSubjectRepositoryFindAll_nameFilterGiven_matchesSubstringIgnoringCase(t *testing.T) {
	db := newFullCatalogDB(t)
	seedSubject(t, db, "s-1", "Fantasy")
	seedSubject(t, db, "s-2", "Science Fiction")

	total, got, err := newSubjectRepo(db).FindAll(context.Background(), option.NewSubjectQueryOption(
		option.SubjectQueryOption{}.WithName(strPtr("science")),
	))
	if err != nil {
		t.Fatalf("FindAll: %v", err)
	}

	if total != 1 || !reflect.DeepEqual(subjectUids(got), []string{"s-2"}) {
		t.Errorf("total = %d, uids = %v, want 1 and [s-2]", total, subjectUids(got))
	}
}
