package pgsql

import (
	"context"

	"github.com/codejsha/shared-library-go/pkg/database"
	"github.com/codejsha/shared-library-go/pkg/pagination"

	"github.com/codejsha/bookstore-microservices/inventory/generated/infrastructure/port/dao"
	genrepo "github.com/codejsha/bookstore-microservices/inventory/generated/infrastructure/port/repo"
	"github.com/codejsha/bookstore-microservices/inventory/internal/application/port/repo"
	"github.com/codejsha/bookstore-microservices/inventory/internal/infrastructure/support/utils"
)

var _ repo.StockHistoryRepo = (*stockHistoryRepository)(nil)

type stockHistoryRepository struct {
	q *dao.Query
	genrepo.StockHistoryRepo
}

func NewStockHistoryRepository(dataSource *database.DataSource) repo.StockHistoryRepo {
	db := dataSource.DB()
	return &stockHistoryRepository{
		q:                dao.Use(db),
		StockHistoryRepo: genrepo.NewStockHistoryRepo(db),
	}
}

func (r *stockHistoryRepository) FindAll(ctx context.Context, stockUid string, page pagination.PageOption) (int64, []*repo.StockHistoryResult, error) {
	h := r.q.StockHistoryEntity
	q := h.WithContext(ctx).Where(h.StockUid.Eq(stockUid)).Order(h.CreatedAt.Desc())

	offset, limit := utils.PageOffsetLimit(page)
	entities, total, err := q.FindByPage(offset, limit)
	if err != nil {
		return 0, nil, err
	}
	results := make([]*repo.StockHistoryResult, len(entities))
	for i, e := range entities {
		results[i] = &repo.StockHistoryResult{
			Uid:        e.Uid,
			ChangeType: e.ChangeType,
			Reason:     e.Reason,
			ChangeQty:  e.ChangeQty,
			CreatedAt:  e.CreatedAt,
		}
	}
	return total, results, nil
}
