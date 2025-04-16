package handler

import (
	"context"
	"fmt"

	"github.com/codejsha/bookstore-microservices/customer/internal/application/port/repo"
	"github.com/codejsha/bookstore-microservices/customer/internal/application/usecase"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/aggregate/entity"
)

var _ WishlistHandler = (*wishlistHandlerImpl)(nil)

func NewWishlistHandler(
	wishlistRepo repo.WishlistRepo,
) WishlistHandler {
	return &wishlistHandlerImpl{
		wishlistRepo: wishlistRepo,
	}
}

type WishlistHandler interface {
	Handle(ctx context.Context, command usecase.WishlistCommand) (entity.WishlistEntity, error)
	HandleFindAll(ctx context.Context, query usecase.WishlistFindAllQuery) (int64, []entity.WishlistEntity, error)
	HandleFind(ctx context.Context, query usecase.WishlistFindQuery) (entity.WishlistEntity, error)
}

type wishlistHandlerImpl struct {
	wishlistRepo repo.WishlistRepo
}

func (h wishlistHandlerImpl) Handle(ctx context.Context, command usecase.WishlistCommand) (entity.WishlistEntity, error) {
	switch command.(type) {
	case usecase.WishlistCreateCommand:
		c := command.(usecase.WishlistCreateCommand)
		return h.wishlistRepo.Create(ctx, c.Ent)
	case usecase.WishlistUpdateCommand:
		c := command.(usecase.WishlistUpdateCommand)
		return h.wishlistRepo.Update(ctx, c.Ent)
	case usecase.WishlistDeleteCommand:
		c := command.(usecase.WishlistDeleteCommand)
		return entity.WishlistEntity{}, h.wishlistRepo.Delete(ctx, c.Id)
	default:
		return entity.WishlistEntity{}, fmt.Errorf("unknown command")
	}
}

func (h wishlistHandlerImpl) HandleFindAll(ctx context.Context, query usecase.WishlistFindAllQuery) (int64, []entity.WishlistEntity, error) {
	return h.wishlistRepo.FindAll(ctx, query.Cond)
}

func (h wishlistHandlerImpl) HandleFind(ctx context.Context, query usecase.WishlistFindQuery) (entity.WishlistEntity, error) {
	return h.wishlistRepo.Find(ctx, query.Id)
}
