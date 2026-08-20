package aggregate

import "time"

type AuditStatus string

const (
	AuditStatusInProgress AuditStatus = "IN_PROGRESS"
	AuditStatusCompleted  AuditStatus = "COMPLETED"
)

type StockAuditAggregate struct {
	Uid          string
	WarehouseUid string
	Status       AuditStatus
	Items        []*StockAuditItem
	Notes        *string
	CompletedAt  *time.Time
	CreatedAt    time.Time
	UpdatedAt    *time.Time
}

type StockAuditItem struct {
	EditionUid     string
	SystemQuantity int32
	ActualQuantity int32
	Difference     int32
}
