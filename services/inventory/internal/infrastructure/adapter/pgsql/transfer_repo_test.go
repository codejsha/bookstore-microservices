package pgsql

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/codejsha/shared-library-go/pkg/pagination"

	"github.com/codejsha/bookstore-microservices/inventory/internal/domain/model/option"
)

const transferSchemaDDL = `CREATE TABLE IF NOT EXISTS stock_transfer (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	uid TEXT, edition_id INTEGER, edition_uid TEXT,
	source_warehouse_id INTEGER, source_warehouse_uid TEXT,
	target_warehouse_id INTEGER, target_warehouse_uid TEXT,
	quantity INTEGER, status TEXT, reason TEXT,
	completed_at DATETIME, created_at DATETIME, updated_at DATETIME, deleted_at DATETIME,
	actor INTEGER, version INTEGER
)`

func TestTransferFindAll_secondPage_skipsOneFullPage(t *testing.T) {
	_, db := newTestRepo(t)
	if err := db.Exec(transferSchemaDDL).Error; err != nil {
		t.Fatalf("create transfer schema: %v", err)
	}
	for i := int32(1); i <= 7; i++ {
		if err := db.Exec(
			"INSERT INTO stock_transfer (uid, edition_id, edition_uid, source_warehouse_id, source_warehouse_uid, target_warehouse_id, target_warehouse_uid, quantity, status, created_at, actor, version) VALUES (?, 1, ?, 1, ?, 2, ?, ?, 'PENDING', ?, 1, 1)",
			uuid.NewString(), editionUidFor(1), warehouseUidFor(1), warehouseUidFor(2), i, time.Now().UTC(),
		).Error; err != nil {
			t.Fatalf("seed transfer: %v", err)
		}
	}
	r := &stockTransferRepository{db: db}
	size, page, sort := int32(3), int32(1), "quantity,asc"

	total, got, err := r.FindAll(context.Background(), option.NewTransferQueryOption(
		option.TransferQueryOption{}.WithPage(pagination.NewPageOption(&size, &page, &sort)),
	))
	if err != nil {
		t.Fatalf("find all: %v", err)
	}

	if total != 7 {
		t.Errorf("total = %d, want 7", total)
	}
	quantities := make([]int32, len(got))
	for i, tr := range got {
		quantities[i] = tr.Quantity
	}
	if want := []int32{4, 5, 6}; !reflect.DeepEqual(quantities, want) {
		t.Fatalf("quantities = %v, want %v", quantities, want)
	}
}
