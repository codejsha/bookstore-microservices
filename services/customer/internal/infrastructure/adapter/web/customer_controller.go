package web

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/codejsha/bookstore-microservices/customer/internal/application/port/openapi"
	"github.com/codejsha/bookstore-microservices/customer/internal/application/usecase"

	"github.com/codejsha/bookstore-microservices/commonlib-go/pkg/errorcode"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/model/mapper"
)

var _ openapi.CustomerAPI = (*customerController)(nil)

func NewCustomerController(
	customerUseCase usecase.CustomerUseCase,
	userMapper *mapper.UserMapper,
) openapi.CustomerAPI {
	return &customerController{
		customerUseCase: customerUseCase,
		userMapper:      userMapper,
	}
}

type customerController struct {
	customerUseCase usecase.CustomerUseCase
	userMapper      *mapper.UserMapper
}

func (g customerController) ApiV1CustomersGet(c *gin.Context) {
	// prepare request
	ctx := c.Request.Context()
	var param openapi.CustomerFindAllWebParam
	if err := c.ShouldBindQuery(&param); err != nil {
		errorcode.RenderResponse(c, errorcode.InvalidRequestParams, err)
		return
	}
	req, err := g.userMapper.ToUserFindAllProtoReq(param)
	if err != nil {
		errorcode.RenderResponse(c, errorcode.InvalidRequestParams, err)
		return
	}

	// perform request
	total, customers, err := g.customerUseCase.FindAllCustomers(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// render response
	resps := make([]openapi.CustomerFindWebResp, len(customers))
	for i, customer := range customers {
		resps[i] = customer.ToCustomerFindWebResp()
	}
	response := openapi.CustomerFindAllWebResp{Total: total, Items: resps}
	c.JSON(http.StatusOK, response)
}

func (g customerController) ApiV1CustomersIdGet(c *gin.Context) {
	// prepare request
	ctx := c.Request.Context()
	var param openapi.CustomerFindWebParam
	if err := c.ShouldBindUri(&param); err != nil {
		errorcode.RenderResponse(c, errorcode.InvalidRequestParams, err)
		return
	}
	req, err := g.userMapper.ToUserFindProtoReq(param)
	if err != nil {
		errorcode.RenderResponse(c, errorcode.InvalidRequestParams, err)
		return
	}

	// perform request
	customer, err := g.customerUseCase.FindCustomer(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// render response
	response := customer.ToCustomerFindWebResp()
	c.JSON(http.StatusOK, response)
}

func (g customerController) ApiV1CustomersIdPut(c *gin.Context) {
	// prepare request
	ctx := c.Request.Context()
	var param openapi.CustomerUpdateWebParam
	if err := c.ShouldBindUri(&param); err != nil {
		errorcode.RenderResponse(c, errorcode.InvalidRequestParams, err)
		return
	}
	var req openapi.CustomerUpdateWebReq
	if err := c.ShouldBindJSON(&req); err != nil {
		errorcode.RenderResponse(c, errorcode.InvalidRequestPayload, err)
		return
	}
	protoReq, err := g.userMapper.ToUserUpdateProtoReq(req, param.Id)
	if err != nil {
		errorcode.RenderResponse(c, errorcode.InvalidRequestPayload, err)
		return
	}

	// perform request
	customer, err := g.customerUseCase.UpdateCustomer(ctx, protoReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// render response
	response := customer.ToCustomerFindWebResp()
	c.JSON(http.StatusOK, response)
}

func (g customerController) ApiV1CustomersIdDelete(c *gin.Context) {
	// prepare request
	ctx := c.Request.Context()
	var param openapi.CustomerDeleteWebParam
	if err := c.ShouldBindUri(&param); err != nil {
		errorcode.RenderResponse(c, errorcode.InvalidRequestParams, err)
		return
	}

	// perform request
	err := g.customerUseCase.DeleteCustomer(ctx, param.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// render response
	c.Status(http.StatusNoContent)
}
