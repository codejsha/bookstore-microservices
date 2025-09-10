package aggregate

import "time"

type EditionAggregate struct {
	Uid            string
	Title          string
	Isbn10         *string
	Isbn13         *string
	NumberOfPages  *int32
	PublishDate    *string
	CoverUids      []string
	Languages      []string
	PhysicalFormat *string
	Description    *string
	Work           *WorkAggregate
	Publisher      *PublisherAggregate
	OlKey          *string
	CreatedAt      time.Time
	UpdatedAt      *time.Time
}
