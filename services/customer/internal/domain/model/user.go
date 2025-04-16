package model

import (
	"github.com/codejsha/bookstore-microservices/customer/internal/application/port/pb/userpb"
)

type UserProto interface {
	GetId() string
	GetEmail() string
	GetFirstName() string
	GetLastName() string
	GetPhone() string
	GetRoles() []userpb.AuthRole
}
