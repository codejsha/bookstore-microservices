package option

import (
	"github.com/codejsha/shared-library-go/pkg/pagination"
)

type AuditQueryOption struct {
	warehouseUid *string
	status       *string
	page         pagination.PageOption
}

type AuditQueryOptionFunc func(*AuditQueryOption)

func NewAuditQueryOption(opts ...AuditQueryOptionFunc) AuditQueryOption {
	var o AuditQueryOption
	for _, opt := range opts {
		opt(&o)
	}
	return o
}

func (o AuditQueryOption) WarehouseUid() *string       { return o.warehouseUid }
func (o AuditQueryOption) Status() *string             { return o.status }
func (o AuditQueryOption) Page() pagination.PageOption { return o.page }

func (AuditQueryOption) WithWarehouseUid(uid *string) AuditQueryOptionFunc {
	return func(o *AuditQueryOption) { o.warehouseUid = uid }
}
func (AuditQueryOption) WithStatus(status *string) AuditQueryOptionFunc {
	return func(o *AuditQueryOption) { o.status = status }
}
func (AuditQueryOption) WithPage(page pagination.PageOption) AuditQueryOptionFunc {
	return func(o *AuditQueryOption) { o.page = page }
}
