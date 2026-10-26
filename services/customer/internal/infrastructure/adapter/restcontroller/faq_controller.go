package restcontroller

import (
	"context"

	"github.com/codejsha/shared-library-go/pkg/pagination"

	"github.com/codejsha/bookstore-microservices/customer/generated/application/port/openapi"
	"github.com/codejsha/bookstore-microservices/customer/internal/application/usecase"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/model/command"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/model/option"
	"github.com/codejsha/bookstore-microservices/customer/internal/infrastructure/httpx"
)

var _ openapi.FaqApi = (*faqController)(nil)

type faqController struct {
	supportUseCase usecase.SupportUseCase
}

func NewFaqController(supportUseCase usecase.SupportUseCase) openapi.FaqApi {
	return &faqController{supportUseCase: supportUseCase}
}

func (c *faqController) FaqsGetAll(
	ctx context.Context,
	_ openapi.FaqApiCredentials,
	categoryUid *string,
	q *string,
	publishedOnly *bool,
	size *int32,
	page *int32,
	sort *string,
) (openapi.FaqsGetAllResponse, error) {
	caller, err := resolveSupportCaller(ctx)
	if err != nil {
		return nil, err
	}

	onlyPublished := !caller.staff || publishedOnly == nil || *publishedOnly
	opts := []option.FaqQueryOptionFunc{
		option.WithFaqPublishedOnly(onlyPublished),
		option.WithFaqPage(pagination.NewPageOption(size, page, sort)),
	}
	if categoryUid != nil {
		opts = append(opts, option.WithFaqCategoryUid(*categoryUid))
	}
	if q != nil {
		opts = append(opts, option.WithFaqQuery(*q))
	}

	total, faqs, err := c.supportUseCase.FindAllFaqs(ctx, option.NewFaqQueryOption(opts...))
	if err != nil {
		return nil, err
	}

	items := make([]openapi.FaqFindResponse, len(faqs))
	for i, f := range faqs {
		items[i] = toFaqFindResponse(f)
	}
	return openapi.FaqsGetAllOk{Body: openapi.FaqFindAllResponse{Items: items, Total: total}}, nil
}

func (c *faqController) FaqsCreate(
	ctx context.Context,
	_ openapi.FaqApiCredentials,
	req openapi.FaqCreateRequest,
) (openapi.FaqsCreateResponse, error) {
	if _, err := resolveSupportCaller(ctx); err != nil {
		return nil, err
	}
	cmd := command.FaqCreateCommand{
		Question:    req.Question,
		Answer:      req.Answer,
		CategoryUid: req.CategoryUid,
		Published:   req.Published != nil && *req.Published,
	}
	faq, err := c.supportUseCase.CreateFaq(ctx, cmd)
	if err != nil {
		return nil, mapSupportError(ctx, err)
	}
	return openapi.FaqsCreateCreated{Location: "/api/v1/faqs/" + faq.Uid}, nil
}

func (c *faqController) FaqsRead(
	ctx context.Context,
	_ openapi.FaqApiCredentials,
	uid string,
	incrementView *bool,
) (openapi.FaqsReadResponse, error) {
	caller, err := resolveSupportCaller(ctx)
	if err != nil {
		return nil, err
	}
	increment := incrementView == nil || *incrementView
	faq, err := c.supportUseCase.FindFaq(ctx, uid, increment)
	if err != nil {
		return nil, mapSupportError(ctx, err)
	}
	if !faq.Published && !caller.staff {
		return nil, httpx.MapNotFound(ctx, httpx.ErrNotFound)
	}
	return openapi.FaqsReadOk{Body: toFaqFindResponse(faq)}, nil
}

func (c *faqController) FaqsUpdate(
	ctx context.Context,
	_ openapi.FaqApiCredentials,
	uid string,
	req openapi.FaqUpdateRequest,
) (openapi.FaqsUpdateResponse, error) {
	if _, err := resolveSupportCaller(ctx); err != nil {
		return nil, err
	}
	cmd := command.FaqUpdateCommand{
		Question:    req.Question,
		Answer:      req.Answer,
		CategoryUid: req.CategoryUid,
		Published:   req.Published,
	}
	faq, err := c.supportUseCase.UpdateFaq(ctx, uid, cmd)
	if err != nil {
		return nil, mapSupportError(ctx, err)
	}
	return openapi.FaqsUpdateOk{Body: toFaqFindResponse(faq)}, nil
}

func (c *faqController) FaqsRemove(
	ctx context.Context,
	_ openapi.FaqApiCredentials,
	uid string,
) (openapi.FaqsRemoveResponse, error) {
	if _, err := resolveSupportCaller(ctx); err != nil {
		return nil, err
	}
	if err := c.supportUseCase.DeleteFaq(ctx, uid); err != nil {
		return nil, mapSupportError(ctx, err)
	}
	return openapi.FaqsRemoveNoContent{}, nil
}
