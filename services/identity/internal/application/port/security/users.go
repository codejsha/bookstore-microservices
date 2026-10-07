package security

import (
	"context"
	"errors"

	idp "github.com/codejsha/bookstore-microservices/identity/generated/application/port/idpapi"
)

var (
	ErrIdpUserNotFound = errors.New("identity provider: user not found")
	ErrIdpUserConflict = errors.New("identity provider: user conflict")
	ErrIdpBadRequest   = errors.New("identity provider: bad request")
)

type UsersClient interface {
	ListUsers(ctx context.Context, realm string, email string) ([]idp.UserRepresentation, error)
	CreateUser(ctx context.Context, realm string, req idp.UserRepresentation) error
	GetUser(ctx context.Context, realm string, userId string) (idp.UserRepresentation, error)
	UpdateUser(ctx context.Context, realm string, userId string, req idp.UserRepresentation) error
	DeleteUser(ctx context.Context, realm string, userId string) error

	LogoutUser(ctx context.Context, realm string, userId string) error

	GetUserRealmRoles(ctx context.Context, realm string, userId string) ([]string, error)
	SetUserRealmRoles(ctx context.Context, realm string, userId string, roles []string) error
}
