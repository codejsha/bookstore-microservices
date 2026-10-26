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
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/model/option"
)

var _ repo.TicketRepo = (*ticketRepository)(nil)

type ticketRepository struct {
	q  *dao.Query
	db *gorm.DB
	genrepo.CustomerTicketRepo
}

func NewTicketRepository(dataSource *database.DataSource) repo.TicketRepo {
	db := dataSource.DB()
	return &ticketRepository{
		q:                  dao.Use(db),
		db:                 db,
		CustomerTicketRepo: genrepo.NewCustomerTicketRepo(db),
	}
}

func (r *ticketRepository) sortColumns() map[string]sortColumn {
	t := r.q.CustomerTicketEntity
	return map[string]sortColumn{
		"created_at":  {asc: t.CreatedAt, desc: t.CreatedAt.Desc()},
		"updated_at":  {asc: t.UpdatedAt, desc: t.UpdatedAt.Desc()},
		"status":      {asc: t.Status, desc: t.Status.Desc()},
		"priority":    {asc: t.Priority, desc: t.Priority.Desc()},
		"subject":     {asc: t.Subject, desc: t.Subject.Desc()},
		"resolved_at": {asc: t.ResolvedAt, desc: t.ResolvedAt.Desc()},
	}
}

func (r *ticketRepository) FindAll(ctx context.Context, opt option.TicketQueryOption) (int64, []*repo.TicketResult, error) {
	t := r.q.CustomerTicketEntity
	q := t.WithContext(ctx)

	if v := opt.CustomerUid(); v != nil && *v != "" {
		q = q.Where(t.CustomerUid.Eq(*v))
	}
	if v := opt.AssigneeUid(); v != nil && *v != "" {
		q = q.Where(t.AssigneeUid.Eq(*v))
	}
	if v := opt.Status(); v != nil {
		q = q.Where(t.Status.Eq(string(*v)))
	}
	if v := opt.Priority(); v != nil {
		q = q.Where(t.Priority.Eq(string(*v)))
	}
	if v := opt.CategoryUid(); v != nil && *v != "" {
		categoryId, err := categoryIdByUid(ctx, r.q, *v, true)
		if err != nil {
			return 0, nil, err
		}
		if categoryId == nil {
			return 0, []*repo.TicketResult{}, nil
		}
		q = q.Where(t.CategoryId.Eq(*categoryId))
	}
	q = q.Order(orderExprs(opt.Page().Sort(), r.sortColumns(), t.CreatedAt.Desc())...)

	offset, limit := opt.Page().OffsetLimit()
	entities, total, err := q.FindByPage(offset, limit)
	if err != nil {
		return 0, nil, err
	}

	ids := make([]*int64, len(entities))
	for i, e := range entities {
		ids[i] = e.CategoryId
	}
	categoryUids, err := categoryUidsByIds(ctx, r.q, collectCategoryIds(ids))
	if err != nil {
		return 0, nil, err
	}

	results := make([]*repo.TicketResult, len(entities))
	for i, e := range entities {
		results[i] = toTicketResult(e, lookupCategoryUid(categoryUids, e.CategoryId))
	}
	return total, results, nil
}

func (r *ticketRepository) FindByUid(ctx context.Context, uid string) (*repo.TicketResult, error) {
	t := r.q.CustomerTicketEntity
	e, err := t.WithContext(ctx).Where(t.Uid.Eq(uid)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	categoryUid, err := categoryUidById(ctx, r.q, e.CategoryId)
	if err != nil {
		return nil, err
	}
	return toTicketResult(e, categoryUid), nil
}

func (r *ticketRepository) Insert(ctx context.Context, ticket repo.TicketCreate) (*repo.TicketResult, error) {
	now := time.Now()
	ent := &entity.CustomerTicketEntity{
		Uid:         uuid.Must(uuid.NewV7()).String(),
		CustomerUid: ticket.CustomerUid,
		CategoryId:  ticket.CategoryId,
		Subject:     ticket.Subject,
		Description: ticket.Description,
		Status:      ticket.Status,
		Priority:    ticket.Priority,
		CreatedAt:   now,
		UpdatedAt:   &now,
		Version:     1,
	}
	if err := r.q.CustomerTicketEntity.WithContext(ctx).Create(ent); err != nil {
		return nil, err
	}
	categoryUid, err := categoryUidById(ctx, r.q, ent.CategoryId)
	if err != nil {
		return nil, err
	}
	return toTicketResult(ent, categoryUid), nil
}

func (r *ticketRepository) UpdateByUid(ctx context.Context, uid string, apply func(*repo.TicketResult)) (*repo.TicketResult, error) {
	var updated *entity.CustomerTicketEntity
	err := dao.Use(r.db.WithContext(ctx)).Transaction(func(tx *dao.Query) error {
		t := tx.CustomerTicketEntity
		e, err := t.WithContext(ctx).
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where(t.Uid.Eq(uid)).
			First()
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil
			}
			return err
		}

		result := toTicketResult(e, nil)
		apply(result)

		now := time.Now()
		e.CategoryId = result.CategoryId
		e.AssigneeUid = result.AssigneeUid
		e.Subject = result.Subject
		e.Description = result.Description
		e.Status = result.Status
		e.Priority = result.Priority
		e.ResolvedAt = result.ResolvedAt
		e.UpdatedAt = &now
		e.Version++
		if err := t.WithContext(ctx).Save(e); err != nil {
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
	categoryUid, err := categoryUidById(ctx, r.q, updated.CategoryId)
	if err != nil {
		return nil, err
	}
	return toTicketResult(updated, categoryUid), nil
}

func toTicketResult(e *entity.CustomerTicketEntity, categoryUid *string) *repo.TicketResult {
	return &repo.TicketResult{
		Id:          e.Id,
		Uid:         e.Uid,
		CustomerUid: e.CustomerUid,
		CategoryId:  e.CategoryId,
		CategoryUid: categoryUid,
		AssigneeUid: e.AssigneeUid,
		Subject:     e.Subject,
		Description: e.Description,
		Status:      e.Status,
		Priority:    e.Priority,
		ResolvedAt:  e.ResolvedAt,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}
