package pgsql

import (
	"context"
	"fmt"
	"reflect"
	"sync/atomic"
	"testing"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/codejsha/shared-library-go/pkg/pagination"

	"github.com/codejsha/bookstore-microservices/inventory/internal/domain/model/option"
)

var auditSchemaDDL = []string{
	`CREATE TABLE IF NOT EXISTS stock_audit (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		uid TEXT, warehouse_id INTEGER, warehouse_uid TEXT, status TEXT, notes TEXT,
		completed_at DATETIME, created_at DATETIME, updated_at DATETIME, deleted_at DATETIME,
		actor INTEGER, version INTEGER
	)`,
	`CREATE TABLE IF NOT EXISTS stock_audit_item (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		audit_id INTEGER, edition_id INTEGER, edition_uid TEXT,
		system_quantity INTEGER, actual_quantity INTEGER, difference INTEGER
	)`,
}

func TestAuditFindAll_secondPage_skipsOneFullPage(t *testing.T) {
	_, db := newTestRepo(t)
	for _, ddl := range auditSchemaDDL {
		if err := db.Exec(ddl).Error; err != nil {
			t.Fatalf("create audit schema: %v", err)
		}
	}
	for i := 1; i <= 7; i++ {
		if err := db.Exec(
			"INSERT INTO stock_audit (uid, warehouse_id, warehouse_uid, status, notes, created_at, actor, version) VALUES (?, 1, ?, 'IN_PROGRESS', ?, ?, 1, 1)",
			uuid.NewString(), warehouseUidFor(1), fmt.Sprintf("audit-%d", i), time.Now().UTC(),
		).Error; err != nil {
			t.Fatalf("seed audit: %v", err)
		}
	}
	r := &stockAuditRepository{db: db}
	size, page, sort := int32(3), int32(1), "note,asc"

	total, got, err := r.FindAll(context.Background(), option.NewAuditQueryOption(
		option.AuditQueryOption{}.WithPage(pagination.NewPageOption(&size, &page, &sort)),
	))
	if err != nil {
		t.Fatalf("find all: %v", err)
	}

	if total != 7 {
		t.Errorf("total = %d, want 7", total)
	}
	notes := make([]string, len(got))
	for i, a := range got {
		if a.Notes != nil {
			notes[i] = *a.Notes
		}
	}
	if want := []string{"audit-4", "audit-5", "audit-6"}; !reflect.DeepEqual(notes, want) {
		t.Fatalf("notes = %v, want %v", notes, want)
	}
}

func countTableReads(db *gorm.DB, table string) *atomic.Int32 {
	var reads atomic.Int32
	_ = db.Callback().Query().After("gorm:query").Register("test:count_"+table+"_query", func(tx *gorm.DB) {
		if tx.Statement.Table == table {
			reads.Add(1)
		}
	})
	return &reads
}

func TestAuditFindAll_severalAudits_readsItemsOnce(t *testing.T) {
	_, db := newTestRepo(t)
	for _, ddl := range auditSchemaDDL {
		if err := db.Exec(ddl).Error; err != nil {
			t.Fatalf("create audit schema: %v", err)
		}
	}
	for i := int64(1); i <= 3; i++ {
		if err := db.Exec(
			"INSERT INTO stock_audit (uid, warehouse_id, warehouse_uid, status, notes, created_at, actor, version) VALUES (?, 1, ?, 'IN_PROGRESS', ?, ?, 1, 1)",
			uuid.NewString(), warehouseUidFor(1), fmt.Sprintf("audit-%d", i), time.Now().UTC(),
		).Error; err != nil {
			t.Fatalf("seed audit: %v", err)
		}
		for edition := int64(1); edition <= 2; edition++ {
			if err := db.Exec(
				"INSERT INTO stock_audit_item (audit_id, edition_id, edition_uid, system_quantity, actual_quantity, difference) VALUES (?, ?, ?, 0, 0, 0)",
				i, edition, editionUidFor(edition),
			).Error; err != nil {
				t.Fatalf("seed audit item: %v", err)
			}
		}
	}
	reads := countTableReads(db, "stock_audit_item")
	r := &stockAuditRepository{db: db}
	size, page := int32(3), int32(0)

	_, got, err := r.FindAll(context.Background(), option.NewAuditQueryOption(
		option.AuditQueryOption{}.WithPage(pagination.NewPageOption(&size, &page, nil)),
	))
	if err != nil {
		t.Fatalf("find all: %v", err)
	}

	if reads.Load() != 1 {
		t.Errorf("stock_audit_item reads = %d, want 1", reads.Load())
	}
	if len(got) != 3 {
		t.Fatalf("audits = %d, want 3", len(got))
	}
	for _, a := range got {
		if len(a.Items) != 2 {
			t.Errorf("audit %s items = %d, want 2", a.Uid, len(a.Items))
		}
	}
}
