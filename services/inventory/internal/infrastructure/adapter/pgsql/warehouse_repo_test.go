package pgsql

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/codejsha/shared-library-go/pkg/pagination"

	"github.com/codejsha/bookstore-microservices/inventory/internal/domain/model/option"
)

func TestWarehouseFindAll_secondPage_skipsOneFullPage(t *testing.T) {
	_, db := newTestRepo(t)
	for i := 1; i <= 7; i++ {
		if err := db.Exec(
			"INSERT INTO warehouse (uid, name, address, capacity, created_at, actor, version) VALUES (?, ?, 'addr', 100, ?, 1, 1)",
			uuid.NewString(), fmt.Sprintf("wh-%d", i), time.Now().UTC(),
		).Error; err != nil {
			t.Fatalf("seed warehouse: %v", err)
		}
	}
	r := warehouseRepository{db: db}
	size, page, sort := int32(3), int32(1), "name,asc"

	total, got, err := r.FindAll(context.Background(), option.NewWarehouseQueryOption(
		option.WarehouseQueryOption{}.WithPage(pagination.NewPageOption(&size, &page, &sort)),
	))
	if err != nil {
		t.Fatalf("find all: %v", err)
	}

	if total != 7 {
		t.Errorf("total = %d, want 7", total)
	}
	names := make([]string, len(got))
	for i, w := range got {
		names[i] = w.Name
		if w.CreatedAt.IsZero() {
			t.Errorf("warehouse %s created_at is zero", w.Name)
		}
	}
	if want := []string{"wh-4", "wh-5", "wh-6"}; !reflect.DeepEqual(names, want) {
		t.Fatalf("names = %v, want %v", names, want)
	}
}

func TestWarehouseFindAll_sortOmitted_newestFirst(t *testing.T) {
	_, db := newTestRepo(t)
	base := time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)
	for i := 1; i <= 3; i++ {
		if err := db.Exec(
			"INSERT INTO warehouse (uid, name, address, capacity, created_at, actor, version) VALUES (?, ?, 'addr', 100, ?, 1, 1)",
			uuid.NewString(), fmt.Sprintf("wh-%d", i), base.Add(time.Duration(i)*time.Hour),
		).Error; err != nil {
			t.Fatalf("seed warehouse: %v", err)
		}
	}
	r := warehouseRepository{db: db}
	size, page := int32(3), int32(0)

	_, got, err := r.FindAll(context.Background(), option.NewWarehouseQueryOption(
		option.WarehouseQueryOption{}.WithPage(pagination.NewPageOption(&size, &page, nil)),
	))
	if err != nil {
		t.Fatalf("find all: %v", err)
	}

	names := make([]string, len(got))
	for i, w := range got {
		names[i] = w.Name
	}
	if want := []string{"wh-3", "wh-2", "wh-1"}; !reflect.DeepEqual(names, want) {
		t.Fatalf("names = %v, want %v", names, want)
	}
}

func seedWarehousesHourly(t *testing.T, db *gorm.DB, n int) {
	t.Helper()
	base := time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)
	for i := 1; i <= n; i++ {
		if err := db.Exec(
			"INSERT INTO warehouse (uid, name, address, capacity, created_at, actor, version) VALUES (?, ?, 'addr', 100, ?, 1, 1)",
			uuid.NewString(), fmt.Sprintf("wh-%d", i), base.Add(time.Duration(i)*time.Hour),
		).Error; err != nil {
			t.Fatalf("seed warehouse: %v", err)
		}
	}
}

func warehouseNamesSorted(t *testing.T, db *gorm.DB, sort string) []string {
	t.Helper()
	r := warehouseRepository{db: db}
	size, page := int32(10), int32(0)
	_, got, err := r.FindAll(context.Background(), option.NewWarehouseQueryOption(
		option.WarehouseQueryOption{}.WithPage(pagination.NewPageOption(&size, &page, &sort)),
	))
	if err != nil {
		t.Fatalf("find all: %v", err)
	}
	names := make([]string, len(got))
	for i, w := range got {
		names[i] = w.Name
	}
	return names
}

func TestWarehouseFindAll_snakeCaseSort_ordersByWireName(t *testing.T) {
	_, db := newTestRepo(t)
	seedWarehousesHourly(t, db, 3)

	if got, want := warehouseNamesSorted(t, db, "created_at,asc"), []string{"wh-1", "wh-2", "wh-3"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("names = %v, want %v", got, want)
	}
}

func TestWarehouseFindAll_camelCaseSort_ordersByAlias(t *testing.T) {
	_, db := newTestRepo(t)
	seedWarehousesHourly(t, db, 3)

	if got, want := warehouseNamesSorted(t, db, "createdAt,asc"), []string{"wh-1", "wh-2", "wh-3"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("names = %v, want %v", got, want)
	}
}

func TestWarehouseFindAll_unknownSortOnly_defaultNewestFirst(t *testing.T) {
	_, db := newTestRepo(t)
	seedWarehousesHourly(t, db, 3)

	if got, want := warehouseNamesSorted(t, db, "bogus,asc;nope"), []string{"wh-3", "wh-2", "wh-1"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("names = %v, want %v", got, want)
	}
}

func TestWarehouseFindAll_unknownMixedWithKnownSort_keepsKnown(t *testing.T) {
	_, db := newTestRepo(t)
	seedWarehousesHourly(t, db, 3)

	if got, want := warehouseNamesSorted(t, db, "bogus,desc;name,asc"), []string{"wh-1", "wh-2", "wh-3"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("names = %v, want %v", got, want)
	}
}
