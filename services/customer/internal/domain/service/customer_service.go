package service

import (
	"context"
	"fmt"
	"time"

	"github.com/codejsha/shared-library-go/pkg/pagination"

	"github.com/codejsha/bookstore-microservices/customer/internal/application/port/protostub"
	"github.com/codejsha/bookstore-microservices/customer/internal/application/port/repo"
	"github.com/codejsha/bookstore-microservices/customer/internal/application/usecase"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/aggregate"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/model/command"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/model/option"
)

var _ usecase.CustomerUseCase = (*customerService)(nil)

const grpcCallTimeout = 5 * time.Second

func withGrpcTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	if dl, ok := ctx.Deadline(); ok && time.Until(dl) < grpcCallTimeout {
		return ctx, func() {}
	}
	return context.WithTimeout(ctx, grpcCallTimeout)
}

type customerService struct {
	pointRepo        repo.PointRepo
	pointHistoryRepo repo.PointHistoryRepo
	reviewRepo       repo.ReviewRepo
	wishlistRepo     repo.WishlistRepo
	userClient       protostub.UserClient
	orderClient      protostub.OrderClient
	paymentClient    protostub.PaymentClient
	deliveryClient   protostub.DeliveryClient
}

func NewCustomerService(
	pointRepo repo.PointRepo,
	pointHistoryRepo repo.PointHistoryRepo,
	reviewRepo repo.ReviewRepo,
	wishlistRepo repo.WishlistRepo,
	userClient protostub.UserClient,
	orderClient protostub.OrderClient,
	paymentClient protostub.PaymentClient,
	deliveryClient protostub.DeliveryClient,
) usecase.CustomerUseCase {
	return &customerService{
		pointRepo:        pointRepo,
		pointHistoryRepo: pointHistoryRepo,
		reviewRepo:       reviewRepo,
		wishlistRepo:     wishlistRepo,
		userClient:       userClient,
		orderClient:      orderClient,
		paymentClient:    paymentClient,
		deliveryClient:   deliveryClient,
	}
}

// ─── Customer profile (via identity gRPC) ───────────────────────────

func (s *customerService) FindAllCustomers(ctx context.Context, opt option.UserQueryOption) (int64, []*aggregate.CustomerAggregate, error) {
	callCtx, cancel := withGrpcTimeout(ctx)
	defer cancel()
	total, out, err := s.userClient.ListUsers(
		callCtx, derefString(opt.Email()), derefString(opt.Name()), derefString(opt.Phone()), opt.Page().GetSize())
	if err != nil {
		return 0, nil, fmt.Errorf("list users via grpc: %w", err)
	}
	return total, out, nil
}

func (s *customerService) FindCustomer(ctx context.Context, uid string) (*aggregate.CustomerAggregate, error) {
	callCtx, cancel := withGrpcTimeout(ctx)
	defer cancel()
	customer, err := s.userClient.FindUser(callCtx, uid)
	if err != nil {
		return nil, fmt.Errorf("find user via grpc: %w", err)
	}
	return customer, nil
}

func (s *customerService) UpdateCustomer(ctx context.Context, cmd command.CustomerUpdateCommand) (*aggregate.CustomerAggregate, error) {
	return nil, fmt.Errorf("update customer: %w (UserService has no UpdateUser RPC)", usecase.ErrNotImplemented)
}

func (s *customerService) DeleteCustomer(ctx context.Context, uid string) error {
	return fmt.Errorf("delete customer: %w (UserService has no DeleteUser RPC)", usecase.ErrNotImplemented)
}

// ─── Orders (via order gRPC) ────────────────────────────────────────

func (s *customerService) FindAllCustomerOrders(ctx context.Context, opt option.OrderQueryOption) (int64, []*aggregate.OrderAggregate, error) {
	callCtx, cancel := withGrpcTimeout(ctx)
	defer cancel()
	total, out, err := s.orderClient.ListOrders(callCtx, derefString(opt.UserUid()), opt.Status(), opt.Page().GetSize())
	if err != nil {
		return 0, nil, fmt.Errorf("list orders via grpc: %w", err)
	}
	return total, out, nil
}

func (s *customerService) FindCustomerOrder(ctx context.Context, uid string, orderUid string) (*aggregate.OrderAggregate, error) {
	callCtx, cancel := withGrpcTimeout(ctx)
	defer cancel()
	o, err := s.orderClient.FindOrder(callCtx, orderUid)
	if err != nil {
		return nil, fmt.Errorf("find order via grpc: %w", err)
	}
	if o == nil {
		return nil, nil
	}
	if o.UserUid != uid {
		return nil, nil
	}
	return o, nil
}

// ─── Payments (via payment gRPC) ────────────────────────────────────

func (s *customerService) FindAllCustomerPayments(ctx context.Context, opt option.PaymentQueryOption) (int64, []*aggregate.PaymentAggregate, error) {
	callCtx, cancel := withGrpcTimeout(ctx)
	defer cancel()
	total, payments, err := s.paymentClient.ListPayments(callCtx, derefString(opt.UserUid()), opt.Page().GetSize())
	if err != nil {
		return 0, nil, fmt.Errorf("list payments via grpc: %w", err)
	}
	out := make([]*aggregate.PaymentAggregate, 0, len(payments))
	for _, p := range payments {
		if orderUid := opt.OrderUid(); orderUid != nil && *orderUid != "" {
			_ = orderUid
		}
		out = append(out, p)
	}
	return total, out, nil
}

func (s *customerService) FindCustomerPayment(ctx context.Context, uid string, paymentUid string) (*aggregate.PaymentAggregate, error) {
	callCtx, cancel := withGrpcTimeout(ctx)
	defer cancel()
	p, err := s.paymentClient.FindPayment(callCtx, paymentUid)
	if err != nil {
		return nil, fmt.Errorf("find payment via grpc: %w", err)
	}
	if p == nil {
		return nil, nil
	}
	if p.CustomerUid != uid {
		return nil, nil
	}
	return p, nil
}

// ─── Shipments (via delivery gRPC) ──────────────────────────────────

func (s *customerService) ownsOrder(ctx context.Context, orderUid string, caller usecase.Caller) (bool, error) {
	if caller.IsAdmin {
		return true, nil
	}
	callCtx, cancel := withGrpcTimeout(ctx)
	defer cancel()
	o, err := s.orderClient.FindOrder(callCtx, orderUid)
	if err != nil {
		return false, fmt.Errorf("resolve order owner via grpc: %w", err)
	}
	if o == nil {
		return false, nil
	}
	return caller.Owns(o.UserUid), nil
}

func (s *customerService) TrackShipment(ctx context.Context, uid string, caller usecase.Caller) (*aggregate.ShipmentAggregate, error) {
	callCtx, cancel := withGrpcTimeout(ctx)
	defer cancel()
	sh, err := s.deliveryClient.TrackShipment(callCtx, uid, caller.UserUid)
	if err != nil {
		return nil, fmt.Errorf("track shipment via grpc: %w", err)
	}
	if sh == nil {
		return nil, nil
	}

	owns, err := s.ownsOrder(ctx, sh.OrderUid, caller)
	if err != nil {
		return nil, err
	}
	if !owns {
		return nil, nil
	}
	return sh, nil
}

func (s *customerService) ListOrderShipments(ctx context.Context, orderUid string, status *string, size *int32, caller usecase.Caller) (int64, []*aggregate.ShipmentAggregate, error) {
	owns, err := s.ownsOrder(ctx, orderUid, caller)
	if err != nil {
		return 0, nil, err
	}
	if !owns {
		return 0, nil, nil
	}

	callCtx, cancel := withGrpcTimeout(ctx)
	defer cancel()
	total, out, err := s.deliveryClient.ListShipments(callCtx, orderUid, status, size, caller.UserUid)
	if err != nil {
		return 0, nil, fmt.Errorf("list shipments via grpc: %w", err)
	}
	return total, out, nil
}

// ─── Wishlist ───────────────────────────────────────────────────────

func (s *customerService) GetWishlist(ctx context.Context, userUid string) (*aggregate.WishlistAggregate, error) {
	opt := option.NewWishlistQueryOption(option.WishlistQueryOption{}.WithUserUid(&userUid))
	_, results, err := s.wishlistRepo.FindAll(ctx, opt)
	if err != nil {
		return nil, err
	}
	return s.toWishlistAggregate(userUid, results), nil
}

func (s *customerService) AddBooksToWishlist(ctx context.Context, userUid string, bookUids []string) (*aggregate.WishlistAggregate, error) {
	if err := s.wishlistRepo.AddBooks(ctx, userUid, bookUids); err != nil {
		return nil, fmt.Errorf("add books to wishlist: %w", err)
	}
	return s.GetWishlist(ctx, userUid)
}

func (s *customerService) RemoveBooksFromWishlist(ctx context.Context, userUid string, bookUids []string) (*aggregate.WishlistAggregate, error) {
	rows, err := s.wishlistRepo.FetchByUserUid(ctx, userUid)
	if err != nil {
		return nil, err
	}
	wanted := make(map[string]struct{}, len(bookUids))
	for _, b := range bookUids {
		wanted[b] = struct{}{}
	}
	for _, row := range rows {
		if _, hit := wanted[row.BookUid]; !hit {
			continue
		}
		if err := s.wishlistRepo.Delete(ctx, row); err != nil {
			return nil, err
		}
	}
	return s.GetWishlist(ctx, userUid)
}

// ─── Points ─────────────────────────────────────────────────────────

func (s *customerService) GetPointBalance(ctx context.Context, userUid string) (*aggregate.PointAggregate, error) {
	result, err := s.pointRepo.FindByUserUid(ctx, userUid)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return &aggregate.PointAggregate{UserUid: userUid, Balance: 0}, nil
	}
	return s.toPointAggregate(result), nil
}

func (s *customerService) EarnPoints(ctx context.Context, cmd command.PointEarnCommand) (*aggregate.PointAggregate, error) {
	if cmd.Amount <= 0 {
		return nil, fmt.Errorf("earn amount must be positive: got %d", cmd.Amount)
	}
	if err := s.pointRepo.EnsureBalance(ctx, cmd.UserUid); err != nil {
		return nil, fmt.Errorf("ensure point balance: %w", err)
	}
	return s.changePoints(ctx, cmd.UserUid, cmd.Amount, "EARN", cmd.Reason)
}

func (s *customerService) SpendPoints(ctx context.Context, cmd command.PointSpendCommand) (*aggregate.PointAggregate, error) {
	if cmd.Amount <= 0 {
		return nil, fmt.Errorf("spend amount must be positive: got %d", cmd.Amount)
	}
	return s.changePoints(ctx, cmd.UserUid, -cmd.Amount, "SPEND", cmd.Reason)
}

func (s *customerService) GetPointHistory(ctx context.Context, userUid string, page pagination.PageOption) (int64, []*aggregate.PointHistoryEntry, error) {
	total, results, err := s.pointHistoryRepo.FindAll(ctx, userUid, page)
	if err != nil {
		return 0, nil, err
	}
	out := make([]*aggregate.PointHistoryEntry, len(results))
	for i, r := range results {
		out[i] = toPointHistoryEntry(r)
	}
	return total, out, nil
}

// ─── Reviews ────────────────────────────────────────────────────────

func (s *customerService) GetCustomerReviews(ctx context.Context, userUid string, opt option.ReviewQueryOption) (int64, []*aggregate.ReviewAggregate, error) {
	merged := option.NewReviewQueryOption(
		option.WithReviewUserUid(userUid),
		option.WithReviewPage(opt.Page()),
	)
	if v := opt.Rating(); v != nil {
		merged = option.NewReviewQueryOption(
			option.WithReviewUserUid(userUid),
			option.WithReviewRating(*v),
			option.WithReviewPage(opt.Page()),
		)
	}
	total, results, err := s.reviewRepo.FindAll(ctx, merged)
	if err != nil {
		return 0, nil, err
	}
	out := make([]*aggregate.ReviewAggregate, len(results))
	for i, r := range results {
		out[i] = toReviewAggregate(r)
	}
	return total, out, nil
}

func (s *customerService) GetBookReviews(ctx context.Context, bookUid string, opt option.ReviewQueryOption) (int64, []*aggregate.ReviewAggregate, error) {
	merged := option.NewReviewQueryOption(
		option.WithReviewBookUid(bookUid),
		option.WithReviewPage(opt.Page()),
	)
	if v := opt.Rating(); v != nil {
		merged = option.NewReviewQueryOption(
			option.WithReviewBookUid(bookUid),
			option.WithReviewRating(*v),
			option.WithReviewPage(opt.Page()),
		)
	}
	total, results, err := s.reviewRepo.FindAll(ctx, merged)
	if err != nil {
		return 0, nil, err
	}
	out := make([]*aggregate.ReviewAggregate, len(results))
	for i, r := range results {
		out[i] = toReviewAggregate(r)
	}
	return total, out, nil
}

func (s *customerService) GetReview(ctx context.Context, userUid string, reviewUid string) (*aggregate.ReviewAggregate, error) {
	r, err := s.reviewRepo.FindByUid(ctx, reviewUid)
	if err != nil {
		return nil, err
	}
	if r == nil {
		return nil, nil
	}
	if r.UserUid != userUid {
		return nil, nil
	}
	return toReviewAggregate(r), nil
}

func (s *customerService) WriteReview(ctx context.Context, cmd command.ReviewWriteCommand) (*aggregate.ReviewAggregate, error) {
	result, err := s.reviewRepo.Insert(ctx, repo.ReviewCreate{
		UserUid: cmd.UserUid,
		BookUid: cmd.BookUid,
		Rating:  cmd.Rating,
		Title:   cmd.Title,
		Content: cmd.Content,
	})
	if err != nil {
		return nil, fmt.Errorf("write review: %w", err)
	}
	return toReviewAggregate(result), nil
}

func (s *customerService) EditReview(ctx context.Context, userUid string, reviewUid string, cmd command.ReviewEditCommand) (*aggregate.ReviewAggregate, error) {
	rows, err := s.reviewRepo.FetchByUid(ctx, reviewUid)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, nil
	}
	ent := rows[0]
	if ent.UserUid != userUid {
		return nil, nil
	}
	if cmd.Rating != nil {
		ent.Rating = int16(*cmd.Rating)
	}
	if cmd.Title != nil {
		ent.Title = cmd.Title
	}
	if cmd.Content != nil {
		ent.Content = cmd.Content
	}
	now := time.Now()
	ent.UpdatedAt = &now
	if err := s.reviewRepo.Update(ctx, ent); err != nil {
		return nil, err
	}
	updated, err := s.reviewRepo.FindByUid(ctx, reviewUid)
	if err != nil {
		return nil, err
	}
	return toReviewAggregate(updated), nil
}

func (s *customerService) RemoveReview(ctx context.Context, userUid string, reviewUid string) error {
	rows, err := s.reviewRepo.FetchByUid(ctx, reviewUid)
	if err != nil {
		return err
	}
	if len(rows) == 0 {
		return nil
	}
	ent := rows[0]
	if ent.UserUid != userUid {
		return nil
	}
	return s.reviewRepo.Delete(ctx, ent)
}

// ─── points helper ──────────────────────────────────────────────────

func (s *customerService) changePoints(ctx context.Context, userUid string, delta int32, changeType string, reason *string) (*aggregate.PointAggregate, error) {
	result, err := s.pointRepo.ApplyPointChange(ctx, userUid, delta, changeType, reason)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, fmt.Errorf("change points: %w (no balance for user)", repo.ErrInsufficientPoints)
	}
	return s.toPointAggregate(result), nil
}

// ─── mapping helpers ────────────────────────────────────────────────

func (s *customerService) toWishlistAggregate(userUid string, results []*repo.WishlistResult) *aggregate.WishlistAggregate {
	bookUids := make([]string, len(results))
	for i, r := range results {
		bookUids[i] = r.BookUid
	}
	return &aggregate.WishlistAggregate{
		UserUid:  userUid,
		BookUids: bookUids,
	}
}

func (s *customerService) toPointAggregate(r *repo.PointResult) *aggregate.PointAggregate {
	return &aggregate.PointAggregate{
		Uid:     r.Uid,
		UserUid: r.UserUid,
		Balance: r.Balance,
	}
}

func toReviewAggregate(r *repo.ReviewResult) *aggregate.ReviewAggregate {
	if r == nil {
		return nil
	}
	return &aggregate.ReviewAggregate{
		Uid:       r.Uid,
		UserUid:   r.UserUid,
		BookUid:   r.BookUid,
		Rating:    r.Rating,
		Title:     r.Title,
		Content:   r.Content,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}

func toPointHistoryEntry(r *repo.PointHistoryResult) *aggregate.PointHistoryEntry {
	return &aggregate.PointHistoryEntry{
		Uid:        r.Uid,
		ChangeType: aggregate.PointChangeType(r.ChangeType),
		Amount:     r.Amount,
		Reason:     r.Reason,
		CreatedAt:  r.CreatedAt,
	}
}

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
