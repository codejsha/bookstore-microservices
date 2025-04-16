package stockmapper

import (
	"fmt"

	"github.com/codejsha/bookstore-microservices/commonlib-go/pkg/database/expression"
	"github.com/codejsha/bookstore-microservices/inventory/internal/application/port/openapi"
	"github.com/codejsha/bookstore-microservices/inventory/internal/application/port/pb/stockpb"
	"github.com/codejsha/bookstore-microservices/inventory/internal/domain/model/condition"
)

func ToFilterCondition(value interface{}) (condition.StockCondition, error) {
	switch value.(type) {
	case openapi.StockFindAllWebParam:
		v := value.(openapi.StockFindAllWebParam)
		cond := condition.StockCondition{
			BookId:      v.BookId,
			WarehouseId: v.WarehouseId,
			QueryFilter: expression.QueryFilter{
				Sort:   v.Sort,
				Order:  v.Order,
				Limit:  v.Limit,
				Offset: v.Offset,
			},
		}
		return cond, nil
	case *stockpb.StockFindAllProtoReq:
		v := value.(*stockpb.StockFindAllProtoReq)
		cond := condition.StockCondition{
			BookId:      v.BookId,
			WarehouseId: v.WarehouseId,
			QueryFilter: expression.QueryFilter{
				Sort:   v.Filter.Sort,
				Order:  v.Filter.Order,
				Limit:  v.Filter.Limit,
				Offset: v.Filter.Offset,
			},
		}
		return cond, nil
	default:
		return condition.StockCondition{}, fmt.Errorf("failed to convert value to StockCondition: %v", value)
	}
}
