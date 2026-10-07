package utils

import (
	"math"

	"github.com/codejsha/shared-library-go/pkg/pagination"
)

const (
	defaultPageSize = 20
	maxPageSize     = 100
	defaultSort     = "createdat,desc;uid,asc"
)

func PageOffsetLimit(p pagination.PageOption) (int, int) {
	page, size := normalizePage(p)
	return int(page) * int(size), int(size)
}

func RowOffsetPage(p pagination.PageOption) pagination.PageOption {
	page, size := normalizePage(p)
	offset := int64(page) * int64(size)
	if offset > math.MaxInt32 {
		offset = math.MaxInt32
	}
	rowOffset := int32(offset)
	return p.WithOverrides(&size, &rowOffset, nil)
}

func normalizePage(p pagination.PageOption) (int32, int32) {
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
	return page, size
}

func WithDefaultSort(p pagination.PageOption) pagination.PageOption {
	if p.GetSort() != "" {
		return p
	}
	sort := defaultSort
	return p.WithOverrides(nil, nil, []*string{&sort})
}
