package pgsql

import (
	"strings"

	"gorm.io/gen/field"
)

func buildOrderExprs(sort string, whitelist map[string]field.OrderExpr, tiebreaker field.OrderExpr) []field.Expr {
	var exprs []field.Expr
	for _, term := range strings.Split(sort, ";") {
		name, desc := parseSortTerm(term)
		col, ok := whitelist[name]
		if !ok {
			continue
		}
		if desc {
			exprs = append(exprs, col.Desc())
		} else {
			exprs = append(exprs, col.Asc())
		}
	}
	return append(exprs, tiebreaker.Asc())
}

func parseSortTerm(term string) (name string, desc bool) {
	term = strings.TrimSpace(term)
	sep := strings.IndexAny(term, ",:")
	if sep < 0 {
		return term, false
	}
	name = strings.TrimSpace(term[:sep])
	dir := strings.ToLower(strings.TrimSpace(term[sep+1:]))
	return name, dir == "desc"
}
