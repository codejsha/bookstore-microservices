package pgsql

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gen"
	"gorm.io/gorm"

	"github.com/codejsha/shared-library-go/pkg/database"

	"github.com/codejsha/bookstore-microservices/inventory/generated/infrastructure/port/dao"
	"github.com/codejsha/bookstore-microservices/inventory/generated/infrastructure/port/entity"
	genrepo "github.com/codejsha/bookstore-microservices/inventory/generated/infrastructure/port/repo"
	"github.com/codejsha/bookstore-microservices/inventory/internal/application/port/repo"
	"github.com/codejsha/bookstore-microservices/inventory/internal/domain/model/option"
)

var _ repo.StockRepo = (*stockRepository)(nil)

type stockRepository struct {
	db *gorm.DB
	genrepo.StockRepo
}

func NewStockRepository(dataSource *database.DataSource) repo.StockRepo {
	db := dataSource.DB()
	return &stockRepository{
		db:        db,
		StockRepo: genrepo.NewStockRepo(db),
	}
}

type stockRow struct {
	Id            int64  `gorm:"column:id"`
	Uid           string `gorm:"column:uid"`
	EditionUid    string `gorm:"column:edition_uid"`
	WarehouseUid  string `gorm:"column:warehouse_uid"`
	WarehouseName string `gorm:"column:warehouse_name"`
	Quantity      int32  `gorm:"column:quantity"`
}

func (r stockRepository) FindAll(ctx context.Context, opt option.StockQueryOption) (int64, []*repo.StockResult, error) {
	q := dao.Use(r.db)
	page := opt.Page()

	total, err := q.StockEntity.WithContext(ctx).
		Scopes(r.buildWhereScope(q, opt)).
		Distinct(q.StockEntity.EditionUid).
		Count()
	if err != nil {
		return 0, nil, err
	}

	editionQuery := q.StockEntity.WithContext(ctx).
		Scopes(r.buildWhereScope(q, opt)).
		Group(q.StockEntity.EditionUid).
		Order(q.StockEntity.EditionUid)
	if limit := page.GetSize(); limit > 0 {
		editionQuery = editionQuery.Limit(int(limit))
	}
	if offset := page.GetPage(); offset > 0 {
		editionQuery = editionQuery.Offset(int(offset))
	}
	var editionUids []string
	if err := editionQuery.Pluck(q.StockEntity.EditionUid, &editionUids); err != nil {
		return 0, nil, err
	}
	if len(editionUids) == 0 {
		return total, nil, nil
	}

	var rows []stockRow
	err = q.StockEntity.WithContext(ctx).
		LeftJoin(q.WarehouseEntity, q.WarehouseEntity.Uid.EqCol(q.StockEntity.WarehouseUid)).
		Scopes(r.buildWhereScope(q, opt)).
		Where(q.StockEntity.EditionUid.In(editionUids...)).
		Select(
			q.StockEntity.Id,
			q.StockEntity.Uid,
			q.StockEntity.EditionUid,
			q.StockEntity.WarehouseUid,
			q.WarehouseEntity.Name.As("warehouse_name"),
			q.StockEntity.Quantity,
		).
		Order(q.StockEntity.EditionUid, q.StockEntity.Id).
		Scan(&rows)
	if err != nil {
		return 0, nil, err
	}

	results := make([]*repo.StockResult, len(rows))
	for i, row := range rows {
		results[i] = &repo.StockResult{
			Id:            row.Id,
			Uid:           row.Uid,
			EditionUid:    row.EditionUid,
			WarehouseUid:  row.WarehouseUid,
			WarehouseName: row.WarehouseName,
			Quantity:      row.Quantity,
		}
	}

	return total, results, nil
}

func (r stockRepository) buildWhereScope(q *dao.Query, opt option.StockQueryOption) func(gen.Dao) gen.Dao {
	conds := make([]gen.Condition, 0)
	if uid := opt.EditionUid(); uid != nil && *uid != "" {
		conds = append(conds, q.StockEntity.EditionUid.Eq(*uid))
	}
	if uid := opt.WarehouseUid(); uid != nil && *uid != "" {
		conds = append(conds, q.StockEntity.WarehouseUid.Eq(*uid))
	}

	return func(dao gen.Dao) gen.Dao {
		if len(conds) > 0 {
			return dao.Where(conds...)
		}
		return dao
	}
}

func (r stockRepository) ApplyChange(ctx context.Context, p repo.ApplyChangeParams) error {
	for attempt := 0; attempt < maxOptimisticRetries; attempt++ {
		var retry bool
		err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			var target entity.StockEntity
			e := tx.Where("edition_uid = ? AND warehouse_uid = ?", p.EditionUid, p.WarehouseUid).
				First(&target).Error
			if errors.Is(e, gorm.ErrRecordNotFound) {
				return fmt.Errorf("stock not found for edition %s in warehouse %s: %w", p.EditionUid, p.WarehouseUid, repo.ErrStockNotFound)
			}
			if e != nil {
				return e
			}

			newQty := target.Quantity + p.Delta
			if newQty < 0 {
				return fmt.Errorf("insufficient stock for edition %s in warehouse %s: have %d, requested %d: %w",
					p.EditionUid, p.WarehouseUid, target.Quantity, -p.Delta, repo.ErrInsufficientStock)
			}

			now := time.Now()
			res := tx.Model(&entity.StockEntity{}).
				Where("id = ? AND version = ?", target.Id, target.Version).
				Updates(map[string]any{
					"quantity":   newQty,
					"version":    target.Version + 1,
					"updated_at": now,
				})
			if res.Error != nil {
				return res.Error
			}
			if res.RowsAffected == 0 {
				retry = true
				return errVersionConflict
			}

			return tx.Create(&entity.StockHistoryEntity{
				Uid:        uuid.Must(uuid.NewV7()).String(),
				StockId:    target.Id,
				StockUid:   target.Uid,
				ChangeType: p.ChangeType,
				Reason:     p.Reason,
				ChangeQty:  p.Delta,
				CreatedAt:  now,
				Version:    1,
			}).Error
		})
		if err == nil {
			return nil
		}
		if retry {
			continue
		}
		return err
	}
	return fmt.Errorf("apply change for edition %s: exhausted %d optimistic-lock retries", p.EditionUid, maxOptimisticRetries)
}

func (r stockRepository) FindOne(ctx context.Context, id int64) (*repo.StockResult, error) {
	q := dao.Use(r.db)

	var row stockRow
	err := q.StockEntity.WithContext(ctx).
		LeftJoin(q.WarehouseEntity, q.WarehouseEntity.Uid.EqCol(q.StockEntity.WarehouseUid)).
		Where(q.StockEntity.Id.Eq(id)).
		Select(
			q.StockEntity.Id,
			q.StockEntity.Uid,
			q.StockEntity.EditionUid,
			q.StockEntity.WarehouseUid,
			q.WarehouseEntity.Name.As("warehouse_name"),
			q.StockEntity.Quantity,
		).
		Scan(&row)
	if err != nil {
		return nil, err
	}
	if row.Id == 0 {
		return nil, nil
	}
	return &repo.StockResult{
		Id:            row.Id,
		Uid:           row.Uid,
		EditionUid:    row.EditionUid,
		WarehouseUid:  row.WarehouseUid,
		WarehouseName: row.WarehouseName,
		Quantity:      row.Quantity,
	}, nil
}
