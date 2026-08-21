package entity

type BookAuthorMappingEntity struct {
	BookId   int64 `gorm:"primaryKey"`
	AuthorId int64 `gorm:"primaryKey"`
}

func (e BookAuthorMappingEntity) TableName() string {
	return "book_author_mapping"
}
