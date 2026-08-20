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

var _ openapi.TransferApi = (*transferController)(nil)

type transferController struct {
	inventoryUseCase usecase.InventoryUseCase
}

func NewTransferController(inventoryUseCase usecase.InventoryUseCase) openapi.TransferApi {
	return &transferController{inventoryUseCase: inventoryUseCase}
}

func (c *transferController) TransfersGetAll(
	ctx context.Context,
	editionUid *string,
	sourceWarehouseUid *string,
	targetWarehouseUid *string,
	status *openapi.TransferStatus,
	size *int32,
	page *int32,
	sort *string,
) (*openapi.TransferFindAllResponse, error) {
	opt := option.NewTransferQueryOption(
		option.TransferQueryOption{}.WithEditionUid(editionUid),
		option.TransferQueryOption{}.WithSourceWarehouseUid(sourceWarehouseUid),
		option.TransferQueryOption{}.WithTargetWarehouseUid(targetWarehouseUid),
		option.TransferQueryOption{}.WithStatus(statusString(status)),
		option.TransferQueryOption{}.WithPage(pagination.NewPageOption(size, page, sort)),
	)

	total, transfers, err := c.inventoryUseCase.FindAllTransfers(ctx, opt)
	if err != nil {
		return nil, err
	}

	items := make([]openapi.TransferFindResponse, len(transfers))
	for i, t := range transfers {
		items[i] = toTransferFindResponse(t)
	}
	return &openapi.TransferFindAllResponse{Items: items, Total: total}, nil
}

func (c *transferController) TransfersCreate(ctx context.Context, req openapi.TransferCreateRequest) (*openapi.TransferFindResponse, error) {
	cmd := command.StockTransferCommand{
		EditionUid:         req.EditionUid,
		SourceWarehouseUid: req.SourceWarehouseUid,
		TargetWarehouseUid: req.TargetWarehouseUid,
		Quantity:           req.Quantity,
		Reason:             req.Reason,
	}

	transfer, err := c.inventoryUseCase.TransferStock(ctx, cmd)
	if err != nil {
		return nil, httpx.MapBusinessError(ctx, err)
	}

	go runSideEffects(context.WithoutCancel(ctx), "transfer", transfer.Uid, "created", logrus.Fields{
		"edition_uid":          req.EditionUid,
		"source_warehouse_uid": req.SourceWarehouseUid,
		"target_warehouse_uid": req.TargetWarehouseUid,
		"quantity":             req.Quantity,
	})
	resp := toTransferFindResponse(transfer)
	return &resp, nil
}

func (c *transferController) TransfersRead(ctx context.Context, uid string) (*openapi.TransferFindResponse, error) {
	transfer, err := c.inventoryUseCase.FindTransfer(ctx, uid)
	if err != nil {
		return nil, httpx.MapNotFound(ctx, err)
	}
	if transfer == nil {
		return nil, httpx.MapNotFound(ctx, httpx.ErrNotFound)
	}
	resp := toTransferFindResponse(transfer)
	return &resp, nil
}

func (c *transferController) TransfersComplete(ctx context.Context, uid string) (*openapi.TransferFindResponse, error) {
	transfer, err := c.inventoryUseCase.CompleteTransfer(ctx, uid)
	if err != nil {
		return nil, err
	}
	if transfer == nil {
		return nil, httpx.MapNotFound(ctx, httpx.ErrNotFound)
	}
	go runSideEffects(context.WithoutCancel(ctx), "transfer", uid, "completed", logrus.Fields{})
	resp := toTransferFindResponse(transfer)
	return &resp, nil
}

func (c *transferController) TransfersCancel(ctx context.Context, uid string) (*openapi.TransferFindResponse, error) {
	transfer, err := c.inventoryUseCase.CancelTransfer(ctx, uid)
	if err != nil {
		return nil, err
	}
	if transfer == nil {
		return nil, httpx.MapNotFound(ctx, httpx.ErrNotFound)
	}
	go runSideEffects(context.WithoutCancel(ctx), "transfer", uid, "cancelled", logrus.Fields{})
	resp := toTransferFindResponse(transfer)
	return &resp, nil
}

// ─── mapping helpers ────────────────────────────────────────────────────────

func toTransferFindResponse(t *aggregate.StockTransferAggregate) openapi.TransferFindResponse {
	return openapi.TransferFindResponse{
		Uid:                t.Uid,
		EditionUid:         t.EditionUid,
		SourceWarehouseUid: t.SourceWarehouseUid,
		TargetWarehouseUid: t.TargetWarehouseUid,
		Quantity:           t.Quantity,
		Status:             openapi.TransferStatus(t.Status),
		Reason:             t.Reason,
		CompletedAt:        t.CompletedAt,
		CreatedAt:          t.CreatedAt,
		UpdatedAt:          t.UpdatedAt,
	}
}

func statusString[T ~string](v *T) *string {
	if v == nil {
		return nil
	}
	s := string(*v)
	return &s
}
