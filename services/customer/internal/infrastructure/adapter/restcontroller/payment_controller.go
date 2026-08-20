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

var _ openapi.PaymentApi = (*paymentController)(nil)

type paymentController struct {
	customerUseCase usecase.CustomerUseCase
}

func NewPaymentController(customerUseCase usecase.CustomerUseCase) openapi.PaymentApi {
	return &paymentController{customerUseCase: customerUseCase}
}

func (c *paymentController) CustomerPaymentsGetAll(
	ctx context.Context,
	uid string,
	orderUid *string,
	size *int32,
	page *int32,
	sort *string,
) (*openapi.PaymentFindAllResponse, error) {
	opt := option.NewPaymentQueryOption(
		option.PaymentQueryOption{}.WithUserUid(&uid),
		option.PaymentQueryOption{}.WithOrderUid(orderUid),
		option.PaymentQueryOption{}.WithPage(pagination.NewPageOption(size, page, sort)),
	)

	total, payments, err := c.customerUseCase.FindAllCustomerPayments(ctx, opt)
	if err != nil {
		return nil, httpx.MapGrpcStatus(ctx, err)
	}

	items := make([]openapi.PaymentFindResponse, len(payments))
	for i, payment := range payments {
		items[i] = toPaymentFindResponse(payment)
	}
	return &openapi.PaymentFindAllResponse{Items: items, Total: total}, nil
}

func (c *paymentController) CustomerPaymentsRead(
	ctx context.Context,
	uid string,
	paymentUid string,
) (*openapi.PaymentFindResponse, error) {
	payment, err := c.customerUseCase.FindCustomerPayment(ctx, uid, paymentUid)
	if err != nil {
		return nil, httpx.MapGrpcStatus(ctx, httpx.MapNotFound(ctx, err))
	}
	if payment == nil {
		return nil, httpx.MapNotFound(ctx, httpx.ErrNotFound)
	}

	resp := toPaymentFindResponse(payment)
	return &resp, nil
}

// ─── mapping helpers ────────────────────────────────────────────────────────

func toPaymentFindResponse(a *aggregate.PaymentAggregate) openapi.PaymentFindResponse {
	userUid := a.CustomerUid
	return openapi.PaymentFindResponse{
		Uid:      a.PaymentUid,
		UserUid:  &userUid,
		Currency: &a.Currency,
		Amount:   &a.Amount,
		Status:   &a.Status,
		Reason:   a.ErrorMessage,
	}
}
