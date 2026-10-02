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
	}
	if want := []string{"wh-4", "wh-5", "wh-6"}; !reflect.DeepEqual(names, want) {
		t.Fatalf("names = %v, want %v", names, want)
	}
}
