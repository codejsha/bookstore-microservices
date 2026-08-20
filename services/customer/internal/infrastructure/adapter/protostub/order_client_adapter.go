package protostub

import (
	"context"

	"github.com/codejsha/bookstore-microservices/customer/generated/application/port/pb/orderpb"
	port "github.com/codejsha/bookstore-microservices/customer/internal/application/port/protostub"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/aggregate"
)

var _ port.OrderClient = (*orderClient)(nil)

type orderClient struct {
	client orderpb.OrderServiceClient
}

func NewOrderClient(c *OrderGrpcClient) port.OrderClient {
	return &orderClient{client: c.Client}
}

func (g *orderClient) ListOrders(ctx context.Context, userUid string, status *string, pageSize int32) (int64, []*aggregate.OrderAggregate, error) {
	ctx = withActor(ctx)
	resp, err := g.client.ListOrders(ctx, &orderpb.ListOrdersRequest{
		UserUid:  userUid,
		Status:   stringToOrderStatusProto(status),
		PageSize: pageSize,
	})
	if err != nil {
		return 0, nil, err
	}
	out := make([]*aggregate.OrderAggregate, len(resp.GetOrders()))
	for i, o := range resp.GetOrders() {
		out[i] = toOrderAggregate(o)
	}
	return int64(resp.GetTotalSize()), out, nil
}

func (g *orderClient) FindOrder(ctx context.Context, uid string) (*aggregate.OrderAggregate, error) {
	ctx = withActor(ctx)
	resp, err := g.client.FindOrder(ctx, &orderpb.FindOrderRequest{Uid: uid})
	if err != nil {
		return nil, err
	}
	o := resp.GetOrder()
	if o == nil {
		return nil, nil
	}
	return toOrderAggregate(o), nil
}

func toOrderAggregate(o *orderpb.Order) *aggregate.OrderAggregate {
	lineItems := o.GetLineItems()
	items := make([]aggregate.OrderItem, len(lineItems))
	for i, item := range lineItems {
		items[i] = aggregate.OrderItem{
			BookUid:  item.GetProductUid(),
			Quantity: item.GetQuantity(),
		}
	}
	return &aggregate.OrderAggregate{
		Uid:         o.GetUid(),
		UserUid:     o.GetUserUid(),
		OrderItems:  items,
		TotalAmount: o.GetTotalAmount(),
		Status:      toOrderStatus(o.GetStatus()),
	}
}

func toOrderStatus(s orderpb.OrderStatus) aggregate.OrderStatus {
	switch s {
	case orderpb.OrderStatus_ORDER_STATUS_PENDING:
		return aggregate.ORDERSTATUS_PENDING
	case orderpb.OrderStatus_ORDER_STATUS_CONFIRMED, orderpb.OrderStatus_ORDER_STATUS_PROCESSING:
		return aggregate.ORDERSTATUS_PAID
	case orderpb.OrderStatus_ORDER_STATUS_SHIPPED:
		return aggregate.ORDERSTATUS_SHIPPING
	case orderpb.OrderStatus_ORDER_STATUS_DELIVERED:
		return aggregate.ORDERSTATUS_COMPLETED
	case orderpb.OrderStatus_ORDER_STATUS_CANCELLED, orderpb.OrderStatus_ORDER_STATUS_REFUNDED:
		return aggregate.ORDERSTATUS_CANCELLED
	default:
		return aggregate.ORDERSTATUS_UNKNOWN
	}
}

func stringToOrderStatusProto(s *string) orderpb.OrderStatus {
	if s == nil {
		return orderpb.OrderStatus_ORDER_STATUS_UNSPECIFIED
	}
	switch *s {
	case "PENDING":
		return orderpb.OrderStatus_ORDER_STATUS_PENDING
	case "PAID":
		return orderpb.OrderStatus_ORDER_STATUS_CONFIRMED
	case "SHIPPING":
		return orderpb.OrderStatus_ORDER_STATUS_SHIPPED
	case "COMPLETED":
		return orderpb.OrderStatus_ORDER_STATUS_DELIVERED
	case "CANCELLED":
		return orderpb.OrderStatus_ORDER_STATUS_CANCELLED
	default:
		return orderpb.OrderStatus_ORDER_STATUS_UNSPECIFIED
	}
}
