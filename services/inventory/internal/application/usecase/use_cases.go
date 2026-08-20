package usecase

import (
	"context"

	"github.com/codejsha/shared-library-go/pkg/pagination"

	"github.com/codejsha/bookstore-microservices/inventory/internal/domain/aggregate"
	"github.com/codejsha/bookstore-microservices/inventory/internal/domain/model/command"
	"github.com/codejsha/bookstore-microservices/inventory/internal/domain/model/option"
)

type InventoryUseCase interface {
	// ─── Stock query ────────────────────────────────────────────────────
	FindAllStocks(ctx context.Context, opt option.StockQueryOption) (int64, []*aggregate.StockAggregate, error)
	FindStock(ctx context.Context, editionUid string) (*aggregate.StockAggregate, error)
	FindStockHistory(ctx context.Context, editionUid string, warehouseUid *string, page pagination.PageOption) (int64, []*aggregate.StockHistoryEntry, error)

	// ─── Stock operations ───────────────────────────────────────────────
	ReceiveStock(ctx context.Context, cmd command.StockReceiveCommand) (*aggregate.StockAggregate, error)
	ReleaseStock(ctx context.Context, cmd command.StockReleaseCommand) (*aggregate.StockAggregate, error)
	AdjustStock(ctx context.Context, cmd command.StockAdjustCommand) (*aggregate.StockAggregate, error)
	ReserveStock(ctx context.Context, cmd command.StockReserveCommand) (*aggregate.StockAggregate, error)

	// ─── Saga (Temporal) stock operations ───────────────────────────────
	ReserveStockForOrder(ctx context.Context, cmd command.StockOrderReserveCommand) (*aggregate.StockAggregate, error)
	ReleaseStockForOrder(ctx context.Context, cmd command.StockOrderReleaseCommand) (*aggregate.StockAggregate, error)

	// ─── Warehouse management ───────────────────────────────────────────
	FindAllWarehouses(ctx context.Context, opt option.WarehouseQueryOption) (int64, []*aggregate.WarehouseAggregate, error)
	FindWarehouse(ctx context.Context, uid string) (*aggregate.WarehouseAggregate, error)
	RegisterWarehouse(ctx context.Context, cmd command.WarehouseCreateCommand) (*aggregate.WarehouseAggregate, error)
	UpdateWarehouse(ctx context.Context, uid string, cmd command.WarehouseUpdateCommand) (*aggregate.WarehouseAggregate, error)

	// ─── Stock transfer ────────────────────────────────────────────────
	TransferStock(ctx context.Context, cmd command.StockTransferCommand) (*aggregate.StockTransferAggregate, error)
	FindAllTransfers(ctx context.Context, opt option.TransferQueryOption) (int64, []*aggregate.StockTransferAggregate, error)
	FindTransfer(ctx context.Context, uid string) (*aggregate.StockTransferAggregate, error)
	CompleteTransfer(ctx context.Context, uid string) (*aggregate.StockTransferAggregate, error)
	CancelTransfer(ctx context.Context, uid string) (*aggregate.StockTransferAggregate, error)

	// ─── Stock audit ───────────────────────────────────────────────────
	CreateAudit(ctx context.Context, cmd command.StockAuditCreateCommand) (*aggregate.StockAuditAggregate, error)
	FindAllAudits(ctx context.Context, opt option.AuditQueryOption) (int64, []*aggregate.StockAuditAggregate, error)
	FindAudit(ctx context.Context, uid string) (*aggregate.StockAuditAggregate, error)
	CompleteAudit(ctx context.Context, uid string) (*aggregate.StockAuditAggregate, error)

	// ─── Monthly closing ───────────────────────────────────────────────
	CreateMonthlyClosing(ctx context.Context, cmd command.MonthlyClosingCommand) (*aggregate.MonthlyClosingAggregate, error)
	FindAllClosings(ctx context.Context, opt option.ClosingQueryOption) (int64, []*aggregate.MonthlyClosingAggregate, error)
	FindClosing(ctx context.Context, uid string) (*aggregate.MonthlyClosingAggregate, error)

	// ─── Stock balance ─────────────────────────────────────────────────
	FindStockBalance(ctx context.Context, opt option.BalanceQueryOption) (int64, []*aggregate.StockBalanceEntry, error)
}
