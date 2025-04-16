package mapper

import (
	"fmt"

	"github.com/codejsha/bookstore-microservices/customer/internal/application/port/openapi"
	"github.com/codejsha/bookstore-microservices/customer/internal/application/port/pb/orderpb"
)

func NewOrderMapper() *OrderMapper {
	return &OrderMapper{}
}

type OrderMapper struct{}

func (m OrderMapper) ToOrderFindAllProtoReq(value interface{}) (*orderpb.OrderFindAllProtoReq, error) {
	switch value.(type) {
	case openapi.OrderFindAllWebParam:
		v := value.(openapi.OrderFindAllWebParam)
		status, err := orderStatusFromValue(v.Status)
		if err != nil {
			return nil, fmt.Errorf("failed to convert request to OrderFindAllRequest: %v", value)
		}
		request := &orderpb.OrderFindAllProtoReq{
			UserId: &v.Id,
			Status: &status,
			Filter: &orderpb.QueryFilter{
				Sort:   v.Sort,
				Order:  v.Order,
				Limit:  v.Limit,
				Offset: v.Offset,
			},
		}
		return request, nil
	default:
		return nil, fmt.Errorf("failed to convert request to OrderFindAllRequest: %v", value)
	}
}

func (m OrderMapper) ToOrderFindProtoReq(value interface{}) (*orderpb.OrderFindProtoReq, error) {
	switch value.(type) {
	case openapi.OrderFindWebParam:
		v := value.(openapi.OrderFindWebParam)
		request := &orderpb.OrderFindProtoReq{
			Id: v.OrderId,
		}
		return request, nil
	default:
		return nil, fmt.Errorf("failed to convert request to OrderFindRequest: %v", value)
	}
}

func orderStatusFromValue(value interface{}) (orderpb.OrderStatus, error) {
	switch value.(type) {
	case openapi.OrderStatus:
		v := value.(openapi.OrderStatus)
		switch v {
		case openapi.ORDERSTATUS_UNKNOWN:
			return orderpb.OrderStatus_UNKNOWN, nil
		case openapi.ORDERSTATUS_PENDING:
			return orderpb.OrderStatus_PENDING, nil
		case openapi.ORDERSTATUS_SHIPPING:
			return orderpb.OrderStatus_SHIPPING, nil
		case openapi.ORDERSTATUS_COMPLETED:
			return orderpb.OrderStatus_COMPLETED, nil
		case openapi.ORDERSTATUS_CANCELLED:
			return orderpb.OrderStatus_CANCELLED, nil
		default:
			return orderpb.OrderStatus_UNKNOWN, fmt.Errorf("failed to convert value to OrderStatus: %v", value)
		}
	default:
		return orderpb.OrderStatus_UNKNOWN, fmt.Errorf("failed to convert value to OrderStatus: %v", value)
	}
}
