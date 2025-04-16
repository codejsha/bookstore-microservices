package web

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/codejsha/bookstore-microservices/catalog/internal/application/port/openapi"
	"github.com/codejsha/bookstore-microservices/catalog/internal/application/usecase"
	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/aggregate/entity"
	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/handler"
	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/model/condition"
	"github.com/codejsha/bookstore-microservices/commonlib-go/pkg/errorcode"
)

var _ openapi.AuthorAPI = (*authorController)(nil)

func NewAuthorController(
	authorHandler handler.AuthorHandler,
) openapi.AuthorAPI {
	return &authorController{
		authorHandler: authorHandler,
	}
}

type authorController struct {
	authorHandler handler.AuthorHandler
}

func (g authorController) ApiV1AuthorsGet(c *gin.Context) {
	// prepare request
	ctx := c.Request.Context()
	var param openapi.AuthorFindAllWebParam
	if err := c.ShouldBindQuery(&param); err != nil {
		errorcode.RenderResponse(c, errorcode.InvalidRequestParams, err)
		return
	}
	var authorCondition condition.AuthorCondition
	cond, err := authorCondition.FromValue(param)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// perform request
	total, authors, err := g.authorHandler.HandleFindAll(ctx, usecase.AuthorFindAllQuery{Cond: cond})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// render response
	resps := make([]openapi.AuthorFindWebResp, len(authors))
	for i, author := range authors {
		resps[i] = author.ToAuthorFindWebResp()
	}
	response := openapi.AuthorFindAllWebResp{Total: total, Items: resps}
	c.JSON(http.StatusOK, response)
}

func (g authorController) ApiV1AuthorsIdGet(c *gin.Context) {
	// prepare request
	ctx := c.Request.Context()
	var param openapi.AuthorFindWebParam
	if err := c.ShouldBindUri(&param); err != nil {
		errorcode.RenderResponse(c, errorcode.InvalidRequestParams, err)
		return
	}

	// perform request
	author, err := g.authorHandler.HandleFind(ctx, usecase.AuthorFindQuery{Id: param.Id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// render response
	response := author.ToAuthorFindWebResp()
	c.JSON(http.StatusOK, response)
}

func (g authorController) ApiV1AuthorsPost(c *gin.Context) {
	// prepare request
	ctx := c.Request.Context()
	var req openapi.AuthorCreateWebReq
	if err := c.ShouldBindJSON(&req); err != nil {
		errorcode.RenderResponse(c, errorcode.InvalidRequestPayload, err)
		return
	}
	var authorEntity entity.AuthorEntity
	ent, err := authorEntity.FromValue(req, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// perform request
	author, err := g.authorHandler.Handle(ctx, usecase.AuthorCreateCommand{Ent: ent})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// render response
	c.Header("Location", c.Request.RequestURI+"/"+strconv.FormatInt(*author.Id, 10))
	c.Status(http.StatusCreated)
}

func (g authorController) ApiV1AuthorsIdPut(c *gin.Context) {
	// prepare request
	ctx := c.Request.Context()
	var param openapi.AuthorUpdateWebParam
	if err := c.ShouldBindUri(&param); err != nil {
		errorcode.RenderResponse(c, errorcode.InvalidRequestParams, err)
		return
	}
	var req openapi.AuthorUpdateWebReq
	if err := c.ShouldBindJSON(&req); err != nil {
		errorcode.RenderResponse(c, errorcode.InvalidRequestPayload, err)
		return
	}
	var authorEntity entity.AuthorEntity
	ent, err := authorEntity.FromValue(req, &param.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// perform request
	resp, err := g.authorHandler.Handle(ctx, usecase.AuthorUpdateCommand{Ent: ent})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// render response
	response := resp.ToAuthorFindWebResp()
	c.JSON(http.StatusOK, response)
}

func (g authorController) ApiV1AuthorsIdDelete(c *gin.Context) {
	// prepare request
	ctx := c.Request.Context()
	var param openapi.AuthorDeleteWebParam
	if err := c.ShouldBindUri(&param); err != nil {
		errorcode.RenderResponse(c, errorcode.InvalidRequestParams, err)
		return
	}

	// perform request
	_, err := g.authorHandler.Handle(ctx, usecase.AuthorDeleteCommand{Id: param.Id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// render response
	c.Status(http.StatusNoContent)
}
