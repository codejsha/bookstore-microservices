package entity

import (
	"time"
)

type WishlistEntity struct {
	Id        *int64    `gorm:"primaryKey;autoIncrement"`
	UserId    string    `gorm:"type:varchar(36);not null"`
	BookId    int64     `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (e WishlistEntity) TableName() string {
	return "wishlist"
}
