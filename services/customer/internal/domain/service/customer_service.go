package service

import (
	"context"

	"github.com/codejsha/bookstore-microservices/customer/internal/application/port/pb/userpb"
	"github.com/codejsha/bookstore-microservices/customer/internal/application/usecase"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/aggregate"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/model/mapper"
	"github.com/codejsha/bookstore-microservices/customer/internal/infrastructure/adapter/protostub"
)

var _ usecase.CustomerUseCase = (*customerService)(nil)

func NewCustomerService(
	userGrpcClient *protostub.UserGrpcClient,
	userMapper *mapper.UserMapper,
) usecase.CustomerUseCase {
	return &customerService{
		userServiceClient: userGrpcClient.Client,
		userMapper:        userMapper,
	}
}

type customerService struct {
	userServiceClient userpb.UserServiceClient
	userMapper        *mapper.UserMapper
}

func (s customerService) FindAllCustomers(ctx context.Context, req *userpb.UserFindAllProtoReq) (int64, []*aggregate.CustomerAggregate, error) {
	// find customers
	resp, err := s.userServiceClient.FindAllUsers(ctx, req)
	if err != nil {
		return 0, nil, err
	}

	// return response
	customerAggs := make([]*aggregate.CustomerAggregate, 0, len(resp.Items))
	for i, item := range resp.Items {
		customerAggs[i] = aggregate.NewCustomerAggregate(item, s.userMapper)
	}
	return resp.Total, customerAggs, nil
}

func (s customerService) FindCustomer(ctx context.Context, req *userpb.UserFindProtoReq) (*aggregate.CustomerAggregate, error) {
	// find customer
	resp, err := s.userServiceClient.FindUser(ctx, req)
	if err != nil {
		return nil, err
	}
	agg := aggregate.NewCustomerAggregate(resp, s.userMapper)
	return agg, nil
}

func (s customerService) UpdateCustomer(ctx context.Context, req *userpb.UserUpdateProtoReq) (*aggregate.CustomerAggregate, error) {
	// update customer
	resp, err := s.userServiceClient.UpdateUser(ctx, req)
	if err != nil {
		return nil, err
	}
	agg := aggregate.NewCustomerAggregate(resp, s.userMapper)
	return agg, nil
}

func (s customerService) DeleteCustomer(ctx context.Context, id string) error {
	// delete customer
	_, err := s.userServiceClient.DeleteUser(ctx, &userpb.UserDeleteProtoReq{Id: id})
	if err != nil {
		return err
	}
	return nil
}
