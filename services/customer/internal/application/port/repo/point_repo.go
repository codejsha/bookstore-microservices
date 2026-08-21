package repo

import (
	"context"

	"github.com/codejsha/bookstore-microservices/customer/internal/domain/aggregate/entity"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/model/condition"
)

type PointRepo interface {
	FindAll(ctx context.Context, cond condition.PointCondition) (int64, []entity.PointEntity, error)
	Find(ctx context.Context, id int64) (entity.PointEntity, error)
	Update(ctx context.Context, ent entity.PointEntity) (entity.PointEntity, error)
}
