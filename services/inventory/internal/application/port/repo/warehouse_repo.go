package repo

import (
	"context"

	"github.com/codejsha/bookstore-microservices/inventory/internal/domain/aggregate/entity"
	"github.com/codejsha/bookstore-microservices/inventory/internal/domain/model/condition"
)

type WarehouseRepo interface {
	FindAll(ctx context.Context, conds condition.WarehouseCondition) (int64, []entity.WarehouseEntity, error)
	Find(ctx context.Context, id int64) (entity.WarehouseEntity, error)
	Create(ctx context.Context, ent entity.WarehouseEntity) (entity.WarehouseEntity, error)
	Update(ctx context.Context, ent entity.WarehouseEntity) (entity.WarehouseEntity, error)
	Delete(ctx context.Context, id int64) error
}
