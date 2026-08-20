package option

import (
	"github.com/codejsha/shared-library-go/pkg/pagination"
)

type PaymentQueryOption struct {
	userUid  *string
	orderUid *string
	page     pagination.PageOption
}

type PaymentQueryOptionFunc func(*PaymentQueryOption)

func NewPaymentQueryOption(opts ...PaymentQueryOptionFunc) PaymentQueryOption {
	var o PaymentQueryOption
	for _, opt := range opts {
		opt(&o)
	}
	return o
}

func (o PaymentQueryOption) UserUid() *string {
	return o.userUid
}

func (o PaymentQueryOption) OrderUid() *string {
	return o.orderUid
}

func (o PaymentQueryOption) Page() pagination.PageOption {
	return o.page
}

func (PaymentQueryOption) WithUserUid(userUid *string) PaymentQueryOptionFunc {
	return func(o *PaymentQueryOption) {
		o.userUid = userUid
	}
}

func (PaymentQueryOption) WithOrderUid(orderUid *string) PaymentQueryOptionFunc {
	return func(o *PaymentQueryOption) {
		o.orderUid = orderUid
	}
}

func (PaymentQueryOption) WithPage(page pagination.PageOption) PaymentQueryOptionFunc {
	return func(o *PaymentQueryOption) {
		o.page = page
	}
}
