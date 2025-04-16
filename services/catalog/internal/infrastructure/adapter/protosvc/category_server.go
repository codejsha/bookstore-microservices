package protosvc

import (
	"context"

	"github.com/codejsha/bookstore-microservices/catalog/internal/application/port/pb/categorypb"
	"github.com/codejsha/bookstore-microservices/catalog/internal/application/usecase"
	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/handler"
	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/model/condition"
)

var _ categorypb.CategoryServiceServer = (*CategoryGrpcServer)(nil)

func NewCategoryGrpcServer(
	categoryHandler handler.CategoryHandler,
) *CategoryGrpcServer {
	return &CategoryGrpcServer{
		categoryHandler: categoryHandler,
	}
}

type CategoryGrpcServer struct {
	categoryHandler handler.CategoryHandler
	categorypb.UnimplementedCategoryServiceServer
}

func (s CategoryGrpcServer) FindAllCategories(ctx context.Context, request *categorypb.CategoryFindAllProtoReq) (*categorypb.CategoryFindAllProtoResp, error) {
	var categoryCondition condition.CategoryCondition
	cond, err := categoryCondition.FromValue(request)
	if err != nil {
		return nil, err
	}

	total, categories, err := s.categoryHandler.HandleFindAll(ctx, usecase.CategoryFindAllQuery{Cond: cond})
	if err != nil {
		return nil, err
	}

	resps := make([]*categorypb.CategoryFindProtoResp, len(categories))
	for i, category := range categories {
		resps[i] = category.ToCategoryFindProtoResp()
	}
	response := &categorypb.CategoryFindAllProtoResp{Total: total, Items: resps}
	return response, nil
}
