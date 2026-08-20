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

var _ repo.StockAuditRepo = (*stockAuditRepository)(nil)

type stockAuditRepository struct {
	db *gorm.DB
	genrepo.StockAuditRepo
}

func NewStockAuditRepository(dataSource *database.DataSource) repo.StockAuditRepo {
	db := dataSource.DB()
	return &stockAuditRepository{
		db:             db,
		StockAuditRepo: genrepo.NewStockAuditRepo(db),
	}
}

func (r *stockAuditRepository) FindAll(ctx context.Context, opt option.AuditQueryOption) (int64, []*repo.AuditResult, error) {
	q := dao.Use(r.db)
	sortable := map[string]field.OrderExpr{
		"uid":          q.StockAuditEntity.Uid,
		"warehouseuid": q.StockAuditEntity.WarehouseUid,
		"status":       q.StockAuditEntity.Status,
		"note":         q.StockAuditEntity.Note,
		"completedat":  q.StockAuditEntity.CompletedAt,
		"createdat":    q.StockAuditEntity.CreatedAt,
		"updatedat":    q.StockAuditEntity.UpdatedAt,
	}

	base := q.StockAuditEntity.WithContext(ctx).Scopes(r.buildWhereScope(q, opt))

	audits, err := base.
		Scopes(gormutils.BuildPageScope(opt.Page(), sortable)).
		Find()
	if err != nil {
		return 0, nil, err
	}

	total, err := base.Count()
	if err != nil {
		return 0, nil, err
	}

	results := make([]*repo.AuditResult, len(audits))
	for i, a := range audits {
		items, err := r.loadItems(ctx, r.db, a.Id)
		if err != nil {
			return 0, nil, err
		}
		results[i] = toAuditResult(a, items)
	}
	return total, results, nil
}

func (r *stockAuditRepository) buildWhereScope(q *dao.Query, opt option.AuditQueryOption) func(gen.Dao) gen.Dao {
	conds := make([]gen.Condition, 0)
	if uid := opt.WarehouseUid(); uid != nil && *uid != "" {
		conds = append(conds, q.StockAuditEntity.WarehouseUid.Eq(*uid))
	}
	if status := opt.Status(); status != nil && *status != "" {
		conds = append(conds, q.StockAuditEntity.Status.Eq(*status))
	}

	return func(dao gen.Dao) gen.Dao {
		if len(conds) > 0 {
			return dao.Where(conds...)
		}
		return dao
	}
}

func (r *stockAuditRepository) FindByUid(ctx context.Context, uid string) (*repo.AuditResult, error) {
	var a entity.StockAuditEntity
	err := r.db.WithContext(ctx).Where("uid = ?", uid).First(&a).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	items, err := r.loadItems(ctx, r.db, a.Id)
	if err != nil {
		return nil, err
	}
	return toAuditResult(&a, items), nil
}

func (r *stockAuditRepository) Create(ctx context.Context, p repo.AuditCreateParams) (*repo.AuditResult, error) {
	if len(p.Items) == 0 {
		return nil, fmt.Errorf("audit must contain at least one item")
	}

	var out *repo.AuditResult
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var warehouse entity.WarehouseEntity
		e := tx.Where("uid = ?", p.WarehouseUid).First(&warehouse).Error
		if errors.Is(e, gorm.ErrRecordNotFound) {
			return fmt.Errorf("warehouse %s not found", p.WarehouseUid)
		}
		if e != nil {
			return e
		}

		now := time.Now()
		audit := &entity.StockAuditEntity{
			Uid:          uuid.Must(uuid.NewV7()).String(),
			WarehouseId:  warehouse.Id,
			WarehouseUid: warehouse.Uid,
			Status:       auditStatusInProgress,
			Note:         p.Notes,
			CreatedAt:    now,
			Version:      1,
		}
		if err := tx.Create(audit).Error; err != nil {
			return err
		}

		items := make([]entity.StockAuditItemEntity, 0, len(p.Items))
		for _, it := range p.Items {
			var stock entity.StockEntity
			e := tx.Where("edition_uid = ? AND warehouse_uid = ?", it.EditionUid, warehouse.Uid).
				First(&stock).Error
			if errors.Is(e, gorm.ErrRecordNotFound) {
				return fmt.Errorf("no stock for edition %s in warehouse %s", it.EditionUid, warehouse.Uid)
			}
			if e != nil {
				return e
			}
			item := entity.StockAuditItemEntity{
				AuditId:        audit.Id,
				EditionId:      stock.EditionId,
				EditionUid:     stock.EditionUid,
				SystemQuantity: stock.Quantity,
				ActualQuantity: it.ActualQuantity,
				Difference:     it.ActualQuantity - stock.Quantity,
			}
			if err := tx.Create(&item).Error; err != nil {
				return err
			}
			items = append(items, item)
		}
		out = toAuditResult(audit, items)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (r *stockAuditRepository) Complete(ctx context.Context, uid string) (*repo.AuditResult, error) {
	for attempt := 0; attempt < maxOptimisticRetries; attempt++ {
		var out *repo.AuditResult
		var missing, retry bool

		err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			var audit entity.StockAuditEntity
			e := tx.Where("uid = ?", uid).First(&audit).Error
			if errors.Is(e, gorm.ErrRecordNotFound) {
				missing = true
				return nil
			}
			if e != nil {
				return e
			}
			if audit.Status == auditStatusCompleted {
				return fmt.Errorf("audit %s is already completed", uid)
			}

			items, err := r.loadItems(ctx, tx, audit.Id)
			if err != nil {
				return err
			}
			for _, it := range items {
				if it.Difference == 0 {
					continue
				}
				reason := fmt.Sprintf("audit %s", audit.Uid)
				if err := applyStockDeltaTx(tx, stockDelta{
					editionID:    it.EditionId,
					editionUID:   it.EditionUid,
					warehouseID:  audit.WarehouseId,
					warehouseUID: audit.WarehouseUid,
					delta:        it.Difference,
					changeType:   changeTypeAdjustment,
					reason:       &reason,
					allowCreate:  true,
				}); err != nil {
					if errors.Is(err, errVersionConflict) {
						retry = true
					}
					return err
				}
			}

			now := time.Now()
			res := tx.Model(&entity.StockAuditEntity{}).
				Where("id = ? AND version = ?", audit.Id, audit.Version).
				Updates(map[string]any{
					"status":       auditStatusCompleted,
					"completed_at": now,
					"updated_at":   now,
					"version":      audit.Version + 1,
				})
			if res.Error != nil {
				return res.Error
			}
			if res.RowsAffected == 0 {
				retry = true
				return errVersionConflict
			}
			audit.Status = auditStatusCompleted
			audit.CompletedAt = &now
			audit.UpdatedAt = &now
			out = toAuditResult(&audit, items)
			return nil
		})

		switch {
		case err == nil && missing:
			return nil, nil
		case err == nil:
			return out, nil
		case retry:
			continue
		default:
			return nil, err
		}
	}
	return nil, fmt.Errorf("complete audit %s: exhausted %d optimistic-lock retries", uid, maxOptimisticRetries)
}

func (r *stockAuditRepository) loadItems(ctx context.Context, db *gorm.DB, auditID int64) ([]entity.StockAuditItemEntity, error) {
	var items []entity.StockAuditItemEntity
	if err := db.WithContext(ctx).Where("audit_id = ?", auditID).Order("id").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func toAuditResult(a *entity.StockAuditEntity, items []entity.StockAuditItemEntity) *repo.AuditResult {
	out := make([]repo.AuditItemResult, len(items))
	for i, it := range items {
		out[i] = repo.AuditItemResult{
			EditionUid:     it.EditionUid,
			SystemQuantity: it.SystemQuantity,
			ActualQuantity: it.ActualQuantity,
			Difference:     it.Difference,
		}
	}
	return &repo.AuditResult{
		Uid:          a.Uid,
		WarehouseUid: a.WarehouseUid,
		Status:       a.Status,
		Items:        out,
		Notes:        a.Note,
		CompletedAt:  a.CompletedAt,
		CreatedAt:    a.CreatedAt,
		UpdatedAt:    a.UpdatedAt,
	}
}
