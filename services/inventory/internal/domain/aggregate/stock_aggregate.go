package aggregate

import (
	"github.com/codejsha/bookstore-microservices/inventory/internal/application/port/openapi"
	"github.com/codejsha/bookstore-microservices/inventory/internal/application/port/pb/stockpb"
	"github.com/codejsha/bookstore-microservices/inventory/internal/domain/aggregate/entity"
)

func NewStockAggregate(
	stock entity.StockEntity,
	warehouse entity.WarehouseEntity,
) *StockAggregate {
	return &StockAggregate{
		Id:        stock.Id,
		Stock:     stock,
		Warehouse: warehouse,
	}
}

type StockAggregate struct {
	// stock id
	Id        int64
	Stock     entity.StockEntity
	Warehouse entity.WarehouseEntity
}

func (a *StockAggregate) ToStockFindWebResp() openapi.StockFindWebResp {
	resp := openapi.StockFindWebResp{
		Id:               a.Id,
		BookId:           a.Stock.BookId,
		WarehouseId:      a.Warehouse.Id,
		Quantity:         a.Stock.Quantity,
		ReservedQuantity: a.Stock.ReservedQuantity,
	}
	return resp
}

func (a *StockAggregate) ToStockFindProtoResp() *stockpb.StockFindProtoResp {
	resp := &stockpb.StockFindProtoResp{
		Id:               a.Id,
		BookId:           a.Stock.BookId,
		WarehouseId:      a.Warehouse.Id,
		Quantity:         a.Stock.Quantity,
		ReservedQuantity: a.Stock.ReservedQuantity,
	}
	return resp
}
