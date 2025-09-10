package option

import (
	"github.com/codejsha/shared-library-go/pkg/pagination"
)

type StockQueryOption struct {
	editionUid   *string
	warehouseUid *string
	page         pagination.PageOption
}

type StockQueryOptionFunc func(*StockQueryOption)

func NewStockQueryOption(opts ...StockQueryOptionFunc) StockQueryOption {
	var o StockQueryOption
	for _, opt := range opts {
		opt(&o)
	}
	return o
}

func (o StockQueryOption) EditionUid() *string {
	return o.editionUid
}

func (o StockQueryOption) WarehouseUid() *string {
	return o.warehouseUid
}

func (o StockQueryOption) Page() pagination.PageOption {
	return o.page
}

func (StockQueryOption) WithEditionUid(editionUid *string) StockQueryOptionFunc {
	return func(o *StockQueryOption) {
		o.editionUid = editionUid
	}
}

func (StockQueryOption) WithWarehouseUid(warehouseUid *string) StockQueryOptionFunc {
	return func(o *StockQueryOption) {
		o.warehouseUid = warehouseUid
	}
}

func (StockQueryOption) WithPage(page pagination.PageOption) StockQueryOptionFunc {
	return func(o *StockQueryOption) {
		o.page = page
	}
}
