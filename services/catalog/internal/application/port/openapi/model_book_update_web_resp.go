// Code generated by OpenAPI Generator (https://openapi-generator.tech) (7.12.0); DO NOT EDIT.
//
// Catalog service API
//
// Catalog service API of Bookstore microservices
//
// API version: 0.1.0
// Contact: admin@example.com

package openapi

type BookUpdateWebResp struct {
	Id int64 `json:"id,omitempty"`

	Title string `json:"title,omitempty"`

	Isbn string `json:"isbn,omitempty"`

	Price float64 `json:"price,omitempty"`

	Description string `json:"description,omitempty"`

	Category CategoryFindWebResp `json:"category,omitempty"`

	Publisher PublisherFindWebResp `json:"publisher,omitempty"`

	Authors []AuthorFindWebResp `json:"authors,omitempty"`

	Quantity int64 `json:"quantity,omitempty"`

	ReservedQuantity int64 `json:"reserved_quantity,omitempty"`
}
