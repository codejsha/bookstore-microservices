package usecase

import (
	"context"

	"github.com/codejsha/bookstore-microservices/customer/internal/application/port/pb/orderpb"
	"github.com/codejsha/bookstore-microservices/customer/internal/application/port/pb/paymentpb"
	"github.com/codejsha/bookstore-microservices/customer/internal/application/port/pb/userpb"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/aggregate"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/aggregate/entity"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/model/condition"
)

type CustomerUseCase interface {
	FindAllCustomers(ctx context.Context, req *userpb.UserFindAllProtoReq) (int64, []*aggregate.CustomerAggregate, error)
	FindCustomer(ctx context.Context, req *userpb.UserFindProtoReq) (*aggregate.CustomerAggregate, error)
	UpdateCustomer(ctx context.Context, req *userpb.UserUpdateProtoReq) (*aggregate.CustomerAggregate, error)
	DeleteCustomer(ctx context.Context, id string) error
}

type OrderUseCase interface {
	FindAllCustomerOrders(ctx context.Context, req *orderpb.OrderFindAllProtoReq) (int64, []*aggregate.OrderAggregate, error)
	FindCustomerOrder(ctx context.Context, req *orderpb.OrderFindProtoReq) (*aggregate.OrderAggregate, error)
}

type PaymentUseCase interface {
	FindAllCustomerPayments(ctx context.Context, req *paymentpb.PaymentFindAllProtoReq) (int64, []*aggregate.PaymentAggregate, error)
	FindCustomerPayment(ctx context.Context, req *paymentpb.PaymentFindProtoReq) (*aggregate.PaymentAggregate, error)
}

type PointUseCase interface {
	FindPoint(ctx context.Context, cond condition.PointCondition) (entity.PointEntity, error)
	UpdatePoint(ctx context.Context, ent entity.PointEntity) (entity.PointEntity, error)
}

type WishlistUseCase interface {
	FindWishlist(ctx context.Context, cond condition.WishlistCondition) (*aggregate.WishlistAggregate, error)
	UpdateWishlist(ctx context.Context, ents []entity.WishlistEntity) (*aggregate.WishlistAggregate, error)
	DeleteWishlist(ctx context.Context, ids []int64) error
}
