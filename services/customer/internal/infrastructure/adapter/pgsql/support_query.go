package pgsql

import (
	"context"
	"errors"
	"strings"

	"gorm.io/gen/field"
	"gorm.io/gorm"

	"github.com/codejsha/bookstore-microservices/customer/generated/infrastructure/port/dao"
)

type sortColumn struct {
	asc  field.Expr
	desc field.Expr
}

func orderExprs(sorts []*string, columns map[string]sortColumn, fallback field.Expr) []field.Expr {
	var out []field.Expr
	for _, s := range sorts {
		if s == nil {
			continue
		}
		name, dir, _ := strings.Cut(strings.TrimSpace(*s), ",")
		col, ok := columns[strings.TrimSpace(name)]
		if !ok {
			continue
		}
		if strings.EqualFold(strings.TrimSpace(dir), "desc") {
			out = append(out, col.desc)
		} else {
			out = append(out, col.asc)
		}
	}
	if len(out) == 0 {
		return []field.Expr{fallback}
	}
	return out
}

func categoryIdByUid(ctx context.Context, q *dao.Query, uid string, includeDeleted bool) (*int64, error) {
	c := q.CustomerTicketCategoryEntity
	do := c.WithContext(ctx)
	if includeDeleted {
		do = do.Unscoped()
	}
	e, err := do.Where(c.Uid.Eq(uid)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &e.Id, nil
}

func categoryUidsByIds(ctx context.Context, q *dao.Query, ids []int64) (map[int64]string, error) {
	out := make(map[int64]string, len(ids))
	if len(ids) == 0 {
		return out, nil
	}
	c := q.CustomerTicketCategoryEntity
	rows, err := c.WithContext(ctx).Unscoped().Where(c.Id.In(ids...)).Find()
	if err != nil {
		return nil, err
	}
	for _, row := range rows {
		out[row.Id] = row.Uid
	}
	return out, nil
}

func categoryUidById(ctx context.Context, q *dao.Query, id *int64) (*string, error) {
	if id == nil {
		return nil, nil
	}
	uids, err := categoryUidsByIds(ctx, q, []int64{*id})
	if err != nil {
		return nil, err
	}
	if uid, ok := uids[*id]; ok {
		return &uid, nil
	}
	return nil, nil
}

func collectCategoryIds(ids []*int64) []int64 {
	seen := make(map[int64]struct{}, len(ids))
	out := make([]int64, 0, len(ids))
	for _, id := range ids {
		if id == nil {
			continue
		}
		if _, dup := seen[*id]; dup {
			continue
		}
		seen[*id] = struct{}{}
		out = append(out, *id)
	}
	return out
}

func lookupCategoryUid(uids map[int64]string, id *int64) *string {
	if id == nil {
		return nil
	}
	if uid, ok := uids[*id]; ok {
		return &uid
	}
	return nil
}
