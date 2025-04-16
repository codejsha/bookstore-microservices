package handler

import (
	"context"
	"fmt"

	"github.com/codejsha/bookstore-microservices/catalog/internal/application/port/repo"
	"github.com/codejsha/bookstore-microservices/catalog/internal/application/usecase"
	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/aggregate/entity"
)

var _ PublisherHandler = (*publisherHandlerImpl)(nil)

func NewPublisherHandler(
	publisherRepo repo.PublisherRepo,
) PublisherHandler {
	return &publisherHandlerImpl{
		publisherRepo: publisherRepo,
	}
}

type PublisherHandler interface {
	Handle(ctx context.Context, command usecase.PublisherCommand) (entity.PublisherEntity, error)
	HandleFindAll(ctx context.Context, query usecase.PublisherFindAllQuery) (int64, []entity.PublisherEntity, error)
	HandleFind(ctx context.Context, query usecase.PublisherFindQuery) (entity.PublisherEntity, error)
}

type publisherHandlerImpl struct {
	publisherRepo repo.PublisherRepo
}

func (h publisherHandlerImpl) Handle(ctx context.Context, command usecase.PublisherCommand) (entity.PublisherEntity, error) {
	switch command.(type) {
	case usecase.PublisherCreateCommand:
		c := command.(usecase.PublisherCreateCommand)
		return h.publisherRepo.Create(ctx, c.Ent)
	case usecase.PublisherUpdateCommand:
		c := command.(usecase.PublisherUpdateCommand)
		return h.publisherRepo.Update(ctx, c.Ent)
	case usecase.PublisherDeleteCommand:
		c := command.(usecase.PublisherDeleteCommand)
		return entity.PublisherEntity{}, h.publisherRepo.Delete(ctx, c.Id)
	default:
		return entity.PublisherEntity{}, fmt.Errorf("unknown command")
	}
}

func (h publisherHandlerImpl) HandleFindAll(ctx context.Context, query usecase.PublisherFindAllQuery) (int64, []entity.PublisherEntity, error) {
	return h.publisherRepo.FindAll(ctx, query.Cond)
}

func (h publisherHandlerImpl) HandleFind(ctx context.Context, query usecase.PublisherFindQuery) (entity.PublisherEntity, error) {
	return h.publisherRepo.Find(ctx, query.Id)
}
