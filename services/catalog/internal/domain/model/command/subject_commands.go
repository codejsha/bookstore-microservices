package command

type SubjectCreateCommand struct {
	Name string
}

type SubjectUpdateCommand struct {
	Name *string
}

func (c SubjectCreateCommand) Validate() error {
	if err := requireNonBlank("name", c.Name); err != nil {
		return err
	}
	return requireMaxLen("name", c.Name, 255)
}

func (c SubjectUpdateCommand) Validate() error {
	if err := optionalNonBlank("name", c.Name); err != nil {
		return err
	}
	return optionalMaxLen("name", c.Name, 255)
}
