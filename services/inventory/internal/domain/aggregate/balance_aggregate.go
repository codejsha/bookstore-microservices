package aggregate

type StockBalanceEntry struct {
	EditionUid       string
	WarehouseUid     string
	WarehouseName    string
	InboundQuantity  int32
	OutboundQuantity int32
	AdjustQuantity   int32
	CurrentQuantity  int32
}
