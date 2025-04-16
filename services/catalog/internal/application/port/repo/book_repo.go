package repo

import (
	"context"

	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/aggregate/entity"
	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/model/condition"
)

type BookRepo interface {
	FindAll(ctx context.Context, cond condition.BookCondition) (int64, []entity.BookEntity, error)
	Find(ctx context.Context, id int64) (entity.BookEntity, error)
	Create(ctx context.Context, ent entity.BookEntity) (entity.BookEntity, error)
	Update(ctx context.Context, ent entity.BookEntity) (entity.BookEntity, error)
	Delete(ctx context.Context, id int64) error
}
