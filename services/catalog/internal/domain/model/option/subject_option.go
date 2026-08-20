package option

import (
	"github.com/codejsha/shared-library-go/pkg/pagination"
)

type SubjectQueryOption struct {
	name *string
	page pagination.PageOption
}

type SubjectQueryOptionFunc func(*SubjectQueryOption)

func NewSubjectQueryOption(opts ...SubjectQueryOptionFunc) SubjectQueryOption {
	var o SubjectQueryOption
	for _, opt := range opts {
		opt(&o)
	}
	return o
}

func (o SubjectQueryOption) Name() *string {
	return o.name
}

func (o SubjectQueryOption) Page() pagination.PageOption {
	return o.page
}

func (SubjectQueryOption) WithName(name *string) SubjectQueryOptionFunc {
	return func(o *SubjectQueryOption) {
		o.name = name
	}
}

func (SubjectQueryOption) WithPage(page pagination.PageOption) SubjectQueryOptionFunc {
	return func(o *SubjectQueryOption) {
		o.page = page
	}
}
