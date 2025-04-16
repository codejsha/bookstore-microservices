package web

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/codejsha/bookstore-microservices/commonlib-go/pkg/errorcode"
	"github.com/codejsha/bookstore-microservices/customer/internal/application/port/openapi"
	"github.com/codejsha/bookstore-microservices/customer/internal/application/usecase"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/aggregate/entity"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/model/condition"
)

var _ openapi.PointAPI = (*pointController)(nil)

func NewPointController(
	pointUseCase usecase.PointUseCase,
) openapi.PointAPI {
	return &pointController{
		pointUseCase: pointUseCase,
	}
}

type pointController struct {
	pointUseCase usecase.PointUseCase
}

func (g pointController) ApiV1CustomersIdPointGet(c *gin.Context) {
	// prepare request
	ctx := c.Request.Context()
	var param openapi.PointFindAllWebParam
	if err := c.ShouldBindUri(&param); err != nil {
		errorcode.RenderResponse(c, errorcode.InvalidRequestParams, err)
		return
	}
	var pointCondition condition.PointCondition
	cond, err := pointCondition.FromValue(param)
	if err != nil {
		errorcode.RenderResponse(c, errorcode.InvalidRequestParams, err)
		return
	}

	// perform request
	point, err := g.pointUseCase.FindPoint(ctx, cond)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// render response
	response := point.ToPointFindResponse()
	c.JSON(http.StatusOK, response)
}

func (g pointController) ApiV1CustomersIdPointPut(c *gin.Context) {
	// prepare request
	ctx := c.Request.Context()
	var param openapi.PointUpdateWebParam
	if err := c.ShouldBindQuery(&param); err != nil {
		errorcode.RenderResponse(c, errorcode.InvalidRequestParams, err)
		return
	}
	var req openapi.PointUpdateWebReq
	if err := c.ShouldBindJSON(&req); err != nil {
		errorcode.RenderResponse(c, errorcode.InvalidRequestPayload, err)
		return
	}
	ent := entity.PointEntity{
		Id:     nil,
		UserId: param.Id,
		Point:  req.Point,
	}

	// perform request
	point, err := g.pointUseCase.UpdatePoint(ctx, ent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// render response
	response := point.ToPointFindResponse()
	c.JSON(http.StatusOK, response)
}
