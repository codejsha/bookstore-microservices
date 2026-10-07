package aggregate

import (
	"time"

	"github.com/codejsha/bookstore-microservices/identity/internal/application/port/repo"
	"github.com/codejsha/bookstore-microservices/identity/internal/domain/constant"
)

type UserAggregate struct {
	id          int64
	idpId       *string
	email       string
	firstName   string
	lastName    string
	phone       *string
	roles       []constant.AuthRole
	status      string
	lastLoginAt *time.Time
	createdAt   time.Time
	updatedAt   *time.Time
}

func NewUserAggregate(result *repo.UserResult) *UserAggregate {
	roles := make([]constant.AuthRole, len(result.Roles))
	for i, role := range result.Roles {
		roles[i] = constant.AuthRoleFromString(role)
	}
	return &UserAggregate{
		id:          result.Id,
		idpId:       result.IdpId,
		email:       result.Email,
		firstName:   result.FirstName,
		lastName:    result.LastName,
		phone:       result.Phone,
		roles:       roles,
		status:      result.Status,
		lastLoginAt: result.LastLoginAt,
		createdAt:   result.CreatedAt,
		updatedAt:   result.UpdatedAt,
	}
}

func (a *UserAggregate) Uid() string {
	if a.idpId == nil {
		return ""
	}
	return *a.idpId
}

func (a *UserAggregate) IdpId() *string {
	return a.idpId
}

func (a *UserAggregate) Email() string {
	return a.email
}

func (a *UserAggregate) FirstName() string {
	return a.firstName
}

func (a *UserAggregate) LastName() string {
	return a.lastName
}

func (a *UserAggregate) Phone() *string {
	return a.phone
}

func (a *UserAggregate) Roles() []constant.AuthRole {
	return a.roles
}

func (a *UserAggregate) Status() string {
	return a.status
}

func (a *UserAggregate) LastLoginAt() *time.Time {
	return a.lastLoginAt
}

func (a *UserAggregate) CreatedAt() time.Time {
	return a.createdAt
}

func (a *UserAggregate) UpdatedAt() *time.Time {
	return a.updatedAt
}
