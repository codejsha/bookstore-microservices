package mysql

import (
	"context"

	"github.com/codejsha/bookstore-microservices/catalog/internal/application/port/repo"
	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/aggregate/entity"
	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/model/condition"
	"github.com/codejsha/bookstore-microservices/commonlib-go/pkg/database"
)

var _ repo.PublisherRepo = (*publisherRepository)(nil)

func NewPublisherRepository(dataSource *database.DataSource) repo.PublisherRepo {
	return &publisherRepository{
		base: database.NewBaseRepository[entity.PublisherEntity](dataSource),
	}
}

type publisherRepository struct {
	base *database.BaseRepository[entity.PublisherEntity]
}

func (r publisherRepository) FindAll(ctx context.Context, cond condition.PublisherCondition) (int64, []entity.PublisherEntity, error) {
	total, items, err := r.base.FindAll(ctx, cond.BuildCondition())
	if err != nil {
		return 0, nil, err
	}
	return total, items, nil
}

func (r publisherRepository) Find(ctx context.Context, id int64) (entity.PublisherEntity, error) {
	item, err := r.base.Find(ctx, id)
	if err != nil {
		return entity.PublisherEntity{}, err
	}
	return item, nil
}

func (r publisherRepository) Create(ctx context.Context, ent entity.PublisherEntity) (entity.PublisherEntity, error) {
	err := r.base.Create(ctx, &ent)
	if err != nil {
		return entity.PublisherEntity{}, err
	}
	return ent, nil
}

func (r publisherRepository) Update(ctx context.Context, ent entity.PublisherEntity) (entity.PublisherEntity, error) {
	err := r.base.Save(ctx, &ent)
	if err != nil {
		return entity.PublisherEntity{}, err
	}
	return ent, nil
}

func (r publisherRepository) Delete(ctx context.Context, id int64) error {
	return r.base.Delete(ctx, id)
}
