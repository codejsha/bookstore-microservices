package command

const (
	maxEmailLen = 100
	maxNameLen  = 50
	maxPhoneLen = 30
)

type UserRegisterCommand struct {
	Email     string
	Password  string
	FirstName string
	LastName  string
	Phone     *string
	Roles     []string
}

type UserUpdateCommand struct {
	FirstName *string
	LastName  *string
	Phone     *string
}

type UserRolesCommand struct {
	Roles []string
}

func (c UserRegisterCommand) Validate() error {
	if err := requireEmail("email", c.Email, maxEmailLen); err != nil {
		return err
	}
	if err := requireNonBlank("password", c.Password); err != nil {
		return err
	}
	if err := requireNonBlank("first_name", c.FirstName); err != nil {
		return err
	}
	if err := requireMaxLen("first_name", c.FirstName, maxNameLen); err != nil {
		return err
	}
	if err := requireNonBlank("last_name", c.LastName); err != nil {
		return err
	}
	if err := requireMaxLen("last_name", c.LastName, maxNameLen); err != nil {
		return err
	}
	if err := optionalNonBlank("phone", c.Phone); err != nil {
		return err
	}
	if err := optionalMaxLen("phone", c.Phone, maxPhoneLen); err != nil {
		return err
	}
	return requireKnownRoles("roles", c.Roles)
}

func (c UserUpdateCommand) Validate() error {
	if err := optionalNonBlank("first_name", c.FirstName); err != nil {
		return err
	}
	if err := optionalMaxLen("first_name", c.FirstName, maxNameLen); err != nil {
		return err
	}
	if err := optionalNonBlank("last_name", c.LastName); err != nil {
		return err
	}
	if err := optionalMaxLen("last_name", c.LastName, maxNameLen); err != nil {
		return err
	}
	if err := optionalNonBlank("phone", c.Phone); err != nil {
		return err
	}
	return optionalMaxLen("phone", c.Phone, maxPhoneLen)
}

func (c UserRolesCommand) Validate() error {
	if len(c.Roles) == 0 {
		return invalidf("roles must not be empty")
	}
	return requireKnownRoles("roles", c.Roles)
}
