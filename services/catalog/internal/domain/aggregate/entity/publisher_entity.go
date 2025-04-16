package entity

import (
	"fmt"
	"time"

	"github.com/codejsha/bookstore-microservices/catalog/internal/application/port/openapi"
	"github.com/codejsha/bookstore-microservices/catalog/internal/application/port/pb/publisherpb"
)

type PublisherEntity struct {
	Id        *int64    `gorm:"primaryKey;autoIncrement"`
	Name      string    `gorm:"not null"`
	Address   string    `gorm:"null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (e PublisherEntity) TableName() string {
	return "publisher"
}

func (e PublisherEntity) ToPublisherFindWebResp() openapi.PublisherFindWebResp {
	resp := openapi.PublisherFindWebResp{
		Id:      *e.Id,
		Name:    e.Name,
		Address: e.Address,
	}
	return resp
}

func (e PublisherEntity) ToPublisherFindProtoResp() *publisherpb.PublisherFindProtoResp {
	resp := &publisherpb.PublisherFindProtoResp{
		Id:      *e.Id,
		Name:    e.Name,
		Address: &e.Address,
	}
	return resp
}

func (e PublisherEntity) FromValue(value interface{}, id *int64) (PublisherEntity, error) {
	switch value.(type) {
	case openapi.PublisherCreateWebReq:
		v := value.(openapi.PublisherCreateWebReq)
		ent := PublisherEntity{
			Name:    v.Name,
			Address: v.Address,
		}
		return ent, nil
	case openapi.PublisherUpdateWebReq:
		v := value.(openapi.PublisherUpdateWebReq)
		ent := PublisherEntity{
			Id:      id,
			Name:    v.Name,
			Address: v.Address,
		}
		return ent, nil
	default:
		return PublisherEntity{}, fmt.Errorf("failed to convert request to PublisherEntity: %v", value)
	}
}
