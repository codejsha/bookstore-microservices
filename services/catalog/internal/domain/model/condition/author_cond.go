package condition

import (
	"fmt"

	"gorm.io/gorm/clause"

	"github.com/codejsha/bookstore-microservices/catalog/internal/application/port/openapi"
	"github.com/codejsha/bookstore-microservices/catalog/internal/application/port/pb/authorpb"
	"github.com/codejsha/bookstore-microservices/commonlib-go/pkg/database/expression"
)

type AuthorCondition struct {
	Name   *string
	BookId *int64
	expression.QueryFilter
}

func (c AuthorCondition) BuildCondition() []clause.Expression {
	conds := expression.NewQueryConditionBuilder(c.QueryFilter)

	if c.Name != nil {
		conds.AddWhereExprs(expression.Like("name", c.Name))
	}

	if c.BookId != nil {
		conds.AddExprs(expression.InnerJoin("book_author_mapping", "book_author_mapping.author_id = author.id"))
		conds.AddWhereExprs(expression.Equal("book_author_mapping.book_id", c.BookId))
	}

	return conds.Build()
}

func (c AuthorCondition) FromValue(value interface{}) (AuthorCondition, error) {
	switch value.(type) {
	case openapi.AuthorFindAllWebParam:
		v := value.(openapi.AuthorFindAllWebParam)
		cond := AuthorCondition{
			Name:   v.Name,
			BookId: v.BookId,
			QueryFilter: expression.QueryFilter{
				Sort:   v.Sort,
				Order:  v.Order,
				Limit:  v.Limit,
				Offset: v.Offset,
			},
		}
		return cond, nil
	case *authorpb.AuthorFindAllProtoReq:
		v := value.(*authorpb.AuthorFindAllProtoReq)
		cond := AuthorCondition{
			Name:   v.Name,
			BookId: v.BookId,
			QueryFilter: expression.QueryFilter{
				Sort:   v.Filter.Sort,
				Order:  v.Filter.Order,
				Limit:  v.Filter.Limit,
				Offset: v.Filter.Offset,
			},
		}
		return cond, nil
	default:
		return AuthorCondition{}, fmt.Errorf("failed to convert request to AuthorCondition: %v", value)
	}
}
