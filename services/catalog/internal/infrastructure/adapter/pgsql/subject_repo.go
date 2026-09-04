package pgsql

import (
	"context"
	"errors"
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

var _ repo.SubjectRepo = (*subjectRepository)(nil)

type subjectRepository struct {
	q *dao.Query
}

func NewSubjectRepository(dataSource *database.DataSource) repo.SubjectRepo {
	return &subjectRepository{q: dao.Use(dataSource.DB())}
}

func (r *subjectRepository) FindAll(ctx context.Context, opt option.SubjectQueryOption) (int64, []*repo.SubjectResult, error) {
	s := r.q.SubjectEntity
	q := s.WithContext(ctx)

	if v := opt.Name(); v != nil && *v != "" {
		q = q.Where(gen.Cond(gorm.Expr("subject.name ILIKE ?", "%"+*v+"%"))...)
	}

	q = q.Order(buildOrderExprs(opt.Page().GetSort(), map[string]field.OrderExpr{
		"name":       s.Name,
		"created_at": s.CreatedAt,
		"updated_at": s.UpdatedAt,
	}, s.Id)...)

	offset, limit := utils.PageOffsetLimit(opt.Page())
	entities, total, err := q.FindByPage(offset, limit)
	if err != nil {
		return 0, nil, err
	}

	results := make([]*repo.SubjectResult, len(entities))
	for i, e := range entities {
		results[i] = toSubjectResult(e)
	}
	return total, results, nil
}

func (r *subjectRepository) FindOne(ctx context.Context, id int64) (*repo.SubjectResult, error) {
	s := r.q.SubjectEntity
	e, err := s.WithContext(ctx).Where(s.Id.Eq(id)).First()
	if err != nil {
		return nil, err
	}
	return toSubjectResult(e), nil
}

func (r *subjectRepository) FindByUid(ctx context.Context, uid string) (*repo.SubjectResult, error) {
	s := r.q.SubjectEntity
	e, err := s.WithContext(ctx).Where(s.Uid.Eq(uid)).First()
	if err != nil {
		return nil, err
	}
	return toSubjectResult(e), nil
}

func (r *subjectRepository) FindByName(ctx context.Context, name string) (*repo.SubjectResult, error) {
	s := r.q.SubjectEntity
	e, err := s.WithContext(ctx).Where(s.Name.Eq(name)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return toSubjectResult(e), nil
}

func (r *subjectRepository) FindOrCreateByName(ctx context.Context, name string) (int64, error) {
	s := r.q.SubjectEntity
	e, err := s.WithContext(ctx).Where(s.Name.Eq(name)).First()
	if err == nil {
		return e.Id, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, err
	}
	e = &entity.SubjectEntity{
		Uid:       uuid.Must(uuid.NewV7()).String(),
		Name:      name,
		CreatedAt: time.Now(),
		Version:   1,
	}
	if err := s.WithContext(ctx).Create(e); err != nil {
		return 0, err
	}
	return e.Id, nil
}

func (r *subjectRepository) Create(ctx context.Context, cmd command.SubjectCreateCommand) (int64, error) {
	e := &entity.SubjectEntity{
		Uid:       uuid.Must(uuid.NewV7()).String(),
		Name:      cmd.Name,
		CreatedAt: time.Now(),
		Version:   1,
	}
	if err := r.q.SubjectEntity.WithContext(ctx).Create(e); err != nil {
		return 0, err
	}
	return e.Id, nil
}

func (r *subjectRepository) Update(ctx context.Context, id int64, cmd command.SubjectUpdateCommand) error {
	s := r.q.SubjectEntity
	assigns := []field.AssignExpr{s.UpdatedAt.Value(time.Now())}
	if cmd.Name != nil {
		assigns = append(assigns, s.Name.Value(*cmd.Name))
	}
	_, err := s.WithContext(ctx).Where(s.Id.Eq(id)).UpdateSimple(assigns...)
	return err
}

func toSubjectResult(e *entity.SubjectEntity) *repo.SubjectResult {
	return &repo.SubjectResult{
		Id:   e.Id,
		Uid:  e.Uid,
		Name: e.Name,
	}
}
