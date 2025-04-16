package protosvc

import (
	"context"

	"github.com/codejsha/bookstore-microservices/catalog/internal/application/port/pb/authorpb"
	"github.com/codejsha/bookstore-microservices/catalog/internal/application/usecase"
	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/handler"
	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/model/condition"
)

var _ authorpb.AuthorServiceServer = (*AuthorGrpcServer)(nil)

func NewAuthorGrpcServer(
	authorHandler handler.AuthorHandler,
) *AuthorGrpcServer {
	return &AuthorGrpcServer{
		authorHandler: authorHandler,
	}
}

type AuthorGrpcServer struct {
	authorHandler handler.AuthorHandler
	authorpb.UnimplementedAuthorServiceServer
}

func (s AuthorGrpcServer) FindAllAuthors(ctx context.Context, request *authorpb.AuthorFindAllProtoReq) (*authorpb.AuthorFindAllProtoResp, error) {
	var authorCondition condition.AuthorCondition
	cond, err := authorCondition.FromValue(request)
	if err != nil {
		return nil, err
	}

	total, authors, err := s.authorHandler.HandleFindAll(ctx, usecase.AuthorFindAllQuery{Cond: cond})
	if err != nil {
		return nil, err
	}

	resps := make([]*authorpb.AuthorFindProtoResp, len(authors))
	for i, author := range authors {
		resps[i] = author.ToAuthorFindProtoResp()
	}
	response := &authorpb.AuthorFindAllProtoResp{Total: total, Items: resps}
	return response, nil
}

func (s AuthorGrpcServer) FindAuthor(ctx context.Context, request *authorpb.AuthorFindProtoReq) (*authorpb.AuthorFindProtoResp, error) {
	author, err := s.authorHandler.HandleFind(ctx, usecase.AuthorFindQuery{Id: request.Id})
	if err != nil {
		return nil, err
	}

	response := author.ToAuthorFindProtoResp()
	return response, nil

}
