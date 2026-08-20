package constant

type OrderStatus int32

const (
	OrderStatus_UNKNOWN   OrderStatus = 0
	OrderStatus_PENDING   OrderStatus = 1
	OrderStatus_PAID      OrderStatus = 2
	OrderStatus_SHIPPING  OrderStatus = 3
	OrderStatus_COMPLETED OrderStatus = 4
	OrderStatus_CANCELLED OrderStatus = 5
)
