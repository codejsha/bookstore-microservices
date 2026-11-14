package pgsql

import (
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/codejsha/shared-library-go/pkg/pagination"

	"github.com/codejsha/bookstore-microservices/catalog/generated/infrastructure/port/entity"
)

const workAuthorMappingTableDDL = `CREATE TABLE work_author_mapping (
	work_id INTEGER NOT NULL, author_id INTEGER NOT NULL,
	PRIMARY KEY (work_id, author_id)
)`

var catalogSchemaDDL = []string{
	workTableDDL, authorTableDDL, publisherTableDDL, subjectTableDDL,
	editionTableDDL, workAuthorMappingTableDDL, workSubjectMappingTableDDL,
}

var seededAt = time.Date(2026, 10, 1, 9, 0, 0, 0, time.UTC)

func newCatalogDB(t *testing.T, ddl ...string) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
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
	return db
}

func newFullCatalogDB(t *testing.T) *gorm.DB {
	t.Helper()
	return newCatalogDB(t, catalogSchemaDDL...)
}

func pageOf(size, page int32, sort string) pagination.PageOption {
	return pagination.NewPageOption(&size, &page, &sort)
}

func strPtr(s string) *string { return &s }

func seedWork(t *testing.T, db *gorm.DB, uid, title string) *entity.WorkEntity {
	t.Helper()
	return seedRow(t, db, &entity.WorkEntity{Uid: uid, Title: title, CreatedAt: seededAt, Version: 1})
}

func seedAuthor(t *testing.T, db *gorm.DB, uid, name string) *entity.AuthorEntity {
	t.Helper()
	return seedRow(t, db, &entity.AuthorEntity{Uid: uid, Name: name, CreatedAt: seededAt, Version: 1})
}

func seedSubject(t *testing.T, db *gorm.DB, uid, name string) *entity.SubjectEntity {
	t.Helper()
	return seedRow(t, db, &entity.SubjectEntity{Uid: uid, Name: name, CreatedAt: seededAt, Version: 1})
}

func seedPublisher(t *testing.T, db *gorm.DB, uid, name string) *entity.PublisherEntity {
	t.Helper()
	return seedRow(t, db, &entity.PublisherEntity{Uid: uid, Name: name, CreatedAt: seededAt, Version: 1})
}

func seedEdition(t *testing.T, db *gorm.DB, uid, title string, work *entity.WorkEntity, publisher *entity.PublisherEntity) *entity.EditionEntity {
	t.Helper()
	e := &entity.EditionEntity{Uid: uid, Title: title, WorkId: work.Id, WorkUid: work.Uid, CreatedAt: seededAt, Version: 1}
	if publisher != nil {
		e.PublisherId = &publisher.Id
		e.PublisherUid = &publisher.Uid
	}
	return seedRow(t, db, e)
}

func linkAuthor(t *testing.T, db *gorm.DB, workID, authorID int64) {
	t.Helper()
	seedRow(t, db, &entity.WorkAuthorMappingEntity{WorkId: workID, AuthorId: authorID})
}

func linkSubject(t *testing.T, db *gorm.DB, workID, subjectID int64) {
	t.Helper()
	seedRow(t, db, &entity.WorkSubjectMappingEntity{WorkId: workID, SubjectId: subjectID})
}

func softDelete(t *testing.T, db *gorm.DB, table string, id int64) {
	t.Helper()
	if err := db.Exec("UPDATE "+table+" SET deleted_at = ? WHERE id = ?", seededAt, id).Error; err != nil {
		t.Fatalf("soft delete %s %d: %v", table, id, err)
	}
}

func setColumn(t *testing.T, db *gorm.DB, table string, id int64, column string, value any) {
	t.Helper()
	if err := db.Exec("UPDATE "+table+" SET "+column+" = ? WHERE id = ?", value, id).Error; err != nil {
		t.Fatalf("set %s.%s on %d: %v", table, column, id, err)
	}
}

func countRows(t *testing.T, db *gorm.DB, table, where string, args ...any) int64 {
	t.Helper()
	var n int64
	if err := db.Table(table).Where(where, args...).Count(&n).Error; err != nil {
		t.Fatalf("count %s: %v", table, err)
	}
	return n
}
