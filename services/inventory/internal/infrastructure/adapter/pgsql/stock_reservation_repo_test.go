package pgsql

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/codejsha/bookstore-microservices/inventory/generated/infrastructure/port/entity"
	"github.com/codejsha/bookstore-microservices/inventory/internal/application/port/repo"
)

var testSchemaDDL = []string{
	`CREATE TABLE warehouse (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		uid TEXT, name TEXT, address TEXT, capacity INTEGER,
		created_at DATETIME, updated_at DATETIME, deleted_at DATETIME,
		actor INTEGER, version INTEGER
	)`,
	`CREATE TABLE stock (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		uid TEXT, edition_id INTEGER, edition_uid TEXT,
		warehouse_id INTEGER, warehouse_uid TEXT, quantity INTEGER,
		created_at DATETIME, updated_at DATETIME, deleted_at DATETIME,
		actor INTEGER, version INTEGER
	)`,
	`CREATE TABLE stock_history (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		uid TEXT, stock_id INTEGER, stock_uid TEXT, change_type TEXT, reason TEXT,
		change_qty INTEGER, created_at DATETIME, updated_at DATETIME, deleted_at DATETIME,
		actor INTEGER, version INTEGER
	)`,
	`CREATE TABLE stock_reservation (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		uid TEXT, reservation_key TEXT, order_ref TEXT,
		edition_id INTEGER, edition_uid TEXT, warehouse_id INTEGER, warehouse_uid TEXT,
		quantity INTEGER, status TEXT, released_at DATETIME,
		created_at DATETIME, updated_at DATETIME, deleted_at DATETIME,
		actor INTEGER, version INTEGER,
		CONSTRAINT uk_stock_reservation_key_warehouse UNIQUE (reservation_key, warehouse_id)
	)`,
}

func newTestRepo(t *testing.T) (*stockReservationRepository, *gorm.DB) {
	t.Helper()
	dsn := "file:" + uuid.NewString() + "?mode=memory&cache=shared"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	for _, ddl := range testSchemaDDL {
		if err := db.Exec(ddl).Error; err != nil {
			t.Fatalf("create schema: %v", err)
		}
	}
	t.Cleanup(func() {
		sqlDB, _ := db.DB()
		_ = sqlDB.Close()
	})
	return &stockReservationRepository{db: db}, db
}

func seedStock(t *testing.T, db *gorm.DB, editionId, warehouseId int64, qty int32) *entity.StockEntity {
	t.Helper()
	s := &entity.StockEntity{
		Uid:          uuid.NewString(),
		EditionId:    editionId,
		EditionUid:   editionUidFor(editionId),
		WarehouseId:  warehouseId,
		WarehouseUid: warehouseUidFor(warehouseId),
		Quantity:     qty,
		CreatedAt:    time.Now(),
		Actor:        1,
		Version:      1,
	}
	if err := db.Create(s).Error; err != nil {
		t.Fatalf("seed stock: %v", err)
	}
	return s
}

func editionUidFor(id int64) string {
	return uuid.NewSHA1(uuid.NameSpaceOID, []byte{'e', byte(id)}).String()
}

func warehouseUidFor(id int64) string {
	return uuid.NewSHA1(uuid.NameSpaceOID, []byte{'w', byte(id)}).String()
}

func stockQty(t *testing.T, db *gorm.DB, editionId, warehouseId int64) int32 {
	t.Helper()
	var s entity.StockEntity
	if err := db.Where("edition_id = ? AND warehouse_id = ?", editionId, warehouseId).First(&s).Error; err != nil {
		t.Fatalf("read stock: %v", err)
	}
	return s.Quantity
}

func countReservations(t *testing.T, db *gorm.DB, key string) int64 {
	t.Helper()
	var n int64
	if err := db.Model(&stockReservationEntity{}).Where("reservation_key = ?", key).Count(&n).Error; err != nil {
		t.Fatalf("count reservations: %v", err)
	}
	return n
}

func reserveParams(key string, editionId int64, qty int32) repo.ReserveParams {
	reason := "test"
	return repo.ReserveParams{
		ReservationKey: key,
		OrderRef:       "order-1",
		EditionId:      editionId,
		Quantity:       qty,
		Reason:         &reason,
	}
}

func TestReserve_WhenNoWarehouseHoldsEnough_SplitsAcrossWarehouses(t *testing.T) {
	r, db := newTestRepo(t)
	seedStock(t, db, 100, 1, 3)
	seedStock(t, db, 100, 2, 4)

	got, err := r.Reserve(context.Background(), reserveParams("k1", 100, 6))
	if err != nil {
		t.Fatalf("reserve: %v", err)
	}
	if got.Quantity != 6 {
		t.Errorf("result quantity = %d, want 6 (sum of slices)", got.Quantity)
	}
	if n := countReservations(t, db, "k1"); n != 2 {
		t.Errorf("reservation rows = %d, want 2 (one per warehouse slice)", n)
	}
	if q := stockQty(t, db, 100, 2); q != 0 {
		t.Errorf("wh2 remaining = %d, want 0", q)
	}
	if q := stockQty(t, db, 100, 1); q != 1 {
		t.Errorf("wh1 remaining = %d, want 1", q)
	}
}

func TestReserve_WhenOneWarehouseFitsExactly_DoesNotSplit(t *testing.T) {
	r, db := newTestRepo(t)
	seedStock(t, db, 100, 1, 2)
	seedStock(t, db, 100, 2, 5)

	got, err := r.Reserve(context.Background(), reserveParams("k1", 100, 5))
	if err != nil {
		t.Fatalf("reserve: %v", err)
	}
	if got.Quantity != 5 {
		t.Errorf("result quantity = %d, want 5", got.Quantity)
	}
	if n := countReservations(t, db, "k1"); n != 1 {
		t.Errorf("reservation rows = %d, want 1 (single warehouse covers it)", n)
	}
	if q := stockQty(t, db, 100, 2); q != 0 {
		t.Errorf("wh2 remaining = %d, want 0", q)
	}
	if q := stockQty(t, db, 100, 1); q != 2 {
		t.Errorf("wh1 must be untouched, remaining = %d, want 2", q)
	}
}

func TestReserve_WhenTotalStockInsufficient_WritesNothing(t *testing.T) {
	r, db := newTestRepo(t)
	seedStock(t, db, 100, 1, 3)
	seedStock(t, db, 100, 2, 4)

	_, err := r.Reserve(context.Background(), reserveParams("k1", 100, 10))
	if !errors.Is(err, repo.ErrInsufficientStock) {
		t.Fatalf("err = %v, want ErrInsufficientStock", err)
	}
	if n := countReservations(t, db, "k1"); n != 0 {
		t.Errorf("reservation rows = %d, want 0 (all-or-nothing rollback)", n)
	}
	if q := stockQty(t, db, 100, 1); q != 3 {
		t.Errorf("wh1 must be unchanged, got %d", q)
	}
	if q := stockQty(t, db, 100, 2); q != 4 {
		t.Errorf("wh2 must be unchanged, got %d", q)
	}
}

func TestReserve_WhenReservationKeyReplayed_WritesNothingTwice(t *testing.T) {
	r, db := newTestRepo(t)
	seedStock(t, db, 100, 1, 3)
	seedStock(t, db, 100, 2, 4)

	first, err := r.Reserve(context.Background(), reserveParams("k1", 100, 6))
	if err != nil {
		t.Fatalf("first reserve: %v", err)
	}
	if first.AlreadyApplied {
		t.Errorf("first reserve should not be a replay")
	}

	second, err := r.Reserve(context.Background(), reserveParams("k1", 100, 6))
	if err != nil {
		t.Fatalf("replay reserve: %v", err)
	}
	if !second.AlreadyApplied {
		t.Errorf("second reserve must be a replay (AlreadyApplied)")
	}
	if second.Quantity != 6 {
		t.Errorf("replay quantity = %d, want 6", second.Quantity)
	}
	if q := stockQty(t, db, 100, 2); q != 0 {
		t.Errorf("wh2 = %d, want 0 (no double-decrement)", q)
	}
	if q := stockQty(t, db, 100, 1); q != 1 {
		t.Errorf("wh1 = %d, want 1 (no double-decrement)", q)
	}
	if n := countReservations(t, db, "k1"); n != 2 {
		t.Errorf("reservation rows = %d, want 2 (replay adds none)", n)
	}
}

func TestRelease_WhenReservationSpansWarehouses_CreditsEachOfThem(t *testing.T) {
	r, db := newTestRepo(t)
	seedStock(t, db, 100, 1, 3)
	seedStock(t, db, 100, 2, 4)

	if _, err := r.Reserve(context.Background(), reserveParams("k1", 100, 6)); err != nil {
		t.Fatalf("reserve: %v", err)
	}

	reason := "compensate"
	got, err := r.Release(context.Background(), repo.ReleaseParams{ReservationKey: "k1", Reason: &reason})
	if err != nil {
		t.Fatalf("release: %v", err)
	}
	if got == nil || got.Status != reservationStatusReleased {
		t.Fatalf("release result = %+v, want RELEASED", got)
	}
	if q := stockQty(t, db, 100, 1); q != 3 {
		t.Errorf("wh1 after release = %d, want 3 (fully restored)", q)
	}
	if q := stockQty(t, db, 100, 2); q != 4 {
		t.Errorf("wh2 after release = %d, want 4 (fully restored)", q)
	}

	replay, err := r.Release(context.Background(), repo.ReleaseParams{ReservationKey: "k1", Reason: &reason})
	if err != nil {
		t.Fatalf("release replay: %v", err)
	}
	if replay == nil || !replay.AlreadyApplied {
		t.Errorf("release replay = %+v, want AlreadyApplied", replay)
	}
	if q := stockQty(t, db, 100, 1); q != 3 {
		t.Errorf("wh1 after replay = %d, want 3 (no double-credit)", q)
	}
	if q := stockQty(t, db, 100, 2); q != 4 {
		t.Errorf("wh2 after replay = %d, want 4 (no double-credit)", q)
	}
}

func TestRelease_WhenReservationKeyUnknown_WritesNothing(t *testing.T) {
	r, _ := newTestRepo(t)
	got, err := r.Release(context.Background(), repo.ReleaseParams{ReservationKey: "missing"})
	if err != nil {
		t.Fatalf("release: %v", err)
	}
	if got != nil {
		t.Errorf("release of unknown key = %+v, want nil no-op", got)
	}
}

func TestAllocateReservation_WhenWarehousesRanked_AllocatesInThatOrder(t *testing.T) {
	rows := []entity.StockEntity{
		{Id: 1, Quantity: 2},
		{Id: 2, Quantity: 5},
		{Id: 3, Quantity: 5},
	}
	allocs, err := allocateReservation(rows, 100, 8)
	if err != nil {
		t.Fatalf("allocate: %v", err)
	}
	if len(allocs) != 2 {
		t.Fatalf("allocs = %d, want 2", len(allocs))
	}
	if allocs[0].stock.Id != 2 || allocs[0].quantity != 5 {
		t.Errorf("alloc[0] = id %d qty %d, want id 2 qty 5", allocs[0].stock.Id, allocs[0].quantity)
	}
	if allocs[1].stock.Id != 3 || allocs[1].quantity != 3 {
		t.Errorf("alloc[1] = id %d qty %d, want id 3 qty 3", allocs[1].stock.Id, allocs[1].quantity)
	}
}

func TestAllocateReservation_WhenStockInsufficient_ReturnsErrInsufficientStock(t *testing.T) {
	rows := []entity.StockEntity{{Id: 1, Quantity: 2}, {Id: 2, Quantity: 1}}
	if _, err := allocateReservation(rows, 100, 5); !errors.Is(err, repo.ErrInsufficientStock) {
		t.Errorf("err = %v, want ErrInsufficientStock", err)
	}
}
