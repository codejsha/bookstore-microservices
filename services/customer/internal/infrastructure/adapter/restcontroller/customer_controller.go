package restcontroller

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/codejsha/shared-library-go/pkg/pagination"

	"github.com/codejsha/bookstore-microservices/customer/generated/application/port/openapi"
	"github.com/codejsha/bookstore-microservices/customer/internal/application/usecase"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/aggregate"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/constant"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/model/command"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/model/option"
	"github.com/codejsha/bookstore-microservices/customer/internal/infrastructure/httpx"
)

var _ openapi.CustomerApi = (*customerController)(nil)

type customerController struct {
	customerUseCase usecase.CustomerUseCase
}

func NewCustomerController(customerUseCase usecase.CustomerUseCase) openapi.CustomerApi {
	return &customerController{customerUseCase: customerUseCase}
}

func (c *customerController) CustomersGetAll(
	ctx context.Context,
	email *string,
	name *string,
	phone *string,
	size *int32,
	page *int32,
	sort *string,
) (*openapi.CustomerFindAllResponse, error) {
	opt := option.NewUserQueryOption(
		option.UserQueryOption{}.WithEmail(email),
		option.UserQueryOption{}.WithName(name),
		option.UserQueryOption{}.WithPhone(phone),
		option.UserQueryOption{}.WithPage(pagination.NewPageOption(size, page, sort)),
	)

	total, customers, err := c.customerUseCase.FindAllCustomers(ctx, opt)
	if err != nil {
		return nil, httpx.MapGrpcStatus(ctx, err)
	}

	items := make([]openapi.CustomerFindResponse, len(customers))
	for i, customer := range customers {
		items[i] = toCustomerFindResponse(customer)
	}
	return &openapi.CustomerFindAllResponse{Items: items, Total: total}, nil
}

func (c *customerController) CustomersRead(ctx context.Context, uid string) (*openapi.CustomerFindResponse, error) {
	customer, err := c.customerUseCase.FindCustomer(ctx, uid)
	if err != nil {
		return nil, httpx.MapGrpcStatus(ctx, httpx.MapNotFound(ctx, err))
	}
	if customer == nil {
		return nil, httpx.MapNotFound(ctx, httpx.ErrNotFound)
	}

	resp := toCustomerFindResponse(customer)
	return &resp, nil
}

func (c *customerController) CustomersUpdate(
	ctx context.Context,
	uid string,
	req openapi.CustomerUpdateRequest,
) (*openapi.CustomerUpdateResponse, error) {
	cmd := command.CustomerUpdateCommand{
		Uid:       uid,
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		Roles:     toAuthRoles(req.Roles),
	}

	customer, err := c.customerUseCase.UpdateCustomer(ctx, cmd)
	if err != nil {
		return nil, httpx.MapBusinessError(ctx, httpx.MapNotImplemented(ctx, err))
	}

	dispatchSideEffects(context.WithoutCancel(ctx), "customer", uid, "updated", logrus.Fields{})

	resp := toCustomerUpdateResponse(customer)
	return &resp, nil
}

func (c *customerController) CustomersDelete(ctx context.Context, uid string) error {
	if err := c.customerUseCase.DeleteCustomer(ctx, uid); err != nil {
		return httpx.MapNotImplemented(ctx, err)
	}

	dispatchSideEffects(context.WithoutCancel(ctx), "customer", uid, "deleted", logrus.Fields{})
	return nil
}

// ─── mapping helpers ────────────────────────────────────────────────────────

func toCustomerFindResponse(a *aggregate.CustomerAggregate) openapi.CustomerFindResponse {
	roles := authRolesToRest(a.Roles)
	return openapi.CustomerFindResponse{
		Uid:       a.Uid,
		Email:     &a.Email,
		FirstName: &a.FirstName,
		LastName:  &a.LastName,
		Phone:     a.Phone,
		Roles:     &roles,
	}
}

func toCustomerUpdateResponse(a *aggregate.CustomerAggregate) openapi.CustomerUpdateResponse {
	roles := authRolesToRest(a.Roles)
	return openapi.CustomerUpdateResponse{
		Uid:       a.Uid,
		Email:     &a.Email,
		FirstName: &a.FirstName,
		LastName:  &a.LastName,
		Phone:     a.Phone,
		Roles:     &roles,
	}
}

func toAuthRoles(roles *[]openapi.AuthRole) []constant.AuthRole {
	if roles == nil {
		return nil
	}
	out := make([]constant.AuthRole, len(*roles))
	for i, r := range *roles {
		out[i] = constant.AuthRoleFromString(string(r))
	}
	return out
}

func authRolesToRest(roles []constant.AuthRole) []openapi.AuthRole {
	out := make([]openapi.AuthRole, len(roles))
	for i, r := range roles {
		out[i] = r.ToAuthRoleRest()
	}
	return out
}
