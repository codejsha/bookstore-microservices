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

func TestClassifyReservationError_WhenErrorGiven_ReturnsRetryClassification(t *testing.T) {
	cases := []struct {
		name        string
		in          error
		wantAppErr  bool
		wantType    string
		nonRetryabl bool
	}{
		{
			name:        "whenStockInsufficient_returnsNonRetryable",
			in:          fmt.Errorf("no warehouse: %w", repo.ErrInsufficientStock),
			wantAppErr:  true,
			wantType:    errTypeInsufficientStock,
			nonRetryabl: true,
		},
		{
			name:        "whenQuantityInvalid_returnsNonRetryable",
			in:          fmt.Errorf("bad qty: %w", repo.ErrInvalidQuantity),
			wantAppErr:  true,
			wantType:    errTypeInvalidQuantity,
			nonRetryabl: true,
		},
		{
			name:       "whenDbErrorTransient_returnsRetryable",
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

func TestReserveStock_WhenStockInsufficient_ReturnsNonRetryableApplicationError(t *testing.T) {
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

func TestReserveStock_WhenDbErrorTransient_ReturnsRetryableError(t *testing.T) {
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

func TestClassifyReservationError_WhenCommandInvalid_ReturnsNonRetryableApplicationError(t *testing.T) {
	got := classifyReservationError(fmt.Errorf("order_uid must not be blank: %w", command.ErrInvalidCommand))

	var appErr *temporal.ApplicationError
	if !errors.As(got, &appErr) {
		t.Fatalf("expected ApplicationError, got %T: %v", got, got)
	}
	if !appErr.NonRetryable() {
		t.Errorf("invalid command must be non-retryable")
	}
	if appErr.Type() != errTypeInvalidCommand {
		t.Errorf("type = %q, want %q", appErr.Type(), errTypeInvalidCommand)
	}
}

func TestReleaseStock_WhenCommandInvalid_ReturnsNonRetryableApplicationError(t *testing.T) {
	var suite testsuite.WorkflowTestSuite
	env := suite.NewTestActivityEnvironment()

	uc := &stubInventoryUC{
		releaseFn: func(_ context.Context, cmd command.StockOrderReleaseCommand) (*aggregate.StockAggregate, error) {
			return nil, cmd.Validate()
		},
	}
	act := NewStockActivities(uc)
	env.RegisterActivity(act.ReleaseStock)

	_, err := env.ExecuteActivity(act.ReleaseStock, ReleaseStockRequest{
		OrderUid: "order-3",
		Items:    []StockReservationItem{{ProductID: 0, Quantity: 1}},
	})
	if err == nil {
		t.Fatal("expected error")
	}

	var appErr *temporal.ApplicationError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected ApplicationError, got %T: %v", err, err)
	}
	if !appErr.NonRetryable() {
		t.Errorf("invalid release input must fail fast instead of retrying forever")
	}
	if appErr.Type() != errTypeInvalidCommand {
		t.Errorf("type = %q, want %q", appErr.Type(), errTypeInvalidCommand)
	}
}

func TestReleaseStock_WhenDbErrorTransient_ReturnsRetryableError(t *testing.T) {
	var suite testsuite.WorkflowTestSuite
	env := suite.NewTestActivityEnvironment()

	uc := &stubInventoryUC{
		releaseFn: func(context.Context, command.StockOrderReleaseCommand) (*aggregate.StockAggregate, error) {
			return nil, errors.New("dial tcp: connection refused")
		},
	}
	act := NewStockActivities(uc)
	env.RegisterActivity(act.ReleaseStock)

	_, err := env.ExecuteActivity(act.ReleaseStock, ReleaseStockRequest{
		OrderUid: "order-4",
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

// ─── Per-order-line aggregation ─────────────────────────────────────────────

func TestAggregateReservationItems_DuplicateEditions_SumsQuantitiesOnce(t *testing.T) {
	got := aggregateReservationItems([]StockReservationItem{
		{ProductID: 7, Quantity: 2},
		{ProductID: 9, Quantity: 1},
		{ProductID: 7, Quantity: 3},
	})

	want := []StockReservationItem{{ProductID: 7, Quantity: 5}, {ProductID: 9, Quantity: 1}}
	if len(got) != len(want) {
		t.Fatalf("items = %+v, want %+v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("item[%d] = %+v, want %+v", i, got[i], want[i])
		}
	}
}

func TestAggregateReservationItems_ItemOrderVaries_KeepsFirstSeenOrder(t *testing.T) {
	forward := aggregateReservationItems([]StockReservationItem{{ProductID: 7, Quantity: 2}, {ProductID: 9, Quantity: 1}})
	reverse := aggregateReservationItems([]StockReservationItem{{ProductID: 9, Quantity: 1}, {ProductID: 7, Quantity: 2}})

	if forward[0].ProductID != 7 || reverse[0].ProductID != 9 {
		t.Fatalf("first-seen order not preserved: forward=%+v reverse=%+v", forward, reverse)
	}
}

func TestReserveStock_TwoLinesSameEdition_ReservesCombinedQuantity(t *testing.T) {
	var suite testsuite.WorkflowTestSuite
	env := suite.NewTestActivityEnvironment()

	var cmds []command.StockOrderReserveCommand
	uc := &stubInventoryUC{
		reserveFn: func(_ context.Context, cmd command.StockOrderReserveCommand) (*aggregate.StockAggregate, error) {
			cmds = append(cmds, cmd)
			return nil, nil
		},
	}
	act := NewStockActivities(uc)
	env.RegisterActivity(act.ReserveStock)

	if _, err := env.ExecuteActivity(act.ReserveStock, ReserveStockRequest{
		OrderUid: "order-1",
		Items:    []StockReservationItem{{ProductID: 7, Quantity: 2}, {ProductID: 7, Quantity: 3}},
	}); err != nil {
		t.Fatalf("reserve: %v", err)
	}

	if len(cmds) != 1 {
		t.Fatalf("reserve calls = %d (%+v), want 1 aggregated call", len(cmds), cmds)
	}
	if cmds[0].Quantity != 5 {
		t.Errorf("quantity = %d, want 5 (both lines), otherwise the second line is silently dropped", cmds[0].Quantity)
	}
	if cmds[0].EditionId != 7 || cmds[0].OrderUid != "order-1" {
		t.Errorf("cmd = %+v, want edition 7 of order-1", cmds[0])
	}
}

func TestReserveStock_DistinctEditions_ReservesEachSeparately(t *testing.T) {
	var suite testsuite.WorkflowTestSuite
	env := suite.NewTestActivityEnvironment()

	var cmds []command.StockOrderReserveCommand
	uc := &stubInventoryUC{
		reserveFn: func(_ context.Context, cmd command.StockOrderReserveCommand) (*aggregate.StockAggregate, error) {
			cmds = append(cmds, cmd)
			return nil, nil
		},
	}
	act := NewStockActivities(uc)
	env.RegisterActivity(act.ReserveStock)

	if _, err := env.ExecuteActivity(act.ReserveStock, ReserveStockRequest{
		OrderUid: "order-1",
		Items:    []StockReservationItem{{ProductID: 7, Quantity: 2}, {ProductID: 9, Quantity: 4}},
	}); err != nil {
		t.Fatalf("reserve: %v", err)
	}

	if len(cmds) != 2 {
		t.Fatalf("reserve calls = %d, want 2 (one per edition)", len(cmds))
	}
	if cmds[0].EditionId != 7 || cmds[0].Quantity != 2 || cmds[1].EditionId != 9 || cmds[1].Quantity != 4 {
		t.Errorf("cmds = %+v, want edition 7 qty 2 then edition 9 qty 4", cmds)
	}
}

func TestReserveStock_SameRequestRetried_IssuesIdenticalCommands(t *testing.T) {
	var suite testsuite.WorkflowTestSuite
	env := suite.NewTestActivityEnvironment()

	var cmds []command.StockOrderReserveCommand
	uc := &stubInventoryUC{
		reserveFn: func(_ context.Context, cmd command.StockOrderReserveCommand) (*aggregate.StockAggregate, error) {
			cmds = append(cmds, cmd)
			return nil, nil
		},
	}
	act := NewStockActivities(uc)
	env.RegisterActivity(act.ReserveStock)

	req := ReserveStockRequest{
		OrderUid: "order-1",
		Items:    []StockReservationItem{{ProductID: 7, Quantity: 2}, {ProductID: 7, Quantity: 3}},
	}
	for attempt := 0; attempt < 2; attempt++ {
		if _, err := env.ExecuteActivity(act.ReserveStock, req); err != nil {
			t.Fatalf("reserve attempt %d: %v", attempt, err)
		}
	}

	if len(cmds) != 2 {
		t.Fatalf("reserve calls = %d, want 2 (one per attempt)", len(cmds))
	}
	if cmds[0].OrderUid != cmds[1].OrderUid || cmds[0].EditionId != cmds[1].EditionId ||
		cmds[0].Quantity != cmds[1].Quantity {
		t.Errorf("retry issued a different reservation: %+v then %+v", cmds[0], cmds[1])
	}
	if cmds[1].Quantity != 5 {
		t.Errorf("retry quantity = %d, want 5 (same aggregated line, so the repo replay is a no-op)", cmds[1].Quantity)
	}
}

func TestReleaseStock_TwoLinesSameEdition_ReleasesCombinedQuantity(t *testing.T) {
	var suite testsuite.WorkflowTestSuite
	env := suite.NewTestActivityEnvironment()

	var cmds []command.StockOrderReleaseCommand
	uc := &stubInventoryUC{
		releaseFn: func(_ context.Context, cmd command.StockOrderReleaseCommand) (*aggregate.StockAggregate, error) {
			cmds = append(cmds, cmd)
			return nil, nil
		},
	}
	act := NewStockActivities(uc)
	env.RegisterActivity(act.ReleaseStock)

	if _, err := env.ExecuteActivity(act.ReleaseStock, ReleaseStockRequest{
		OrderUid: "order-1",
		Items:    []StockReservationItem{{ProductID: 7, Quantity: 2}, {ProductID: 7, Quantity: 3}},
	}); err != nil {
		t.Fatalf("release: %v", err)
	}

	if len(cmds) != 1 {
		t.Fatalf("release calls = %d (%+v), want 1 aggregated call matching the reservation", len(cmds), cmds)
	}
	if cmds[0].Quantity != 5 || cmds[0].EditionId != 7 {
		t.Errorf("cmd = %+v, want edition 7 qty 5 so release mirrors reserve", cmds[0])
	}
}
