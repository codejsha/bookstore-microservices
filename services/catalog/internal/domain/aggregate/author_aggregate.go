package aggregate

import "time"

type AuthorAggregate struct {
	Uid            string
	Name           string
	Bio            *string
	BirthDate      *string
	DeathDate      *string
	PhotoUids      []string
	AlternateNames []string
	OlKey          *string
	CreatedAt      time.Time
	UpdatedAt      *time.Time
}
