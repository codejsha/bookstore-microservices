package pgsql

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/codejsha/shared-library-go/pkg/pagination"

	"github.com/codejsha/bookstore-microservices/inventory/internal/application/port/repo"
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

func seedTransfer(t *testing.T, db *gorm.DB, status string) string {
	t.Helper()
	if err := db.Exec(transferSchemaDDL).Error; err != nil {
		t.Fatalf("create transfer schema: %v", err)
	}
	uid := uuid.NewString()
	if err := db.Exec(
		"INSERT INTO stock_transfer (uid, edition_id, edition_uid, source_warehouse_id, source_warehouse_uid, target_warehouse_id, target_warehouse_uid, quantity, status, created_at, actor, version) VALUES (?, 7, ?, 1, ?, 2, ?, 3, ?, ?, 1, 1)",
		uid, editionUidFor(7), warehouseUidFor(1), warehouseUidFor(2), status, time.Now().UTC(),
	).Error; err != nil {
		t.Fatalf("seed transfer: %v", err)
	}
	return uid
}

func TestTransferComplete_inTransit_completesAndCreditsTarget(t *testing.T) {
	_, db := newTestRepo(t)
	seedStock(t, db, 7, 1, 10)
	uid := seedTransfer(t, db, "IN_TRANSIT")
	r := &stockTransferRepository{db: db}

	got, err := r.Complete(context.Background(), uid)
	if err != nil {
		t.Fatalf("complete: %v", err)
	}
	if got.Status != transferStatusCompleted {
		t.Errorf("status = %s, want COMPLETED", got.Status)
	}
	if qty := stockQty(t, db, 7, 2); qty != 3 {
		t.Errorf("target qty = %d, want 3", qty)
	}
}

func TestTransferComplete_alreadyCompleted_invalidTransition(t *testing.T) {
	_, db := newTestRepo(t)
	uid := seedTransfer(t, db, "COMPLETED")
	r := &stockTransferRepository{db: db}

	if _, err := r.Complete(context.Background(), uid); !errors.Is(err, repo.ErrInvalidTransition) {
		t.Fatalf("err = %v, want ErrInvalidTransition", err)
	}
}

func TestTransferCancel_inTransit_invalidTransition(t *testing.T) {
	_, db := newTestRepo(t)
	uid := seedTransfer(t, db, "IN_TRANSIT")
	r := &stockTransferRepository{db: db}

	if _, err := r.Cancel(context.Background(), uid); !errors.Is(err, repo.ErrInvalidTransition) {
		t.Fatalf("err = %v, want ErrInvalidTransition", err)
	}
}
