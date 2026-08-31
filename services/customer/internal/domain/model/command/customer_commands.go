package command

import (
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/constant"
)

type CustomerUpdateCommand struct {
	Uid       string
	Email     *string
	Password  *string
	FirstName *string
	LastName  *string
	Phone     *string
	Roles     []constant.AuthRole
}

func (c CustomerUpdateCommand) Validate() error {
	if err := requireNonBlank("uid", c.Uid); err != nil {
		return err
	}
	if err := optionalEmail("email", c.Email); err != nil {
		return err
	}
	if err := optionalNonBlank("password", c.Password); err != nil {
		return err
	}
	if err := optionalNonBlank("first_name", c.FirstName); err != nil {
		return err
	}
	if err := optionalNonBlank("last_name", c.LastName); err != nil {
		return err
	}
	if err := optionalNonBlank("phone", c.Phone); err != nil {
		return err
	}
	return requireKnownRoles("roles", c.Roles)
}
