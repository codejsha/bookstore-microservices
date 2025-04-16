package mapper

import (
	"fmt"

	"github.com/codejsha/bookstore-microservices/customer/internal/application/port/openapi"
	"github.com/codejsha/bookstore-microservices/customer/internal/application/port/pb/userpb"
)

func NewUserMapper() *UserMapper {
	return &UserMapper{}
}

type UserMapper struct{}

func (m UserMapper) ToUserFindAllProtoReq(value interface{}) (*userpb.UserFindAllProtoReq, error) {
	switch value.(type) {
	case openapi.CustomerFindAllWebParam:
		v := value.(openapi.CustomerFindAllWebParam)
		request := &userpb.UserFindAllProtoReq{
			Email: v.Email,
			Name:  v.Name,
			Phone: v.Phone,
			Filter: &userpb.QueryFilter{
				Sort:   v.Sort,
				Order:  v.Order,
				Limit:  v.Limit,
				Offset: v.Offset,
			},
		}
		return request, nil
	default:
		return nil, fmt.Errorf("failed to convert request to UserFindAllRequest: %v", value)
	}
}

func (m UserMapper) ToUserFindProtoReq(value interface{}) (*userpb.UserFindProtoReq, error) {
	switch value.(type) {
	case openapi.CustomerFindWebParam:
		v := value.(openapi.CustomerFindWebParam)
		request := &userpb.UserFindProtoReq{
			Id: v.Id,
		}
		return request, nil
	default:
		return nil, fmt.Errorf("failed to convert request to UserFindRequest: %v", value)
	}
}

func (m UserMapper) ToAuthRolesPb(values []openapi.AuthRole) []userpb.AuthRole {
	roles := make([]userpb.AuthRole, len(values))
	for i, role := range values {
		r := m.ToAuthRolePb(role)
		roles[i] = r
	}
	return roles
}

func (m UserMapper) ToAuthRolePb(value openapi.AuthRole) userpb.AuthRole {
	switch value {
	case openapi.AUTHROLE_UNKNOWN:
		return userpb.AuthRole_UNKNOWN
	case openapi.AUTHROLE_SYSTEM:
		return userpb.AuthRole_SYSTEM
	case openapi.AUTHROLE_MANAGE:
		return userpb.AuthRole_MANAGE
	case openapi.AUTHROLE_PROFILE:
		return userpb.AuthRole_PROFILE
	case openapi.AUTHROLE_ORDER:
		return userpb.AuthRole_ORDER
	case openapi.AUTHROLE_VIEW:
		return userpb.AuthRole_VIEW
	default:
		return userpb.AuthRole_UNKNOWN
	}
}

func (m UserMapper) ToAuthRoles(values []userpb.AuthRole) []openapi.AuthRole {
	roles := make([]openapi.AuthRole, len(values))
	for i, role := range values {
		r := m.ToAuthRole(role)
		roles[i] = r
	}
	return roles
}

func (m UserMapper) ToAuthRole(value userpb.AuthRole) openapi.AuthRole {
	switch value {
	case userpb.AuthRole_UNKNOWN:
		return openapi.AUTHROLE_UNKNOWN
	case userpb.AuthRole_SYSTEM:
		return openapi.AUTHROLE_SYSTEM
	case userpb.AuthRole_MANAGE:
		return openapi.AUTHROLE_MANAGE
	case userpb.AuthRole_PROFILE:
		return openapi.AUTHROLE_PROFILE
	case userpb.AuthRole_ORDER:
		return openapi.AUTHROLE_ORDER
	case userpb.AuthRole_VIEW:
		return openapi.AUTHROLE_VIEW
	default:
		return openapi.AUTHROLE_UNKNOWN
	}
}

func (m UserMapper) ToUserUpdateProtoReq(value interface{}, id string) (*userpb.UserUpdateProtoReq, error) {
	switch value.(type) {
	case openapi.CustomerUpdateWebReq:
		v := value.(openapi.CustomerUpdateWebReq)
		roles := m.ToAuthRolesPb(v.Roles)
		request := &userpb.UserUpdateProtoReq{
			Id:        id,
			Email:     v.Email,
			FirstName: v.FirstName,
			LastName:  v.LastName,
			Phone:     v.Phone,
			Roles:     roles,
		}
		return request, nil
	default:
		return nil, fmt.Errorf("failed to convert request to UserUpdateRequest: %v", value)
	}
}
