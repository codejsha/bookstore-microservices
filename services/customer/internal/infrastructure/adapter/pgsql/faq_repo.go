package pgsql

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gen/field"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/codejsha/shared-library-go/pkg/database"

	"github.com/codejsha/bookstore-microservices/customer/generated/infrastructure/port/dao"
	"github.com/codejsha/bookstore-microservices/customer/generated/infrastructure/port/entity"
	genrepo "github.com/codejsha/bookstore-microservices/customer/generated/infrastructure/port/repo"
	"github.com/codejsha/bookstore-microservices/customer/internal/application/port/repo"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/model/option"
)

var _ repo.FaqRepo = (*faqRepository)(nil)

type faqRepository struct {
	q  *dao.Query
	db *gorm.DB
	genrepo.CustomerFaqRepo
}

func NewFaqRepository(dataSource *database.DataSource) repo.FaqRepo {
	db := dataSource.DB()
	return &faqRepository{
		q:               dao.Use(db),
		db:              db,
		CustomerFaqRepo: genrepo.NewCustomerFaqRepo(db),
	}
}

func (r *faqRepository) sortColumns() map[string]sortColumn {
	f := r.q.CustomerFaqEntity
	return map[string]sortColumn{
		"view_count": {asc: f.ViewCount, desc: f.ViewCount.Desc()},
		"created_at": {asc: f.CreatedAt, desc: f.CreatedAt.Desc()},
		"updated_at": {asc: f.UpdatedAt, desc: f.UpdatedAt.Desc()},
		"question":   {asc: f.Question, desc: f.Question.Desc()},
		"published":  {asc: f.Published, desc: f.Published.Desc()},
	}
}

func (r *faqRepository) FindAll(ctx context.Context, opt option.FaqQueryOption) (int64, []*repo.FaqResult, error) {
	f := r.q.CustomerFaqEntity
	q := f.WithContext(ctx)

	if opt.PublishedOnly() {
		q = q.Where(f.Published.Is(true))
	}
	if v := opt.CategoryUid(); v != nil && *v != "" {
		categoryId, err := categoryIdByUid(ctx, r.q, *v, true)
		if err != nil {
			return 0, nil, err
		}
		if categoryId == nil {
			return 0, []*repo.FaqResult{}, nil
		}
		q = q.Where(f.CategoryId.Eq(*categoryId))
	}
	if v := opt.Query(); v != nil && strings.TrimSpace(*v) != "" {
		pattern := "%" + strings.ToLower(strings.TrimSpace(*v)) + "%"
		q = q.Where(field.Or(f.Question.Lower().Like(pattern), f.Answer.Lower().Like(pattern)))
	}
	q = q.Order(orderExprs(opt.Page().Sort(), r.sortColumns(), f.ViewCount.Desc())...)

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

	results := make([]*repo.FaqResult, len(entities))
	for i, e := range entities {
		results[i] = toFaqResult(e, lookupCategoryUid(categoryUids, e.CategoryId))
	}
	return total, results, nil
}

func (r *faqRepository) FindByUid(ctx context.Context, uid string) (*repo.FaqResult, error) {
	f := r.q.CustomerFaqEntity
	e, err := f.WithContext(ctx).Where(f.Uid.Eq(uid)).First()
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
	return toFaqResult(e, categoryUid), nil
}

func (r *faqRepository) Insert(ctx context.Context, faq repo.FaqCreate) (*repo.FaqResult, error) {
	now := time.Now()
	ent := &entity.CustomerFaqEntity{
		Uid:        uuid.Must(uuid.NewV7()).String(),
		CategoryId: faq.CategoryId,
		Question:   faq.Question,
		Answer:     faq.Answer,
		ViewCount:  0,
		Published:  faq.Published,
		CreatedAt:  now,
		UpdatedAt:  &now,
		Version:    1,
	}
	if err := r.q.CustomerFaqEntity.WithContext(ctx).Create(ent); err != nil {
		return nil, err
	}
	categoryUid, err := categoryUidById(ctx, r.q, ent.CategoryId)
	if err != nil {
		return nil, err
	}
	return toFaqResult(ent, categoryUid), nil
}

func (r *faqRepository) UpdateByUid(ctx context.Context, uid string, apply func(*repo.FaqResult)) (*repo.FaqResult, error) {
	var updated *entity.CustomerFaqEntity
	err := dao.Use(r.db.WithContext(ctx)).Transaction(func(tx *dao.Query) error {
		f := tx.CustomerFaqEntity
		e, err := f.WithContext(ctx).
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where(f.Uid.Eq(uid)).
			First()
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil
			}
			return err
		}

		result := toFaqResult(e, nil)
		apply(result)

		now := time.Now()
		e.CategoryId = result.CategoryId
		e.Question = result.Question
		e.Answer = result.Answer
		e.Published = result.Published
		e.UpdatedAt = &now
		e.Version++
		if err := f.WithContext(ctx).Save(e); err != nil {
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
	return toFaqResult(updated, categoryUid), nil
}

func (r *faqRepository) SoftDeleteByUid(ctx context.Context, uid string) (bool, error) {
	f := r.q.CustomerFaqEntity
	info, err := f.WithContext(ctx).Where(f.Uid.Eq(uid)).Delete()
	if err != nil {
		return false, err
	}
	return info.RowsAffected > 0, nil
}

func (r *faqRepository) IncrementViewCount(ctx context.Context, uid string) error {
	f := r.q.CustomerFaqEntity
	_, err := f.WithContext(ctx).Where(f.Uid.Eq(uid)).UpdateSimple(f.ViewCount.Add(1))
	return err
}

func toFaqResult(e *entity.CustomerFaqEntity, categoryUid *string) *repo.FaqResult {
	return &repo.FaqResult{
		Id:          e.Id,
		Uid:         e.Uid,
		CategoryId:  e.CategoryId,
		CategoryUid: categoryUid,
		Question:    e.Question,
		Answer:      e.Answer,
		ViewCount:   e.ViewCount,
		Published:   e.Published,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}
