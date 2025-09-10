package command

type StockTransferCommand struct {
	EditionUid         string
	SourceWarehouseUid string
	TargetWarehouseUid string
	Quantity           int32
	Reason             *string
}
