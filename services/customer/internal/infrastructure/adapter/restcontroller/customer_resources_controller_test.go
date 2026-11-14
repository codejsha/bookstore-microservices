package restcontroller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	"github.com/codejsha/shared-library-go/pkg/pagination"

	"github.com/codejsha/bookstore-microservices/customer/generated/application/port/openapi"
	"github.com/codejsha/bookstore-microservices/customer/internal/application/port/repo"
	"github.com/codejsha/bookstore-microservices/customer/internal/application/usecase"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/aggregate"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/model/command"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/model/option"
	"github.com/codejsha/bookstore-microservices/customer/internal/infrastructure/httpx"
)

const (
	bookUid     = "018f0000-0000-7000-8000-0000000000b1"
	reviewUid   = "018f0000-0000-7000-8000-0000000000b7"
	shipmentUid = "018f0000-0000-7000-8000-0000000000b8"
	orderUid    = "018f0000-0000-7000-8000-0000000000b9"
)

var errBackend = errors.New("pq: connection refused at 10.0.0.7")

type stubCustomerUseCase struct {
	usecase.CustomerUseCase
	err   error
	calls int

	total     int64
	reviews   []*aggregate.ReviewAggregate
	review    *aggregate.ReviewAggregate
	point     *aggregate.PointAggregate
	history   []*aggregate.PointHistoryEntry
	wishlist  *aggregate.WishlistAggregate
	shipment  *aggregate.ShipmentAggregate
	shipments []*aggregate.ShipmentAggregate

	gotUserUid    string
	gotBookUid    string
	gotReviewUid  string
	gotReviewOpt  option.ReviewQueryOption
	gotWrite      command.ReviewWriteCommand
	gotEdit       command.ReviewEditCommand
	gotEarn       command.PointEarnCommand
	gotSpend      command.PointSpendCommand
	gotPage       pagination.PageOption
	gotWishAdd    command.WishlistAddCommand
	gotWishRemove command.WishlistRemoveCommand
	gotShipment   string
	gotOrder      string
	gotStatus     *string
	gotSize       *int32
	gotCaller     usecase.Caller
}

func (s *stubCustomerUseCase) GetBookReviews(_ context.Context, bookUid string, opt option.ReviewQueryOption) (int64, []*aggregate.ReviewAggregate, error) {
	s.calls++
	s.gotBookUid, s.gotReviewOpt = bookUid, opt
	return s.total, s.reviews, s.err
}

func (s *stubCustomerUseCase) GetCustomerReviews(_ context.Context, userUid string, opt option.ReviewQueryOption) (int64, []*aggregate.ReviewAggregate, error) {
	s.calls++
	s.gotUserUid, s.gotReviewOpt = userUid, opt
	return s.total, s.reviews, s.err
}

func (s *stubCustomerUseCase) GetReview(_ context.Context, userUid string, reviewUid string) (*aggregate.ReviewAggregate, error) {
	s.calls++
	s.gotUserUid, s.gotReviewUid = userUid, reviewUid
	return s.review, s.err
}

func (s *stubCustomerUseCase) WriteReview(_ context.Context, cmd command.ReviewWriteCommand) (*aggregate.ReviewAggregate, error) {
	s.calls++
	s.gotWrite = cmd
	if s.err != nil {
		return nil, s.err
	}
	return s.review, nil
}

func (s *stubCustomerUseCase) EditReview(_ context.Context, userUid string, reviewUid string, cmd command.ReviewEditCommand) (*aggregate.ReviewAggregate, error) {
	s.calls++
	s.gotUserUid, s.gotReviewUid, s.gotEdit = userUid, reviewUid, cmd
	return s.review, s.err
}

func (s *stubCustomerUseCase) GetPointBalance(_ context.Context, userUid string) (*aggregate.PointAggregate, error) {
	s.calls++
	s.gotUserUid = userUid
	return s.point, s.err
}

func (s *stubCustomerUseCase) EarnPoints(_ context.Context, cmd command.PointEarnCommand) (*aggregate.PointAggregate, error) {
	s.calls++
	s.gotEarn = cmd
	return s.point, s.err
}

func (s *stubCustomerUseCase) SpendPoints(_ context.Context, cmd command.PointSpendCommand) (*aggregate.PointAggregate, error) {
	s.calls++
	s.gotSpend = cmd
	return s.point, s.err
}

func (s *stubCustomerUseCase) GetPointHistory(_ context.Context, userUid string, page pagination.PageOption) (int64, []*aggregate.PointHistoryEntry, error) {
	s.calls++
	s.gotUserUid, s.gotPage = userUid, page
	return s.total, s.history, s.err
}

func (s *stubCustomerUseCase) GetWishlist(_ context.Context, userUid string) (*aggregate.WishlistAggregate, error) {
	s.calls++
	s.gotUserUid = userUid
	return s.wishlist, s.err
}

func (s *stubCustomerUseCase) AddBooksToWishlist(_ context.Context, cmd command.WishlistAddCommand) (*aggregate.WishlistAggregate, error) {
	s.calls++
	s.gotWishAdd = cmd
	return s.wishlist, s.err
}

func (s *stubCustomerUseCase) RemoveBooksFromWishlist(_ context.Context, cmd command.WishlistRemoveCommand) (*aggregate.WishlistAggregate, error) {
	s.calls++
	s.gotWishRemove = cmd
	return s.wishlist, s.err
}

func (s *stubCustomerUseCase) TrackShipment(_ context.Context, uid string, caller usecase.Caller) (*aggregate.ShipmentAggregate, error) {
	s.calls++
	s.gotShipment, s.gotCaller = uid, caller
	return s.shipment, s.err
}

func (s *stubCustomerUseCase) ListOrderShipments(_ context.Context, orderUid string, status *string, size *int32, caller usecase.Caller) (int64, []*aggregate.ShipmentAggregate, error) {
	s.calls++
	s.gotOrder, s.gotStatus, s.gotSize, s.gotCaller = orderUid, status, size, caller
	return s.total, s.shipments, s.err
}

func newCustomerRouter(uc usecase.CustomerUseCase) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(httpx.GinResponseMapping())
	r.Use(testCallerMiddleware())
	openapi.RegisterReviewApiGinRoutes(r, NewReviewController(uc))
	openapi.RegisterPointApiGinRoutes(r, NewPointController(uc))
	openapi.RegisterWishlistApiGinRoutes(r, NewWishlistController(uc))
	openapi.RegisterDeliveryApiGinRoutes(r, NewDeliveryController(uc))
	return r
}

func callWithoutBearer(r *gin.Engine, method, path string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(method, path, strings.NewReader("")))
	return w
}

func decodeBody(t *testing.T, w *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("body %q is not a json object: %v", w.Body.String(), err)
	}
	return body
}

func itemsOf(t *testing.T, body map[string]any) []map[string]any {
	t.Helper()
	raw, ok := body["items"].([]any)
	if !ok {
		t.Fatalf("items = %v, want an array", body["items"])
	}
	out := make([]map[string]any, len(raw))
	for i, item := range raw {
		out[i] = item.(map[string]any)
	}
	return out
}

func assertProblem(t *testing.T, w *httptest.ResponseRecorder, wantStatus int, wantDetail string) {
	t.Helper()
	if w.Code != wantStatus {
		t.Fatalf("status = %d, want %d: %s", w.Code, wantStatus, w.Body.String())
	}
	if ct := w.Header().Get("Content-Type"); !strings.HasPrefix(ct, "application/problem+json") {
		t.Errorf("content type = %q, want application/problem+json", ct)
	}
	body := decodeBody(t, w)
	if body["status"] != float64(wantStatus) || body["title"] != http.StatusText(wantStatus) {
		t.Errorf("problem = %v, want status %d and its standard title", body, wantStatus)
	}
	if wantDetail != "" && body["detail"] != wantDetail {
		t.Errorf("detail = %v, want %q", body["detail"], wantDetail)
	}
	if strings.Contains(w.Body.String(), "10.0.0.7") {
		t.Errorf("body %q leaks the backend error text", w.Body.String())
	}
}

func invalidCommand(msg string) error {
	return fmt.Errorf("%w: %s", command.ErrInvalidCommand, msg)
}

func sampleReview() *aggregate.ReviewAggregate {
	title := "Great"
	return &aggregate.ReviewAggregate{
		Uid: reviewUid, UserUid: ownerUid, BookUid: bookUid, Rating: 4, Title: &title,
		CreatedAt: time.Date(2026, 11, 1, 9, 0, 0, 0, time.UTC),
	}
}

func customerPath(suffix string) string {
	return "/api/v1/customers/" + ownerUid + suffix
}

func TestBookReviewsGetAll_anonymousWithRatingAndPaging_filtersForwardedAndPageRendered(t *testing.T) {
	uc := &stubCustomerUseCase{total: 7, reviews: []*aggregate.ReviewAggregate{sampleReview()}}
	r := newCustomerRouter(uc)

	w := callWithoutBearer(r, http.MethodGet, "/api/v1/books/"+bookUid+"/reviews?rating=4&size=5&page=2&sort=created_at,desc")

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", w.Code, w.Body.String())
	}
	opt := uc.gotReviewOpt
	if uc.gotBookUid != bookUid || opt.BookUid() == nil || *opt.BookUid() != bookUid {
		t.Errorf("book filter = %q/%v, want %s", uc.gotBookUid, opt.BookUid(), bookUid)
	}
	if opt.Rating() == nil || *opt.Rating() != 4 {
		t.Errorf("rating filter = %v, want 4", opt.Rating())
	}
	if opt.UserUid() != nil {
		t.Errorf("user filter = %v, want none on the public book listing", opt.UserUid())
	}
	if opt.Page().GetSize() != 5 || opt.Page().GetPage() != 2 || opt.Page().GetSort() != "created_at,desc" {
		t.Errorf("page = %d/%d/%q, want 5/2/created_at,desc", opt.Page().GetSize(), opt.Page().GetPage(), opt.Page().GetSort())
	}
	body := decodeBody(t, w)
	items := itemsOf(t, body)
	if body["total"] != float64(7) || len(items) != 1 {
		t.Fatalf("body = %v, want total 7 with one item", body)
	}
	item := items[0]
	if item["uid"] != reviewUid || item["user_uid"] != ownerUid || item["book_uid"] != bookUid || item["rating"] != float64(4) || item["title"] != "Great" {
		t.Errorf("item = %v", item)
	}
	if _, leaked := item["id"]; leaked {
		t.Error("response exposes the internal id")
	}
	if _, present := item["updated_at"]; present {
		t.Error("updated_at rendered for a review that was never edited")
	}
}

func TestBookReviewsGetAll_ratingOmitted_noRatingFilterAndEmptyItemsArray(t *testing.T) {
	uc := &stubCustomerUseCase{}
	r := newCustomerRouter(uc)

	w := callWithoutBearer(r, http.MethodGet, "/api/v1/books/"+bookUid+"/reviews")

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", w.Code, w.Body.String())
	}
	if uc.gotReviewOpt.Rating() != nil {
		t.Errorf("rating filter = %v, want none", uc.gotReviewOpt.Rating())
	}
	if items := itemsOf(t, decodeBody(t, w)); len(items) != 0 {
		t.Errorf("items = %v, want an empty array", items)
	}
}

func TestBookReviewsGetAll_ratingNotNumeric_400WithoutCallingUseCase(t *testing.T) {
	uc := &stubCustomerUseCase{}
	r := newCustomerRouter(uc)

	w := callWithoutBearer(r, http.MethodGet, "/api/v1/books/"+bookUid+"/reviews?rating=five")

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", w.Code)
	}
	if uc.calls != 0 {
		t.Error("use case invoked for an unparseable rating")
	}
}

func TestBookReviewsGetAll_useCaseFails_500ProblemWithoutBackendText(t *testing.T) {
	r := newCustomerRouter(&stubCustomerUseCase{err: errBackend})

	w := callWithoutBearer(r, http.MethodGet, "/api/v1/books/"+bookUid+"/reviews")

	assertProblem(t, w, http.StatusInternalServerError, "internal server error")
}

func TestCustomerReviewsGetAll_pathCustomer_userFilterPinnedToPathUid(t *testing.T) {
	uc := &stubCustomerUseCase{total: 1, reviews: []*aggregate.ReviewAggregate{sampleReview()}}
	r := newCustomerRouter(uc)

	w := call(t, r, owner, http.MethodGet, customerPath("/reviews?size=3&page=1"), "")

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", w.Code, w.Body.String())
	}
	opt := uc.gotReviewOpt
	if uc.gotUserUid != ownerUid || opt.UserUid() == nil || *opt.UserUid() != ownerUid {
		t.Errorf("user filter = %q/%v, want %s", uc.gotUserUid, opt.UserUid(), ownerUid)
	}
	if opt.BookUid() != nil || opt.Rating() != nil {
		t.Errorf("book/rating filters = %v/%v, want none", opt.BookUid(), opt.Rating())
	}
	if opt.Page().GetSize() != 3 || opt.Page().GetPage() != 1 {
		t.Errorf("page = %d/%d, want 3/1", opt.Page().GetSize(), opt.Page().GetPage())
	}
	if body := decodeBody(t, w); body["total"] != float64(1) || itemsOf(t, body)[0]["uid"] != reviewUid {
		t.Errorf("body = %v", body)
	}
}

func TestCustomerReviewsGetAll_useCaseFails_500Problem(t *testing.T) {
	r := newCustomerRouter(&stubCustomerUseCase{err: errBackend})

	w := call(t, r, owner, http.MethodGet, customerPath("/reviews"), "")

	assertProblem(t, w, http.StatusInternalServerError, "internal server error")
}

func TestCustomerRoutes_bearerMissing_401WithoutCallingUseCase(t *testing.T) {
	uc := &stubCustomerUseCase{}
	r := newCustomerRouter(uc)

	for _, path := range []string{
		customerPath("/reviews"),
		customerPath("/reviews/" + reviewUid),
		customerPath("/points"),
		customerPath("/points/history"),
		customerPath("/wishlist"),
		"/api/v1/shipments/" + shipmentUid,
		"/api/v1/orders/" + orderUid + "/shipments",
	} {
		if w := callWithoutBearer(r, http.MethodGet, path); w.Code != http.StatusUnauthorized {
			t.Errorf("GET %s = %d, want 401", path, w.Code)
		}
	}
	if uc.calls != 0 {
		t.Errorf("use case invoked %d times without credentials", uc.calls)
	}
}

func TestCustomerReviewsWrite_validReview_201WithLocationOfNewReview(t *testing.T) {
	uc := &stubCustomerUseCase{review: sampleReview()}
	r := newCustomerRouter(uc)

	w := call(t, r, owner, http.MethodPost, customerPath("/reviews"), `{"book_uid":"`+bookUid+`","rating":4,"title":"Great","content":"Loved it"}`)
	WaitForSideEffects(context.Background())

	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201: %s", w.Code, w.Body.String())
	}
	if got, want := w.Header().Get("Location"), customerPath("/reviews/"+reviewUid); got != want {
		t.Errorf("location = %q, want %q", got, want)
	}
	cmd := uc.gotWrite
	if cmd.UserUid != ownerUid || cmd.BookUid != bookUid || cmd.Rating != 4 || cmd.Title == nil || *cmd.Title != "Great" || cmd.Content == nil || *cmd.Content != "Loved it" {
		t.Errorf("command = %+v", cmd)
	}
}

func TestCustomerReviewsWrite_bookAlreadyReviewed_409Problem(t *testing.T) {
	r := newCustomerRouter(&stubCustomerUseCase{err: fmt.Errorf("write review: %w", repo.ErrReviewExists)})

	w := call(t, r, owner, http.MethodPost, customerPath("/reviews"), `{"book_uid":"`+bookUid+`","rating":4}`)

	assertProblem(t, w, http.StatusConflict, "review already exists for this book")
}

func TestCustomerReviewsWrite_commandRejected_400ProblemWithReason(t *testing.T) {
	r := newCustomerRouter(&stubCustomerUseCase{err: invalidCommand("book_uid must be a uuid")})

	w := call(t, r, owner, http.MethodPost, customerPath("/reviews"), `{"book_uid":"not-a-uuid","rating":4}`)

	assertProblem(t, w, http.StatusBadRequest, "invalid command: book_uid must be a uuid")
}

func TestCustomerReviewsWrite_bodyViolatesContract_400WithoutCallingUseCase(t *testing.T) {
	cases := map[string]string{
		"ratingAboveFive": `{"book_uid":"` + bookUid + `","rating":6}`,
		"ratingMissing":   `{"book_uid":"` + bookUid + `"}`,
		"bookMissing":     `{"rating":3}`,
		"malformedJson":   `{"book_uid":`,
	}
	for name, body := range cases {
		t.Run(name, func(t *testing.T) {
			uc := &stubCustomerUseCase{review: sampleReview()}
			r := newCustomerRouter(uc)

			w := call(t, r, owner, http.MethodPost, customerPath("/reviews"), body)

			if w.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want 400: %s", w.Code, w.Body.String())
			}
			if uc.calls != 0 {
				t.Error("use case invoked for a body that violates the contract")
			}
		})
	}
}

func TestCustomerReviewsWrite_useCaseFails_500Problem(t *testing.T) {
	r := newCustomerRouter(&stubCustomerUseCase{err: errBackend})

	w := call(t, r, owner, http.MethodPost, customerPath("/reviews"), `{"book_uid":"`+bookUid+`","rating":4}`)

	assertProblem(t, w, http.StatusInternalServerError, "internal server error")
}

func TestCustomerReviewsRead_ownReview_200WithReview(t *testing.T) {
	uc := &stubCustomerUseCase{review: sampleReview()}
	r := newCustomerRouter(uc)

	w := call(t, r, owner, http.MethodGet, customerPath("/reviews/"+reviewUid), "")

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", w.Code, w.Body.String())
	}
	if uc.gotUserUid != ownerUid || uc.gotReviewUid != reviewUid {
		t.Errorf("lookup = %q/%q, want the path customer and review", uc.gotUserUid, uc.gotReviewUid)
	}
	body := decodeBody(t, w)
	if body["uid"] != reviewUid || body["rating"] != float64(4) || body["created_at"] != "2026-11-01T09:00:00Z" {
		t.Errorf("body = %v", body)
	}
}

func TestCustomerReviewsRead_reviewMissingOrForeign_404Problem(t *testing.T) {
	cases := map[string]*stubCustomerUseCase{
		"useCaseNotFound":  {err: fmt.Errorf("get review: %w", usecase.ErrNotFound)},
		"recordNotFound":   {err: fmt.Errorf("get review: %w", gorm.ErrRecordNotFound)},
		"nilReviewNoError": {},
	}
	for name, uc := range cases {
		t.Run(name, func(t *testing.T) {
			w := call(t, newCustomerRouter(uc), owner, http.MethodGet, customerPath("/reviews/"+missingUid), "")

			assertProblem(t, w, http.StatusNotFound, "resource not found")
		})
	}
}

func TestCustomerReviewsRead_useCaseFails_500Problem(t *testing.T) {
	r := newCustomerRouter(&stubCustomerUseCase{err: errBackend})

	w := call(t, r, owner, http.MethodGet, customerPath("/reviews/"+reviewUid), "")

	assertProblem(t, w, http.StatusInternalServerError, "internal server error")
}

func TestCustomerReviewsEdit_partialUpdate_200WithOnlySentFieldsInCommand(t *testing.T) {
	edited := sampleReview()
	edited.Rating = 2
	updatedAt := time.Date(2026, 11, 2, 9, 0, 0, 0, time.UTC)
	edited.UpdatedAt = &updatedAt
	uc := &stubCustomerUseCase{review: edited}
	r := newCustomerRouter(uc)

	w := call(t, r, owner, http.MethodPut, customerPath("/reviews/"+reviewUid), `{"rating":2}`)
	WaitForSideEffects(context.Background())

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", w.Code, w.Body.String())
	}
	if uc.gotUserUid != ownerUid || uc.gotReviewUid != reviewUid {
		t.Errorf("target = %q/%q, want the path customer and review", uc.gotUserUid, uc.gotReviewUid)
	}
	if uc.gotEdit.Rating == nil || *uc.gotEdit.Rating != 2 || uc.gotEdit.Title != nil || uc.gotEdit.Content != nil {
		t.Errorf("command = %+v, want only the rating set", uc.gotEdit)
	}
	body := decodeBody(t, w)
	if body["rating"] != float64(2) || body["updated_at"] != "2026-11-02T09:00:00Z" {
		t.Errorf("body = %v", body)
	}
}

func TestCustomerReviewsEdit_reviewMissingOrForeign_404Problem(t *testing.T) {
	cases := map[string]*stubCustomerUseCase{
		"useCaseNotFound":  {err: fmt.Errorf("edit review: %w", usecase.ErrNotFound)},
		"nilReviewNoError": {},
	}
	for name, uc := range cases {
		t.Run(name, func(t *testing.T) {
			w := call(t, newCustomerRouter(uc), owner, http.MethodPut, customerPath("/reviews/"+missingUid), `{"rating":2}`)

			assertProblem(t, w, http.StatusNotFound, "resource not found")
		})
	}
}

func TestCustomerReviewsEdit_commandRejected_400ProblemWithReason(t *testing.T) {
	r := newCustomerRouter(&stubCustomerUseCase{err: invalidCommand("nothing to update")})

	w := call(t, r, owner, http.MethodPut, customerPath("/reviews/"+reviewUid), `{}`)

	assertProblem(t, w, http.StatusBadRequest, "invalid command: nothing to update")
}

func TestCustomerReviewsEdit_ratingOutOfRange_400WithoutCallingUseCase(t *testing.T) {
	uc := &stubCustomerUseCase{review: sampleReview()}
	r := newCustomerRouter(uc)

	w := call(t, r, owner, http.MethodPut, customerPath("/reviews/"+reviewUid), `{"rating":0}`)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", w.Code)
	}
	if uc.calls != 0 {
		t.Error("use case invoked for an out-of-range rating")
	}
}

func TestPointsGetBalance_existingBalance_200WithPathCustomerAndBalance(t *testing.T) {
	uc := &stubCustomerUseCase{point: &aggregate.PointAggregate{Uid: "p-1", UserUid: ownerUid, Balance: 1500}}
	r := newCustomerRouter(uc)

	w := call(t, r, owner, http.MethodGet, customerPath("/points"), "")

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", w.Code, w.Body.String())
	}
	if uc.gotUserUid != ownerUid {
		t.Errorf("user = %q, want %s", uc.gotUserUid, ownerUid)
	}
	if body := decodeBody(t, w); body["user_uid"] != ownerUid || body["balance"] != float64(1500) {
		t.Errorf("body = %v", body)
	}
}

func TestPointsGetBalance_useCaseFails_500Problem(t *testing.T) {
	r := newCustomerRouter(&stubCustomerUseCase{err: errBackend})

	w := call(t, r, owner, http.MethodGet, customerPath("/points"), "")

	assertProblem(t, w, http.StatusInternalServerError, "internal server error")
}

func TestPointsEarn_validAmount_200WithNewBalance(t *testing.T) {
	uc := &stubCustomerUseCase{point: &aggregate.PointAggregate{UserUid: ownerUid, Balance: 160}}
	r := newCustomerRouter(uc)

	w := call(t, r, owner, http.MethodPost, customerPath("/points/earn"), `{"amount":60,"reason":"promotion"}`)
	WaitForSideEffects(context.Background())

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", w.Code, w.Body.String())
	}
	if cmd := uc.gotEarn; cmd.UserUid != ownerUid || cmd.Amount != 60 || cmd.Reason == nil || *cmd.Reason != "promotion" {
		t.Errorf("command = %+v", cmd)
	}
	if body := decodeBody(t, w); body["user_uid"] != ownerUid || body["balance"] != float64(160) {
		t.Errorf("body = %v", body)
	}
}

func TestPointsEarn_balanceLimitExceeded_400Problem(t *testing.T) {
	r := newCustomerRouter(&stubCustomerUseCase{err: fmt.Errorf("earn: %w: have 2147483640, adding 60", repo.ErrPointBalanceOverflow)})

	w := call(t, r, owner, http.MethodPost, customerPath("/points/earn"), `{"amount":60}`)

	assertProblem(t, w, http.StatusBadRequest, "point balance limit exceeded")
}

func TestPointsEarn_commandRejected_400ProblemWithReason(t *testing.T) {
	r := newCustomerRouter(&stubCustomerUseCase{err: invalidCommand("user_uid must be a uuid")})

	w := call(t, r, owner, http.MethodPost, customerPath("/points/earn"), `{"amount":60}`)

	assertProblem(t, w, http.StatusBadRequest, "invalid command: user_uid must be a uuid")
}

func TestPointsSpend_validAmount_200WithNewBalance(t *testing.T) {
	uc := &stubCustomerUseCase{point: &aggregate.PointAggregate{UserUid: ownerUid, Balance: 40}}
	r := newCustomerRouter(uc)

	w := call(t, r, owner, http.MethodPost, customerPath("/points/spend"), `{"amount":60}`)
	WaitForSideEffects(context.Background())

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", w.Code, w.Body.String())
	}
	if cmd := uc.gotSpend; cmd.UserUid != ownerUid || cmd.Amount != 60 || cmd.Reason != nil {
		t.Errorf("command = %+v", cmd)
	}
	if body := decodeBody(t, w); body["balance"] != float64(40) {
		t.Errorf("body = %v", body)
	}
}

func TestPointsSpend_balanceTooLow_400Problem(t *testing.T) {
	r := newCustomerRouter(&stubCustomerUseCase{err: fmt.Errorf("spend: %w: have 10, requested 60", repo.ErrInsufficientPoints)})

	w := call(t, r, owner, http.MethodPost, customerPath("/points/spend"), `{"amount":60}`)

	assertProblem(t, w, http.StatusBadRequest, "insufficient points")
}

func TestPointsSpend_commandRejected_400ProblemWithReason(t *testing.T) {
	r := newCustomerRouter(&stubCustomerUseCase{err: invalidCommand("amount must be positive")})

	w := call(t, r, owner, http.MethodPost, customerPath("/points/spend"), `{"amount":60}`)

	assertProblem(t, w, http.StatusBadRequest, "invalid command: amount must be positive")
}

func TestPointsSpend_amountZero_400WithoutCallingUseCase(t *testing.T) {
	uc := &stubCustomerUseCase{point: &aggregate.PointAggregate{}}
	r := newCustomerRouter(uc)

	w := call(t, r, owner, http.MethodPost, customerPath("/points/spend"), `{"amount":0}`)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", w.Code)
	}
	if uc.calls != 0 {
		t.Error("use case invoked for a zero amount")
	}
}

func TestPointsHistory_pagedRequest_200WithEntriesAndPagingForwarded(t *testing.T) {
	reason := "promotion"
	uc := &stubCustomerUseCase{total: 12, history: []*aggregate.PointHistoryEntry{
		{Uid: "018f0000-0000-7000-8000-000000000a01", ChangeType: aggregate.POINTCHANGE_EARN, Amount: 60, Reason: &reason, CreatedAt: time.Date(2026, 11, 1, 9, 0, 0, 0, time.UTC)},
		{Uid: "018f0000-0000-7000-8000-000000000a02", ChangeType: aggregate.POINTCHANGE_SPEND, Amount: -20, CreatedAt: time.Date(2026, 11, 1, 8, 0, 0, 0, time.UTC)},
	}}
	r := newCustomerRouter(uc)

	w := call(t, r, owner, http.MethodGet, customerPath("/points/history?size=2&page=3"), "")

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", w.Code, w.Body.String())
	}
	if uc.gotUserUid != ownerUid || uc.gotPage.GetSize() != 2 || uc.gotPage.GetPage() != 3 {
		t.Errorf("request = %q %d/%d, want the path customer and page 2/3", uc.gotUserUid, uc.gotPage.GetSize(), uc.gotPage.GetPage())
	}
	body := decodeBody(t, w)
	items := itemsOf(t, body)
	if body["total"] != float64(12) || len(items) != 2 {
		t.Fatalf("body = %v, want total 12 with two items", body)
	}
	if items[0]["change_type"] != "EARN" || items[0]["amount"] != float64(60) || items[0]["reason"] != "promotion" {
		t.Errorf("first entry = %v", items[0])
	}
	if _, present := items[1]["reason"]; present || items[1]["change_type"] != "SPEND" || items[1]["amount"] != float64(-20) {
		t.Errorf("second entry = %v", items[1])
	}
}

func TestPointsHistory_useCaseFails_500Problem(t *testing.T) {
	r := newCustomerRouter(&stubCustomerUseCase{err: errBackend})

	w := call(t, r, owner, http.MethodGet, customerPath("/points/history"), "")

	assertProblem(t, w, http.StatusInternalServerError, "internal server error")
}

func TestWishlistsGet_customerWithBooks_200WithBookUids(t *testing.T) {
	uc := &stubCustomerUseCase{wishlist: &aggregate.WishlistAggregate{UserUid: ownerUid, BookUids: []string{bookUid, orderUid}}}
	r := newCustomerRouter(uc)

	w := call(t, r, owner, http.MethodGet, customerPath("/wishlist"), "")

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", w.Code, w.Body.String())
	}
	if uc.gotUserUid != ownerUid {
		t.Errorf("user = %q, want %s", uc.gotUserUid, ownerUid)
	}
	if got := decodeBody(t, w)["book_uids"]; !reflect.DeepEqual(got, []any{bookUid, orderUid}) {
		t.Errorf("book_uids = %v", got)
	}
}

func TestWishlistsGet_useCaseFails_500Problem(t *testing.T) {
	r := newCustomerRouter(&stubCustomerUseCase{err: errBackend})

	w := call(t, r, owner, http.MethodGet, customerPath("/wishlist"), "")

	assertProblem(t, w, http.StatusInternalServerError, "internal server error")
}

func TestWishlistsAddBooks_validBooks_200WithUpdatedWishlist(t *testing.T) {
	uc := &stubCustomerUseCase{wishlist: &aggregate.WishlistAggregate{UserUid: ownerUid, BookUids: []string{bookUid}}}
	r := newCustomerRouter(uc)

	w := call(t, r, owner, http.MethodPost, customerPath("/wishlist/add"), `{"book_uids":["`+bookUid+`"]}`)
	WaitForSideEffects(context.Background())

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", w.Code, w.Body.String())
	}
	if cmd := uc.gotWishAdd; cmd.UserUid != ownerUid || !reflect.DeepEqual(cmd.BookUids, []string{bookUid}) {
		t.Errorf("command = %+v", cmd)
	}
	if got := decodeBody(t, w)["book_uids"]; !reflect.DeepEqual(got, []any{bookUid}) {
		t.Errorf("book_uids = %v", got)
	}
}

func TestWishlistsRemoveBooks_validBooks_200WithUpdatedWishlist(t *testing.T) {
	uc := &stubCustomerUseCase{wishlist: &aggregate.WishlistAggregate{UserUid: ownerUid, BookUids: []string{}}}
	r := newCustomerRouter(uc)

	w := call(t, r, owner, http.MethodPost, customerPath("/wishlist/remove"), `{"book_uids":["`+bookUid+`"]}`)
	WaitForSideEffects(context.Background())

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", w.Code, w.Body.String())
	}
	if cmd := uc.gotWishRemove; cmd.UserUid != ownerUid || !reflect.DeepEqual(cmd.BookUids, []string{bookUid}) {
		t.Errorf("command = %+v", cmd)
	}
	if got := decodeBody(t, w)["book_uids"]; !reflect.DeepEqual(got, []any{}) {
		t.Errorf("book_uids = %v, want an empty array", got)
	}
}

func TestWishlistsChange_commandRejected_400ProblemWithReason(t *testing.T) {
	for _, action := range []string{"add", "remove"} {
		t.Run(action, func(t *testing.T) {
			r := newCustomerRouter(&stubCustomerUseCase{err: invalidCommand("book_uids[0] must be a uuid")})

			w := call(t, r, owner, http.MethodPost, customerPath("/wishlist/"+action), `{"book_uids":["nope"]}`)

			assertProblem(t, w, http.StatusBadRequest, "invalid command: book_uids[0] must be a uuid")
		})
	}
}

func TestWishlistsChange_emptyBookList_400WithoutCallingUseCase(t *testing.T) {
	for _, action := range []string{"add", "remove"} {
		t.Run(action, func(t *testing.T) {
			uc := &stubCustomerUseCase{wishlist: &aggregate.WishlistAggregate{}}
			r := newCustomerRouter(uc)

			w := call(t, r, owner, http.MethodPost, customerPath("/wishlist/"+action), `{"book_uids":[]}`)

			if w.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want 400: %s", w.Code, w.Body.String())
			}
			if uc.calls != 0 {
				t.Error("use case invoked for an empty book list")
			}
		})
	}
}

func sampleShipment() *aggregate.ShipmentAggregate {
	return &aggregate.ShipmentAggregate{
		Uid: shipmentUid, OrderUid: orderUid, Status: aggregate.SHIPMENTSTATUS_IN_TRANSIT,
		TrackingNumber: "1Z-TRACK-9", DestinationCity: "Seattle",
	}
}

func TestDeliveryTrack_callerPrincipal_forwardedToUseCase(t *testing.T) {
	cases := []struct {
		name string
		who  principal
		want usecase.Caller
	}{
		{"customer", owner, usecase.Caller{UserUid: ownerUid}},
		{"staff", staff, usecase.Caller{UserUid: staffUid, IsAdmin: true}},
		{"system", system, usecase.Caller{UserUid: staffUid, IsAdmin: true}},
		{"noPrincipal", anonymous, usecase.Caller{}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			uc := &stubCustomerUseCase{shipment: sampleShipment()}
			r := newCustomerRouter(uc)

			w := call(t, r, c.who, http.MethodGet, "/api/v1/shipments/"+shipmentUid, "")

			if w.Code != http.StatusOK {
				t.Fatalf("status = %d, want 200: %s", w.Code, w.Body.String())
			}
			if uc.gotShipment != shipmentUid || uc.gotCaller != c.want {
				t.Errorf("use case got %q as %+v, want %q as %+v", uc.gotShipment, uc.gotCaller, shipmentUid, c.want)
			}
		})
	}
}

func TestDeliveryTrack_shipmentFound_200WithStatusNameAndEmptyFieldsOmitted(t *testing.T) {
	r := newCustomerRouter(&stubCustomerUseCase{shipment: sampleShipment()})

	w := call(t, r, owner, http.MethodGet, "/api/v1/shipments/"+shipmentUid, "")

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", w.Code, w.Body.String())
	}
	body := decodeBody(t, w)
	if body["uid"] != shipmentUid || body["order_uid"] != orderUid || body["status"] != "IN_TRANSIT" ||
		body["tracking_number"] != "1Z-TRACK-9" || body["destination_city"] != "Seattle" {
		t.Errorf("body = %v", body)
	}
	for _, key := range []string{"carrier_uid", "destination_state", "actual_delivery_at", "updated_at", "id"} {
		if _, present := body[key]; present {
			t.Errorf("%s rendered although the shipment has no value for it", key)
		}
	}
}

func TestDeliveryTrack_shipmentMissing_404Problem(t *testing.T) {
	cases := map[string]*stubCustomerUseCase{
		"useCaseNotFound":    {err: fmt.Errorf("track shipment: %w", usecase.ErrNotFound)},
		"nilShipmentNoError": {},
		"upstreamNotFound":   {err: fmt.Errorf("track shipment: %w", status.Error(codes.NotFound, "no such shipment"))},
	}
	for name, uc := range cases {
		t.Run(name, func(t *testing.T) {
			w := call(t, newCustomerRouter(uc), owner, http.MethodGet, "/api/v1/shipments/"+missingUid, "")

			assertProblem(t, w, http.StatusNotFound, "resource not found")
		})
	}
}

func TestDelivery_upstreamGrpcFailure_translatedToHttpProblem(t *testing.T) {
	cases := []struct {
		code       codes.Code
		wantStatus int
		wantDetail string
	}{
		{codes.InvalidArgument, http.StatusBadRequest, "invalid argument"},
		{codes.Unauthenticated, http.StatusUnauthorized, "unauthenticated"},
		{codes.PermissionDenied, http.StatusForbidden, "forbidden"},
		{codes.NotFound, http.StatusNotFound, "resource not found"},
		{codes.Unavailable, http.StatusServiceUnavailable, "upstream service unavailable"},
		{codes.DeadlineExceeded, http.StatusGatewayTimeout, "upstream service timed out"},
		{codes.Internal, http.StatusInternalServerError, "internal server error"},
	}
	paths := map[string]string{
		"track": "/api/v1/shipments/" + shipmentUid,
		"list":  "/api/v1/orders/" + orderUid + "/shipments",
	}
	for _, c := range cases {
		for name, path := range paths {
			t.Run(name+"_"+c.code.String(), func(t *testing.T) {
				uc := &stubCustomerUseCase{err: fmt.Errorf("fulfillment: %w", status.Error(c.code, "upstream said 10.0.0.7 is down"))}

				w := call(t, newCustomerRouter(uc), owner, http.MethodGet, path, "")

				assertProblem(t, w, c.wantStatus, c.wantDetail)
			})
		}
	}
}

func TestDeliveryList_statusAndSizeGiven_filtersAndCallerForwardedAndPageRendered(t *testing.T) {
	delivered := sampleShipment()
	delivered.Status = aggregate.SHIPMENTSTATUS_DELIVERED
	uc := &stubCustomerUseCase{total: 4, shipments: []*aggregate.ShipmentAggregate{delivered}}
	r := newCustomerRouter(uc)

	w := call(t, r, staff, http.MethodGet, "/api/v1/orders/"+orderUid+"/shipments?status=DELIVERED&size=5", "")

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", w.Code, w.Body.String())
	}
	if uc.gotOrder != orderUid || uc.gotStatus == nil || *uc.gotStatus != "DELIVERED" || uc.gotSize == nil || *uc.gotSize != 5 {
		t.Errorf("request = %q status %v size %v, want the order, DELIVERED and 5", uc.gotOrder, uc.gotStatus, uc.gotSize)
	}
	if want := (usecase.Caller{UserUid: staffUid, IsAdmin: true}); uc.gotCaller != want {
		t.Errorf("caller = %+v, want %+v", uc.gotCaller, want)
	}
	body := decodeBody(t, w)
	items := itemsOf(t, body)
	if body["total"] != float64(4) || len(items) != 1 || items[0]["uid"] != shipmentUid || items[0]["status"] != "DELIVERED" {
		t.Errorf("body = %v", body)
	}
}

func TestDeliveryList_noFilters_nilStatusAndSizeAndEmptyItemsArray(t *testing.T) {
	uc := &stubCustomerUseCase{}
	r := newCustomerRouter(uc)

	w := call(t, r, owner, http.MethodGet, "/api/v1/orders/"+orderUid+"/shipments", "")

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", w.Code, w.Body.String())
	}
	if uc.gotStatus != nil || uc.gotSize != nil {
		t.Errorf("status/size = %v/%v, want both absent", uc.gotStatus, uc.gotSize)
	}
	if want := (usecase.Caller{UserUid: ownerUid}); uc.gotCaller != want {
		t.Errorf("caller = %+v, want %+v", uc.gotCaller, want)
	}
	if items := itemsOf(t, decodeBody(t, w)); len(items) != 0 {
		t.Errorf("items = %v, want an empty array", items)
	}
}

func TestDeliveryList_statusNotInContract_400WithoutCallingUseCase(t *testing.T) {
	uc := &stubCustomerUseCase{}
	r := newCustomerRouter(uc)

	w := call(t, r, owner, http.MethodGet, "/api/v1/orders/"+orderUid+"/shipments?status=TELEPORTED", "")

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400: %s", w.Code, w.Body.String())
	}
	if uc.calls != 0 {
		t.Error("use case invoked for a status the contract does not declare")
	}
}

func TestDeliveryList_plainUseCaseFailure_500Problem(t *testing.T) {
	r := newCustomerRouter(&stubCustomerUseCase{err: errBackend})

	w := call(t, r, owner, http.MethodGet, "/api/v1/orders/"+orderUid+"/shipments", "")

	assertProblem(t, w, http.StatusInternalServerError, "internal server error")
}
