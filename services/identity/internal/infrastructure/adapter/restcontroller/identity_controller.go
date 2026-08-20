package restcontroller

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/codejsha/shared-library-go/pkg/pagination"

	"github.com/codejsha/bookstore-microservices/identity/generated/application/port/openapi"
	"github.com/codejsha/bookstore-microservices/identity/internal/application/usecase"
	"github.com/codejsha/bookstore-microservices/identity/internal/domain/aggregate"
	"github.com/codejsha/bookstore-microservices/identity/internal/domain/model/command"
	"github.com/codejsha/bookstore-microservices/identity/internal/domain/model/option"
	"github.com/codejsha/bookstore-microservices/identity/internal/infrastructure/httpx"
)

var _ openapi.UserApi = (*userController)(nil)

type userController struct {
	identityUseCase usecase.IdentityUseCase
}

func NewUserController(identityUseCase usecase.IdentityUseCase) openapi.UserApi {
	return &userController{identityUseCase: identityUseCase}
}

func (c *userController) UsersGetAll(
	ctx context.Context,
	email *string,
	name *string,
	phone *string,
	_ *openapi.UserStatus,
	size *int32,
	page *int32,
	sort *string,
) (*openapi.UserFindAllResponse, error) {
	opt := option.NewUserQueryOption(
		option.UserQueryOption{}.WithEmail(email),
		option.UserQueryOption{}.WithName(name),
		option.UserQueryOption{}.WithPhone(phone),
		option.UserQueryOption{}.WithPage(pagination.NewPageOption(size, page, sort)),
	)

	total, users, err := c.identityUseCase.FindAllUsers(ctx, opt)
	if err != nil {
		return nil, httpx.MapError(ctx, err)
	}

	items := make([]openapi.UserFindResponse, len(users))
	for i, u := range users {
		items[i] = toUserFindResponse(u)
	}
	return &openapi.UserFindAllResponse{Items: items, Total: total}, nil
}

func (c *userController) UsersRegister(
	ctx context.Context,
	req openapi.UserRegisterRequest,
) (*openapi.UserFindResponse, error) {
	cmd := command.UserRegisterCommand{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		Roles:     authRolesToStrings(req.Roles),
	}

	user, err := c.identityUseCase.RegisterUser(ctx, cmd)
	if err != nil {
		return nil, httpx.MapError(ctx, err)
	}

	go runSideEffects(context.WithoutCancel(ctx), "user", user.Email(), "registered", logrus.Fields{})

	resp := toUserFindResponse(user)
	return &resp, nil
}

func (c *userController) UsersReadByEmail(
	ctx context.Context,
	email string,
) (*openapi.UserFindResponse, error) {
	user, err := c.identityUseCase.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, httpx.MapError(ctx, err)
	}
	resp := toUserFindResponse(user)
	return &resp, nil
}

func (c *userController) UsersSyncFromIdp(
	ctx context.Context,
	idpUid string,
) (*openapi.UserFindResponse, error) {
	user, err := c.identityUseCase.SyncUserFromIdp(ctx, idpUid)
	if err != nil {
		return nil, httpx.MapError(ctx, err)
	}

	go runSideEffects(context.WithoutCancel(ctx), "user", idpUid, "synced", logrus.Fields{})

	resp := toUserFindResponse(user)
	return &resp, nil
}

func (c *userController) UsersRead(
	ctx context.Context,
	uid string,
) (*openapi.UserFindResponse, error) {
	user, err := c.identityUseCase.FindUser(ctx, uid)
	if err != nil {
		return nil, httpx.MapError(ctx, err)
	}
	resp := toUserFindResponse(user)
	return &resp, nil
}

func (c *userController) UsersUpdate(
	ctx context.Context,
	uid string,
	req openapi.UserUpdateRequest,
) (*openapi.UserFindResponse, error) {
	cmd := command.UserUpdateCommand{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
	}

	user, err := c.identityUseCase.UpdateUser(ctx, uid, cmd)
	if err != nil {
		return nil, httpx.MapError(ctx, err)
	}

	go runSideEffects(context.WithoutCancel(ctx), "user", uid, "updated", logrus.Fields{})

	resp := toUserFindResponse(user)
	return &resp, nil
}

func (c *userController) UsersDeactivate(
	ctx context.Context,
	uid string,
) (*openapi.UserFindResponse, error) {
	user, err := c.identityUseCase.DeactivateUser(ctx, uid)
	if err != nil {
		return nil, httpx.MapError(ctx, err)
	}

	go runSideEffects(context.WithoutCancel(ctx), "user", uid, "deactivated", logrus.Fields{})

	resp := toUserFindResponseWithStatus(user, openapi.USERSTATUS_DEACTIVATED)
	return &resp, nil
}

func (c *userController) UsersReactivate(
	ctx context.Context,
	uid string,
) (*openapi.UserFindResponse, error) {
	user, err := c.identityUseCase.ReactivateUser(ctx, uid)
	if err != nil {
		return nil, httpx.MapError(ctx, err)
	}

	go runSideEffects(context.WithoutCancel(ctx), "user", uid, "reactivated", logrus.Fields{})

	resp := toUserFindResponseWithStatus(user, openapi.USERSTATUS_ACTIVE)
	return &resp, nil
}

func (c *userController) UsersUpdateRoles(
	ctx context.Context,
	uid string,
	req openapi.UserRolesRequest,
) (*openapi.UserFindResponse, error) {
	cmd := command.UserRolesCommand{Roles: authRoleSliceToStrings(req.Roles)}

	user, err := c.identityUseCase.UpdateUserRoles(ctx, uid, cmd)
	if err != nil {
		return nil, httpx.MapError(ctx, err)
	}

	go runSideEffects(context.WithoutCancel(ctx), "user", uid, "roles_updated", logrus.Fields{"roles": req.Roles})

	resp := toUserFindResponse(user)
	return &resp, nil
}

func (c *userController) UsersSuspend(
	ctx context.Context,
	uid string,
) (*openapi.UserFindResponse, error) {
	user, err := c.identityUseCase.SuspendUser(ctx, uid)
	if err != nil {
		return nil, httpx.MapError(ctx, err)
	}

	go runSideEffects(context.WithoutCancel(ctx), "user", uid, "suspended", logrus.Fields{})

	resp := toUserFindResponseWithStatus(user, openapi.USERSTATUS_SUSPENDED)
	return &resp, nil
}

// ─── mapping helpers ────────────────────────────────────────────────────────

func toUserFindResponse(a *aggregate.UserAggregate) openapi.UserFindResponse {
	return toUserFindResponseWithStatus(a, parseStatusOrDefault(a.Status()))
}

func toUserFindResponseWithStatus(a *aggregate.UserAggregate, status openapi.UserStatus) openapi.UserFindResponse {
	roles := a.Roles()
	out := make([]openapi.AuthRole, len(roles))
	for i, r := range roles {
		out[i] = r.ToAuthRoleRest()
	}
	return openapi.UserFindResponse{
		Uid:         a.Uid(),
		Email:       a.Email(),
		FirstName:   a.FirstName(),
		LastName:    a.LastName(),
		Phone:       a.Phone(),
		Roles:       out,
		Status:      status,
		LastLoginAt: a.LastLoginAt(),
		CreatedAt:   a.CreatedAt(),
		UpdatedAt:   a.UpdatedAt(),
	}
}

func parseStatusOrDefault(s string) openapi.UserStatus {
	if s == "" {
		return openapi.USERSTATUS_ACTIVE
	}
	v, err := openapi.ParseUserStatus(s)
	if err != nil {
		return openapi.USERSTATUS_ACTIVE
	}
	return v
}

func authRolesToStrings(roles *[]openapi.AuthRole) []string {
	if roles == nil {
		return nil
	}
	return authRoleSliceToStrings(*roles)
}

func authRoleSliceToStrings(roles []openapi.AuthRole) []string {
	out := make([]string, len(roles))
	for i, r := range roles {
		out[i] = string(r)
	}
	return out
}
