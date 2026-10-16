package restcontroller

import (
	"context"

	"github.com/codejsha/shared-library-go/pkg/pagination"

	"github.com/codejsha/bookstore-microservices/customer/generated/application/port/openapi"
	"github.com/codejsha/bookstore-microservices/customer/internal/application/usecase"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/aggregate"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/model/option"
	"github.com/codejsha/bookstore-microservices/customer/internal/infrastructure/httpx"
)

var _ openapi.OrderApi = (*orderController)(nil)

type orderController struct {
	customerUseCase usecase.CustomerUseCase
}

func NewOrderController(customerUseCase usecase.CustomerUseCase) openapi.OrderApi {
	return &orderController{customerUseCase: customerUseCase}
}

func (c *orderController) CustomerOrdersGetAll(
	ctx context.Context,
	_ openapi.OrderApiCredentials,
	uid string,
	status *openapi.OrderStatus,
	size *int32,
	page *int32,
	sort *string,
) (openapi.CustomerOrdersGetAllResponse, error) {
	opt := option.NewOrderQueryOption(
		option.OrderQueryOption{}.WithUserUid(&uid),
		option.OrderQueryOption{}.WithStatusRest(status),
		option.OrderQueryOption{}.WithPage(pagination.NewPageOption(size, page, sort)),
	)

	total, orders, err := c.customerUseCase.FindAllCustomerOrders(ctx, opt)
	if err != nil {
		return nil, httpx.MapGrpcStatus(ctx, err)
	}

	items := make([]openapi.OrderFindResponse, len(orders))
	for i, order := range orders {
		items[i] = toOrderFindResponse(order)
	}
	return openapi.CustomerOrdersGetAllOk{Body: openapi.OrderFindAllResponse{Items: items, Total: total}}, nil
}

func (c *orderController) CustomerOrdersRead(
	ctx context.Context,
	_ openapi.OrderApiCredentials,
	uid string,
	orderUid string,
) (openapi.CustomerOrdersReadResponse, error) {
	order, err := c.customerUseCase.FindCustomerOrder(ctx, uid, orderUid)
	if err != nil {
		return nil, httpx.MapGrpcStatus(ctx, httpx.MapNotFound(ctx, err))
	}
	if order == nil {
		return nil, httpx.MapNotFound(ctx, httpx.ErrNotFound)
	}

	return openapi.CustomerOrdersReadOk{Body: toOrderFindResponse(order)}, nil
}

// ─── mapping helpers ────────────────────────────────────────────────────────

func toOrderFindResponse(a *aggregate.OrderAggregate) openapi.OrderFindResponse {
	out := make([]openapi.OrderItem, len(a.OrderItems))
	for i, item := range a.OrderItems {
		out[i] = openapi.OrderItem{BookUid: item.BookUid, Quantity: item.Quantity}
	}
	status := toOrderStatusRest(a.Status)
	userUid := a.UserUid
	return openapi.OrderFindResponse{
		Uid:        a.Uid,
		UserUid:    &userUid,
		OrderItems: &out,
		TotalPrice: &a.TotalAmount,
		Status:     &status,
	}
}

func toOrderStatusRest(s aggregate.OrderStatus) openapi.OrderStatus {
	switch s {
	case aggregate.ORDERSTATUS_PENDING:
		return openapi.OrderStatusPending
	case aggregate.ORDERSTATUS_PAID:
		return openapi.OrderStatusPaid
	case aggregate.ORDERSTATUS_SHIPPING:
		return openapi.OrderStatusShipping
	case aggregate.ORDERSTATUS_COMPLETED:
		return openapi.OrderStatusCompleted
	case aggregate.ORDERSTATUS_CANCELLED:
		return openapi.OrderStatusCancelled
	default:
		return openapi.OrderStatusUnknown
	}
}
