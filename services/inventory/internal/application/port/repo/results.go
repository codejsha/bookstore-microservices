package repo

import "time"

type StockResult struct {
	Id            int64
	Uid           string
	EditionUid    string
	WarehouseUid  string
	WarehouseName string
	Quantity      int32
}

type ApplyChangeParams struct {
	EditionUid   string
	WarehouseUid string
	Delta        int32
	ChangeType   string
	Reason       *string
}

// ─── Warehouse ───────────────────────────────────────────────────────────────

type WarehouseCreateParams struct {
	Name     string
	Address  *string
	Capacity int32
}

type WarehouseUpdateParams struct {
	Uid      string
	Name     *string
	Address  *string
	Capacity *int32
}

type WarehouseResult struct {
	Uid       string
	Name      string
	Address   *string
	Capacity  int32
	CreatedAt time.Time
	UpdatedAt *time.Time
}

// ─── Stock history ───────────────────────────────────────────────────────────

type StockHistoryResult struct {
	Uid        string
	ChangeType string
	Reason     *string
	ChangeQty  int32
	CreatedAt  time.Time
}

type ReserveParams struct {
	ReservationKey string
	OrderRef       string
	EditionId      int64
	Quantity       int32
	Reason         *string
}

type ReleaseParams struct {
	ReservationKey string
	Reason         *string
}

type ReservationResult struct {
	Uid            string
	ReservationKey string
	EditionUid     string
	WarehouseUid   string
	Quantity       int32
	Status         string
	AlreadyApplied bool
}

// ─── Transfer ────────────────────────────────────────────────────────────────

type TransferCreateParams struct {
	EditionUid         string
	SourceWarehouseUid string
	TargetWarehouseUid string
	Quantity           int32
	Reason             *string
}

type TransferResult struct {
	Uid                string
	EditionUid         string
	SourceWarehouseUid string
	TargetWarehouseUid string
	Quantity           int32
	Status             string
	Reason             *string
	CompletedAt        *time.Time
	CreatedAt          time.Time
	UpdatedAt          *time.Time
}

// ─── Audit ───────────────────────────────────────────────────────────────────

type AuditCreateParams struct {
	WarehouseUid string
	Items        []AuditItemParam
	Notes        *string
}

type AuditItemParam struct {
	EditionUid     string
	ActualQuantity int32
}

type AuditResult struct {
	Uid          string
	WarehouseUid string
	Status       string
	Items        []AuditItemResult
	Notes        *string
	CompletedAt  *time.Time
	CreatedAt    time.Time
	UpdatedAt    *time.Time
}

type AuditItemResult struct {
	EditionUid     string
	SystemQuantity int32
	ActualQuantity int32
	Difference     int32
}

// ─── Monthly closing ─────────────────────────────────────────────────────────

type ClosingCreateParams struct {
	WarehouseUid string
	Year         int32
	Month        int32
}

type ClosingResult struct {
	Uid          string
	WarehouseUid string
	Year         int32
	Month        int32
	Status       string
	Items        []ClosingItemResult
	ClosedAt     *time.Time
	CreatedAt    time.Time
	UpdatedAt    *time.Time
}

type ClosingItemResult struct {
	EditionUid       string
	OpeningQuantity  int32
	InboundQuantity  int32
	OutboundQuantity int32
	AdjustQuantity   int32
	ClosingQuantity  int32
}

// ─── Balance ─────────────────────────────────────────────────────────────────

type BalanceResult struct {
	EditionUid       string
	WarehouseUid     string
	WarehouseName    string
	InboundQuantity  int32
	OutboundQuantity int32
	AdjustQuantity   int32
	CurrentQuantity  int32
}
