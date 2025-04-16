package usecase

import (
	"context"

	"github.com/codejsha/bookstore-microservices/inventory/internal/domain/aggregate"
	"github.com/codejsha/bookstore-microservices/inventory/internal/domain/model/condition"
)

type StockUseCase interface {
	FindAllStocks(ctx context.Context, cond condition.StockCondition) (int64, []*aggregate.StockAggregate, error)
	FindStock(ctx context.Context, id int64) (*aggregate.StockAggregate, error)
}
