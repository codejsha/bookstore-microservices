package command

type StockReceiveCommand struct {
	EditionUid   string
	WarehouseUid string
	Quantity     int32
	Reason       *string
}

type StockReleaseCommand struct {
	EditionUid   string
	WarehouseUid string
	Quantity     int32
	Reason       *string
}

type StockAdjustCommand struct {
	EditionUid   string
	WarehouseUid string
	Quantity     int32
	Reason       string
}

type StockReserveCommand struct {
	EditionUid   string
	WarehouseUid string
	Quantity     int32
	Reason       *string
}

type StockOrderReserveCommand struct {
	OrderUid  string
	EditionId int64
	Quantity  int32
	Reason    *string
}

type StockOrderReleaseCommand struct {
	OrderUid  string
	EditionId int64
	Quantity  int32
	Reason    *string
}
