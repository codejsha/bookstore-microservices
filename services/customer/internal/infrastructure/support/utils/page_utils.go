package utils

import "github.com/codejsha/shared-library-go/pkg/pagination"

const (
	defaultPageSize = 20
	maxPageSize     = 100
)

func PageOffsetLimit(p pagination.PageOption) (int, int) {
	size := int(p.GetSize())
	page := int(p.GetPage())
	if size <= 0 {
		size = defaultPageSize
	} else if size > maxPageSize {
		size = maxPageSize
	}
	if page <= 0 {
		page = 1
	}
	return (page - 1) * size, size
}
