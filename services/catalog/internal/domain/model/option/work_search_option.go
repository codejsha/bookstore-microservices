package option

import (
	"github.com/codejsha/shared-library-go/pkg/pagination"
)

type WorkSearchOption struct {
	query      *string
	authorUid  *string
	subjectUid *string
	page       pagination.PageOption
}

type WorkSearchOptionFunc func(*WorkSearchOption)

func NewWorkSearchOption(opts ...WorkSearchOptionFunc) WorkSearchOption {
	var o WorkSearchOption
	for _, opt := range opts {
		opt(&o)
	}
	return o
}

func (o WorkSearchOption) Query() *string      { return o.query }
func (o WorkSearchOption) AuthorUid() *string  { return o.authorUid }
func (o WorkSearchOption) SubjectUid() *string { return o.subjectUid }
func (o WorkSearchOption) Page() pagination.PageOption {
	return o.page
}

func (WorkSearchOption) WithQuery(q *string) WorkSearchOptionFunc {
	return func(o *WorkSearchOption) { o.query = q }
}

func (WorkSearchOption) WithAuthorUid(uid *string) WorkSearchOptionFunc {
	return func(o *WorkSearchOption) { o.authorUid = uid }
}

func (WorkSearchOption) WithSubjectUid(uid *string) WorkSearchOptionFunc {
	return func(o *WorkSearchOption) { o.subjectUid = uid }
}

func (WorkSearchOption) WithPage(page pagination.PageOption) WorkSearchOptionFunc {
	return func(o *WorkSearchOption) { o.page = page }
}
