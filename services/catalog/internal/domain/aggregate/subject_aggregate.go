package aggregate

import "time"

type SubjectAggregate struct {
	Uid       string
	Name      string
	CreatedAt time.Time
	UpdatedAt *time.Time
}
