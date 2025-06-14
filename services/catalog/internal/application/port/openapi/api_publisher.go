// Code generated by OpenAPI Generator (https://openapi-generator.tech) (7.12.0); DO NOT EDIT.
//
// Catalog service API
//
// Catalog service API of Bookstore microservices
//
// API version: 0.1.0
// Contact: admin@example.com

package openapi

import (
	"github.com/gin-gonic/gin"
)

type PublisherAPI interface {

	// ApiV1PublishersGet Get /api/v1/publishers
	// Find all publishers
	ApiV1PublishersGet(c *gin.Context)

	// ApiV1PublishersIdDelete Delete /api/v1/publishers/:id
	// Delete publisher
	ApiV1PublishersIdDelete(c *gin.Context)

	// ApiV1PublishersIdGet Get /api/v1/publishers/:id
	// Find publisher
	ApiV1PublishersIdGet(c *gin.Context)

	// ApiV1PublishersIdPut Put /api/v1/publishers/:id
	// Update publisher
	ApiV1PublishersIdPut(c *gin.Context)

	// ApiV1PublishersPost Post /api/v1/publishers
	// Add new publisher
	ApiV1PublishersPost(c *gin.Context)
}
