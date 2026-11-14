package pgsql

import (
	"strings"

	"gorm.io/gen/field"
)

func containsFold(column field.String, value string) field.Expr {
	return column.Lower().Like("%" + strings.ToLower(value) + "%")
}
