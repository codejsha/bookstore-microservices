package restcontroller

import (
	"context"

	"github.com/codejsha/bookstore-microservices/customer/generated/application/port/openapi"
	"github.com/codejsha/bookstore-microservices/customer/internal/application/usecase"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/model/command"
)

var _ openapi.CategoryApi = (*categoryController)(nil)

type categoryController struct {
	supportUseCase usecase.SupportUseCase
}

func NewCategoryController(supportUseCase usecase.SupportUseCase) openapi.CategoryApi {
	return &categoryController{supportUseCase: supportUseCase}
}

func (c *categoryController) CategoriesGetAll(
	ctx context.Context,
	_ openapi.CategoryApiCredentials,
) (openapi.CategoriesGetAllResponse, error) {
	if _, err := resolveSupportCaller(ctx); err != nil {
		return nil, err
	}
	categories, err := c.supportUseCase.FindAllCategories(ctx)
	if err != nil {
		return nil, err
	}
	items := make([]openapi.CategoryFindResponse, len(categories))
	for i, cat := range categories {
		items[i] = toCategoryFindResponse(cat)
	}
	return openapi.CategoriesGetAllOk{Body: openapi.CategoryFindAllResponse{Items: items}}, nil
}

func (c *categoryController) CategoriesCreate(
	ctx context.Context,
	_ openapi.CategoryApiCredentials,
	req openapi.CategoryCreateRequest,
) (openapi.CategoriesCreateResponse, error) {
	if _, err := resolveSupportCaller(ctx); err != nil {
		return nil, err
	}
	cmd := command.CategoryCreateCommand{
		Name:        req.Name,
		Description: req.Description,
		ParentUid:   req.ParentUid,
	}
	category, err := c.supportUseCase.CreateCategory(ctx, cmd)
	if err != nil {
		return nil, mapSupportError(ctx, err)
	}
	return openapi.CategoriesCreateCreated{Location: "/api/v1/categories/" + category.Uid}, nil
}

func (c *categoryController) CategoriesRead(
	ctx context.Context,
	_ openapi.CategoryApiCredentials,
	uid string,
) (openapi.CategoriesReadResponse, error) {
	if _, err := resolveSupportCaller(ctx); err != nil {
		return nil, err
	}
	category, err := c.supportUseCase.FindCategory(ctx, uid)
	if err != nil {
		return nil, mapSupportError(ctx, err)
	}
	return openapi.CategoriesReadOk{Body: toCategoryFindResponse(category)}, nil
}

func (c *categoryController) CategoriesUpdate(
	ctx context.Context,
	_ openapi.CategoryApiCredentials,
	uid string,
	req openapi.CategoryUpdateRequest,
) (openapi.CategoriesUpdateResponse, error) {
	if _, err := resolveSupportCaller(ctx); err != nil {
		return nil, err
	}
	cmd := command.CategoryUpdateCommand{
		Name:        req.Name,
		Description: req.Description,
		ParentUid:   req.ParentUid,
	}
	category, err := c.supportUseCase.UpdateCategory(ctx, uid, cmd)
	if err != nil {
		return nil, mapSupportError(ctx, err)
	}
	return openapi.CategoriesUpdateOk{Body: toCategoryFindResponse(category)}, nil
}

func (c *categoryController) CategoriesRemove(
	ctx context.Context,
	_ openapi.CategoryApiCredentials,
	uid string,
) (openapi.CategoriesRemoveResponse, error) {
	if _, err := resolveSupportCaller(ctx); err != nil {
		return nil, err
	}
	if err := c.supportUseCase.DeleteCategory(ctx, uid); err != nil {
		return nil, mapSupportError(ctx, err)
	}
	return openapi.CategoriesRemoveNoContent{}, nil
}
