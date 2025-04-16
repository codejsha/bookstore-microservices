package aggregate

import (
	"github.com/codejsha/bookstore-microservices/catalog/internal/application/port/openapi"
	"github.com/codejsha/bookstore-microservices/catalog/internal/application/port/pb/authorpb"
	"github.com/codejsha/bookstore-microservices/catalog/internal/application/port/pb/bookpb"
	"github.com/codejsha/bookstore-microservices/catalog/internal/application/port/pb/stockpb"
	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/aggregate/entity"
	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/constant"
)

func NewBookAggregate(
	book entity.BookEntity,
	publisher entity.PublisherEntity,
	authors []entity.AuthorEntity,
	stocks *stockpb.StockFindAllProtoResp,
) *BookAggregate {
	category := constant.Category{}.GetCategoryById(book.CategoryId)

	quantity := int64(0)
	reservedQuantity := int64(0)
	for _, stock := range stocks.Items {
		quantity += stock.Quantity
		reservedQuantity += stock.ReservedQuantity
	}

	return &BookAggregate{
		Id:               *book.Id,
		Book:             book,
		Category:         category,
		Publisher:        publisher,
		Authors:          authors,
		Quantity:         quantity,
		ReservedQuantity: reservedQuantity,
	}
}

type BookAggregate struct {
	// book id
	Id               int64
	Book             entity.BookEntity
	Category         constant.Category
	Publisher        entity.PublisherEntity
	Authors          []entity.AuthorEntity
	Quantity         int64
	ReservedQuantity int64
}

func (a *BookAggregate) ToBookFindWebResp() openapi.BookFindWebResp {
	authorResps := make([]openapi.AuthorFindWebResp, len(a.Authors))
	for i, author := range a.Authors {
		authorResps[i] = author.ToAuthorFindWebResp()
	}
	resp := openapi.BookFindWebResp{
		Id:               a.Id,
		Title:            a.Book.Title,
		Isbn:             a.Book.Isbn,
		Price:            a.Book.Price,
		Description:      a.Book.Description,
		Category:         a.Category.ToCategoryFindWebResp(),
		Publisher:        a.Publisher.ToPublisherFindWebResp(),
		Authors:          authorResps,
		Quantity:         a.Quantity,
		ReservedQuantity: a.ReservedQuantity,
	}
	return resp
}

func (a *BookAggregate) ToBookFindResponse() *bookpb.BookFindProtoResp {
	authorResps := make([]*authorpb.AuthorFindProtoResp, len(a.Authors))
	for i, author := range a.Authors {
		authorResps[i] = author.ToAuthorFindProtoResp()
	}
	resp := &bookpb.BookFindProtoResp{
		Id:               a.Id,
		Title:            a.Book.Title,
		Isbn:             &a.Book.Isbn,
		Price:            &a.Book.Price,
		Description:      &a.Book.Description,
		Category:         a.Category.ToCategoryFindProtoResp(),
		Publisher:        a.Publisher.ToPublisherFindProtoResp(),
		Authors:          authorResps,
		Quantity:         &a.Quantity,
		ReservedQuantity: &a.ReservedQuantity,
	}
	return resp
}
