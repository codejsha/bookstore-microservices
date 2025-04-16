package usecase

import (
	"context"

	"github.com/codejsha/bookstore-microservices/identity/internal/application/port/openapi"
	"github.com/codejsha/bookstore-microservices/identity/internal/application/port/pb/userpb"
	"github.com/codejsha/bookstore-microservices/identity/internal/domain/aggregate"
)

type IdentityUseCase interface {
	SignIn(ctx context.Context, req openapi.SignInWebReq) (*aggregate.TokenAggregate, error)
	SignOut(ctx context.Context, req openapi.SignOutWebParam) error
	SignUp(ctx context.Context, req openapi.SignUpWebReq) (*aggregate.UserAggregate, error)
	ExchangeToken(ctx context.Context, req openapi.TokenExchangeWebParam) (*aggregate.TokenAggregate, error)
	FindAllUsers(ctx context.Context, req *userpb.UserFindAllProtoReq) (int64, []*aggregate.UserAggregate, error)
}
