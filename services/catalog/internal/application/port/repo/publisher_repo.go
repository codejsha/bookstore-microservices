package repo

import (
	"context"

	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/aggregate/entity"
	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/model/condition"
)

type PublisherRepo interface {
	FindAll(ctx context.Context, cond condition.PublisherCondition) (int64, []entity.PublisherEntity, error)
	Find(ctx context.Context, id int64) (entity.PublisherEntity, error)
	Create(ctx context.Context, ent entity.PublisherEntity) (entity.PublisherEntity, error)
	Update(ctx context.Context, ent entity.PublisherEntity) (entity.PublisherEntity, error)
	Delete(ctx context.Context, id int64) error
}
