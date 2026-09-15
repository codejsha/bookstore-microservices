package pgsql

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/codejsha/bookstore-microservices/inventory/internal/application/port/repo"
)

var closingSchemaDDL = []string{
	`CREATE TABLE monthly_closing (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		uid TEXT, warehouse_id INTEGER, warehouse_uid TEXT, year INTEGER, month INTEGER, status TEXT,
		closed_at DATETIME, created_at DATETIME, updated_at DATETIME, deleted_at DATETIME,
		actor INTEGER, version INTEGER
	)`,
	`CREATE TABLE monthly_closing_item (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		closing_id INTEGER, edition_id INTEGER, edition_uid TEXT,
		opening_quantity INTEGER, inbound_quantity INTEGER, outbound_quantity INTEGER,
		adjust_quantity INTEGER, closing_quantity INTEGER
	)`,
}

func seedHistory(t *testing.T, db *gorm.DB, stockID int64, changeType string, qty int32, at time.Time) {
	t.Helper()
	err := db.Exec(
		"INSERT INTO stock_history (uid, stock_id, change_type, change_qty, created_at, actor, version) VALUES (?, ?, ?, ?, ?, 1, 1)",
		uuid.NewString(), stockID, changeType, qty, at,
	).Error
	if err != nil {
		t.Fatalf("seed history: %v", err)
	}
}

func countStockHistoryReads(db *gorm.DB) *atomic.Int32 {
	var reads atomic.Int32
	count := func(tx *gorm.DB) {
		if tx.Statement.Table == "stock_history" {
			reads.Add(1)
		}
	}
	_ = db.Callback().Query().After("gorm:query").Register("test:count_history_query", count)
	_ = db.Callback().Row().After("gorm:row").Register("test:count_history_row", count)
	return &reads
}

func TestMonthlyClosingCreate_severalStocks_readsHistoryOnceAndComputesItems(t *testing.T) {
	_, db := newTestRepo(t)
	for _, ddl := range closingSchemaDDL {
		if err := db.Exec(ddl).Error; err != nil {
			t.Fatalf("create closing schema: %v", err)
		}
	}
	const warehouseID = 7
	if err := db.Exec(
		"INSERT INTO warehouse (id, uid, name, created_at, actor, version) VALUES (?, ?, 'main', ?, 1, 1)",
		warehouseID, warehouseUidFor(warehouseID), time.Now().UTC(),
	).Error; err != nil {
		t.Fatalf("seed warehouse: %v", err)
	}

	moving := seedStock(t, db, 101, warehouseID, 50)
	reserved := seedStock(t, db, 102, warehouseID, 20)
	idle := seedStock(t, db, 103, warehouseID, 5)
	elsewhere := seedStock(t, db, 104, 8, 99)

	august := func(day int) time.Time { return time.Date(2026, time.August, day, 12, 0, 0, 0, time.UTC) }
	seedHistory(t, db, moving.Id, "INBOUND", 30, august(10))
	seedHistory(t, db, moving.Id, "OUTBOUND", -10, august(15))
	seedHistory(t, db, moving.Id, "ADJUSTMENT", 2, august(20))
	seedHistory(t, db, moving.Id, "INBOUND", 100, time.Date(2026, time.September, 2, 12, 0, 0, 0, time.UTC))
	seedHistory(t, db, reserved.Id, "RESERVATION", -5, august(5))
	seedHistory(t, db, elsewhere.Id, "INBOUND", 7, august(9))

	reads := countStockHistoryReads(db)
	r := &monthlyClosingRepository{db: db}

	res, err := r.Create(context.Background(), repo.ClosingCreateParams{
		WarehouseUid: warehouseUidFor(warehouseID),
		Year:         2026,
		Month:        8,
	})
	if err != nil {
		t.Fatalf("create closing: %v", err)
	}

	if got := reads.Load(); got != 1 {
		t.Fatalf("stock_history read %d times, want 1 for the whole warehouse", got)
	}

	want := map[string]repo.ClosingItemResult{
		moving.EditionUid:   {EditionUid: moving.EditionUid, OpeningQuantity: 28, InboundQuantity: 30, OutboundQuantity: 10, AdjustQuantity: 2, ClosingQuantity: 50},
		reserved.EditionUid: {EditionUid: reserved.EditionUid, OpeningQuantity: 25, InboundQuantity: 0, OutboundQuantity: 5, AdjustQuantity: 0, ClosingQuantity: 20},
		idle.EditionUid:     {EditionUid: idle.EditionUid, OpeningQuantity: 5, ClosingQuantity: 5},
	}
	if len(res.Items) != len(want) {
		t.Fatalf("items = %d, want %d", len(res.Items), len(want))
	}
	for _, item := range res.Items {
		if item != want[item.EditionUid] {
			t.Fatalf("item %s = %+v, want %+v", item.EditionUid, item, want[item.EditionUid])
		}
	}

	var stored int64
	if err := db.Table("monthly_closing_item").Count(&stored).Error; err != nil {
		t.Fatalf("count stored items: %v", err)
	}
	if stored != int64(len(want)) {
		t.Fatalf("stored items = %d, want %d", stored, len(want))
	}
}
