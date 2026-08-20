package pgsql

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/codejsha/shared-library-go/pkg/database"

	"github.com/codejsha/bookstore-microservices/inventory/internal/application/port/repo"
	"github.com/codejsha/bookstore-microservices/inventory/internal/domain/model/option"
	"github.com/codejsha/bookstore-microservices/inventory/internal/infrastructure/support/utils"
)

var _ repo.StockBalanceRepo = (*stockBalanceRepository)(nil)

type stockBalanceRepository struct {
	db *gorm.DB
}

func NewStockBalanceRepository(dataSource *database.DataSource) repo.StockBalanceRepo {
	return &stockBalanceRepository{db: dataSource.DB()}
}

type balanceRow struct {
	EditionUid       string `gorm:"column:edition_uid"`
	WarehouseUid     string `gorm:"column:warehouse_uid"`
	WarehouseName    string `gorm:"column:warehouse_name"`
	InboundQuantity  int32  `gorm:"column:inbound_quantity"`
	OutboundQuantity int32  `gorm:"column:outbound_quantity"`
	AdjustQuantity   int32  `gorm:"column:adjust_quantity"`
	CurrentQuantity  int32  `gorm:"column:current_quantity"`
}

func (r *stockBalanceRepository) FindAll(ctx context.Context, opt option.BalanceQueryOption) (int64, []*repo.BalanceResult, error) {
	where, args, err := r.buildFilter(opt)
	if err != nil {
		return 0, nil, err
	}

	histJoin, histArgs, err := r.historyJoin(opt)
	if err != nil {
		return 0, nil, err
	}

	offset, limit := utils.PageOffsetLimit(opt.Page())

	query := fmt.Sprintf(`
SELECT
	s.edition_uid   AS edition_uid,
	s.warehouse_uid AS warehouse_uid,
	w.name          AS warehouse_name,
	s.quantity      AS current_quantity,
	COALESCE(SUM(CASE WHEN h.change_type IN ('INBOUND','RELEASE') THEN h.change_qty ELSE 0 END), 0)    AS inbound_quantity,
	COALESCE(SUM(CASE WHEN h.change_type IN ('OUTBOUND','RESERVATION') THEN -h.change_qty ELSE 0 END), 0) AS outbound_quantity,
	COALESCE(SUM(CASE WHEN h.change_type = 'ADJUSTMENT' THEN h.change_qty ELSE 0 END), 0)              AS adjust_quantity
FROM stock s
JOIN warehouse w ON w.id = s.warehouse_id
LEFT JOIN stock_history h ON h.stock_id = s.id AND h.deleted_at IS NULL%s
WHERE s.deleted_at IS NULL%s
GROUP BY s.id, s.edition_uid, s.warehouse_uid, w.name, s.quantity
ORDER BY s.id
LIMIT ? OFFSET ?`, histJoin, where)

	scanArgs := make([]any, 0, len(histArgs)+len(args)+2)
	scanArgs = append(scanArgs, histArgs...)
	scanArgs = append(scanArgs, args...)
	scanArgs = append(scanArgs, limit, offset)

	var rows []balanceRow
	if err := r.db.WithContext(ctx).Raw(query, scanArgs...).Scan(&rows).Error; err != nil {
		return 0, nil, err
	}

	total, err := r.count(ctx, where, args)
	if err != nil {
		return 0, nil, err
	}

	results := make([]*repo.BalanceResult, len(rows))
	for i, row := range rows {
		results[i] = &repo.BalanceResult{
			EditionUid:       row.EditionUid,
			WarehouseUid:     row.WarehouseUid,
			WarehouseName:    row.WarehouseName,
			InboundQuantity:  row.InboundQuantity,
			OutboundQuantity: row.OutboundQuantity,
			AdjustQuantity:   row.AdjustQuantity,
			CurrentQuantity:  row.CurrentQuantity,
		}
	}
	return total, results, nil
}

func (r *stockBalanceRepository) buildFilter(opt option.BalanceQueryOption) (string, []any, error) {
	where := ""
	args := make([]any, 0, 2)
	if uid := opt.EditionUid(); uid != nil && *uid != "" {
		where += " AND s.edition_uid = ?"
		args = append(args, *uid)
	}
	if uid := opt.WarehouseUid(); uid != nil && *uid != "" {
		where += " AND s.warehouse_uid = ?"
		args = append(args, *uid)
	}
	return where, args, nil
}

func (r *stockBalanceRepository) historyJoin(opt option.BalanceQueryOption) (string, []any, error) {
	ym := opt.YearMonth()
	if ym == nil || *ym == "" {
		return "", nil, nil
	}
	start, err := time.Parse("2006-01", *ym)
	if err != nil {
		return "", nil, fmt.Errorf("invalid year_month %q (want YYYY-MM): %w", *ym, err)
	}
	end := start.AddDate(0, 1, 0)
	return " AND h.created_at >= ? AND h.created_at < ?", []any{start, end}, nil
}

func (r *stockBalanceRepository) count(ctx context.Context, where string, args []any) (int64, error) {
	query := fmt.Sprintf(`
SELECT COUNT(*)
FROM stock s
JOIN warehouse w ON w.id = s.warehouse_id
WHERE s.deleted_at IS NULL%s`, where)

	var total int64
	if err := r.db.WithContext(ctx).Raw(query, args...).Scan(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}
