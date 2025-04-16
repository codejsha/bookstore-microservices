package web

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/codejsha/bookstore-microservices/commonlib-go/pkg/errorcode"
	"github.com/codejsha/bookstore-microservices/customer/internal/application/port/openapi"
	"github.com/codejsha/bookstore-microservices/customer/internal/application/usecase"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/model/mapper"
)

var _ openapi.OrderAPI = (*orderController)(nil)

func NewOrderController(
	orderUseCase usecase.OrderUseCase,
	orderMapper *mapper.OrderMapper,
) openapi.OrderAPI {
	return &orderController{
		orderUseCase: orderUseCase,
		orderMapper:  orderMapper,
	}
}

type orderController struct {
	orderUseCase usecase.OrderUseCase
	orderMapper  *mapper.OrderMapper
}

func (g orderController) ApiV1CustomersIdOrdersGet(c *gin.Context) {
	// prepare request
	ctx := c.Request.Context()
	var param openapi.OrderFindAllWebParam
	if err := c.ShouldBindUri(&param); err != nil {
		errorcode.RenderResponse(c, errorcode.InvalidRequestParams, err)
		return
	}
	if err := c.ShouldBindQuery(&param); err != nil {
		errorcode.RenderResponse(c, errorcode.InvalidRequestParams, err)
		return
	}
	req, err := g.orderMapper.ToOrderFindAllProtoReq(param)
	if err != nil {
		errorcode.RenderResponse(c, errorcode.InvalidRequestParams, err)
		return
	}

	// perform request
	total, orders, err := g.orderUseCase.FindAllCustomerOrders(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// render response
	resps := make([]openapi.OrderFindWebResp, len(orders))
	for _, order := range orders {
		resps = append(resps, order.ToOrderFindResponse())
	}
	response := openapi.OrderFindAllWebResp{Total: total, Items: resps}
	c.JSON(http.StatusOK, response)
}

func (g orderController) ApiV1CustomersIdOrdersOrderIdGet(c *gin.Context) {
	// prepare request
	ctx := c.Request.Context()
	var param openapi.OrderFindWebParam
	if err := c.ShouldBindUri(&param); err != nil {
		errorcode.RenderResponse(c, errorcode.InvalidRequestParams, err)
		return
	}
	req, err := g.orderMapper.ToOrderFindProtoReq(param)
	if err != nil {
		errorcode.RenderResponse(c, errorcode.InvalidRequestParams, err)
		return
	}

	// perform request
	order, err := g.orderUseCase.FindCustomerOrder(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// render response
	response := order.ToOrderFindResponse()
	c.JSON(http.StatusOK, response)
}
