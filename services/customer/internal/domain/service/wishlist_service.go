package service

import (
	"context"

	"github.com/codejsha/bookstore-microservices/customer/internal/application/usecase"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/aggregate"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/aggregate/entity"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/handler"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/model/condition"
)

var _ usecase.WishlistUseCase = (*wishlistService)(nil)

func NewWishlistService(
	wishlistHandler handler.WishlistHandler,
) usecase.WishlistUseCase {
	return &wishlistService{
		wishlistHandler: wishlistHandler,
	}
}

type wishlistService struct {
	wishlistHandler handler.WishlistHandler
}

func (s wishlistService) FindWishlist(ctx context.Context, cond condition.WishlistCondition) (*aggregate.WishlistAggregate, error) {
	// find wishlists
	_, wishlists, err := s.wishlistHandler.HandleFindAll(ctx, usecase.WishlistFindAllQuery{Cond: cond})
	if err != nil {
		return nil, err
	}
	wishlistAgg := aggregate.NewWishlistAggregate(*cond.UserId, wishlists)
	return wishlistAgg, nil
}

func (s wishlistService) UpdateWishlist(ctx context.Context, ents []entity.WishlistEntity) (*aggregate.WishlistAggregate, error) {
	// find wishlists
	userId := ents[0].UserId
	cond := condition.WishlistCondition{UserId: &userId}
	_, wishlists, err := s.wishlistHandler.HandleFindAll(ctx, usecase.WishlistFindAllQuery{Cond: cond})
	if err != nil {
		return nil, err
	}

	// update wishlists
	for _, ent := range ents {
		found := false
		for _, ent := range wishlists {
			if ent.BookId == ent.BookId {
				found = true
				break
			}
		}
		if !found {
			ent, err := s.wishlistHandler.Handle(ctx, usecase.WishlistCreateCommand{Ent: ent})
			if err != nil {
				return nil, err
			}
			wishlists = append(wishlists, ent)
		}
	}

	wishlistAgg := aggregate.NewWishlistAggregate(userId, wishlists)
	return wishlistAgg, nil
}

func (s wishlistService) DeleteWishlist(ctx context.Context, bookIds []int64) error {
	// find wishlists
	cond := condition.WishlistCondition{BookIds: bookIds}
	_, wishlists, err := s.wishlistHandler.HandleFindAll(ctx, usecase.WishlistFindAllQuery{Cond: cond})
	if err != nil {
		return err
	}

	// delete wishlists
	for _, ent := range wishlists {
		_, err := s.wishlistHandler.Handle(ctx, usecase.WishlistDeleteCommand{Id: *ent.Id})
		if err != nil {
			return err
		}
	}

	return nil
}
