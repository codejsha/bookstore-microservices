package service

import (
	"context"
	"errors"
	"testing"

	"github.com/codejsha/bookstore-microservices/inventory/internal/application/port/repo"
	"github.com/codejsha/bookstore-microservices/inventory/internal/domain/model/command"
	"github.com/codejsha/bookstore-microservices/inventory/internal/domain/model/option"
)

// ─── stub repos ─────────────────────────────────────────────────────────────

type stubTransferRepo struct {
	findAllFn   func(context.Context, option.TransferQueryOption) (int64, []*repo.TransferResult, error)
	findByUidFn func(context.Context, string) (*repo.TransferResult, error)
	createFn    func(context.Context, repo.TransferCreateParams) (*repo.TransferResult, error)
	completeFn  func(context.Context, string) (*repo.TransferResult, error)
	cancelFn    func(context.Context, string) (*repo.TransferResult, error)
}

func (s *stubTransferRepo) FindAll(ctx context.Context, opt option.TransferQueryOption) (int64, []*repo.TransferResult, error) {
	return s.findAllFn(ctx, opt)
}
func (s *stubTransferRepo) FindByUid(ctx context.Context, uid string) (*repo.TransferResult, error) {
	return s.findByUidFn(ctx, uid)
}
func (s *stubTransferRepo) Create(ctx context.Context, p repo.TransferCreateParams) (*repo.TransferResult, error) {
	return s.createFn(ctx, p)
}
func (s *stubTransferRepo) Complete(ctx context.Context, uid string) (*repo.TransferResult, error) {
	return s.completeFn(ctx, uid)
}
func (s *stubTransferRepo) Cancel(ctx context.Context, uid string) (*repo.TransferResult, error) {
	return s.cancelFn(ctx, uid)
}

type stubAuditRepo struct {
	findAllFn   func(context.Context, option.AuditQueryOption) (int64, []*repo.AuditResult, error)
	findByUidFn func(context.Context, string) (*repo.AuditResult, error)
	createFn    func(context.Context, repo.AuditCreateParams) (*repo.AuditResult, error)
	completeFn  func(context.Context, string) (*repo.AuditResult, error)
}

func (s *stubAuditRepo) FindAll(ctx context.Context, opt option.AuditQueryOption) (int64, []*repo.AuditResult, error) {
	return s.findAllFn(ctx, opt)
}
func (s *stubAuditRepo) FindByUid(ctx context.Context, uid string) (*repo.AuditResult, error) {
	return s.findByUidFn(ctx, uid)
}
func (s *stubAuditRepo) Create(ctx context.Context, p repo.AuditCreateParams) (*repo.AuditResult, error) {
	return s.createFn(ctx, p)
}
func (s *stubAuditRepo) Complete(ctx context.Context, uid string) (*repo.AuditResult, error) {
	return s.completeFn(ctx, uid)
}

type stubClosingRepo struct {
	findAllFn   func(context.Context, option.ClosingQueryOption) (int64, []*repo.ClosingResult, error)
	findByUidFn func(context.Context, string) (*repo.ClosingResult, error)
	createFn    func(context.Context, repo.ClosingCreateParams) (*repo.ClosingResult, error)
}

func (s *stubClosingRepo) FindAll(ctx context.Context, opt option.ClosingQueryOption) (int64, []*repo.ClosingResult, error) {
	return s.findAllFn(ctx, opt)
}
func (s *stubClosingRepo) FindByUid(ctx context.Context, uid string) (*repo.ClosingResult, error) {
	return s.findByUidFn(ctx, uid)
}
func (s *stubClosingRepo) Create(ctx context.Context, p repo.ClosingCreateParams) (*repo.ClosingResult, error) {
	return s.createFn(ctx, p)
}

type stubBalanceRepo struct {
	findAllFn func(context.Context, option.BalanceQueryOption) (int64, []*repo.BalanceResult, error)
}

func (s *stubBalanceRepo) FindAll(ctx context.Context, opt option.BalanceQueryOption) (int64, []*repo.BalanceResult, error) {
	return s.findAllFn(ctx, opt)
}

// ─── Transfer tests ─────────────────────────────────────────────────────────

func TestTransferStock_WhenCommandValid_CreatesTransfer(t *testing.T) {
	var captured repo.TransferCreateParams
	tr := &stubTransferRepo{
		createFn: func(_ context.Context, p repo.TransferCreateParams) (*repo.TransferResult, error) {
			captured = p
			return &repo.TransferResult{Uid: uidTransfer1, EditionUid: p.EditionUid, Status: "PENDING", Quantity: p.Quantity}, nil
		},
	}
	svc := &stockService{transferRepo: tr}
	got, err := svc.TransferStock(context.Background(), command.StockTransferCommand{
		EditionUid: uidEdition1, SourceWarehouseUid: uidWarehouse1, TargetWarehouseUid: uidWarehouse2, Quantity: 5,
	})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if captured.EditionUid != uidEdition1 || captured.SourceWarehouseUid != uidWarehouse1 || captured.TargetWarehouseUid != uidWarehouse2 || captured.Quantity != 5 {
		t.Errorf("captured = %+v", captured)
	}
	if got == nil || got.Uid != uidTransfer1 || got.Status != "PENDING" {
		t.Errorf("got = %+v", got)
	}
}

func TestTransferStock_WhenQuantityNotPositive_ReturnsErrInvalidCommand(t *testing.T) {
	svc := &stockService{transferRepo: &stubTransferRepo{}}
	if _, err := svc.TransferStock(context.Background(), command.StockTransferCommand{
		EditionUid: uidEdition3, SourceWarehouseUid: uidWarehouse1, TargetWarehouseUid: uidWarehouse2, Quantity: 0,
	}); err == nil {
		t.Errorf("expected error for non-positive transfer quantity")
	}
}

func TestTransferStock_WhenWarehousesMatch_ReturnsErrInvalidCommand(t *testing.T) {
	svc := &stockService{transferRepo: &stubTransferRepo{}}
	if _, err := svc.TransferStock(context.Background(), command.StockTransferCommand{
		EditionUid: uidEdition3, SourceWarehouseUid: uidWarehouse1, TargetWarehouseUid: uidWarehouse1, Quantity: 3,
	}); err == nil {
		t.Errorf("expected error when source and target warehouses match")
	}
}

func TestFindTransfer_WhenTransferMissing_ReturnsNotFoundError(t *testing.T) {
	tr := &stubTransferRepo{findByUidFn: func(context.Context, string) (*repo.TransferResult, error) { return nil, nil }}
	svc := &stockService{transferRepo: tr}
	got, err := svc.FindTransfer(context.Background(), "missing")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if got != nil {
		t.Errorf("got = %+v, want nil", got)
	}
}

func TestCompleteTransfer_WhenTransferPending_MovesStockAndCompletes(t *testing.T) {
	var completed string
	tr := &stubTransferRepo{
		findByUidFn: func(_ context.Context, uid string) (*repo.TransferResult, error) {
			return &repo.TransferResult{Uid: uid, Status: "PENDING"}, nil
		},
		completeFn: func(_ context.Context, uid string) (*repo.TransferResult, error) {
			completed = uid
			return &repo.TransferResult{Uid: uid, Status: "COMPLETED"}, nil
		},
	}
	svc := &stockService{transferRepo: tr}
	got, err := svc.CompleteTransfer(context.Background(), uidTransfer1)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if completed != uidTransfer1 || got == nil || got.Status != "COMPLETED" {
		t.Errorf("completed=%q got=%+v", completed, got)
	}
}

func TestCompleteTransfer_WhenTransferNotPending_ReturnsStateError(t *testing.T) {
	tr := &stubTransferRepo{
		findByUidFn: func(_ context.Context, uid string) (*repo.TransferResult, error) {
			return &repo.TransferResult{Uid: uid, Status: "CANCELLED"}, nil
		},
		completeFn: func(context.Context, string) (*repo.TransferResult, error) {
			t.Errorf("Complete must not be called for a non-completable transfer")
			return nil, nil
		},
	}
	svc := &stockService{transferRepo: tr}
	if _, err := svc.CompleteTransfer(context.Background(), uidTransfer1); err == nil {
		t.Errorf("expected error completing a CANCELLED transfer")
	}
}

func TestCancelTransfer_WhenTransferNotPending_ReturnsStateError(t *testing.T) {
	tr := &stubTransferRepo{
		findByUidFn: func(_ context.Context, uid string) (*repo.TransferResult, error) {
			return &repo.TransferResult{Uid: uid, Status: "IN_TRANSIT"}, nil
		},
		cancelFn: func(context.Context, string) (*repo.TransferResult, error) {
			t.Errorf("Cancel must not be called for a non-cancellable transfer")
			return nil, nil
		},
	}
	svc := &stockService{transferRepo: tr}
	if _, err := svc.CancelTransfer(context.Background(), uidTransfer1); err == nil {
		t.Errorf("expected error cancelling an IN_TRANSIT transfer")
	}
}

func TestCancelTransfer_WhenTransferMissing_ReturnsNotFoundError(t *testing.T) {
	tr := &stubTransferRepo{findByUidFn: func(context.Context, string) (*repo.TransferResult, error) { return nil, nil }}
	svc := &stockService{transferRepo: tr}
	got, err := svc.CancelTransfer(context.Background(), "missing")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if got != nil {
		t.Errorf("got = %+v, want nil (404)", got)
	}
}

// ─── Audit tests ────────────────────────────────────────────────────────────

func TestCreateAudit_WhenItemsGiven_ComputesDifferencesFromThem(t *testing.T) {
	var captured repo.AuditCreateParams
	ar := &stubAuditRepo{
		createFn: func(_ context.Context, p repo.AuditCreateParams) (*repo.AuditResult, error) {
			captured = p
			return &repo.AuditResult{
				Uid: uidAudit1, WarehouseUid: p.WarehouseUid, Status: "IN_PROGRESS",
				Items: []repo.AuditItemResult{{EditionUid: uidEdition1, SystemQuantity: 10, ActualQuantity: 8, Difference: -2}},
			}, nil
		},
	}
	svc := &stockService{auditRepo: ar}
	got, err := svc.CreateAudit(context.Background(), command.StockAuditCreateCommand{
		WarehouseUid: uidWarehouse1,
		Items:        []command.StockAuditItemCommand{{EditionUid: uidEdition1, ActualQuantity: 8}},
	})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if captured.WarehouseUid != uidWarehouse1 || len(captured.Items) != 1 || captured.Items[0].ActualQuantity != 8 {
		t.Errorf("captured = %+v", captured)
	}
	if got == nil || len(got.Items) != 1 || got.Items[0].Difference != -2 {
		t.Errorf("got = %+v", got)
	}
}

func TestCreateAudit_WhenItemsEmpty_ReturnsErrInvalidCommand(t *testing.T) {
	svc := &stockService{auditRepo: &stubAuditRepo{}}
	if _, err := svc.CreateAudit(context.Background(), command.StockAuditCreateCommand{WarehouseUid: uidWarehouse1}); err == nil {
		t.Errorf("expected error for empty audit items")
	}
}

func TestCompleteAudit_WhenAuditAlreadyCompleted_ReturnsStateError(t *testing.T) {
	ar := &stubAuditRepo{
		findByUidFn: func(_ context.Context, uid string) (*repo.AuditResult, error) {
			return &repo.AuditResult{Uid: uid, Status: "COMPLETED"}, nil
		},
		completeFn: func(context.Context, string) (*repo.AuditResult, error) {
			t.Errorf("Complete must not be called for an already-completed audit")
			return nil, nil
		},
	}
	svc := &stockService{auditRepo: ar}
	if _, err := svc.CompleteAudit(context.Background(), uidAudit1); err == nil {
		t.Errorf("expected error completing an already-completed audit")
	}
}

func TestCompleteAudit_WhenAuditOpen_AppliesAdjustmentsAndCompletes(t *testing.T) {
	ar := &stubAuditRepo{
		findByUidFn: func(_ context.Context, uid string) (*repo.AuditResult, error) {
			return &repo.AuditResult{Uid: uid, Status: "IN_PROGRESS"}, nil
		},
		completeFn: func(_ context.Context, uid string) (*repo.AuditResult, error) {
			return &repo.AuditResult{Uid: uid, Status: "COMPLETED"}, nil
		},
	}
	svc := &stockService{auditRepo: ar}
	got, err := svc.CompleteAudit(context.Background(), uidAudit1)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if got == nil || got.Status != "COMPLETED" {
		t.Errorf("got = %+v", got)
	}
}

// ─── Closing tests ──────────────────────────────────────────────────────────

func TestCreateMonthlyClosing_WhenPeriodValid_CreatesClosing(t *testing.T) {
	var captured repo.ClosingCreateParams
	cr := &stubClosingRepo{
		createFn: func(_ context.Context, p repo.ClosingCreateParams) (*repo.ClosingResult, error) {
			captured = p
			return &repo.ClosingResult{
				Uid: uidClosing1, WarehouseUid: p.WarehouseUid, Year: p.Year, Month: p.Month, Status: "CLOSED",
				Items: []repo.ClosingItemResult{{EditionUid: uidEdition1, OpeningQuantity: 3, InboundQuantity: 5, OutboundQuantity: 2, ClosingQuantity: 6}},
			}, nil
		},
	}
	svc := &stockService{closingRepo: cr}
	got, err := svc.CreateMonthlyClosing(context.Background(), command.MonthlyClosingCommand{
		WarehouseUid: uidWarehouse1, Year: 2026, Month: 7,
	})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if captured.Year != 2026 || captured.Month != 7 {
		t.Errorf("captured = %+v", captured)
	}
	if got == nil || got.Status != "CLOSED" || len(got.Items) != 1 || got.Items[0].ClosingQuantity != 6 {
		t.Errorf("got = %+v", got)
	}
}

func TestCreateMonthlyClosing_WhenMonthOutOfRange_ReturnsErrInvalidCommand(t *testing.T) {
	svc := &stockService{closingRepo: &stubClosingRepo{}}
	if _, err := svc.CreateMonthlyClosing(context.Background(), command.MonthlyClosingCommand{
		WarehouseUid: uidWarehouse1, Year: 2026, Month: 13,
	}); err == nil {
		t.Errorf("expected error for month 13")
	}
}

func TestFindClosing_WhenClosingMissing_ReturnsNotFoundError(t *testing.T) {
	cr := &stubClosingRepo{findByUidFn: func(context.Context, string) (*repo.ClosingResult, error) { return nil, nil }}
	svc := &stockService{closingRepo: cr}
	got, err := svc.FindClosing(context.Background(), "missing")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if got != nil {
		t.Errorf("got = %+v, want nil", got)
	}
}

// ─── Balance tests ──────────────────────────────────────────────────────────

func TestFindStockBalance_WhenRepoReturnsRows_MapsThemToAggregates(t *testing.T) {
	br := &stubBalanceRepo{
		findAllFn: func(_ context.Context, _ option.BalanceQueryOption) (int64, []*repo.BalanceResult, error) {
			return 1, []*repo.BalanceResult{{
				EditionUid: uidEdition1, WarehouseUid: uidWarehouse1, WarehouseName: "Main",
				InboundQuantity: 10, OutboundQuantity: 4, AdjustQuantity: -1, CurrentQuantity: 5,
			}}, nil
		},
	}
	svc := &stockService{balanceRepo: br}
	total, entries, err := svc.FindStockBalance(context.Background(), option.NewBalanceQueryOption())
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if total != 1 || len(entries) != 1 {
		t.Fatalf("entries = %+v (total=%d)", entries, total)
	}
	if entries[0].WarehouseName != "Main" || entries[0].CurrentQuantity != 5 || entries[0].OutboundQuantity != 4 {
		t.Errorf("entries[0] = %+v", entries[0])
	}
}

func TestFindStockBalance_WhenRepoFails_ReturnsRepoError(t *testing.T) {
	wantErr := errors.New("report failed")
	br := &stubBalanceRepo{
		findAllFn: func(context.Context, option.BalanceQueryOption) (int64, []*repo.BalanceResult, error) {
			return 0, nil, wantErr
		},
	}
	svc := &stockService{balanceRepo: br}
	if _, _, err := svc.FindStockBalance(context.Background(), option.NewBalanceQueryOption()); !errors.Is(err, wantErr) {
		t.Errorf("err = %v, want %v", err, wantErr)
	}
}
