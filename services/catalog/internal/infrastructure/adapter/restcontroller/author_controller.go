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

var _ openapi.AuthorApi = (*authorController)(nil)

type authorController struct {
	catalogUseCase usecase.CatalogUseCase
}

func NewAuthorController(catalogUseCase usecase.CatalogUseCase) openapi.AuthorApi {
	return &authorController{catalogUseCase: catalogUseCase}
}

func (c *authorController) AuthorsGetAll(
	ctx context.Context,
	name *string,
	size *int32,
	page *int32,
	sort *string,
) (*openapi.AuthorFindAllResponse, error) {
	opt := option.NewAuthorQueryOption(
		option.AuthorQueryOption{}.WithName(name),
		option.AuthorQueryOption{}.WithPage(pagination.NewPageOption(size, page, sort)),
	)

	total, authors, err := c.catalogUseCase.FindAllAuthors(ctx, opt)
	if err != nil {
		return nil, err
	}

	items := make([]openapi.AuthorItem, len(authors))
	for i, a := range authors {
		items[i] = toAuthorItem(a)
	}
	return &openapi.AuthorFindAllResponse{Items: items, Total: total}, nil
}

func (c *authorController) AuthorsCreate(ctx context.Context, req openapi.AuthorCreateRequest) error {
	cmd := command.AuthorCreateCommand{
		Name:           req.Name,
		Bio:            req.Bio,
		BirthDate:      req.BirthDate,
		DeathDate:      req.DeathDate,
		PhotoUids:      derefStringSlice(req.PhotoUids),
		AlternateNames: derefStringSlice(req.AlternateNames),
		OlKey:          req.OlKey,
	}

	author, err := c.catalogUseCase.CreateAuthor(ctx, cmd)
	if err != nil {
		return err
	}

	dispatchSideEffects(context.WithoutCancel(ctx), "author", author.Uid, "created", logrus.Fields{"name": req.Name})
	return nil
}

func (c *authorController) AuthorsRead(ctx context.Context, uid string) (*openapi.AuthorFindResponse, error) {
	author, err := c.catalogUseCase.FindAuthor(ctx, uid)
	if err != nil {
		return nil, httpx.MapNotFound(ctx, err)
	}

	return &openapi.AuthorFindResponse{
		Uid:            author.Uid,
		Name:           author.Name,
		Bio:            author.Bio,
		BirthDate:      author.BirthDate,
		DeathDate:      author.DeathDate,
		PhotoUids:      ptrStringSlice(author.PhotoUids),
		AlternateNames: ptrStringSlice(author.AlternateNames),
		OlKey:          author.OlKey,
	}, nil
}

func (c *authorController) AuthorsUpdate(
	ctx context.Context,
	uid string,
	req openapi.AuthorUpdateRequest,
) (*openapi.AuthorUpdateResponse, error) {
	cmd := command.AuthorUpdateCommand{
		Name:           req.Name,
		Bio:            req.Bio,
		BirthDate:      req.BirthDate,
		DeathDate:      req.DeathDate,
		PhotoUids:      derefStringSlice(req.PhotoUids),
		AlternateNames: derefStringSlice(req.AlternateNames),
		OlKey:          req.OlKey,
	}

	author, err := c.catalogUseCase.UpdateAuthor(ctx, uid, cmd)
	if err != nil {
		return nil, httpx.MapNotFound(ctx, err)
	}

	dispatchSideEffects(context.WithoutCancel(ctx), "author", uid, "updated", logrus.Fields{})

	return &openapi.AuthorUpdateResponse{
		Uid:            author.Uid,
		Name:           author.Name,
		Bio:            author.Bio,
		BirthDate:      author.BirthDate,
		DeathDate:      author.DeathDate,
		PhotoUids:      ptrStringSlice(author.PhotoUids),
		AlternateNames: ptrStringSlice(author.AlternateNames),
		OlKey:          author.OlKey,
	}, nil
}
