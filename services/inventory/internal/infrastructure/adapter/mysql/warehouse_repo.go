package mysql

import (
	"context"

	"github.com/codejsha/bookstore-microservices/commonlib-go/pkg/database"
	"github.com/codejsha/bookstore-microservices/inventory/internal/application/port/repo"
	"github.com/codejsha/bookstore-microservices/inventory/internal/domain/aggregate/entity"
	"github.com/codejsha/bookstore-microservices/inventory/internal/domain/model/condition"
)

var _ repo.WarehouseRepo = (*warehouseRepository)(nil)

func NewWarehouseRepository(dataSource *database.DataSource) repo.WarehouseRepo {
	return &warehouseRepository{
		base: database.NewBaseRepository[entity.WarehouseEntity](dataSource),
	}
}

type warehouseRepository struct {
	base *database.BaseRepository[entity.WarehouseEntity]
}

func (r warehouseRepository) FindAll(ctx context.Context, conds condition.WarehouseCondition) (int64, []entity.WarehouseEntity, error) {
	total, items, err := r.base.FindAll(ctx, conds.BuildCondition())
	if err != nil {
		return 0, nil, err
	}
	return total, items, nil
}

func (r warehouseRepository) Find(ctx context.Context, id int64) (entity.WarehouseEntity, error) {
	item, err := r.base.Find(ctx, id)
	if err != nil {
		return entity.WarehouseEntity{}, err
	}
	return item, nil
}

func (r warehouseRepository) Create(ctx context.Context, ent entity.WarehouseEntity) (entity.WarehouseEntity, error) {
	err := r.base.Create(ctx, &ent)
	if err != nil {
		return entity.WarehouseEntity{}, err
	}
	return ent, nil
}

func (r warehouseRepository) Update(ctx context.Context, ent entity.WarehouseEntity) (entity.WarehouseEntity, error) {
	err := r.base.Save(ctx, &ent)
	if err != nil {
		return entity.WarehouseEntity{}, err
	}
	return ent, nil
}

func (r warehouseRepository) Delete(ctx context.Context, id int64) error {
	return r.base.Delete(ctx, id)
}
