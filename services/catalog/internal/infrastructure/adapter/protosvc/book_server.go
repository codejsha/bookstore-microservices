package protosvc

import (
	"context"

	"github.com/codejsha/bookstore-microservices/catalog/internal/application/port/pb/bookpb"
	"github.com/codejsha/bookstore-microservices/catalog/internal/application/usecase"
	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/model/condition"
)

var _ bookpb.BookServiceServer = (*BookGrpcServer)(nil)

func NewBookGrpcServer(
	bookUseCase usecase.BookUseCase,
) *BookGrpcServer {
	return &BookGrpcServer{
		bookUseCase: bookUseCase,
	}
}

type BookGrpcServer struct {
	bookUseCase usecase.BookUseCase
	bookpb.UnimplementedBookServiceServer
}

func (s BookGrpcServer) FindAllBooks(ctx context.Context, request *bookpb.BookFindAllProtoReq) (*bookpb.BookFindAllProtoResp, error) {
	var bookCondition condition.BookCondition
	cond, err := bookCondition.FromValue(request)
	if err != nil {
		return nil, err
	}

	total, bookAggs, err := s.bookUseCase.FindAllBooks(ctx, cond)
	if err != nil {
		return nil, err
	}

	resps := make([]*bookpb.BookFindProtoResp, len(bookAggs))
	for i, bookAgg := range bookAggs {
		resps[i] = bookAgg.ToBookFindResponse()
	}
	response := &bookpb.BookFindAllProtoResp{Total: total, Items: resps}
	return response, nil
}

func (s BookGrpcServer) FindBook(ctx context.Context, request *bookpb.BookFindProtoReq) (*bookpb.BookFindProtoResp, error) {
	bookAgg, err := s.bookUseCase.FindBook(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	response := bookAgg.ToBookFindResponse()
	return response, nil
}
