package restcontroller

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/codejsha/shared-library-go/pkg/pagination"

	"github.com/codejsha/bookstore-microservices/catalog/generated/application/port/openapi"
	"github.com/codejsha/bookstore-microservices/catalog/internal/application/usecase"
	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/model/command"
	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/model/option"
	"github.com/codejsha/bookstore-microservices/catalog/internal/infrastructure/httpx"
)

var _ openapi.EditionApi = (*editionController)(nil)

type editionController struct {
	catalogUseCase usecase.CatalogUseCase
}

func NewEditionController(catalogUseCase usecase.CatalogUseCase) openapi.EditionApi {
	return &editionController{catalogUseCase: catalogUseCase}
}

func (c *editionController) EditionsSearch(
	ctx context.Context,
	title *string,
	isbn *string,
	workUid *string,
	publisherUid *string,
	language *string,
	olKey *string,
	size *int32,
	page *int32,
	sort *string,
) (*openapi.EditionFindAllResponse, error) {
	opt := option.NewEditionQueryOption(
		option.EditionQueryOption{}.WithTitle(title),
		option.EditionQueryOption{}.WithWorkUid(workUid),
		option.EditionQueryOption{}.WithPublisherUid(publisherUid),
		option.EditionQueryOption{}.WithIsbn(isbn),
		option.EditionQueryOption{}.WithOlKey(olKey),
		option.EditionQueryOption{}.WithLanguage(language),
		option.EditionQueryOption{}.WithPage(pagination.NewPageOption(size, page, sort)),
	)

	total, editions, err := c.catalogUseCase.SearchEditions(ctx, opt)
	if err != nil {
		return nil, err
	}

	items := make([]openapi.EditionItem, len(editions))
	for i, e := range editions {
		items[i] = toEditionItem(e)
	}
	return &openapi.EditionFindAllResponse{Items: items, Total: total}, nil
}

func (c *editionController) EditionsCreate(ctx context.Context, req openapi.EditionCreateRequest) error {
	cmd := command.EditionCreateCommand{
		Title:          req.Title,
		Isbn10:         req.Isbn10,
		Isbn13:         req.Isbn13,
		NumberOfPages:  req.NumberOfPages,
		PublishDate:    req.PublishDate,
		CoverUids:      derefStringSlice(req.CoverUids),
		Languages:      derefStringSlice(req.Languages),
		PhysicalFormat: req.PhysicalFormat,
		Description:    req.Description,
		WorkUid:        req.WorkUid,
		PublisherUid:   req.PublisherUid,
		OlKey:          req.OlKey,
	}

	edition, err := c.catalogUseCase.CreateEdition(ctx, cmd)
	if err != nil {
		return httpx.MapConflict(ctx, httpx.MapBusinessError(ctx, httpx.MapNotFound(ctx, err)))
	}

	dispatchSideEffects(context.WithoutCancel(ctx), "edition", edition.Uid, "created", logrus.Fields{"title": req.Title})
	return nil
}

func (c *editionController) EditionsRead(ctx context.Context, uid string) (*openapi.EditionFindResponse, error) {
	edition, err := c.catalogUseCase.FindEdition(ctx, uid)
	if err != nil {
		return nil, httpx.MapNotFound(ctx, err)
	}

	resp := &openapi.EditionFindResponse{
		Uid:            edition.Uid,
		Title:          edition.Title,
		Isbn10:         edition.Isbn10,
		Isbn13:         edition.Isbn13,
		NumberOfPages:  edition.NumberOfPages,
		PublishDate:    edition.PublishDate,
		CoverUids:      ptrStringSlice(edition.CoverUids),
		Languages:      ptrStringSlice(edition.Languages),
		PhysicalFormat: edition.PhysicalFormat,
		Description:    edition.Description,
		OlKey:          edition.OlKey,
		CreatedAt:      edition.CreatedAt,
		UpdatedAt:      edition.UpdatedAt,
	}
	if edition.Work != nil {
		resp.Work = toWorkItem(edition.Work)
	}
	if edition.Publisher != nil {
		pub := toPublisherItem(edition.Publisher)
		resp.Publisher = &pub
	}
	return resp, nil
}

func (c *editionController) EditionsUpdate(
	ctx context.Context,
	uid string,
	req openapi.EditionUpdateRequest,
) (*openapi.EditionUpdateResponse, error) {
	cmd := command.EditionUpdateCommand{
		Title:          req.Title,
		Isbn10:         req.Isbn10,
		Isbn13:         req.Isbn13,
		NumberOfPages:  req.NumberOfPages,
		PublishDate:    req.PublishDate,
		CoverUids:      derefStringSlice(req.CoverUids),
		Languages:      derefStringSlice(req.Languages),
		PhysicalFormat: req.PhysicalFormat,
		Description:    req.Description,
		WorkUid:        req.WorkUid,
		PublisherUid:   req.PublisherUid,
		OlKey:          req.OlKey,
	}

	edition, err := c.catalogUseCase.UpdateEdition(ctx, uid, cmd)
	if err != nil {
		return nil, httpx.MapConflict(ctx, httpx.MapBusinessError(ctx, httpx.MapNotFound(ctx, err)))
	}

	dispatchSideEffects(context.WithoutCancel(ctx), "edition", uid, "updated", logrus.Fields{})

	resp := &openapi.EditionUpdateResponse{
		Uid:            edition.Uid,
		Title:          edition.Title,
		Isbn10:         edition.Isbn10,
		Isbn13:         edition.Isbn13,
		NumberOfPages:  edition.NumberOfPages,
		PublishDate:    edition.PublishDate,
		CoverUids:      ptrStringSlice(edition.CoverUids),
		Languages:      ptrStringSlice(edition.Languages),
		PhysicalFormat: edition.PhysicalFormat,
		Description:    edition.Description,
		OlKey:          edition.OlKey,
		CreatedAt:      edition.CreatedAt,
		UpdatedAt:      edition.UpdatedAt,
	}
	if edition.Work != nil {
		resp.Work = toWorkItem(edition.Work)
	}
	if edition.Publisher != nil {
		pub := toPublisherItem(edition.Publisher)
		resp.Publisher = &pub
	}
	return resp, nil
}
