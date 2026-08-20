package repo

import (
	"context"

	"github.com/codejsha/bookstore-microservices/identity/internal/domain/model/option"
)

type UserRepo interface {
	FindAll(ctx context.Context, opt option.UserQueryOption) (int64, []*UserResult, error)
	FindByUid(ctx context.Context, uid string) (*UserResult, error)
	FetchByEmail(ctx context.Context, email string) ([]*UserResult, error)

	Upsert(ctx context.Context, profile UserUpsert) (*UserResult, error)

	UpdateProfile(ctx context.Context, idpUid string, patch UserProfileUpdate) error

	UpdateStatus(ctx context.Context, idpUid string, status string) error

	UpdateRoles(ctx context.Context, idpUid string, roles []string) error

	SoftDelete(ctx context.Context, idpUid string) error
}

type UserUpsert struct {
	IdpUid    string
	Email     string
	FirstName string
	LastName  string
	Phone     *string
	Roles     []string
	Status    string
}

type UserProfileUpdate struct {
	FirstName *string
	LastName  *string
	Phone     *string
}
