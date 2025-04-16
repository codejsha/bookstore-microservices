package service

import (
	"context"

	"github.com/codejsha/bookstore-microservices/customer/internal/application/port/pb/orderpb"
	"github.com/codejsha/bookstore-microservices/customer/internal/application/usecase"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/aggregate"
	"github.com/codejsha/bookstore-microservices/customer/internal/infrastructure/adapter/protostub"
)

var _ usecase.OrderUseCase = (*orderService)(nil)

func NewOrderService(
	orderGrpcClient *protostub.OrderGrpcClient,
) usecase.OrderUseCase {
	return &orderService{
		orderServiceClient: orderGrpcClient.Client,
	}
}

type orderService struct {
	orderServiceClient orderpb.OrderServiceClient
}

func (s orderService) FindAllCustomerOrders(ctx context.Context, req *orderpb.OrderFindAllProtoReq) (int64, []*aggregate.OrderAggregate, error) {
	// find orders
	resp, err := s.orderServiceClient.FindAllOrders(ctx, req)
	if err != nil {
		return 0, nil, err
	}

	// return response
	orderAggs := make([]*aggregate.OrderAggregate, 0, len(resp.Items))
	for i, item := range resp.Items {
		orderAggs[i] = aggregate.NewOrderAggregate(item)
	}
	return resp.Total, orderAggs, nil
}

func (s orderService) FindCustomerOrder(ctx context.Context, req *orderpb.OrderFindProtoReq) (*aggregate.OrderAggregate, error) {
	// find order
	resp, err := s.orderServiceClient.FindOrder(ctx, req)
	if err != nil {
		return nil, err
	}
	return aggregate.NewOrderAggregate(resp), nil
}
