package web

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/codejsha/bookstore-microservices/commonlib-go/pkg/errorcode"
	"github.com/codejsha/bookstore-microservices/customer/internal/application/port/openapi"
	"github.com/codejsha/bookstore-microservices/customer/internal/application/usecase"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/model/mapper"
)

var _ openapi.PaymentAPI = (*paymentController)(nil)

func NewPaymentController(
	paymentUseCase usecase.PaymentUseCase,
	paymentMapper *mapper.PaymentMapper,
) openapi.PaymentAPI {
	return &paymentController{
		paymentUseCase: paymentUseCase,
		paymentMapper:  paymentMapper,
	}
}

type paymentController struct {
	paymentUseCase usecase.PaymentUseCase
	paymentMapper  *mapper.PaymentMapper
}

func (g paymentController) ApiV1CustomersIdPaymentsGet(c *gin.Context) {
	// prepare request
	ctx := c.Request.Context()
	var param openapi.PaymentFindAllWebParam
	if err := c.ShouldBindUri(&param); err != nil {
		errorcode.RenderResponse(c, errorcode.InvalidRequestParams, err)
		return
	}
	if err := c.ShouldBindQuery(&param); err != nil {
		errorcode.RenderResponse(c, errorcode.InvalidRequestParams, err)
		return
	}
	req, err := g.paymentMapper.ToPaymentFindAllProtoReq(param)
	if err != nil {
		errorcode.RenderResponse(c, errorcode.InvalidRequestParams, err)
		return
	}

	// perform request
	total, payments, err := g.paymentUseCase.FindAllCustomerPayments(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// render response
	resps := make([]openapi.PaymentFindWebResp, len(payments))
	for _, payment := range payments {
		resps = append(resps, payment.ToPaymentFindResponse())
	}
	response := openapi.PaymentFindAllWebResp{Total: total, Items: resps}
	c.JSON(http.StatusOK, response)
}

func (g paymentController) ApiV1CustomersIdPaymentsPaymentIdGet(c *gin.Context) {
	// prepare request
	ctx := c.Request.Context()
	var param openapi.PaymentFindWebParam
	if err := c.ShouldBindUri(&param); err != nil {
		errorcode.RenderResponse(c, errorcode.InvalidRequestParams, err)
		return
	}
	req, err := g.paymentMapper.ToPaymentFindProtoReq(param)
	if err != nil {
		errorcode.RenderResponse(c, errorcode.InvalidRequestParams, err)
		return
	}

	// perform request
	payment, err := g.paymentUseCase.FindCustomerPayment(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// render response
	response := payment.ToPaymentFindResponse()
	c.JSON(http.StatusOK, response)
}
