package repo

import (
	"context"
	"errors"
	"time"

	genrepo "github.com/codejsha/bookstore-microservices/customer/generated/infrastructure/port/repo"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/model/option"
)

var ErrCategoryNameExists = errors.New("category name already exists")

var ErrUnknownReference = errors.New("referenced resource not found")

type TicketRepo interface {
	FindAll(ctx context.Context, opt option.TicketQueryOption) (int64, []*TicketResult, error)
	FindByUid(ctx context.Context, uid string) (*TicketResult, error)
	Insert(ctx context.Context, ticket TicketCreate) (*TicketResult, error)
	UpdateByUid(ctx context.Context, uid string, apply func(*TicketResult)) (*TicketResult, error)
	genrepo.CustomerTicketRepo
}

type TicketCommentRepo interface {
	FindAllByTicketId(ctx context.Context, ticketId int64, includeInternal bool) ([]*TicketCommentResult, error)
	Insert(ctx context.Context, comment TicketCommentCreate) (*TicketCommentResult, error)
	genrepo.CustomerTicketCommentRepo
}

type TicketCategoryRepo interface {
	FindAll(ctx context.Context) ([]*TicketCategoryResult, error)
	FindByUid(ctx context.Context, uid string) (*TicketCategoryResult, error)
	FindByName(ctx context.Context, name string) (*TicketCategoryResult, error)
	FindIdByUid(ctx context.Context, uid string) (*int64, error)
	Insert(ctx context.Context, category TicketCategoryCreate) (*TicketCategoryResult, error)
	UpdateByUid(ctx context.Context, uid string, apply func(*TicketCategoryResult)) (*TicketCategoryResult, error)
	SoftDeleteByUid(ctx context.Context, uid string) (bool, error)
	genrepo.CustomerTicketCategoryRepo
}

type FaqRepo interface {
	FindAll(ctx context.Context, opt option.FaqQueryOption) (int64, []*FaqResult, error)
	FindByUid(ctx context.Context, uid string) (*FaqResult, error)
	Insert(ctx context.Context, faq FaqCreate) (*FaqResult, error)
	UpdateByUid(ctx context.Context, uid string, apply func(*FaqResult)) (*FaqResult, error)
	SoftDeleteByUid(ctx context.Context, uid string) (bool, error)
	IncrementViewCount(ctx context.Context, uid string) error
	genrepo.CustomerFaqRepo
}

type TicketCreate struct {
	CustomerUid string
	CategoryId  *int64
	Subject     string
	Description string
	Status      string
	Priority    string
}

type TicketResult struct {
	Id          int64
	Uid         string
	CustomerUid string
	CategoryId  *int64
	CategoryUid *string
	AssigneeUid *string
	Subject     string
	Description string
	Status      string
	Priority    string
	ResolvedAt  *time.Time
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}

type TicketCommentCreate struct {
	TicketId   int64
	AuthorUid  string
	AuthorRole string
	Body       string
	Internal   bool
}

type TicketCommentResult struct {
	Id         int64
	Uid        string
	TicketId   int64
	AuthorUid  string
	AuthorRole string
	Body       string
	Internal   bool
	CreatedAt  time.Time
	UpdatedAt  *time.Time
}

type TicketCategoryCreate struct {
	Name        string
	Description *string
	ParentId    *int64
}

type TicketCategoryResult struct {
	Id          int64
	Uid         string
	Name        string
	Description *string
	ParentId    *int64
	ParentUid   *string
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}

type FaqCreate struct {
	CategoryId *int64
	Question   string
	Answer     string
	Published  bool
}

type FaqResult struct {
	Id          int64
	Uid         string
	CategoryId  *int64
	CategoryUid *string
	Question    string
	Answer      string
	ViewCount   int64
	Published   bool
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}
