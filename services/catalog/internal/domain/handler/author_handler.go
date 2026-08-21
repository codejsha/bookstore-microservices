package handler

import (
	"context"
	"fmt"

	"github.com/codejsha/bookstore-microservices/catalog/internal/application/port/repo"
	"github.com/codejsha/bookstore-microservices/catalog/internal/application/usecase"
	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/aggregate/entity"
)

var _ AuthorHandler = (*authorHandlerImpl)(nil)

func NewAuthorHandler(
	authorRepo repo.AuthorRepo,
) AuthorHandler {
	return &authorHandlerImpl{
		authorRepo: authorRepo,
	}
}

type AuthorHandler interface {
	Handle(ctx context.Context, command usecase.AuthorCommand) (entity.AuthorEntity, error)
	HandleFindAll(ctx context.Context, query usecase.AuthorFindAllQuery) (int64, []entity.AuthorEntity, error)
	HandleFind(ctx context.Context, query usecase.AuthorFindQuery) (entity.AuthorEntity, error)
}

type authorHandlerImpl struct {
	authorRepo repo.AuthorRepo
}

func (h authorHandlerImpl) Handle(ctx context.Context, command usecase.AuthorCommand) (entity.AuthorEntity, error) {
	switch command.(type) {
	case usecase.AuthorCreateCommand:
		c := command.(usecase.AuthorCreateCommand)
		return h.authorRepo.Create(ctx, c.Ent)
	case usecase.AuthorUpdateCommand:
		c := command.(usecase.AuthorUpdateCommand)
		return h.authorRepo.Update(ctx, c.Ent)
	case usecase.AuthorDeleteCommand:
		c := command.(usecase.AuthorDeleteCommand)
		return entity.AuthorEntity{}, h.authorRepo.Delete(ctx, c.Id)
	default:
		return entity.AuthorEntity{}, fmt.Errorf("unknown command")
	}
}

func (h authorHandlerImpl) HandleFindAll(ctx context.Context, query usecase.AuthorFindAllQuery) (int64, []entity.AuthorEntity, error) {
	return h.authorRepo.FindAll(ctx, query.Cond)
}

func (h authorHandlerImpl) HandleFind(ctx context.Context, query usecase.AuthorFindQuery) (entity.AuthorEntity, error) {
	return h.authorRepo.Find(ctx, query.Id)
}
