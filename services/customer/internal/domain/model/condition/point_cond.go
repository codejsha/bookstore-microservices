package condition

import (
	"fmt"

	"gorm.io/gorm/clause"

	"github.com/codejsha/bookstore-microservices/commonlib-go/pkg/database/expression"
	"github.com/codejsha/bookstore-microservices/customer/internal/application/port/openapi"
)

type PointCondition struct {
	UserId *string
	expression.QueryFilter
}

func (c PointCondition) BuildCondition() []clause.Expression {
	conds := expression.NewQueryConditionBuilder(c.QueryFilter)

	if c.UserId != nil {
		conds.AddWhereExprs(expression.Equal("user_id", c.UserId))
	}

	return conds.Build()
}

func (c PointCondition) FromValue(value interface{}) (PointCondition, error) {
	switch value.(type) {
	case openapi.PointFindAllWebParam:
		v := value.(openapi.PointFindAllWebParam)
		cond := PointCondition{
			UserId: &v.Id,
		}
		return cond, nil
	default:
		return PointCondition{}, fmt.Errorf("failed to convert request to PointCondition: %v", value)
	}
}
