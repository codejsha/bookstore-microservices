package protostub

import (
	"context"

	"github.com/codejsha/bookstore-microservices/customer/generated/application/port/pb/paymentpb"
	port "github.com/codejsha/bookstore-microservices/customer/internal/application/port/protostub"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/aggregate"
)

var _ port.PaymentClient = (*paymentClient)(nil)

type paymentClient struct {
	client paymentpb.PaymentServiceClient
}

func NewPaymentClient(c *PaymentGrpcClient) port.PaymentClient {
	return &paymentClient{client: c.Client}
}

func (g *paymentClient) ListPayments(ctx context.Context, customerUid string, pageSize int32) (int64, []*aggregate.PaymentAggregate, error) {
	ctx = withActor(ctx)
	resp, err := g.client.ListPayments(ctx, &paymentpb.ListPaymentsRequest{
		CustomerUid: customerUid,
		PageSize:    pageSize,
	})
	if err != nil {
		return 0, nil, err
	}
	out := make([]*aggregate.PaymentAggregate, len(resp.GetPayments()))
	for i, p := range resp.GetPayments() {
		out[i] = toPaymentAggregate(p)
	}
	return int64(resp.GetTotalSize()), out, nil
}

func (g *paymentClient) FindPayment(ctx context.Context, uid string) (*aggregate.PaymentAggregate, error) {
	ctx = withActor(ctx)
	resp, err := g.client.FindPayment(ctx, &paymentpb.FindPaymentRequest{Uid: uid})
	if err != nil {
		return nil, err
	}
	p := resp.GetPayment()
	if p == nil {
		return nil, nil
	}
	return toPaymentAggregate(p), nil
}

func toPaymentAggregate(p *paymentpb.Payment) *aggregate.PaymentAggregate {
	return &aggregate.PaymentAggregate{
		Uid:           p.GetUid(),
		PaymentUid:    p.GetPaymentUid(),
		CustomerUid:   p.GetCustomerUid(),
		Currency:      p.GetCurrency(),
		Amount:        float64(p.GetAmount()),
		Status:        p.GetStatus().String(),
		PaymentMethod: optionalString(p.GetPaymentMethod()),
		ErrorCode:     optionalString(p.GetErrorCode()),
		ErrorMessage:  optionalString(p.GetErrorMessage()),
	}
}

func optionalString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
