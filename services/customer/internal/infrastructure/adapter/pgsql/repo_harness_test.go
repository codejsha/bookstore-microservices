package pgsql

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/codejsha/shared-library-go/pkg/pagination"

	"github.com/codejsha/bookstore-microservices/customer/generated/infrastructure/port/dao"
	"github.com/codejsha/bookstore-microservices/customer/generated/infrastructure/port/entity"
)

const (
	userA     = "018f0000-0000-7000-8000-0000000000a1"
	userB     = "018f0000-0000-7000-8000-0000000000a2"
	bookA     = "018f0000-0000-7000-8000-0000000000b1"
	bookB     = "018f0000-0000-7000-8000-0000000000b2"
	bookC     = "018f0000-0000-7000-8000-0000000000b3"
	absentUid = "018f0000-0000-7000-8000-0000000000ff"
)

var errInjected = errors.New("injected database failure")

var schemaDDL = []string{
	`CREATE TABLE customer_wishlist (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		uid TEXT NOT NULL UNIQUE,
		user_id BIGINT NOT NULL, user_uid TEXT NOT NULL,
		book_id BIGINT NOT NULL, book_uid TEXT NOT NULL,
		created_at DATETIME NOT NULL, updated_at DATETIME, deleted_at DATETIME,
		actor BIGINT NOT NULL, version BIGINT NOT NULL
	)`,
	`CREATE UNIQUE INDEX uq_customer_wishlist_user_book_uid ON customer_wishlist (user_uid, book_uid) WHERE deleted_at IS NULL`,
	`CREATE TABLE customer_point (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		uid TEXT NOT NULL UNIQUE,
		user_id BIGINT NOT NULL, user_uid TEXT NOT NULL UNIQUE,
		balance INTEGER NOT NULL,
		created_at DATETIME NOT NULL, updated_at DATETIME, deleted_at DATETIME,
		actor BIGINT NOT NULL, version BIGINT NOT NULL
	)`,
	`CREATE TABLE customer_point_history (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		uid TEXT NOT NULL UNIQUE,
		user_id BIGINT NOT NULL, user_uid TEXT NOT NULL,
		change_type TEXT NOT NULL, amount INTEGER NOT NULL, reason TEXT,
		created_at DATETIME NOT NULL
	)`,
	`CREATE TABLE customer_review (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		uid TEXT NOT NULL UNIQUE,
		user_id BIGINT NOT NULL, user_uid TEXT NOT NULL,
		book_id BIGINT NOT NULL, book_uid TEXT NOT NULL,
		rating SMALLINT NOT NULL, title TEXT, content TEXT,
		created_at DATETIME NOT NULL, updated_at DATETIME, deleted_at DATETIME,
		actor BIGINT NOT NULL, version BIGINT NOT NULL
	)`,
	`CREATE UNIQUE INDEX uq_review_user_book_uid ON customer_review (user_uid, book_uid) WHERE deleted_at IS NULL`,
	`CREATE TABLE customer_ticket_category (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		uid TEXT NOT NULL UNIQUE,
		name TEXT NOT NULL, description TEXT, parent_id BIGINT,
		created_at DATETIME NOT NULL, updated_at DATETIME, deleted_at DATETIME,
		actor BIGINT NOT NULL, version BIGINT NOT NULL
	)`,
	`CREATE UNIQUE INDEX uq_ticket_category_name ON customer_ticket_category (name) WHERE deleted_at IS NULL`,
	`CREATE TABLE customer_ticket (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		uid TEXT NOT NULL UNIQUE,
		customer_uid TEXT NOT NULL, category_id BIGINT, assignee_uid TEXT,
		subject TEXT NOT NULL, description TEXT NOT NULL,
		status TEXT NOT NULL, priority TEXT NOT NULL, resolved_at DATETIME,
		created_at DATETIME NOT NULL, updated_at DATETIME, deleted_at DATETIME,
		actor BIGINT NOT NULL, version BIGINT NOT NULL
	)`,
	`CREATE TABLE customer_ticket_comment (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		uid TEXT NOT NULL UNIQUE,
		ticket_id BIGINT NOT NULL, author_uid TEXT NOT NULL, author_role TEXT NOT NULL,
		body TEXT NOT NULL, internal BOOLEAN NOT NULL,
		created_at DATETIME NOT NULL, updated_at DATETIME, deleted_at DATETIME,
		actor BIGINT NOT NULL, version BIGINT NOT NULL
	)`,
	`CREATE TABLE customer_faq (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		uid TEXT NOT NULL UNIQUE,
		category_id BIGINT, question TEXT NOT NULL, answer TEXT NOT NULL,
		view_count BIGINT NOT NULL, published BOOLEAN NOT NULL,
		created_at DATETIME NOT NULL, updated_at DATETIME, deleted_at DATETIME,
		actor BIGINT NOT NULL, version BIGINT NOT NULL
	)`,
}

var baseTime = time.Date(2026, 11, 1, 0, 0, 0, 0, time.UTC)

func hoursAfterBase(h int) time.Time {
	return baseTime.Add(time.Duration(h) * time.Hour)
}

func newSchemaDB(t *testing.T) *gorm.DB {
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
	for _, ddl := range schemaDDL {
		if err := db.Exec(ddl).Error; err != nil {
			t.Fatalf("apply schema: %v", err)
		}
	}
	return db
}

func closeDatabase(t *testing.T, db *gorm.DB) {
	t.Helper()
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("sql db: %v", err)
	}
	if err := sqlDB.Close(); err != nil {
		t.Fatalf("close db: %v", err)
	}
}

func failQueriesOn(t *testing.T, db *gorm.DB, table string) {
	t.Helper()
	failQueriesAfter(t, db, table, 0)
}

func failQueriesAfter(t *testing.T, db *gorm.DB, table string, allowed int) {
	t.Helper()
	seen := 0
	err := db.Callback().Query().Before("gorm:query").Register("test:fail_query_"+table, func(tx *gorm.DB) {
		if tx.Statement.Table != table {
			return
		}
		seen++
		if seen > allowed {
			_ = tx.AddError(errInjected)
		}
	})
	if err != nil {
		t.Fatalf("register query failure: %v", err)
	}
}

func failCreatesOn(t *testing.T, db *gorm.DB, table string, failure error) {
	t.Helper()
	err := db.Callback().Create().Before("gorm:create").Register("test:fail_create_"+table, func(tx *gorm.DB) {
		if tx.Statement.Table == table {
			_ = tx.AddError(failure)
		}
	})
	if err != nil {
		t.Fatalf("register create failure: %v", err)
	}
}

func uniqueViolation() error {
	return &pgconn.PgError{Code: "23505", Message: "duplicate key value violates unique constraint"}
}

func pageOf(size, page int32, sort string) pagination.PageOption {
	if sort == "" {
		return pagination.NewPageOption(&size, &page, nil)
	}
	return pagination.NewPageOption(&size, &page, &sort)
}

func i64p(v int64) *int64 { return &v }

func mustCreate(t *testing.T, db *gorm.DB, row any) {
	t.Helper()
	if err := db.Create(row).Error; err != nil {
		t.Fatalf("seed %T: %v", row, err)
	}
}

func softDelete(t *testing.T, db *gorm.DB, row any) {
	t.Helper()
	if err := db.Delete(row).Error; err != nil {
		t.Fatalf("soft delete %T: %v", row, err)
	}
}

func seedCategory(t *testing.T, db *gorm.DB, uid, name string, parentId *int64) *entity.CustomerTicketCategoryEntity {
	t.Helper()
	row := &entity.CustomerTicketCategoryEntity{Uid: uid, Name: name, ParentId: parentId, CreatedAt: baseTime, Version: 1}
	mustCreate(t, db, row)
	return row
}

func assertInjected(t *testing.T, err error) {
	t.Helper()
	if !errors.Is(err, errInjected) {
		t.Fatalf("err = %v, want the injected database failure", err)
	}
}

func assertFails(t *testing.T, name string, err error) {
	t.Helper()
	if err == nil {
		t.Errorf("%s: err = nil, want the database failure surfaced", name)
	}
}

func bg() context.Context { return context.Background() }

func queryOf(db *gorm.DB) *dao.Query { return dao.Use(db) }
