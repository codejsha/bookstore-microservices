package aggregate

import "time"

type TransferStatus string

const (
	TransferStatusPending   TransferStatus = "PENDING"
	TransferStatusInTransit TransferStatus = "IN_TRANSIT"
	TransferStatusCompleted TransferStatus = "COMPLETED"
	TransferStatusCancelled TransferStatus = "CANCELLED"
)

type StockTransferAggregate struct {
	Uid                string
	EditionUid         string
	SourceWarehouseUid string
	TargetWarehouseUid string
	Quantity           int32
	Status             TransferStatus
	Reason             *string
	CompletedAt        *time.Time
	CreatedAt          time.Time
	UpdatedAt          *time.Time
}

func (t *StockTransferAggregate) CanComplete() bool {
	return t.Status == TransferStatusPending || t.Status == TransferStatusInTransit
}

func (t *StockTransferAggregate) CanCancel() bool {
	return t.Status == TransferStatusPending
}
