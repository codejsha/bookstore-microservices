package condition

import (
	"fmt"

	"gorm.io/gorm/clause"

	"github.com/codejsha/bookstore-microservices/catalog/internal/application/port/openapi"
	"github.com/codejsha/bookstore-microservices/catalog/internal/application/port/pb/bookpb"
	"github.com/codejsha/bookstore-microservices/commonlib-go/pkg/database/expression"
)

type BookCondition struct {
	Title     *string
	Publisher *string
	AuthorIds []int64
	expression.QueryFilter
}

func (c BookCondition) BuildCondition() []clause.Expression {
	conds := expression.NewQueryConditionBuilder(c.QueryFilter)

	if c.Title != nil {
		conds.AddWhereExprs(expression.Like("title", c.Title))
	}

	if c.AuthorIds != nil {
		conds.AddExprs(expression.InnerJoin("book_author_mapping", "book_author_mapping.author_id = author.id"))
		conds.AddWhereExprs(expression.In("author_ids", c.AuthorIds))
	}

	if c.Publisher != nil {
		conds.AddWhereExprs(expression.Like("publisher", c.Publisher))
	}

	return conds.Build()
}

func (c BookCondition) FromValue(value interface{}) (BookCondition, error) {
	switch value.(type) {
	case openapi.BookFindAllWebParam:
		v := value.(openapi.BookFindAllWebParam)
		cond := BookCondition{
			Title:     v.Title,
			Publisher: v.Publisher,
			AuthorIds: v.AuthorIds,
			QueryFilter: expression.QueryFilter{
				Sort:   v.Sort,
				Order:  v.Order,
				Limit:  v.Limit,
				Offset: v.Offset,
			},
		}
		return cond, nil
	case *bookpb.BookFindAllProtoReq:
		v := value.(*bookpb.BookFindAllProtoReq)
		cond := BookCondition{
			Title:     v.Title,
			Publisher: v.Publisher,
			AuthorIds: v.AuthorIds,
			QueryFilter: expression.QueryFilter{
				Sort:   v.Filter.Sort,
				Order:  v.Filter.Order,
				Limit:  v.Filter.Limit,
				Offset: v.Filter.Offset,
			},
		}
		return cond, nil
	default:
		return BookCondition{}, fmt.Errorf("failed to convert request to BookCondition: %v", value)
	}
}
