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

func (c StockAuditCreateCommand) Validate() error {
	if err := requireNonBlank("warehouse_uid", c.WarehouseUid); err != nil {
		return err
	}
	if err := requireUUID("warehouse_uid", c.WarehouseUid); err != nil {
		return err
	}
	if len(c.Items) == 0 {
		return invalidf("items must not be empty")
	}
	seen := make(map[string]struct{}, len(c.Items))
	for _, it := range c.Items {
		if err := requireNonBlank("items.edition_uid", it.EditionUid); err != nil {
			return err
		}
		if err := requireUUID("items.edition_uid", it.EditionUid); err != nil {
			return err
		}
		if it.ActualQuantity < 0 {
			return invalidf("items.actual_quantity must not be negative: got %d", it.ActualQuantity)
		}
		if _, ok := seen[it.EditionUid]; ok {
			return invalidf("items must not contain duplicate edition_uid: %s", it.EditionUid)
		}
		seen[it.EditionUid] = struct{}{}
	}
	if err := optionalNonBlank("notes", c.Notes); err != nil {
		return err
	}
	return optionalMaxLen("notes", c.Notes, 500)
}
