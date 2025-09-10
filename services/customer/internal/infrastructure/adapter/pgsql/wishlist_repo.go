package pgsql

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm/clause"

	"github.com/codejsha/shared-library-go/pkg/database"

	"github.com/codejsha/bookstore-microservices/customer/generated/infrastructure/port/dao"
	"github.com/codejsha/bookstore-microservices/customer/generated/infrastructure/port/entity"
	genrepo "github.com/codejsha/bookstore-microservices/customer/generated/infrastructure/port/repo"
	"github.com/codejsha/bookstore-microservices/customer/internal/application/port/repo"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/model/option"
	"github.com/codejsha/bookstore-microservices/customer/internal/infrastructure/support/utils"
)

var _ repo.WishlistRepo = (*wishlistRepository)(nil)

type wishlistRepository struct {
	q *dao.Query
	genrepo.CustomerWishlistRepo
}

func NewWishlistRepository(dataSource *database.DataSource) repo.WishlistRepo {
	db := dataSource.DB()
	return &wishlistRepository{
		q:                    dao.Use(db),
		CustomerWishlistRepo: genrepo.NewCustomerWishlistRepo(db),
	}
}

func (r *wishlistRepository) FindAll(ctx context.Context, opt option.WishlistQueryOption) (int64, []*repo.WishlistResult, error) {
	w := r.q.CustomerWishlistEntity
	q := w.WithContext(ctx)

	if v := opt.UserUid(); v != nil && *v != "" {
		q = q.Where(w.UserUid.Eq(*v))
	}
	if v := opt.BookUids(); v != nil && len(*v) > 0 {
		q = q.Where(w.BookUid.In(*v...))
	}
	q = q.Order(w.CreatedAt.Desc())

	offset, limit := utils.PageOffsetLimit(opt.Page())
	entities, total, err := q.FindByPage(offset, limit)
	if err != nil {
		return 0, nil, err
	}

	results := make([]*repo.WishlistResult, len(entities))
	for i, e := range entities {
		results[i] = toWishlistResult(e)
	}
	return total, results, nil
}

func (r *wishlistRepository) AddBooks(ctx context.Context, userUid string, bookUids []string) error {
	if len(bookUids) == 0 {
		return nil
	}
	now := time.Now()
	rows := make([]*entity.CustomerWishlistEntity, 0, len(bookUids))
	for _, bookUid := range bookUids {
		if bookUid == "" {
			continue
		}
		rows = append(rows, &entity.CustomerWishlistEntity{
			Uid:       uuid.Must(uuid.NewV7()).String(),
			UserUid:   userUid,
			BookUid:   bookUid,
			CreatedAt: now,
			Version:   1,
		})
	}
	if len(rows) == 0 {
		return nil
	}
	return r.q.CustomerWishlistEntity.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:     []clause.Column{{Name: "user_uid"}, {Name: "book_uid"}},
			TargetWhere: clause.Where{Exprs: []clause.Expression{clause.Expr{SQL: "deleted_at IS NULL"}}},
			DoNothing:   true,
		}).
		Create(rows...)
}

func (r *wishlistRepository) FindOne(ctx context.Context, id int64) (*repo.WishlistResult, error) {
	w := r.q.CustomerWishlistEntity
	e, err := w.WithContext(ctx).Where(w.Id.Eq(id)).First()
	if err != nil {
		return nil, err
	}
	return toWishlistResult(e), nil
}

func toWishlistResult(e *entity.CustomerWishlistEntity) *repo.WishlistResult {
	return &repo.WishlistResult{
		Id:        e.Id,
		Uid:       e.Uid,
		UserId:    e.UserId,
		UserUid:   e.UserUid,
		BookId:    e.BookId,
		BookUid:   e.BookUid,
		CreatedAt: e.CreatedAt,
	}
}
