package aggregate

type OrderAggregate struct {
	Uid         string
	UserUid     string
	OrderItems  []OrderItem
	TotalAmount float64
	Status      OrderStatus
}

type OrderItem struct {
	BookUid  string
	Quantity int32
}

type OrderStatus int32

const (
	ORDERSTATUS_UNKNOWN   OrderStatus = 0
	ORDERSTATUS_PENDING   OrderStatus = 1
	ORDERSTATUS_PAID      OrderStatus = 2
	ORDERSTATUS_SHIPPING  OrderStatus = 3
	ORDERSTATUS_COMPLETED OrderStatus = 4
	ORDERSTATUS_CANCELLED OrderStatus = 5
)
