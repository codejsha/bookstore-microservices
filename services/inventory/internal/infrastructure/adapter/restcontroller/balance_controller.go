package restcontroller

import (
	"context"

	"github.com/codejsha/shared-library-go/pkg/pagination"

	"github.com/codejsha/bookstore-microservices/inventory/generated/application/port/openapi"
	"github.com/codejsha/bookstore-microservices/inventory/internal/application/usecase"
	"github.com/codejsha/bookstore-microservices/inventory/internal/domain/aggregate"
	"github.com/codejsha/bookstore-microservices/inventory/internal/domain/model/option"
)

var _ openapi.BalanceApi = (*balanceController)(nil)

type balanceController struct {
	inventoryUseCase usecase.InventoryUseCase
}

func NewBalanceController(inventoryUseCase usecase.InventoryUseCase) openapi.BalanceApi {
	return &balanceController{inventoryUseCase: inventoryUseCase}
}

func (c *balanceController) BalanceGetAll(
	ctx context.Context,
	editionUid *string,
	warehouseUid *string,
	yearMonth *string,
	size *int32,
	page *int32,
	sort *string,
) (*openapi.BalanceFindAllResponse, error) {
	if err := optionalUidParam(ctx, "edition_uid", editionUid); err != nil {
		return nil, err
	}
	if err := optionalUidParam(ctx, "warehouse_uid", warehouseUid); err != nil {
		return nil, err
	}

	opt := option.NewBalanceQueryOption(
		option.BalanceQueryOption{}.WithEditionUid(editionUid),
		option.BalanceQueryOption{}.WithWarehouseUid(warehouseUid),
		option.BalanceQueryOption{}.WithYearMonth(yearMonth),
		option.BalanceQueryOption{}.WithPage(pagination.NewPageOption(size, page, sort)),
	)

	total, entries, err := c.inventoryUseCase.FindStockBalance(ctx, opt)
	if err != nil {
		return nil, err
	}

	items := make([]openapi.BalanceFindResponse, len(entries))
	for i, e := range entries {
		items[i] = toBalanceFindResponse(e)
	}
	return &openapi.BalanceFindAllResponse{Items: items, Total: total}, nil
}

// ─── mapping helpers ────────────────────────────────────────────────────────

func toBalanceFindResponse(e *aggregate.StockBalanceEntry) openapi.BalanceFindResponse {
	return openapi.BalanceFindResponse{
		EditionUid:       e.EditionUid,
		WarehouseUid:     e.WarehouseUid,
		WarehouseName:    e.WarehouseName,
		InboundQuantity:  e.InboundQuantity,
		OutboundQuantity: e.OutboundQuantity,
		AdjustQuantity:   e.AdjustQuantity,
		CurrentQuantity:  e.CurrentQuantity,
	}
}
