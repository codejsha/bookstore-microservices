package utils

import (
	"math"

	"github.com/codejsha/shared-library-go/pkg/pagination"
)

const (
	defaultPageSize = 20
	maxPageSize     = 100
)

func RowOffsetPage(p pagination.PageOption) pagination.PageOption {
	size := p.GetSize()
	if size <= 0 {
		size = defaultPageSize
	} else if size > maxPageSize {
		size = maxPageSize
	}
	page := p.GetPage()
	if page < 0 {
		page = 0
	}
	offset := int64(page) * int64(size)
	if offset > math.MaxInt32 {
		offset = math.MaxInt32
	}
	rowOffset := int32(offset)
	return p.WithOverrides(&size, &rowOffset, nil)
}
