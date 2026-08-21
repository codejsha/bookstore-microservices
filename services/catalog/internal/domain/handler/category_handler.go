package handler

import (
	"context"

	"github.com/codejsha/bookstore-microservices/catalog/internal/application/usecase"
	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/constant"
)

var _ CategoryHandler = (*categoryHandlerImpl)(nil)

func NewCategoryHandler() CategoryHandler {
	return &categoryHandlerImpl{}
}

type CategoryHandler interface {
	HandleFindAll(ctx context.Context, query usecase.CategoryFindAllQuery) (int64, []constant.Category, error)
}

type categoryHandlerImpl struct{}

func (c categoryHandlerImpl) HandleFindAll(ctx context.Context, query usecase.CategoryFindAllQuery) (int64, []constant.Category, error) {
	categories := constant.Category{}.FindAllCategories(query.Cond.Sort, query.Cond.Order)
	return int64(len(categories)), categories, nil
}
