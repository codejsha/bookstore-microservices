package constant

import (
	"github.com/codejsha/bookstore-microservices/customer/generated/application/port/openapi"
)

type AuthRole int32

const (
	AUTHROLE_UNKNOWN AuthRole = 0
	AUTHROLE_SYSTEM  AuthRole = 1
	AUTHROLE_MANAGE  AuthRole = 2
	AUTHROLE_PROFILE AuthRole = 3
	AUTHROLE_ORDER   AuthRole = 4
	AUTHROLE_VIEW    AuthRole = 5
)

type AuthRoleValue string

const (
	AUTHROLE_UNKNOWN_VALUE AuthRoleValue = "UNKNOWN"
	AUTHROLE_SYSTEM_VALUE  AuthRoleValue = "SYSTEM"
	AUTHROLE_MANAGE_VALUE  AuthRoleValue = "MANAGE"
	AUTHROLE_PROFILE_VALUE AuthRoleValue = "PROFILE"
	AUTHROLE_ORDER_VALUE   AuthRoleValue = "ORDER"
	AUTHROLE_VIEW_VALUE    AuthRoleValue = "VIEW"
)

func (a AuthRole) ToAuthRoleRest() openapi.AuthRole {
	switch a {
	case AUTHROLE_UNKNOWN:
		return openapi.AUTHROLE_UNKNOWN
	case AUTHROLE_SYSTEM:
		return openapi.AUTHROLE_SYSTEM
	case AUTHROLE_MANAGE:
		return openapi.AUTHROLE_MANAGE
	case AUTHROLE_PROFILE:
		return openapi.AUTHROLE_PROFILE
	case AUTHROLE_ORDER:
		return openapi.AUTHROLE_ORDER
	case AUTHROLE_VIEW:
		return openapi.AUTHROLE_VIEW
	default:
		return openapi.AUTHROLE_UNKNOWN
	}
}

func AuthRoleFromString(value string) AuthRole {
	switch value {
	case string(AUTHROLE_UNKNOWN_VALUE):
		return AUTHROLE_UNKNOWN
	case string(AUTHROLE_SYSTEM_VALUE):
		return AUTHROLE_SYSTEM
	case string(AUTHROLE_MANAGE_VALUE):
		return AUTHROLE_MANAGE
	case string(AUTHROLE_PROFILE_VALUE):
		return AUTHROLE_PROFILE
	case string(AUTHROLE_ORDER_VALUE):
		return AUTHROLE_ORDER
	case string(AUTHROLE_VIEW_VALUE):
		return AUTHROLE_VIEW
	default:
		return AUTHROLE_UNKNOWN
	}
}
