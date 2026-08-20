package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/codejsha/shared-library-go/pkg/pagination"

	"github.com/codejsha/bookstore-microservices/inventory/internal/application/port/repo"
	"github.com/codejsha/bookstore-microservices/inventory/internal/domain/model/command"
	"github.com/codejsha/bookstore-microservices/inventory/internal/domain/model/option"
)

// ─── stub StockRepo ─────────────────────────────────────────────────────────

type stubStockRepo struct {
	findAllFn     func(ctx context.Context, opt option.StockQueryOption) (int64, []*repo.StockResult, error)
	findOneFn     func(ctx context.Context, id int64) (*repo.StockResult, error)
	applyChangeFn func(ctx context.Context, p repo.ApplyChangeParams) error
}

func (s *stubStockRepo) FindAll(ctx context.Context, opt option.StockQueryOption) (int64, []*repo.StockResult, error) {
	return s.findAllFn(ctx, opt)
}

func (s *stubStockRepo) FindOne(ctx context.Context, id int64) (*repo.StockResult, error) {
	return s.findOneFn(ctx, id)
}

func (s *stubStockRepo) ApplyChange(ctx context.Context, p repo.ApplyChangeParams) error {
	if s.applyChangeFn != nil {
		return s.applyChangeFn(ctx, p)
	}
	return nil
}

// ─── stub StockReservationRepo ──────────────────────────────────────────────

type stubReservationRepo struct {
	reserveFn func(ctx context.Context, p repo.ReserveParams) (*repo.ReservationResult, error)
	releaseFn func(ctx context.Context, p repo.ReleaseParams) (*repo.ReservationResult, error)
}

func (s *stubReservationRepo) Reserve(ctx context.Context, p repo.ReserveParams) (*repo.ReservationResult, error) {
	return s.reserveFn(ctx, p)
}

func (s *stubReservationRepo) Release(ctx context.Context, p repo.ReleaseParams) (*repo.ReservationResult, error) {
	return s.releaseFn(ctx, p)
}

// ─── stub WarehouseRepo / StockHistoryRepo ─────────────────────────────────

type stubWarehouseRepo struct {
	findAllFn func(ctx context.Context, opt option.WarehouseQueryOption) (int64, []*repo.WarehouseResult, error)
	findByUid func(ctx context.Context, uid string) (*repo.WarehouseResult, error)
	createFn  func(ctx context.Context, p repo.WarehouseCreateParams) (*repo.WarehouseResult, error)
	updateFn  func(ctx context.Context, p repo.WarehouseUpdateParams) (*repo.WarehouseResult, error)
}

func (s *stubWarehouseRepo) FindAll(ctx context.Context, opt option.WarehouseQueryOption) (int64, []*repo.WarehouseResult, error) {
	return s.findAllFn(ctx, opt)
}

func (s *stubWarehouseRepo) FindByUid(ctx context.Context, uid string) (*repo.WarehouseResult, error) {
	if s.findByUid != nil {
		return s.findByUid(ctx, uid)
	}
	return nil, nil
}

func (s *stubWarehouseRepo) Create(ctx context.Context, p repo.WarehouseCreateParams) (*repo.WarehouseResult, error) {
	if s.createFn != nil {
		return s.createFn(ctx, p)
	}
	return nil, nil
}

func (s *stubWarehouseRepo) Update(ctx context.Context, p repo.WarehouseUpdateParams) (*repo.WarehouseResult, error) {
	if s.updateFn != nil {
		return s.updateFn(ctx, p)
	}
	return nil, nil
}

type stubStockHistoryRepo struct {
	findAllFn func(ctx context.Context, stockUid string, page pagination.PageOption) (int64, []*repo.StockHistoryResult, error)
}

func (s *stubStockHistoryRepo) FindAll(ctx context.Context, stockUid string, page pagination.PageOption) (int64, []*repo.StockHistoryResult, error) {
	return s.findAllFn(ctx, stockUid, page)
}

// ─── Stock query tests ──────────────────────────────────────────────────────

func TestFindStock_Found(t *testing.T) {
	stockRepo := &stubStockRepo{
		findAllFn: func(_ context.Context, opt option.StockQueryOption) (int64, []*repo.StockResult, error) {
			if opt.EditionUid() == nil || *opt.EditionUid() != "ed-1" {
				t.Errorf("EditionUid = %v, want ed-1", opt.EditionUid())
			}
			return 2, []*repo.StockResult{
				{Id: 1, Uid: "stk-1", EditionUid: "ed-1", WarehouseUid: "wh-1", WarehouseName: "W1", Quantity: 10},
				{Id: 1, Uid: "stk-1", EditionUid: "ed-1", WarehouseUid: "wh-2", WarehouseName: "W2", Quantity: 5},
			}, nil
		},
	}
	svc := &stockService{stockRepo: stockRepo}
	got, err := svc.FindStock(context.Background(), "ed-1")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if got == nil {
		t.Fatal("got nil aggregate, want populated")
	}
	if got.EditionUid != "ed-1" {
		t.Errorf("EditionUid = %q, want ed-1", got.EditionUid)
	}
	if got.TotalQuantity != 15 {
		t.Errorf("TotalQuantity = %d, want 15", got.TotalQuantity)
	}
	if len(got.Warehouses) != 2 {
		t.Fatalf("Warehouses len = %d, want 2", len(got.Warehouses))
	}
	if got.Warehouses[0].WarehouseUid != "wh-1" || got.Warehouses[0].Quantity != 10 {
		t.Errorf("Warehouses[0] = %+v", got.Warehouses[0])
	}
}

func TestFindStock_NotFound(t *testing.T) {
	stockRepo := &stubStockRepo{
		findAllFn: func(context.Context, option.StockQueryOption) (int64, []*repo.StockResult, error) {
			return 0, nil, nil
		},
	}
	svc := &stockService{stockRepo: stockRepo}
	got, err := svc.FindStock(context.Background(), "missing")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if got != nil {
		t.Errorf("got = %+v, want nil", got)
	}
}

func TestFindStock_RepoError(t *testing.T) {
	wantErr := errors.New("db down")
	stockRepo := &stubStockRepo{
		findAllFn: func(context.Context, option.StockQueryOption) (int64, []*repo.StockResult, error) {
			return 0, nil, wantErr
		},
	}
	svc := &stockService{stockRepo: stockRepo}
	_, err := svc.FindStock(context.Background(), "ed-1")
	if !errors.Is(err, wantErr) {
		t.Errorf("err = %v, want %v", err, wantErr)
	}
}

func TestFindAllStocks_GroupsByEditionPreservesOrder(t *testing.T) {
	stockRepo := &stubStockRepo{
		findAllFn: func(context.Context, option.StockQueryOption) (int64, []*repo.StockResult, error) {
			return 2, []*repo.StockResult{
				{Id: 1, Uid: "stk-A", EditionUid: "ed-A", WarehouseUid: "wh-1", Quantity: 3},
				{Id: 2, Uid: "stk-B", EditionUid: "ed-B", WarehouseUid: "wh-1", Quantity: 7},
				{Id: 3, Uid: "stk-A", EditionUid: "ed-A", WarehouseUid: "wh-2", Quantity: 2},
				{Id: 4, Uid: "stk-B", EditionUid: "ed-B", WarehouseUid: "wh-2", Quantity: 1},
			}, nil
		},
	}
	svc := &stockService{stockRepo: stockRepo}
	total, aggs, err := svc.FindAllStocks(context.Background(), option.NewStockQueryOption())
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if total != 2 {
		t.Errorf("total = %d, want 2 (distinct-edition count)", total)
	}
	if len(aggs) != 2 {
		t.Fatalf("aggs len = %d, want 2", len(aggs))
	}
	if aggs[0].EditionUid != "ed-A" || aggs[0].TotalQuantity != 5 {
		t.Errorf("aggs[0] = %+v, want edA total=5", aggs[0])
	}
	if aggs[1].EditionUid != "ed-B" || aggs[1].TotalQuantity != 8 {
		t.Errorf("aggs[1] = %+v, want edB total=8", aggs[1])
	}
}

func TestFindAllStocks_RepoError(t *testing.T) {
	wantErr := errors.New("boom")
	stockRepo := &stubStockRepo{
		findAllFn: func(context.Context, option.StockQueryOption) (int64, []*repo.StockResult, error) {
			return 0, nil, wantErr
		},
	}
	svc := &stockService{stockRepo: stockRepo}
	_, _, err := svc.FindAllStocks(context.Background(), option.NewStockQueryOption())
	if !errors.Is(err, wantErr) {
		t.Errorf("err = %v, want %v", err, wantErr)
	}
}

func TestToStockAggregate_Empty(t *testing.T) {
	if got := toStockAggregate(nil); got != nil {
		t.Errorf("got = %+v, want nil", got)
	}
}

func TestToStockAggregate_SumsQuantities(t *testing.T) {
	got := toStockAggregate([]*repo.StockResult{
		{Uid: "stk", EditionUid: "ed", WarehouseUid: "w1", Quantity: 4},
		{Uid: "stk", EditionUid: "ed", WarehouseUid: "w2", Quantity: 6},
	})
	if got.TotalQuantity != 10 {
		t.Errorf("TotalQuantity = %d, want 10", got.TotalQuantity)
	}
	if len(got.Warehouses) != 2 {
		t.Errorf("Warehouses = %+v, want len 2", got.Warehouses)
	}
}

// ─── Warehouse CRUD tests ──────────────────────────────────────────────────

func TestRegisterWarehouse(t *testing.T) {
	var captured repo.WarehouseCreateParams
	wh := &stubWarehouseRepo{
		createFn: func(_ context.Context, p repo.WarehouseCreateParams) (*repo.WarehouseResult, error) {
			captured = p
			return &repo.WarehouseResult{Uid: "wh-new", Name: p.Name, Address: p.Address, Capacity: p.Capacity}, nil
		},
	}
	svc := &stockService{warehouseRepo: wh}
	addr := "Mars"
	got, err := svc.RegisterWarehouse(context.Background(), command.WarehouseCreateCommand{
		Name: "New", Address: &addr, Capacity: 500,
	})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if captured.Name != "New" || captured.Capacity != 500 {
		t.Fatalf("captured = %+v", captured)
	}
	if captured.Address == nil || *captured.Address != "Mars" {
		t.Errorf("Address = %v, want Mars", captured.Address)
	}
	if got.Name != "New" || got.Capacity != 500 {
		t.Errorf("got = %+v", got)
	}
	if got.Uid == "" {
		t.Errorf("Uid should be assigned by the repository")
	}
}

func TestUpdateWarehouse_AppliesPatch(t *testing.T) {
	var captured repo.WarehouseUpdateParams
	wh := &stubWarehouseRepo{
		updateFn: func(_ context.Context, p repo.WarehouseUpdateParams) (*repo.WarehouseResult, error) {
			captured = p
			return &repo.WarehouseResult{Uid: p.Uid, Name: *p.Name, Capacity: *p.Capacity}, nil
		},
	}
	svc := &stockService{warehouseRepo: wh}
	name := "New"
	cap := int32(200)
	got, err := svc.UpdateWarehouse(context.Background(), "wh-1", command.WarehouseUpdateCommand{
		Name: &name, Capacity: &cap,
	})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if captured.Uid != "wh-1" {
		t.Errorf("Uid = %q, want wh-1", captured.Uid)
	}
	if captured.Name == nil || *captured.Name != "New" {
		t.Errorf("Name = %v, want New", captured.Name)
	}
	if captured.Capacity == nil || *captured.Capacity != 200 {
		t.Errorf("Capacity = %v, want 200", captured.Capacity)
	}
	if got.Name != "New" || got.Capacity != 200 {
		t.Errorf("got = %+v", got)
	}
}

func TestUpdateWarehouse_NotFound(t *testing.T) {
	wh := &stubWarehouseRepo{
		updateFn: func(context.Context, repo.WarehouseUpdateParams) (*repo.WarehouseResult, error) {
			return nil, nil
		},
	}
	svc := &stockService{warehouseRepo: wh}
	got, err := svc.UpdateWarehouse(context.Background(), "wh-x", command.WarehouseUpdateCommand{})
	if err != nil {
		t.Fatalf("err = %v, want nil for missing warehouse", err)
	}
	if got != nil {
		t.Errorf("got = %+v, want nil aggregate for missing warehouse", got)
	}
}

func TestFindWarehouse(t *testing.T) {
	wh := &stubWarehouseRepo{
		findByUid: func(_ context.Context, uid string) (*repo.WarehouseResult, error) {
			return &repo.WarehouseResult{Uid: uid, Name: "Main", Capacity: 100}, nil
		},
	}
	svc := &stockService{warehouseRepo: wh}
	got, err := svc.FindWarehouse(context.Background(), "wh-1")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if got == nil || got.Name != "Main" {
		t.Errorf("got = %+v", got)
	}
}

func TestFindWarehouse_NotFound(t *testing.T) {
	wh := &stubWarehouseRepo{
		findByUid: func(context.Context, string) (*repo.WarehouseResult, error) { return nil, nil },
	}
	svc := &stockService{warehouseRepo: wh}
	got, err := svc.FindWarehouse(context.Background(), "missing")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if got != nil {
		t.Errorf("got = %+v, want nil", got)
	}
}

func TestFindAllWarehouses(t *testing.T) {
	wh := &stubWarehouseRepo{
		findAllFn: func(_ context.Context, opt option.WarehouseQueryOption) (int64, []*repo.WarehouseResult, error) {
			return 1, []*repo.WarehouseResult{{Uid: "wh-1", Name: "Main", Capacity: 100}}, nil
		},
	}
	svc := &stockService{warehouseRepo: wh}
	total, got, err := svc.FindAllWarehouses(context.Background(), option.NewWarehouseQueryOption())
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if total != 1 || len(got) != 1 || got[0].Name != "Main" {
		t.Errorf("got = %+v (total=%d)", got, total)
	}
}

// ─── Stock operation tests ────────────────────────────────────────────────

func newStockSvc(stk *stubStockRepo, hist *stubStockHistoryRepo) *stockService {
	if hist == nil {
		hist = &stubStockHistoryRepo{}
	}
	return &stockService{stockRepo: stk, stockHistoryRepo: hist}
}

func TestReceiveStock_IncrementsAndRecordsHistory(t *testing.T) {
	var applied *repo.ApplyChangeParams
	stk := &stubStockRepo{
		applyChangeFn: func(_ context.Context, p repo.ApplyChangeParams) error {
			applied = &p
			return nil
		},
		findAllFn: func(_ context.Context, _ option.StockQueryOption) (int64, []*repo.StockResult, error) {
			return 1, []*repo.StockResult{{Id: 7, Uid: "stk-1", EditionUid: "ed-1", WarehouseUid: "wh-1", Quantity: 8}}, nil
		},
	}
	svc := newStockSvc(stk, nil)
	reason := "shipment received"
	got, err := svc.ReceiveStock(context.Background(), command.StockReceiveCommand{
		EditionUid: "ed-1", WarehouseUid: "wh-1", Quantity: 3, Reason: &reason,
	})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if applied == nil || applied.Delta != 3 || applied.ChangeType != "INBOUND" {
		t.Errorf("applied = %+v", applied)
	}
	if applied.EditionUid != "ed-1" || applied.WarehouseUid != "wh-1" {
		t.Errorf("applied target = %+v", applied)
	}
	if got == nil || got.TotalQuantity != 8 {
		t.Errorf("got = %+v", got)
	}
}

func TestReserveStock_DecrementsAndRecordsReservation(t *testing.T) {
	var applied *repo.ApplyChangeParams
	stk := &stubStockRepo{
		applyChangeFn: func(_ context.Context, p repo.ApplyChangeParams) error { applied = &p; return nil },
		findAllFn: func(context.Context, option.StockQueryOption) (int64, []*repo.StockResult, error) {
			return 1, []*repo.StockResult{{Uid: "stk-1", EditionUid: "ed-1", WarehouseUid: "wh-1", Quantity: 7}}, nil
		},
	}
	svc := newStockSvc(stk, nil)
	if _, err := svc.ReserveStock(context.Background(), command.StockReserveCommand{
		EditionUid: "ed-1", WarehouseUid: "wh-1", Quantity: 3,
	}); err != nil {
		t.Fatalf("err: %v", err)
	}
	if applied == nil || applied.Delta != -3 || applied.ChangeType != "RESERVATION" {
		t.Errorf("applied = %+v", applied)
	}
}

func TestReserveStock_RejectsNonPositive(t *testing.T) {
	svc := newStockSvc(&stubStockRepo{}, nil)
	if _, err := svc.ReserveStock(context.Background(), command.StockReserveCommand{
		EditionUid: "ed", WarehouseUid: "wh", Quantity: 0,
	}); err == nil {
		t.Errorf("expected error for non-positive quantity")
	}
}

func TestReceiveStock_RejectsNonPositive(t *testing.T) {
	svc := newStockSvc(&stubStockRepo{}, nil)
	if _, err := svc.ReceiveStock(context.Background(), command.StockReceiveCommand{
		EditionUid: "ed", WarehouseUid: "wh", Quantity: 0,
	}); err == nil {
		t.Errorf("expected error for non-positive receive quantity")
	}
}

func TestReleaseStock_RejectsNonPositive(t *testing.T) {
	svc := newStockSvc(&stubStockRepo{}, nil)
	if _, err := svc.ReleaseStock(context.Background(), command.StockReleaseCommand{
		EditionUid: "ed", WarehouseUid: "wh", Quantity: -1,
	}); err == nil {
		t.Errorf("expected error for non-positive release quantity")
	}
}

func TestAdjustStock_RejectsZero(t *testing.T) {
	svc := newStockSvc(&stubStockRepo{}, nil)
	if _, err := svc.AdjustStock(context.Background(), command.StockAdjustCommand{
		EditionUid: "ed", WarehouseUid: "wh", Quantity: 0, Reason: "noop",
	}); err == nil {
		t.Errorf("expected error for zero adjustment")
	}
}

func TestApplyStockChange_RepoError(t *testing.T) {
	wantErr := errors.New("stock not found for edition ed-1 in warehouse wh-missing")
	stk := &stubStockRepo{
		applyChangeFn: func(context.Context, repo.ApplyChangeParams) error { return wantErr },
	}
	svc := newStockSvc(stk, nil)
	if _, err := svc.ReceiveStock(context.Background(), command.StockReceiveCommand{
		EditionUid: "ed-1", WarehouseUid: "wh-missing", Quantity: 1,
	}); !errors.Is(err, wantErr) {
		t.Errorf("err = %v, want %v", err, wantErr)
	}
}

func TestAdjustStock_AppliesNegativeDelta(t *testing.T) {
	var applied *repo.ApplyChangeParams
	stk := &stubStockRepo{
		applyChangeFn: func(_ context.Context, p repo.ApplyChangeParams) error { applied = &p; return nil },
		findAllFn: func(context.Context, option.StockQueryOption) (int64, []*repo.StockResult, error) {
			return 1, []*repo.StockResult{{Uid: "stk", EditionUid: "ed", WarehouseUid: "wh", Quantity: 6}}, nil
		},
	}
	svc := newStockSvc(stk, nil)
	if _, err := svc.AdjustStock(context.Background(), command.StockAdjustCommand{
		EditionUid: "ed", WarehouseUid: "wh", Quantity: -4, Reason: "damaged",
	}); err != nil {
		t.Fatalf("err: %v", err)
	}
	if applied == nil || applied.Delta != -4 || applied.ChangeType != "ADJUSTMENT" {
		t.Errorf("applied = %+v", applied)
	}
	if applied.Reason == nil || *applied.Reason != "damaged" {
		t.Errorf("Reason = %v", applied.Reason)
	}
}

// ─── Saga reservation tests ─────────────────────────────────────────────────

func TestReserveStockForOrder_UsesDeterministicKey(t *testing.T) {
	var gotParams repo.ReserveParams
	stk := &stubStockRepo{
		findAllFn: func(context.Context, option.StockQueryOption) (int64, []*repo.StockResult, error) {
			return 1, []*repo.StockResult{{Uid: "stk", EditionUid: "ed-uid", WarehouseUid: "wh", Quantity: 7}}, nil
		},
	}
	resv := &stubReservationRepo{
		reserveFn: func(_ context.Context, p repo.ReserveParams) (*repo.ReservationResult, error) {
			gotParams = p
			return &repo.ReservationResult{EditionUid: "ed-uid", WarehouseUid: "wh", Quantity: 3, Status: "RESERVED"}, nil
		},
	}
	svc := &stockService{stockRepo: stk, reservationRepo: resv}
	if _, err := svc.ReserveStockForOrder(context.Background(), command.StockOrderReserveCommand{
		OrderUid: "order-9", EditionId: 42, Quantity: 3,
	}); err != nil {
		t.Fatalf("err: %v", err)
	}
	if gotParams.ReservationKey != "order-9:42" {
		t.Errorf("ReservationKey = %q, want order-9:42", gotParams.ReservationKey)
	}
	if gotParams.OrderRef != "order-9" || gotParams.EditionId != 42 || gotParams.Quantity != 3 {
		t.Errorf("params = %+v", gotParams)
	}
}

func TestReserveStockForOrder_RejectsNonPositive(t *testing.T) {
	svc := &stockService{reservationRepo: &stubReservationRepo{}}
	if _, err := svc.ReserveStockForOrder(context.Background(), command.StockOrderReserveCommand{
		OrderUid: "o", EditionId: 1, Quantity: 0,
	}); err == nil {
		t.Errorf("expected error for non-positive reserve quantity")
	}
}

func TestReleaseStockForOrder_SameKeyAsReserve(t *testing.T) {
	var gotKey string
	stk := &stubStockRepo{
		findAllFn: func(context.Context, option.StockQueryOption) (int64, []*repo.StockResult, error) {
			return 1, []*repo.StockResult{{Uid: "stk", EditionUid: "ed-uid", WarehouseUid: "wh", Quantity: 10}}, nil
		},
	}
	resv := &stubReservationRepo{
		releaseFn: func(_ context.Context, p repo.ReleaseParams) (*repo.ReservationResult, error) {
			gotKey = p.ReservationKey
			return &repo.ReservationResult{EditionUid: "ed-uid", WarehouseUid: "wh", Quantity: 3, Status: "RELEASED"}, nil
		},
	}
	svc := &stockService{stockRepo: stk, reservationRepo: resv}
	if _, err := svc.ReleaseStockForOrder(context.Background(), command.StockOrderReleaseCommand{
		OrderUid: "order-9", EditionId: 42, Quantity: 3,
	}); err != nil {
		t.Fatalf("err: %v", err)
	}
	if gotKey != "order-9:42" {
		t.Errorf("ReservationKey = %q, want order-9:42 (must match reserve)", gotKey)
	}
}

func TestReleaseStockForOrder_NoReservationIsNoOp(t *testing.T) {
	resv := &stubReservationRepo{
		releaseFn: func(context.Context, repo.ReleaseParams) (*repo.ReservationResult, error) {
			return nil, nil
		},
	}
	svc := &stockService{reservationRepo: resv}
	got, err := svc.ReleaseStockForOrder(context.Background(), command.StockOrderReleaseCommand{
		OrderUid: "order-x", EditionId: 1, Quantity: 1,
	})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if got != nil {
		t.Errorf("got = %+v, want nil (no-op release)", got)
	}
}

// ─── FindStockHistory ─────────────────────────────────────────────────────

func TestFindStockHistory(t *testing.T) {
	now := time.Now()
	stk := &stubStockRepo{
		findAllFn: func(_ context.Context, opt option.StockQueryOption) (int64, []*repo.StockResult, error) {
			if opt.EditionUid() == nil || *opt.EditionUid() != "ed-1" {
				t.Errorf("EditionUid = %v", opt.EditionUid())
			}
			if opt.WarehouseUid() == nil || *opt.WarehouseUid() != "wh-1" {
				t.Errorf("WarehouseUid = %v", opt.WarehouseUid())
			}
			return 1, []*repo.StockResult{{Id: 7, Uid: "stk-1", EditionUid: "ed-1", WarehouseUid: "wh-1", Quantity: 4}}, nil
		},
	}
	hist := &stubStockHistoryRepo{
		findAllFn: func(_ context.Context, stockUid string, _ pagination.PageOption) (int64, []*repo.StockHistoryResult, error) {
			if stockUid != "stk-1" {
				t.Errorf("stockUid = %q", stockUid)
			}
			reason := "restock"
			return 2, []*repo.StockHistoryResult{
				{Uid: "h1", ChangeType: "INBOUND", ChangeQty: 5, Reason: &reason, CreatedAt: now},
				{Uid: "h2", ChangeType: "RESERVATION", ChangeQty: -1, CreatedAt: now},
			}, nil
		},
	}
	svc := newStockSvc(stk, hist)
	wh := "wh-1"
	total, entries, err := svc.FindStockHistory(context.Background(), "ed-1", &wh, pagination.NewPageOptionDefault())
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if total != 2 || len(entries) != 2 {
		t.Fatalf("entries = %+v (total=%d)", entries, total)
	}
	if entries[0].ChangeType != "INBOUND" || entries[0].ChangeQty != 5 {
		t.Errorf("entries[0] = %+v", entries[0])
	}
}

func TestFindStockHistory_StockNotFound(t *testing.T) {
	stk := &stubStockRepo{
		findAllFn: func(context.Context, option.StockQueryOption) (int64, []*repo.StockResult, error) {
			return 0, nil, nil
		},
	}
	svc := newStockSvc(stk, nil)
	total, entries, err := svc.FindStockHistory(context.Background(), "ed-x", nil, pagination.NewPageOptionDefault())
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if total != 0 || entries != nil {
		t.Errorf("got total=%d entries=%+v", total, entries)
	}
}
