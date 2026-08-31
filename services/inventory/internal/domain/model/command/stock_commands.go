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

func (c StockReceiveCommand) Validate() error {
	if err := validateStockTarget(c.EditionUid, c.WarehouseUid); err != nil {
		return err
	}
	if c.Quantity <= 0 {
		return invalidf("quantity must be positive: got %d", c.Quantity)
	}
	if err := optionalNonBlank("reason", c.Reason); err != nil {
		return err
	}
	return optionalMaxLen("reason", c.Reason, 255)
}

func (c StockReleaseCommand) Validate() error {
	if err := validateStockTarget(c.EditionUid, c.WarehouseUid); err != nil {
		return err
	}
	if c.Quantity <= 0 {
		return invalidf("quantity must be positive: got %d", c.Quantity)
	}
	if err := optionalNonBlank("reason", c.Reason); err != nil {
		return err
	}
	return optionalMaxLen("reason", c.Reason, 255)
}

func (c StockAdjustCommand) Validate() error {
	if err := validateStockTarget(c.EditionUid, c.WarehouseUid); err != nil {
		return err
	}
	if c.Quantity == 0 {
		return invalidf("quantity must be non-zero")
	}
	if err := requireNonBlank("reason", c.Reason); err != nil {
		return err
	}
	return requireMaxLen("reason", c.Reason, 255)
}

func (c StockReserveCommand) Validate() error {
	if err := validateStockTarget(c.EditionUid, c.WarehouseUid); err != nil {
		return err
	}
	if c.Quantity <= 0 {
		return invalidf("quantity must be positive: got %d", c.Quantity)
	}
	if err := optionalNonBlank("reason", c.Reason); err != nil {
		return err
	}
	return optionalMaxLen("reason", c.Reason, 255)
}

func (c StockOrderReserveCommand) Validate() error {
	if err := validateOrderReference(c.OrderUid, c.EditionId); err != nil {
		return err
	}
	if err := optionalNonBlank("reason", c.Reason); err != nil {
		return err
	}
	return optionalMaxLen("reason", c.Reason, 255)
}

func (c StockOrderReleaseCommand) Validate() error {
	if err := validateOrderReference(c.OrderUid, c.EditionId); err != nil {
		return err
	}
	if err := optionalNonBlank("reason", c.Reason); err != nil {
		return err
	}
	return optionalMaxLen("reason", c.Reason, 255)
}

func validateStockTarget(editionUid, warehouseUid string) error {
	if err := requireNonBlank("edition_uid", editionUid); err != nil {
		return err
	}
	if err := requireUUID("edition_uid", editionUid); err != nil {
		return err
	}
	if err := requireNonBlank("warehouse_uid", warehouseUid); err != nil {
		return err
	}
	return requireUUID("warehouse_uid", warehouseUid)
}

func validateOrderReference(orderUid string, editionId int64) error {
	if err := requireNonBlank("order_uid", orderUid); err != nil {
		return err
	}
	if err := requireMaxLen("order_uid", orderUid, 200); err != nil {
		return err
	}
	if editionId <= 0 {
		return invalidf("edition_id must be positive: got %d", editionId)
	}
	return nil
}
