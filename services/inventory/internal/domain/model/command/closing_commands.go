package command

type MonthlyClosingCommand struct {
	WarehouseUid string
	Year         int32
	Month        int32
}

func (c MonthlyClosingCommand) Validate() error {
	if err := requireNonBlank("warehouse_uid", c.WarehouseUid); err != nil {
		return err
	}
	if err := requireUUID("warehouse_uid", c.WarehouseUid); err != nil {
		return err
	}
	if c.Year < 2000 || c.Year > 2100 {
		return invalidf("year must be between 2000 and 2100: got %d", c.Year)
	}
	if c.Month < 1 || c.Month > 12 {
		return invalidf("month must be between 1 and 12: got %d", c.Month)
	}
	return nil
}
