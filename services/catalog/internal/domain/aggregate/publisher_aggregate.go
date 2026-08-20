package aggregate

import "time"

type PublisherAggregate struct {
	Uid       string
	Name      string
	Address   *string
	OlKey     *string
	CreatedAt time.Time
	UpdatedAt *time.Time
}
