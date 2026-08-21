package repo

import (
	"context"

	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/aggregate/entity"
	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/model/condition"
)

type AuthorRepo interface {
	FindAll(ctx context.Context, cond condition.AuthorCondition) (int64, []entity.AuthorEntity, error)
	Find(ctx context.Context, id int64) (entity.AuthorEntity, error)
	Create(ctx context.Context, ent entity.AuthorEntity) (entity.AuthorEntity, error)
	Update(ctx context.Context, ent entity.AuthorEntity) (entity.AuthorEntity, error)
	Delete(ctx context.Context, id int64) error
}
