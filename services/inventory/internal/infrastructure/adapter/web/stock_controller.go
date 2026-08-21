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
	"github.com/codejsha/bookstore-microservices/inventory/internal/domain/model/mapper/stockmapper"
)

var _ openapi.StocksAPI = (*stockController)(nil)

func NewStockController(
	stockHandler handler.StockHandler,
) openapi.StocksAPI {
	return &stockController{
		stockHandler: stockHandler,
	}
}

type stockController struct {
	stockHandler handler.StockHandler
}

func (g stockController) ApiV1StocksGet(c *gin.Context) {
	// prepare request
	ctx := c.Request.Context()
	var param openapi.StockFindAllWebParam
	if err := c.ShouldBindQuery(&param); err != nil {
		errorcode.RenderResponse(c, errorcode.InvalidRequestParams, err)
		return
	}
	cond, err := stockmapper.ToFilterCondition(param)
	if err != nil {
		errorcode.RenderResponse(c, errorcode.InvalidRequestParams, err)
		return
	}

	// perform request
	total, stocks, err := g.stockHandler.HandleFindAll(ctx, usecase.StockFindAllQuery{Cond: cond})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// render response
	resps := make([]openapi.StockFindWebResp, len(stocks))
	for _, stock := range stocks {
		resps = append(resps, stock.ToStockFindWebResp())
	}
	response := openapi.StockFindAllWebResp{Total: total, Items: resps}
	c.JSON(http.StatusOK, response)
}

func (g stockController) ApiV1StocksIdGet(c *gin.Context) {
	// prepare request
	ctx := c.Request.Context()
	var param openapi.StockFindWebParam
	if err := c.ShouldBindUri(&param); err != nil {
		errorcode.RenderResponse(c, errorcode.InvalidRequestParams, err)
		return
	}

	// perform request
	stock, err := g.stockHandler.HandleFind(ctx, usecase.StockFindQuery{Id: param.Id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// render response
	response := stock.ToStockFindWebResp()
	c.JSON(http.StatusOK, response)
}

func (g stockController) ApiV1StocksPost(c *gin.Context) {
	// prepare request
	ctx := c.Request.Context()
	var req openapi.StockCreateWebReq
	if err := c.ShouldBindJSON(&req); err != nil {
		errorcode.RenderResponse(c, errorcode.InvalidRequestPayload, err)
		return
	}
	var stockEntity entity.StockEntity
	ent, err := stockEntity.FromValue(req, nil)
	if err != nil {
		errorcode.RenderResponse(c, errorcode.InvalidRequestPayload, err)
		return
	}

	// perform request
	stock, err := g.stockHandler.Handle(ctx, usecase.StockCreateCommand{Ent: ent})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// render response
	c.Header("Location", c.Request.RequestURI+"/"+strconv.FormatInt(stock.Id, 10))
	c.Status(http.StatusCreated)
}

func (g stockController) ApiV1StocksIdPut(c *gin.Context) {
	// prepare request
	ctx := c.Request.Context()
	var param openapi.StockUpdateWebParam
	if err := c.ShouldBindUri(&param); err != nil {
		errorcode.RenderResponse(c, errorcode.InvalidRequestParams, err)
		return
	}
	var req openapi.StockUpdateWebReq
	if err := c.ShouldBindJSON(&req); err != nil {
		errorcode.RenderResponse(c, errorcode.InvalidRequestPayload, err)
		return
	}
	var stockEntity entity.StockEntity
	ent, err := stockEntity.FromValue(req, &param.Id)
	if err != nil {
		errorcode.RenderResponse(c, errorcode.InvalidRequestPayload, err)
		return
	}

	// perform request
	stock, err := g.stockHandler.Handle(ctx, usecase.StockUpdateCommand{Ent: ent})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// render response
	response := stock.ToStockFindWebResp()
	c.JSON(http.StatusOK, response)
}

func (g stockController) ApiV1StocksIdDelete(c *gin.Context) {
	// prepare request
	ctx := c.Request.Context()
	var param openapi.StockDeleteWebParam
	if err := c.ShouldBindUri(&param); err != nil {
		errorcode.RenderResponse(c, errorcode.InvalidRequestParams, err)
		return
	}

	// perform request
	_, err := g.stockHandler.Handle(ctx, usecase.StockDeleteCommand{Id: param.Id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// render response
	c.Status(http.StatusNoContent)
}
