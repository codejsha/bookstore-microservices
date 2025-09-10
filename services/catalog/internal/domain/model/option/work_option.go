package option

import (
	"github.com/codejsha/shared-library-go/pkg/pagination"
)

type WorkQueryOption struct {
	title      *string
	authorUid  *string
	subjectUid *string
	olKey      *string
	page       pagination.PageOption
}

type WorkQueryOptionFunc func(*WorkQueryOption)

func NewWorkQueryOption(opts ...WorkQueryOptionFunc) WorkQueryOption {
	var o WorkQueryOption
	for _, opt := range opts {
		opt(&o)
	}
	return o
}

func (o WorkQueryOption) Title() *string {
	return o.title
}

func (o WorkQueryOption) AuthorUid() *string {
	return o.authorUid
}

func (o WorkQueryOption) SubjectUid() *string {
	return o.subjectUid
}

func (o WorkQueryOption) OlKey() *string {
	return o.olKey
}

func (o WorkQueryOption) Page() pagination.PageOption {
	return o.page
}

func (WorkQueryOption) WithTitle(title *string) WorkQueryOptionFunc {
	return func(o *WorkQueryOption) {
		o.title = title
	}
}

func (WorkQueryOption) WithAuthorUid(authorUid *string) WorkQueryOptionFunc {
	return func(o *WorkQueryOption) {
		o.authorUid = authorUid
	}
}

func (WorkQueryOption) WithSubjectUid(subjectUid *string) WorkQueryOptionFunc {
	return func(o *WorkQueryOption) {
		o.subjectUid = subjectUid
	}
}

func (WorkQueryOption) WithOlKey(olKey *string) WorkQueryOptionFunc {
	return func(o *WorkQueryOption) {
		o.olKey = olKey
	}
}

func (WorkQueryOption) WithPage(page pagination.PageOption) WorkQueryOptionFunc {
	return func(o *WorkQueryOption) {
		o.page = page
	}
}
