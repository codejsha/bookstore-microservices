package aggregate

import "time"

type WorkAggregate struct {
	Uid              string
	Title            string
	Description      *string
	CoverUids        []string
	FirstPublishDate *string
	OlKey            *string
	Authors          []*AuthorAggregate
	Subjects         []*SubjectAggregate
	Editions         []*EditionAggregate
	CreatedAt        time.Time
	UpdatedAt        *time.Time
}
