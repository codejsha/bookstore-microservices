package option

import (
	"github.com/codejsha/shared-library-go/pkg/pagination"
)

type AuthorQueryOption struct {
	name    *string
	workUid *string
	olKey   *string
	page    pagination.PageOption
}

type AuthorQueryOptionFunc func(*AuthorQueryOption)

func NewAuthorQueryOption(opts ...AuthorQueryOptionFunc) AuthorQueryOption {
	var o AuthorQueryOption
	for _, opt := range opts {
		opt(&o)
	}
	return o
}

func (o AuthorQueryOption) Name() *string {
	return o.name
}

func (o AuthorQueryOption) WorkUid() *string {
	return o.workUid
}

func (o AuthorQueryOption) OlKey() *string {
	return o.olKey
}

func (o AuthorQueryOption) Page() pagination.PageOption {
	return o.page
}

func (AuthorQueryOption) WithName(name *string) AuthorQueryOptionFunc {
	return func(o *AuthorQueryOption) {
		o.name = name
	}
}

func (AuthorQueryOption) WithWorkUid(workUid *string) AuthorQueryOptionFunc {
	return func(o *AuthorQueryOption) {
		o.workUid = workUid
	}
}

func (AuthorQueryOption) WithOlKey(olKey *string) AuthorQueryOptionFunc {
	return func(o *AuthorQueryOption) {
		o.olKey = olKey
	}
}

func (AuthorQueryOption) WithPage(page pagination.PageOption) AuthorQueryOptionFunc {
	return func(o *AuthorQueryOption) {
		o.page = page
	}
}
