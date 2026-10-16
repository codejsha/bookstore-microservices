package option

import (
	"github.com/codejsha/shared-library-go/pkg/pagination"

	"github.com/codejsha/bookstore-microservices/customer/generated/application/port/openapi"
)

type OrderQueryOption struct {
	userUid *string
	status  *string
	page    pagination.PageOption
}

type OrderQueryOptionFunc func(*OrderQueryOption)

func NewOrderQueryOption(opts ...OrderQueryOptionFunc) OrderQueryOption {
	var o OrderQueryOption
	for _, opt := range opts {
		opt(&o)
	}
	return o
}

func (o OrderQueryOption) UserUid() *string {
	return o.userUid
}

func (o OrderQueryOption) Status() *string {
	return o.status
}

func (o OrderQueryOption) Page() pagination.PageOption {
	return o.page
}

func (OrderQueryOption) WithUserUid(userUid *string) OrderQueryOptionFunc {
	return func(o *OrderQueryOption) {
		o.userUid = userUid
	}
}

func (OrderQueryOption) WithStatus(status *string) OrderQueryOptionFunc {
	return func(o *OrderQueryOption) {
		o.status = status
	}
}

func (OrderQueryOption) WithStatusRest(status *openapi.OrderStatus) OrderQueryOptionFunc {
	return func(o *OrderQueryOption) {
		var orderStatus string
		if status != nil {
			switch *status {
			case openapi.OrderStatusUnknown:
				orderStatus = "UNKNOWN"
			case openapi.OrderStatusPending:
				orderStatus = "PENDING"
			case openapi.OrderStatusPaid:
				orderStatus = "PAID"
			case openapi.OrderStatusShipping:
				orderStatus = "SHIPPING"
			case openapi.OrderStatusCompleted:
				orderStatus = "COMPLETED"
			case openapi.OrderStatusCancelled:
				orderStatus = "CANCELLED"
			default:
				orderStatus = "UNKNOWN"
			}
		}
		o.status = &orderStatus
	}
}

func (OrderQueryOption) WithPage(page pagination.PageOption) OrderQueryOptionFunc {
	return func(o *OrderQueryOption) {
		o.page = page
	}
}
