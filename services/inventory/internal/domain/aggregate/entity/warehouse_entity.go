package entity

import (
	"fmt"
	"time"

	"github.com/codejsha/bookstore-microservices/inventory/internal/application/port/openapi"
)

type WarehouseEntity struct {
	Id        int64     `gorm:"primaryKey;autoIncrement"`
	Name      string    `gorm:"not null"`
	Address   string    `gorm:"null"`
	Capacity  int64     `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (e WarehouseEntity) TableName() string {
	return "warehouse"
}

func (e WarehouseEntity) ToWarehouseFindResponse() openapi.WarehouseFindWebResp {
	resp := openapi.WarehouseFindWebResp{
		Id:       e.Id,
		Name:     e.Name,
		Address:  e.Address,
		Capacity: e.Capacity,
	}
	return resp
}

func (e WarehouseEntity) FromValue(value interface{}, params interface{}) (WarehouseEntity, error) {
	switch value.(type) {
	case openapi.WarehouseCreateWebReq:
		v := value.(openapi.WarehouseCreateWebReq)
		ent := WarehouseEntity{
			Name:     v.Name,
			Address:  v.Address,
			Capacity: v.Capacity,
		}
		return ent, nil
	case openapi.WarehouseUpdateWebReq:
		v := value.(openapi.WarehouseUpdateWebReq)
		p := params.(openapi.WarehouseUpdateWebParam)
		ent := WarehouseEntity{
			Id:       p.Id,
			Name:     v.Name,
			Address:  v.Address,
			Capacity: v.Capacity,
		}
		return ent, nil
	default:
		return WarehouseEntity{}, fmt.Errorf("failed to convert request to WarehouseEntity: %v", value)
	}
}
