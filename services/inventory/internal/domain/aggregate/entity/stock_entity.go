package entity

import (
	"fmt"
	"time"

	"github.com/codejsha/bookstore-microservices/inventory/internal/application/port/openapi"
)

type StockEntity struct {
	Id               int64     `gorm:"primaryKey;autoIncrement"`
	BookId           int64     `gorm:"not null"`
	WarehouseId      int64     `gorm:"not null"`
	Quantity         int64     `gorm:"not null"`
	ReservedQuantity int64     `gorm:"not null"`
	CreatedAt        time.Time `gorm:"autoCreateTime"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime"`
}

func (e StockEntity) TableName() string {
	return "stock"
}

func (e StockEntity) ToStockFindWebResp() openapi.StockFindWebResp {
	resp := openapi.StockFindWebResp{
		Id:               e.Id,
		BookId:           e.BookId,
		WarehouseId:      e.WarehouseId,
		Quantity:         e.Quantity,
		ReservedQuantity: e.ReservedQuantity,
	}
	return resp
}

func (e StockEntity) FromValue(value interface{}, params interface{}) (StockEntity, error) {
	switch value.(type) {
	case openapi.StockCreateWebReq:
		v := value.(openapi.StockCreateWebReq)
		ent := StockEntity{
			BookId:           v.BookId,
			WarehouseId:      v.WarehouseId,
			Quantity:         v.Quantity,
			ReservedQuantity: v.ReservedQuantity,
		}
		return ent, nil
	case openapi.StockUpdateWebReq:
		v := value.(openapi.StockUpdateWebReq)
		p := params.(openapi.StockUpdateWebParam)
		ent := StockEntity{
			Id:               p.Id,
			BookId:           v.BookId,
			WarehouseId:      v.WarehouseId,
			Quantity:         v.Quantity,
			ReservedQuantity: v.ReservedQuantity,
		}
		return ent, nil
	default:
		return StockEntity{}, fmt.Errorf("failed to convert request to StockEntity: %v", value)
	}
}
