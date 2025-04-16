package warehousemapper

import (
	"fmt"

	"github.com/codejsha/bookstore-microservices/commonlib-go/pkg/database/expression"
	"github.com/codejsha/bookstore-microservices/inventory/internal/application/port/openapi"
	"github.com/codejsha/bookstore-microservices/inventory/internal/domain/model/condition"
)

func ToFilterCondition(value interface{}) (condition.WarehouseCondition, error) {
	switch value.(type) {
	case openapi.WarehouseFindAllWebParam:
		v := value.(openapi.WarehouseFindAllWebParam)
		cond := condition.WarehouseCondition{
			Name: v.Name,
			QueryFilter: expression.QueryFilter{
				Sort:   v.Sort,
				Order:  v.Order,
				Limit:  v.Limit,
				Offset: v.Offset,
			},
		}
		return cond, nil
	default:
		return condition.WarehouseCondition{}, fmt.Errorf("failed to convert value to WarehouseCondition: %v", value)
	}
}
