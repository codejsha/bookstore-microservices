package aggregate

import (
	"github.com/codejsha/bookstore-microservices/customer/internal/application/port/openapi"
	"github.com/codejsha/bookstore-microservices/customer/internal/application/port/pb/orderpb"
)

func NewOrderAggregate(
	resp *orderpb.OrderFindProtoResp,
) *OrderAggregate {
	items := make([]openapi.OrderItem, len(resp.OrderItems))
	for i, item := range resp.OrderItems {
		orderItem := openapi.OrderItem{
			BookId:   item.BookId,
			Quantity: item.Quantity,
		}
		items[i] = orderItem
	}
	agg := &OrderAggregate{
		Id:         resp.Id,
		UserId:     resp.UserId,
		OrderItems: items,
		TotalPrice: resp.TotalPrice,
		Status:     openapi.OrderStatus(resp.Status),
	}
	return agg
}

type OrderAggregate struct {
	// order id
	Id         int64
	UserId     string
	OrderItems []openapi.OrderItem
	TotalPrice float64
	Status     openapi.OrderStatus
}

func (a *OrderAggregate) ToOrderFindResponse() openapi.OrderFindWebResp {
	resp := openapi.OrderFindWebResp{
		Id:         a.Id,
		UserId:     a.UserId,
		OrderItems: a.OrderItems,
		TotalPrice: a.TotalPrice,
		Status:     a.Status,
	}
	return resp
}
