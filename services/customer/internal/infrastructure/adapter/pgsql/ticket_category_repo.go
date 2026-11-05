package pgsql

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/codejsha/shared-library-go/pkg/database"

	"github.com/codejsha/bookstore-microservices/customer/generated/infrastructure/port/dao"
	"github.com/codejsha/bookstore-microservices/customer/generated/infrastructure/port/entity"
	genrepo "github.com/codejsha/bookstore-microservices/customer/generated/infrastructure/port/repo"
	"github.com/codejsha/bookstore-microservices/customer/internal/application/port/repo"
)

var _ repo.TicketCategoryRepo = (*ticketCategoryRepository)(nil)

type ticketCategoryRepository struct {
	q  *dao.Query
	db *gorm.DB
	genrepo.CustomerTicketCategoryRepo
}

func NewTicketCategoryRepository(dataSource *database.DataSource) repo.TicketCategoryRepo {
	db := dataSource.DB()
	return &ticketCategoryRepository{
		q:                          dao.Use(db),
		db:                         db,
		CustomerTicketCategoryRepo: genrepo.NewCustomerTicketCategoryRepo(db),
	}
}

func (r *ticketCategoryRepository) FindAll(ctx context.Context) ([]*repo.TicketCategoryResult, error) {
	c := r.q.CustomerTicketCategoryEntity
	entities, err := c.WithContext(ctx).Order(c.Id).Find()
	if err != nil {
		return nil, err
	}
	ids := make([]*int64, len(entities))
	for i, e := range entities {
		ids[i] = e.ParentId
	}
	parentUids, err := categoryUidsByIds(ctx, r.q, collectCategoryIds(ids))
	if err != nil {
		return nil, err
	}
	results := make([]*repo.TicketCategoryResult, len(entities))
	for i, e := range entities {
		results[i] = toTicketCategoryResult(e, lookupCategoryUid(parentUids, e.ParentId))
	}
	return results, nil
}

func (r *ticketCategoryRepository) findOne(ctx context.Context, e *entity.CustomerTicketCategoryEntity, err error) (*repo.TicketCategoryResult, error) {
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	parentUid, err := categoryUidById(ctx, r.q, e.ParentId)
	if err != nil {
		return nil, err
	}
	return toTicketCategoryResult(e, parentUid), nil
}

func (r *ticketCategoryRepository) FindByUid(ctx context.Context, uid string) (*repo.TicketCategoryResult, error) {
	c := r.q.CustomerTicketCategoryEntity
	e, err := c.WithContext(ctx).Where(c.Uid.Eq(uid)).First()
	return r.findOne(ctx, e, err)
}

func (r *ticketCategoryRepository) FindByName(ctx context.Context, name string) (*repo.TicketCategoryResult, error) {
	c := r.q.CustomerTicketCategoryEntity
	e, err := c.WithContext(ctx).Where(c.Name.Eq(name)).First()
	return r.findOne(ctx, e, err)
}

func (r *ticketCategoryRepository) FindIdByUid(ctx context.Context, uid string) (*int64, error) {
	return categoryIdByUid(ctx, r.q, uid, false)
}

func (r *ticketCategoryRepository) Insert(ctx context.Context, category repo.TicketCategoryCreate) (*repo.TicketCategoryResult, error) {
	now := time.Now().UTC()
	ent := &entity.CustomerTicketCategoryEntity{
		Uid:         uuid.Must(uuid.NewV7()).String(),
		Name:        category.Name,
		Description: category.Description,
		ParentId:    category.ParentId,
		CreatedAt:   now,
		UpdatedAt:   &now,
		Version:     1,
	}
	if err := r.q.CustomerTicketCategoryEntity.WithContext(ctx).Create(ent); err != nil {
		if isUniqueViolation(err) {
			return nil, repo.ErrCategoryNameExists
		}
		return nil, err
	}
	parentUid, err := categoryUidById(ctx, r.q, ent.ParentId)
	if err != nil {
		return nil, err
	}
	return toTicketCategoryResult(ent, parentUid), nil
}

func (r *ticketCategoryRepository) UpdateByUid(ctx context.Context, uid string, apply func(*repo.TicketCategoryResult)) (*repo.TicketCategoryResult, error) {
	var updated *entity.CustomerTicketCategoryEntity
	err := dao.Use(r.db.WithContext(ctx)).Transaction(func(tx *dao.Query) error {
		c := tx.CustomerTicketCategoryEntity
		e, err := c.WithContext(ctx).
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where(c.Uid.Eq(uid)).
			First()
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil
			}
			return err
		}

		result := toTicketCategoryResult(e, nil)
		apply(result)

		now := time.Now().UTC()
		e.Name = result.Name
		e.Description = result.Description
		e.ParentId = result.ParentId
		e.UpdatedAt = &now
		e.Version++
		if err := c.WithContext(ctx).Save(e); err != nil {
			if isUniqueViolation(err) {
				return repo.ErrCategoryNameExists
			}
			return err
		}
		updated = e
		return nil
	})
	if err != nil {
		return nil, err
	}
	if updated == nil {
		return nil, nil
	}
	parentUid, err := categoryUidById(ctx, r.q, updated.ParentId)
	if err != nil {
		return nil, err
	}
	return toTicketCategoryResult(updated, parentUid), nil
}

func (r *ticketCategoryRepository) SoftDeleteByUid(ctx context.Context, uid string) (bool, error) {
	c := r.q.CustomerTicketCategoryEntity
	info, err := c.WithContext(ctx).Where(c.Uid.Eq(uid)).Delete()
	if err != nil {
		return false, err
	}
	return info.RowsAffected > 0, nil
}

func toTicketCategoryResult(e *entity.CustomerTicketCategoryEntity, parentUid *string) *repo.TicketCategoryResult {
	return &repo.TicketCategoryResult{
		Id:          e.Id,
		Uid:         e.Uid,
		Name:        e.Name,
		Description: e.Description,
		ParentId:    e.ParentId,
		ParentUid:   parentUid,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}
