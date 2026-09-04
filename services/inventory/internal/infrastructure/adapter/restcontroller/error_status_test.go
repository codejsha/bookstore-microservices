package restcontroller

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/codejsha/bookstore-microservices/inventory/generated/application/port/openapi"
	"github.com/codejsha/bookstore-microservices/inventory/internal/application/port/repo"
	"github.com/codejsha/bookstore-microservices/inventory/internal/domain/aggregate"
	"github.com/codejsha/bookstore-microservices/inventory/internal/domain/model/command"
	"github.com/codejsha/bookstore-microservices/inventory/internal/infrastructure/httpx"
)

func stockNotFoundErr() error {
	return fmt.Errorf("stock not found for edition %s in warehouse %s: %w", uidEdition1, uidWarehouse1, repo.ErrStockNotFound)
}

func warehouseNotFoundErr(warehouseUid string) error {
	return fmt.Errorf("warehouse %s not found: %w", warehouseUid, repo.ErrWarehouseNotFound)
}

func serveController(t *testing.T, invoke func(ctx context.Context) error) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(httpx.GinResponseMapping())
	r.POST("/x", func(c *gin.Context) {
		if err := invoke(c.Request.Context()); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/x", nil))
	return w
}

func assertStockNotFoundResponse(t *testing.T, w *httptest.ResponseRecorder) {
	t.Helper()
	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", w.Code)
	}
	if body := w.Body.String(); body != `{"title":"Not Found","status":404,"detail":"stock not found"}` {
		t.Errorf("body = %q, want stock-not-found json", body)
	}
}

func assertWarehouseNotFoundResponse(t *testing.T, w *httptest.ResponseRecorder) {
	t.Helper()
	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", w.Code)
	}
	if body := w.Body.String(); body != `{"title":"Not Found","status":404,"detail":"warehouse not found"}` {
		t.Errorf("body = %q, want warehouse-not-found json", body)
	}
}

func TestTransferController_WhenCompleteHitsMissingStock_Returns404(t *testing.T) {
	use := &stubUseCase{
		completeTransfer: func(context.Context, string) (*aggregate.StockTransferAggregate, error) {
			return nil, stockNotFoundErr()
		},
	}
	ctrl := NewTransferController(use)
	assertStockNotFoundResponse(t, serveController(t, func(ctx context.Context) error {
		_, err := ctrl.TransfersComplete(ctx, uidTransfer1)
		return err
	}))
}

func TestTransferController_WhenCancelHitsMissingStock_Returns404(t *testing.T) {
	use := &stubUseCase{
		cancelTransfer: func(context.Context, string) (*aggregate.StockTransferAggregate, error) {
			return nil, stockNotFoundErr()
		},
	}
	ctrl := NewTransferController(use)
	assertStockNotFoundResponse(t, serveController(t, func(ctx context.Context) error {
		_, err := ctrl.TransfersCancel(ctx, uidTransfer1)
		return err
	}))
}

func TestTransferController_WhenCreateHitsMissingStock_Returns404(t *testing.T) {
	use := &stubUseCase{
		transferStock: func(context.Context, command.StockTransferCommand) (*aggregate.StockTransferAggregate, error) {
			return nil, fmt.Errorf("no stock for edition %s in source warehouse %s: %w", uidEdition1, uidWarehouse1, repo.ErrStockNotFound)
		},
	}
	ctrl := NewTransferController(use)
	assertStockNotFoundResponse(t, serveController(t, func(ctx context.Context) error {
		_, err := ctrl.TransfersCreate(ctx, openapi.TransferCreateRequest{
			EditionUid: uidEdition1, SourceWarehouseUid: uidWarehouse1, TargetWarehouseUid: uidWarehouse2, Quantity: 3,
		})
		return err
	}))
}

func TestAuditController_WhenCompleteHitsMissingStock_Returns404(t *testing.T) {
	use := &stubUseCase{
		completeAudit: func(context.Context, string) (*aggregate.StockAuditAggregate, error) {
			return nil, stockNotFoundErr()
		},
	}
	ctrl := NewAuditController(use)
	assertStockNotFoundResponse(t, serveController(t, func(ctx context.Context) error {
		_, err := ctrl.AuditsComplete(ctx, uidAudit1)
		return err
	}))
}

func TestAuditController_WhenCreateHitsMissingStock_Returns404(t *testing.T) {
	use := &stubUseCase{
		createAudit: func(context.Context, command.StockAuditCreateCommand) (*aggregate.StockAuditAggregate, error) {
			return nil, stockNotFoundErr()
		},
	}
	ctrl := NewAuditController(use)
	assertStockNotFoundResponse(t, serveController(t, func(ctx context.Context) error {
		_, err := ctrl.AuditsCreate(ctx, openapi.AuditCreateRequest{
			WarehouseUid: uidWarehouse1,
			Items:        []openapi.AuditItemRequest{{EditionUid: uidEdition1, ActualQuantity: 2}},
		})
		return err
	}))
}

// Every stock write is driven through the same subtest, named after the method under test.
func TestStockController_WhenWriteHitsMissingStock_Returns404(t *testing.T) {
	cases := []struct {
		name   string
		use    *stubUseCase
		invoke func(ctrl openapi.StockApi, ctx context.Context) error
	}{
		{
			name: "receive",
			use: &stubUseCase{receiveStock: func(context.Context, command.StockReceiveCommand) (*aggregate.StockAggregate, error) {
				return nil, stockNotFoundErr()
			}},
			invoke: func(ctrl openapi.StockApi, ctx context.Context) error {
				_, err := ctrl.StocksReceive(ctx, openapi.StockReceiveRequest{EditionUid: uidEdition1, WarehouseUid: uidWarehouse1, Quantity: 1})
				return err
			},
		},
		{
			name: "release",
			use: &stubUseCase{releaseStock: func(context.Context, command.StockReleaseCommand) (*aggregate.StockAggregate, error) {
				return nil, stockNotFoundErr()
			}},
			invoke: func(ctrl openapi.StockApi, ctx context.Context) error {
				_, err := ctrl.StocksRelease(ctx, openapi.StockReleaseRequest{EditionUid: uidEdition1, WarehouseUid: uidWarehouse1, Quantity: 1})
				return err
			},
		},
		{
			name: "adjust",
			use: &stubUseCase{adjustStock: func(context.Context, command.StockAdjustCommand) (*aggregate.StockAggregate, error) {
				return nil, stockNotFoundErr()
			}},
			invoke: func(ctrl openapi.StockApi, ctx context.Context) error {
				_, err := ctrl.StocksAdjust(ctx, openapi.StockAdjustRequest{EditionUid: uidEdition1, WarehouseUid: uidWarehouse1, Quantity: -1, Reason: "damaged"})
				return err
			},
		},
		{
			name: "reserve",
			use: &stubUseCase{reserveStock: func(context.Context, command.StockReserveCommand) (*aggregate.StockAggregate, error) {
				return nil, stockNotFoundErr()
			}},
			invoke: func(ctrl openapi.StockApi, ctx context.Context) error {
				_, err := ctrl.StocksReserve(ctx, openapi.StockReserveRequest{EditionUid: uidEdition1, WarehouseUid: uidWarehouse1, Quantity: 1})
				return err
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := NewStockController(tc.use)
			assertStockNotFoundResponse(t, serveController(t, func(ctx context.Context) error {
				return tc.invoke(ctrl, ctx)
			}))
		})
	}
}

func TestTransferController_WhenCompleteFailsUnexpectedly_Returns500(t *testing.T) {
	wantErr := errors.New("transfer cannot be completed in status COMPLETED")
	use := &stubUseCase{
		completeTransfer: func(context.Context, string) (*aggregate.StockTransferAggregate, error) { return nil, wantErr },
	}
	ctrl := NewTransferController(use)
	w := serveController(t, func(ctx context.Context) error {
		_, err := ctrl.TransfersComplete(ctx, uidTransfer1)
		return err
	})
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", w.Code)
	}
}

func TestTransferController_WhenCreateHitsMissingTargetWarehouse_Returns404(t *testing.T) {
	use := &stubUseCase{
		transferStock: func(context.Context, command.StockTransferCommand) (*aggregate.StockTransferAggregate, error) {
			return nil, fmt.Errorf("target warehouse %s not found: %w", uidWarehouse2, repo.ErrWarehouseNotFound)
		},
	}
	ctrl := NewTransferController(use)
	assertWarehouseNotFoundResponse(t, serveController(t, func(ctx context.Context) error {
		_, err := ctrl.TransfersCreate(ctx, openapi.TransferCreateRequest{
			EditionUid: uidEdition1, SourceWarehouseUid: uidWarehouse1, TargetWarehouseUid: uidWarehouse2, Quantity: 3,
		})
		return err
	}))
}

func TestAuditController_WhenCreateHitsMissingWarehouse_Returns404(t *testing.T) {
	use := &stubUseCase{
		createAudit: func(context.Context, command.StockAuditCreateCommand) (*aggregate.StockAuditAggregate, error) {
			return nil, warehouseNotFoundErr(uidWarehouse1)
		},
	}
	ctrl := NewAuditController(use)
	assertWarehouseNotFoundResponse(t, serveController(t, func(ctx context.Context) error {
		_, err := ctrl.AuditsCreate(ctx, openapi.AuditCreateRequest{
			WarehouseUid: uidWarehouse1,
			Items:        []openapi.AuditItemRequest{{EditionUid: uidEdition1, ActualQuantity: 2}},
		})
		return err
	}))
}

func TestClosingController_WhenCreateHitsMissingWarehouse_Returns404(t *testing.T) {
	use := &stubUseCase{
		createClosing: func(context.Context, command.MonthlyClosingCommand) (*aggregate.MonthlyClosingAggregate, error) {
			return nil, warehouseNotFoundErr(uidWarehouse1)
		},
	}
	ctrl := NewClosingController(use)
	assertWarehouseNotFoundResponse(t, serveController(t, func(ctx context.Context) error {
		_, err := ctrl.ClosingsCreate(ctx, openapi.ClosingCreateRequest{
			WarehouseUid: uidWarehouse1, Year: 2026, Month: 7,
		})
		return err
	}))
}
