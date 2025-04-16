package entity

import (
	"time"

	"github.com/codejsha/bookstore-microservices/customer/internal/application/port/openapi"
)

type PointEntity struct {
	Id        *int64    `gorm:"primaryKey;autoIncrement"`
	UserId    string    `gorm:"primaryKey;autoIncrement"`
	Point     int64     `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (e PointEntity) TableName() string {
	return "point"
}

func (e PointEntity) ToPointFindResponse() openapi.PointFindWebResp {
	resp := openapi.PointFindWebResp{
		Id:    *e.Id,
		Point: e.Point,
	}
	return resp
}
