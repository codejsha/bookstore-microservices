package option

import (
	"github.com/codejsha/shared-library-go/pkg/pagination"

	"github.com/codejsha/bookstore-microservices/customer/internal/domain/constant"
)

type TicketQueryOption struct {
	customerUid *string
	assigneeUid *string
	categoryUid *string
	status      *constant.TicketStatus
	priority    *constant.TicketPriority
	page        pagination.PageOption
}

type TicketQueryOptionFunc func(*TicketQueryOption)

func NewTicketQueryOption(opts ...TicketQueryOptionFunc) TicketQueryOption {
	o := TicketQueryOption{}
	for _, opt := range opts {
		opt(&o)
	}
	return o
}

func WithTicketCustomerUid(customerUid string) TicketQueryOptionFunc {
	return func(o *TicketQueryOption) { o.customerUid = &customerUid }
}

func WithTicketAssigneeUid(assigneeUid string) TicketQueryOptionFunc {
	return func(o *TicketQueryOption) { o.assigneeUid = &assigneeUid }
}

func WithTicketCategoryUid(categoryUid string) TicketQueryOptionFunc {
	return func(o *TicketQueryOption) { o.categoryUid = &categoryUid }
}

func WithTicketStatus(status constant.TicketStatus) TicketQueryOptionFunc {
	return func(o *TicketQueryOption) { o.status = &status }
}

func WithTicketPriority(priority constant.TicketPriority) TicketQueryOptionFunc {
	return func(o *TicketQueryOption) { o.priority = &priority }
}

func WithTicketPage(page pagination.PageOption) TicketQueryOptionFunc {
	return func(o *TicketQueryOption) { o.page = page }
}

func (o TicketQueryOption) CustomerUid() *string               { return o.customerUid }
func (o TicketQueryOption) AssigneeUid() *string               { return o.assigneeUid }
func (o TicketQueryOption) CategoryUid() *string               { return o.categoryUid }
func (o TicketQueryOption) Status() *constant.TicketStatus     { return o.status }
func (o TicketQueryOption) Priority() *constant.TicketPriority { return o.priority }
func (o TicketQueryOption) Page() pagination.PageOption        { return o.page }

type FaqQueryOption struct {
	categoryUid   *string
	query         *string
	publishedOnly bool
	page          pagination.PageOption
}

type FaqQueryOptionFunc func(*FaqQueryOption)

func NewFaqQueryOption(opts ...FaqQueryOptionFunc) FaqQueryOption {
	o := FaqQueryOption{publishedOnly: true}
	for _, opt := range opts {
		opt(&o)
	}
	return o
}

func WithFaqCategoryUid(categoryUid string) FaqQueryOptionFunc {
	return func(o *FaqQueryOption) { o.categoryUid = &categoryUid }
}

func WithFaqQuery(query string) FaqQueryOptionFunc {
	return func(o *FaqQueryOption) { o.query = &query }
}

func WithFaqPublishedOnly(publishedOnly bool) FaqQueryOptionFunc {
	return func(o *FaqQueryOption) { o.publishedOnly = publishedOnly }
}

func WithFaqPage(page pagination.PageOption) FaqQueryOptionFunc {
	return func(o *FaqQueryOption) { o.page = page }
}

func (o FaqQueryOption) CategoryUid() *string        { return o.categoryUid }
func (o FaqQueryOption) Query() *string              { return o.query }
func (o FaqQueryOption) PublishedOnly() bool         { return o.publishedOnly }
func (o FaqQueryOption) Page() pagination.PageOption { return o.page }
