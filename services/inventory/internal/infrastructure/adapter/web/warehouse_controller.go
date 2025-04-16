package web

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/codejsha/bookstore-microservices/commonlib-go/pkg/errorcode"
	"github.com/codejsha/bookstore-microservices/inventory/internal/application/port/openapi"
	"github.com/codejsha/bookstore-microservices/inventory/internal/application/usecase"
	"github.com/codejsha/bookstore-microservices/inventory/internal/domain/aggregate/entity"
	"github.com/codejsha/bookstore-microservices/inventory/internal/domain/handler"
	"github.com/codejsha/bookstore-microservices/inventory/internal/domain/model/condition"
)

var _ openapi.WarehouseAPI = (*warehouseController)(nil)

func NewWarehouseController(
	warehouseHandler handler.WarehouseHandler,
) openapi.WarehouseAPI {
	return &warehouseController{
		warehouseHandler: warehouseHandler,
	}
}

type warehouseController struct {
	warehouseHandler handler.WarehouseHandler
}

func (g warehouseController) ApiV1WarehousesGet(c *gin.Context) {
	// prepare request
	ctx := c.Request.Context()
	var param openapi.WarehouseFindAllWebParam
	if err := c.ShouldBindUri(&param); err != nil {
		errorcode.RenderResponse(c, errorcode.InvalidRequestParams, err)
		return
	}
	var warehouseCondition condition.WarehouseCondition
	cond, err := warehouseCondition.FromValue(param)
	if err != nil {
		errorcode.RenderResponse(c, errorcode.InvalidRequestParams, err)
		return
	}

	// perform request
	total, warehouses, err := g.warehouseHandler.HandleFindAll(ctx, usecase.WarehouseFindAllQuery{Cond: cond})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// render response
	resps := make([]openapi.WarehouseFindWebResp, len(warehouses))
	for _, warehouse := range warehouses {
		resps = append(resps, warehouse.ToWarehouseFindResponse())
	}
	response := openapi.WarehouseFindAllWebResp{Total: total, Items: resps}
	c.JSON(http.StatusOK, response)
}

func (g warehouseController) ApiV1WarehousesIdGet(c *gin.Context) {
	// prepare request
	ctx := c.Request.Context()
	var params openapi.WarehouseFindWebParam
	if err := c.ShouldBindUri(&params); err != nil {
		errorcode.RenderResponse(c, errorcode.InvalidRequestParams, err)
		return
	}

	// perform request
	warehouse, err := g.warehouseHandler.HandleFind(ctx, usecase.WarehouseFindQuery{Id: params.Id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// render response
	response := warehouse.ToWarehouseFindResponse()
	c.JSON(http.StatusOK, response)
}

func (g warehouseController) ApiV1WarehousesPost(c *gin.Context) {
	// prepare request
	ctx := c.Request.Context()
	var req openapi.WarehouseCreateWebReq
	if err := c.ShouldBindJSON(&req); err != nil {
		errorcode.RenderResponse(c, errorcode.InvalidRequestPayload, err)
		return
	}
	var warehouseEntity entity.WarehouseEntity
	ent, err := warehouseEntity.FromValue(req, nil)
	if err != nil {
		errorcode.RenderResponse(c, errorcode.InvalidRequestPayload, err)
		return
	}

	// perform request
	warehouse, err := g.warehouseHandler.Handle(ctx, usecase.WarehouseCreateCommand{Ent: ent})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// render response
	c.Header("Location", c.Request.RequestURI+"/"+strconv.FormatInt(warehouse.Id, 10))
	c.Status(http.StatusCreated)
}

func (g warehouseController) ApiV1WarehousesIdPut(c *gin.Context) {
	// prepare request
	ctx := c.Request.Context()
	var param openapi.WarehouseUpdateWebParam
	if err := c.ShouldBindUri(&param); err != nil {
		errorcode.RenderResponse(c, errorcode.InvalidRequestParams, err)
		return
	}
	var req openapi.WarehouseUpdateWebReq
	if err := c.ShouldBindJSON(&req); err != nil {
		errorcode.RenderResponse(c, errorcode.InvalidRequestPayload, err)
		return
	}
	var warehouseEntity entity.WarehouseEntity
	ent, err := warehouseEntity.FromValue(req, &param.Id)
	if err != nil {
		errorcode.RenderResponse(c, errorcode.InvalidRequestPayload, err)
		return
	}

	// perform request
	warehouse, err := g.warehouseHandler.Handle(ctx, usecase.WarehouseUpdateCommand{Ent: ent})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// render response
	response := warehouse.ToWarehouseFindResponse()
	c.JSON(http.StatusOK, response)
}

func (g warehouseController) ApiV1WarehousesIdDelete(c *gin.Context) {
	// prepare request
	ctx := c.Request.Context()
	var param openapi.WarehouseDeleteWebParam
	if err := c.ShouldBindUri(&param); err != nil {
		errorcode.RenderResponse(c, errorcode.InvalidRequestParams, err)
		return
	}

	// perform request
	_, err := g.warehouseHandler.Handle(ctx, usecase.WarehouseDeleteCommand{Id: param.Id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// render response
	c.Status(http.StatusNoContent)
}
