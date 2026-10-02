package pgsql

import (
	"context"
	"reflect"
	"sort"
	"testing"

	"github.com/codejsha/shared-library-go/pkg/pagination"

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
