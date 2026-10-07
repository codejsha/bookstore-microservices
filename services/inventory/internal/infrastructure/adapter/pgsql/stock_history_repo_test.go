package pgsql

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/codejsha/shared-library-go/pkg/pagination"

	"github.com/codejsha/bookstore-microservices/inventory/generated/infrastructure/port/dao"
	genrepo "github.com/codejsha/bookstore-microservices/inventory/generated/infrastructure/port/repo"
)

func TestStockHistoryFindAll_secondPage_skipsOneFullPage(t *testing.T) {
	_, db := newTestRepo(t)
	stockUid := uuid.NewString()
	base := time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)
	for i := int32(1); i <= 7; i++ {
		if err := db.Exec(
			"INSERT INTO stock_history (uid, stock_id, stock_uid, change_type, change_qty, created_at, actor, version) VALUES (?, 1, ?, 'INBOUND', ?, ?, 1, 1)",
			uuid.NewString(), stockUid, i, base.Add(time.Duration(i)*time.Minute),
		).Error; err != nil {
			t.Fatalf("seed history: %v", err)
		}
	}
	r := &stockHistoryRepository{q: dao.Use(db), StockHistoryRepo: genrepo.NewStockHistoryRepo(db)}
	size, page := int32(3), int32(1)

	total, got, err := r.FindAll(context.Background(), []string{stockUid}, pagination.NewPageOption(&size, &page, nil))
	if err != nil {
		t.Fatalf("find all: %v", err)
	}

	if total != 7 {
		t.Errorf("total = %d, want 7", total)
	}
	quantities := make([]int32, len(got))
	for i, h := range got {
		quantities[i] = h.ChangeQty
	}
	if want := []int32{4, 3, 2}; !reflect.DeepEqual(quantities, want) {
		t.Fatalf("change quantities = %v, want %v", quantities, want)
	}
}

func TestStockHistoryFindAll_twoStocks_pagesBothHistories(t *testing.T) {
	_, db := newTestRepo(t)
	first, second := uuid.NewString(), uuid.NewString()
	base := time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)
	seed := []struct {
		stockUid string
		qty      int32
	}{{first, 1}, {second, 2}, {first, 3}, {second, 4}, {uuid.NewString(), 99}}
	for i, s := range seed {
		if err := db.Exec(
			"INSERT INTO stock_history (uid, stock_id, stock_uid, change_type, change_qty, created_at, actor, version) VALUES (?, 1, ?, 'INBOUND', ?, ?, 1, 1)",
			uuid.NewString(), s.stockUid, s.qty, base.Add(time.Duration(i)*time.Minute),
		).Error; err != nil {
			t.Fatalf("seed history: %v", err)
		}
	}
	r := &stockHistoryRepository{q: dao.Use(db), StockHistoryRepo: genrepo.NewStockHistoryRepo(db)}
	size, page := int32(3), int32(0)

	total, got, err := r.FindAll(context.Background(), []string{first, second}, pagination.NewPageOption(&size, &page, nil))
	if err != nil {
		t.Fatalf("find all: %v", err)
	}

	if total != 4 {
		t.Errorf("total = %d, want 4", total)
	}
	quantities := make([]int32, len(got))
	for i, h := range got {
		quantities[i] = h.ChangeQty
	}
	if want := []int32{4, 3, 2}; !reflect.DeepEqual(quantities, want) {
		t.Fatalf("change quantities = %v, want %v", quantities, want)
	}
}
