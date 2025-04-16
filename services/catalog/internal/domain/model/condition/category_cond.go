package condition

import (
	"fmt"

	"github.com/codejsha/bookstore-microservices/catalog/internal/application/port/openapi"
	"github.com/codejsha/bookstore-microservices/catalog/internal/application/port/pb/categorypb"
)

type CategoryCondition struct {
	Sort  string
	Order string
}

func (c CategoryCondition) FromValue(value interface{}) (CategoryCondition, error) {
	switch value.(type) {
	case openapi.CategoryFindAllWebParam:
		v := value.(openapi.CategoryFindAllWebParam)
		cond := CategoryCondition{
			Sort:  v.Sort,
			Order: v.Order,
		}
		return cond, nil
	case *categorypb.CategoryFindAllProtoReq:
		v := value.(*categorypb.CategoryFindAllProtoReq)
		cond := CategoryCondition{
			Sort:  v.Filter.Sort,
			Order: v.Filter.Order,
		}
		return cond, nil
	default:
		return CategoryCondition{}, fmt.Errorf("failed to convert request to CategoryCondition: %v", value)
	}
}
