package pgsql

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"

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
