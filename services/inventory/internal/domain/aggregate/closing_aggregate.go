package aggregate

import "time"

type ClosingStatus string

const (
	ClosingStatusOpen   ClosingStatus = "OPEN"
	ClosingStatusClosed ClosingStatus = "CLOSED"
)

type MonthlyClosingAggregate struct {
	Uid          string
	WarehouseUid string
	Year         int32
	Month        int32
	Status       ClosingStatus
	Items        []*MonthlyClosingItem
	ClosedAt     *time.Time
	CreatedAt    time.Time
	UpdatedAt    *time.Time
}

type MonthlyClosingItem struct {
	EditionUid       string
	OpeningQuantity  int32
	InboundQuantity  int32
	OutboundQuantity int32
	AdjustQuantity   int32
	ClosingQuantity  int32
}
