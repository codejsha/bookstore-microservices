package entity

import (
	"fmt"
	"time"

	"github.com/codejsha/bookstore-microservices/catalog/internal/application/port/openapi"
)

type BookEntity struct {
	Id          *int64    `gorm:"primaryKey;autoIncrement"`
	Title       string    `gorm:"not null"`
	Isbn        string    `gorm:"null"`
	Price       float64   `gorm:"null"`
	Description string    `gorm:"null"`
	CategoryId  int64     `gorm:"not null"`
	PublisherId int64     `gorm:"not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

func (e BookEntity) TableName() string {
	return "book"
}

func (e BookEntity) FromValue(value interface{}, id *int64) (BookEntity, error) {
	switch value.(type) {
	case openapi.BookCreateWebReq:
		v := value.(openapi.BookCreateWebReq)
		dto := BookEntity{
			Title:       v.Title,
			Isbn:        v.Isbn,
			Price:       v.Price,
			Description: v.Description,
			CategoryId:  v.CategoryId,
			PublisherId: v.PublisherId,
		}
		return dto, nil
	case openapi.BookUpdateWebReq:
		v := value.(openapi.BookUpdateWebReq)
		ent := BookEntity{
			Id:          id,
			Title:       v.Title,
			Isbn:        v.Isbn,
			Price:       v.Price,
			Description: v.Description,
			CategoryId:  v.CategoryId,
			PublisherId: v.PublisherId,
		}
		return ent, nil
	default:
		return BookEntity{}, fmt.Errorf("failed to convert request to BookEntity: %v", value)
	}
}
