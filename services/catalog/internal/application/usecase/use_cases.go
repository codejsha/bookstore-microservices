package usecase

import (
	"context"

	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/aggregate"
	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/model/condition"
)

type BookUseCase interface {
	FindAllBooks(ctx context.Context, cond condition.BookCondition) (int64, []*aggregate.BookAggregate, error)
	FindBook(ctx context.Context, id int64) (*aggregate.BookAggregate, error)
}
