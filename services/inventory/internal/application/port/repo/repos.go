package repo

import (
	"context"

	"github.com/codejsha/shared-library-go/pkg/pagination"

	"github.com/codejsha/bookstore-microservices/inventory/internal/domain/model/option"
)

type StockRepo interface {
	FindAll(ctx context.Context, opt option.StockQueryOption) (int64, []*StockResult, error)
	FindOne(ctx context.Context, id int64) (*StockResult, error)
	ApplyChange(ctx context.Context, p ApplyChangeParams) error
}

type StockReservationRepo interface {
	Reserve(ctx context.Context, p ReserveParams) (*ReservationResult, error)
	Release(ctx context.Context, p ReleaseParams) (*ReservationResult, error)
}

type WarehouseRepo interface {
	FindAll(ctx context.Context, opt option.WarehouseQueryOption) (int64, []*WarehouseResult, error)
	FindByUid(ctx context.Context, uid string) (*WarehouseResult, error)
	Create(ctx context.Context, p WarehouseCreateParams) (*WarehouseResult, error)
	Update(ctx context.Context, p WarehouseUpdateParams) (*WarehouseResult, error)
}

type StockHistoryRepo interface {
	FindAll(ctx context.Context, stockUid string, page pagination.PageOption) (int64, []*StockHistoryResult, error)
}

type StockTransferRepo interface {
	FindAll(ctx context.Context, opt option.TransferQueryOption) (int64, []*TransferResult, error)
	FindByUid(ctx context.Context, uid string) (*TransferResult, error)
	Create(ctx context.Context, p TransferCreateParams) (*TransferResult, error)
	Complete(ctx context.Context, uid string) (*TransferResult, error)
	Cancel(ctx context.Context, uid string) (*TransferResult, error)
}

type StockAuditRepo interface {
	FindAll(ctx context.Context, opt option.AuditQueryOption) (int64, []*AuditResult, error)
	FindByUid(ctx context.Context, uid string) (*AuditResult, error)
	Create(ctx context.Context, p AuditCreateParams) (*AuditResult, error)
	Complete(ctx context.Context, uid string) (*AuditResult, error)
}

type MonthlyClosingRepo interface {
	FindAll(ctx context.Context, opt option.ClosingQueryOption) (int64, []*ClosingResult, error)
	FindByUid(ctx context.Context, uid string) (*ClosingResult, error)
	Create(ctx context.Context, p ClosingCreateParams) (*ClosingResult, error)
}

type StockBalanceRepo interface {
	FindAll(ctx context.Context, opt option.BalanceQueryOption) (int64, []*BalanceResult, error)
}
