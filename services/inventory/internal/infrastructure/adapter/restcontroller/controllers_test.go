package restcontroller

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/codejsha/shared-library-go/pkg/pagination"

	"github.com/codejsha/bookstore-microservices/inventory/generated/application/port/openapi"
	"github.com/codejsha/bookstore-microservices/inventory/internal/application/usecase"
	"github.com/codejsha/bookstore-microservices/inventory/internal/domain/aggregate"
	"github.com/codejsha/bookstore-microservices/inventory/internal/domain/model/command"
	"github.com/codejsha/bookstore-microservices/inventory/internal/domain/model/option"
)

type stubUseCase struct {
	findAllStocks     func(context.Context, option.StockQueryOption) (int64, []*aggregate.StockAggregate, error)
	findStock         func(context.Context, string) (*aggregate.StockAggregate, error)
	findStockHistory  func(context.Context, string, *string, pagination.PageOption) (int64, []*aggregate.StockHistoryEntry, error)
	receiveStock      func(context.Context, command.StockReceiveCommand) (*aggregate.StockAggregate, error)
	releaseStock      func(context.Context, command.StockReleaseCommand) (*aggregate.StockAggregate, error)
	adjustStock       func(context.Context, command.StockAdjustCommand) (*aggregate.StockAggregate, error)
	reserveStock      func(context.Context, command.StockReserveCommand) (*aggregate.StockAggregate, error)
	findAllWarehouses func(context.Context, option.WarehouseQueryOption) (int64, []*aggregate.WarehouseAggregate, error)
	findWarehouse     func(context.Context, string) (*aggregate.WarehouseAggregate, error)
	registerWarehouse func(context.Context, command.WarehouseCreateCommand) (*aggregate.WarehouseAggregate, error)
	updateWarehouse   func(context.Context, string, command.WarehouseUpdateCommand) (*aggregate.WarehouseAggregate, error)

	transferStock    func(context.Context, command.StockTransferCommand) (*aggregate.StockTransferAggregate, error)
	findAllTransfers func(context.Context, option.TransferQueryOption) (int64, []*aggregate.StockTransferAggregate, error)
	findTransfer     func(context.Context, string) (*aggregate.StockTransferAggregate, error)
	completeTransfer func(context.Context, string) (*aggregate.StockTransferAggregate, error)
	cancelTransfer   func(context.Context, string) (*aggregate.StockTransferAggregate, error)
	createAudit      func(context.Context, command.StockAuditCreateCommand) (*aggregate.StockAuditAggregate, error)
	findAllAudits    func(context.Context, option.AuditQueryOption) (int64, []*aggregate.StockAuditAggregate, error)
	findAudit        func(context.Context, string) (*aggregate.StockAuditAggregate, error)
	completeAudit    func(context.Context, string) (*aggregate.StockAuditAggregate, error)
	createClosing    func(context.Context, command.MonthlyClosingCommand) (*aggregate.MonthlyClosingAggregate, error)
	findAllClosings  func(context.Context, option.ClosingQueryOption) (int64, []*aggregate.MonthlyClosingAggregate, error)
	findClosing      func(context.Context, string) (*aggregate.MonthlyClosingAggregate, error)
	findStockBalance func(context.Context, option.BalanceQueryOption) (int64, []*aggregate.StockBalanceEntry, error)
}

var _ usecase.InventoryUseCase = (*stubUseCase)(nil)

func (s *stubUseCase) FindAllStocks(ctx context.Context, opt option.StockQueryOption) (int64, []*aggregate.StockAggregate, error) {
	return s.findAllStocks(ctx, opt)
}
func (s *stubUseCase) FindStock(ctx context.Context, editionUid string) (*aggregate.StockAggregate, error) {
	return s.findStock(ctx, editionUid)
}
func (s *stubUseCase) FindStockHistory(ctx context.Context, editionUid string, warehouseUid *string, page pagination.PageOption) (int64, []*aggregate.StockHistoryEntry, error) {
	return s.findStockHistory(ctx, editionUid, warehouseUid, page)
}
func (s *stubUseCase) ReceiveStock(ctx context.Context, cmd command.StockReceiveCommand) (*aggregate.StockAggregate, error) {
	return s.receiveStock(ctx, cmd)
}
func (s *stubUseCase) ReleaseStock(ctx context.Context, cmd command.StockReleaseCommand) (*aggregate.StockAggregate, error) {
	return s.releaseStock(ctx, cmd)
}
func (s *stubUseCase) AdjustStock(ctx context.Context, cmd command.StockAdjustCommand) (*aggregate.StockAggregate, error) {
	return s.adjustStock(ctx, cmd)
}
func (s *stubUseCase) ReserveStock(ctx context.Context, cmd command.StockReserveCommand) (*aggregate.StockAggregate, error) {
	return s.reserveStock(ctx, cmd)
}
func (s *stubUseCase) ReserveStockForOrder(context.Context, command.StockOrderReserveCommand) (*aggregate.StockAggregate, error) {
	return nil, nil
}
func (s *stubUseCase) ReleaseStockForOrder(context.Context, command.StockOrderReleaseCommand) (*aggregate.StockAggregate, error) {
	return nil, nil
}
func (s *stubUseCase) FindAllWarehouses(ctx context.Context, opt option.WarehouseQueryOption) (int64, []*aggregate.WarehouseAggregate, error) {
	return s.findAllWarehouses(ctx, opt)
}
func (s *stubUseCase) FindWarehouse(ctx context.Context, uid string) (*aggregate.WarehouseAggregate, error) {
	return s.findWarehouse(ctx, uid)
}
func (s *stubUseCase) RegisterWarehouse(ctx context.Context, cmd command.WarehouseCreateCommand) (*aggregate.WarehouseAggregate, error) {
	return s.registerWarehouse(ctx, cmd)
}
func (s *stubUseCase) UpdateWarehouse(ctx context.Context, uid string, cmd command.WarehouseUpdateCommand) (*aggregate.WarehouseAggregate, error) {
	return s.updateWarehouse(ctx, uid, cmd)
}

func (s *stubUseCase) TransferStock(ctx context.Context, cmd command.StockTransferCommand) (*aggregate.StockTransferAggregate, error) {
	return s.transferStock(ctx, cmd)
}
func (s *stubUseCase) FindAllTransfers(ctx context.Context, opt option.TransferQueryOption) (int64, []*aggregate.StockTransferAggregate, error) {
	return s.findAllTransfers(ctx, opt)
}
func (s *stubUseCase) FindTransfer(ctx context.Context, uid string) (*aggregate.StockTransferAggregate, error) {
	return s.findTransfer(ctx, uid)
}
func (s *stubUseCase) CompleteTransfer(ctx context.Context, uid string) (*aggregate.StockTransferAggregate, error) {
	return s.completeTransfer(ctx, uid)
}
func (s *stubUseCase) CancelTransfer(ctx context.Context, uid string) (*aggregate.StockTransferAggregate, error) {
	return s.cancelTransfer(ctx, uid)
}
func (s *stubUseCase) CreateAudit(ctx context.Context, cmd command.StockAuditCreateCommand) (*aggregate.StockAuditAggregate, error) {
	return s.createAudit(ctx, cmd)
}
func (s *stubUseCase) FindAllAudits(ctx context.Context, opt option.AuditQueryOption) (int64, []*aggregate.StockAuditAggregate, error) {
	return s.findAllAudits(ctx, opt)
}
func (s *stubUseCase) FindAudit(ctx context.Context, uid string) (*aggregate.StockAuditAggregate, error) {
	return s.findAudit(ctx, uid)
}
func (s *stubUseCase) CompleteAudit(ctx context.Context, uid string) (*aggregate.StockAuditAggregate, error) {
	return s.completeAudit(ctx, uid)
}
func (s *stubUseCase) CreateMonthlyClosing(ctx context.Context, cmd command.MonthlyClosingCommand) (*aggregate.MonthlyClosingAggregate, error) {
	return s.createClosing(ctx, cmd)
}
func (s *stubUseCase) FindAllClosings(ctx context.Context, opt option.ClosingQueryOption) (int64, []*aggregate.MonthlyClosingAggregate, error) {
	return s.findAllClosings(ctx, opt)
}
func (s *stubUseCase) FindClosing(ctx context.Context, uid string) (*aggregate.MonthlyClosingAggregate, error) {
	return s.findClosing(ctx, uid)
}
func (s *stubUseCase) FindStockBalance(ctx context.Context, opt option.BalanceQueryOption) (int64, []*aggregate.StockBalanceEntry, error) {
	return s.findStockBalance(ctx, opt)
}

func ptrStr(s string) *string { return &s }

// ─── Stock controller ─────────────────────────────────────────────────────

func TestStockController_StocksGetAll(t *testing.T) {
	use := &stubUseCase{
		findAllStocks: func(_ context.Context, opt option.StockQueryOption) (int64, []*aggregate.StockAggregate, error) {
			if opt.EditionUid() == nil || *opt.EditionUid() != "ed-1" {
				t.Errorf("EditionUid = %v", opt.EditionUid())
			}
			return 1, []*aggregate.StockAggregate{
				{
					Uid: "stk-1", EditionUid: "ed-1", TotalQuantity: 9,
					Warehouses: []*aggregate.StockWarehouse{
						{WarehouseUid: "wh-1", WarehouseName: "Main", Quantity: 9},
					},
				},
			}, nil
		},
	}
	ctrl := NewStockController(use)
	editionUid := "ed-1"
	resp, err := ctrl.StocksGetAll(context.Background(), &editionUid, nil, nil, nil, nil)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp.Total != 1 || len(resp.Items) != 1 {
		t.Fatalf("resp = %+v", resp)
	}
	if resp.Items[0].TotalQuantity != 9 || resp.Items[0].Warehouses[0].WarehouseName != "Main" {
		t.Errorf("Items[0] = %+v", resp.Items[0])
	}
}

func TestStockController_StocksRead(t *testing.T) {
	use := &stubUseCase{
		findStock: func(_ context.Context, editionUid string) (*aggregate.StockAggregate, error) {
			if editionUid != "ed-1" {
				t.Errorf("editionUid = %q", editionUid)
			}
			return &aggregate.StockAggregate{Uid: "stk-1", EditionUid: editionUid, TotalQuantity: 5}, nil
		},
	}
	ctrl := NewStockController(use)
	resp, err := ctrl.StocksRead(context.Background(), "ed-1")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp.EditionUid != "ed-1" || resp.TotalQuantity != 5 {
		t.Errorf("resp = %+v", resp)
	}
}

func TestStockController_StocksReceive(t *testing.T) {
	captured := command.StockReceiveCommand{}
	use := &stubUseCase{
		receiveStock: func(_ context.Context, cmd command.StockReceiveCommand) (*aggregate.StockAggregate, error) {
			captured = cmd
			return &aggregate.StockAggregate{Uid: "stk-1", EditionUid: cmd.EditionUid, TotalQuantity: 10}, nil
		},
	}
	ctrl := NewStockController(use)
	reason := "restock"
	resp, err := ctrl.StocksReceive(context.Background(), openapi.StockReceiveRequest{
		EditionUid:   "ed-1",
		WarehouseUid: "wh-1",
		Quantity:     10,
		Reason:       &reason,
	})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if captured.EditionUid != "ed-1" || captured.WarehouseUid != "wh-1" || captured.Quantity != 10 {
		t.Errorf("captured = %+v", captured)
	}
	if captured.Reason == nil || *captured.Reason != reason {
		t.Errorf("Reason = %v", captured.Reason)
	}
	if resp.TotalQuantity != 10 {
		t.Errorf("TotalQuantity = %d", resp.TotalQuantity)
	}
}

func TestStockController_StocksReceive_PropagatesError(t *testing.T) {
	wantErr := errors.New("oversold")
	use := &stubUseCase{
		receiveStock: func(context.Context, command.StockReceiveCommand) (*aggregate.StockAggregate, error) {
			return nil, wantErr
		},
	}
	ctrl := NewStockController(use)
	_, err := ctrl.StocksReceive(context.Background(), openapi.StockReceiveRequest{EditionUid: "x", WarehouseUid: "w", Quantity: 1})
	if !errors.Is(err, wantErr) {
		t.Errorf("err = %v, want %v", err, wantErr)
	}
}

func TestStockController_StocksAdjust_RequiredReason(t *testing.T) {
	captured := command.StockAdjustCommand{}
	use := &stubUseCase{
		adjustStock: func(_ context.Context, cmd command.StockAdjustCommand) (*aggregate.StockAggregate, error) {
			captured = cmd
			return &aggregate.StockAggregate{Uid: "stk", EditionUid: cmd.EditionUid}, nil
		},
	}
	ctrl := NewStockController(use)
	_, err := ctrl.StocksAdjust(context.Background(), openapi.StockAdjustRequest{
		EditionUid: "ed", WarehouseUid: "wh", Quantity: -5, Reason: "damaged",
	})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if captured.Quantity != -5 || captured.Reason != "damaged" {
		t.Errorf("captured = %+v", captured)
	}
}

func TestStockController_StocksReserveAndRelease(t *testing.T) {
	use := &stubUseCase{
		reserveStock: func(_ context.Context, cmd command.StockReserveCommand) (*aggregate.StockAggregate, error) {
			if cmd.Quantity != 3 {
				t.Errorf("reserve qty = %d", cmd.Quantity)
			}
			return &aggregate.StockAggregate{Uid: "stk-1"}, nil
		},
		releaseStock: func(_ context.Context, cmd command.StockReleaseCommand) (*aggregate.StockAggregate, error) {
			if cmd.Quantity != 2 {
				t.Errorf("release qty = %d", cmd.Quantity)
			}
			return &aggregate.StockAggregate{Uid: "stk-1"}, nil
		},
	}
	ctrl := NewStockController(use)
	if _, err := ctrl.StocksReserve(context.Background(), openapi.StockReserveRequest{EditionUid: "ed", WarehouseUid: "wh", Quantity: 3}); err != nil {
		t.Errorf("reserve err: %v", err)
	}
	if _, err := ctrl.StocksRelease(context.Background(), openapi.StockReleaseRequest{EditionUid: "ed", WarehouseUid: "wh", Quantity: 2}); err != nil {
		t.Errorf("release err: %v", err)
	}
}

func TestStockController_StocksHistory(t *testing.T) {
	use := &stubUseCase{
		findStockHistory: func(_ context.Context, editionUid string, warehouseUid *string, _ pagination.PageOption) (int64, []*aggregate.StockHistoryEntry, error) {
			if editionUid != "ed-1" {
				t.Errorf("editionUid = %q", editionUid)
			}
			if warehouseUid == nil || *warehouseUid != "wh-1" {
				t.Errorf("warehouseUid = %v", warehouseUid)
			}
			reason := "restock"
			return 1, []*aggregate.StockHistoryEntry{{
				Uid: "sh-1", ChangeType: aggregate.StockChangeInbound, Reason: &reason, ChangeQty: 5, CreatedAt: time.Now(),
			}}, nil
		},
	}
	ctrl := NewStockController(use)
	wh := "wh-1"
	resp, err := ctrl.StocksHistory(context.Background(), "ed-1", &wh, nil, nil, nil)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp.Total != 1 || len(resp.Items) != 1 {
		t.Fatalf("resp = %+v", resp)
	}
	if resp.Items[0].ChangeType != openapi.StockChangeType("INBOUND") {
		t.Errorf("ChangeType = %v", resp.Items[0].ChangeType)
	}
}

// ─── Warehouse controller ─────────────────────────────────────────────────

func TestWarehouseController_GetAll(t *testing.T) {
	use := &stubUseCase{
		findAllWarehouses: func(_ context.Context, opt option.WarehouseQueryOption) (int64, []*aggregate.WarehouseAggregate, error) {
			if opt.Name() == nil || *opt.Name() != "Main" {
				t.Errorf("Name = %v", opt.Name())
			}
			return 1, []*aggregate.WarehouseAggregate{{Uid: "wh-1", Name: "Main", Capacity: 1000}}, nil
		},
	}
	ctrl := NewWarehouseController(use)
	name := "Main"
	resp, err := ctrl.WarehousesGetAll(context.Background(), &name, nil, nil, nil)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp.Total != 1 || resp.Items[0].Capacity != 1000 {
		t.Errorf("resp = %+v", resp)
	}
}

func TestWarehouseController_Create(t *testing.T) {
	captured := command.WarehouseCreateCommand{}
	use := &stubUseCase{
		registerWarehouse: func(_ context.Context, cmd command.WarehouseCreateCommand) (*aggregate.WarehouseAggregate, error) {
			captured = cmd
			return &aggregate.WarehouseAggregate{Uid: "wh-new", Name: cmd.Name}, nil
		},
	}
	ctrl := NewWarehouseController(use)
	addr := "Mars"
	if err := ctrl.WarehousesCreate(context.Background(), openapi.WarehouseCreateRequest{
		Name: "New", Address: &addr, Capacity: 500,
	}); err != nil {
		t.Fatalf("err: %v", err)
	}
	if captured.Name != "New" || captured.Capacity != 500 {
		t.Errorf("captured = %+v", captured)
	}
	if captured.Address == nil || *captured.Address != "Mars" {
		t.Errorf("Address = %v", captured.Address)
	}
}

func TestWarehouseController_Update(t *testing.T) {
	use := &stubUseCase{
		updateWarehouse: func(_ context.Context, uid string, cmd command.WarehouseUpdateCommand) (*aggregate.WarehouseAggregate, error) {
			if uid != "wh-1" {
				t.Errorf("uid = %q", uid)
			}
			if cmd.Name == nil || *cmd.Name != "Renamed" {
				t.Errorf("Name = %v", cmd.Name)
			}
			return &aggregate.WarehouseAggregate{Uid: uid, Name: *cmd.Name, Capacity: 200}, nil
		},
	}
	ctrl := NewWarehouseController(use)
	name := "Renamed"
	resp, err := ctrl.WarehousesUpdate(context.Background(), "wh-1", openapi.WarehouseUpdateRequest{Name: &name})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp.Name != "Renamed" {
		t.Errorf("Name = %q", resp.Name)
	}
}

// ─── mapping helpers ───────────────────────────────────────────────────────

func TestToStockFindResponse(t *testing.T) {
	got := toStockFindResponse(&aggregate.StockAggregate{
		Uid: "stk-1", EditionUid: "ed-1", TotalQuantity: 7,
		Warehouses: []*aggregate.StockWarehouse{
			{WarehouseUid: "wh-1", WarehouseName: "Main", Quantity: 4},
			{WarehouseUid: "wh-2", WarehouseName: "Backup", Quantity: 3},
		},
	})
	if len(got.Warehouses) != 2 || got.Warehouses[1].WarehouseName != "Backup" {
		t.Errorf("Warehouses = %+v", got.Warehouses)
	}
}

func TestToStockHistoryItem(t *testing.T) {
	now := time.Now()
	reason := "x"
	got := toStockHistoryItem(&aggregate.StockHistoryEntry{
		Uid: "sh", ChangeType: aggregate.StockChangeReservation, Reason: &reason, ChangeQty: -2, CreatedAt: now,
	})
	if got.ChangeType != openapi.StockChangeType("RESERVATION") || got.ChangeQty != -2 {
		t.Errorf("got = %+v", got)
	}
	if got.Reason == nil || *got.Reason != reason {
		t.Errorf("Reason = %v", got.Reason)
	}
	if !got.CreatedAt.Equal(now) {
		t.Errorf("CreatedAt = %v", got.CreatedAt)
	}
}

func TestToWarehouseFindResponse(t *testing.T) {
	addr := "Earth"
	now := time.Now()
	got := toWarehouseFindResponse(&aggregate.WarehouseAggregate{
		Uid: "wh-1", Name: "Main", Address: &addr, Capacity: 100, CreatedAt: now,
	})
	if got.Uid != "wh-1" || got.Capacity != 100 {
		t.Errorf("got = %+v", got)
	}
	if got.Address == nil || *got.Address != addr {
		t.Errorf("Address = %v", got.Address)
	}
}

var _ = ptrStr
