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

var _ openapi.ClosingApi = (*closingController)(nil)

type closingController struct {
	inventoryUseCase usecase.InventoryUseCase
}

func NewClosingController(inventoryUseCase usecase.InventoryUseCase) openapi.ClosingApi {
	return &closingController{inventoryUseCase: inventoryUseCase}
}

func (c *closingController) ClosingsGetAll(
	ctx context.Context,
	warehouseUid *string,
	year *int32,
	month *int32,
	status *openapi.ClosingStatus,
	size *int32,
	page *int32,
	sort *string,
) (*openapi.ClosingFindAllResponse, error) {
	opt := option.NewClosingQueryOption(
		option.ClosingQueryOption{}.WithWarehouseUid(warehouseUid),
		option.ClosingQueryOption{}.WithYear(year),
		option.ClosingQueryOption{}.WithMonth(month),
		option.ClosingQueryOption{}.WithStatus(statusString(status)),
		option.ClosingQueryOption{}.WithPage(pagination.NewPageOption(size, page, sort)),
	)

	total, closings, err := c.inventoryUseCase.FindAllClosings(ctx, opt)
	if err != nil {
		return nil, err
	}

	items := make([]openapi.ClosingFindResponse, len(closings))
	for i, cl := range closings {
		items[i] = toClosingFindResponse(cl)
	}
	return &openapi.ClosingFindAllResponse{Items: items, Total: total}, nil
}

func (c *closingController) ClosingsCreate(ctx context.Context, req openapi.ClosingCreateRequest) (*openapi.ClosingFindResponse, error) {
	cmd := command.MonthlyClosingCommand{
		WarehouseUid: req.WarehouseUid,
		Year:         req.Year,
		Month:        req.Month,
	}

	closing, err := c.inventoryUseCase.CreateMonthlyClosing(ctx, cmd)
	if err != nil {
		return nil, httpx.MapBusinessError(ctx, err)
	}

	go runSideEffects(context.WithoutCancel(ctx), "closing", closing.Uid, "created", logrus.Fields{
		"warehouse_uid": req.WarehouseUid,
		"year":          req.Year,
		"month":         req.Month,
	})
	resp := toClosingFindResponse(closing)
	return &resp, nil
}

func (c *closingController) ClosingsRead(ctx context.Context, uid string) (*openapi.ClosingFindResponse, error) {
	closing, err := c.inventoryUseCase.FindClosing(ctx, uid)
	if err != nil {
		return nil, httpx.MapNotFound(ctx, err)
	}
	if closing == nil {
		return nil, httpx.MapNotFound(ctx, httpx.ErrNotFound)
	}
	resp := toClosingFindResponse(closing)
	return &resp, nil
}

// ─── mapping helpers ────────────────────────────────────────────────────────

func toClosingFindResponse(c *aggregate.MonthlyClosingAggregate) openapi.ClosingFindResponse {
	items := make([]openapi.ClosingItemResponse, len(c.Items))
	for i, it := range c.Items {
		items[i] = openapi.ClosingItemResponse{
			EditionUid:       it.EditionUid,
			OpeningQuantity:  it.OpeningQuantity,
			InboundQuantity:  it.InboundQuantity,
			OutboundQuantity: it.OutboundQuantity,
			AdjustQuantity:   it.AdjustQuantity,
			ClosingQuantity:  it.ClosingQuantity,
		}
	}
	return openapi.ClosingFindResponse{
		Uid:          c.Uid,
		WarehouseUid: c.WarehouseUid,
		Year:         c.Year,
		Month:        c.Month,
		Status:       openapi.ClosingStatus(c.Status),
		Items:        items,
		ClosedAt:     c.ClosedAt,
		CreatedAt:    c.CreatedAt,
		UpdatedAt:    c.UpdatedAt,
	}
}
