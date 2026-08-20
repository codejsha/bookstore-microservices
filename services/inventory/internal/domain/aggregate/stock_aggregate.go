package aggregate

import "time"

type StockChangeType string

const (
	StockChangeInbound     StockChangeType = "INBOUND"
	StockChangeOutbound    StockChangeType = "OUTBOUND"
	StockChangeAdjustment  StockChangeType = "ADJUSTMENT"
	StockChangeReservation StockChangeType = "RESERVATION"
	StockChangeRelease     StockChangeType = "RELEASE"
)

type StockAggregate struct {
	Uid           string
	EditionUid    string
	TotalQuantity int32
	Warehouses    []*StockWarehouse
}

type StockWarehouse struct {
	WarehouseUid  string
	WarehouseName string
	Quantity      int32
}

type StockHistoryEntry struct {
	Uid        string
	ChangeType StockChangeType
	Reason     *string
	ChangeQty  int32
	CreatedAt  time.Time
}
