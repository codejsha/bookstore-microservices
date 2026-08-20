package option

import (
	"github.com/codejsha/shared-library-go/pkg/pagination"
)

type EditionQueryOption struct {
	title        *string
	workUid      *string
	publisherUid *string
	isbn         *string
	olKey        *string
	language     *string
	page         pagination.PageOption
}

type EditionQueryOptionFunc func(*EditionQueryOption)

func NewEditionQueryOption(opts ...EditionQueryOptionFunc) EditionQueryOption {
	var o EditionQueryOption
	for _, opt := range opts {
		opt(&o)
	}
	return o
}

func (o EditionQueryOption) Title() *string {
	return o.title
}

func (o EditionQueryOption) WorkUid() *string {
	return o.workUid
}

func (o EditionQueryOption) PublisherUid() *string {
	return o.publisherUid
}

func (o EditionQueryOption) Isbn() *string {
	return o.isbn
}

func (o EditionQueryOption) OlKey() *string {
	return o.olKey
}

func (o EditionQueryOption) Language() *string {
	return o.language
}

func (o EditionQueryOption) Page() pagination.PageOption {
	return o.page
}

func (EditionQueryOption) WithTitle(title *string) EditionQueryOptionFunc {
	return func(o *EditionQueryOption) {
		o.title = title
	}
}

func (EditionQueryOption) WithWorkUid(workUid *string) EditionQueryOptionFunc {
	return func(o *EditionQueryOption) {
		o.workUid = workUid
	}
}

func (EditionQueryOption) WithPublisherUid(publisherUid *string) EditionQueryOptionFunc {
	return func(o *EditionQueryOption) {
		o.publisherUid = publisherUid
	}
}

func (EditionQueryOption) WithIsbn(isbn *string) EditionQueryOptionFunc {
	return func(o *EditionQueryOption) {
		o.isbn = isbn
	}
}

func (EditionQueryOption) WithOlKey(olKey *string) EditionQueryOptionFunc {
	return func(o *EditionQueryOption) {
		o.olKey = olKey
	}
}

func (EditionQueryOption) WithLanguage(language *string) EditionQueryOptionFunc {
	return func(o *EditionQueryOption) {
		o.language = language
	}
}

func (EditionQueryOption) WithPage(page pagination.PageOption) EditionQueryOptionFunc {
	return func(o *EditionQueryOption) {
		o.page = page
	}
}
