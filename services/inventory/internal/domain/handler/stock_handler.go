package handler

import (
	"context"
	"fmt"

	"github.com/codejsha/bookstore-microservices/inventory/internal/application/port/repo"
	"github.com/codejsha/bookstore-microservices/inventory/internal/application/usecase"
	"github.com/codejsha/bookstore-microservices/inventory/internal/domain/aggregate/entity"
)

var _ StockHandler = (*stockHandlerImpl)(nil)

func NewStockHandler(
	stockRepo repo.StockRepo,
) StockHandler {
	return &stockHandlerImpl{
		stockRepo: stockRepo,
	}
}

type StockHandler interface {
	Handle(ctx context.Context, command usecase.StockCommand) (entity.StockEntity, error)
	HandleFindAll(ctx context.Context, query usecase.StockFindAllQuery) (int64, []entity.StockEntity, error)
	HandleFind(ctx context.Context, query usecase.StockFindQuery) (entity.StockEntity, error)
}

type stockHandlerImpl struct {
	stockRepo repo.StockRepo
}

func (h stockHandlerImpl) Handle(ctx context.Context, command usecase.StockCommand) (entity.StockEntity, error) {
	switch command.(type) {
	case usecase.StockCreateCommand:
		c := command.(usecase.StockCreateCommand)
		return h.stockRepo.Update(ctx, c.Ent)
	case usecase.StockUpdateCommand:
		c := command.(usecase.StockUpdateCommand)
		return h.stockRepo.Update(ctx, c.Ent)
	case usecase.StockDeleteCommand:
		c := command.(usecase.StockDeleteCommand)
		return entity.StockEntity{}, h.stockRepo.Delete(ctx, c.Id)
	default:
		return entity.StockEntity{}, fmt.Errorf("unknown command")
	}
}

func (h stockHandlerImpl) HandleFindAll(ctx context.Context, query usecase.StockFindAllQuery) (int64, []entity.StockEntity, error) {
	return h.stockRepo.FindAll(ctx, query.Cond)
}

func (h stockHandlerImpl) HandleFind(ctx context.Context, query usecase.StockFindQuery) (entity.StockEntity, error) {
	return h.stockRepo.Find(ctx, query.Id)
}
