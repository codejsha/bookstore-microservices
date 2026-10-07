package pgsql

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"testing"
	"time"

	"gorm.io/gorm"

	"github.com/codejsha/shared-library-go/pkg/pagination"

	"github.com/codejsha/bookstore-microservices/inventory/generated/infrastructure/port/entity"
	"github.com/codejsha/bookstore-microservices/inventory/internal/application/port/repo"
	"github.com/codejsha/bookstore-microservices/inventory/internal/domain/model/option"
)

func TestStockFindAll_secondPage_skipsOneFullPageOfEditions(t *testing.T) {
	_, db := newTestRepo(t)
	editionUids := make([]string, 0, 5)
	for editionId := int64(1); editionId <= 5; editionId++ {
		seedStock(t, db, editionId, 1, 10)
		editionUids = append(editionUids, editionUidFor(editionId))
	}
	sort.Strings(editionUids)
	r := stockRepository{db: db}
	size, page := int32(2), int32(1)

	total, got, err := r.FindAll(context.Background(), option.NewStockQueryOption(
		option.StockQueryOption{}.WithPage(pagination.NewPageOption(&size, &page, nil)),
	))
	if err != nil {
		t.Fatalf("find all: %v", err)
	}

	if total != 5 {
		t.Errorf("total = %d, want 5", total)
	}
	gotUids := make([]string, len(got))
	for i, s := range got {
		gotUids[i] = s.EditionUid
	}
	if want := editionUids[2:4]; !reflect.DeepEqual(gotUids, want) {
		t.Fatalf("edition uids = %v, want %v", gotUids, want)
	}
}

func seedWarehouse(t *testing.T, db *gorm.DB, id int64) {
	t.Helper()
	if err := db.Exec(
		"INSERT INTO warehouse (id, uid, name, address, capacity, created_at, actor, version) VALUES (?, ?, ?, 'addr', 100, ?, 1, 1)",
		id, warehouseUidFor(id), fmt.Sprintf("wh-%d", id), time.Now().UTC(),
	).Error; err != nil {
		t.Fatalf("seed warehouse: %v", err)
	}
}

func TestStockApplyChange_receiveIntoEmptyWarehouse_createsRowWithReceivedQuantity(t *testing.T) {
	_, db := newTestRepo(t)
	seedStock(t, db, 7, 1, 10)
	seedWarehouse(t, db, 2)
	r := stockRepository{db: db}

	err := r.ApplyChange(context.Background(), repo.ApplyChangeParams{
		EditionUid:   editionUidFor(7),
		WarehouseUid: warehouseUidFor(2),
		Delta:        5,
		ChangeType:   changeTypeInbound,
	})
	if err != nil {
		t.Fatalf("apply change: %v", err)
	}

	var created entity.StockEntity
	if err := db.Where("edition_uid = ? AND warehouse_uid = ?", editionUidFor(7), warehouseUidFor(2)).First(&created).Error; err != nil {
		t.Fatalf("read created stock: %v", err)
	}
	if created.Quantity != 5 || created.EditionId != 7 || created.WarehouseId != 2 {
		t.Errorf("created = qty %d edition %d warehouse %d, want 5/7/2", created.Quantity, created.EditionId, created.WarehouseId)
	}
	var history int64
	if err := db.Model(&entity.StockHistoryEntity{}).Where("stock_uid = ? AND change_type = ? AND change_qty = 5", created.Uid, changeTypeInbound).Count(&history).Error; err != nil {
		t.Fatalf("count history: %v", err)
	}
	if history != 1 {
		t.Errorf("history rows = %d, want 1", history)
	}
}

func TestStockApplyChange_receiveIntoMissingWarehouse_warehouseNotFound(t *testing.T) {
	_, db := newTestRepo(t)
	seedStock(t, db, 7, 1, 10)
	r := stockRepository{db: db}

	err := r.ApplyChange(context.Background(), repo.ApplyChangeParams{
		EditionUid:   editionUidFor(7),
		WarehouseUid: warehouseUidFor(9),
		Delta:        5,
		ChangeType:   changeTypeInbound,
	})
	if !errors.Is(err, repo.ErrWarehouseNotFound) {
		t.Fatalf("err = %v, want ErrWarehouseNotFound", err)
	}
}

func TestStockApplyChange_releaseIntoMissingRow_stockNotFound(t *testing.T) {
	_, db := newTestRepo(t)
	seedStock(t, db, 7, 1, 10)
	seedWarehouse(t, db, 2)
	r := stockRepository{db: db}

	err := r.ApplyChange(context.Background(), repo.ApplyChangeParams{
		EditionUid:   editionUidFor(7),
		WarehouseUid: warehouseUidFor(2),
		Delta:        5,
		ChangeType:   "RELEASE",
	})
	if !errors.Is(err, repo.ErrStockNotFound) {
		t.Fatalf("err = %v, want ErrStockNotFound", err)
	}
}
