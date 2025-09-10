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
