package restcontroller

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/codejsha/shared-library-go/pkg/pagination"

	"github.com/codejsha/bookstore-microservices/customer/generated/application/port/openapi"
	"github.com/codejsha/bookstore-microservices/customer/internal/application/usecase"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/aggregate"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/model/command"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/model/option"
	"github.com/codejsha/bookstore-microservices/customer/internal/infrastructure/httpx"
)

var _ openapi.ReviewApi = (*reviewController)(nil)

type reviewController struct {
	customerUseCase usecase.CustomerUseCase
}

func NewReviewController(customerUseCase usecase.CustomerUseCase) openapi.ReviewApi {
	return &reviewController{customerUseCase: customerUseCase}
}

func (c *reviewController) BookReviewsGetAll(
	ctx context.Context,
	bookUid string,
	rating *int32,
	size *int32,
	page *int32,
	sort *string,
) (*openapi.ReviewFindAllResponse, error) {
	opts := []option.ReviewQueryOptionFunc{
		option.WithReviewBookUid(bookUid),
		option.WithReviewPage(pagination.NewPageOption(size, page, sort)),
	}
	if rating != nil {
		opts = append(opts, option.WithReviewRating(*rating))
	}
	opt := option.NewReviewQueryOption(opts...)

	total, reviews, err := c.customerUseCase.GetBookReviews(ctx, bookUid, opt)
	if err != nil {
		return nil, err
	}

	items := make([]openapi.ReviewFindResponse, len(reviews))
	for i, r := range reviews {
		items[i] = toReviewFindResponse(r)
	}
	return &openapi.ReviewFindAllResponse{Items: items, Total: total}, nil
}

func (c *reviewController) CustomerReviewsGetAll(
	ctx context.Context,
	uid string,
	size *int32,
	page *int32,
	sort *string,
) (*openapi.ReviewFindAllResponse, error) {
	opt := option.NewReviewQueryOption(
		option.WithReviewUserUid(uid),
		option.WithReviewPage(pagination.NewPageOption(size, page, sort)),
	)

	total, reviews, err := c.customerUseCase.GetCustomerReviews(ctx, uid, opt)
	if err != nil {
		return nil, err
	}

	items := make([]openapi.ReviewFindResponse, len(reviews))
	for i, r := range reviews {
		items[i] = toReviewFindResponse(r)
	}
	return &openapi.ReviewFindAllResponse{Items: items, Total: total}, nil
}

func (c *reviewController) CustomerReviewsWrite(
	ctx context.Context,
	uid string,
	req openapi.ReviewCreateRequest,
) error {
	cmd := command.ReviewWriteCommand{
		UserUid: uid,
		BookUid: req.BookUid,
		Rating:  req.Rating,
		Title:   req.Title,
		Content: req.Content,
	}

	review, err := c.customerUseCase.WriteReview(ctx, cmd)
	if err != nil {
		return httpx.MapBusinessError(ctx, httpx.MapReviewExists(ctx, err))
	}

	dispatchSideEffects(context.WithoutCancel(ctx), "review", review.Uid, "created", logrus.Fields{"book_uid": req.BookUid})
	return nil
}

func (c *reviewController) CustomerReviewsRead(
	ctx context.Context,
	uid string,
	reviewUid string,
) (*openapi.ReviewFindResponse, error) {
	review, err := c.customerUseCase.GetReview(ctx, uid, reviewUid)
	if err != nil {
		return nil, httpx.MapNotFound(ctx, err)
	}
	if review == nil {
		return nil, httpx.MapNotFound(ctx, httpx.ErrNotFound)
	}

	resp := toReviewFindResponse(review)
	return &resp, nil
}

func (c *reviewController) CustomerReviewsEdit(
	ctx context.Context,
	uid string,
	reviewUid string,
	req openapi.ReviewUpdateRequest,
) (*openapi.ReviewFindResponse, error) {
	cmd := command.ReviewEditCommand{
		Rating:  req.Rating,
		Title:   req.Title,
		Content: req.Content,
	}

	review, err := c.customerUseCase.EditReview(ctx, uid, reviewUid, cmd)
	if err != nil {
		return nil, httpx.MapBusinessError(ctx, httpx.MapNotFound(ctx, err))
	}
	if review == nil {
		return nil, httpx.MapNotFound(ctx, httpx.ErrNotFound)
	}

	dispatchSideEffects(context.WithoutCancel(ctx), "review", reviewUid, "updated", logrus.Fields{})

	resp := toReviewFindResponse(review)
	return &resp, nil
}

func (c *reviewController) CustomerReviewsRemove(ctx context.Context, uid string, reviewUid string) error {
	if err := c.customerUseCase.RemoveReview(ctx, uid, reviewUid); err != nil {
		return err
	}

	dispatchSideEffects(context.WithoutCancel(ctx), "review", reviewUid, "deleted", logrus.Fields{})
	return nil
}

// ─── mapping helpers ────────────────────────────────────────────────────────

func toReviewFindResponse(a *aggregate.ReviewAggregate) openapi.ReviewFindResponse {
	return openapi.ReviewFindResponse{
		Uid:       a.Uid,
		UserUid:   a.UserUid,
		BookUid:   a.BookUid,
		Rating:    a.Rating,
		Title:     a.Title,
		Content:   a.Content,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
}
