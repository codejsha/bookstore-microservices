package pgsql

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/codejsha/shared-library-go/pkg/database"
	"github.com/codejsha/shared-library-go/pkg/pagination"

	"github.com/codejsha/bookstore-microservices/customer/generated/infrastructure/port/dao"
	"github.com/codejsha/bookstore-microservices/customer/generated/infrastructure/port/entity"
	genrepo "github.com/codejsha/bookstore-microservices/customer/generated/infrastructure/port/repo"
	"github.com/codejsha/bookstore-microservices/customer/internal/application/port/repo"
	"github.com/codejsha/bookstore-microservices/customer/internal/infrastructure/support/utils"
)

var (
	_ repo.PointRepo        = (*pointRepository)(nil)
	_ repo.PointHistoryRepo = (*pointHistoryRepository)(nil)
)

type pointRepository struct {
	q *dao.Query
	genrepo.CustomerPointRepo
}

func NewPointRepository(dataSource *database.DataSource) repo.PointRepo {
	db := dataSource.DB()
	return &pointRepository{
		q:                 dao.Use(db),
		CustomerPointRepo: genrepo.NewCustomerPointRepo(db),
	}
}

func (r *pointRepository) FindOne(ctx context.Context, id int64) (*repo.PointResult, error) {
	p := r.q.CustomerPointEntity
	e, err := p.WithContext(ctx).Where(p.Id.Eq(id)).First()
	if err != nil {
		return nil, err
	}
	return toPointResult(e), nil
}

func (r *pointRepository) FindByUserUid(ctx context.Context, userUid string) (*repo.PointResult, error) {
	p := r.q.CustomerPointEntity
	e, err := p.WithContext(ctx).Where(p.UserUid.Eq(userUid)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return toPointResult(e), nil
}

func (r *pointRepository) EnsureBalance(ctx context.Context, userUid string) error {
	ent := &entity.CustomerPointEntity{
		Uid:       uuid.Must(uuid.NewV7()).String(),
		UserUid:   userUid,
		Balance:   0,
		CreatedAt: time.Now(),
		Version:   1,
	}
	return r.q.CustomerPointEntity.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_uid"}},
			DoNothing: true,
		}).
		Create(ent)
}

func (r *pointRepository) ApplyPointChange(
	ctx context.Context,
	userUid string,
	delta int32,
	changeType string,
	reason *string,
) (*repo.PointResult, error) {
	var result *repo.PointResult
	err := r.q.Transaction(func(tx *dao.Query) error {
		p := tx.CustomerPointEntity
		e, err := p.WithContext(ctx).
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where(p.UserUid.Eq(userUid)).
			First()
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil
			}
			return err
		}

		newBalance := e.Balance + delta
		if newBalance < 0 {
			return fmt.Errorf("%w: have %d, requested %d", repo.ErrInsufficientPoints, e.Balance, -delta)
		}

		now := time.Now()
		e.Balance = newBalance
		e.UpdatedAt = &now
		if err := p.WithContext(ctx).Save(e); err != nil {
			return err
		}

		history := &entity.CustomerPointHistoryEntity{
			Uid:        uuid.Must(uuid.NewV7()).String(),
			UserId:     e.UserId,
			UserUid:    e.UserUid,
			ChangeType: changeType,
			Amount:     delta,
			Reason:     reason,
			CreatedAt:  now,
		}
		if err := tx.CustomerPointHistoryEntity.WithContext(ctx).Create(history); err != nil {
			return err
		}

		result = toPointResult(e)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

type pointHistoryRepository struct {
	q *dao.Query
	genrepo.CustomerPointHistoryRepo
}

func NewPointHistoryRepository(dataSource *database.DataSource) repo.PointHistoryRepo {
	db := dataSource.DB()
	return &pointHistoryRepository{
		q:                        dao.Use(db),
		CustomerPointHistoryRepo: genrepo.NewCustomerPointHistoryRepo(db),
	}
}

func (r *pointHistoryRepository) FindAll(ctx context.Context, userUid string, page pagination.PageOption) (int64, []*repo.PointHistoryResult, error) {
	h := r.q.CustomerPointHistoryEntity
	q := h.WithContext(ctx).Where(h.UserUid.Eq(userUid)).Order(h.CreatedAt.Desc())

	offset, limit := utils.PageOffsetLimit(page)
	entities, total, err := q.FindByPage(offset, limit)
	if err != nil {
		return 0, nil, err
	}

	results := make([]*repo.PointHistoryResult, len(entities))
	for i, e := range entities {
		results[i] = toPointHistoryResult(e)
	}
	return total, results, nil
}

func toPointResult(e *entity.CustomerPointEntity) *repo.PointResult {
	return &repo.PointResult{
		Id:      e.Id,
		Uid:     e.Uid,
		UserId:  e.UserId,
		UserUid: e.UserUid,
		Balance: e.Balance,
	}
}

func toPointHistoryResult(e *entity.CustomerPointHistoryEntity) *repo.PointHistoryResult {
	return &repo.PointHistoryResult{
		Id:         e.Id,
		Uid:        e.Uid,
		UserId:     e.UserId,
		UserUid:    e.UserUid,
		ChangeType: e.ChangeType,
		Amount:     e.Amount,
		Reason:     e.Reason,
		CreatedAt:  e.CreatedAt,
	}
}
