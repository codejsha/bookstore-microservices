package service

import (
	"context"

	"github.com/codejsha/bookstore-microservices/customer/internal/application/port/pb/paymentpb"
	"github.com/codejsha/bookstore-microservices/customer/internal/application/usecase"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/aggregate"
	"github.com/codejsha/bookstore-microservices/customer/internal/infrastructure/adapter/protostub"
)

var _ usecase.PaymentUseCase = (*paymentService)(nil)

func NewPaymentService(
	paymentGrpcClient *protostub.PaymentGrpcClient,
) usecase.PaymentUseCase {
	return &paymentService{
		paymentServiceClient: paymentGrpcClient.Client,
	}
}

type paymentService struct {
	paymentServiceClient paymentpb.PaymentServiceClient
}

func (s paymentService) FindAllCustomerPayments(ctx context.Context, req *paymentpb.PaymentFindAllProtoReq) (int64, []*aggregate.PaymentAggregate, error) {
	// find payments
	resp, err := s.paymentServiceClient.FindAllPayments(ctx, req)
	if err != nil {
		return 0, nil, err
	}

	// return response
	paymentAggs := make([]*aggregate.PaymentAggregate, 0, len(resp.Items))
	for i, item := range resp.Items {
		paymentAggs[i] = aggregate.NewPaymentAggregate(item)
	}
	return resp.Total, paymentAggs, nil
}

func (s paymentService) FindCustomerPayment(ctx context.Context, req *paymentpb.PaymentFindProtoReq) (*aggregate.PaymentAggregate, error) {
	// find payment
	resp, err := s.paymentServiceClient.FindPayment(ctx, req)
	if err != nil {
		return nil, err
	}
	return aggregate.NewPaymentAggregate(resp), nil
}
