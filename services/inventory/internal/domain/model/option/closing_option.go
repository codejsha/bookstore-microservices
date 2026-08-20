package option

import (
	"github.com/codejsha/shared-library-go/pkg/pagination"
)

type ClosingQueryOption struct {
	warehouseUid *string
	year         *int32
	month        *int32
	status       *string
	page         pagination.PageOption
}

type ClosingQueryOptionFunc func(*ClosingQueryOption)

func NewClosingQueryOption(opts ...ClosingQueryOptionFunc) ClosingQueryOption {
	var o ClosingQueryOption
	for _, opt := range opts {
		opt(&o)
	}
	return o
}

func (o ClosingQueryOption) WarehouseUid() *string       { return o.warehouseUid }
func (o ClosingQueryOption) Year() *int32                { return o.year }
func (o ClosingQueryOption) Month() *int32               { return o.month }
func (o ClosingQueryOption) Status() *string             { return o.status }
func (o ClosingQueryOption) Page() pagination.PageOption { return o.page }

func (ClosingQueryOption) WithWarehouseUid(uid *string) ClosingQueryOptionFunc {
	return func(o *ClosingQueryOption) { o.warehouseUid = uid }
}
func (ClosingQueryOption) WithYear(year *int32) ClosingQueryOptionFunc {
	return func(o *ClosingQueryOption) { o.year = year }
}
func (ClosingQueryOption) WithMonth(month *int32) ClosingQueryOptionFunc {
	return func(o *ClosingQueryOption) { o.month = month }
}
func (ClosingQueryOption) WithStatus(status *string) ClosingQueryOptionFunc {
	return func(o *ClosingQueryOption) { o.status = status }
}
func (ClosingQueryOption) WithPage(page pagination.PageOption) ClosingQueryOptionFunc {
	return func(o *ClosingQueryOption) { o.page = page }
}
