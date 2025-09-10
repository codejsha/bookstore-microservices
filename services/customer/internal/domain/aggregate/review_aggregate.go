package aggregate

import (
	"time"
)

type ReviewAggregate struct {
	Uid       string
	UserUid   string
	BookUid   string
	Rating    int32
	Title     *string
	Content   *string
	CreatedAt time.Time
	UpdatedAt *time.Time
}
