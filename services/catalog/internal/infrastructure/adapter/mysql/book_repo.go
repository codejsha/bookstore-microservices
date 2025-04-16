package mysql

import (
	"context"

	"github.com/codejsha/bookstore-microservices/catalog/internal/application/port/repo"
	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/aggregate/entity"
	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/model/condition"
	"github.com/codejsha/bookstore-microservices/commonlib-go/pkg/database"
)

var _ repo.BookRepo = (*bookRepository)(nil)

func NewBookRepository(dataSource *database.DataSource) repo.BookRepo {
	return &bookRepository{
		base: database.NewBaseRepository[entity.BookEntity](dataSource),
	}
}

type bookRepository struct {
	base *database.BaseRepository[entity.BookEntity]
}

func (r bookRepository) FindAll(ctx context.Context, cond condition.BookCondition) (int64, []entity.BookEntity, error) {
	total, items, err := r.base.FindAll(ctx, cond.BuildCondition())
	if err != nil {
		return 0, nil, err
	}
	return total, items, nil
}

func (r bookRepository) Find(ctx context.Context, id int64) (entity.BookEntity, error) {
	item, err := r.base.Find(ctx, id)
	if err != nil {
		return entity.BookEntity{}, err
	}
	return item, nil
}

func (r bookRepository) Create(ctx context.Context, ent entity.BookEntity) (entity.BookEntity, error) {
	err := r.base.Create(ctx, &ent)
	if err != nil {
		return entity.BookEntity{}, err
	}
	return ent, nil
}

func (r bookRepository) Update(ctx context.Context, ent entity.BookEntity) (entity.BookEntity, error) {
	err := r.base.Save(ctx, &ent)
	if err != nil {
		return entity.BookEntity{}, err
	}
	return ent, nil
}

func (r bookRepository) Delete(ctx context.Context, id int64) error {
	return r.base.Delete(ctx, id)
}
