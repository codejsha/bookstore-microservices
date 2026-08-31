package usecase

import (
	"context"
	"errors"

	"github.com/codejsha/shared-library-go/pkg/pagination"

	"github.com/codejsha/bookstore-microservices/customer/internal/domain/aggregate"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/model/command"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/model/option"
)

var ErrNotImplemented = errors.New("operation not supported by upstream service")

type Caller struct {
	UserUid string
	IsAdmin bool
}

func (c Caller) Owns(ownerUid string) bool {
	return c.IsAdmin || (c.UserUid != "" && c.UserUid == ownerUid)
}

type CustomerUseCase interface {
	// ─── Customer profile ──────────────
	FindAllCustomers(ctx context.Context, opt option.UserQueryOption) (int64, []*aggregate.CustomerAggregate, error)
	FindCustomer(ctx context.Context, uid string) (*aggregate.CustomerAggregate, error)
	UpdateCustomer(ctx context.Context, cmd command.CustomerUpdateCommand) (*aggregate.CustomerAggregate, error)
	DeleteCustomer(ctx context.Context, uid string) error

	// ─── Orders ─────────────────────────
	FindAllCustomerOrders(ctx context.Context, opt option.OrderQueryOption) (int64, []*aggregate.OrderAggregate, error)
	FindCustomerOrder(ctx context.Context, uid string, orderUid string) (*aggregate.OrderAggregate, error)

	// ─── Payments ───────────────────────
	FindAllCustomerPayments(ctx context.Context, opt option.PaymentQueryOption) (int64, []*aggregate.PaymentAggregate, error)
	FindCustomerPayment(ctx context.Context, uid string, paymentUid string) (*aggregate.PaymentAggregate, error)

	// ─── Shipments ──────────────────────
	TrackShipment(ctx context.Context, uid string, caller Caller) (*aggregate.ShipmentAggregate, error)
	ListOrderShipments(ctx context.Context, orderUid string, status *string, size *int32, caller Caller) (int64, []*aggregate.ShipmentAggregate, error)

	// ─── Wishlist ───────────────────────
	GetWishlist(ctx context.Context, userUid string) (*aggregate.WishlistAggregate, error)
	AddBooksToWishlist(ctx context.Context, cmd command.WishlistAddCommand) (*aggregate.WishlistAggregate, error)
	RemoveBooksFromWishlist(ctx context.Context, cmd command.WishlistRemoveCommand) (*aggregate.WishlistAggregate, error)

	// ─── Points ─────────────────────────
	GetPointBalance(ctx context.Context, userUid string) (*aggregate.PointAggregate, error)
	EarnPoints(ctx context.Context, cmd command.PointEarnCommand) (*aggregate.PointAggregate, error)
	SpendPoints(ctx context.Context, cmd command.PointSpendCommand) (*aggregate.PointAggregate, error)
	GetPointHistory(ctx context.Context, userUid string, page pagination.PageOption) (int64, []*aggregate.PointHistoryEntry, error)

	// ─── Reviews ───────────────────────
	GetCustomerReviews(ctx context.Context, userUid string, opt option.ReviewQueryOption) (int64, []*aggregate.ReviewAggregate, error)
	GetBookReviews(ctx context.Context, bookUid string, opt option.ReviewQueryOption) (int64, []*aggregate.ReviewAggregate, error)
	GetReview(ctx context.Context, userUid string, reviewUid string) (*aggregate.ReviewAggregate, error)
	WriteReview(ctx context.Context, cmd command.ReviewWriteCommand) (*aggregate.ReviewAggregate, error)
	EditReview(ctx context.Context, userUid string, reviewUid string, cmd command.ReviewEditCommand) (*aggregate.ReviewAggregate, error)
	RemoveReview(ctx context.Context, userUid string, reviewUid string) error
}
