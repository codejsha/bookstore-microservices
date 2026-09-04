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

var _ openapi.AuditApi = (*auditController)(nil)

type auditController struct {
	inventoryUseCase usecase.InventoryUseCase
}

func NewAuditController(inventoryUseCase usecase.InventoryUseCase) openapi.AuditApi {
	return &auditController{inventoryUseCase: inventoryUseCase}
}

func (c *auditController) AuditsGetAll(
	ctx context.Context,
	warehouseUid *string,
	status *openapi.AuditStatus,
	size *int32,
	page *int32,
	sort *string,
) (*openapi.AuditFindAllResponse, error) {
	if err := optionalUidParam(ctx, "warehouse_uid", warehouseUid); err != nil {
		return nil, err
	}

	opt := option.NewAuditQueryOption(
		option.AuditQueryOption{}.WithWarehouseUid(warehouseUid),
		option.AuditQueryOption{}.WithStatus(statusString(status)),
		option.AuditQueryOption{}.WithPage(pagination.NewPageOption(size, page, sort)),
	)

	total, audits, err := c.inventoryUseCase.FindAllAudits(ctx, opt)
	if err != nil {
		return nil, err
	}

	items := make([]openapi.AuditFindResponse, len(audits))
	for i, a := range audits {
		items[i] = toAuditFindResponse(a)
	}
	return &openapi.AuditFindAllResponse{Items: items, Total: total}, nil
}

func (c *auditController) AuditsCreate(ctx context.Context, req openapi.AuditCreateRequest) (*openapi.AuditFindResponse, error) {
	items := make([]command.StockAuditItemCommand, len(req.Items))
	for i, it := range req.Items {
		items[i] = command.StockAuditItemCommand{
			EditionUid:     it.EditionUid,
			ActualQuantity: it.ActualQuantity,
		}
	}
	cmd := command.StockAuditCreateCommand{
		WarehouseUid: req.WarehouseUid,
		Items:        items,
		Notes:        req.Notes,
	}

	audit, err := c.inventoryUseCase.CreateAudit(ctx, cmd)
	if err != nil {
		return nil, httpx.MapBusinessError(ctx, err)
	}

	go runSideEffects(context.WithoutCancel(ctx), "audit", audit.Uid, "created", logrus.Fields{
		"warehouse_uid": req.WarehouseUid,
		"item_count":    len(req.Items),
	})
	resp := toAuditFindResponse(audit)
	return &resp, nil
}

func (c *auditController) AuditsRead(ctx context.Context, uid string) (*openapi.AuditFindResponse, error) {
	if err := requireUidParam(ctx, "uid", uid); err != nil {
		return nil, err
	}

	audit, err := c.inventoryUseCase.FindAudit(ctx, uid)
	if err != nil {
		return nil, httpx.MapNotFound(ctx, err)
	}
	if audit == nil {
		return nil, httpx.MapNotFound(ctx, httpx.ErrNotFound)
	}
	resp := toAuditFindResponse(audit)
	return &resp, nil
}

func (c *auditController) AuditsComplete(ctx context.Context, uid string) (*openapi.AuditFindResponse, error) {
	if err := requireUidParam(ctx, "uid", uid); err != nil {
		return nil, err
	}

	audit, err := c.inventoryUseCase.CompleteAudit(ctx, uid)
	if err != nil {
		return nil, httpx.MapBusinessError(ctx, err)
	}
	if audit == nil {
		return nil, httpx.MapNotFound(ctx, httpx.ErrNotFound)
	}
	go runSideEffects(context.WithoutCancel(ctx), "audit", uid, "completed", logrus.Fields{})
	resp := toAuditFindResponse(audit)
	return &resp, nil
}

// ─── mapping helpers ────────────────────────────────────────────────────────

func toAuditFindResponse(a *aggregate.StockAuditAggregate) openapi.AuditFindResponse {
	items := make([]openapi.AuditItemResponse, len(a.Items))
	for i, it := range a.Items {
		items[i] = openapi.AuditItemResponse{
			EditionUid:     it.EditionUid,
			SystemQuantity: it.SystemQuantity,
			ActualQuantity: it.ActualQuantity,
			Difference:     it.Difference,
		}
	}
	return openapi.AuditFindResponse{
		Uid:          a.Uid,
		WarehouseUid: a.WarehouseUid,
		Status:       openapi.AuditStatus(a.Status),
		Items:        items,
		Notes:        a.Notes,
		CompletedAt:  a.CompletedAt,
		CreatedAt:    a.CreatedAt,
		UpdatedAt:    a.UpdatedAt,
	}
}
