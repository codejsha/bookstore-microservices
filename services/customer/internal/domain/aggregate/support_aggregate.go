package aggregate

import (
	"time"

	"github.com/codejsha/bookstore-microservices/customer/internal/domain/constant"
)

type TicketAggregate struct {
	Uid         string
	CustomerUid string
	CategoryUid *string
	AssigneeUid *string
	Subject     string
	Description string
	Status      constant.TicketStatus
	Priority    constant.TicketPriority
	ResolvedAt  *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type TicketCommentAggregate struct {
	Uid        string
	TicketUid  string
	AuthorUid  string
	AuthorRole constant.AuthorRole
	Body       string
	Internal   bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type TicketCategoryAggregate struct {
	Uid         string
	Name        string
	Description *string
	ParentUid   *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type FaqAggregate struct {
	Uid         string
	CategoryUid *string
	Question    string
	Answer      string
	ViewCount   int64
	Published   bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
