package workflow

import (
	"context"
	"errors"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/temporal"

	"github.com/codejsha/bookstore-microservices/inventory/internal/application/port/repo"
	"github.com/codejsha/bookstore-microservices/inventory/internal/application/usecase"
	"github.com/codejsha/bookstore-microservices/inventory/internal/domain/model/command"
)

const (
	errTypeInsufficientStock = "InsufficientStock"
	errTypeInvalidQuantity   = "InvalidQuantity"
	errTypeInvalidCommand    = "InvalidCommand"
)

func classifyReservationError(err error) error {
	switch {
	case errors.Is(err, repo.ErrInsufficientStock):
		return temporal.NewNonRetryableApplicationError(err.Error(), errTypeInsufficientStock, err)
	case errors.Is(err, repo.ErrInvalidQuantity):
		return temporal.NewNonRetryableApplicationError(err.Error(), errTypeInvalidQuantity, err)
	case errors.Is(err, command.ErrInvalidCommand):
		return temporal.NewNonRetryableApplicationError(err.Error(), errTypeInvalidCommand, err)
	default:
		return err
	}
}

type StockActivities struct {
	inventoryUC usecase.InventoryUseCase
}

func NewStockActivities(inventoryUC usecase.InventoryUseCase) *StockActivities {
	return &StockActivities{inventoryUC: inventoryUC}
}

func (a *StockActivities) ReserveStock(ctx context.Context, req ReserveStockRequest) error {
	logger := activity.GetLogger(ctx)
	reason := "temporal-saga-reserve"
	orderRef := a.orderRef(ctx, req)

	for _, item := range req.Items {
		if _, err := a.inventoryUC.ReserveStockForOrder(ctx, command.StockOrderReserveCommand{
			OrderUid:  orderRef,
			EditionId: item.ProductID,
			Quantity:  item.Quantity,
			Reason:    &reason,
		}); err != nil {
			logger.Error("Failed to reserve stock", "orderRef", orderRef, "productId", item.ProductID, "error", err)
			return classifyReservationError(err)
		}
		logger.Info("Reserved stock", "orderRef", orderRef, "productId", item.ProductID, "quantity", item.Quantity)
	}
	return nil
}

func (a *StockActivities) ReleaseStock(ctx context.Context, req ReleaseStockRequest) error {
	logger := activity.GetLogger(ctx)
	reason := "temporal-saga-compensation"
	orderRef := a.orderRef(ctx, req)

	for _, item := range req.Items {
		if _, err := a.inventoryUC.ReleaseStockForOrder(ctx, command.StockOrderReleaseCommand{
			OrderUid:  orderRef,
			EditionId: item.ProductID,
			Quantity:  item.Quantity,
			Reason:    &reason,
		}); err != nil {
			logger.Error("Failed to release stock", "orderRef", orderRef, "productId", item.ProductID, "error", err)
			return classifyReservationError(err)
		}
		logger.Info("Released stock", "orderRef", orderRef, "productId", item.ProductID, "quantity", item.Quantity)
	}
	return nil
}

func (a *StockActivities) orderRef(ctx context.Context, req ReserveStockRequest) string {
	if req.OrderUid != "" {
		return req.OrderUid
	}
	return activity.GetInfo(ctx).WorkflowExecution.ID
}
