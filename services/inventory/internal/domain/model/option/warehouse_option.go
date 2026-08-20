package option

import (
	"github.com/codejsha/shared-library-go/pkg/pagination"
)

type WarehouseQueryOption struct {
	name *string
	page pagination.PageOption
}

type WarehouseQueryOptionFunc func(*WarehouseQueryOption)

func NewWarehouseQueryOption(opts ...WarehouseQueryOptionFunc) WarehouseQueryOption {
	var o WarehouseQueryOption
	for _, opt := range opts {
		opt(&o)
	}
	return o
}

func (o WarehouseQueryOption) Name() *string {
	return o.name
}

func (o WarehouseQueryOption) Page() pagination.PageOption {
	return o.page
}

func (WarehouseQueryOption) WithName(name *string) WarehouseQueryOptionFunc {
	return func(o *WarehouseQueryOption) {
		o.name = name
	}
}
func (WarehouseQueryOption) WithPage(page pagination.PageOption) WarehouseQueryOptionFunc {
	return func(o *WarehouseQueryOption) {
		o.page = page
	}
}
