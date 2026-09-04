package restcontroller

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/codejsha/shared-library-go/pkg/pagination"

	"github.com/codejsha/bookstore-microservices/inventory/generated/application/port/openapi"
	"github.com/codejsha/bookstore-microservices/inventory/internal/application/usecase"
	"github.com/codejsha/bookstore-microservices/inventory/internal/domain/aggregate"
	"github.com/codejsha/bookstore-microservices/inventory/internal/domain/model/command"
	"github.com/codejsha/bookstore-microservices/inventory/internal/domain/model/option"
	"github.com/codejsha/bookstore-microservices/inventory/internal/infrastructure/httpx"
)

var _ openapi.StockApi = (*stockController)(nil)

type stockController struct {
	inventoryUseCase usecase.InventoryUseCase
}

func NewStockController(inventoryUseCase usecase.InventoryUseCase) openapi.StockApi {
	return &stockController{inventoryUseCase: inventoryUseCase}
}

func (c *stockController) StocksGetAll(
	ctx context.Context,
	editionUid *string,
	warehouseUid *string,
	size *int32,
	page *int32,
	sort *string,
) (*openapi.StockFindAllResponse, error) {
	if err := optionalUidParam(ctx, "edition_uid", editionUid); err != nil {
		return nil, err
	}
	if err := optionalUidParam(ctx, "warehouse_uid", warehouseUid); err != nil {
		return nil, err
	}

	opt := option.NewStockQueryOption(
		option.StockQueryOption{}.WithEditionUid(editionUid),
		option.StockQueryOption{}.WithWarehouseUid(warehouseUid),
		option.StockQueryOption{}.WithPage(pagination.NewPageOption(size, page, sort)),
	)

	total, stocks, err := c.inventoryUseCase.FindAllStocks(ctx, opt)
	if err != nil {
		return nil, err
	}

	items := make([]openapi.StockFindResponse, len(stocks))
	for i, s := range stocks {
		items[i] = toStockFindResponse(s)
	}
	return &openapi.StockFindAllResponse{Items: items, Total: total}, nil
}

func (c *stockController) StocksReceive(ctx context.Context, req openapi.StockReceiveRequest) (*openapi.StockFindResponse, error) {
	cmd := command.StockReceiveCommand{
		EditionUid:   req.EditionUid,
		WarehouseUid: req.WarehouseUid,
		Quantity:     req.Quantity,
		Reason:       req.Reason,
	}

	stock, err := c.inventoryUseCase.ReceiveStock(ctx, cmd)
	if err != nil {
		return nil, httpx.MapBusinessError(ctx, err)
	}

	go runSideEffects(context.WithoutCancel(ctx), "stock", stock.Uid, "received", logrus.Fields{
		"edition_uid":   req.EditionUid,
		"warehouse_uid": req.WarehouseUid,
		"quantity":      req.Quantity,
	})
	resp := toStockFindResponse(stock)
	return &resp, nil
}

func (c *stockController) StocksRelease(ctx context.Context, req openapi.StockReleaseRequest) (*openapi.StockFindResponse, error) {
	cmd := command.StockReleaseCommand{
		EditionUid:   req.EditionUid,
		WarehouseUid: req.WarehouseUid,
		Quantity:     req.Quantity,
		Reason:       req.Reason,
	}

	stock, err := c.inventoryUseCase.ReleaseStock(ctx, cmd)
	if err != nil {
		return nil, httpx.MapBusinessError(ctx, err)
	}

	go runSideEffects(context.WithoutCancel(ctx), "stock", stock.Uid, "released", logrus.Fields{
		"edition_uid":   req.EditionUid,
		"warehouse_uid": req.WarehouseUid,
		"quantity":      req.Quantity,
	})
	resp := toStockFindResponse(stock)
	return &resp, nil
}

func (c *stockController) StocksAdjust(ctx context.Context, req openapi.StockAdjustRequest) (*openapi.StockFindResponse, error) {
	cmd := command.StockAdjustCommand{
		EditionUid:   req.EditionUid,
		WarehouseUid: req.WarehouseUid,
		Quantity:     req.Quantity,
		Reason:       req.Reason,
	}

	stock, err := c.inventoryUseCase.AdjustStock(ctx, cmd)
	if err != nil {
		return nil, httpx.MapBusinessError(ctx, err)
	}

	go runSideEffects(context.WithoutCancel(ctx), "stock", stock.Uid, "adjusted", logrus.Fields{
		"edition_uid":   req.EditionUid,
		"warehouse_uid": req.WarehouseUid,
		"quantity":      req.Quantity,
		"reason":        req.Reason,
	})
	resp := toStockFindResponse(stock)
	return &resp, nil
}

func (c *stockController) StocksReserve(ctx context.Context, req openapi.StockReserveRequest) (*openapi.StockFindResponse, error) {
	cmd := command.StockReserveCommand{
		EditionUid:   req.EditionUid,
		WarehouseUid: req.WarehouseUid,
		Quantity:     req.Quantity,
		Reason:       req.Reason,
	}

	stock, err := c.inventoryUseCase.ReserveStock(ctx, cmd)
	if err != nil {
		return nil, httpx.MapBusinessError(ctx, err)
	}

	go runSideEffects(context.WithoutCancel(ctx), "stock", stock.Uid, "reserved", logrus.Fields{
		"edition_uid":   req.EditionUid,
		"warehouse_uid": req.WarehouseUid,
		"quantity":      req.Quantity,
	})
	resp := toStockFindResponse(stock)
	return &resp, nil
}

func (c *stockController) StocksRead(ctx context.Context, editionUid string) (*openapi.StockFindResponse, error) {
	if err := requireUidParam(ctx, "edition_uid", editionUid); err != nil {
		return nil, err
	}

	stock, err := c.inventoryUseCase.FindStock(ctx, editionUid)
	if err != nil {
		return nil, httpx.MapNotFound(ctx, err)
	}
	if stock == nil {
		return nil, httpx.MapNotFound(ctx, httpx.ErrNotFound)
	}
	resp := toStockFindResponse(stock)
	return &resp, nil
}

func (c *stockController) StocksHistory(
	ctx context.Context,
	editionUid string,
	warehouseUid *string,
	size *int32,
	page *int32,
	sort *string,
) (*openapi.StockHistoryFindAllResponse, error) {
	if err := requireUidParam(ctx, "edition_uid", editionUid); err != nil {
		return nil, err
	}
	if err := optionalUidParam(ctx, "warehouse_uid", warehouseUid); err != nil {
		return nil, err
	}

	total, entries, err := c.inventoryUseCase.FindStockHistory(
		ctx,
		editionUid,
		warehouseUid,
		pagination.NewPageOption(size, page, sort),
	)
	if err != nil {
		return nil, err
	}

	items := make([]openapi.StockHistoryFindResponse, len(entries))
	for i, e := range entries {
		items[i] = toStockHistoryItem(e)
	}
	return &openapi.StockHistoryFindAllResponse{Items: items, Total: total}, nil
}

// ─── mapping helpers ────────────────────────────────────────────────────────

func toStockFindResponse(s *aggregate.StockAggregate) openapi.StockFindResponse {
	warehouses := make([]openapi.StockWarehouseItem, len(s.Warehouses))
	for i, w := range s.Warehouses {
		warehouses[i] = openapi.StockWarehouseItem{
			WarehouseUid:  w.WarehouseUid,
			WarehouseName: w.WarehouseName,
			Quantity:      w.Quantity,
		}
	}
	return openapi.StockFindResponse{
		Uid:           s.Uid,
		EditionUid:    s.EditionUid,
		TotalQuantity: s.TotalQuantity,
		Warehouses:    warehouses,
	}
}

func toStockHistoryItem(e *aggregate.StockHistoryEntry) openapi.StockHistoryFindResponse {
	return openapi.StockHistoryFindResponse{
		Uid:        e.Uid,
		ChangeType: openapi.StockChangeType(e.ChangeType),
		Reason:     e.Reason,
		ChangeQty:  e.ChangeQty,
		CreatedAt:  e.CreatedAt,
	}
}
