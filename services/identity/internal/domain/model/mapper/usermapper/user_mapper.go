package usermapper

import (
	"github.com/codejsha/bookstore-microservices/identity/internal/application/port/openapi"
	"github.com/codejsha/bookstore-microservices/identity/internal/application/port/pb/userpb"
)

func ToAuthRolesPb(values []openapi.AuthRole) []userpb.AuthRole {
	roles := make([]userpb.AuthRole, len(values))
	for i, role := range values {
		r := ToAuthRolePb(role)
		roles[i] = r
	}
	return roles
}

func ToAuthRolePb(value openapi.AuthRole) userpb.AuthRole {
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

func ToAuthRoles(values []userpb.AuthRole) []openapi.AuthRole {
	roles := make([]openapi.AuthRole, len(values))
	for i, role := range values {
		r := ToAuthRole(role)
		roles[i] = r
	}
	return roles
}

func ToAuthRole(value userpb.AuthRole) openapi.AuthRole {
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
