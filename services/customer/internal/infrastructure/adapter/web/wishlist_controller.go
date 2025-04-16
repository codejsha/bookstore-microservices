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

var _ openapi.WishlistAPI = (*wishlistController)(nil)

func NewWishlistController(
	wishlistUseCase usecase.WishlistUseCase,
) openapi.WishlistAPI {
	return &wishlistController{
		wishlistUseCase: wishlistUseCase,
	}
}

type wishlistController struct {
	wishlistUseCase usecase.WishlistUseCase
}

func (g wishlistController) ApiV1CustomersIdWishlistGet(c *gin.Context) {
	// prepare request
	ctx := c.Request.Context()
	var param openapi.WishlistFindAllWebParam
	if err := c.ShouldBindUri(&param); err != nil {
		errorcode.RenderResponse(c, errorcode.InvalidRequestParams, err)
		return
	}
	var wishlistCondition condition.WishlistCondition
	cond, err := wishlistCondition.FromValue(param)
	if err != nil {
		errorcode.RenderResponse(c, errorcode.InvalidRequestParams, err)
		return
	}

	// perform request
	wishlist, err := g.wishlistUseCase.FindWishlist(ctx, cond)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// render response
	response := wishlist.ToWishlistFindResponse()
	c.JSON(http.StatusOK, response)
}

func (g wishlistController) ApiV1CustomersIdWishlistPost(c *gin.Context) {
	// TODO implement me
	panic("implement me")
}

func (g wishlistController) ApiV1CustomersIdWishlistPut(c *gin.Context) {
	// prepare request
	ctx := c.Request.Context()
	var param openapi.WishlistUpdateWebParam
	if err := c.ShouldBindUri(&param); err != nil {
		errorcode.RenderResponse(c, errorcode.InvalidRequestParams, err)
		return
	}
	var req openapi.WishlistUpdateWebReq
	if err := c.ShouldBindJSON(&req); err != nil {
		errorcode.RenderResponse(c, errorcode.InvalidRequestPayload, err)
		return
	}
	ents := make([]entity.WishlistEntity, len(req.BookIds))
	for i, bookId := range req.BookIds {
		ents[i] = entity.WishlistEntity{
			Id:     nil,
			UserId: param.Id,
			BookId: bookId,
		}
	}

	// perform request
	wishlist, err := g.wishlistUseCase.UpdateWishlist(ctx, ents)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// render response
	response := wishlist.ToWishlistFindResponse()
	c.JSON(http.StatusOK, response)
}

func (g wishlistController) ApiV1CustomersIdWishlistDelete(c *gin.Context) {
	// prepare request
	ctx := c.Request.Context()
	var param openapi.WishlistDeleteWebParam
	if err := c.ShouldBindUri(&param); err != nil {
		errorcode.RenderResponse(c, errorcode.InvalidRequestParams, err)
		return
	}
	var req openapi.WishlistDeleteWebReq
	if err := c.ShouldBindJSON(&req); err != nil {
		errorcode.RenderResponse(c, errorcode.InvalidRequestPayload, err)
		return
	}

	// perform request
	err := g.wishlistUseCase.DeleteWishlist(ctx, req.BookIds)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// render response
	c.Status(http.StatusNoContent)
}
