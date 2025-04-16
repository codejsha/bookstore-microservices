package handler

import (
	"context"

	"github.com/codejsha/bookstore-microservices/inventory/internal/application/port/repo"
	"github.com/codejsha/bookstore-microservices/inventory/internal/application/usecase"
	"github.com/codejsha/bookstore-microservices/inventory/internal/domain/aggregate/entity"
)

var _ WarehouseHandler = (*warehouseHandlerImpl)(nil)

func NewWarehouseHandler(
	warehouseRepo repo.WarehouseRepo,
) WarehouseHandler {
	return &warehouseHandlerImpl{
		warehouseRepo: warehouseRepo,
	}
}

type WarehouseHandler interface {
	Handle(ctx context.Context, command usecase.WarehouseCommand) (entity.WarehouseEntity, error)
	HandleFindAll(ctx context.Context, query usecase.WarehouseFindAllQuery) (int64, []entity.WarehouseEntity, error)
	HandleFind(ctx context.Context, query usecase.WarehouseFindQuery) (entity.WarehouseEntity, error)
}

type warehouseHandlerImpl struct {
	warehouseRepo repo.WarehouseRepo
}

func (h warehouseHandlerImpl) Handle(ctx context.Context, command usecase.WarehouseCommand) (entity.WarehouseEntity, error) {
	switch command.(type) {
	case usecase.WarehouseCreateCommand:
		c := command.(usecase.WarehouseCreateCommand)
		return h.warehouseRepo.Create(ctx, c.Ent)
	case usecase.WarehouseUpdateCommand:
		c := command.(usecase.WarehouseUpdateCommand)
		return h.warehouseRepo.Update(ctx, c.Ent)
	case usecase.WarehouseDeleteCommand:
		c := command.(usecase.WarehouseDeleteCommand)
		return entity.WarehouseEntity{}, h.warehouseRepo.Delete(ctx, c.Id)
	default:
		return entity.WarehouseEntity{}, nil
	}
}

func (h warehouseHandlerImpl) HandleFindAll(ctx context.Context, query usecase.WarehouseFindAllQuery) (int64, []entity.WarehouseEntity, error) {
	return h.warehouseRepo.FindAll(ctx, query.Cond)
}

func (h warehouseHandlerImpl) HandleFind(ctx context.Context, query usecase.WarehouseFindQuery) (entity.WarehouseEntity, error) {
	return h.warehouseRepo.Find(ctx, query.Id)
}
