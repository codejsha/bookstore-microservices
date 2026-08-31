package pgsql

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/codejsha/bookstore-microservices/inventory/internal/application/port/repo"
)

func TestAuditCreate_WhenWarehouseMissing_ReturnsWrappedWarehouseNotFound(t *testing.T) {
	_, db := newTestRepo(t)
	r := &stockAuditRepository{db: db}

	_, err := r.Create(context.Background(), repo.AuditCreateParams{
		WarehouseUid: warehouseUidFor(9),
		Items:        []repo.AuditItemParam{{EditionUid: editionUidFor(7), ActualQuantity: 2}},
	})
	if !errors.Is(err, repo.ErrWarehouseNotFound) {
		t.Fatalf("err = %v, want ErrWarehouseNotFound", err)
	}
	if !strings.Contains(err.Error(), warehouseUidFor(9)) {
		t.Errorf("err = %q, want descriptive warehouse context", err)
	}
}

func TestClosingCreate_WhenWarehouseMissing_ReturnsWrappedWarehouseNotFound(t *testing.T) {
	_, db := newTestRepo(t)
	r := &monthlyClosingRepository{db: db}

	_, err := r.Create(context.Background(), repo.ClosingCreateParams{
		WarehouseUid: warehouseUidFor(9),
		Year:         2026,
		Month:        7,
	})
	if !errors.Is(err, repo.ErrWarehouseNotFound) {
		t.Fatalf("err = %v, want ErrWarehouseNotFound", err)
	}
	if !strings.Contains(err.Error(), warehouseUidFor(9)) {
		t.Errorf("err = %q, want descriptive warehouse context", err)
	}
}

func TestTransferCreate_WhenTargetWarehouseMissing_ReturnsWrappedWarehouseNotFound(t *testing.T) {
	_, db := newTestRepo(t)
	seedStock(t, db, 7, 9, 10)
	r := &stockTransferRepository{db: db}

	_, err := r.Create(context.Background(), repo.TransferCreateParams{
		EditionUid:         editionUidFor(7),
		SourceWarehouseUid: warehouseUidFor(9),
		TargetWarehouseUid: warehouseUidFor(11),
		Quantity:           3,
	})
	if !errors.Is(err, repo.ErrWarehouseNotFound) {
		t.Fatalf("err = %v, want ErrWarehouseNotFound", err)
	}
	if !strings.Contains(err.Error(), warehouseUidFor(11)) {
		t.Errorf("err = %q, want descriptive warehouse context", err)
	}
}

func TestTransferCreate_WhenSourceStockMissing_ReturnsWrappedStockNotFound(t *testing.T) {
	_, db := newTestRepo(t)
	r := &stockTransferRepository{db: db}

	_, err := r.Create(context.Background(), repo.TransferCreateParams{
		EditionUid:         editionUidFor(7),
		SourceWarehouseUid: warehouseUidFor(9),
		TargetWarehouseUid: warehouseUidFor(11),
		Quantity:           3,
	})
	if !errors.Is(err, repo.ErrStockNotFound) {
		t.Fatalf("err = %v, want ErrStockNotFound", err)
	}
}
