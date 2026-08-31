package restcontroller

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/codejsha/bookstore-microservices/customer/generated/application/port/openapi"
	"github.com/codejsha/bookstore-microservices/customer/internal/application/usecase"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/aggregate"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/model/command"
	"github.com/codejsha/bookstore-microservices/customer/internal/infrastructure/httpx"
)

var _ openapi.WishlistApi = (*wishlistController)(nil)

type wishlistController struct {
	customerUseCase usecase.CustomerUseCase
}

func NewWishlistController(customerUseCase usecase.CustomerUseCase) openapi.WishlistApi {
	return &wishlistController{customerUseCase: customerUseCase}
}

func (c *wishlistController) WishlistsGet(ctx context.Context, uid string) (*openapi.WishlistResponse, error) {
	wishlist, err := c.customerUseCase.GetWishlist(ctx, uid)
	if err != nil {
		return nil, err
	}

	resp := toWishlistResponse(wishlist)
	return &resp, nil
}

func (c *wishlistController) WishlistsAddBooks(
	ctx context.Context,
	uid string,
	req openapi.WishlistAddRequest,
) (*openapi.WishlistResponse, error) {
	cmd := command.WishlistAddCommand{
		UserUid:  uid,
		BookUids: req.BookUids,
	}

	wishlist, err := c.customerUseCase.AddBooksToWishlist(ctx, cmd)
	if err != nil {
		return nil, httpx.MapBusinessError(ctx, err)
	}

	dispatchSideEffects(context.WithoutCancel(ctx), "wishlist", uid, "added", logrus.Fields{"books": len(req.BookUids)})

	resp := toWishlistResponse(wishlist)
	return &resp, nil
}

func (c *wishlistController) WishlistsRemoveBooks(
	ctx context.Context,
	uid string,
	req openapi.WishlistRemoveRequest,
) (*openapi.WishlistResponse, error) {
	cmd := command.WishlistRemoveCommand{
		UserUid:  uid,
		BookUids: req.BookUids,
	}

	wishlist, err := c.customerUseCase.RemoveBooksFromWishlist(ctx, cmd)
	if err != nil {
		return nil, httpx.MapBusinessError(ctx, err)
	}

	dispatchSideEffects(context.WithoutCancel(ctx), "wishlist", uid, "removed", logrus.Fields{"books": len(req.BookUids)})

	resp := toWishlistResponse(wishlist)
	return &resp, nil
}

// ─── mapping helpers ────────────────────────────────────────────────────────

func toWishlistResponse(a *aggregate.WishlistAggregate) openapi.WishlistResponse {
	return openapi.WishlistResponse{BookUids: a.BookUids}
}
