package pgsql

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"

	"github.com/codejsha/shared-library-go/pkg/database"
	"github.com/codejsha/shared-library-go/pkg/database/gormutils"

	"github.com/codejsha/bookstore-microservices/inventory/generated/infrastructure/port/dao"
	"github.com/codejsha/bookstore-microservices/inventory/generated/infrastructure/port/entity"
	genrepo "github.com/codejsha/bookstore-microservices/inventory/generated/infrastructure/port/repo"
	"github.com/codejsha/bookstore-microservices/inventory/internal/application/port/repo"
	"github.com/codejsha/bookstore-microservices/inventory/internal/domain/model/option"
)

var _ repo.MonthlyClosingRepo = (*monthlyClosingRepository)(nil)

type monthlyClosingRepository struct {
	db *gorm.DB
	genrepo.MonthlyClosingRepo
}

func NewMonthlyClosingRepository(dataSource *database.DataSource) repo.MonthlyClosingRepo {
	db := dataSource.DB()
	return &monthlyClosingRepository{
		db:                 db,
		MonthlyClosingRepo: genrepo.NewMonthlyClosingRepo(db),
	}
}

func (r *monthlyClosingRepository) FindAll(ctx context.Context, opt option.ClosingQueryOption) (int64, []*repo.ClosingResult, error) {
	q := dao.Use(r.db)
	sortable := map[string]field.OrderExpr{
		"uid":          q.MonthlyClosingEntity.Uid,
		"warehouseuid": q.MonthlyClosingEntity.WarehouseUid,
		"year":         q.MonthlyClosingEntity.Year,
		"month":        q.MonthlyClosingEntity.Month,
		"status":       q.MonthlyClosingEntity.Status,
		"closedat":     q.MonthlyClosingEntity.ClosedAt,
		"createdat":    q.MonthlyClosingEntity.CreatedAt,
		"updatedat":    q.MonthlyClosingEntity.UpdatedAt,
	}

	base := q.MonthlyClosingEntity.WithContext(ctx).Scopes(r.buildWhereScope(q, opt))

	closings, err := base.
		Scopes(gormutils.BuildPageScope(opt.Page(), sortable)).
		Find()
	if err != nil {
		return 0, nil, err
	}

	total, err := base.Count()
	if err != nil {
		return 0, nil, err
	}

	results := make([]*repo.ClosingResult, len(closings))
	for i, c := range closings {
		items, err := r.loadItems(ctx, r.db, c.Id)
		if err != nil {
			return 0, nil, err
		}
		results[i] = toClosingResult(c, items)
	}
	return total, results, nil
}

func (r *monthlyClosingRepository) buildWhereScope(q *dao.Query, opt option.ClosingQueryOption) func(gen.Dao) gen.Dao {
	conds := make([]gen.Condition, 0)
	if uid := opt.WarehouseUid(); uid != nil && *uid != "" {
		conds = append(conds, q.MonthlyClosingEntity.WarehouseUid.Eq(*uid))
	}
	if year := opt.Year(); year != nil {
		conds = append(conds, q.MonthlyClosingEntity.Year.Eq(*year))
	}
	if month := opt.Month(); month != nil {
		conds = append(conds, q.MonthlyClosingEntity.Month.Eq(*month))
	}
	if status := opt.Status(); status != nil && *status != "" {
		conds = append(conds, q.MonthlyClosingEntity.Status.Eq(*status))
	}

	return func(dao gen.Dao) gen.Dao {
		if len(conds) > 0 {
			return dao.Where(conds...)
		}
		return dao
	}
}

func (r *monthlyClosingRepository) FindByUid(ctx context.Context, uid string) (*repo.ClosingResult, error) {
	var c entity.MonthlyClosingEntity
	err := r.db.WithContext(ctx).Where("uid = ?", uid).First(&c).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	items, err := r.loadItems(ctx, r.db, c.Id)
	if err != nil {
		return nil, err
	}
	return toClosingResult(&c, items), nil
}

func (r *monthlyClosingRepository) Create(ctx context.Context, p repo.ClosingCreateParams) (*repo.ClosingResult, error) {
	if p.Month < 1 || p.Month > 12 {
		return nil, fmt.Errorf("invalid closing month: %d", p.Month)
	}
	start := time.Date(int(p.Year), time.Month(p.Month), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0)

	var out *repo.ClosingResult
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var warehouse entity.WarehouseEntity
		e := tx.Where("uid = ?", p.WarehouseUid).First(&warehouse).Error
		if errors.Is(e, gorm.ErrRecordNotFound) {
			return fmt.Errorf("warehouse %s not found", p.WarehouseUid)
		}
		if e != nil {
			return e
		}

		var stocks []entity.StockEntity
		if err := tx.Where("warehouse_id = ?", warehouse.Id).Order("id").Find(&stocks).Error; err != nil {
			return err
		}

		now := time.Now()
		closing := &entity.MonthlyClosingEntity{
			Uid:          uuid.Must(uuid.NewV7()).String(),
			WarehouseId:  warehouse.Id,
			WarehouseUid: warehouse.Uid,
			Year:         p.Year,
			Month:        p.Month,
			Status:       closingStatusClosed,
			ClosedAt:     &now,
			CreatedAt:    now,
			Version:      1,
		}
		if err := tx.Create(closing).Error; err != nil {
			if isUniqueViolation(err, "uq_closing_warehouse_period") {
				return fmt.Errorf("closing already exists for warehouse %s %d-%02d: %w",
					p.WarehouseUid, p.Year, p.Month, repo.ErrDuplicateClosing)
			}
			return err
		}

		items := make([]entity.MonthlyClosingItemEntity, 0, len(stocks))
		for _, st := range stocks {
			sums, err := r.periodMovement(tx, st.Id, start, end)
			if err != nil {
				return err
			}
			net := sums.inbound - sums.outbound + sums.adjust
			item := entity.MonthlyClosingItemEntity{
				ClosingId:        closing.Id,
				EditionId:        st.EditionId,
				EditionUid:       st.EditionUid,
				OpeningQuantity:  st.Quantity - net,
				InboundQuantity:  sums.inbound,
				OutboundQuantity: sums.outbound,
				AdjustQuantity:   sums.adjust,
				ClosingQuantity:  st.Quantity,
			}
			if err := tx.Create(&item).Error; err != nil {
				return err
			}
			items = append(items, item)
		}
		out = toClosingResult(closing, items)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

type movementSums struct {
	inbound  int32
	outbound int32
	adjust   int32
}

func (r *monthlyClosingRepository) periodMovement(tx *gorm.DB, stockID int64, start, end time.Time) (movementSums, error) {
	var row struct {
		Inbound  int32
		Outbound int32
		Adjust   int32
	}
	err := tx.Model(&entity.StockHistoryEntity{}).
		Select(historyBucketSelect).
		Where("stock_id = ? AND created_at >= ? AND created_at < ?", stockID, start, end).
		Scan(&row).Error
	if err != nil {
		return movementSums{}, err
	}
	return movementSums{inbound: row.Inbound, outbound: row.Outbound, adjust: row.Adjust}, nil
}

func (r *monthlyClosingRepository) loadItems(ctx context.Context, db *gorm.DB, closingID int64) ([]entity.MonthlyClosingItemEntity, error) {
	var items []entity.MonthlyClosingItemEntity
	if err := db.WithContext(ctx).Where("closing_id = ?", closingID).Order("id").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func toClosingResult(c *entity.MonthlyClosingEntity, items []entity.MonthlyClosingItemEntity) *repo.ClosingResult {
	out := make([]repo.ClosingItemResult, len(items))
	for i, it := range items {
		out[i] = repo.ClosingItemResult{
			EditionUid:       it.EditionUid,
			OpeningQuantity:  it.OpeningQuantity,
			InboundQuantity:  it.InboundQuantity,
			OutboundQuantity: it.OutboundQuantity,
			AdjustQuantity:   it.AdjustQuantity,
			ClosingQuantity:  it.ClosingQuantity,
		}
	}
	return &repo.ClosingResult{
		Uid:          c.Uid,
		WarehouseUid: c.WarehouseUid,
		Year:         c.Year,
		Month:        c.Month,
		Status:       c.Status,
		Items:        out,
		ClosedAt:     c.ClosedAt,
		CreatedAt:    c.CreatedAt,
		UpdatedAt:    c.UpdatedAt,
	}
}
