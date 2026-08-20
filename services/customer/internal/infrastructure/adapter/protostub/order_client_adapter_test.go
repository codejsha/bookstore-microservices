package protostub

import (
	"reflect"
	"testing"

	"github.com/codejsha/bookstore-microservices/customer/generated/application/port/pb/orderpb"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/aggregate"
)

func TestToOrderAggregate(t *testing.T) {
	got := toOrderAggregate(&orderpb.Order{
		Uid:         "o-1",
		UserUid:     "u-1",
		TotalAmount: 12.5,
		Status:      orderpb.OrderStatus_ORDER_STATUS_SHIPPED,
		LineItems: []*orderpb.OrderLineItem{
			{ProductUid: "b-1", Quantity: 2},
			{ProductUid: "b-2", Quantity: 5},
		},
	})
	if got.Uid != "o-1" || got.UserUid != "u-1" || got.TotalAmount != 12.5 {
		t.Errorf("scalar fields = %+v", got)
	}
	if got.Status != aggregate.ORDERSTATUS_SHIPPING {
		t.Errorf("status = %v, want SHIPPING", got.Status)
	}
	want := []aggregate.OrderItem{{BookUid: "b-1", Quantity: 2}, {BookUid: "b-2", Quantity: 5}}
	if !reflect.DeepEqual(got.OrderItems, want) {
		t.Errorf("items = %+v, want %+v", got.OrderItems, want)
	}
}

func TestToOrderStatus(t *testing.T) {
	cases := []struct {
		in   orderpb.OrderStatus
		want aggregate.OrderStatus
	}{
		{orderpb.OrderStatus_ORDER_STATUS_PENDING, aggregate.ORDERSTATUS_PENDING},
		{orderpb.OrderStatus_ORDER_STATUS_CONFIRMED, aggregate.ORDERSTATUS_PAID},
		{orderpb.OrderStatus_ORDER_STATUS_PROCESSING, aggregate.ORDERSTATUS_PAID},
		{orderpb.OrderStatus_ORDER_STATUS_SHIPPED, aggregate.ORDERSTATUS_SHIPPING},
		{orderpb.OrderStatus_ORDER_STATUS_DELIVERED, aggregate.ORDERSTATUS_COMPLETED},
		{orderpb.OrderStatus_ORDER_STATUS_CANCELLED, aggregate.ORDERSTATUS_CANCELLED},
		{orderpb.OrderStatus_ORDER_STATUS_REFUNDED, aggregate.ORDERSTATUS_CANCELLED},
		{orderpb.OrderStatus_ORDER_STATUS_UNSPECIFIED, aggregate.ORDERSTATUS_UNKNOWN},
	}
	for _, c := range cases {
		if got := toOrderStatus(c.in); got != c.want {
			t.Errorf("toOrderStatus(%v) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestStringToOrderStatusProto(t *testing.T) {
	cases := []struct {
		in   *string
		want orderpb.OrderStatus
	}{
		{nil, orderpb.OrderStatus_ORDER_STATUS_UNSPECIFIED},
		{ptrStr("PENDING"), orderpb.OrderStatus_ORDER_STATUS_PENDING},
		{ptrStr("PAID"), orderpb.OrderStatus_ORDER_STATUS_CONFIRMED},
		{ptrStr("SHIPPING"), orderpb.OrderStatus_ORDER_STATUS_SHIPPED},
		{ptrStr("COMPLETED"), orderpb.OrderStatus_ORDER_STATUS_DELIVERED},
		{ptrStr("CANCELLED"), orderpb.OrderStatus_ORDER_STATUS_CANCELLED},
		{ptrStr("garbage"), orderpb.OrderStatus_ORDER_STATUS_UNSPECIFIED},
	}
	for _, c := range cases {
		if got := stringToOrderStatusProto(c.in); got != c.want {
			t.Errorf("stringToOrderStatusProto(%v) = %v, want %v", c.in, got, c.want)
		}
	}
}

func ptrStr(s string) *string { return &s }
