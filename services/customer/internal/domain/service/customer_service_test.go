package service

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/codejsha/shared-library-go/pkg/pagination"

	"github.com/codejsha/bookstore-microservices/customer/generated/infrastructure/port/entity"
	"github.com/codejsha/bookstore-microservices/customer/internal/application/port/protostub"
	"github.com/codejsha/bookstore-microservices/customer/internal/application/port/repo"
	"github.com/codejsha/bookstore-microservices/customer/internal/application/usecase"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/aggregate"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/model/command"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/model/option"
)

// ─── client stubs ───────────────────────────────────────────────────────────

var (
	_ protostub.UserClient     = (*stubUserClient)(nil)
	_ protostub.OrderClient    = (*stubOrderClient)(nil)
	_ protostub.PaymentClient  = (*stubPaymentClient)(nil)
	_ protostub.DeliveryClient = (*stubDeliveryClient)(nil)
)

type stubUserClient struct {
	listUsersFn func(ctx context.Context, email, name, phone string, pageSize int32) (int64, []*aggregate.CustomerAggregate, error)
	findUserFn  func(ctx context.Context, uid string) (*aggregate.CustomerAggregate, error)
}

func (s *stubUserClient) ListUsers(ctx context.Context, email, name, phone string, pageSize int32) (int64, []*aggregate.CustomerAggregate, error) {
	return s.listUsersFn(ctx, email, name, phone, pageSize)
}

func (s *stubUserClient) FindUser(ctx context.Context, uid string) (*aggregate.CustomerAggregate, error) {
	return s.findUserFn(ctx, uid)
}

type stubOrderClient struct {
	listOrdersFn func(ctx context.Context, userUid string, status *string, pageSize int32) (int64, []*aggregate.OrderAggregate, error)
	findOrderFn  func(ctx context.Context, uid string) (*aggregate.OrderAggregate, error)
}

func (s *stubOrderClient) ListOrders(ctx context.Context, userUid string, status *string, pageSize int32) (int64, []*aggregate.OrderAggregate, error) {
	return s.listOrdersFn(ctx, userUid, status, pageSize)
}

func (s *stubOrderClient) FindOrder(ctx context.Context, uid string) (*aggregate.OrderAggregate, error) {
	return s.findOrderFn(ctx, uid)
}

type stubPaymentClient struct {
	listPaymentsFn func(ctx context.Context, customerUid string, pageSize int32) (int64, []*aggregate.PaymentAggregate, error)
	findPaymentFn  func(ctx context.Context, uid string) (*aggregate.PaymentAggregate, error)
}

func (s *stubPaymentClient) ListPayments(ctx context.Context, customerUid string, pageSize int32) (int64, []*aggregate.PaymentAggregate, error) {
	return s.listPaymentsFn(ctx, customerUid, pageSize)
}

func (s *stubPaymentClient) FindPayment(ctx context.Context, uid string) (*aggregate.PaymentAggregate, error) {
	return s.findPaymentFn(ctx, uid)
}

type stubDeliveryClient struct {
	trackShipmentFn func(ctx context.Context, uid string, callerUid string) (*aggregate.ShipmentAggregate, error)
	listShipmentsFn func(ctx context.Context, orderUid string, status *string, pageSize *int32, callerUid string) (int64, []*aggregate.ShipmentAggregate, error)
}

func (s *stubDeliveryClient) TrackShipment(ctx context.Context, uid string, callerUid string) (*aggregate.ShipmentAggregate, error) {
	return s.trackShipmentFn(ctx, uid, callerUid)
}

func (s *stubDeliveryClient) ListShipments(ctx context.Context, orderUid string, status *string, pageSize *int32, callerUid string) (int64, []*aggregate.ShipmentAggregate, error) {
	return s.listShipmentsFn(ctx, orderUid, status, pageSize, callerUid)
}

// ─── repository stubs ─────────────────────────────────────────────────────────

type stubPointRepo struct {
	repo.PointRepo
	findByUserUidFn    func(ctx context.Context, userUid string) (*repo.PointResult, error)
	ensureBalanceFn    func(ctx context.Context, userUid string) error
	applyPointChangeFn func(ctx context.Context, userUid string, delta int32, changeType string, reason *string) (*repo.PointResult, error)
}

func (s *stubPointRepo) FindByUserUid(ctx context.Context, userUid string) (*repo.PointResult, error) {
	return s.findByUserUidFn(ctx, userUid)
}

func (s *stubPointRepo) EnsureBalance(ctx context.Context, userUid string) error {
	if s.ensureBalanceFn == nil {
		return nil
	}
	return s.ensureBalanceFn(ctx, userUid)
}

func (s *stubPointRepo) ApplyPointChange(
	ctx context.Context, userUid string, delta int32, changeType string, reason *string,
) (*repo.PointResult, error) {
	return s.applyPointChangeFn(ctx, userUid, delta, changeType, reason)
}

type stubPointHistoryRepo struct {
	repo.PointHistoryRepo
	findAllFn func(ctx context.Context, userUid string, page pagination.PageOption) (int64, []*repo.PointHistoryResult, error)
	createFn  func(ctx context.Context, ent *entity.CustomerPointHistoryEntity) error
}

func (s *stubPointHistoryRepo) FindAll(
	ctx context.Context, userUid string, page pagination.PageOption,
) (int64, []*repo.PointHistoryResult, error) {
	return s.findAllFn(ctx, userUid, page)
}

func (s *stubPointHistoryRepo) Create(ctx context.Context, ent *entity.CustomerPointHistoryEntity) error {
	return s.createFn(ctx, ent)
}

type stubWishlistRepo struct {
	repo.WishlistRepo
	findAllFn        func(ctx context.Context, opt option.WishlistQueryOption) (int64, []*repo.WishlistResult, error)
	fetchByUserUidFn func(ctx context.Context, userUid string) ([]*entity.CustomerWishlistEntity, error)
	deleteFn         func(ctx context.Context, ent *entity.CustomerWishlistEntity) error
	addBooksFn       func(ctx context.Context, userUid string, bookUids []string) error
}

func (s *stubWishlistRepo) AddBooks(ctx context.Context, userUid string, bookUids []string) error {
	return s.addBooksFn(ctx, userUid, bookUids)
}

func (s *stubWishlistRepo) FindAll(
	ctx context.Context, opt option.WishlistQueryOption,
) (int64, []*repo.WishlistResult, error) {
	return s.findAllFn(ctx, opt)
}

func (s *stubWishlistRepo) FetchByUserUid(
	ctx context.Context, userUid string,
) ([]*entity.CustomerWishlistEntity, error) {
	return s.fetchByUserUidFn(ctx, userUid)
}

func (s *stubWishlistRepo) Delete(ctx context.Context, ent *entity.CustomerWishlistEntity) error {
	return s.deleteFn(ctx, ent)
}

type stubReviewRepo struct {
	repo.ReviewRepo
	findAllFn    func(ctx context.Context, opt option.ReviewQueryOption) (int64, []*repo.ReviewResult, error)
	findByUidFn  func(ctx context.Context, uid string) (*repo.ReviewResult, error)
	fetchByUidFn func(ctx context.Context, uid string) ([]*entity.CustomerReviewEntity, error)
	updateFn     func(ctx context.Context, ent *entity.CustomerReviewEntity) error
	deleteFn     func(ctx context.Context, ent *entity.CustomerReviewEntity) error
	insertFn     func(ctx context.Context, review repo.ReviewCreate) (*repo.ReviewResult, error)
}

func (s *stubReviewRepo) Insert(ctx context.Context, review repo.ReviewCreate) (*repo.ReviewResult, error) {
	return s.insertFn(ctx, review)
}

func (s *stubReviewRepo) FindAll(
	ctx context.Context, opt option.ReviewQueryOption,
) (int64, []*repo.ReviewResult, error) {
	return s.findAllFn(ctx, opt)
}

func (s *stubReviewRepo) FindByUid(ctx context.Context, uid string) (*repo.ReviewResult, error) {
	return s.findByUidFn(ctx, uid)
}

func (s *stubReviewRepo) FetchByUid(ctx context.Context, uid string) ([]*entity.CustomerReviewEntity, error) {
	return s.fetchByUidFn(ctx, uid)
}

func (s *stubReviewRepo) Update(ctx context.Context, ent *entity.CustomerReviewEntity) error {
	return s.updateFn(ctx, ent)
}

func (s *stubReviewRepo) Delete(ctx context.Context, ent *entity.CustomerReviewEntity) error {
	return s.deleteFn(ctx, ent)
}

// ─── Customer profile (identity gRPC) ───────────────────────────────────────

func TestFindAllCustomers(t *testing.T) {
	userClient := &stubUserClient{
		listUsersFn: func(context.Context, string, string, string, int32) (int64, []*aggregate.CustomerAggregate, error) {
			return 2, []*aggregate.CustomerAggregate{
				{Uid: "u-1", Email: "a@b.com"},
				{Uid: "u-2", Email: "c@d.com"},
			}, nil
		},
	}
	svc := &customerService{userClient: userClient}
	total, aggs, err := svc.FindAllCustomers(context.Background(), option.NewUserQueryOption())
	if err != nil {
		t.Fatalf("FindAllCustomers err: %v", err)
	}
	if total != 2 || len(aggs) != 2 {
		t.Fatalf("total=%d len=%d, want 2/2", total, len(aggs))
	}
	if aggs[0].Uid != "u-1" || aggs[1].Uid != "u-2" {
		t.Errorf("aggs = %+v", aggs)
	}
}

func TestFindAllCustomersRepoError(t *testing.T) {
	wantErr := errors.New("identity down")
	svc := &customerService{userClient: &stubUserClient{
		listUsersFn: func(context.Context, string, string, string, int32) (int64, []*aggregate.CustomerAggregate, error) {
			return 0, nil, wantErr
		},
	}}
	if _, _, err := svc.FindAllCustomers(context.Background(), option.NewUserQueryOption()); !errors.Is(err, wantErr) {
		t.Errorf("err = %v, want %v", err, wantErr)
	}
}

func TestFindCustomer(t *testing.T) {
	svc := &customerService{userClient: &stubUserClient{
		findUserFn: func(_ context.Context, uid string) (*aggregate.CustomerAggregate, error) {
			if uid != "u-1" {
				t.Errorf("uid = %q, want u-1", uid)
			}
			return &aggregate.CustomerAggregate{Uid: "u-1", Email: "a@b.com"}, nil
		},
	}}
	got, err := svc.FindCustomer(context.Background(), "u-1")
	if err != nil {
		t.Fatalf("FindCustomer err: %v", err)
	}
	if got == nil || got.Uid != "u-1" {
		t.Errorf("got = %+v, want uid=u-1", got)
	}
}

func TestFindCustomerNotFound(t *testing.T) {
	svc := &customerService{userClient: &stubUserClient{
		findUserFn: func(context.Context, string) (*aggregate.CustomerAggregate, error) {
			return nil, nil
		},
	}}
	got, err := svc.FindCustomer(context.Background(), "u-x")
	if err != nil {
		t.Fatalf("err = %v, want nil", err)
	}
	if got != nil {
		t.Errorf("got = %+v, want nil", got)
	}
}

func TestUpdateAndDeleteCustomerNotImplemented(t *testing.T) {
	svc := &customerService{}
	if _, err := svc.UpdateCustomer(context.Background(), command.CustomerUpdateCommand{Uid: "u-1"}); !errors.Is(err, usecase.ErrNotImplemented) {
		t.Errorf("UpdateCustomer err = %v, want ErrNotImplemented", err)
	}
	if err := svc.DeleteCustomer(context.Background(), "u-1"); !errors.Is(err, usecase.ErrNotImplemented) {
		t.Errorf("DeleteCustomer err = %v, want ErrNotImplemented", err)
	}
}

// ─── Orders / Payments (gRPC) ───────────────────────────────────────────────

func TestFindAllCustomerOrders(t *testing.T) {
	svc := &customerService{orderClient: &stubOrderClient{
		listOrdersFn: func(context.Context, string, *string, int32) (int64, []*aggregate.OrderAggregate, error) {
			return 1, []*aggregate.OrderAggregate{{Uid: "o-1", UserUid: "u-1"}}, nil
		},
	}}
	total, aggs, err := svc.FindAllCustomerOrders(context.Background(), option.NewOrderQueryOption())
	if err != nil {
		t.Fatalf("FindAllCustomerOrders err: %v", err)
	}
	if total != 1 || len(aggs) != 1 || aggs[0].Uid != "o-1" {
		t.Errorf("got total=%d aggs=%+v", total, aggs)
	}
}

func TestFindCustomerOrder(t *testing.T) {
	t.Run("owned by caller", func(t *testing.T) {
		svc := &customerService{orderClient: &stubOrderClient{
			findOrderFn: func(_ context.Context, uid string) (*aggregate.OrderAggregate, error) {
				if uid != "o-1" {
					t.Errorf("order uid = %q, want o-1", uid)
				}
				return &aggregate.OrderAggregate{Uid: "o-1", UserUid: "u-1"}, nil
			},
		}}
		got, err := svc.FindCustomerOrder(context.Background(), "u-1", "o-1")
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if got == nil || got.Uid != "o-1" {
			t.Errorf("got = %+v, want uid=o-1", got)
		}
	})

	t.Run("owned by another customer is hidden", func(t *testing.T) {
		svc := &customerService{orderClient: &stubOrderClient{
			findOrderFn: func(context.Context, string) (*aggregate.OrderAggregate, error) {
				return &aggregate.OrderAggregate{Uid: "o-1", UserUid: "someone-else"}, nil
			},
		}}
		got, err := svc.FindCustomerOrder(context.Background(), "u-1", "o-1")
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if got != nil {
			t.Errorf("got = %+v, want nil (other customer's order must not leak)", got)
		}
	})
}

func TestFindAllCustomerPayments(t *testing.T) {
	svc := &customerService{paymentClient: &stubPaymentClient{
		listPaymentsFn: func(context.Context, string, int32) (int64, []*aggregate.PaymentAggregate, error) {
			return 1, []*aggregate.PaymentAggregate{{Uid: "p-1", CustomerUid: "u-1"}}, nil
		},
	}}
	total, aggs, err := svc.FindAllCustomerPayments(context.Background(), option.NewPaymentQueryOption())
	if err != nil {
		t.Fatalf("FindAllCustomerPayments err: %v", err)
	}
	if total != 1 || len(aggs) != 1 || aggs[0].Uid != "p-1" {
		t.Errorf("got total=%d aggs=%+v", total, aggs)
	}
}

func TestFindCustomerPaymentOwnershipCheck(t *testing.T) {
	svc := &customerService{paymentClient: &stubPaymentClient{
		findPaymentFn: func(context.Context, string) (*aggregate.PaymentAggregate, error) {
			return &aggregate.PaymentAggregate{Uid: "p-1", CustomerUid: "other"}, nil
		},
	}}
	got, err := svc.FindCustomerPayment(context.Background(), "u-1", "p-1")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if got != nil {
		t.Errorf("got = %+v, want nil (other customer's payment must not leak)", got)
	}
}

// ─── Wishlist ───────────────────────────────────────────────────────────────

func TestGetWishlist(t *testing.T) {
	svc := &customerService{wishlistRepo: &stubWishlistRepo{
		findAllFn: func(_ context.Context, _ option.WishlistQueryOption) (int64, []*repo.WishlistResult, error) {
			return 2, []*repo.WishlistResult{
				{Uid: "w-1", BookUid: "b-1"},
				{Uid: "w-2", BookUid: "b-2"},
			}, nil
		},
	}}
	got, err := svc.GetWishlist(context.Background(), "u-1")
	if err != nil {
		t.Fatalf("GetWishlist err: %v", err)
	}
	if got.UserUid != "u-1" || !reflect.DeepEqual(got.BookUids, []string{"b-1", "b-2"}) {
		t.Errorf("got = %+v", got)
	}
}

func TestAddBooksToWishlist(t *testing.T) {
	var gotUser string
	var gotBooks []string
	wishlistRepo := &stubWishlistRepo{
		addBooksFn: func(_ context.Context, userUid string, bookUids []string) error {
			gotUser = userUid
			gotBooks = bookUids
			return nil
		},
		findAllFn: func(context.Context, option.WishlistQueryOption) (int64, []*repo.WishlistResult, error) {
			return 2, []*repo.WishlistResult{
				{Uid: "w-1", BookUid: "b-1"},
				{Uid: "w-2", BookUid: "b-2"},
			}, nil
		},
	}
	svc := &customerService{wishlistRepo: wishlistRepo}
	got, err := svc.AddBooksToWishlist(context.Background(), "u-1", []string{"b-1", "b-2"})
	if err != nil {
		t.Fatalf("AddBooksToWishlist err: %v", err)
	}
	if gotUser != "u-1" || !reflect.DeepEqual(gotBooks, []string{"b-1", "b-2"}) {
		t.Errorf("AddBooks called with user=%q books=%v", gotUser, gotBooks)
	}
	if got == nil || len(got.BookUids) != 2 {
		t.Errorf("got = %+v, want refreshed wishlist with 2 books", got)
	}
}

func TestAddBooksToWishlistRepoError(t *testing.T) {
	wantErr := errors.New("insert failed")
	svc := &customerService{wishlistRepo: &stubWishlistRepo{
		addBooksFn: func(context.Context, string, []string) error { return wantErr },
	}}
	if _, err := svc.AddBooksToWishlist(context.Background(), "u-1", []string{"b-1"}); !errors.Is(err, wantErr) {
		t.Errorf("err = %v, want wrapped %v", err, wantErr)
	}
}

func TestRemoveBooksFromWishlist(t *testing.T) {
	rows := []*entity.CustomerWishlistEntity{
		{Id: 1, Uid: "w-1", UserUid: "u-1", BookUid: "b-1"},
		{Id: 2, Uid: "w-2", UserUid: "u-1", BookUid: "b-2"},
		{Id: 3, Uid: "w-3", UserUid: "u-1", BookUid: "b-3"},
	}
	var deleted []string
	wishlistRepo := &stubWishlistRepo{
		fetchByUserUidFn: func(context.Context, string) ([]*entity.CustomerWishlistEntity, error) {
			return rows, nil
		},
		deleteFn: func(_ context.Context, ent *entity.CustomerWishlistEntity) error {
			deleted = append(deleted, ent.BookUid)
			return nil
		},
		findAllFn: func(context.Context, option.WishlistQueryOption) (int64, []*repo.WishlistResult, error) {
			return 1, []*repo.WishlistResult{{Uid: "w-2", BookUid: "b-2"}}, nil
		},
	}
	svc := &customerService{wishlistRepo: wishlistRepo}
	got, err := svc.RemoveBooksFromWishlist(context.Background(), "u-1", []string{"b-1", "b-3"})
	if err != nil {
		t.Fatalf("RemoveBooksFromWishlist err: %v", err)
	}
	if !reflect.DeepEqual(deleted, []string{"b-1", "b-3"}) {
		t.Errorf("deleted = %v, want only b-1 and b-3", deleted)
	}
	if got == nil || len(got.BookUids) != 1 {
		t.Errorf("got = %+v, want refreshed wishlist", got)
	}
}

// ─── Points ─────────────────────────────────────────────────────────────────

func TestGetPointBalance(t *testing.T) {
	t.Run("existing balance", func(t *testing.T) {
		svc := &customerService{pointRepo: &stubPointRepo{
			findByUserUidFn: func(context.Context, string) (*repo.PointResult, error) {
				return &repo.PointResult{Uid: "pt-1", UserUid: "u-1", Balance: 750}, nil
			},
		}}
		got, err := svc.GetPointBalance(context.Background(), "u-1")
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if got == nil || got.Balance != 750 {
			t.Errorf("got = %+v, want balance=750", got)
		}
	})

	t.Run("no point row is a zero balance", func(t *testing.T) {
		svc := &customerService{pointRepo: &stubPointRepo{
			findByUserUidFn: func(context.Context, string) (*repo.PointResult, error) { return nil, nil },
		}}
		got, err := svc.GetPointBalance(context.Background(), "u-1")
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if got == nil || got.Balance != 0 || got.UserUid != "u-1" {
			t.Errorf("got = %+v, want zero-balance aggregate for u-1", got)
		}
	})
}

func TestEarnPoints(t *testing.T) {
	var gotDelta int32
	var gotType string
	ensureCalled := false
	pointRepo := &stubPointRepo{
		ensureBalanceFn: func(_ context.Context, userUid string) error {
			ensureCalled = true
			if userUid != "u-1" {
				t.Errorf("EnsureBalance userUid = %q, want u-1", userUid)
			}
			return nil
		},
		applyPointChangeFn: func(_ context.Context, userUid string, delta int32, changeType string, _ *string) (*repo.PointResult, error) {
			gotDelta = delta
			gotType = changeType
			return &repo.PointResult{Uid: "pt-1", UserId: 1, UserUid: userUid, Balance: 100 + delta}, nil
		},
	}
	svc := &customerService{pointRepo: pointRepo}
	got, err := svc.EarnPoints(context.Background(), command.PointEarnCommand{UserUid: "u-1", Amount: 50})
	if err != nil {
		t.Fatalf("EarnPoints err: %v", err)
	}
	if !ensureCalled {
		t.Error("EnsureBalance was not called before applying the earn")
	}
	if gotDelta != 50 || gotType != "EARN" {
		t.Errorf("ApplyPointChange called with delta=%d type=%s, want 50/EARN", gotDelta, gotType)
	}
	if got == nil || got.Balance != 150 {
		t.Errorf("got = %+v, want balance=150", got)
	}
}

func TestEarnPointsRejectsNonPositive(t *testing.T) {
	svc := &customerService{}
	if _, err := svc.EarnPoints(context.Background(), command.PointEarnCommand{UserUid: "u-1", Amount: 0}); err == nil {
		t.Error("EarnPoints(0) err = nil, want positive-amount error")
	}
}

func TestSpendPointsAppliesNegativeDelta(t *testing.T) {
	var gotDelta int32
	var gotType string
	svc := &customerService{pointRepo: &stubPointRepo{
		applyPointChangeFn: func(_ context.Context, userUid string, delta int32, changeType string, _ *string) (*repo.PointResult, error) {
			gotDelta = delta
			gotType = changeType
			return &repo.PointResult{Uid: "pt-1", UserUid: userUid, Balance: 60 + delta}, nil
		},
	}}
	got, err := svc.SpendPoints(context.Background(), command.PointSpendCommand{UserUid: "u-1", Amount: 40})
	if err != nil {
		t.Fatalf("SpendPoints err: %v", err)
	}
	if gotDelta != -40 || gotType != "SPEND" {
		t.Errorf("ApplyPointChange called with delta=%d type=%s, want -40/SPEND", gotDelta, gotType)
	}
	if got == nil || got.Balance != 20 {
		t.Errorf("got = %+v, want balance=20", got)
	}
}

func TestSpendPointsRejectsNonPositive(t *testing.T) {
	svc := &customerService{}
	if _, err := svc.SpendPoints(context.Background(), command.PointSpendCommand{UserUid: "u-1", Amount: 0}); err == nil {
		t.Error("SpendPoints(0) err = nil, want positive-amount error")
	}
}

func TestSpendPointsInsufficientBalance(t *testing.T) {
	svc := &customerService{pointRepo: &stubPointRepo{
		applyPointChangeFn: func(context.Context, string, int32, string, *string) (*repo.PointResult, error) {
			return nil, fmt.Errorf("%w: have 10, requested 50", repo.ErrInsufficientPoints)
		},
	}}
	_, err := svc.SpendPoints(context.Background(), command.PointSpendCommand{UserUid: "u-1", Amount: 50})
	if !errors.Is(err, repo.ErrInsufficientPoints) {
		t.Errorf("err = %v, want ErrInsufficientPoints", err)
	}
}

func TestSpendPointsNoExistingRow(t *testing.T) {
	svc := &customerService{pointRepo: &stubPointRepo{
		applyPointChangeFn: func(context.Context, string, int32, string, *string) (*repo.PointResult, error) {
			return nil, nil
		},
	}}
	if _, err := svc.SpendPoints(context.Background(), command.PointSpendCommand{UserUid: "u-1", Amount: 50}); !errors.Is(err, repo.ErrInsufficientPoints) {
		t.Errorf("err = %v, want ErrInsufficientPoints (no balance to spend)", err)
	}
}

func TestGetPointHistory(t *testing.T) {
	reason := "signup bonus"
	svc := &customerService{pointHistoryRepo: &stubPointHistoryRepo{
		findAllFn: func(_ context.Context, userUid string, _ pagination.PageOption) (int64, []*repo.PointHistoryResult, error) {
			if userUid != "u-1" {
				t.Errorf("userUid = %q, want u-1", userUid)
			}
			return 1, []*repo.PointHistoryResult{
				{Uid: "h-1", ChangeType: "EARN", Amount: 100, Reason: &reason},
			}, nil
		},
	}}
	total, entries, err := svc.GetPointHistory(context.Background(), "u-1", pagination.NewPageOptionDefault())
	if err != nil {
		t.Fatalf("GetPointHistory err: %v", err)
	}
	if total != 1 || len(entries) != 1 {
		t.Fatalf("total=%d len=%d, want 1/1", total, len(entries))
	}
	if entries[0].Uid != "h-1" || entries[0].Amount != 100 {
		t.Errorf("entries[0] = %+v", entries[0])
	}
}

// ─── Reviews ────────────────────────────────────────────────────────────────

func TestGetCustomerReviews(t *testing.T) {
	svc := &customerService{reviewRepo: &stubReviewRepo{
		findAllFn: func(_ context.Context, opt option.ReviewQueryOption) (int64, []*repo.ReviewResult, error) {
			if opt.UserUid() == nil || *opt.UserUid() != "u-1" {
				t.Errorf("filter userUid = %v, want u-1 (forced regardless of caller opts)", opt.UserUid())
			}
			return 1, []*repo.ReviewResult{{Uid: "r-1", UserUid: "u-1", BookUid: "b-1", Rating: 5}}, nil
		},
	}}
	total, aggs, err := svc.GetCustomerReviews(context.Background(), "u-1", option.NewReviewQueryOption())
	if err != nil {
		t.Fatalf("GetCustomerReviews err: %v", err)
	}
	if total != 1 || len(aggs) != 1 || aggs[0].Uid != "r-1" {
		t.Errorf("got total=%d aggs=%+v", total, aggs)
	}
}

func TestGetBookReviews(t *testing.T) {
	svc := &customerService{reviewRepo: &stubReviewRepo{
		findAllFn: func(_ context.Context, opt option.ReviewQueryOption) (int64, []*repo.ReviewResult, error) {
			if opt.BookUid() == nil || *opt.BookUid() != "b-1" {
				t.Errorf("filter bookUid = %v, want b-1", opt.BookUid())
			}
			return 2, []*repo.ReviewResult{
				{Uid: "r-1", BookUid: "b-1", Rating: 4},
				{Uid: "r-2", BookUid: "b-1", Rating: 5},
			}, nil
		},
	}}
	total, aggs, err := svc.GetBookReviews(context.Background(), "b-1", option.NewReviewQueryOption())
	if err != nil {
		t.Fatalf("GetBookReviews err: %v", err)
	}
	if total != 2 || len(aggs) != 2 {
		t.Errorf("got total=%d len=%d, want 2/2", total, len(aggs))
	}
}

func TestGetReview(t *testing.T) {
	t.Run("owned by caller", func(t *testing.T) {
		svc := &customerService{reviewRepo: &stubReviewRepo{
			findByUidFn: func(context.Context, string) (*repo.ReviewResult, error) {
				return &repo.ReviewResult{Uid: "r-1", UserUid: "u-1", BookUid: "b-1", Rating: 5}, nil
			},
		}}
		got, err := svc.GetReview(context.Background(), "u-1", "r-1")
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if got == nil || got.Uid != "r-1" {
			t.Errorf("got = %+v, want uid=r-1", got)
		}
	})

	t.Run("owned by another user is hidden", func(t *testing.T) {
		svc := &customerService{reviewRepo: &stubReviewRepo{
			findByUidFn: func(context.Context, string) (*repo.ReviewResult, error) {
				return &repo.ReviewResult{Uid: "r-1", UserUid: "other", BookUid: "b-1"}, nil
			},
		}}
		got, err := svc.GetReview(context.Background(), "u-1", "r-1")
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if got != nil {
			t.Errorf("got = %+v, want nil", got)
		}
	})
}

func TestWriteReview(t *testing.T) {
	var got repo.ReviewCreate
	svc := &customerService{reviewRepo: &stubReviewRepo{
		insertFn: func(_ context.Context, review repo.ReviewCreate) (*repo.ReviewResult, error) {
			got = review
			return &repo.ReviewResult{Uid: "r-1", UserUid: review.UserUid, BookUid: review.BookUid, Rating: review.Rating}, nil
		},
	}}
	title := "Solid"
	agg, err := svc.WriteReview(context.Background(), command.ReviewWriteCommand{
		UserUid: "u-1", BookUid: "b-1", Rating: 5, Title: &title,
	})
	if err != nil {
		t.Fatalf("WriteReview err: %v", err)
	}
	if got.UserUid != "u-1" || got.BookUid != "b-1" || got.Rating != 5 || got.Title == nil || *got.Title != "Solid" {
		t.Errorf("Insert called with %+v", got)
	}
	if agg == nil || agg.Uid != "r-1" || agg.Rating != 5 {
		t.Errorf("got = %+v, want uid=r-1 rating=5", agg)
	}
}

func TestWriteReviewDuplicate(t *testing.T) {
	svc := &customerService{reviewRepo: &stubReviewRepo{
		insertFn: func(context.Context, repo.ReviewCreate) (*repo.ReviewResult, error) {
			return nil, repo.ErrReviewExists
		},
	}}
	_, err := svc.WriteReview(context.Background(), command.ReviewWriteCommand{UserUid: "u-1", BookUid: "b-1", Rating: 5})
	if !errors.Is(err, repo.ErrReviewExists) {
		t.Errorf("err = %v, want ErrReviewExists", err)
	}
}

func TestEditReview(t *testing.T) {
	t.Run("applies patch when owned by caller", func(t *testing.T) {
		var updated *entity.CustomerReviewEntity
		svc := &customerService{reviewRepo: &stubReviewRepo{
			fetchByUidFn: func(context.Context, string) ([]*entity.CustomerReviewEntity, error) {
				return []*entity.CustomerReviewEntity{{Id: 1, Uid: "r-1", UserUid: "u-1", Rating: 3}}, nil
			},
			updateFn: func(_ context.Context, ent *entity.CustomerReviewEntity) error {
				updated = ent
				return nil
			},
			findByUidFn: func(context.Context, string) (*repo.ReviewResult, error) {
				return &repo.ReviewResult{Uid: "r-1", UserUid: "u-1", Rating: 5}, nil
			},
		}}
		newRating := int32(5)
		got, err := svc.EditReview(context.Background(), "u-1", "r-1", command.ReviewEditCommand{Rating: &newRating})
		if err != nil {
			t.Fatalf("EditReview err: %v", err)
		}
		if updated == nil || updated.Rating != 5 {
			t.Errorf("updated entity rating = %v, want 5", updated)
		}
		if got == nil || got.Rating != 5 {
			t.Errorf("got = %+v, want rating=5", got)
		}
	})

	t.Run("review of another user is not editable", func(t *testing.T) {
		svc := &customerService{reviewRepo: &stubReviewRepo{
			fetchByUidFn: func(context.Context, string) ([]*entity.CustomerReviewEntity, error) {
				return []*entity.CustomerReviewEntity{{Id: 1, Uid: "r-1", UserUid: "other"}}, nil
			},
			updateFn: func(context.Context, *entity.CustomerReviewEntity) error {
				t.Error("Update must not run for another user's review")
				return nil
			},
		}}
		newRating := int32(1)
		got, err := svc.EditReview(context.Background(), "u-1", "r-1", command.ReviewEditCommand{Rating: &newRating})
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if got != nil {
			t.Errorf("got = %+v, want nil", got)
		}
	})

	t.Run("missing review returns nil", func(t *testing.T) {
		svc := &customerService{reviewRepo: &stubReviewRepo{
			fetchByUidFn: func(context.Context, string) ([]*entity.CustomerReviewEntity, error) {
				return nil, nil
			},
		}}
		got, err := svc.EditReview(context.Background(), "u-1", "missing", command.ReviewEditCommand{})
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if got != nil {
			t.Errorf("got = %+v, want nil", got)
		}
	})
}

func TestRemoveReview(t *testing.T) {
	t.Run("deletes when owned by caller", func(t *testing.T) {
		deleteCalled := false
		svc := &customerService{reviewRepo: &stubReviewRepo{
			fetchByUidFn: func(context.Context, string) ([]*entity.CustomerReviewEntity, error) {
				return []*entity.CustomerReviewEntity{{Id: 1, Uid: "r-1", UserUid: "u-1"}}, nil
			},
			deleteFn: func(context.Context, *entity.CustomerReviewEntity) error {
				deleteCalled = true
				return nil
			},
		}}
		if err := svc.RemoveReview(context.Background(), "u-1", "r-1"); err != nil {
			t.Fatalf("RemoveReview err: %v", err)
		}
		if !deleteCalled {
			t.Error("expected reviewRepo.Delete to be called")
		}
	})

	t.Run("review of another user is not deletable", func(t *testing.T) {
		svc := &customerService{reviewRepo: &stubReviewRepo{
			fetchByUidFn: func(context.Context, string) ([]*entity.CustomerReviewEntity, error) {
				return []*entity.CustomerReviewEntity{{Id: 1, Uid: "r-1", UserUid: "other"}}, nil
			},
			deleteFn: func(context.Context, *entity.CustomerReviewEntity) error {
				t.Error("Delete must not run for another user's review")
				return nil
			},
		}}
		if err := svc.RemoveReview(context.Background(), "u-1", "r-1"); err != nil {
			t.Errorf("err = %v, want nil", err)
		}
	})
}

// ─── package-level helpers ──────────────────────────────────────────────────

func TestToWishlistAggregate(t *testing.T) {
	s := &customerService{}
	results := []*repo.WishlistResult{
		{Uid: "w-1", BookUid: "b-1"},
		{Uid: "w-2", BookUid: "b-2"},
	}
	got := s.toWishlistAggregate("u-1", results)
	if got.UserUid != "u-1" {
		t.Errorf("userUid = %q, want u-1", got.UserUid)
	}
	if !reflect.DeepEqual(got.BookUids, []string{"b-1", "b-2"}) {
		t.Errorf("bookUids = %v", got.BookUids)
	}
}

func TestToWishlistAggregateEmpty(t *testing.T) {
	s := &customerService{}
	got := s.toWishlistAggregate("u-1", nil)
	if got.UserUid != "u-1" {
		t.Errorf("userUid = %q, want u-1", got.UserUid)
	}
	if len(got.BookUids) != 0 {
		t.Errorf("bookUids = %v, want empty", got.BookUids)
	}
}

func TestToPointAggregate(t *testing.T) {
	s := &customerService{}
	got := s.toPointAggregate(&repo.PointResult{Uid: "p-1", UserUid: "u-1", Balance: 1500})
	if got.Uid != "p-1" || got.UserUid != "u-1" || got.Balance != 1500 {
		t.Errorf("got = %+v", got)
	}
}

func TestToReviewAggregate(t *testing.T) {
	if got := toReviewAggregate(nil); got != nil {
		t.Errorf("toReviewAggregate(nil) = %+v, want nil", got)
	}
	title := "Great"
	got := toReviewAggregate(&repo.ReviewResult{Uid: "r-1", UserUid: "u-1", BookUid: "b-1", Rating: 4, Title: &title})
	if got.Uid != "r-1" || got.Rating != 4 || got.Title == nil || *got.Title != "Great" {
		t.Errorf("got = %+v", got)
	}
}

func TestToPointHistoryEntry(t *testing.T) {
	got := toPointHistoryEntry(&repo.PointHistoryResult{Uid: "h-1", ChangeType: "SPEND", Amount: -30})
	if got.Uid != "h-1" || got.Amount != -30 {
		t.Errorf("got = %+v", got)
	}
	if got.ChangeType != aggregate.PointChangeType("SPEND") {
		t.Errorf("changeType = %v, want SPEND", got.ChangeType)
	}
}

func TestDerefString(t *testing.T) {
	if got := derefString(nil); got != "" {
		t.Errorf("derefString(nil) = %q, want empty", got)
	}
	if got := derefString(ptrStrVal("x")); got != "x" {
		t.Errorf("derefString(\"x\") = %q, want x", got)
	}
}

func ptrStrVal(s string) *string { return &s }
