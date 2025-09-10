package pgsql

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"

	"github.com/codejsha/shared-library-go/pkg/database"

	"github.com/codejsha/bookstore-microservices/customer/generated/infrastructure/port/dao"
	"github.com/codejsha/bookstore-microservices/customer/generated/infrastructure/port/entity"
	genrepo "github.com/codejsha/bookstore-microservices/customer/generated/infrastructure/port/repo"
	"github.com/codejsha/bookstore-microservices/customer/internal/application/port/repo"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/model/option"
	"github.com/codejsha/bookstore-microservices/customer/internal/infrastructure/support/utils"
)

var _ repo.ReviewRepo = (*reviewRepository)(nil)

type reviewRepository struct {
	q *dao.Query
	genrepo.CustomerReviewRepo
}

func NewReviewRepository(dataSource *database.DataSource) repo.ReviewRepo {
	db := dataSource.DB()
	return &reviewRepository{
		q:                  dao.Use(db),
		CustomerReviewRepo: genrepo.NewCustomerReviewRepo(db),
	}
}

func (r *reviewRepository) FindAll(ctx context.Context, opt option.ReviewQueryOption) (int64, []*repo.ReviewResult, error) {
	rv := r.q.CustomerReviewEntity
	q := rv.WithContext(ctx)

	if v := opt.UserUid(); v != nil && *v != "" {
		q = q.Where(rv.UserUid.Eq(*v))
	}
	if v := opt.BookUid(); v != nil && *v != "" {
		q = q.Where(rv.BookUid.Eq(*v))
	}
	if v := opt.Rating(); v != nil {
		q = q.Where(rv.Rating.Eq(int16(*v)))
	}
	q = q.Order(rv.CreatedAt.Desc())

	offset, limit := utils.PageOffsetLimit(opt.Page())
	entities, total, err := q.FindByPage(offset, limit)
	if err != nil {
		return 0, nil, err
	}

	results := make([]*repo.ReviewResult, len(entities))
	for i, e := range entities {
		results[i] = toReviewResult(e)
	}
	return total, results, nil
}

func (r *reviewRepository) Insert(ctx context.Context, review repo.ReviewCreate) (*repo.ReviewResult, error) {
	rv := r.q.CustomerReviewEntity

	existing, err := rv.WithContext(ctx).
		Where(rv.UserUid.Eq(review.UserUid), rv.BookUid.Eq(review.BookUid)).
		First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existing != nil {
		return nil, repo.ErrReviewExists
	}

	ent := &entity.CustomerReviewEntity{
		Uid:       uuid.Must(uuid.NewV7()).String(),
		UserUid:   review.UserUid,
		BookUid:   review.BookUid,
		Rating:    int16(review.Rating),
		Title:     review.Title,
		Content:   review.Content,
		CreatedAt: time.Now(),
		Version:   1,
	}
	if err := rv.WithContext(ctx).Create(ent); err != nil {
		if isUniqueViolation(err) {
			return nil, repo.ErrReviewExists
		}
		return nil, err
	}
	return toReviewResult(ent), nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}

func (r *reviewRepository) FindOne(ctx context.Context, id int64) (*repo.ReviewResult, error) {
	rv := r.q.CustomerReviewEntity
	e, err := rv.WithContext(ctx).Where(rv.Id.Eq(id)).First()
	if err != nil {
		return nil, err
	}
	return toReviewResult(e), nil
}

func (r *reviewRepository) FindByUid(ctx context.Context, uid string) (*repo.ReviewResult, error) {
	rv := r.q.CustomerReviewEntity
	e, err := rv.WithContext(ctx).Where(rv.Uid.Eq(uid)).First()
	if err != nil {
		return nil, err
	}
	return toReviewResult(e), nil
}

func toReviewResult(e *entity.CustomerReviewEntity) *repo.ReviewResult {
	return &repo.ReviewResult{
		Id:        e.Id,
		Uid:       e.Uid,
		UserId:    e.UserId,
		UserUid:   e.UserUid,
		BookId:    e.BookId,
		BookUid:   e.BookUid,
		Rating:    int32(e.Rating),
		Title:     e.Title,
		Content:   e.Content,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}
