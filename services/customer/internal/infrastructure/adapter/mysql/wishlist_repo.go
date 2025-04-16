package mysql

import (
	"context"

	"github.com/codejsha/bookstore-microservices/commonlib-go/pkg/database"
	"github.com/codejsha/bookstore-microservices/customer/internal/application/port/repo"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/aggregate/entity"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/model/condition"
)

var _ repo.WishlistRepo = (*wishlistRepository)(nil)

func NewWishlistRepository(
	dataSource *database.DataSource,
) repo.WishlistRepo {
	return &wishlistRepository{
		base: database.NewBaseRepository[entity.WishlistEntity](dataSource),
	}
}

type wishlistRepository struct {
	base *database.BaseRepository[entity.WishlistEntity]
}

func (r wishlistRepository) FindAll(ctx context.Context, cond condition.WishlistCondition) (int64, []entity.WishlistEntity, error) {
	total, items, err := r.base.FindAll(ctx, cond.BuildCondition())
	if err != nil {
		return 0, nil, err
	}
	return total, items, nil
}

func (r wishlistRepository) Find(ctx context.Context, id int64) (entity.WishlistEntity, error) {
	item, err := r.base.Find(ctx, id)
	if err != nil {
		return entity.WishlistEntity{}, err
	}
	return item, nil
}

func (r wishlistRepository) Create(ctx context.Context, ent entity.WishlistEntity) (entity.WishlistEntity, error) {
	err := r.base.Create(ctx, &ent)
	if err != nil {
		return entity.WishlistEntity{}, err
	}
	return ent, nil
}

func (r wishlistRepository) Update(ctx context.Context, ent entity.WishlistEntity) (entity.WishlistEntity, error) {
	err := r.base.Save(ctx, &ent)
	if err != nil {
		return entity.WishlistEntity{}, err
	}
	return ent, nil
}

func (r wishlistRepository) Delete(ctx context.Context, id int64) error {
	return r.base.Delete(ctx, id)
}
