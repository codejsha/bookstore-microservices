package pgsql

import (
	"strings"

	"gorm.io/gen/field"

	"github.com/codejsha/shared-library-go/pkg/pagination"
)

const defaultSort = "created_at,desc;uid,asc"

func sortableWithAliases(byWireName map[string]field.OrderExpr) map[string]field.OrderExpr {
	out := make(map[string]field.OrderExpr, len(byWireName)*2)
	for name, expr := range byWireName {
		key := strings.ToLower(name)
		out[key] = expr
		out[strings.ReplaceAll(key, "_", "")] = expr
	}
	return out
}

func resolveSort(page pagination.PageOption, sortable map[string]field.OrderExpr) pagination.PageOption {
	kept := make([]*string, 0)
	for _, item := range strings.Split(page.GetSort(), ";") {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		key := strings.ToLower(strings.TrimSpace(strings.SplitN(item, ",", 2)[0]))
		if _, ok := sortable[key]; !ok {
			continue
		}
		kept = append(kept, &item)
	}
	if len(kept) == 0 {
		fallback := defaultSort
		kept = append(kept, &fallback)
	}
	return page.WithOverrides(nil, nil, kept)
}
