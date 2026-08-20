package protostub

import (
	"context"

	"github.com/codejsha/bookstore-microservices/customer/internal/domain/aggregate"
)

type UserClient interface {
	ListUsers(ctx context.Context, email, name, phone string, pageSize int32) (int64, []*aggregate.CustomerAggregate, error)
	FindUser(ctx context.Context, uid string) (*aggregate.CustomerAggregate, error)
}

type OrderClient interface {
	ListOrders(ctx context.Context, userUid string, status *string, pageSize int32) (int64, []*aggregate.OrderAggregate, error)
	FindOrder(ctx context.Context, uid string) (*aggregate.OrderAggregate, error)
}

type PaymentClient interface {
	ListPayments(ctx context.Context, customerUid string, pageSize int32) (int64, []*aggregate.PaymentAggregate, error)
	FindPayment(ctx context.Context, uid string) (*aggregate.PaymentAggregate, error)
}

type DeliveryClient interface {
	TrackShipment(ctx context.Context, uid string, callerUid string) (*aggregate.ShipmentAggregate, error)
	ListShipments(ctx context.Context, orderUid string, status *string, pageSize *int32, callerUid string) (int64, []*aggregate.ShipmentAggregate, error)
}
