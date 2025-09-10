package command

type WarehouseCreateCommand struct {
	Name     string
	Address  *string
	Capacity int32
}

type WarehouseUpdateCommand struct {
	Name     *string
	Address  *string
	Capacity *int32
}
