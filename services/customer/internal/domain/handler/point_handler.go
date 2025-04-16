package handler

import (
	"context"
	"fmt"

	"github.com/codejsha/bookstore-microservices/customer/internal/application/port/repo"
	"github.com/codejsha/bookstore-microservices/customer/internal/application/usecase"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/aggregate/entity"
)

var _ PointHandler = (*pointHandlerImpl)(nil)

func NewPointHandler(
	pointRepo repo.PointRepo,
) PointHandler {
	return &pointHandlerImpl{
		pointRepo: pointRepo,
	}
}

type PointHandler interface {
	Handle(ctx context.Context, command usecase.PointCommand) (entity.PointEntity, error)
	HandleFindAll(ctx context.Context, query usecase.PointFindAllQuery) (int64, []entity.PointEntity, error)
	HandleFind(ctx context.Context, query usecase.PointFindQuery) (entity.PointEntity, error)
}

type pointHandlerImpl struct {
	pointRepo repo.PointRepo
}

func (h pointHandlerImpl) Handle(ctx context.Context, command usecase.PointCommand) (entity.PointEntity, error) {
	switch command.(type) {
	case usecase.PointUpdateCommand:
		c := command.(usecase.PointUpdateCommand)
		return h.pointRepo.Update(ctx, c.Ent)
	default:
		return entity.PointEntity{}, fmt.Errorf("unknown command")
	}
}

func (h pointHandlerImpl) HandleFindAll(ctx context.Context, query usecase.PointFindAllQuery) (int64, []entity.PointEntity, error) {
	return h.pointRepo.FindAll(ctx, query.Cond)
}

func (h pointHandlerImpl) HandleFind(ctx context.Context, query usecase.PointFindQuery) (entity.PointEntity, error) {
	return h.pointRepo.Find(ctx, query.Id)
}
