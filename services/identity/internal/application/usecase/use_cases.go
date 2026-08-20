package usecase

import (
	"context"

	"github.com/codejsha/bookstore-microservices/identity/internal/domain/aggregate"
	"github.com/codejsha/bookstore-microservices/identity/internal/domain/model/command"
	"github.com/codejsha/bookstore-microservices/identity/internal/domain/model/option"
)

type IdentityUseCase interface {
	// ─── User registration ──────────────────────────────────────────────
	RegisterUser(ctx context.Context, cmd command.UserRegisterCommand) (*aggregate.UserAggregate, error)

	// ─── User query ─────────────────────────────────────────────────────
	FindAllUsers(ctx context.Context, opt option.UserQueryOption) (int64, []*aggregate.UserAggregate, error)
	FindUser(ctx context.Context, uid string) (*aggregate.UserAggregate, error)
	FindUserByEmail(ctx context.Context, email string) (*aggregate.UserAggregate, error)

	// ─── User management ────────────────────────────────────────────────
	UpdateUser(ctx context.Context, uid string, cmd command.UserUpdateCommand) (*aggregate.UserAggregate, error)
	UpdateUserRoles(ctx context.Context, uid string, cmd command.UserRolesCommand) (*aggregate.UserAggregate, error)
	SuspendUser(ctx context.Context, uid string) (*aggregate.UserAggregate, error)
	ReactivateUser(ctx context.Context, uid string) (*aggregate.UserAggregate, error)
	DeactivateUser(ctx context.Context, uid string) (*aggregate.UserAggregate, error)

	// ─── Keycloak sync ──────────────────────────────────────────────────
	SyncUserFromIdp(ctx context.Context, idpId string) (*aggregate.UserAggregate, error)
}
