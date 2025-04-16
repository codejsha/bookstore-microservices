package usecase

import (
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/aggregate/entity"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/model/condition"
)

type (
	PointCommand       interface{ isPointCommand() }
	PointUpdateCommand struct{ Ent entity.PointEntity }
)

func (PointUpdateCommand) isPointCommand() {}

type (
	PointQuery        interface{ isPointQuery() }
	PointFindAllQuery struct{ Cond condition.PointCondition }
	PointFindQuery    struct{ Id int64 }
)

func (PointFindAllQuery) isPointQuery() {}
func (PointFindQuery) isPointQuery()    {}

type (
	WishlistCommand       interface{ isWishlistCommand() }
	WishlistCreateCommand struct{ Ent entity.WishlistEntity }
	WishlistUpdateCommand struct{ Ent entity.WishlistEntity }
	WishlistDeleteCommand struct{ Id int64 }
)

func (WishlistCreateCommand) isWishlistCommand() {}
func (WishlistUpdateCommand) isWishlistCommand() {}
func (WishlistDeleteCommand) isWishlistCommand() {}

type (
	WishlistQuery        interface{ isWishlistQuery() }
	WishlistFindAllQuery struct{ Cond condition.WishlistCondition }
	WishlistFindQuery    struct{ Id int64 }
)

func (WishlistFindAllQuery) isWishlistQuery() {}
func (WishlistFindQuery) isWishlistQuery()    {}
