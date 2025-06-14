// Code generated by OpenAPI Generator (https://openapi-generator.tech) (7.12.0); DO NOT EDIT.
//
// Customer service API
//
// Customer service API of Bookstore microservices
//
// API version: 0.1.0
// Contact: admin@example.com

package openapi

import (
	"time"
)

type PaymentWebReq struct {
	OrderId int64 `json:"order_id,omitempty"`

	UserId string `json:"user_id,omitempty"`

	PaymentType string `json:"PaymentType,omitempty"`

	CardNumber string `json:"card_number,omitempty"`

	Amount float64 `json:"amount,omitempty"`

	PaymentDate time.Time `json:"payment_date,omitempty"`
}
