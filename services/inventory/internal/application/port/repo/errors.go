package repo

import "errors"

var (
	ErrInvalidQuantity = errors.New("invalid quantity")

	ErrInsufficientStock = errors.New("insufficient stock")

	ErrSameWarehouse = errors.New("source and target warehouse must differ")

	ErrStockNotFound = errors.New("stock not found")

	ErrWarehouseNotFound = errors.New("warehouse not found")

	ErrDuplicateClosing = errors.New("monthly closing already exists for this warehouse and period")
)
