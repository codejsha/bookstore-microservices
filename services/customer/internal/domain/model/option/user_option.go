package option

import (
	"github.com/codejsha/shared-library-go/pkg/pagination"
)

type UserQueryOption struct {
	email *string
	name  *string
	phone *string
	page  pagination.PageOption
}

type UserQueryOptionFunc func(*UserQueryOption)

func NewUserQueryOption(opts ...UserQueryOptionFunc) UserQueryOption {
	var o UserQueryOption
	for _, opt := range opts {
		opt(&o)
	}
	return o
}

func (o UserQueryOption) Email() *string {
	return o.email
}

func (o UserQueryOption) Name() *string {
	return o.name
}

func (o UserQueryOption) Phone() *string {
	return o.phone
}

func (o UserQueryOption) Page() pagination.PageOption {
	return o.page
}

func (UserQueryOption) WithEmail(email *string) UserQueryOptionFunc {
	return func(o *UserQueryOption) {
		o.email = email
	}
}

func (UserQueryOption) WithName(name *string) UserQueryOptionFunc {
	return func(o *UserQueryOption) {
		o.name = name
	}
}

func (UserQueryOption) WithPhone(phone *string) UserQueryOptionFunc {
	return func(o *UserQueryOption) {
		o.phone = phone
	}
}

func (UserQueryOption) WithPage(page pagination.PageOption) UserQueryOptionFunc {
	return func(o *UserQueryOption) {
		o.page = page
	}
}
