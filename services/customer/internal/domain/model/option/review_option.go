package option

import "github.com/codejsha/shared-library-go/pkg/pagination"

type ReviewQueryOption struct {
	userUid *string
	bookUid *string
	rating  *int32
	page    pagination.PageOption
}

type ReviewQueryOptionFunc func(*ReviewQueryOption)

func NewReviewQueryOption(opts ...ReviewQueryOptionFunc) ReviewQueryOption {
	o := ReviewQueryOption{}
	for _, opt := range opts {
		opt(&o)
	}
	return o
}

func WithReviewUserUid(userUid string) ReviewQueryOptionFunc {
	return func(o *ReviewQueryOption) { o.userUid = &userUid }
}

func WithReviewBookUid(bookUid string) ReviewQueryOptionFunc {
	return func(o *ReviewQueryOption) { o.bookUid = &bookUid }
}

func WithReviewRating(rating int32) ReviewQueryOptionFunc {
	return func(o *ReviewQueryOption) { o.rating = &rating }
}

func WithReviewPage(page pagination.PageOption) ReviewQueryOptionFunc {
	return func(o *ReviewQueryOption) { o.page = page }
}

func (o ReviewQueryOption) UserUid() *string            { return o.userUid }
func (o ReviewQueryOption) BookUid() *string            { return o.bookUid }
func (o ReviewQueryOption) Rating() *int32              { return o.rating }
func (o ReviewQueryOption) Page() pagination.PageOption { return o.page }
