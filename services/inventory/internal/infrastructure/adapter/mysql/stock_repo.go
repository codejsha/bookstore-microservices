package mysql

import (
	"context"

	"github.com/codejsha/bookstore-microservices/commonlib-go/pkg/database"
	"github.com/codejsha/bookstore-microservices/inventory/internal/application/port/repo"
	"github.com/codejsha/bookstore-microservices/inventory/internal/domain/aggregate/entity"
	"github.com/codejsha/bookstore-microservices/inventory/internal/domain/model/condition"
)

var _ repo.StockRepo = (*stockRepository)(nil)

func NewStockRepository(dataSource *database.DataSource) repo.StockRepo {
	return &stockRepository{
		base: database.NewBaseRepository[entity.StockEntity](dataSource),
	}
}

type stockRepository struct {
	base *database.BaseRepository[entity.StockEntity]
}

func (r stockRepository) FindAll(ctx context.Context, cond condition.StockCondition) (int64, []entity.StockEntity, error) {
	total, items, err := r.base.FindAll(ctx, cond.BuildCondition())
	if err != nil {
		return 0, nil, err
	}
	return total, items, nil
}

func (r stockRepository) Find(ctx context.Context, id int64) (entity.StockEntity, error) {
	item, err := r.base.Find(ctx, id)
	if err != nil {
		return entity.StockEntity{}, err
	}
	return item, nil
}

func (r stockRepository) Create(ctx context.Context, ent entity.StockEntity) (entity.StockEntity, error) {
	err := r.base.Create(ctx, &ent)
	if err != nil {
		return entity.StockEntity{}, err
	}
	return ent, nil
}

func (r stockRepository) Update(ctx context.Context, ent entity.StockEntity) (entity.StockEntity, error) {
	err := r.base.Save(ctx, &ent)
	if err != nil {
		return entity.StockEntity{}, err
	}
	return ent, nil
}

func (r stockRepository) Delete(ctx context.Context, id int64) error {
	return r.base.Delete(ctx, id)
}
