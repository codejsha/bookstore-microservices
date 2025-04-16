package repo

import (
	"context"

	"github.com/codejsha/bookstore-microservices/customer/internal/domain/aggregate/entity"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/model/condition"
)

type WishlistRepo interface {
	FindAll(ctx context.Context, cond condition.WishlistCondition) (int64, []entity.WishlistEntity, error)
	Find(ctx context.Context, id int64) (entity.WishlistEntity, error)
	Create(ctx context.Context, ent entity.WishlistEntity) (entity.WishlistEntity, error)
	Update(ctx context.Context, ent entity.WishlistEntity) (entity.WishlistEntity, error)
	Delete(ctx context.Context, id int64) error
}
