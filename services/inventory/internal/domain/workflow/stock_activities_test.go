package workflow

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/testsuite"

	"github.com/codejsha/bookstore-microservices/inventory/internal/application/port/repo"
	"github.com/codejsha/bookstore-microservices/inventory/internal/application/usecase"
	"github.com/codejsha/bookstore-microservices/inventory/internal/domain/aggregate"
	"github.com/codejsha/bookstore-microservices/inventory/internal/domain/model/command"
)

type stubInventoryUC struct {
	usecase.InventoryUseCase
	reserveFn func(context.Context, command.StockOrderReserveCommand) (*aggregate.StockAggregate, error)
	releaseFn func(context.Context, command.StockOrderReleaseCommand) (*aggregate.StockAggregate, error)
}

func (s *stubInventoryUC) ReserveStockForOrder(ctx context.Context, cmd command.StockOrderReserveCommand) (*aggregate.StockAggregate, error) {
	return s.reserveFn(ctx, cmd)
}

func (s *stubInventoryUC) ReleaseStockForOrder(ctx context.Context, cmd command.StockOrderReleaseCommand) (*aggregate.StockAggregate, error) {
	return s.releaseFn(ctx, cmd)
}

func TestClassifyReservationError(t *testing.T) {
	cases := []struct {
		name        string
		in          error
		wantAppErr  bool
		wantType    string
		nonRetryabl bool
	}{
		{
			name:        "insufficient stock is permanent",
			in:          fmt.Errorf("no warehouse: %w", repo.ErrInsufficientStock),
			wantAppErr:  true,
			wantType:    errTypeInsufficientStock,
			nonRetryabl: true,
		},
		{
			name:        "invalid quantity is permanent",
			in:          fmt.Errorf("bad qty: %w", repo.ErrInvalidQuantity),
			wantAppErr:  true,
			wantType:    errTypeInvalidQuantity,
			nonRetryabl: true,
		},
		{
			name:       "transient DB error stays retryable",
			in:         errors.New("dial tcp: connection refused"),
			wantAppErr: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := classifyReservationError(tc.in)

			var appErr *temporal.ApplicationError
			isAppErr := errors.As(got, &appErr)
			if isAppErr != tc.wantAppErr {
				t.Fatalf("errors.As ApplicationError = %v, want %v (err=%v)", isAppErr, tc.wantAppErr, got)
			}
			if !tc.wantAppErr {
				if !errors.Is(got, tc.in) {
					t.Errorf("transient error not passed through: got %v", got)
				}
				return
			}
			if appErr.Type() != tc.wantType {
				t.Errorf("type = %q, want %q", appErr.Type(), tc.wantType)
			}
			if appErr.NonRetryable() != tc.nonRetryabl {
				t.Errorf("NonRetryable = %v, want %v", appErr.NonRetryable(), tc.nonRetryabl)
			}
			if !errors.Is(got, tc.in) {
				t.Errorf("cause lost: %v does not wrap %v", got, tc.in)
			}
		})
	}
}

func TestReserveStock_InsufficientIsNonRetryable(t *testing.T) {
	var suite testsuite.WorkflowTestSuite
	env := suite.NewTestActivityEnvironment()

	uc := &stubInventoryUC{
		reserveFn: func(context.Context, command.StockOrderReserveCommand) (*aggregate.StockAggregate, error) {
			return nil, fmt.Errorf("no warehouse holds enough: %w", repo.ErrInsufficientStock)
		},
	}
	act := NewStockActivities(uc)
	env.RegisterActivity(act.ReserveStock)

	_, err := env.ExecuteActivity(act.ReserveStock, ReserveStockRequest{
		OrderUid: "order-1",
		Items:    []StockReservationItem{{ProductID: 7, Quantity: 5}},
	})
	if err == nil {
		t.Fatal("expected error")
	}

	var appErr *temporal.ApplicationError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected ApplicationError, got %T: %v", err, err)
	}
	if !appErr.NonRetryable() {
		t.Errorf("insufficient-stock activity failure must be non-retryable so the saga compensates immediately")
	}
	if appErr.Type() != errTypeInsufficientStock {
		t.Errorf("type = %q, want %q", appErr.Type(), errTypeInsufficientStock)
	}
}

func TestReserveStock_TransientStaysRetryable(t *testing.T) {
	var suite testsuite.WorkflowTestSuite
	env := suite.NewTestActivityEnvironment()

	uc := &stubInventoryUC{
		reserveFn: func(context.Context, command.StockOrderReserveCommand) (*aggregate.StockAggregate, error) {
			return nil, errors.New("dial tcp: connection refused")
		},
	}
	act := NewStockActivities(uc)
	env.RegisterActivity(act.ReserveStock)

	_, err := env.ExecuteActivity(act.ReserveStock, ReserveStockRequest{
		OrderUid: "order-2",
		Items:    []StockReservationItem{{ProductID: 7, Quantity: 5}},
	})
	if err == nil {
		t.Fatal("expected error")
	}
	var appErr *temporal.ApplicationError
	if errors.As(err, &appErr) && appErr.NonRetryable() {
		t.Errorf("transient DB error must stay retryable, got non-retryable ApplicationError: %v", err)
	}
}
