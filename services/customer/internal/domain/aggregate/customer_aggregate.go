package aggregate

import (
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/constant"
)

type CustomerAggregate struct {
	Uid       string
	Email     string
	FirstName string
	LastName  string
	Phone     *string
	Roles     []constant.AuthRole
}
