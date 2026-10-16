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

var _ openapi.PublisherApi = (*publisherController)(nil)

type publisherController struct {
	catalogUseCase usecase.CatalogUseCase
}

func NewPublisherController(catalogUseCase usecase.CatalogUseCase) openapi.PublisherApi {
	return &publisherController{catalogUseCase: catalogUseCase}
}

func (c *publisherController) PublishersGetAll(
	ctx context.Context,
	_ openapi.PublisherApiCredentials,
	name *string,
	size *int32,
	page *int32,
	sort *string,
) (openapi.PublishersGetAllResponse, error) {
	opt := option.NewPublisherQueryOption(
		option.PublisherQueryOption{}.WithName(name),
		option.PublisherQueryOption{}.WithPage(pagination.NewPageOption(size, page, sort)),
	)

	total, publishers, err := c.catalogUseCase.FindAllPublishers(ctx, opt)
	if err != nil {
		return nil, err
	}

	items := make([]openapi.PublisherItem, len(publishers))
	for i, p := range publishers {
		items[i] = toPublisherItem(p)
	}
	return openapi.PublishersGetAllOk{Body: openapi.PublisherFindAllResponse{Items: items, Total: total}}, nil
}

func (c *publisherController) PublishersCreate(ctx context.Context, _ openapi.PublisherApiCredentials, req openapi.PublisherCreateRequest) (openapi.PublishersCreateResponse, error) {
	cmd := command.PublisherCreateCommand{
		Name:    req.Name,
		Address: req.Address,
		OlKey:   req.OlKey,
	}

	publisher, err := c.catalogUseCase.CreatePublisher(ctx, cmd)
	if err != nil {
		return nil, httpx.MapConflict(ctx, httpx.MapBusinessError(ctx, err))
	}

	dispatchSideEffects(context.WithoutCancel(ctx), "publisher", publisher.Uid, "created", logrus.Fields{"name": req.Name})
	return openapi.PublishersCreateCreated{Location: "/api/v1/publishers/" + publisher.Uid}, nil
}

func (c *publisherController) PublishersRead(ctx context.Context, _ openapi.PublisherApiCredentials, uid string) (openapi.PublishersReadResponse, error) {
	publisher, err := c.catalogUseCase.FindPublisher(ctx, uid)
	if err != nil {
		return nil, httpx.MapNotFound(ctx, err)
	}

	return openapi.PublishersReadOk{Body: openapi.PublisherFindResponse{
		Uid:     publisher.Uid,
		Name:    publisher.Name,
		Address: publisher.Address,
		OlKey:   publisher.OlKey,
	}}, nil
}

func (c *publisherController) PublishersUpdate(
	ctx context.Context,
	_ openapi.PublisherApiCredentials,
	uid string,
	req openapi.PublisherUpdateRequest,
) (openapi.PublishersUpdateResponse, error) {
	cmd := command.PublisherUpdateCommand{
		Name:    req.Name,
		Address: req.Address,
		OlKey:   req.OlKey,
	}

	publisher, err := c.catalogUseCase.UpdatePublisher(ctx, uid, cmd)
	if err != nil {
		return nil, httpx.MapConflict(ctx, httpx.MapBusinessError(ctx, httpx.MapNotFound(ctx, err)))
	}

	dispatchSideEffects(context.WithoutCancel(ctx), "publisher", uid, "updated", logrus.Fields{})

	return openapi.PublishersUpdateOk{Body: openapi.PublisherUpdateResponse{
		Uid:     publisher.Uid,
		Name:    publisher.Name,
		Address: publisher.Address,
		OlKey:   publisher.OlKey,
	}}, nil
}
