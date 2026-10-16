package pgsql

import (
	"context"
	"strings"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/codejsha/bookstore-microservices/catalog/generated/infrastructure/port/dao"
	"github.com/codejsha/bookstore-microservices/catalog/generated/infrastructure/port/entity"
	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/model/command"
)

const workTableDDL = `CREATE TABLE work (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	uid TEXT NOT NULL, title TEXT NOT NULL, description TEXT,
	cover_ids TEXT, cover_uids TEXT, first_publish_date TEXT, ol_key TEXT,
	created_at DATETIME NOT NULL, updated_at DATETIME, deleted_at DATETIME,
	actor INTEGER NOT NULL, version INTEGER NOT NULL
)`

const authorTableDDL = `CREATE TABLE author (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	uid TEXT NOT NULL, name TEXT NOT NULL, bio TEXT, birth_date TEXT, death_date TEXT,
	photo_ids TEXT, photo_uids TEXT, alternate_names TEXT, ol_key TEXT,
	created_at DATETIME NOT NULL, updated_at DATETIME, deleted_at DATETIME,
	actor INTEGER NOT NULL, version INTEGER NOT NULL
)`

const publisherTableDDL = `CREATE TABLE publisher (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	uid TEXT NOT NULL, name TEXT NOT NULL, address TEXT, ol_key TEXT,
	created_at DATETIME NOT NULL, updated_at DATETIME, deleted_at DATETIME,
	actor INTEGER NOT NULL, version INTEGER NOT NULL
)`

const subjectTableDDL = `CREATE TABLE subject (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	uid TEXT NOT NULL, name TEXT NOT NULL,
	created_at DATETIME NOT NULL, updated_at DATETIME, deleted_at DATETIME,
	actor INTEGER NOT NULL, version INTEGER NOT NULL
)`

const workSubjectMappingTableDDL = `CREATE TABLE work_subject_mapping (
	work_id INTEGER NOT NULL, subject_id INTEGER NOT NULL,
	PRIMARY KEY (work_id, subject_id)
)`

type sqlRecorder struct {
	logger.Interface
	statements []string
}

func (r *sqlRecorder) Trace(_ context.Context, _ time.Time, fc func() (string, int64), _ error) {
	sql, _ := fc()
	r.statements = append(r.statements, sql)
}

func (r *sqlRecorder) updatesOf(table string) []string {
	var updates []string
	for _, s := range r.statements {
		if strings.HasPrefix(s, "UPDATE `"+table+"`") {
			updates = append(updates, s)
		}
	}
	return updates
}

func newRecordedDB(t *testing.T, ddl ...string) (*gorm.DB, *sqlRecorder) {
	t.Helper()
	rec := &sqlRecorder{Interface: logger.Default.LogMode(logger.Silent)}
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: rec})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("sql db: %v", err)
	}
	sqlDB.SetMaxOpenConns(1)
	t.Cleanup(func() { _ = sqlDB.Close() })
	for _, stmt := range ddl {
		if err := db.Exec(stmt).Error; err != nil {
			t.Fatalf("create table: %v", err)
		}
	}
	return db, rec
}

func seedRow[T any](t *testing.T, db *gorm.DB, row *T) *T {
	t.Helper()
	if err := db.Create(row).Error; err != nil {
		t.Fatalf("seed %T: %v", row, err)
	}
	return row
}

func loadRow[T any](t *testing.T, db *gorm.DB, id int64) T {
	t.Helper()
	var got T
	if err := db.Unscoped().Where("id = ?", id).Take(&got).Error; err != nil {
		t.Fatalf("load %T %d: %v", got, id, err)
	}
	return got
}

func assertUpdatedAtAssignedOnce(t *testing.T, rec *sqlRecorder, table string) {
	t.Helper()
	updates := rec.updatesOf(table)
	if len(updates) != 1 {
		t.Fatalf("UPDATE statements on %s = %d, want 1: %q", table, len(updates), updates)
	}
	if n := strings.Count(updates[0], "`updated_at`"); n != 1 {
		t.Errorf("updated_at assigned %d times, want 1: %s", n, updates[0])
	}
}

func TestWorkRepository_updateTitle_assignsUpdatedAtOnce(t *testing.T) {
	db, rec := newRecordedDB(t, workTableDDL)
	work := seedRow(t, db, &entity.WorkEntity{Uid: testWorkUid, Title: "before", CreatedAt: time.Now(), Version: 1})
	r := &workRepository{q: dao.Use(db), db: db}
	title := "after"

	if err := r.Update(context.Background(), work.Id, command.WorkUpdateCommand{Title: &title}); err != nil {
		t.Fatalf("update: %v", err)
	}

	assertUpdatedAtAssignedOnce(t, rec, "work")
	got := loadRow[entity.WorkEntity](t, db, work.Id)
	if got.Title != "after" || got.UpdatedAt == nil {
		t.Errorf("title = %q, updated_at = %v", got.Title, got.UpdatedAt)
	}
}

func TestWorkRepository_updateSubjectsOnly_touchesUpdatedAtOnce(t *testing.T) {
	db, rec := newRecordedDB(t, workTableDDL, subjectTableDDL, workSubjectMappingTableDDL)
	work := seedRow(t, db, &entity.WorkEntity{Uid: testWorkUid, Title: "before", CreatedAt: time.Now(), Version: 1})
	r := &workRepository{q: dao.Use(db), db: db}

	if err := r.Update(context.Background(), work.Id, command.WorkUpdateCommand{SubjectNames: []string{"history"}}); err != nil {
		t.Fatalf("update: %v", err)
	}

	assertUpdatedAtAssignedOnce(t, rec, "work")
	got := loadRow[entity.WorkEntity](t, db, work.Id)
	if got.Title != "before" || got.UpdatedAt == nil {
		t.Errorf("title = %q, updated_at = %v", got.Title, got.UpdatedAt)
	}
	var mappings int64
	if err := db.Table("work_subject_mapping").Where("work_id = ?", work.Id).Count(&mappings).Error; err != nil {
		t.Fatalf("count mappings: %v", err)
	}
	if mappings != 1 {
		t.Errorf("subject mappings = %d, want 1", mappings)
	}
}

func TestWorkRepository_updateWithoutChanges_skipsRowUpdate(t *testing.T) {
	db, rec := newRecordedDB(t, workTableDDL)
	seeded := time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)
	work := seedRow(t, db, &entity.WorkEntity{Uid: testWorkUid, Title: "before", CreatedAt: seeded, UpdatedAt: &seeded, Version: 1})
	r := &workRepository{q: dao.Use(db), db: db}

	if err := r.Update(context.Background(), work.Id, command.WorkUpdateCommand{}); err != nil {
		t.Fatalf("update: %v", err)
	}

	if updates := rec.updatesOf("work"); len(updates) != 0 {
		t.Errorf("UPDATE statements = %q, want none", updates)
	}
	if got := loadRow[entity.WorkEntity](t, db, work.Id); got.UpdatedAt == nil || !got.UpdatedAt.Equal(seeded) {
		t.Errorf("updated_at = %v, want %v", got.UpdatedAt, seeded)
	}
}

func TestAuthorRepository_updateName_assignsUpdatedAtOnce(t *testing.T) {
	db, rec := newRecordedDB(t, authorTableDDL)
	author := seedRow(t, db, &entity.AuthorEntity{Uid: testWorkUid, Name: "before", CreatedAt: time.Now(), Version: 1})
	r := &authorRepository{q: dao.Use(db)}
	name := "after"

	if err := r.Update(context.Background(), author.Id, command.AuthorUpdateCommand{Name: &name}); err != nil {
		t.Fatalf("update: %v", err)
	}

	assertUpdatedAtAssignedOnce(t, rec, "author")
	got := loadRow[entity.AuthorEntity](t, db, author.Id)
	if got.Name != "after" || got.UpdatedAt == nil {
		t.Errorf("name = %q, updated_at = %v", got.Name, got.UpdatedAt)
	}
}

func TestAuthorRepository_updateWithoutChanges_skipsRowUpdate(t *testing.T) {
	db, rec := newRecordedDB(t, authorTableDDL)
	seeded := time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)
	author := seedRow(t, db, &entity.AuthorEntity{Uid: testWorkUid, Name: "before", CreatedAt: seeded, UpdatedAt: &seeded, Version: 1})
	r := &authorRepository{q: dao.Use(db)}

	if err := r.Update(context.Background(), author.Id, command.AuthorUpdateCommand{}); err != nil {
		t.Fatalf("update: %v", err)
	}

	if updates := rec.updatesOf("author"); len(updates) != 0 {
		t.Errorf("UPDATE statements = %q, want none", updates)
	}
	if got := loadRow[entity.AuthorEntity](t, db, author.Id); got.UpdatedAt == nil || !got.UpdatedAt.Equal(seeded) {
		t.Errorf("updated_at = %v, want %v", got.UpdatedAt, seeded)
	}
}

func TestEditionRepository_updateTitle_assignsUpdatedAtOnce(t *testing.T) {
	db, rec := newRecordedDB(t, editionTableDDL)
	edition := seedRow(t, db, &entity.EditionEntity{Uid: testOtherWorkUid, Title: "before", WorkId: 1, WorkUid: testWorkUid, CreatedAt: time.Now(), Version: 1})
	r := &editionRepository{q: dao.Use(db)}
	title := "after"

	if err := r.Update(context.Background(), edition.Id, command.EditionUpdateCommand{Title: &title}, nil, nil); err != nil {
		t.Fatalf("update: %v", err)
	}

	assertUpdatedAtAssignedOnce(t, rec, "edition")
	got := loadRow[entity.EditionEntity](t, db, edition.Id)
	if got.Title != "after" || got.UpdatedAt == nil {
		t.Errorf("title = %q, updated_at = %v", got.Title, got.UpdatedAt)
	}
}

func TestPublisherRepository_updateName_assignsUpdatedAtOnce(t *testing.T) {
	db, rec := newRecordedDB(t, publisherTableDDL)
	publisher := seedRow(t, db, &entity.PublisherEntity{Uid: testPublisherUid, Name: "before", CreatedAt: time.Now(), Version: 1})
	r := &publisherRepository{q: dao.Use(db)}
	name := "after"

	if err := r.Update(context.Background(), publisher.Id, command.PublisherUpdateCommand{Name: &name}); err != nil {
		t.Fatalf("update: %v", err)
	}

	assertUpdatedAtAssignedOnce(t, rec, "publisher")
	got := loadRow[entity.PublisherEntity](t, db, publisher.Id)
	if got.Name != "after" || got.UpdatedAt == nil {
		t.Errorf("name = %q, updated_at = %v", got.Name, got.UpdatedAt)
	}
}

func TestSubjectRepository_updateName_assignsUpdatedAtOnce(t *testing.T) {
	db, rec := newRecordedDB(t, subjectTableDDL)
	subject := seedRow(t, db, &entity.SubjectEntity{Uid: testWorkUid, Name: "before", CreatedAt: time.Now(), Version: 1})
	r := &subjectRepository{q: dao.Use(db)}
	name := "after"

	if err := r.Update(context.Background(), subject.Id, command.SubjectUpdateCommand{Name: &name}); err != nil {
		t.Fatalf("update: %v", err)
	}

	assertUpdatedAtAssignedOnce(t, rec, "subject")
	got := loadRow[entity.SubjectEntity](t, db, subject.Id)
	if got.Name != "after" || got.UpdatedAt == nil {
		t.Errorf("name = %q, updated_at = %v", got.Name, got.UpdatedAt)
	}
}
