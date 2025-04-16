package aggregate

import (
	"time"

	"github.com/codejsha/bookstore-microservices/customer/internal/application/port/openapi"
	"github.com/codejsha/bookstore-microservices/customer/internal/application/port/pb/paymentpb"
)

func NewPaymentAggregate(
	resp *paymentpb.PaymentFindProtoResp,
) *PaymentAggregate {
	agg := &PaymentAggregate{
		Id:          resp.Id,
		UserId:      resp.UserId,
		PaymentType: resp.PaymentType,
		CardNumber:  resp.CardNumber,
		Amount:      resp.Amount,
		PaymentDate: resp.PaymentDate.AsTime(),
	}
	return agg
}

type PaymentAggregate struct {
	// payment id
	Id          int64
	UserId      string
	PaymentType string
	CardNumber  string
	Amount      float64
	PaymentDate time.Time
}

func (a *PaymentAggregate) ToPaymentFindResponse() openapi.PaymentFindWebResp {
	resp := openapi.PaymentFindWebResp{
		Id:          a.Id,
		UserId:      a.UserId,
		PaymentType: a.PaymentType,
		CardNumber:  a.CardNumber,
		Amount:      a.Amount,
		PaymentDate: a.PaymentDate,
	}
	return resp
}
