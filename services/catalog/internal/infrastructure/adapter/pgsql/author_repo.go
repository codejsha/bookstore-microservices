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

var _ repo.AuthorRepo = (*authorRepository)(nil)

type authorRepository struct {
	q *dao.Query
}

func NewAuthorRepository(dataSource *database.DataSource) repo.AuthorRepo {
	return &authorRepository{q: dao.Use(dataSource.DB())}
}

func (r *authorRepository) FindAll(ctx context.Context, opt option.AuthorQueryOption) (int64, []*repo.AuthorResult, error) {
	a := r.q.AuthorEntity
	q := a.WithContext(ctx)

	if v := opt.Name(); v != nil && *v != "" {
		q = q.Where(gen.Cond(gorm.Expr("author.name ILIKE ?", "%"+*v+"%"))...)
	}
	if v := opt.OlKey(); v != nil && *v != "" {
		q = q.Where(a.OlKey.Eq(*v))
	}
	if uid := opt.WorkUid(); uid != nil && *uid != "" {
		mp := r.q.WorkAuthorMappingEntity
		w := r.q.WorkEntity
		sub := mp.WithContext(ctx).
			Join(w, w.Id.EqCol(mp.WorkId), w.DeletedAt.IsNull()).
			Where(w.Uid.Eq(*uid)).
			Select(mp.AuthorId)
		q = q.Where(a.WithContext(ctx).Columns(a.Id).In(sub))
	}

	q = q.Order(buildOrderExprs(opt.Page().GetSort(), map[string]field.OrderExpr{
		"name":       a.Name,
		"created_at": a.CreatedAt,
		"updated_at": a.UpdatedAt,
	}, a.Id)...)

	offset, limit := utils.PageOffsetLimit(opt.Page())
	entities, total, err := q.FindByPage(offset, limit)
	if err != nil {
		return 0, nil, err
	}

	results := make([]*repo.AuthorResult, len(entities))
	for i, e := range entities {
		results[i] = toAuthorResult(e)
	}
	return total, results, nil
}

func (r *authorRepository) FindOne(ctx context.Context, id int64) (*repo.AuthorResult, error) {
	a := r.q.AuthorEntity
	e, err := a.WithContext(ctx).Where(a.Id.Eq(id)).First()
	if err != nil {
		return nil, err
	}
	return toAuthorResult(e), nil
}

func (r *authorRepository) FindByUid(ctx context.Context, uid string) (*repo.AuthorResult, error) {
	a := r.q.AuthorEntity
	e, err := a.WithContext(ctx).Where(a.Uid.Eq(uid)).First()
	if err != nil {
		return nil, err
	}
	return toAuthorResult(e), nil
}

func (r *authorRepository) Create(ctx context.Context, cmd command.AuthorCreateCommand) (int64, error) {
	e := &entity.AuthorEntity{
		Uid:           uuid.Must(uuid.NewV7()).String(),
		Name:          cmd.Name,
		Bio:           cmd.Bio,
		BirthDate:     cmd.BirthDate,
		DeathDate:     cmd.DeathDate,
		PhotoUid:      utils.ToJsonString(cmd.PhotoUids),
		AlternateName: utils.ToJsonString(cmd.AlternateNames),
		OlKey:         cmd.OlKey,
		CreatedAt:     time.Now(),
		Version:       1,
	}
	if err := r.q.AuthorEntity.WithContext(ctx).Create(e); err != nil {
		return 0, err
	}
	return e.Id, nil
}

func (r *authorRepository) Update(ctx context.Context, id int64, cmd command.AuthorUpdateCommand) error {
	a := r.q.AuthorEntity
	assigns := []field.AssignExpr{a.UpdatedAt.Value(time.Now())}
	if cmd.Name != nil {
		assigns = append(assigns, a.Name.Value(*cmd.Name))
	}
	if cmd.Bio != nil {
		assigns = append(assigns, a.Bio.Value(*cmd.Bio))
	}
	if cmd.BirthDate != nil {
		assigns = append(assigns, a.BirthDate.Value(*cmd.BirthDate))
	}
	if cmd.DeathDate != nil {
		assigns = append(assigns, a.DeathDate.Value(*cmd.DeathDate))
	}
	if cmd.PhotoUids != nil {
		if s := utils.ToJsonString(cmd.PhotoUids); s != nil {
			assigns = append(assigns, a.PhotoUid.Value(*s))
		}
	}
	if cmd.AlternateNames != nil {
		if s := utils.ToJsonString(cmd.AlternateNames); s != nil {
			assigns = append(assigns, a.AlternateName.Value(*s))
		}
	}
	if cmd.OlKey != nil {
		assigns = append(assigns, a.OlKey.Value(*cmd.OlKey))
	}
	_, err := a.WithContext(ctx).Where(a.Id.Eq(id)).UpdateSimple(assigns...)
	return err
}

func toAuthorResult(e *entity.AuthorEntity) *repo.AuthorResult {
	return &repo.AuthorResult{
		Id:             e.Id,
		Uid:            e.Uid,
		Name:           e.Name,
		Bio:            e.Bio,
		BirthDate:      e.BirthDate,
		DeathDate:      e.DeathDate,
		PhotoUids:      utils.ParseJsonStringArray(e.PhotoUid),
		AlternateNames: utils.ParseJsonStringArray(e.AlternateName),
		OlKey:          e.OlKey,
	}
}
