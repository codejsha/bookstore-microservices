package mysql

import (
	"context"

	"github.com/codejsha/bookstore-microservices/commonlib-go/pkg/database"
	"github.com/codejsha/bookstore-microservices/customer/internal/application/port/repo"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/aggregate/entity"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/model/condition"
)

var _ repo.PointRepo = (*pointRepository)(nil)

func NewPointRepository(
	dataSource *database.DataSource,
) repo.PointRepo {
	return &pointRepository{
		base: database.NewBaseRepository[entity.PointEntity](dataSource),
	}
}

type pointRepository struct {
	base *database.BaseRepository[entity.PointEntity]
}

func (r pointRepository) FindAll(ctx context.Context, cond condition.PointCondition) (int64, []entity.PointEntity, error) {
	total, items, err := r.base.FindAll(ctx, cond.BuildCondition())
	if err != nil {
		return 0, nil, err
	}
	return total, items, nil
}

func (r pointRepository) Find(ctx context.Context, id int64) (entity.PointEntity, error) {
	item, err := r.base.Find(ctx, id)
	if err != nil {
		return entity.PointEntity{}, err
	}
	return item, nil
}

func (r pointRepository) Update(ctx context.Context, ent entity.PointEntity) (entity.PointEntity, error) {
	err := r.base.Save(ctx, &ent)
	if err != nil {
		return entity.PointEntity{}, err
	}
	return ent, nil
}
