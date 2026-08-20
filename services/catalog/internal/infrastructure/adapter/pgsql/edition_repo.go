package pgsql

import (
	"context"
	"encoding/json"
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

var _ repo.EditionRepo = (*editionRepository)(nil)

type editionRepository struct {
	q *dao.Query
}

func NewEditionRepository(dataSource *database.DataSource) repo.EditionRepo {
	return &editionRepository{q: dao.Use(dataSource.DB())}
}

type editionRow struct {
	Id             int64      `gorm:"column:id"`
	Uid            string     `gorm:"column:uid"`
	Title          string     `gorm:"column:title"`
	Isbn10         *string    `gorm:"column:isbn10"`
	Isbn13         *string    `gorm:"column:isbn13"`
	NumberOfPage   *int32     `gorm:"column:number_of_pages"`
	PublishDate    *string    `gorm:"column:publish_date"`
	CoverUid       *string    `gorm:"column:cover_uids"`
	Language       *string    `gorm:"column:languages"`
	PhysicalFormat *string    `gorm:"column:physical_format"`
	Description    *string    `gorm:"column:description"`
	WorkUid        string     `gorm:"column:work_uid"`
	WorkTitle      string     `gorm:"column:work_title"`
	PublisherUid   *string    `gorm:"column:publisher_uid"`
	PublisherName  *string    `gorm:"column:publisher_name"`
	OlKey          *string    `gorm:"column:ol_key"`
	CreatedAt      time.Time  `gorm:"column:created_at"`
	UpdatedAt      *time.Time `gorm:"column:updated_at"`
}

func (r *editionRepository) FindAll(ctx context.Context, opt option.EditionQueryOption) (int64, []*repo.EditionResult, error) {
	e := r.q.EditionEntity
	w := r.q.WorkEntity
	p := r.q.PublisherEntity

	q := e.WithContext(ctx).
		Join(w, w.Id.EqCol(e.WorkId), w.DeletedAt.IsNull()).
		LeftJoin(p, p.Id.EqCol(e.PublisherId), p.DeletedAt.IsNull())

	if v := opt.Title(); v != nil && *v != "" {
		q = q.Where(gen.Cond(gorm.Expr("edition.title ILIKE ?", "%"+*v+"%"))...)
	}
	if uid := opt.WorkUid(); uid != nil && *uid != "" {
		q = q.Where(w.Uid.Eq(*uid))
	}
	if uid := opt.PublisherUid(); uid != nil && *uid != "" {
		q = q.Where(p.Uid.Eq(*uid))
	}
	if v := opt.Isbn(); v != nil && *v != "" {
		isbn := "%" + *v + "%"
		q = q.Where(gen.Cond(gorm.Expr("(edition.isbn10 ILIKE ? OR edition.isbn13 ILIKE ?)", isbn, isbn))...)
	}
	if v := opt.OlKey(); v != nil && *v != "" {
		q = q.Where(e.OlKey.Eq(*v))
	}
	if v := opt.Language(); v != nil && *v != "" {
		langJSON, err := json.Marshal([]string{*v})
		if err != nil {
			return 0, nil, err
		}
		q = q.Where(gen.Cond(gorm.Expr("edition.languages @> ?::jsonb", string(langJSON)))...)
	}

	q = q.Select(r.editionSelectExprs()...)

	q = q.Order(buildOrderExprs(opt.Page().GetSort(), map[string]field.OrderExpr{
		"title":      e.Title,
		"created_at": e.CreatedAt,
		"updated_at": e.UpdatedAt,
	}, e.Id)...)

	offset, limit := utils.PageOffsetLimit(opt.Page())
	var rows []editionRow
	total, err := q.ScanByPage(&rows, offset, limit)
	if err != nil {
		return 0, nil, err
	}

	results := make([]*repo.EditionResult, len(rows))
	for i, row := range rows {
		results[i] = toEditionResult(&row)
	}
	return total, results, nil
}

func (r *editionRepository) FindOne(ctx context.Context, id int64) (*repo.EditionResult, error) {
	return r.findOneBy(ctx, r.q.EditionEntity.Id.Eq(id))
}

func (r *editionRepository) FindByUid(ctx context.Context, uid string) (*repo.EditionResult, error) {
	return r.findOneBy(ctx, r.q.EditionEntity.Uid.Eq(uid))
}

func (r *editionRepository) findOneBy(ctx context.Context, cond field.Expr) (*repo.EditionResult, error) {
	e := r.q.EditionEntity
	w := r.q.WorkEntity
	p := r.q.PublisherEntity

	var row editionRow
	err := e.WithContext(ctx).
		Join(w, w.Id.EqCol(e.WorkId), w.DeletedAt.IsNull()).
		LeftJoin(p, p.Id.EqCol(e.PublisherId), p.DeletedAt.IsNull()).
		Select(r.editionSelectExprs()...).
		Where(cond).
		Limit(1).
		Scan(&row)
	if err != nil {
		return nil, err
	}
	if row.Id == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return toEditionResult(&row), nil
}

func (r *editionRepository) Create(ctx context.Context, cmd command.EditionCreateCommand, workId int64, publisherId *int64) (int64, error) {
	e := &entity.EditionEntity{
		Uid:            uuid.Must(uuid.NewV7()).String(),
		Title:          cmd.Title,
		Isbn10:         cmd.Isbn10,
		Isbn13:         cmd.Isbn13,
		NumberOfPage:   cmd.NumberOfPages,
		PublishDate:    cmd.PublishDate,
		CoverUid:       utils.ToJsonString(cmd.CoverUids),
		Language:       utils.ToJsonString(cmd.Languages),
		PhysicalFormat: cmd.PhysicalFormat,
		Description:    cmd.Description,
		WorkId:         workId,
		PublisherId:    publisherId,
		OlKey:          cmd.OlKey,
		CreatedAt:      time.Now(),
		Version:        1,
	}
	if err := r.q.EditionEntity.WithContext(ctx).Create(e); err != nil {
		return 0, err
	}
	return e.Id, nil
}

func (r *editionRepository) Update(ctx context.Context, id int64, cmd command.EditionUpdateCommand, workId *int64, publisherId *int64) error {
	e := r.q.EditionEntity
	assigns := []field.AssignExpr{e.UpdatedAt.Value(time.Now())}
	if cmd.Title != nil {
		assigns = append(assigns, e.Title.Value(*cmd.Title))
	}
	if cmd.Isbn10 != nil {
		assigns = append(assigns, e.Isbn10.Value(*cmd.Isbn10))
	}
	if cmd.Isbn13 != nil {
		assigns = append(assigns, e.Isbn13.Value(*cmd.Isbn13))
	}
	if cmd.NumberOfPages != nil {
		assigns = append(assigns, e.NumberOfPage.Value(*cmd.NumberOfPages))
	}
	if cmd.PublishDate != nil {
		assigns = append(assigns, e.PublishDate.Value(*cmd.PublishDate))
	}
	if cmd.CoverUids != nil {
		if s := utils.ToJsonString(cmd.CoverUids); s != nil {
			assigns = append(assigns, e.CoverUid.Value(*s))
		}
	}
	if cmd.Languages != nil {
		if s := utils.ToJsonString(cmd.Languages); s != nil {
			assigns = append(assigns, e.Language.Value(*s))
		}
	}
	if cmd.PhysicalFormat != nil {
		assigns = append(assigns, e.PhysicalFormat.Value(*cmd.PhysicalFormat))
	}
	if cmd.Description != nil {
		assigns = append(assigns, e.Description.Value(*cmd.Description))
	}
	if workId != nil {
		assigns = append(assigns, e.WorkId.Value(*workId))
	}
	if publisherId != nil {
		assigns = append(assigns, e.PublisherId.Value(*publisherId))
	}
	if cmd.OlKey != nil {
		assigns = append(assigns, e.OlKey.Value(*cmd.OlKey))
	}
	_, err := e.WithContext(ctx).Where(e.Id.Eq(id)).UpdateSimple(assigns...)
	return err
}

func (r *editionRepository) editionSelectExprs() []field.Expr {
	e := r.q.EditionEntity
	w := r.q.WorkEntity
	p := r.q.PublisherEntity
	return []field.Expr{
		e.Id, e.Uid, e.Title, e.Isbn10, e.Isbn13,
		e.NumberOfPage, e.PublishDate, e.CoverUid,
		e.Language, e.PhysicalFormat, e.Description,
		w.Uid.As("work_uid"), w.Title.As("work_title"),
		p.Uid.As("publisher_uid"), p.Name.As("publisher_name"),
		e.OlKey, e.CreatedAt, e.UpdatedAt,
	}
}

func toEditionResult(row *editionRow) *repo.EditionResult {
	return &repo.EditionResult{
		Id:             row.Id,
		Uid:            row.Uid,
		Title:          row.Title,
		Isbn10:         row.Isbn10,
		Isbn13:         row.Isbn13,
		NumberOfPages:  row.NumberOfPage,
		PublishDate:    row.PublishDate,
		CoverUids:      utils.ParseJsonStringArray(row.CoverUid),
		Languages:      utils.ParseJsonStringArray(row.Language),
		PhysicalFormat: row.PhysicalFormat,
		Description:    row.Description,
		WorkUid:        row.WorkUid,
		WorkTitle:      row.WorkTitle,
		PublisherUid:   row.PublisherUid,
		PublisherName:  row.PublisherName,
		OlKey:          row.OlKey,
		CreatedAt:      row.CreatedAt,
		UpdatedAt:      row.UpdatedAt,
	}
}
