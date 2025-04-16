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

var _ openapi.PublisherAPI = (*publisherController)(nil)

func NewPublisherController(
	publisherHandler handler.PublisherHandler,
) openapi.PublisherAPI {
	return &publisherController{
		publisherHandler: publisherHandler,
	}
}

type publisherController struct {
	publisherHandler handler.PublisherHandler
}

func (g publisherController) ApiV1PublishersGet(c *gin.Context) {
	// prepare request
	ctx := c.Request.Context()
	var param openapi.PublisherFindAllWebParam
	if err := c.ShouldBindQuery(&param); err != nil {
		errorcode.RenderResponse(c, errorcode.InvalidRequestParams, err)
		return
	}
	var publisherCondition condition.PublisherCondition
	cond, err := publisherCondition.FromValue(param)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// perform request
	total, publishers, err := g.publisherHandler.HandleFindAll(ctx, usecase.PublisherFindAllQuery{Cond: cond})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// render response
	resps := make([]openapi.PublisherFindWebResp, len(publishers))
	for i, publisher := range publishers {
		resps[i] = publisher.ToPublisherFindWebResp()
	}
	response := openapi.PublisherFindAllWebResp{Total: total, Items: resps}
	c.JSON(http.StatusOK, response)
}

func (g publisherController) ApiV1PublishersIdGet(c *gin.Context) {
	// prepare request
	ctx := c.Request.Context()
	var param openapi.PublisherFindWebParam
	if err := c.ShouldBindUri(&param); err != nil {
		errorcode.RenderResponse(c, errorcode.InvalidRequestParams, err)
		return
	}

	// perform request
	publisher, err := g.publisherHandler.HandleFind(ctx, usecase.PublisherFindQuery{Id: param.Id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// render response
	response := publisher.ToPublisherFindWebResp()
	c.JSON(http.StatusOK, response)
}

func (g publisherController) ApiV1PublishersPost(c *gin.Context) {
	// prepare request
	ctx := c.Request.Context()
	var req openapi.PublisherCreateWebReq
	if err := c.ShouldBindJSON(&req); err != nil {
		errorcode.RenderResponse(c, errorcode.InvalidRequestPayload, err)
		return
	}
	var publisherEntity entity.PublisherEntity
	ent, err := publisherEntity.FromValue(req, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// perform request
	publisher, err := g.publisherHandler.Handle(ctx, usecase.PublisherCreateCommand{Ent: ent})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// render response
	c.Header("Location", c.Request.RequestURI+"/"+strconv.FormatInt(*publisher.Id, 10))
	c.Status(http.StatusCreated)
}

func (g publisherController) ApiV1PublishersIdPut(c *gin.Context) {
	// prepare request
	ctx := c.Request.Context()
	var param openapi.PublisherUpdateWebParam
	if err := c.ShouldBindUri(&param); err != nil {
		errorcode.RenderResponse(c, errorcode.InvalidRequestParams, err)
		return
	}
	var req openapi.PublisherUpdateWebReq
	if err := c.ShouldBindJSON(&req); err != nil {
		errorcode.RenderResponse(c, errorcode.InvalidRequestPayload, err)
		return
	}
	var publisherEntity entity.PublisherEntity
	ent, err := publisherEntity.FromValue(req, &param.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// perform request
	publisher, err := g.publisherHandler.Handle(ctx, usecase.PublisherUpdateCommand{Ent: ent})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// render response
	response := publisher.ToPublisherFindWebResp()
	c.JSON(http.StatusOK, response)
}

func (g publisherController) ApiV1PublishersIdDelete(c *gin.Context) {
	// prepare request
	ctx := c.Request.Context()
	var param openapi.PublisherDeleteWebParam
	if err := c.ShouldBindUri(&param); err != nil {
		errorcode.RenderResponse(c, errorcode.InvalidRequestParams, err)
		return
	}

	// perform request
	_, err := g.publisherHandler.Handle(ctx, usecase.PublisherDeleteCommand{Id: param.Id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// render response
	c.Status(http.StatusNoContent)
}
