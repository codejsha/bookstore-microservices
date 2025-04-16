package condition

import (
	"fmt"

	"gorm.io/gorm/clause"

	"github.com/codejsha/bookstore-microservices/commonlib-go/pkg/database/expression"
	"github.com/codejsha/bookstore-microservices/customer/internal/application/port/openapi"
)

type WishlistCondition struct {
	UserId  *string
	BookIds []int64
	expression.QueryFilter
}

func (c WishlistCondition) BuildCondition() []clause.Expression {
	conds := expression.NewQueryConditionBuilder(c.QueryFilter)

	if c.UserId != nil {
		conds.AddWhereExprs(expression.Equal("user_id", c.UserId))
	}

	if len(c.BookIds) > 0 {
		conds.AddWhereExprs(expression.In("book_id", c.BookIds))
	}

	return conds.Build()
}

func (c WishlistCondition) FromValue(value interface{}) (WishlistCondition, error) {
	switch value.(type) {
	case openapi.WishlistFindAllWebParam:
		v := value.(openapi.WishlistFindAllWebParam)
		cond := WishlistCondition{
			UserId: &v.Id,
		}
		return cond, nil
	default:
		return WishlistCondition{}, fmt.Errorf("failed to convert value to WishlistCondition: %v", value)
	}
}
