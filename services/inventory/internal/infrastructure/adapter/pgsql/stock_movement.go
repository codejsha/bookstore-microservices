package pgsql

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/codejsha/bookstore-microservices/inventory/generated/infrastructure/port/entity"
)

const (
	changeTypeInbound    = "INBOUND"
	changeTypeOutbound   = "OUTBOUND"
	changeTypeAdjustment = "ADJUSTMENT"
)

const (
	transferStatusPending   = "PENDING"
	transferStatusCompleted = "COMPLETED"
	transferStatusCancelled = "CANCELLED"
)

const (
	auditStatusInProgress = "IN_PROGRESS"
	auditStatusCompleted  = "COMPLETED"
)

const closingStatusClosed = "CLOSED"

const historyBucketSelect = `
	COALESCE(SUM(CASE WHEN change_type IN ('INBOUND','RELEASE') THEN change_qty ELSE 0 END), 0) AS inbound,
	COALESCE(SUM(CASE WHEN change_type IN ('OUTBOUND','RESERVATION') THEN -change_qty ELSE 0 END), 0) AS outbound,
	COALESCE(SUM(CASE WHEN change_type = 'ADJUSTMENT' THEN change_qty ELSE 0 END), 0) AS adjust`

type stockDelta struct {
	editionID    int64
	editionUID   string
	warehouseID  int64
	warehouseUID string
	delta        int32
	changeType   string
	reason       *string
	allowCreate  bool
}

func applyStockDeltaTx(tx *gorm.DB, d stockDelta) error {
	now := time.Now()

	var st entity.StockEntity
	e := tx.Where("edition_uid = ? AND warehouse_uid = ?", d.editionUID, d.warehouseUID).
		First(&st).Error
	switch {
	case errors.Is(e, gorm.ErrRecordNotFound):
		if !d.allowCreate || d.delta < 0 {
			return fmt.Errorf("stock not found for edition %s in warehouse %s", d.editionUID, d.warehouseUID)
		}
		st = entity.StockEntity{
			Uid:          uuid.Must(uuid.NewV7()).String(),
			EditionId:    d.editionID,
			EditionUid:   d.editionUID,
			WarehouseId:  d.warehouseID,
			WarehouseUid: d.warehouseUID,
			Quantity:     0,
			CreatedAt:    now,
			Version:      1,
		}
		if err := tx.Create(&st).Error; err != nil {
			return err
		}
	case e != nil:
		return e
	}

	newQty := st.Quantity + d.delta
	if newQty < 0 {
		return fmt.Errorf("insufficient stock for edition %s in warehouse %s: have %d, requested %d",
			d.editionUID, d.warehouseUID, st.Quantity, -d.delta)
	}

	res := tx.Model(&entity.StockEntity{}).
		Where("id = ? AND version = ?", st.Id, st.Version).
		Updates(map[string]any{
			"quantity":   newQty,
			"version":    st.Version + 1,
			"updated_at": now,
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errVersionConflict
	}

	return tx.Create(&entity.StockHistoryEntity{
		Uid:        uuid.Must(uuid.NewV7()).String(),
		StockId:    st.Id,
		StockUid:   st.Uid,
		ChangeType: d.changeType,
		Reason:     d.reason,
		ChangeQty:  d.delta,
		CreatedAt:  now,
		Version:    1,
	}).Error
}
