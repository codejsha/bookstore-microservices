package service

import (
	"context"
	"fmt"

	"github.com/codejsha/shared-library-go/pkg/pagination"

	"github.com/codejsha/bookstore-microservices/inventory/internal/application/port/repo"
	"github.com/codejsha/bookstore-microservices/inventory/internal/application/usecase"
	"github.com/codejsha/bookstore-microservices/inventory/internal/domain/aggregate"
	"github.com/codejsha/bookstore-microservices/inventory/internal/domain/model/command"
	"github.com/codejsha/bookstore-microservices/inventory/internal/domain/model/option"
)

var _ usecase.InventoryUseCase = (*stockService)(nil)

type stockService struct {
	stockRepo        repo.StockRepo
	warehouseRepo    repo.WarehouseRepo
	stockHistoryRepo repo.StockHistoryRepo
	reservationRepo  repo.StockReservationRepo
	transferRepo     repo.StockTransferRepo
	auditRepo        repo.StockAuditRepo
	closingRepo      repo.MonthlyClosingRepo
	balanceRepo      repo.StockBalanceRepo
}

func NewStockService(
	stockRepo repo.StockRepo,
	warehouseRepo repo.WarehouseRepo,
	stockHistoryRepo repo.StockHistoryRepo,
	reservationRepo repo.StockReservationRepo,
	transferRepo repo.StockTransferRepo,
	auditRepo repo.StockAuditRepo,
	closingRepo repo.MonthlyClosingRepo,
	balanceRepo repo.StockBalanceRepo,
) usecase.InventoryUseCase {
	return &stockService{
		stockRepo:        stockRepo,
		warehouseRepo:    warehouseRepo,
		stockHistoryRepo: stockHistoryRepo,
		reservationRepo:  reservationRepo,
		transferRepo:     transferRepo,
		auditRepo:        auditRepo,
		closingRepo:      closingRepo,
		balanceRepo:      balanceRepo,
	}
}

// ─── Stock query ────────────────────────────────────────────────────────────

func (s *stockService) FindAllStocks(ctx context.Context, opt option.StockQueryOption) (int64, []*aggregate.StockAggregate, error) {
	total, stocks, err := s.stockRepo.FindAll(ctx, opt)
	if err != nil {
		return 0, nil, err
	}

	byEdition := make(map[string][]*repo.StockResult)
	order := make([]string, 0)
	for _, st := range stocks {
		if _, seen := byEdition[st.EditionUid]; !seen {
			order = append(order, st.EditionUid)
		}
		byEdition[st.EditionUid] = append(byEdition[st.EditionUid], st)
	}

	aggs := make([]*aggregate.StockAggregate, 0, len(order))
	for _, key := range order {
		aggs = append(aggs, toStockAggregate(byEdition[key]))
	}
	return total, aggs, nil
}

func (s *stockService) FindStock(ctx context.Context, editionUid string) (*aggregate.StockAggregate, error) {
	uid := editionUid
	opt := option.NewStockQueryOption(option.StockQueryOption{}.WithEditionUid(&uid))
	_, stocks, err := s.stockRepo.FindAll(ctx, opt)
	if err != nil {
		return nil, err
	}
	if len(stocks) == 0 {
		return nil, nil
	}
	return toStockAggregate(stocks), nil
}

func (s *stockService) FindStockHistory(ctx context.Context, editionUid string, warehouseUid *string, page pagination.PageOption) (int64, []*aggregate.StockHistoryEntry, error) {
	stock, err := s.findStockRow(ctx, editionUid, warehouseUid)
	if err != nil {
		return 0, nil, err
	}
	if stock == nil {
		return 0, nil, nil
	}
	total, entries, err := s.stockHistoryRepo.FindAll(ctx, stock.Uid, page)
	if err != nil {
		return 0, nil, err
	}
	out := make([]*aggregate.StockHistoryEntry, len(entries))
	for i, e := range entries {
		out[i] = toStockHistoryEntry(e)
	}
	return total, out, nil
}

// ─── Stock operations ───────────────────────────────────────────────────────

func (s *stockService) ReceiveStock(ctx context.Context, cmd command.StockReceiveCommand) (*aggregate.StockAggregate, error) {
	if err := cmd.Validate(); err != nil {
		return nil, err
	}
	if cmd.Quantity <= 0 {
		return nil, fmt.Errorf("receive quantity must be positive: got %d: %w", cmd.Quantity, repo.ErrInvalidQuantity)
	}
	return s.applyStockChange(ctx, cmd.EditionUid, cmd.WarehouseUid, cmd.Quantity, aggregate.StockChangeInbound, cmd.Reason)
}

func (s *stockService) ReleaseStock(ctx context.Context, cmd command.StockReleaseCommand) (*aggregate.StockAggregate, error) {
	if err := cmd.Validate(); err != nil {
		return nil, err
	}
	if cmd.Quantity <= 0 {
		return nil, fmt.Errorf("release quantity must be positive: got %d: %w", cmd.Quantity, repo.ErrInvalidQuantity)
	}
	return s.applyStockChange(ctx, cmd.EditionUid, cmd.WarehouseUid, cmd.Quantity, aggregate.StockChangeRelease, cmd.Reason)
}

func (s *stockService) AdjustStock(ctx context.Context, cmd command.StockAdjustCommand) (*aggregate.StockAggregate, error) {
	if err := cmd.Validate(); err != nil {
		return nil, err
	}
	if cmd.Quantity == 0 {
		return nil, fmt.Errorf("adjust quantity must be non-zero: %w", repo.ErrInvalidQuantity)
	}
	reason := cmd.Reason
	return s.applyStockChange(ctx, cmd.EditionUid, cmd.WarehouseUid, cmd.Quantity, aggregate.StockChangeAdjustment, &reason)
}

func (s *stockService) ReserveStock(ctx context.Context, cmd command.StockReserveCommand) (*aggregate.StockAggregate, error) {
	if err := cmd.Validate(); err != nil {
		return nil, err
	}
	if cmd.Quantity <= 0 {
		return nil, fmt.Errorf("reserve quantity must be positive: got %d: %w", cmd.Quantity, repo.ErrInvalidQuantity)
	}
	return s.applyStockChange(ctx, cmd.EditionUid, cmd.WarehouseUid, -cmd.Quantity, aggregate.StockChangeReservation, cmd.Reason)
}

// ─── Saga (Temporal) stock operations ─────────────────────────────────────────

func (s *stockService) ReserveStockForOrder(ctx context.Context, cmd command.StockOrderReserveCommand) (*aggregate.StockAggregate, error) {
	if err := cmd.Validate(); err != nil {
		return nil, err
	}
	if cmd.Quantity <= 0 {
		return nil, fmt.Errorf("reserve quantity must be positive: got %d: %w", cmd.Quantity, repo.ErrInvalidQuantity)
	}
	res, err := s.reservationRepo.Reserve(ctx, repo.ReserveParams{
		ReservationKey: reservationKey(cmd.OrderUid, cmd.EditionId),
		OrderRef:       cmd.OrderUid,
		EditionId:      cmd.EditionId,
		Quantity:       cmd.Quantity,
		Reason:         cmd.Reason,
	})
	if err != nil {
		return nil, err
	}
	return s.FindStock(ctx, res.EditionUid)
}

func (s *stockService) ReleaseStockForOrder(ctx context.Context, cmd command.StockOrderReleaseCommand) (*aggregate.StockAggregate, error) {
	if err := cmd.Validate(); err != nil {
		return nil, err
	}
	res, err := s.reservationRepo.Release(ctx, repo.ReleaseParams{
		ReservationKey: reservationKey(cmd.OrderUid, cmd.EditionId),
		Reason:         cmd.Reason,
	})
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}
	return s.FindStock(ctx, res.EditionUid)
}

func reservationKey(orderRef string, editionId int64) string {
	return fmt.Sprintf("%s:%d", orderRef, editionId)
}

// ─── Warehouse management ───────────────────────────────────────────────────

func (s *stockService) FindAllWarehouses(ctx context.Context, opt option.WarehouseQueryOption) (int64, []*aggregate.WarehouseAggregate, error) {
	total, warehouses, err := s.warehouseRepo.FindAll(ctx, opt)
	if err != nil {
		return 0, nil, err
	}
	aggs := make([]*aggregate.WarehouseAggregate, len(warehouses))
	for i, w := range warehouses {
		aggs[i] = toWarehouseAggregate(w)
	}
	return total, aggs, nil
}

func (s *stockService) FindWarehouse(ctx context.Context, uid string) (*aggregate.WarehouseAggregate, error) {
	res, err := s.warehouseRepo.FindByUid(ctx, uid)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}
	return toWarehouseAggregate(res), nil
}

func (s *stockService) RegisterWarehouse(ctx context.Context, cmd command.WarehouseCreateCommand) (*aggregate.WarehouseAggregate, error) {
	if err := cmd.Validate(); err != nil {
		return nil, err
	}
	res, err := s.warehouseRepo.Create(ctx, repo.WarehouseCreateParams{
		Name:     cmd.Name,
		Address:  cmd.Address,
		Capacity: cmd.Capacity,
	})
	if err != nil {
		return nil, err
	}
	return toWarehouseAggregate(res), nil
}

func (s *stockService) UpdateWarehouse(ctx context.Context, uid string, cmd command.WarehouseUpdateCommand) (*aggregate.WarehouseAggregate, error) {
	if err := cmd.Validate(); err != nil {
		return nil, err
	}
	res, err := s.warehouseRepo.Update(ctx, repo.WarehouseUpdateParams{
		Uid:      uid,
		Name:     cmd.Name,
		Address:  cmd.Address,
		Capacity: cmd.Capacity,
	})
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}
	return toWarehouseAggregate(res), nil
}

// ─── Stock transfer ─────────────────────────────────────────────────────────

func (s *stockService) TransferStock(ctx context.Context, cmd command.StockTransferCommand) (*aggregate.StockTransferAggregate, error) {
	if err := cmd.Validate(); err != nil {
		return nil, err
	}
	if cmd.Quantity <= 0 {
		return nil, fmt.Errorf("transfer quantity must be positive: got %d: %w", cmd.Quantity, repo.ErrInvalidQuantity)
	}
	if cmd.SourceWarehouseUid == cmd.TargetWarehouseUid {
		return nil, fmt.Errorf("transfer source and target warehouses must differ: %w", repo.ErrSameWarehouse)
	}
	res, err := s.transferRepo.Create(ctx, repo.TransferCreateParams{
		EditionUid:         cmd.EditionUid,
		SourceWarehouseUid: cmd.SourceWarehouseUid,
		TargetWarehouseUid: cmd.TargetWarehouseUid,
		Quantity:           cmd.Quantity,
		Reason:             cmd.Reason,
	})
	if err != nil {
		return nil, err
	}
	return toTransferAggregate(res), nil
}

func (s *stockService) FindAllTransfers(ctx context.Context, opt option.TransferQueryOption) (int64, []*aggregate.StockTransferAggregate, error) {
	total, transfers, err := s.transferRepo.FindAll(ctx, opt)
	if err != nil {
		return 0, nil, err
	}
	aggs := make([]*aggregate.StockTransferAggregate, len(transfers))
	for i, t := range transfers {
		aggs[i] = toTransferAggregate(t)
	}
	return total, aggs, nil
}

func (s *stockService) FindTransfer(ctx context.Context, uid string) (*aggregate.StockTransferAggregate, error) {
	res, err := s.transferRepo.FindByUid(ctx, uid)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}
	return toTransferAggregate(res), nil
}

func (s *stockService) CompleteTransfer(ctx context.Context, uid string) (*aggregate.StockTransferAggregate, error) {
	res, err := s.transferRepo.FindByUid(ctx, uid)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}
	if agg := toTransferAggregate(res); !agg.CanComplete() {
		return nil, fmt.Errorf("transfer %s cannot be completed in status %s", uid, agg.Status)
	}
	done, err := s.transferRepo.Complete(ctx, uid)
	if err != nil {
		return nil, err
	}
	if done == nil {
		return nil, nil
	}
	return toTransferAggregate(done), nil
}

func (s *stockService) CancelTransfer(ctx context.Context, uid string) (*aggregate.StockTransferAggregate, error) {
	res, err := s.transferRepo.FindByUid(ctx, uid)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}
	if agg := toTransferAggregate(res); !agg.CanCancel() {
		return nil, fmt.Errorf("transfer %s cannot be cancelled in status %s", uid, agg.Status)
	}
	done, err := s.transferRepo.Cancel(ctx, uid)
	if err != nil {
		return nil, err
	}
	if done == nil {
		return nil, nil
	}
	return toTransferAggregate(done), nil
}

// ─── Stock audit ────────────────────────────────────────────────────────────

func (s *stockService) CreateAudit(ctx context.Context, cmd command.StockAuditCreateCommand) (*aggregate.StockAuditAggregate, error) {
	if err := cmd.Validate(); err != nil {
		return nil, err
	}
	items := make([]repo.AuditItemParam, len(cmd.Items))
	for i, it := range cmd.Items {
		items[i] = repo.AuditItemParam{
			EditionUid:     it.EditionUid,
			ActualQuantity: it.ActualQuantity,
		}
	}
	res, err := s.auditRepo.Create(ctx, repo.AuditCreateParams{
		WarehouseUid: cmd.WarehouseUid,
		Items:        items,
		Notes:        cmd.Notes,
	})
	if err != nil {
		return nil, err
	}
	return toAuditAggregate(res), nil
}

func (s *stockService) FindAllAudits(ctx context.Context, opt option.AuditQueryOption) (int64, []*aggregate.StockAuditAggregate, error) {
	total, audits, err := s.auditRepo.FindAll(ctx, opt)
	if err != nil {
		return 0, nil, err
	}
	aggs := make([]*aggregate.StockAuditAggregate, len(audits))
	for i, a := range audits {
		aggs[i] = toAuditAggregate(a)
	}
	return total, aggs, nil
}

func (s *stockService) FindAudit(ctx context.Context, uid string) (*aggregate.StockAuditAggregate, error) {
	res, err := s.auditRepo.FindByUid(ctx, uid)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}
	return toAuditAggregate(res), nil
}

func (s *stockService) CompleteAudit(ctx context.Context, uid string) (*aggregate.StockAuditAggregate, error) {
	res, err := s.auditRepo.FindByUid(ctx, uid)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}
	if res.Status == string(aggregate.AuditStatusCompleted) {
		return nil, fmt.Errorf("audit %s is already completed", uid)
	}
	done, err := s.auditRepo.Complete(ctx, uid)
	if err != nil {
		return nil, err
	}
	if done == nil {
		return nil, nil
	}
	return toAuditAggregate(done), nil
}

// ─── Monthly closing ────────────────────────────────────────────────────────

func (s *stockService) CreateMonthlyClosing(ctx context.Context, cmd command.MonthlyClosingCommand) (*aggregate.MonthlyClosingAggregate, error) {
	if err := cmd.Validate(); err != nil {
		return nil, err
	}
	res, err := s.closingRepo.Create(ctx, repo.ClosingCreateParams{
		WarehouseUid: cmd.WarehouseUid,
		Year:         cmd.Year,
		Month:        cmd.Month,
	})
	if err != nil {
		return nil, err
	}
	return toClosingAggregate(res), nil
}

func (s *stockService) FindAllClosings(ctx context.Context, opt option.ClosingQueryOption) (int64, []*aggregate.MonthlyClosingAggregate, error) {
	total, closings, err := s.closingRepo.FindAll(ctx, opt)
	if err != nil {
		return 0, nil, err
	}
	aggs := make([]*aggregate.MonthlyClosingAggregate, len(closings))
	for i, c := range closings {
		aggs[i] = toClosingAggregate(c)
	}
	return total, aggs, nil
}

func (s *stockService) FindClosing(ctx context.Context, uid string) (*aggregate.MonthlyClosingAggregate, error) {
	res, err := s.closingRepo.FindByUid(ctx, uid)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}
	return toClosingAggregate(res), nil
}

// ─── Stock balance ──────────────────────────────────────────────────────────

func (s *stockService) FindStockBalance(ctx context.Context, opt option.BalanceQueryOption) (int64, []*aggregate.StockBalanceEntry, error) {
	total, rows, err := s.balanceRepo.FindAll(ctx, opt)
	if err != nil {
		return 0, nil, err
	}
	entries := make([]*aggregate.StockBalanceEntry, len(rows))
	for i, r := range rows {
		entries[i] = toBalanceEntry(r)
	}
	return total, entries, nil
}

// ─── helpers ────────────────────────────────────────────────────────────────

func (s *stockService) findStockRow(ctx context.Context, editionUid string, warehouseUid *string) (*repo.StockResult, error) {
	opts := []option.StockQueryOptionFunc{option.StockQueryOption{}.WithEditionUid(&editionUid)}
	if warehouseUid != nil {
		opts = append(opts, option.StockQueryOption{}.WithWarehouseUid(warehouseUid))
	}
	_, stocks, err := s.stockRepo.FindAll(ctx, option.NewStockQueryOption(opts...))
	if err != nil {
		return nil, err
	}
	if len(stocks) == 0 {
		return nil, nil
	}
	return stocks[0], nil
}

func (s *stockService) applyStockChange(
	ctx context.Context,
	editionUid string,
	warehouseUid string,
	delta int32,
	changeType aggregate.StockChangeType,
	reason *string,
) (*aggregate.StockAggregate, error) {
	if err := s.stockRepo.ApplyChange(ctx, repo.ApplyChangeParams{
		EditionUid:   editionUid,
		WarehouseUid: warehouseUid,
		Delta:        delta,
		ChangeType:   string(changeType),
		Reason:       reason,
	}); err != nil {
		return nil, err
	}
	return s.FindStock(ctx, editionUid)
}

// ─── mappers ────────────────────────────────────────────────────────────────

func toStockAggregate(rows []*repo.StockResult) *aggregate.StockAggregate {
	if len(rows) == 0 {
		return nil
	}
	warehouses := make([]*aggregate.StockWarehouse, 0, len(rows))
	var total int32
	for _, r := range rows {
		warehouses = append(warehouses, &aggregate.StockWarehouse{
			WarehouseUid:  r.WarehouseUid,
			WarehouseName: r.WarehouseName,
			Quantity:      r.Quantity,
		})
		total += r.Quantity
	}
	return &aggregate.StockAggregate{
		Uid:           rows[0].Uid,
		EditionUid:    rows[0].EditionUid,
		TotalQuantity: total,
		Warehouses:    warehouses,
	}
}

func toWarehouseAggregate(r *repo.WarehouseResult) *aggregate.WarehouseAggregate {
	return &aggregate.WarehouseAggregate{
		Uid:       r.Uid,
		Name:      r.Name,
		Address:   r.Address,
		Capacity:  r.Capacity,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}

func toStockHistoryEntry(r *repo.StockHistoryResult) *aggregate.StockHistoryEntry {
	return &aggregate.StockHistoryEntry{
		Uid:        r.Uid,
		ChangeType: aggregate.StockChangeType(r.ChangeType),
		Reason:     r.Reason,
		ChangeQty:  r.ChangeQty,
		CreatedAt:  r.CreatedAt,
	}
}

func toTransferAggregate(r *repo.TransferResult) *aggregate.StockTransferAggregate {
	return &aggregate.StockTransferAggregate{
		Uid:                r.Uid,
		EditionUid:         r.EditionUid,
		SourceWarehouseUid: r.SourceWarehouseUid,
		TargetWarehouseUid: r.TargetWarehouseUid,
		Quantity:           r.Quantity,
		Status:             aggregate.TransferStatus(r.Status),
		Reason:             r.Reason,
		CompletedAt:        r.CompletedAt,
		CreatedAt:          r.CreatedAt,
		UpdatedAt:          r.UpdatedAt,
	}
}

func toAuditAggregate(r *repo.AuditResult) *aggregate.StockAuditAggregate {
	items := make([]*aggregate.StockAuditItem, len(r.Items))
	for i, it := range r.Items {
		items[i] = &aggregate.StockAuditItem{
			EditionUid:     it.EditionUid,
			SystemQuantity: it.SystemQuantity,
			ActualQuantity: it.ActualQuantity,
			Difference:     it.Difference,
		}
	}
	return &aggregate.StockAuditAggregate{
		Uid:          r.Uid,
		WarehouseUid: r.WarehouseUid,
		Status:       aggregate.AuditStatus(r.Status),
		Items:        items,
		Notes:        r.Notes,
		CompletedAt:  r.CompletedAt,
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
	}
}

func toClosingAggregate(r *repo.ClosingResult) *aggregate.MonthlyClosingAggregate {
	items := make([]*aggregate.MonthlyClosingItem, len(r.Items))
	for i, it := range r.Items {
		items[i] = &aggregate.MonthlyClosingItem{
			EditionUid:       it.EditionUid,
			OpeningQuantity:  it.OpeningQuantity,
			InboundQuantity:  it.InboundQuantity,
			OutboundQuantity: it.OutboundQuantity,
			AdjustQuantity:   it.AdjustQuantity,
			ClosingQuantity:  it.ClosingQuantity,
		}
	}
	return &aggregate.MonthlyClosingAggregate{
		Uid:          r.Uid,
		WarehouseUid: r.WarehouseUid,
		Year:         r.Year,
		Month:        r.Month,
		Status:       aggregate.ClosingStatus(r.Status),
		Items:        items,
		ClosedAt:     r.ClosedAt,
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
	}
}

func toBalanceEntry(r *repo.BalanceResult) *aggregate.StockBalanceEntry {
	return &aggregate.StockBalanceEntry{
		EditionUid:       r.EditionUid,
		WarehouseUid:     r.WarehouseUid,
		WarehouseName:    r.WarehouseName,
		InboundQuantity:  r.InboundQuantity,
		OutboundQuantity: r.OutboundQuantity,
		AdjustQuantity:   r.AdjustQuantity,
		CurrentQuantity:  r.CurrentQuantity,
	}
}
