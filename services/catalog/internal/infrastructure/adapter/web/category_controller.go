package web

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/codejsha/bookstore-microservices/catalog/internal/application/port/openapi"
	"github.com/codejsha/bookstore-microservices/catalog/internal/application/usecase"
	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/handler"
	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/model/condition"
	"github.com/codejsha/bookstore-microservices/commonlib-go/pkg/errorcode"
)

var _ openapi.CategoryAPI = (*categoryController)(nil)

func NewCategoryController(
	categoryHandler handler.CategoryHandler,
) openapi.CategoryAPI {
	return &categoryController{
		categoryHandler: categoryHandler,
	}
}

type categoryController struct {
	categoryHandler handler.CategoryHandler
}

func (g categoryController) ApiV1CategoriesGet(c *gin.Context) {
	// prepare request
	ctx := c.Request.Context()
	var param openapi.CategoryFindAllWebParam
	if err := c.ShouldBindQuery(&param); err != nil {
		errorcode.RenderResponse(c, errorcode.InvalidRequestParams, err)
		return
	}
	var categoryCondition condition.CategoryCondition
	cond, err := categoryCondition.FromValue(param)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// perform request
	total, categories, err := g.categoryHandler.HandleFindAll(ctx, usecase.CategoryFindAllQuery{Cond: cond})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// render response
	resps := make([]openapi.CategoryFindWebResp, len(categories))
	for i, category := range categories {
		resps[i] = category.ToCategoryFindWebResp()
	}
	response := openapi.CategoryFindAllWebResp{Total: total, Items: resps}
	c.JSON(http.StatusOK, response)
}
