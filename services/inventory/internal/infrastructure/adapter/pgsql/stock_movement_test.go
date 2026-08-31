package pgsql

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/codejsha/bookstore-microservices/inventory/internal/application/port/repo"
)

func TestApplyStockDeltaTx_WhenStockMissing_ReturnsWrappedNotFound(t *testing.T) {
	_, db := newTestRepo(t)

	err := applyStockDeltaTx(db, stockDelta{
		editionID:    7,
		editionUID:   editionUidFor(7),
		warehouseID:  9,
		warehouseUID: warehouseUidFor(9),
		delta:        -3,
		changeType:   changeTypeOutbound,
	})
	if !errors.Is(err, repo.ErrStockNotFound) {
		t.Fatalf("err = %v, want ErrStockNotFound", err)
	}
	if !strings.Contains(err.Error(), editionUidFor(7)) || !strings.Contains(err.Error(), warehouseUidFor(9)) {
		t.Errorf("err = %q, want descriptive edition/warehouse context", err)
	}
}

func TestApplyStockDeltaTx_WhenDeltaNegativeAndRowMissing_ReturnsWrappedNotFound(t *testing.T) {
	_, db := newTestRepo(t)

	err := applyStockDeltaTx(db, stockDelta{
		editionID:    7,
		editionUID:   editionUidFor(7),
		warehouseID:  9,
		warehouseUID: warehouseUidFor(9),
		delta:        -2,
		changeType:   changeTypeAdjustment,
		allowCreate:  true,
	})
	if !errors.Is(err, repo.ErrStockNotFound) {
		t.Fatalf("err = %v, want ErrStockNotFound", err)
	}
}

func TestApplyStockDeltaTx_WhenDeltaInboundAndCreateAllowed_CreatesRow(t *testing.T) {
	_, db := newTestRepo(t)

	if err := applyStockDeltaTx(db, stockDelta{
		editionID:    7,
		editionUID:   editionUidFor(7),
		warehouseID:  9,
		warehouseUID: warehouseUidFor(9),
		delta:        4,
		changeType:   changeTypeInbound,
		allowCreate:  true,
	}); err != nil {
		t.Fatalf("err: %v", err)
	}
	if qty := stockQty(t, db, 7, 9); qty != 4 {
		t.Errorf("quantity = %d, want 4", qty)
	}
}

func TestApplyChange_WhenStockMissing_ReturnsWrappedNotFound(t *testing.T) {
	_, db := newTestRepo(t)
	r := stockRepository{db: db}

	err := r.ApplyChange(context.Background(), repo.ApplyChangeParams{
		EditionUid:   editionUidFor(7),
		WarehouseUid: warehouseUidFor(9),
		Delta:        5,
		ChangeType:   changeTypeInbound,
	})
	if !errors.Is(err, repo.ErrStockNotFound) {
		t.Fatalf("err = %v, want ErrStockNotFound", err)
	}
	if !strings.Contains(err.Error(), editionUidFor(7)) {
		t.Errorf("err = %q, want descriptive edition context", err)
	}
}

func TestApplyChange_WhenStockExists_AppliesDelta(t *testing.T) {
	_, db := newTestRepo(t)
	seedStock(t, db, 7, 9, 10)
	r := stockRepository{db: db}

	if err := r.ApplyChange(context.Background(), repo.ApplyChangeParams{
		EditionUid:   editionUidFor(7),
		WarehouseUid: warehouseUidFor(9),
		Delta:        -4,
		ChangeType:   changeTypeOutbound,
	}); err != nil {
		t.Fatalf("err: %v", err)
	}
	if qty := stockQty(t, db, 7, 9); qty != 6 {
		t.Errorf("quantity = %d, want 6", qty)
	}
}
