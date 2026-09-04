package command

type StockTransferCommand struct {
	EditionUid         string
	SourceWarehouseUid string
	TargetWarehouseUid string
	Quantity           int32
	Reason             *string
}

func (c StockTransferCommand) Validate() error {
	if err := requireNonBlank("edition_uid", c.EditionUid); err != nil {
		return err
	}
	if err := requireUUID("edition_uid", c.EditionUid); err != nil {
		return err
	}
	if err := requireNonBlank("source_warehouse_uid", c.SourceWarehouseUid); err != nil {
		return err
	}
	if err := requireUUID("source_warehouse_uid", c.SourceWarehouseUid); err != nil {
		return err
	}
	if err := requireNonBlank("target_warehouse_uid", c.TargetWarehouseUid); err != nil {
		return err
	}
	if err := requireUUID("target_warehouse_uid", c.TargetWarehouseUid); err != nil {
		return err
	}
	if c.SourceWarehouseUid == c.TargetWarehouseUid {
		return invalidf("source_warehouse_uid and target_warehouse_uid must differ")
	}
	if c.Quantity <= 0 {
		return invalidf("quantity must be positive: got %d", c.Quantity)
	}
	if err := optionalNonBlank("reason", c.Reason); err != nil {
		return err
	}
	return optionalMaxLen("reason", c.Reason, 255)
}
