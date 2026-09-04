package command

type PublisherCreateCommand struct {
	Name    string
	Address *string
	OlKey   *string
}

type PublisherUpdateCommand struct {
	Name    *string
	Address *string
	OlKey   *string
}

func (c PublisherCreateCommand) Validate() error {
	if err := requireNonBlank("name", c.Name); err != nil {
		return err
	}
	if err := requireMaxLen("name", c.Name, 255); err != nil {
		return err
	}
	if err := optionalMaxLen("address", c.Address, 255); err != nil {
		return err
	}
	if err := optionalNonBlank("ol_key", c.OlKey); err != nil {
		return err
	}
	return optionalMaxLen("ol_key", c.OlKey, 64)
}

func (c PublisherUpdateCommand) Validate() error {
	if err := optionalNonBlank("name", c.Name); err != nil {
		return err
	}
	if err := optionalMaxLen("name", c.Name, 255); err != nil {
		return err
	}
	if err := optionalMaxLen("address", c.Address, 255); err != nil {
		return err
	}
	if err := optionalNonBlank("ol_key", c.OlKey); err != nil {
		return err
	}
	return optionalMaxLen("ol_key", c.OlKey, 64)
}
