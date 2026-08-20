package aggregate

import "time"

type WarehouseAggregate struct {
	Uid       string
	Name      string
	Address   *string
	Capacity  int32
	CreatedAt time.Time
	UpdatedAt *time.Time
}
