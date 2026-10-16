package restcontroller

import (
	"context"

	"github.com/codejsha/shared-library-go/pkg/pagination"

	"github.com/codejsha/bookstore-microservices/customer/generated/application/port/openapi"
	"github.com/codejsha/bookstore-microservices/customer/internal/application/usecase"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/aggregate"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/constant"
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
	_ openapi.CustomerApiCredentials,
	email *string,
	name *string,
	phone *string,
	size *int32,
	page *int32,
	sort *string,
) (openapi.CustomersGetAllResponse, error) {
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
	return openapi.CustomersGetAllOk{Body: openapi.CustomerFindAllResponse{Items: items, Total: total}}, nil
}

func (c *customerController) CustomersRead(ctx context.Context, _ openapi.CustomerApiCredentials, uid string) (openapi.CustomersReadResponse, error) {
	customer, err := c.customerUseCase.FindCustomer(ctx, uid)
	if err != nil {
		return nil, httpx.MapGrpcStatus(ctx, httpx.MapNotFound(ctx, err))
	}
	if customer == nil {
		return nil, httpx.MapNotFound(ctx, httpx.ErrNotFound)
	}

	return openapi.CustomersReadOk{Body: toCustomerFindResponse(customer)}, nil
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

func authRolesToRest(roles []constant.AuthRole) []openapi.AuthRole {
	out := make([]openapi.AuthRole, len(roles))
	for i, r := range roles {
		out[i] = r.ToAuthRoleRest()
	}
	return out
}
