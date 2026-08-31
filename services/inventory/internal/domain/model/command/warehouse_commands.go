package command

type WarehouseCreateCommand struct {
	Name     string
	Address  *string
	Capacity int32
}

type WarehouseUpdateCommand struct {
	Name     *string
	Address  *string
	Capacity *int32
}

func (c WarehouseCreateCommand) Validate() error {
	if err := requireNonBlank("name", c.Name); err != nil {
		return err
	}
	if err := requireMaxLen("name", c.Name, 255); err != nil {
		return err
	}
	if err := optionalNonBlank("address", c.Address); err != nil {
		return err
	}
	if err := optionalMaxLen("address", c.Address, 255); err != nil {
		return err
	}
	if c.Capacity < 0 {
		return invalidf("capacity must not be negative: got %d", c.Capacity)
	}
	return nil
}

func (c WarehouseUpdateCommand) Validate() error {
	if err := optionalNonBlank("name", c.Name); err != nil {
		return err
	}
	if err := optionalMaxLen("name", c.Name, 255); err != nil {
		return err
	}
	if err := optionalNonBlank("address", c.Address); err != nil {
		return err
	}
	if err := optionalMaxLen("address", c.Address, 255); err != nil {
		return err
	}
	if c.Capacity != nil && *c.Capacity < 0 {
		return invalidf("capacity must not be negative: got %d", *c.Capacity)
	}
	return nil
}
