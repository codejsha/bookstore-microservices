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

var _ repo.StockTransferRepo = (*stockTransferRepository)(nil)

type stockTransferRepository struct {
	db *gorm.DB
	genrepo.StockTransferRepo
}

func NewStockTransferRepository(dataSource *database.DataSource) repo.StockTransferRepo {
	db := dataSource.DB()
	return &stockTransferRepository{
		db:                db,
		StockTransferRepo: genrepo.NewStockTransferRepo(db),
	}
}

func (r *stockTransferRepository) FindAll(ctx context.Context, opt option.TransferQueryOption) (int64, []*repo.TransferResult, error) {
	q := dao.Use(r.db)
	sortable := map[string]field.OrderExpr{
		"uid":                q.StockTransferEntity.Uid,
		"editionuid":         q.StockTransferEntity.EditionUid,
		"sourcewarehouseuid": q.StockTransferEntity.SourceWarehouseUid,
		"targetwarehouseuid": q.StockTransferEntity.TargetWarehouseUid,
		"quantity":           q.StockTransferEntity.Quantity,
		"status":             q.StockTransferEntity.Status,
		"completedat":        q.StockTransferEntity.CompletedAt,
		"createdat":          q.StockTransferEntity.CreatedAt,
		"updatedat":          q.StockTransferEntity.UpdatedAt,
	}

	base := q.StockTransferEntity.WithContext(ctx).Scopes(r.buildWhereScope(q, opt))

	entities, err := base.
		Scopes(gormutils.BuildPageScope(opt.Page(), sortable)).
		Find()
	if err != nil {
		return 0, nil, err
	}

	total, err := base.Count()
	if err != nil {
		return 0, nil, err
	}

	results := make([]*repo.TransferResult, len(entities))
	for i, e := range entities {
		results[i] = toTransferResult(e)
	}
	return total, results, nil
}

func (r *stockTransferRepository) buildWhereScope(q *dao.Query, opt option.TransferQueryOption) func(gen.Dao) gen.Dao {
	conds := make([]gen.Condition, 0)
	if uid := opt.EditionUid(); uid != nil && *uid != "" {
		conds = append(conds, q.StockTransferEntity.EditionUid.Eq(*uid))
	}
	if uid := opt.SourceWarehouseUid(); uid != nil && *uid != "" {
		conds = append(conds, q.StockTransferEntity.SourceWarehouseUid.Eq(*uid))
	}
	if uid := opt.TargetWarehouseUid(); uid != nil && *uid != "" {
		conds = append(conds, q.StockTransferEntity.TargetWarehouseUid.Eq(*uid))
	}
	if status := opt.Status(); status != nil && *status != "" {
		conds = append(conds, q.StockTransferEntity.Status.Eq(*status))
	}

	return func(dao gen.Dao) gen.Dao {
		if len(conds) > 0 {
			return dao.Where(conds...)
		}
		return dao
	}
}

func (r *stockTransferRepository) FindByUid(ctx context.Context, uid string) (*repo.TransferResult, error) {
	e, err := r.fetchByUid(ctx, r.db, uid)
	if err != nil {
		return nil, err
	}
	if e == nil {
		return nil, nil
	}
	return toTransferResult(e), nil
}

func (r *stockTransferRepository) Create(ctx context.Context, p repo.TransferCreateParams) (*repo.TransferResult, error) {
	if p.SourceWarehouseUid == p.TargetWarehouseUid {
		return nil, fmt.Errorf("transfer source and target warehouses must differ")
	}

	for attempt := 0; attempt < maxOptimisticRetries; attempt++ {
		var out *repo.TransferResult
		var retry bool

		err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			var source entity.StockEntity
			e := tx.Where("edition_uid = ? AND warehouse_uid = ?", p.EditionUid, p.SourceWarehouseUid).
				First(&source).Error
			if errors.Is(e, gorm.ErrRecordNotFound) {
				return fmt.Errorf("no stock for edition %s in source warehouse %s: %w", p.EditionUid, p.SourceWarehouseUid, repo.ErrStockNotFound)
			}
			if e != nil {
				return e
			}

			var target entity.WarehouseEntity
			e = tx.Where("uid = ?", p.TargetWarehouseUid).First(&target).Error
			if errors.Is(e, gorm.ErrRecordNotFound) {
				return fmt.Errorf("target warehouse %s not found: %w", p.TargetWarehouseUid, repo.ErrWarehouseNotFound)
			}
			if e != nil {
				return e
			}

			if err := applyStockDeltaTx(tx, stockDelta{
				editionID:    source.EditionId,
				editionUID:   source.EditionUid,
				warehouseID:  source.WarehouseId,
				warehouseUID: source.WarehouseUid,
				delta:        -p.Quantity,
				changeType:   changeTypeOutbound,
				reason:       p.Reason,
			}); err != nil {
				if isRetryableConflict(err) {
					retry = true
				}
				return err
			}

			now := time.Now()
			transfer := &entity.StockTransferEntity{
				Uid:                uuid.Must(uuid.NewV7()).String(),
				EditionId:          source.EditionId,
				EditionUid:         source.EditionUid,
				SourceWarehouseId:  source.WarehouseId,
				SourceWarehouseUid: source.WarehouseUid,
				TargetWarehouseId:  target.Id,
				TargetWarehouseUid: target.Uid,
				Quantity:           p.Quantity,
				Status:             transferStatusPending,
				Reason:             p.Reason,
				CreatedAt:          now,
				Version:            1,
			}
			if err := tx.Create(transfer).Error; err != nil {
				return err
			}
			out = toTransferResult(transfer)
			return nil
		})

		switch {
		case err == nil:
			return out, nil
		case retry:
			continue
		default:
			return nil, err
		}
	}
	return nil, fmt.Errorf("create transfer for edition %s: exhausted %d optimistic-lock retries", p.EditionUid, maxOptimisticRetries)
}

func (r *stockTransferRepository) Complete(ctx context.Context, uid string) (*repo.TransferResult, error) {
	return r.finalize(ctx, uid, func(tx *gorm.DB, t *entity.StockTransferEntity) error {
		if t.Status != transferStatusPending {
			return fmt.Errorf("transfer %s cannot be completed in status %s", uid, t.Status)
		}
		if err := applyStockDeltaTx(tx, stockDelta{
			editionID:    t.EditionId,
			editionUID:   t.EditionUid,
			warehouseID:  t.TargetWarehouseId,
			warehouseUID: t.TargetWarehouseUid,
			delta:        t.Quantity,
			changeType:   changeTypeInbound,
			reason:       t.Reason,
			allowCreate:  true,
		}); err != nil {
			return err
		}
		now := time.Now()
		return r.updateStatus(tx, t, transferStatusCompleted, &now, &now)
	})
}

func (r *stockTransferRepository) Cancel(ctx context.Context, uid string) (*repo.TransferResult, error) {
	return r.finalize(ctx, uid, func(tx *gorm.DB, t *entity.StockTransferEntity) error {
		if t.Status != transferStatusPending {
			return fmt.Errorf("transfer %s cannot be cancelled in status %s", uid, t.Status)
		}
		if err := applyStockDeltaTx(tx, stockDelta{
			editionID:    t.EditionId,
			editionUID:   t.EditionUid,
			warehouseID:  t.SourceWarehouseId,
			warehouseUID: t.SourceWarehouseUid,
			delta:        t.Quantity,
			changeType:   changeTypeInbound,
			reason:       t.Reason,
			allowCreate:  true,
		}); err != nil {
			return err
		}
		now := time.Now()
		return r.updateStatus(tx, t, transferStatusCancelled, nil, &now)
	})
}

func (r *stockTransferRepository) finalize(
	ctx context.Context,
	uid string,
	apply func(tx *gorm.DB, t *entity.StockTransferEntity) error,
) (*repo.TransferResult, error) {
	for attempt := 0; attempt < maxOptimisticRetries; attempt++ {
		var out *repo.TransferResult
		var missing, retry bool

		err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			var t entity.StockTransferEntity
			e := tx.Where("uid = ?", uid).First(&t).Error
			if errors.Is(e, gorm.ErrRecordNotFound) {
				missing = true
				return nil
			}
			if e != nil {
				return e
			}
			if err := apply(tx, &t); err != nil {
				if isRetryableConflict(err) {
					retry = true
				}
				return err
			}
			out = toTransferResult(&t)
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
	return nil, fmt.Errorf("finalize transfer %s: exhausted %d optimistic-lock retries", uid, maxOptimisticRetries)
}

func (r *stockTransferRepository) updateStatus(
	tx *gorm.DB,
	t *entity.StockTransferEntity,
	status string,
	completedAt *time.Time,
	updatedAt *time.Time,
) error {
	res := tx.Model(&entity.StockTransferEntity{}).
		Where("id = ? AND version = ?", t.Id, t.Version).
		Updates(map[string]any{
			"status":       status,
			"completed_at": completedAt,
			"updated_at":   updatedAt,
			"version":      t.Version + 1,
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errVersionConflict
	}
	t.Status = status
	t.CompletedAt = completedAt
	t.UpdatedAt = updatedAt
	t.Version++
	return nil
}

func (r *stockTransferRepository) fetchByUid(ctx context.Context, db *gorm.DB, uid string) (*entity.StockTransferEntity, error) {
	var t entity.StockTransferEntity
	err := db.WithContext(ctx).Where("uid = ?", uid).First(&t).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func toTransferResult(e *entity.StockTransferEntity) *repo.TransferResult {
	return &repo.TransferResult{
		Uid:                e.Uid,
		EditionUid:         e.EditionUid,
		SourceWarehouseUid: e.SourceWarehouseUid,
		TargetWarehouseUid: e.TargetWarehouseUid,
		Quantity:           e.Quantity,
		Status:             e.Status,
		Reason:             e.Reason,
		CompletedAt:        e.CompletedAt,
		CreatedAt:          e.CreatedAt,
		UpdatedAt:          e.UpdatedAt,
	}
}
