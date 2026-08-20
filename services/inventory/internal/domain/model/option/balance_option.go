package option

import (
	"github.com/codejsha/shared-library-go/pkg/pagination"
)

type BalanceQueryOption struct {
	editionUid   *string
	warehouseUid *string
	yearMonth    *string
	page         pagination.PageOption
}

type BalanceQueryOptionFunc func(*BalanceQueryOption)

func NewBalanceQueryOption(opts ...BalanceQueryOptionFunc) BalanceQueryOption {
	var o BalanceQueryOption
	for _, opt := range opts {
		opt(&o)
	}
	return o
}

func (o BalanceQueryOption) EditionUid() *string         { return o.editionUid }
func (o BalanceQueryOption) WarehouseUid() *string       { return o.warehouseUid }
func (o BalanceQueryOption) YearMonth() *string          { return o.yearMonth }
func (o BalanceQueryOption) Page() pagination.PageOption { return o.page }

func (BalanceQueryOption) WithEditionUid(uid *string) BalanceQueryOptionFunc {
	return func(o *BalanceQueryOption) { o.editionUid = uid }
}
func (BalanceQueryOption) WithWarehouseUid(uid *string) BalanceQueryOptionFunc {
	return func(o *BalanceQueryOption) { o.warehouseUid = uid }
}
func (BalanceQueryOption) WithYearMonth(ym *string) BalanceQueryOptionFunc {
	return func(o *BalanceQueryOption) { o.yearMonth = ym }
}
func (BalanceQueryOption) WithPage(page pagination.PageOption) BalanceQueryOptionFunc {
	return func(o *BalanceQueryOption) { o.page = page }
}
