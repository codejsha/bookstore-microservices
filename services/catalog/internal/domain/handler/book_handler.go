package handler

import (
	"context"
	"fmt"

	"github.com/codejsha/bookstore-microservices/catalog/internal/application/port/repo"
	"github.com/codejsha/bookstore-microservices/catalog/internal/application/usecase"
	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/aggregate/entity"
)

var _ BookHandler = (*bookHandlerImpl)(nil)

func NewBookHandler(
	bookRepo repo.BookRepo,
) BookHandler {
	return &bookHandlerImpl{
		bookRepo: bookRepo,
	}
}

type BookHandler interface {
	Handle(ctx context.Context, command usecase.BookCommand) (entity.BookEntity, error)
	HandleFindAll(ctx context.Context, query usecase.BookFindAllQuery) (int64, []entity.BookEntity, error)
	HandleFind(ctx context.Context, query usecase.BookFindQuery) (entity.BookEntity, error)
}

type bookHandlerImpl struct {
	bookRepo repo.BookRepo
}

func (h bookHandlerImpl) Handle(ctx context.Context, command usecase.BookCommand) (entity.BookEntity, error) {
	switch command.(type) {
	case usecase.BookCreateCommand:
		c := command.(usecase.BookCreateCommand)
		return h.bookRepo.Create(ctx, c.Ent)
	case usecase.BookUpdateCommand:
		c := command.(usecase.BookUpdateCommand)
		return h.bookRepo.Update(ctx, c.Ent)
	case usecase.BookDeleteCommand:
		c := command.(usecase.BookDeleteCommand)
		return entity.BookEntity{}, h.bookRepo.Delete(ctx, c.Id)
	default:
		return entity.BookEntity{}, fmt.Errorf("unknown command")
	}
}

func (h bookHandlerImpl) HandleFindAll(ctx context.Context, query usecase.BookFindAllQuery) (int64, []entity.BookEntity, error) {
	return h.bookRepo.FindAll(ctx, query.Cond)
}

func (h bookHandlerImpl) HandleFind(ctx context.Context, query usecase.BookFindQuery) (entity.BookEntity, error) {
	return h.bookRepo.Find(ctx, query.Id)
}
