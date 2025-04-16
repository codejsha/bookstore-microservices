package expression

import (
	"gorm.io/gorm/clause"

	"github.com/codejsha/bookstore-microservices/commonlib-go/pkg/convert"
)

func Equal(col string, value interface{}) clause.Expression {
	return clause.Eq{Column: col, Value: value}
}

func NotEqual(col string, value interface{}) clause.Expression {
	return clause.Neq{Column: col, Value: value}
}

func GreaterThan(col string, value interface{}) clause.Expression {
	return clause.Gt{Column: col, Value: value}
}

func GreaterThanOrEqual(col string, value interface{}) clause.Expression {
	return clause.Gte{Column: col, Value: value}
}

func LessThan(col string, value interface{}) clause.Expression {
	return clause.Lt{Column: col, Value: value}
}

func LessThanOrEqual(col string, value interface{}) clause.Expression {
	return clause.Lte{Column: col, Value: value}
}

func In(col string, values ...interface{}) clause.Expression {
	return clause.IN{Column: col, Values: values}
}

func NotIn(col string, values ...interface{}) clause.Expression {
	return clause.Not(clause.IN{Column: col, Values: values})
}

func Like(col string, value interface{}) clause.Expression {
	return clause.Like{Column: col, Value: value}
}

func SortBy(sort, order *string) clause.Expression {
	if sort != nil && order != nil {
		return clause.OrderBy{
			Columns: []clause.OrderByColumn{
				{Column: clause.Column{Name: *sort}, Desc: *order == "desc"},
			},
		}
	}
	return nil
}

func Limit(limit, offset *int32) clause.Expression {
	convertedLimit := convert.ConvertInt[int32, int](convert.ToValue(limit))
	convertedOffset := convert.ConvertInt[int32, int](convert.ToValue(offset))
	return clause.Limit{
		Limit:  &convertedLimit,
		Offset: convertedOffset,
	}
}

func GroupBy(groupBy, having string) clause.Expression {
	return clause.GroupBy{
		Columns: []clause.Column{{Name: groupBy}},
		Having:  []clause.Expression{clause.Expr{SQL: having}},
	}
}

func InnerJoin(table, on string) clause.Expression {
	return clause.Join{
		Type:  clause.InnerJoin,
		Table: clause.Table{Name: table},
		ON: clause.Where{
			Exprs: []clause.Expression{clause.Expr{SQL: on}},
		},
	}
}

func LeftJoin(table, on string) clause.Expression {
	return clause.Join{
		Type:  clause.LeftJoin,
		Table: clause.Table{Name: table},
		ON: clause.Where{
			Exprs: []clause.Expression{clause.Expr{SQL: on}},
		},
	}
}

func RightJoin(table, on string) clause.Expression {
	return clause.Join{
		Type:  clause.RightJoin,
		Table: clause.Table{Name: table},
		ON: clause.Where{
			Exprs: []clause.Expression{clause.Expr{SQL: on}},
		},
	}
}

func CrossJoin(joinTable, on string) clause.Expression {
	return clause.Join{
		Type:  clause.CrossJoin,
		Table: clause.Table{Name: joinTable},
		ON: clause.Where{
			Exprs: []clause.Expression{clause.Expr{SQL: on}},
		},
	}
}

func Sum(col, alias string) clause.Expression {
	return clause.Select{
		Expression: clause.Expr{SQL: "SUM(" + col + ") AS " + alias},
	}
}

func Count(col, alias string) clause.Expression {
	return clause.Select{
		Expression: clause.Expr{SQL: "COUNT(" + col + ") AS " + alias},
	}
}

func CustomExpr(sql string, vars []interface{}) clause.Expression {
	return clause.Expr{SQL: sql, Vars: vars}
}
