package entity

import (
	"fmt"
	"time"

	"github.com/codejsha/bookstore-microservices/catalog/internal/application/port/openapi"
	"github.com/codejsha/bookstore-microservices/catalog/internal/application/port/pb/authorpb"
)

type AuthorEntity struct {
	Id        *int64    `gorm:"primaryKey;autoIncrement"`
	Name      string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (e AuthorEntity) TableName() string {
	return "author"
}

func (e AuthorEntity) ToAuthorFindWebResp() openapi.AuthorFindWebResp {
	resp := openapi.AuthorFindWebResp{
		Id:   *e.Id,
		Name: e.Name,
	}
	return resp
}

func (e AuthorEntity) ToAuthorFindProtoResp() *authorpb.AuthorFindProtoResp {
	resp := &authorpb.AuthorFindProtoResp{
		Id:   *e.Id,
		Name: e.Name,
	}
	return resp
}

func (e AuthorEntity) FromValue(value interface{}, id *int64) (AuthorEntity, error) {
	switch value.(type) {
	case openapi.AuthorCreateWebReq:
		v := value.(openapi.AuthorCreateWebReq)
		ent := AuthorEntity{
			Name: v.Name,
		}
		return ent, nil
	case openapi.AuthorUpdateWebReq:
		v := value.(openapi.AuthorUpdateWebReq)
		ent := AuthorEntity{
			Id:   id,
			Name: v.Name,
		}
		return ent, nil
	default:
		return AuthorEntity{}, fmt.Errorf("failed to convert request to AuthorEntity: %v", value)
	}
}
