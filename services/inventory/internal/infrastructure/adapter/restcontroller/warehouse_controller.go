package restcontroller

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/codejsha/shared-library-go/pkg/pagination"

	"github.com/codejsha/bookstore-microservices/inventory/generated/application/port/openapi"
	"github.com/codejsha/bookstore-microservices/inventory/internal/application/usecase"
	"github.com/codejsha/bookstore-microservices/inventory/internal/domain/aggregate"
	"github.com/codejsha/bookstore-microservices/inventory/internal/domain/model/command"
	"github.com/codejsha/bookstore-microservices/inventory/internal/domain/model/option"
	"github.com/codejsha/bookstore-microservices/inventory/internal/infrastructure/httpx"
)

var _ openapi.WarehouseApi = (*warehouseController)(nil)

type warehouseController struct {
	inventoryUseCase usecase.InventoryUseCase
}

func NewWarehouseController(inventoryUseCase usecase.InventoryUseCase) openapi.WarehouseApi {
	return &warehouseController{inventoryUseCase: inventoryUseCase}
}

func (c *warehouseController) WarehousesGetAll(
	ctx context.Context,
	name *string,
	size *int32,
	page *int32,
	sort *string,
) (*openapi.WarehouseFindAllResponse, error) {
	opt := option.NewWarehouseQueryOption(
		option.WarehouseQueryOption{}.WithName(name),
		option.WarehouseQueryOption{}.WithPage(pagination.NewPageOption(size, page, sort)),
	)

	total, warehouses, err := c.inventoryUseCase.FindAllWarehouses(ctx, opt)
	if err != nil {
		return nil, err
	}

	items := make([]openapi.WarehouseFindResponse, len(warehouses))
	for i, w := range warehouses {
		items[i] = toWarehouseFindResponse(w)
	}
	return &openapi.WarehouseFindAllResponse{Items: items, Total: total}, nil
}

func (c *warehouseController) WarehousesCreate(ctx context.Context, req openapi.WarehouseCreateRequest) error {
	cmd := command.WarehouseCreateCommand{
		Name:     req.Name,
		Address:  req.Address,
		Capacity: req.Capacity,
	}

	warehouse, err := c.inventoryUseCase.RegisterWarehouse(ctx, cmd)
	if err != nil {
		return httpx.MapBusinessError(ctx, err)
	}

	go runSideEffects(context.WithoutCancel(ctx), "warehouse", warehouse.Uid, "created", logrus.Fields{"name": req.Name})
	return nil
}

func (c *warehouseController) WarehousesRead(ctx context.Context, uid string) (*openapi.WarehouseFindResponse, error) {
	if err := requireUidParam(ctx, "uid", uid); err != nil {
		return nil, err
	}

	warehouse, err := c.inventoryUseCase.FindWarehouse(ctx, uid)
	if err != nil {
		return nil, httpx.MapNotFound(ctx, err)
	}
	if warehouse == nil {
		return nil, httpx.MapNotFound(ctx, httpx.ErrNotFound)
	}
	resp := toWarehouseFindResponse(warehouse)
	return &resp, nil
}

func (c *warehouseController) WarehousesUpdate(
	ctx context.Context,
	uid string,
	req openapi.WarehouseUpdateRequest,
) (*openapi.WarehouseUpdateResponse, error) {
	if err := requireUidParam(ctx, "uid", uid); err != nil {
		return nil, err
	}

	cmd := command.WarehouseUpdateCommand{
		Name:     req.Name,
		Address:  req.Address,
		Capacity: req.Capacity,
	}

	warehouse, err := c.inventoryUseCase.UpdateWarehouse(ctx, uid, cmd)
	if err != nil {
		return nil, httpx.MapBusinessError(ctx, httpx.MapNotFound(ctx, err))
	}
	if warehouse == nil {
		return nil, httpx.MapNotFound(ctx, httpx.ErrNotFound)
	}

	go runSideEffects(context.WithoutCancel(ctx), "warehouse", uid, "updated", logrus.Fields{})

	return &openapi.WarehouseUpdateResponse{
		Uid:       warehouse.Uid,
		Name:      warehouse.Name,
		Address:   warehouse.Address,
		Capacity:  warehouse.Capacity,
		CreatedAt: warehouse.CreatedAt,
		UpdatedAt: warehouse.UpdatedAt,
	}, nil
}

// ─── mapping helpers ────────────────────────────────────────────────────────

func toWarehouseFindResponse(w *aggregate.WarehouseAggregate) openapi.WarehouseFindResponse {
	return openapi.WarehouseFindResponse{
		Uid:       w.Uid,
		Name:      w.Name,
		Address:   w.Address,
		Capacity:  w.Capacity,
		CreatedAt: w.CreatedAt,
		UpdatedAt: w.UpdatedAt,
	}
}
