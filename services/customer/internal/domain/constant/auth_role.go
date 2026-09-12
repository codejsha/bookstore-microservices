package constant

import (
	"github.com/codejsha/bookstore-microservices/customer/generated/application/port/openapi"
)

type AuthRole int32

const (
	AUTHROLE_UNKNOWN AuthRole = 0
	AUTHROLE_USER    AuthRole = 1
	AUTHROLE_STAFF   AuthRole = 2
	AUTHROLE_MANAGE  AuthRole = 3
	AUTHROLE_SYSTEM  AuthRole = 4
)

type AuthRoleValue string

const (
	AUTHROLE_UNKNOWN_VALUE AuthRoleValue = "UNKNOWN"
	AUTHROLE_USER_VALUE    AuthRoleValue = "USER"
	AUTHROLE_STAFF_VALUE   AuthRoleValue = "STAFF"
	AUTHROLE_MANAGE_VALUE  AuthRoleValue = "MANAGE"
	AUTHROLE_SYSTEM_VALUE  AuthRoleValue = "SYSTEM"
)

func (a AuthRole) ToAuthRoleRest() openapi.AuthRole {
	switch a {
	case AUTHROLE_UNKNOWN:
		return openapi.AUTHROLE_UNKNOWN
	case AUTHROLE_USER:
		return openapi.AUTHROLE_USER
	case AUTHROLE_STAFF:
		return openapi.AUTHROLE_STAFF
	case AUTHROLE_MANAGE:
		return openapi.AUTHROLE_MANAGE
	case AUTHROLE_SYSTEM:
		return openapi.AUTHROLE_SYSTEM
	default:
		return openapi.AUTHROLE_UNKNOWN
	}
}

func AuthRoleFromString(value string) AuthRole {
	switch value {
	case string(AUTHROLE_UNKNOWN_VALUE):
		return AUTHROLE_UNKNOWN
	case string(AUTHROLE_USER_VALUE):
		return AUTHROLE_USER
	case string(AUTHROLE_STAFF_VALUE):
		return AUTHROLE_STAFF
	case string(AUTHROLE_MANAGE_VALUE):
		return AUTHROLE_MANAGE
	case string(AUTHROLE_SYSTEM_VALUE):
		return AUTHROLE_SYSTEM
	default:
		return AUTHROLE_UNKNOWN
	}
}
