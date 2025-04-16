package aggregate

import (
	"github.com/codejsha/bookstore-microservices/identity/internal/application/port/idp"
	"github.com/codejsha/bookstore-microservices/identity/internal/application/port/openapi"
	"github.com/codejsha/bookstore-microservices/identity/internal/application/port/pb/userpb"
	"github.com/codejsha/bookstore-microservices/identity/internal/domain/model/mapper/usermapper"
)

func NewUserAggregate(
	userRepresentation idp.UserRepresentation,
) *UserAggregate {
	roles := make([]openapi.AuthRole, len(userRepresentation.RealmRoles))
	for i, role := range userRepresentation.RealmRoles {
		roles[i] = openapi.AuthRole(role)
	}

	agg := &UserAggregate{
		Id:        userRepresentation.Id,
		Email:     *userRepresentation.Email,
		FirstName: *userRepresentation.FirstName,
		LastName:  *userRepresentation.LastName,
		Phone:     nil,
		Roles:     roles,
	}
	return agg
}

type UserAggregate struct {
	Id        *string
	Email     string
	FirstName string
	LastName  string
	Phone     *string
	Roles     []openapi.AuthRole
}

func (a *UserAggregate) ToSignUpWebResp() openapi.SignUpWebResp {
	return openapi.SignUpWebResp{
		Id: *a.Id,
	}
}

func (a *UserAggregate) ToUserFindProtoResp() *userpb.UserFindProtoResp {
	return &userpb.UserFindProtoResp{
		Id:        *a.Id,
		Email:     a.Email,
		FirstName: a.FirstName,
		LastName:  a.LastName,
		Phone:     a.Phone,
		Roles:     usermapper.ToAuthRolesPb(a.Roles),
	}
}
