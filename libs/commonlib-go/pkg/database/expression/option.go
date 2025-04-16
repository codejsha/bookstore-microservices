package expression

import (
	"gorm.io/gorm/clause"

	"github.com/codejsha/bookstore-microservices/commonlib-go/pkg/convert"
)

type QueryFilter struct {
	Sort   string
	Order  string
	Limit  int32
	Offset int32
}

type BuildConditionFeature interface {
	BuildCondition() []clause.Expression
}

func NewQueryConditionBuilder(conds QueryFilter) *QueryCondition {
	queryConds := QueryCondition{
		Exprs:      make([]clause.Expression, 0),
		WhereExprs: make([]clause.Expression, 0),
	}
	queryConds.addBaseConditions(conds)
	return &queryConds
}

type QueryCondition struct {
	Exprs      []clause.Expression
	WhereExprs []clause.Expression
}

func (o *QueryCondition) addBaseConditions(conds QueryFilter) {
	o.Exprs = append(o.Exprs, SortBy(convert.ToPtr(conds.Sort), convert.ToPtr(conds.Order)))
	o.Exprs = append(o.Exprs, Limit(convert.ToPtr(conds.Limit), convert.ToPtr(conds.Offset)))
}

func (o *QueryCondition) AddWhereExprs(exprs ...clause.Expression) {
	o.WhereExprs = append(o.WhereExprs, exprs...)
}

func (o *QueryCondition) AddExprs(exprs ...clause.Expression) {
	o.Exprs = append(o.Exprs, exprs...)
}

func (o *QueryCondition) Build() []clause.Expression {
	if len(o.WhereExprs) > 0 {
		o.Exprs = append(o.Exprs, clause.Where{Exprs: o.WhereExprs})
	}
	return o.Exprs
}
