package option

import (
	"github.com/codejsha/shared-library-go/pkg/pagination"
)

type WishlistQueryOption struct {
	userUid  *string
	bookUids *[]string
	page     pagination.PageOption
}

type WishlistQueryOptionFunc func(*WishlistQueryOption)

func NewWishlistQueryOption(opts ...WishlistQueryOptionFunc) WishlistQueryOption {
	var o WishlistQueryOption
	for _, opt := range opts {
		opt(&o)
	}
	return o
}

func (o WishlistQueryOption) UserUid() *string {
	return o.userUid
}

func (o WishlistQueryOption) BookUids() *[]string {
	return o.bookUids
}

func (o WishlistQueryOption) Page() pagination.PageOption {
	return o.page
}

func (WishlistQueryOption) WithUserUid(userUid *string) WishlistQueryOptionFunc {
	return func(o *WishlistQueryOption) {
		o.userUid = userUid
	}
}

func (WishlistQueryOption) WithBookUids(bookUids *[]string) WishlistQueryOptionFunc {
	return func(o *WishlistQueryOption) {
		o.bookUids = bookUids
	}
}

func (WishlistQueryOption) WithPage(page pagination.PageOption) WishlistQueryOptionFunc {
	return func(o *WishlistQueryOption) {
		o.page = page
	}
}
