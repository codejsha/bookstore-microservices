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

var _ repo.WorkRepo = (*workRepository)(nil)

type workRepository struct {
	q *dao.Query
}

func NewWorkRepository(dataSource *database.DataSource) repo.WorkRepo {
	return &workRepository{q: dao.Use(dataSource.DB())}
}

type workAuthorRow struct {
	WorkId     int64  `gorm:"column:work_id"`
	AuthorUid  string `gorm:"column:uid"`
	AuthorName string `gorm:"column:name"`
}

type workSubjectRow struct {
	WorkId      int64  `gorm:"column:work_id"`
	SubjectUid  string `gorm:"column:uid"`
	SubjectName string `gorm:"column:name"`
}

func (r *workRepository) FindAll(ctx context.Context, opt option.WorkQueryOption) (int64, []*repo.WorkResult, error) {
	w := r.q.WorkEntity
	q := w.WithContext(ctx)

	if uid := opt.AuthorUid(); uid != nil && *uid != "" {
		mp := r.q.WorkAuthorMappingEntity
		a := r.q.AuthorEntity
		sub := mp.WithContext(ctx).
			Join(a, a.Id.EqCol(mp.AuthorId), a.DeletedAt.IsNull()).
			Where(a.Uid.Eq(*uid)).
			Select(mp.WorkId)
		q = q.Where(w.Columns(w.Id).In(sub))
	}
	if uid := opt.SubjectUid(); uid != nil && *uid != "" {
		mp := r.q.WorkSubjectMappingEntity
		s := r.q.SubjectEntity
		sub := mp.WithContext(ctx).
			Join(s, s.Id.EqCol(mp.SubjectId), s.DeletedAt.IsNull()).
			Where(s.Uid.Eq(*uid)).
			Select(mp.WorkId)
		q = q.Where(w.Columns(w.Id).In(sub))
	}
	if v := opt.Title(); v != nil && *v != "" {
		q = q.Where(gen.Cond(gorm.Expr("work.title ILIKE ?", "%"+*v+"%"))...)
	}
	if v := opt.OlKey(); v != nil && *v != "" {
		q = q.Where(w.OlKey.Eq(*v))
	}

	q = q.Order(buildOrderExprs(opt.Page().GetSort(), map[string]field.OrderExpr{
		"title":      w.Title,
		"created_at": w.CreatedAt,
		"updated_at": w.UpdatedAt,
	}, w.Id)...)

	offset, limit := utils.PageOffsetLimit(opt.Page())
	entities, total, err := q.FindByPage(offset, limit)
	if err != nil {
		return 0, nil, err
	}
	if len(entities) == 0 {
		return total, []*repo.WorkResult{}, nil
	}

	workIds := make([]int64, len(entities))
	for i, e := range entities {
		workIds[i] = e.Id
	}

	authorMap, err := r.batchLoadAuthors(ctx, workIds)
	if err != nil {
		return 0, nil, err
	}
	subjectMap, err := r.batchLoadSubjects(ctx, workIds)
	if err != nil {
		return 0, nil, err
	}

	results := make([]*repo.WorkResult, len(entities))
	for i, e := range entities {
		results[i] = toWorkResult(e, authorMap[e.Id], subjectMap[e.Id])
	}
	return total, results, nil
}

func (r *workRepository) FindOne(ctx context.Context, id int64) (*repo.WorkResult, error) {
	w := r.q.WorkEntity
	e, err := w.WithContext(ctx).Where(w.Id.Eq(id)).First()
	if err != nil {
		return nil, err
	}
	return r.hydrateWork(ctx, e)
}

func (r *workRepository) FindByUid(ctx context.Context, uid string) (*repo.WorkResult, error) {
	w := r.q.WorkEntity
	e, err := w.WithContext(ctx).Where(w.Uid.Eq(uid)).First()
	if err != nil {
		return nil, err
	}
	return r.hydrateWork(ctx, e)
}

func (r *workRepository) hydrateWork(ctx context.Context, e *entity.WorkEntity) (*repo.WorkResult, error) {
	authorMap, err := r.batchLoadAuthors(ctx, []int64{e.Id})
	if err != nil {
		return nil, err
	}
	subjectMap, err := r.batchLoadSubjects(ctx, []int64{e.Id})
	if err != nil {
		return nil, err
	}
	return toWorkResult(e, authorMap[e.Id], subjectMap[e.Id]), nil
}

func (r *workRepository) Create(ctx context.Context, cmd command.WorkCreateCommand) (int64, error) {
	var id int64
	err := r.q.Transaction(func(tx *dao.Query) error {
		now := time.Now()
		work := &entity.WorkEntity{
			Uid:              uuid.Must(uuid.NewV7()).String(),
			Title:            cmd.Title,
			Description:      cmd.Description,
			CoverUid:         utils.ToJsonString(cmd.CoverUids),
			FirstPublishDate: cmd.FirstPublishDate,
			OlKey:            cmd.OlKey,
			CreatedAt:        now,
			Version:          1,
		}
		if err := tx.WorkEntity.WithContext(ctx).Create(work); err != nil {
			return err
		}
		id = work.Id

		ae := tx.AuthorEntity
		for _, authorUid := range cmd.AuthorUids {
			author, err := ae.WithContext(ctx).Where(ae.Uid.Eq(authorUid)).First()
			if err != nil {
				return err
			}
			if err := tx.WorkAuthorMappingEntity.WithContext(ctx).Create(&entity.WorkAuthorMappingEntity{
				WorkId:   id,
				AuthorId: author.Id,
			}); err != nil {
				return err
			}
		}

		se := tx.SubjectEntity
		for _, name := range cmd.SubjectNames {
			subject, err := se.WithContext(ctx).Where(se.Name.Eq(name)).First()
			if err != nil {
				if !errors.Is(err, gorm.ErrRecordNotFound) {
					return err
				}
				subject = &entity.SubjectEntity{
					Uid:       uuid.Must(uuid.NewV7()).String(),
					Name:      name,
					CreatedAt: now,
					Version:   1,
				}
				if err := se.WithContext(ctx).Create(subject); err != nil {
					return err
				}
			}
			if err := tx.WorkSubjectMappingEntity.WithContext(ctx).Create(&entity.WorkSubjectMappingEntity{
				WorkId:    id,
				SubjectId: subject.Id,
			}); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *workRepository) Update(ctx context.Context, id int64, cmd command.WorkUpdateCommand) error {
	return r.q.Transaction(func(tx *dao.Query) error {
		w := tx.WorkEntity
		now := time.Now()

		assigns := []field.AssignExpr{w.UpdatedAt.Value(now)}
		if cmd.Title != nil {
			assigns = append(assigns, w.Title.Value(*cmd.Title))
		}
		if cmd.Description != nil {
			assigns = append(assigns, w.Description.Value(*cmd.Description))
		}
		if cmd.CoverUids != nil {
			if s := utils.ToJsonString(cmd.CoverUids); s != nil {
				assigns = append(assigns, w.CoverUid.Value(*s))
			}
		}
		if cmd.FirstPublishDate != nil {
			assigns = append(assigns, w.FirstPublishDate.Value(*cmd.FirstPublishDate))
		}
		if cmd.OlKey != nil {
			assigns = append(assigns, w.OlKey.Value(*cmd.OlKey))
		}
		if _, err := w.WithContext(ctx).Where(w.Id.Eq(id)).UpdateSimple(assigns...); err != nil {
			return err
		}

		if cmd.AuthorUids != nil {
			if err := replaceWorkAuthors(ctx, tx, id, cmd.AuthorUids); err != nil {
				return err
			}
		}
		if cmd.SubjectNames != nil {
			if err := replaceWorkSubjects(ctx, tx, id, cmd.SubjectNames, now); err != nil {
				return err
			}
		}
		return nil
	})
}

func replaceWorkAuthors(ctx context.Context, tx *dao.Query, workID int64, authorUids []string) error {
	wam := tx.WorkAuthorMappingEntity
	if _, err := wam.WithContext(ctx).Where(wam.WorkId.Eq(workID)).Delete(); err != nil {
		return err
	}
	ae := tx.AuthorEntity
	for _, authorUid := range authorUids {
		author, err := ae.WithContext(ctx).Where(ae.Uid.Eq(authorUid)).First()
		if err != nil {
			return err
		}
		if err := wam.WithContext(ctx).Create(&entity.WorkAuthorMappingEntity{
			WorkId:   workID,
			AuthorId: author.Id,
		}); err != nil {
			return err
		}
	}
	return nil
}

func replaceWorkSubjects(ctx context.Context, tx *dao.Query, workID int64, subjectNames []string, now time.Time) error {
	wsm := tx.WorkSubjectMappingEntity
	if _, err := wsm.WithContext(ctx).Where(wsm.WorkId.Eq(workID)).Delete(); err != nil {
		return err
	}
	se := tx.SubjectEntity
	for _, name := range subjectNames {
		subject, err := se.WithContext(ctx).Where(se.Name.Eq(name)).First()
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			subject = &entity.SubjectEntity{
				Uid:       uuid.Must(uuid.NewV7()).String(),
				Name:      name,
				CreatedAt: now,
				Version:   1,
			}
			if err := se.WithContext(ctx).Create(subject); err != nil {
				return err
			}
		}
		if err := wsm.WithContext(ctx).Create(&entity.WorkSubjectMappingEntity{
			WorkId:    workID,
			SubjectId: subject.Id,
		}); err != nil {
			return err
		}
	}
	return nil
}

func toWorkResult(e *entity.WorkEntity, authors []workAuthorRow, subjects []workSubjectRow) *repo.WorkResult {
	authorUids := make([]string, len(authors))
	authorNames := make([]string, len(authors))
	for i, a := range authors {
		authorUids[i] = a.AuthorUid
		authorNames[i] = a.AuthorName
	}

	subjectUids := make([]string, len(subjects))
	subjectNames := make([]string, len(subjects))
	for i, s := range subjects {
		subjectUids[i] = s.SubjectUid
		subjectNames[i] = s.SubjectName
	}

	return &repo.WorkResult{
		Id:               e.Id,
		Uid:              e.Uid,
		Title:            e.Title,
		Description:      e.Description,
		CoverUids:        utils.ParseJsonStringArray(e.CoverUid),
		FirstPublishDate: e.FirstPublishDate,
		OlKey:            e.OlKey,
		CreatedAt:        e.CreatedAt,
		UpdatedAt:        e.UpdatedAt,
		AuthorUids:       authorUids,
		AuthorNames:      authorNames,
		SubjectUids:      subjectUids,
		SubjectNames:     subjectNames,
	}
}

func (r *workRepository) batchLoadAuthors(ctx context.Context, workIds []int64) (map[int64][]workAuthorRow, error) {
	mp := r.q.WorkAuthorMappingEntity
	a := r.q.AuthorEntity

	var rows []workAuthorRow
	err := mp.WithContext(ctx).
		Join(a, a.Id.EqCol(mp.AuthorId), a.DeletedAt.IsNull()).
		Where(mp.WorkId.In(workIds...)).
		Select(mp.WorkId, a.Uid, a.Name).
		Scan(&rows)
	if err != nil {
		return nil, err
	}

	result := make(map[int64][]workAuthorRow, len(workIds))
	for _, row := range rows {
		result[row.WorkId] = append(result[row.WorkId], row)
	}
	return result, nil
}

func (r *workRepository) batchLoadSubjects(ctx context.Context, workIds []int64) (map[int64][]workSubjectRow, error) {
	mp := r.q.WorkSubjectMappingEntity
	s := r.q.SubjectEntity

	var rows []workSubjectRow
	err := mp.WithContext(ctx).
		Join(s, s.Id.EqCol(mp.SubjectId), s.DeletedAt.IsNull()).
		Where(mp.WorkId.In(workIds...)).
		Select(mp.WorkId, s.Uid, s.Name).
		Scan(&rows)
	if err != nil {
		return nil, err
	}

	result := make(map[int64][]workSubjectRow, len(workIds))
	for _, row := range rows {
		result[row.WorkId] = append(result[row.WorkId], row)
	}
	return result, nil
}
