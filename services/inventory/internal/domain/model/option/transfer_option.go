package option

import (
	"github.com/codejsha/shared-library-go/pkg/pagination"
)

type TransferQueryOption struct {
	editionUid         *string
	sourceWarehouseUid *string
	targetWarehouseUid *string
	status             *string
	page               pagination.PageOption
}

type TransferQueryOptionFunc func(*TransferQueryOption)

func NewTransferQueryOption(opts ...TransferQueryOptionFunc) TransferQueryOption {
	var o TransferQueryOption
	for _, opt := range opts {
		opt(&o)
	}
	return o
}

func (o TransferQueryOption) EditionUid() *string         { return o.editionUid }
func (o TransferQueryOption) SourceWarehouseUid() *string { return o.sourceWarehouseUid }
func (o TransferQueryOption) TargetWarehouseUid() *string { return o.targetWarehouseUid }
func (o TransferQueryOption) Status() *string             { return o.status }
func (o TransferQueryOption) Page() pagination.PageOption { return o.page }

func (TransferQueryOption) WithEditionUid(uid *string) TransferQueryOptionFunc {
	return func(o *TransferQueryOption) { o.editionUid = uid }
}
func (TransferQueryOption) WithSourceWarehouseUid(uid *string) TransferQueryOptionFunc {
	return func(o *TransferQueryOption) { o.sourceWarehouseUid = uid }
}
func (TransferQueryOption) WithTargetWarehouseUid(uid *string) TransferQueryOptionFunc {
	return func(o *TransferQueryOption) { o.targetWarehouseUid = uid }
}
func (TransferQueryOption) WithStatus(status *string) TransferQueryOptionFunc {
	return func(o *TransferQueryOption) { o.status = status }
}
func (TransferQueryOption) WithPage(page pagination.PageOption) TransferQueryOptionFunc {
	return func(o *TransferQueryOption) { o.page = page }
}
