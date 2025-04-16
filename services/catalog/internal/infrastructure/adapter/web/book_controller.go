package web

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/codejsha/bookstore-microservices/catalog/internal/application/port/openapi"
	"github.com/codejsha/bookstore-microservices/catalog/internal/application/usecase"
	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/aggregate/entity"
	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/handler"
	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/model/condition"
	"github.com/codejsha/bookstore-microservices/commonlib-go/pkg/errorcode"
)

var _ openapi.BookAPI = (*bookController)(nil)

func NewBookController(
	bookUseCase usecase.BookUseCase,
	bookHandler handler.BookHandler,
) openapi.BookAPI {
	return &bookController{
		bookUseCase: bookUseCase,
		bookHandler: bookHandler,
	}
}

type bookController struct {
	bookUseCase usecase.BookUseCase
	bookHandler handler.BookHandler
}

func (g bookController) ApiV1BooksGet(c *gin.Context) {
	// prepare request
	ctx := c.Request.Context()
	var param openapi.BookFindAllWebParam
	if err := c.ShouldBindQuery(&param); err != nil {
		errorcode.RenderResponse(c, errorcode.InvalidRequestParams, err)
		return
	}
	var bookCondition condition.BookCondition
	cond, err := bookCondition.FromValue(param)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// perform request
	total, bookAggs, err := g.bookUseCase.FindAllBooks(ctx, cond)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// render response
	resps := make([]openapi.BookFindWebResp, len(bookAggs))
	for i, book := range bookAggs {
		resps[i] = book.ToBookFindWebResp()
	}
	response := openapi.BookFindAllWebResp{Total: total, Items: resps}
	c.JSON(http.StatusOK, response)
}

func (g bookController) ApiV1BooksIdGet(c *gin.Context) {
	// prepare request
	ctx := c.Request.Context()
	var param openapi.BookFindWebParam
	if err := c.ShouldBindUri(&param); err != nil {
		errorcode.RenderResponse(c, errorcode.InvalidRequestParams, err)
		return
	}

	// perform request
	bookAgg, err := g.bookUseCase.FindBook(ctx, param.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// render response
	response := bookAgg.ToBookFindWebResp()
	c.JSON(http.StatusOK, response)
}

func (g bookController) ApiV1BooksPost(c *gin.Context) {
	// prepare request
	ctx := c.Request.Context()
	var req openapi.BookCreateWebReq
	if err := c.ShouldBindJSON(&req); err != nil {
		errorcode.RenderResponse(c, errorcode.InvalidRequestPayload, err)
		return
	}
	var bookEntity entity.BookEntity
	ent, err := bookEntity.FromValue(req, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// perform request
	book, err := g.bookHandler.Handle(ctx, usecase.BookCreateCommand{Ent: ent})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	bookAgg, err := g.bookUseCase.FindBook(ctx, *book.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// render response
	response := bookAgg.ToBookFindWebResp()
	c.JSON(http.StatusCreated, response)
}

func (g bookController) ApiV1BooksIdPut(c *gin.Context) {
	// prepare request
	ctx := c.Request.Context()
	var param openapi.BookUpdateWebParam
	if err := c.ShouldBindUri(&param); err != nil {
		errorcode.RenderResponse(c, errorcode.InvalidRequestParams, err)
		return
	}
	var req openapi.BookUpdateWebReq
	if err := c.ShouldBindJSON(&req); err != nil {
		errorcode.RenderResponse(c, errorcode.InvalidRequestPayload, err)
		return
	}
	var bookEntity entity.BookEntity
	ent, err := bookEntity.FromValue(req, &param.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// perform request
	book, err := g.bookHandler.Handle(ctx, usecase.BookUpdateCommand{Ent: ent})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	bookAgg, err := g.bookUseCase.FindBook(ctx, *book.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// render response
	response := bookAgg.ToBookFindWebResp()
	c.JSON(http.StatusOK, response)
}

func (g bookController) ApiV1BooksIdDelete(c *gin.Context) {
	// prepare request
	ctx := c.Request.Context()
	var param openapi.BookDeleteWebParam
	if err := c.ShouldBindUri(&param); err != nil {
		errorcode.RenderResponse(c, errorcode.InvalidRequestParams, err)
		return
	}

	// perform request
	_, err := g.bookHandler.Handle(ctx, usecase.BookDeleteCommand{Id: param.Id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// render response
	c.Status(http.StatusNoContent)
}
