package condition

import (
	"fmt"

	"gorm.io/gorm/clause"

	"github.com/codejsha/bookstore-microservices/catalog/internal/application/port/openapi"
	"github.com/codejsha/bookstore-microservices/commonlib-go/pkg/database/expression"
)

type PublisherCondition struct {
	Name *string
	expression.QueryFilter
}

func (c PublisherCondition) BuildCondition() []clause.Expression {
	conds := expression.NewQueryConditionBuilder(c.QueryFilter)

	if c.Name != nil {
		conds.AddWhereExprs(expression.Like("name", c.Name))
	}

	return conds.Build()
}
func (c PublisherCondition) FromValue(value interface{}) (PublisherCondition, error) {
	switch value.(type) {
	case openapi.PublisherFindAllWebParam:
		v := value.(openapi.PublisherFindAllWebParam)
		cond := PublisherCondition{
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
		return PublisherCondition{}, fmt.Errorf("failed to convert request to PublisherCondition: %v", value)
	}
}
