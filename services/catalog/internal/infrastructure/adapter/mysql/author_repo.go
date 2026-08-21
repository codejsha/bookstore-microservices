package mysql

import (
	"context"

	"github.com/codejsha/bookstore-microservices/catalog/internal/application/port/repo"
	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/aggregate/entity"
	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/model/condition"
	"github.com/codejsha/bookstore-microservices/commonlib-go/pkg/database"
)

var _ repo.AuthorRepo = (*authorRepository)(nil)

func NewAuthorRepository(dataSource *database.DataSource) repo.AuthorRepo {
	return &authorRepository{
		base: database.NewBaseRepository[entity.AuthorEntity](dataSource),
	}
}

type authorRepository struct {
	base *database.BaseRepository[entity.AuthorEntity]
}

func (r authorRepository) FindAll(ctx context.Context, cond condition.AuthorCondition) (int64, []entity.AuthorEntity, error) {
	total, items, err := r.base.FindAll(ctx, cond.BuildCondition())
	if err != nil {
		return 0, nil, err
	}
	return total, items, nil
}

func (r authorRepository) Find(ctx context.Context, id int64) (entity.AuthorEntity, error) {
	item, err := r.base.Find(ctx, id)
	if err != nil {
		return entity.AuthorEntity{}, err
	}
	return item, nil
}

func (r authorRepository) Create(ctx context.Context, ent entity.AuthorEntity) (entity.AuthorEntity, error) {
	err := r.base.Create(ctx, &ent)
	if err != nil {
		return entity.AuthorEntity{}, err
	}
	return ent, nil
}

func (r authorRepository) Update(ctx context.Context, ent entity.AuthorEntity) (entity.AuthorEntity, error) {
	err := r.base.Save(ctx, &ent)
	if err != nil {
		return entity.AuthorEntity{}, err
	}
	return ent, nil
}

func (r authorRepository) Delete(ctx context.Context, id int64) error {
	return r.base.Delete(ctx, id)
}
