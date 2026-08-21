package repo

import (
	"context"

	"github.com/codejsha/bookstore-microservices/inventory/internal/domain/aggregate/entity"
	"github.com/codejsha/bookstore-microservices/inventory/internal/domain/model/condition"
)

type StockRepo interface {
	FindAll(ctx context.Context, cond condition.StockCondition) (int64, []entity.StockEntity, error)
	Find(ctx context.Context, id int64) (entity.StockEntity, error)
	Create(ctx context.Context, ent entity.StockEntity) (entity.StockEntity, error)
	Update(ctx context.Context, ent entity.StockEntity) (entity.StockEntity, error)
	Delete(ctx context.Context, id int64) error
}
