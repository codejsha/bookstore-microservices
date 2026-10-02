package utils

import "github.com/codejsha/shared-library-go/pkg/pagination"

func RowOffsetPage(p pagination.PageOption) pagination.PageOption {
	offset := p.GetPage() * p.GetSize()
	return p.WithOverrides(nil, &offset, nil)
}
