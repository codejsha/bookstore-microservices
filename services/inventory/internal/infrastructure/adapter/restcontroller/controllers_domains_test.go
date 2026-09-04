package restcontroller

import (
	"context"
	"errors"
	"testing"

	"github.com/codejsha/bookstore-microservices/inventory/generated/application/port/openapi"
	"github.com/codejsha/bookstore-microservices/inventory/internal/domain/aggregate"
	"github.com/codejsha/bookstore-microservices/inventory/internal/domain/model/command"
	"github.com/codejsha/bookstore-microservices/inventory/internal/domain/model/option"
)

// ─── Transfer controller ──────────────────────────────────────────────────

func TestTransferController_WhenFiltersGiven_ReturnsPagedTransfers(t *testing.T) {
	use := &stubUseCase{
		findAllTransfers: func(_ context.Context, opt option.TransferQueryOption) (int64, []*aggregate.StockTransferAggregate, error) {
			if opt.Status() == nil || *opt.Status() != "PENDING" {
				t.Errorf("Status = %v, want PENDING", opt.Status())
			}
			return 1, []*aggregate.StockTransferAggregate{
				{Uid: uidTransfer1, EditionUid: uidEdition1, SourceWarehouseUid: uidWarehouse1, TargetWarehouseUid: uidWarehouse2, Quantity: 4, Status: aggregate.TransferStatusPending},
			}, nil
		},
	}
	ctrl := NewTransferController(use)
	status := openapi.TRANSFERSTATUS_PENDING
	resp, err := ctrl.TransfersGetAll(context.Background(), nil, nil, nil, &status, nil, nil, nil)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp.Total != 1 || len(resp.Items) != 1 || resp.Items[0].Status != openapi.TRANSFERSTATUS_PENDING {
		t.Errorf("resp = %+v", resp)
	}
}

func TestTransferController_WhenRequestValid_PassesCommandToUsecase(t *testing.T) {
	captured := command.StockTransferCommand{}
	use := &stubUseCase{
		transferStock: func(_ context.Context, cmd command.StockTransferCommand) (*aggregate.StockTransferAggregate, error) {
			captured = cmd
			return &aggregate.StockTransferAggregate{Uid: uidTransfer1, EditionUid: cmd.EditionUid, Status: aggregate.TransferStatusPending, Quantity: cmd.Quantity}, nil
		},
	}
	ctrl := NewTransferController(use)
	resp, err := ctrl.TransfersCreate(context.Background(), openapi.TransferCreateRequest{
		EditionUid: uidEdition1, SourceWarehouseUid: uidWarehouse1, TargetWarehouseUid: uidWarehouse2, Quantity: 6,
	})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if captured.SourceWarehouseUid != uidWarehouse1 || captured.TargetWarehouseUid != uidWarehouse2 || captured.Quantity != 6 {
		t.Errorf("captured = %+v", captured)
	}
	if resp.Status != openapi.TRANSFERSTATUS_PENDING {
		t.Errorf("resp = %+v", resp)
	}
}

func TestTransferController_WhenTransferMissing_ReturnsNotFoundError(t *testing.T) {
	use := &stubUseCase{
		findTransfer: func(context.Context, string) (*aggregate.StockTransferAggregate, error) { return nil, nil },
	}
	ctrl := NewTransferController(use)
	if _, err := ctrl.TransfersRead(context.Background(), uidMissing); err == nil {
		t.Errorf("expected not-found error for missing transfer")
	}
}

func TestTransferController_WhenTransferPending_ReturnsCompletedTransfer(t *testing.T) {
	use := &stubUseCase{
		completeTransfer: func(_ context.Context, uid string) (*aggregate.StockTransferAggregate, error) {
			return &aggregate.StockTransferAggregate{Uid: uid, Status: aggregate.TransferStatusCompleted}, nil
		},
	}
	ctrl := NewTransferController(use)
	resp, err := ctrl.TransfersComplete(context.Background(), uidTransfer1)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp.Status != openapi.TRANSFERSTATUS_COMPLETED {
		t.Errorf("resp = %+v", resp)
	}
}

// ─── Audit controller ─────────────────────────────────────────────────────

func TestAuditController_WhenRequestValid_PassesCommandToUsecase(t *testing.T) {
	captured := command.StockAuditCreateCommand{}
	use := &stubUseCase{
		createAudit: func(_ context.Context, cmd command.StockAuditCreateCommand) (*aggregate.StockAuditAggregate, error) {
			captured = cmd
			return &aggregate.StockAuditAggregate{
				Uid: uidAudit1, WarehouseUid: cmd.WarehouseUid, Status: aggregate.AuditStatusInProgress,
				Items: []*aggregate.StockAuditItem{{EditionUid: uidEdition1, SystemQuantity: 10, ActualQuantity: 7, Difference: -3}},
			}, nil
		},
	}
	ctrl := NewAuditController(use)
	resp, err := ctrl.AuditsCreate(context.Background(), openapi.AuditCreateRequest{
		WarehouseUid: uidWarehouse1,
		Items:        []openapi.AuditItemRequest{{EditionUid: uidEdition1, ActualQuantity: 7}},
	})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if captured.WarehouseUid != uidWarehouse1 || len(captured.Items) != 1 || captured.Items[0].ActualQuantity != 7 {
		t.Errorf("captured = %+v", captured)
	}
	if resp.Status != openapi.AUDITSTATUS_IN_PROGRESS || len(resp.Items) != 1 || resp.Items[0].Difference != -3 {
		t.Errorf("resp = %+v", resp)
	}
}

func TestAuditController_WhenAuditMissing_ReturnsNotFoundError(t *testing.T) {
	use := &stubUseCase{
		findAudit: func(context.Context, string) (*aggregate.StockAuditAggregate, error) { return nil, nil },
	}
	ctrl := NewAuditController(use)
	if _, err := ctrl.AuditsRead(context.Background(), uidMissing); err == nil {
		t.Errorf("expected not-found error for missing audit")
	}
}

// ─── Closing controller ───────────────────────────────────────────────────

func TestClosingController_WhenRequestValid_PassesCommandToUsecase(t *testing.T) {
	captured := command.MonthlyClosingCommand{}
	use := &stubUseCase{
		createClosing: func(_ context.Context, cmd command.MonthlyClosingCommand) (*aggregate.MonthlyClosingAggregate, error) {
			captured = cmd
			return &aggregate.MonthlyClosingAggregate{
				Uid: uidClosing1, WarehouseUid: cmd.WarehouseUid, Year: cmd.Year, Month: cmd.Month, Status: aggregate.ClosingStatusClosed,
				Items: []*aggregate.MonthlyClosingItem{{EditionUid: uidEdition1, OpeningQuantity: 2, InboundQuantity: 5, OutboundQuantity: 3, ClosingQuantity: 4}},
			}, nil
		},
	}
	ctrl := NewClosingController(use)
	resp, err := ctrl.ClosingsCreate(context.Background(), openapi.ClosingCreateRequest{
		WarehouseUid: uidWarehouse1, Year: 2026, Month: 7,
	})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if captured.Year != 2026 || captured.Month != 7 {
		t.Errorf("captured = %+v", captured)
	}
	if resp.Status != openapi.CLOSINGSTATUS_CLOSED || len(resp.Items) != 1 || resp.Items[0].ClosingQuantity != 4 {
		t.Errorf("resp = %+v", resp)
	}
}

func TestClosingController_WhenClosingMissing_ReturnsNotFoundError(t *testing.T) {
	use := &stubUseCase{
		findClosing: func(context.Context, string) (*aggregate.MonthlyClosingAggregate, error) { return nil, nil },
	}
	ctrl := NewClosingController(use)
	if _, err := ctrl.ClosingsRead(context.Background(), uidMissing); err == nil {
		t.Errorf("expected not-found error for missing closing")
	}
}

// ─── Balance controller ───────────────────────────────────────────────────

func TestBalanceController_WhenFiltersGiven_ReturnsPagedBalances(t *testing.T) {
	use := &stubUseCase{
		findStockBalance: func(_ context.Context, opt option.BalanceQueryOption) (int64, []*aggregate.StockBalanceEntry, error) {
			if opt.YearMonth() == nil || *opt.YearMonth() != "2026-07" {
				t.Errorf("YearMonth = %v, want 2026-07", opt.YearMonth())
			}
			return 1, []*aggregate.StockBalanceEntry{
				{EditionUid: uidEdition1, WarehouseUid: uidWarehouse1, WarehouseName: "Main", InboundQuantity: 8, OutboundQuantity: 2, AdjustQuantity: 1, CurrentQuantity: 7},
			}, nil
		},
	}
	ctrl := NewBalanceController(use)
	ym := "2026-07"
	resp, err := ctrl.BalanceGetAll(context.Background(), nil, nil, &ym, nil, nil, nil)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp.Total != 1 || len(resp.Items) != 1 || resp.Items[0].CurrentQuantity != 7 || resp.Items[0].WarehouseName != "Main" {
		t.Errorf("resp = %+v", resp)
	}
}

func TestBalanceController_WhenUsecaseFails_PropagatesError(t *testing.T) {
	wantErr := errors.New("report boom")
	use := &stubUseCase{
		findStockBalance: func(context.Context, option.BalanceQueryOption) (int64, []*aggregate.StockBalanceEntry, error) {
			return 0, nil, wantErr
		},
	}
	ctrl := NewBalanceController(use)
	if _, err := ctrl.BalanceGetAll(context.Background(), nil, nil, nil, nil, nil, nil); !errors.Is(err, wantErr) {
		t.Errorf("err = %v, want %v", err, wantErr)
	}
}
