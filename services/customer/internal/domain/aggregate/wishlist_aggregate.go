package aggregate

import (
	"github.com/codejsha/bookstore-microservices/customer/internal/application/port/openapi"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/aggregate/entity"
)

func NewWishlistAggregate(
	userId string,
	wishlists []entity.WishlistEntity,
) *WishlistAggregate {
	agg := &WishlistAggregate{
		Id:        userId,
		Wishlists: wishlists,
	}
	return agg
}

type WishlistAggregate struct {
	// user id
	Id        string
	Wishlists []entity.WishlistEntity
}

func (a *WishlistAggregate) BookIds() []int64 {
	bookIds := make([]int64, len(a.Wishlists))
	for i, w := range a.Wishlists {
		bookIds[i] = w.BookId
	}
	return bookIds
}

func (a *WishlistAggregate) UpdateWishBooks(ents []entity.WishlistEntity) {
	for _, e := range ents {
		found := false
		for _, ent := range a.Wishlists {
			if ent.BookId == e.BookId {
				found = true
				break
			}
		}
		if !found {
			a.Wishlists = append(a.Wishlists, e)
		}
	}
}

func (a *WishlistAggregate) DeleteWishBooks(bookIds []int64) {
	for i, ent := range a.Wishlists {
		for _, bookId := range bookIds {
			if ent.BookId == bookId {
				a.Wishlists = append(a.Wishlists[:i], a.Wishlists[i+1:]...)
				break
			}
		}
	}
}

func (a *WishlistAggregate) ToWishlistFindResponse() openapi.WishlistFindWebResp {
	resp := openapi.WishlistFindWebResp{BookIds: a.BookIds()}
	return resp
}
