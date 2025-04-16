package mapper

import (
	"fmt"

	"github.com/codejsha/bookstore-microservices/customer/internal/application/port/openapi"
	"github.com/codejsha/bookstore-microservices/customer/internal/application/port/pb/paymentpb"
)

func NewPaymentMapper() *PaymentMapper {
	return &PaymentMapper{}
}

type PaymentMapper struct {
}

func (m PaymentMapper) ToPaymentFindAllProtoReq(value interface{}) (*paymentpb.PaymentFindAllProtoReq, error) {
	switch value.(type) {
	case openapi.PaymentFindAllWebParam:
		v := value.(openapi.PaymentFindAllWebParam)
		request := &paymentpb.PaymentFindAllProtoReq{
			UserId:  &v.Id,
			OrderId: v.OrderId,
			Filter: &paymentpb.QueryFilter{
				Sort:   v.Sort,
				Order:  v.Order,
				Limit:  v.Limit,
				Offset: v.Offset,
			},
		}
		return request, nil
	default:
		return nil, fmt.Errorf("failed to convert request to PaymentFindAllRequest: %v", value)
	}
}

func (m PaymentMapper) ToPaymentFindProtoReq(value interface{}) (*paymentpb.PaymentFindProtoReq, error) {
	switch value.(type) {
	case openapi.PaymentFindWebParam:
		v := value.(openapi.PaymentFindWebParam)
		request := &paymentpb.PaymentFindProtoReq{
			Id: v.PaymentId,
		}
		return request, nil
	default:
		return nil, fmt.Errorf("failed to convert request to PaymentFindRequest: %v", value)
	}
}
