package protostub

import (
	"strconv"

	"github.com/codejsha/shared-library-go/pkg/pagination"
)

type pageRequest struct {
	pageSize  int32
	pageToken string
	orderBy   string
}

func toPageRequest(page pagination.PageOption) pageRequest {
	n := page.Normalized()
	req := pageRequest{pageSize: n.GetSize(), orderBy: page.GetSort()}
	if idx := n.GetPage(); idx > 0 {
		req.pageToken = strconv.FormatInt(int64(idx), 10)
	}
	return req
}
