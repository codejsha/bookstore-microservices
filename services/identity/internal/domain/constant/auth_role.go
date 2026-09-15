package constant

import (
	"strings"

	"github.com/codejsha/bookstore-microservices/identity/generated/application/port/openapi"
)

type AuthRole int32

const (
	AUTHROLE_UNKNOWN AuthRole = 0
	AUTHROLE_USER    AuthRole = 1
	AUTHROLE_STAFF   AuthRole = 2
	AUTHROLE_MANAGER AuthRole = 3
	AUTHROLE_SYSTEM  AuthRole = 4
)

type AuthRoleValue string

const (
	AUTHROLE_UNKNOWN_VALUE AuthRoleValue = "UNKNOWN"
	AUTHROLE_USER_VALUE    AuthRoleValue = "USER"
	AUTHROLE_STAFF_VALUE   AuthRoleValue = "STAFF"
	AUTHROLE_MANAGER_VALUE AuthRoleValue = "MANAGER"
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
	case AUTHROLE_MANAGER:
		return openapi.AUTHROLE_MANAGER
	case AUTHROLE_SYSTEM:
		return openapi.AUTHROLE_SYSTEM
	default:
		return openapi.AUTHROLE_UNKNOWN
	}
}

func (a AuthRole) ToAuthRoleProto() string {
	switch a {
	case AUTHROLE_USER:
		return string(AUTHROLE_USER_VALUE)
	case AUTHROLE_STAFF:
		return string(AUTHROLE_STAFF_VALUE)
	case AUTHROLE_MANAGER:
		return string(AUTHROLE_MANAGER_VALUE)
	case AUTHROLE_SYSTEM:
		return string(AUTHROLE_SYSTEM_VALUE)
	default:
		return string(AUTHROLE_UNKNOWN_VALUE)
	}
}

func AuthRoleFromString(value string) AuthRole {
	switch strings.ToUpper(value) {
	case string(AUTHROLE_UNKNOWN_VALUE):
		return AUTHROLE_UNKNOWN
	case string(AUTHROLE_USER_VALUE):
		return AUTHROLE_USER
	case string(AUTHROLE_STAFF_VALUE):
		return AUTHROLE_STAFF
	case string(AUTHROLE_MANAGER_VALUE):
		return AUTHROLE_MANAGER
	case string(AUTHROLE_SYSTEM_VALUE):
		return AUTHROLE_SYSTEM
	default:
		return AUTHROLE_UNKNOWN
	}
}

type AuthRoleJson struct {
	Values []AuthRoleValue `json:"values"`
}
