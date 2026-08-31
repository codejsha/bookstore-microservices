package command

type AuthorCreateCommand struct {
	Name           string
	Bio            *string
	BirthDate      *string
	DeathDate      *string
	PhotoUids      []string
	AlternateNames []string
	OlKey          *string
}

type AuthorUpdateCommand struct {
	Name           *string
	Bio            *string
	BirthDate      *string
	DeathDate      *string
	PhotoUids      []string
	AlternateNames []string
	OlKey          *string
}

func (c AuthorCreateCommand) Validate() error {
	if err := requireNonBlank("name", c.Name); err != nil {
		return err
	}
	if err := requireMaxLen("name", c.Name, 255); err != nil {
		return err
	}
	if err := optionalMaxLen("birth_date", c.BirthDate, 32); err != nil {
		return err
	}
	if err := optionalMaxLen("death_date", c.DeathDate, 32); err != nil {
		return err
	}
	if err := optionalNonBlank("ol_key", c.OlKey); err != nil {
		return err
	}
	if err := optionalMaxLen("ol_key", c.OlKey, 64); err != nil {
		return err
	}
	return validateLifespan(c.BirthDate, c.DeathDate)
}

func (c AuthorUpdateCommand) Validate() error {
	if err := optionalNonBlank("name", c.Name); err != nil {
		return err
	}
	if err := optionalMaxLen("name", c.Name, 255); err != nil {
		return err
	}
	if err := optionalMaxLen("birth_date", c.BirthDate, 32); err != nil {
		return err
	}
	if err := optionalMaxLen("death_date", c.DeathDate, 32); err != nil {
		return err
	}
	if err := optionalNonBlank("ol_key", c.OlKey); err != nil {
		return err
	}
	if err := optionalMaxLen("ol_key", c.OlKey, 64); err != nil {
		return err
	}
	return validateLifespan(c.BirthDate, c.DeathDate)
}
