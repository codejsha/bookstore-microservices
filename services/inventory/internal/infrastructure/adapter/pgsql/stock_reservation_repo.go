package pgsql

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/codejsha/shared-library-go/pkg/database"

	"github.com/codejsha/bookstore-microservices/inventory/generated/infrastructure/port/entity"
	"github.com/codejsha/bookstore-microservices/inventory/internal/application/port/repo"
)

const (
	reservationStatusReserved = "RESERVED"
	reservationStatusReleased = "RELEASED"

	changeTypeReservation = "RESERVATION"
	changeTypeRelease     = "RELEASE"

	reservationKeyConstraint = "uk_stock_reservation_key_warehouse"
)

type stockReservationEntity struct {
	Id             int64          `gorm:"column:id;primaryKey;autoIncrement:true"`
	Uid            string         `gorm:"column:uid"`
	ReservationKey string         `gorm:"column:reservation_key"`
	OrderRef       string         `gorm:"column:order_ref"`
	EditionId      int64          `gorm:"column:edition_id"`
	EditionUid     string         `gorm:"column:edition_uid"`
	WarehouseId    int64          `gorm:"column:warehouse_id"`
	WarehouseUid   string         `gorm:"column:warehouse_uid"`
	Quantity       int32          `gorm:"column:quantity"`
	Status         string         `gorm:"column:status"`
	ReleasedAt     *time.Time     `gorm:"column:released_at"`
	CreatedAt      time.Time      `gorm:"column:created_at"`
	UpdatedAt      *time.Time     `gorm:"column:updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at"`
	Actor          int64          `gorm:"column:actor"`
	Version        int64          `gorm:"column:version"`
}

func (*stockReservationEntity) TableName() string { return "stock_reservation" }

var _ repo.StockReservationRepo = (*stockReservationRepository)(nil)

type stockReservationRepository struct {
	db *gorm.DB
}

func NewStockReservationRepository(dataSource *database.DataSource) repo.StockReservationRepo {
	return &stockReservationRepository{db: dataSource.DB()}
}

type stockAllocation struct {
	stock    *entity.StockEntity
	quantity int32
}

func allocateReservation(rows []entity.StockEntity, editionId int64, quantity int32) ([]stockAllocation, error) {
	ordered := make([]entity.StockEntity, len(rows))
	copy(ordered, rows)
	sort.SliceStable(ordered, func(i, j int) bool {
		if ordered[i].Quantity != ordered[j].Quantity {
			return ordered[i].Quantity > ordered[j].Quantity
		}
		return ordered[i].Id < ordered[j].Id
	})

	remaining := quantity
	allocs := make([]stockAllocation, 0, len(ordered))
	for i := range ordered {
		if remaining == 0 {
			break
		}
		avail := ordered[i].Quantity
		if avail <= 0 {
			continue
		}
		take := avail
		if take > remaining {
			take = remaining
		}
		allocs = append(allocs, stockAllocation{stock: &ordered[i], quantity: take})
		remaining -= take
	}
	if remaining > 0 {
		return nil, fmt.Errorf("insufficient stock for edition_id %d: requested %d: %w", editionId, quantity, repo.ErrInsufficientStock)
	}
	return allocs, nil
}

func (r *stockReservationRepository) Reserve(ctx context.Context, p repo.ReserveParams) (*repo.ReservationResult, error) {
	for attempt := 0; attempt < maxOptimisticRetries; attempt++ {
		var out *repo.ReservationResult
		var retry bool

		err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			var existing []stockReservationEntity
			if e := tx.Where("reservation_key = ?", p.ReservationKey).Find(&existing).Error; e != nil {
				return e
			}
			if len(existing) > 0 {
				out = toReservationResultFromRows(existing, true)
				return nil
			}

			var rows []entity.StockEntity
			if err := tx.Where("edition_id = ?", p.EditionId).Find(&rows).Error; err != nil {
				return err
			}
			if len(rows) == 0 {
				return fmt.Errorf("no stock found for edition_id %d: %w", p.EditionId, repo.ErrInsufficientStock)
			}
			allocs, err := allocateReservation(rows, p.EditionId, p.Quantity)
			if err != nil {
				return err
			}

			now := time.Now()
			resvRows := make([]stockReservationEntity, 0, len(allocs))
			for _, a := range allocs {
				res := tx.Model(&entity.StockEntity{}).
					Where("id = ? AND version = ?", a.stock.Id, a.stock.Version).
					Updates(map[string]any{
						"quantity":   a.stock.Quantity - a.quantity,
						"version":    a.stock.Version + 1,
						"updated_at": now,
					})
				if res.Error != nil {
					return res.Error
				}
				if res.RowsAffected == 0 {
					retry = true
					return errVersionConflict
				}

				if err := tx.Create(&entity.StockHistoryEntity{
					Uid:        uuid.Must(uuid.NewV7()).String(),
					StockId:    a.stock.Id,
					StockUid:   a.stock.Uid,
					ChangeType: changeTypeReservation,
					Reason:     p.Reason,
					ChangeQty:  -a.quantity,
					CreatedAt:  now,
					Version:    1,
				}).Error; err != nil {
					return err
				}

				resv := stockReservationEntity{
					Uid:            uuid.Must(uuid.NewV7()).String(),
					ReservationKey: p.ReservationKey,
					OrderRef:       p.OrderRef,
					EditionId:      p.EditionId,
					EditionUid:     a.stock.EditionUid,
					WarehouseId:    a.stock.WarehouseId,
					WarehouseUid:   a.stock.WarehouseUid,
					Quantity:       a.quantity,
					Status:         reservationStatusReserved,
					CreatedAt:      now,
					Version:        1,
				}
				if err := tx.Create(&resv).Error; err != nil {
					if isUniqueViolation(err, reservationKeyConstraint) {
						return fmt.Errorf("%w: %v", errReservationConflict, err)
					}
					return err
				}
				resvRows = append(resvRows, resv)
			}

			out = toReservationResultFromRows(resvRows, false)
			return nil
		})

		switch {
		case err == nil:
			return out, nil
		case retry:
			continue
		case errors.Is(err, errReservationConflict):
			var existing []stockReservationEntity
			if e := r.db.WithContext(ctx).Where("reservation_key = ?", p.ReservationKey).Find(&existing).Error; e == nil && len(existing) > 0 {
				return toReservationResultFromRows(existing, true), nil
			}
			return nil, err
		default:
			return nil, err
		}
	}
	return nil, fmt.Errorf("reserve %q: exhausted %d optimistic-lock retries", p.ReservationKey, maxOptimisticRetries)
}

func (r *stockReservationRepository) Release(ctx context.Context, p repo.ReleaseParams) (*repo.ReservationResult, error) {
	for attempt := 0; attempt < maxOptimisticRetries; attempt++ {
		var out *repo.ReservationResult
		var retry bool

		err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			var resvRows []stockReservationEntity
			if e := tx.Where("reservation_key = ?", p.ReservationKey).Find(&resvRows).Error; e != nil {
				return e
			}
			if len(resvRows) == 0 {
				out = nil
				return nil
			}

			allReleased := true
			for i := range resvRows {
				if resvRows[i].Status != reservationStatusReleased {
					allReleased = false
					break
				}
			}
			if allReleased {
				out = toReservationResultFromRows(resvRows, true)
				return nil
			}

			now := time.Now()
			for i := range resvRows {
				resv := &resvRows[i]
				if resv.Status == reservationStatusReleased {
					continue
				}

				var target entity.StockEntity
				if err := tx.Where("edition_id = ? AND warehouse_id = ?", resv.EditionId, resv.WarehouseId).
					First(&target).Error; err != nil {
					return err
				}

				res := tx.Model(&entity.StockEntity{}).
					Where("id = ? AND version = ?", target.Id, target.Version).
					Updates(map[string]any{
						"quantity":   target.Quantity + resv.Quantity,
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

				if err := tx.Create(&entity.StockHistoryEntity{
					Uid:        uuid.Must(uuid.NewV7()).String(),
					StockId:    target.Id,
					StockUid:   target.Uid,
					ChangeType: changeTypeRelease,
					Reason:     p.Reason,
					ChangeQty:  resv.Quantity,
					CreatedAt:  now,
					Version:    1,
				}).Error; err != nil {
					return err
				}

				mark := tx.Model(&stockReservationEntity{}).
					Where("id = ? AND version = ?", resv.Id, resv.Version).
					Updates(map[string]any{
						"status":      reservationStatusReleased,
						"released_at": now,
						"version":     resv.Version + 1,
						"updated_at":  now,
					})
				if mark.Error != nil {
					return mark.Error
				}
				if mark.RowsAffected == 0 {
					retry = true
					return errVersionConflict
				}
				resv.Status = reservationStatusReleased
			}

			out = toReservationResultFromRows(resvRows, false)
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
	return nil, fmt.Errorf("release %q: exhausted %d optimistic-lock retries", p.ReservationKey, maxOptimisticRetries)
}

func toReservationResultFromRows(rows []stockReservationEntity, alreadyApplied bool) *repo.ReservationResult {
	if len(rows) == 0 {
		return nil
	}
	var total int32
	for i := range rows {
		total += rows[i].Quantity
	}
	first := &rows[0]
	return &repo.ReservationResult{
		Uid:            first.Uid,
		ReservationKey: first.ReservationKey,
		EditionUid:     first.EditionUid,
		WarehouseUid:   first.WarehouseUid,
		Quantity:       total,
		Status:         first.Status,
		AlreadyApplied: alreadyApplied,
	}
}
