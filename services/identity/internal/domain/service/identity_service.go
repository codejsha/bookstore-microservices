package service

import (
	"context"
	"fmt"

	"github.com/codejsha/bookstore-microservices/identity/internal/application/port/idp"
	"github.com/codejsha/bookstore-microservices/identity/internal/application/port/openapi"
	"github.com/codejsha/bookstore-microservices/identity/internal/application/port/pb/userpb"
	"github.com/codejsha/bookstore-microservices/identity/internal/application/usecase"
	"github.com/codejsha/bookstore-microservices/identity/internal/config"
	"github.com/codejsha/bookstore-microservices/identity/internal/domain/aggregate"
	"github.com/codejsha/bookstore-microservices/identity/internal/domain/model/mapper/oidcmapper"
)

var _ usecase.IdentityUseCase = (*identityService)(nil)

func NewIdentityService(
	cfg *config.Config,
	oidcAPI idp.OpenIDConnectAPI,
	usersAPI idp.UsersAPI,
) usecase.IdentityUseCase {
	return &identityService{
		cfg:         cfg,
		oidcClient:  oidcAPI,
		usersClient: usersAPI,
	}
}

type identityService struct {
	cfg         *config.Config
	oidcClient  idp.OpenIDConnectAPI
	usersClient idp.UsersAPI
}

func (h identityService) SignIn(ctx context.Context, req openapi.SignInWebReq) (*aggregate.TokenAggregate, error) {
	request, err := oidcmapper.ToTokenExchangeParam(h.cfg.Keycloak, req)
	if err != nil {
		return nil, err
	}
	response, err := h.oidcClient.RealmsRealmProtocolOpenidConnectTokenPost(ctx, request)
	if err != nil {
		return nil, err
	}
	var tokenAggregate aggregate.TokenAggregate
	tokenAgg, err := tokenAggregate.FromValue(response)
	if err != nil {
		return nil, err
	}
	return tokenAgg, nil
}

func (h identityService) SignOut(ctx context.Context, req openapi.SignOutWebParam) error {
	request, err := oidcmapper.ToLogoutParam(h.cfg.Keycloak, req)
	if err != nil {
		return err
	}
	err = h.oidcClient.RealmsRealmProtocolOpenidConnectLogoutPost(ctx, request)
	if err != nil {
		return err
	}
	return nil
}

func (h identityService) SignUp(ctx context.Context, req openapi.SignUpWebReq) (*aggregate.UserAggregate, error) {
	// check if user already exists
	usersGetParams := idp.AdminRealmsRealmUsersGetParam{Realm: h.cfg.Keycloak.Realm, Email: req.Email}
	users, err := h.usersClient.AdminRealmsRealmUsersGet(ctx, usersGetParams)
	if err != nil {
		return nil, err
	}
	if len(users) > 0 {
		return nil, fmt.Errorf("failed to create user: user already exists with email %s", req.Email)
	}

	// create user
	params := idp.AdminRealmsRealmUsersPostParam{Realm: h.cfg.Keycloak.Realm}
	request, err := oidcmapper.ToUserRepresentation(req)
	err = h.usersClient.AdminRealmsRealmUsersPost(ctx, params, request)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// get user
	users, err = h.usersClient.AdminRealmsRealmUsersGet(ctx, usersGetParams)
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, fmt.Errorf("failed to create user: user not found after creation")
	}

	// return response
	userAgg := aggregate.NewUserAggregate(users[0])
	return userAgg, nil
}

func (h identityService) ExchangeToken(ctx context.Context, req openapi.TokenExchangeWebParam) (*aggregate.TokenAggregate, error) {
	request, err := oidcmapper.ToTokenExchangeParam(h.cfg.Keycloak, req)
	if err != nil {
		return nil, err
	}
	response, err := h.oidcClient.RealmsRealmProtocolOpenidConnectTokenPost(ctx, request)
	if err != nil {
		return nil, err
	}
	var tokenAggregate aggregate.TokenAggregate
	tokenAgg, err := tokenAggregate.FromValue(response)
	if err != nil {
		return nil, err
	}
	return tokenAgg, nil
}

func (h identityService) FindAllUsers(ctx context.Context, req *userpb.UserFindAllProtoReq) (int64, []*aggregate.UserAggregate, error) {
	// find users
	usersGetParams := idp.AdminRealmsRealmUsersGetParam{
		Realm: h.cfg.Keycloak.Realm,
		Email: *req.Email,
	}
	users, err := h.usersClient.AdminRealmsRealmUsersGet(ctx, usersGetParams)
	if err != nil {
		return 0, nil, err
	}

	// return response
	userAggs := make([]*aggregate.UserAggregate, 0, len(users))
	for i, item := range users {
		userAggs[i] = aggregate.NewUserAggregate(item)
	}
	return int64(len(users)), userAggs, nil
}
