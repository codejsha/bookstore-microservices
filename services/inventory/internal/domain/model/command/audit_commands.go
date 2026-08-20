package command

type StockAuditCreateCommand struct {
	WarehouseUid string
	Items        []StockAuditItemCommand
	Notes        *string
}

type StockAuditItemCommand struct {
	EditionUid     string
	ActualQuantity int32
}
