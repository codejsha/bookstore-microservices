package option

import (
	"github.com/codejsha/shared-library-go/pkg/pagination"
)

type PublisherQueryOption struct {
	name  *string
	olKey *string
	page  pagination.PageOption
}

type PublisherQueryOptionFunc func(*PublisherQueryOption)

func NewPublisherQueryOption(opts ...PublisherQueryOptionFunc) PublisherQueryOption {
	var o PublisherQueryOption
	for _, opt := range opts {
		opt(&o)
	}
	return o
}

func (o PublisherQueryOption) Name() *string {
	return o.name
}

func (o PublisherQueryOption) Page() pagination.PageOption {
	return o.page
}

func (PublisherQueryOption) WithName(name *string) PublisherQueryOptionFunc {
	return func(o *PublisherQueryOption) {
		o.name = name
	}
}

func (PublisherQueryOption) WithOlKey(olKey *string) PublisherQueryOptionFunc {
	return func(o *PublisherQueryOption) {
		o.olKey = olKey
	}
}

func (PublisherQueryOption) WithPage(page pagination.PageOption) PublisherQueryOptionFunc {
	return func(o *PublisherQueryOption) {
		o.page = page
	}
}

func (o PublisherQueryOption) OlKey() *string {
	return o.olKey
}
