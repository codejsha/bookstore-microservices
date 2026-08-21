package service

import (
	"context"

	"github.com/codejsha/bookstore-microservices/inventory/internal/application/usecase"
	"github.com/codejsha/bookstore-microservices/inventory/internal/domain/aggregate"
	"github.com/codejsha/bookstore-microservices/inventory/internal/domain/handler"
	"github.com/codejsha/bookstore-microservices/inventory/internal/domain/model/condition"
)

var _ usecase.StockUseCase = (*stockService)(nil)

func NewStockService(
	stockHandler handler.StockHandler,
	warehouseHandler handler.WarehouseHandler,
) usecase.StockUseCase {
	return &stockService{
		stockHandler:     stockHandler,
		warehouseHandler: warehouseHandler,
	}
}

type stockService struct {
	stockHandler     handler.StockHandler
	warehouseHandler handler.WarehouseHandler
}

func (s stockService) FindAllStocks(ctx context.Context, cond condition.StockCondition) (int64, []*aggregate.StockAggregate, error) {
	// find stocks
	total, stocks, err := s.stockHandler.HandleFindAll(ctx, usecase.StockFindAllQuery{Cond: cond})
	if err != nil {
		return 0, nil, err
	}
	stockAggs := make([]*aggregate.StockAggregate, len(stocks))
	for i := range stocks {
		warehouse, err := s.warehouseHandler.HandleFind(ctx, usecase.WarehouseFindQuery{Id: stocks[i].WarehouseId})
		if err != nil {
			return 0, nil, err
		}
		stockAggs[i] = aggregate.NewStockAggregate(stocks[i], warehouse)
	}
	return total, stockAggs, nil

}

func (s stockService) FindStock(ctx context.Context, id int64) (*aggregate.StockAggregate, error) {
	// TODO implement me
	panic("implement me")
}
