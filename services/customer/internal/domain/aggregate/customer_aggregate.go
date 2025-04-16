package aggregate

import (
	"github.com/codejsha/bookstore-microservices/customer/internal/application/port/openapi"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/model"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/model/mapper"
)

func NewCustomerAggregate(
	customer model.UserProto,
	userMapper *mapper.UserMapper,
) *CustomerAggregate {
	roles := userMapper.ToAuthRoles(customer.GetRoles())
	var phone *string
	if p := customer.GetPhone(); p != "" {
		phone = &p
	}

	return &CustomerAggregate{
		Id:        customer.GetId(),
		Email:     customer.GetEmail(),
		FirstName: customer.GetFirstName(),
		LastName:  customer.GetLastName(),
		Phone:     phone,
		Roles:     roles,
	}
}

type CustomerAggregate struct {
	// user id
	Id        string
	Email     string
	FirstName string
	LastName  string
	Phone     *string
	Roles     []openapi.AuthRole
}

func (a *CustomerAggregate) ToCustomerFindWebResp() openapi.CustomerFindWebResp {
	resp := openapi.CustomerFindWebResp{
		Id:        a.Id,
		Email:     a.Email,
		FirstName: a.FirstName,
		LastName:  a.LastName,
		Phone:     a.Phone,
		Roles:     a.Roles,
	}
	return resp
}
