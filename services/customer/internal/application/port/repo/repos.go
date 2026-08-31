package repo

import (
	"context"
	"errors"

	"github.com/codejsha/shared-library-go/pkg/pagination"

	genrepo "github.com/codejsha/bookstore-microservices/customer/generated/infrastructure/port/repo"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/model/option"
)

var ErrInsufficientPoints = errors.New("insufficient points")

var ErrPointBalanceOverflow = errors.New("point balance limit exceeded")

var ErrReviewExists = errors.New("review already exists for this book")

type PointRepo interface {
	FindOne(ctx context.Context, id int64) (*PointResult, error)
	FindByUserUid(ctx context.Context, userUid string) (*PointResult, error)
	EnsureBalance(ctx context.Context, userUid string) error
	ApplyPointChange(ctx context.Context, userUid string, delta int32, changeType string, reason *string) (*PointResult, error)
	genrepo.CustomerPointRepo
}

type PointHistoryRepo interface {
	FindAll(ctx context.Context, userUid string, page pagination.PageOption) (int64, []*PointHistoryResult, error)
	genrepo.CustomerPointHistoryRepo
}

type ReviewRepo interface {
	FindAll(ctx context.Context, opt option.ReviewQueryOption) (int64, []*ReviewResult, error)
	FindOne(ctx context.Context, id int64) (*ReviewResult, error)
	FindByUid(ctx context.Context, uid string) (*ReviewResult, error)
	Insert(ctx context.Context, review ReviewCreate) (*ReviewResult, error)
	genrepo.CustomerReviewRepo
}

type WishlistRepo interface {
	FindAll(ctx context.Context, opt option.WishlistQueryOption) (int64, []*WishlistResult, error)
	FindOne(ctx context.Context, id int64) (*WishlistResult, error)
	AddBooks(ctx context.Context, userUid string, bookUids []string) error
	genrepo.CustomerWishlistRepo
}
