package service

import (
	"context"

	"github.com/codejsha/bookstore-microservices/catalog/internal/application/port/pb/stockpb"
	"github.com/codejsha/bookstore-microservices/catalog/internal/application/usecase"
	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/aggregate"
	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/handler"
	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/model/condition"
	"github.com/codejsha/bookstore-microservices/catalog/internal/infrastructure/adapter/protostub"
)

var _ usecase.BookUseCase = (*bookService)(nil)

func NewBookService(
	authorHandler handler.AuthorHandler,
	bookHandler handler.BookHandler,
	publisherHandler handler.PublisherHandler,
	categoryHandler handler.CategoryHandler,
	stockGrpcClient *protostub.StockGrpcClient,
) usecase.BookUseCase {
	return &bookService{
		authorHandler:      authorHandler,
		bookHandler:        bookHandler,
		publisherHandler:   publisherHandler,
		categoryHandler:    categoryHandler,
		stockServiceClient: stockGrpcClient.Client,
	}
}

type bookService struct {
	authorHandler      handler.AuthorHandler
	bookHandler        handler.BookHandler
	publisherHandler   handler.PublisherHandler
	categoryHandler    handler.CategoryHandler
	stockServiceClient stockpb.StockServiceClient
}

func (s bookService) FindAllBooks(ctx context.Context, cond condition.BookCondition) (int64, []*aggregate.BookAggregate, error) {
	// find books
	total, books, err := s.bookHandler.HandleFindAll(ctx, usecase.BookFindAllQuery{Cond: cond})
	if err != nil {
		return 0, nil, err
	}

	// find related data
	bookAggs := make([]*aggregate.BookAggregate, len(books))
	for _, book := range books {
		publisher, err := s.publisherHandler.HandleFind(ctx, usecase.PublisherFindQuery{Id: book.PublisherId})
		if err != nil {
			return 0, nil, err
		}
		_, authors, err := s.authorHandler.HandleFindAll(ctx, usecase.AuthorFindAllQuery{Cond: condition.AuthorCondition{BookId: book.Id}})
		if err != nil {
			return 0, nil, err
		}
		stocks, err := s.stockServiceClient.FindAllStocks(ctx, &stockpb.StockFindAllProtoReq{BookId: book.Id})
		if err != nil {
			return 0, nil, err
		}

		bookAgg := aggregate.NewBookAggregate(book, publisher, authors, stocks)
		bookAggs = append(bookAggs, bookAgg)
	}

	// return response
	return total, bookAggs, nil
}

func (s bookService) FindBook(ctx context.Context, id int64) (*aggregate.BookAggregate, error) {
	// find book
	book, err := s.bookHandler.HandleFind(ctx, usecase.BookFindQuery{Id: id})
	if err != nil {
		return nil, err
	}

	// find related data
	publisher, err := s.publisherHandler.HandleFind(ctx, usecase.PublisherFindQuery{Id: book.PublisherId})
	if err != nil {
		return nil, err
	}
	_, authors, err := s.authorHandler.HandleFindAll(ctx, usecase.AuthorFindAllQuery{Cond: condition.AuthorCondition{BookId: book.Id}})
	if err != nil {
		return nil, err
	}
	stocks, err := s.stockServiceClient.FindAllStocks(ctx, &stockpb.StockFindAllProtoReq{BookId: book.Id})
	if err != nil {
		return nil, err
	}

	// return response
	bookAgg := aggregate.NewBookAggregate(book, publisher, authors, stocks)
	return bookAgg, nil
}
