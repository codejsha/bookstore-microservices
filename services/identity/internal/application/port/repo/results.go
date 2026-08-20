package repo

import "time"

type UserResult struct {
	Id          int64
	IdpId       *string
	Email       string
	FirstName   string
	LastName    string
	Phone       *string
	Roles       []string
	Status      string
	LastLoginAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}
