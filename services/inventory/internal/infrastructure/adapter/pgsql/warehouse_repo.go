package pgsql

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"

	"github.com/codejsha/shared-library-go/pkg/database"
	"github.com/codejsha/shared-library-go/pkg/database/gormutils"

	"github.com/codejsha/bookstore-microservices/inventory/generated/infrastructure/port/dao"
	"github.com/codejsha/bookstore-microservices/inventory/generated/infrastructure/port/entity"
	genrepo "github.com/codejsha/bookstore-microservices/inventory/generated/infrastructure/port/repo"
	"github.com/codejsha/bookstore-microservices/inventory/internal/application/port/repo"
	"github.com/codejsha/bookstore-microservices/inventory/internal/domain/model/option"
)

var _ repo.WarehouseRepo = (*warehouseRepository)(nil)

type warehouseRepository struct {
	db *gorm.DB
	genrepo.WarehouseRepo
}

func NewWarehouseRepository(dataSource *database.DataSource) repo.WarehouseRepo {
	db := dataSource.DB()
	return &warehouseRepository{
		db:            db,
		WarehouseRepo: genrepo.NewWarehouseRepo(db),
	}
}

func (r warehouseRepository) FindAll(ctx context.Context, opt option.WarehouseQueryOption) (int64, []*repo.WarehouseResult, error) {
	q := dao.Use(r.db)
	sortable := map[string]field.OrderExpr{
		"uid":       q.WarehouseEntity.Uid,
		"name":      q.WarehouseEntity.Name,
		"address":   q.WarehouseEntity.Address,
		"capacity":  q.WarehouseEntity.Capacity,
		"createdat": q.WarehouseEntity.CreatedAt,
		"updatedat": q.WarehouseEntity.UpdatedAt,
	}

	base := q.WarehouseEntity.WithContext(ctx).
		Scopes(r.buildWhereScope(q, opt))

	var entities []*entity.WarehouseEntity
	err := base.
		Select(
			q.WarehouseEntity.Uid,
			q.WarehouseEntity.Name,
			q.WarehouseEntity.Address,
			q.WarehouseEntity.Capacity,
		).
		Scopes(gormutils.BuildPageScope(opt.Page(), sortable)).
		Scan(&entities)
	if err != nil {
		return 0, nil, err
	}

	total, err := base.Count()
	if err != nil {
		return 0, nil, err
	}

	results := make([]*repo.WarehouseResult, len(entities))
	for i, e := range entities {
		results[i] = toWarehouseResult(e)
	}
	return total, results, nil
}

func (r warehouseRepository) FindByUid(ctx context.Context, uid string) (*repo.WarehouseResult, error) {
	entities, err := r.WarehouseRepo.FetchByUid(ctx, uid)
	if err != nil {
		return nil, err
	}
	if len(entities) == 0 {
		return nil, nil
	}
	return toWarehouseResult(entities[0]), nil
}

func (r warehouseRepository) Create(ctx context.Context, p repo.WarehouseCreateParams) (*repo.WarehouseResult, error) {
	e := &entity.WarehouseEntity{
		Uid:       uuid.Must(uuid.NewV7()).String(),
		Name:      p.Name,
		Address:   p.Address,
		Capacity:  p.Capacity,
		CreatedAt: time.Now(),
		Version:   1,
	}
	if err := r.WarehouseRepo.Create(ctx, e); err != nil {
		return nil, err
	}
	return toWarehouseResult(e), nil
}

func (r warehouseRepository) Update(ctx context.Context, p repo.WarehouseUpdateParams) (*repo.WarehouseResult, error) {
	entities, err := r.WarehouseRepo.FetchByUid(ctx, p.Uid)
	if err != nil {
		return nil, err
	}
	if len(entities) == 0 {
		return nil, nil
	}
	e := entities[0]
	if p.Name != nil {
		e.Name = *p.Name
	}
	if p.Address != nil {
		e.Address = p.Address
	}
	if p.Capacity != nil {
		e.Capacity = *p.Capacity
	}
	now := time.Now()
	e.UpdatedAt = &now
	if err := r.WarehouseRepo.Update(ctx, e); err != nil {
		return nil, err
	}
	return toWarehouseResult(e), nil
}

func toWarehouseResult(e *entity.WarehouseEntity) *repo.WarehouseResult {
	return &repo.WarehouseResult{
		Uid:       e.Uid,
		Name:      e.Name,
		Address:   e.Address,
		Capacity:  e.Capacity,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

func (r warehouseRepository) buildWhereScope(q *dao.Query, opt option.WarehouseQueryOption) func(gen.Dao) gen.Dao {
	conds := make([]gen.Condition, 0)
	if name := opt.Name(); name != nil && *name != "" {
		conds = append(conds, gen.Cond(gorm.Expr("warehouse.name ILIKE ?", "%"+*name+"%"))...)
	}

	return func(dao gen.Dao) gen.Dao {
		if len(conds) > 0 {
			return dao.Where(conds...)
		}
		return dao
	}
}
