package pgsql

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"

	"github.com/codejsha/shared-library-go/pkg/database"

	"github.com/codejsha/bookstore-microservices/catalog/generated/infrastructure/port/dao"
	"github.com/codejsha/bookstore-microservices/catalog/generated/infrastructure/port/entity"
	"github.com/codejsha/bookstore-microservices/catalog/internal/application/port/repo"
	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/model/command"
	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/model/option"
	"github.com/codejsha/bookstore-microservices/catalog/internal/infrastructure/support/utils"
)

var _ repo.PublisherRepo = (*publisherRepository)(nil)

type publisherRepository struct {
	q *dao.Query
}

func NewPublisherRepository(dataSource *database.DataSource) repo.PublisherRepo {
	return &publisherRepository{q: dao.Use(dataSource.DB())}
}

func (r *publisherRepository) FindAll(ctx context.Context, opt option.PublisherQueryOption) (int64, []*repo.PublisherResult, error) {
	p := r.q.PublisherEntity
	q := p.WithContext(ctx)

	if v := opt.Name(); v != nil && *v != "" {
		q = q.Where(gen.Cond(gorm.Expr("publisher.name ILIKE ?", "%"+*v+"%"))...)
	}
	if v := opt.OlKey(); v != nil && *v != "" {
		q = q.Where(p.OlKey.Eq(*v))
	}

	q = q.Order(buildOrderExprs(opt.Page().GetSort(), map[string]field.OrderExpr{
		"name":       p.Name,
		"created_at": p.CreatedAt,
		"updated_at": p.UpdatedAt,
	}, p.Id)...)

	offset, limit := utils.PageOffsetLimit(opt.Page())
	entities, total, err := q.FindByPage(offset, limit)
	if err != nil {
		return 0, nil, err
	}

	results := make([]*repo.PublisherResult, len(entities))
	for i, e := range entities {
		results[i] = toPublisherResult(e)
	}
	return total, results, nil
}

func (r *publisherRepository) FindOne(ctx context.Context, id int64) (*repo.PublisherResult, error) {
	p := r.q.PublisherEntity
	e, err := p.WithContext(ctx).Where(p.Id.Eq(id)).First()
	if err != nil {
		return nil, err
	}
	return toPublisherResult(e), nil
}

func (r *publisherRepository) FindByUid(ctx context.Context, uid string) (*repo.PublisherResult, error) {
	p := r.q.PublisherEntity
	e, err := p.WithContext(ctx).Where(p.Uid.Eq(uid)).First()
	if err != nil {
		return nil, err
	}
	return toPublisherResult(e), nil
}

func (r *publisherRepository) Create(ctx context.Context, cmd command.PublisherCreateCommand) (int64, error) {
	e := &entity.PublisherEntity{
		Uid:       uuid.Must(uuid.NewV7()).String(),
		Name:      cmd.Name,
		Address:   cmd.Address,
		OlKey:     cmd.OlKey,
		CreatedAt: time.Now(),
		Version:   1,
	}
	if err := r.q.PublisherEntity.WithContext(ctx).Create(e); err != nil {
		return 0, err
	}
	return e.Id, nil
}

func (r *publisherRepository) Update(ctx context.Context, id int64, cmd command.PublisherUpdateCommand) error {
	p := r.q.PublisherEntity
	assigns := []field.AssignExpr{p.UpdatedAt.Value(time.Now())}
	if cmd.Name != nil {
		assigns = append(assigns, p.Name.Value(*cmd.Name))
	}
	if cmd.Address != nil {
		assigns = append(assigns, p.Address.Value(*cmd.Address))
	}
	if cmd.OlKey != nil {
		assigns = append(assigns, p.OlKey.Value(*cmd.OlKey))
	}
	_, err := p.WithContext(ctx).Where(p.Id.Eq(id)).UpdateSimple(assigns...)
	return err
}

func toPublisherResult(e *entity.PublisherEntity) *repo.PublisherResult {
	return &repo.PublisherResult{
		Id:      e.Id,
		Uid:     e.Uid,
		Name:    e.Name,
		Address: e.Address,
		OlKey:   e.OlKey,
	}
}
