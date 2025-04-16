package condition

import (
	"fmt"

	"gorm.io/gorm/clause"

	"github.com/codejsha/bookstore-microservices/commonlib-go/pkg/database/expression"
	"github.com/codejsha/bookstore-microservices/inventory/internal/application/port/openapi"
)

type StockCondition struct {
	BookId      *int64
	WarehouseId *int64
	expression.QueryFilter
}

func (c StockCondition) BuildCondition() []clause.Expression {
	conds := expression.NewQueryConditionBuilder(c.QueryFilter)

	if c.BookId != nil {
		conds.AddWhereExprs(expression.Equal("book_id", c.BookId))
	}

	if c.WarehouseId != nil {
		conds.AddWhereExprs(expression.Equal("warehouse_id", c.WarehouseId))
	}

	return conds.Build()
}

func (c StockCondition) FromValue(value interface{}) (StockCondition, error) {
	switch value.(type) {
	case openapi.StockFindAllWebParam:
		v := value.(openapi.StockFindAllWebParam)
		cond := StockCondition{
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
	default:
		return StockCondition{}, fmt.Errorf("failed to convert request to StockCondition: %v", value)
	}
}
