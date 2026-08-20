package restcontroller

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/codejsha/shared-library-go/pkg/pagination"

	"github.com/codejsha/bookstore-microservices/catalog/generated/application/port/openapi"
	"github.com/codejsha/bookstore-microservices/catalog/internal/application/port/repo"
	"github.com/codejsha/bookstore-microservices/catalog/internal/application/usecase"
	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/aggregate"
	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/model/command"
	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/model/option"
	"github.com/codejsha/bookstore-microservices/catalog/internal/infrastructure/httpx"
)

var _ openapi.WorkApi = (*workController)(nil)

type workController struct {
	catalogUseCase usecase.CatalogUseCase
}

func NewWorkController(catalogUseCase usecase.CatalogUseCase) openapi.WorkApi {
	return &workController{catalogUseCase: catalogUseCase}
}

func (c *workController) WorksSearch(
	ctx context.Context,
	title *string,
	authorUid *string,
	subjectUid *string,
	olKey *string,
	size *int32,
	page *int32,
	sort *string,
) (*openapi.WorkFindAllResponse, error) {
	opt := option.NewWorkQueryOption(
		option.WorkQueryOption{}.WithTitle(title),
		option.WorkQueryOption{}.WithAuthorUid(authorUid),
		option.WorkQueryOption{}.WithSubjectUid(subjectUid),
		option.WorkQueryOption{}.WithOlKey(olKey),
		option.WorkQueryOption{}.WithPage(pagination.NewPageOption(size, page, sort)),
	)

	total, works, err := c.catalogUseCase.SearchWorks(ctx, opt)
	if err != nil {
		return nil, err
	}

	items := make([]openapi.WorkItem, len(works))
	for i, w := range works {
		items[i] = toWorkItem(w)
	}
	return &openapi.WorkFindAllResponse{Items: items, Total: total}, nil
}

func (c *workController) WorksCreate(ctx context.Context, req openapi.WorkCreateRequest) error {
	cmd := command.WorkCreateCommand{
		Title:            req.Title,
		Description:      req.Description,
		CoverUids:        derefStringSlice(req.CoverUids),
		FirstPublishDate: req.FirstPublishDate,
		OlKey:            req.OlKey,
		AuthorUids:       req.AuthorUids,
		SubjectNames:     derefStringSlice(req.SubjectNames),
	}

	work, err := c.catalogUseCase.CreateWork(ctx, cmd)
	if err != nil {
		return httpx.MapNotFound(ctx, err)
	}

	dispatchSideEffects(context.WithoutCancel(ctx), "work", work.Uid, "created", logrus.Fields{"title": req.Title})
	return nil
}

func (c *workController) WorksRead(ctx context.Context, uid string) (*openapi.WorkFindResponse, error) {
	work, err := c.catalogUseCase.FindWork(ctx, uid)
	if err != nil {
		return nil, httpx.MapNotFound(ctx, err)
	}

	return &openapi.WorkFindResponse{
		Uid:              work.Uid,
		Title:            work.Title,
		Description:      work.Description,
		CoverUids:        ptrStringSlice(work.CoverUids),
		FirstPublishDate: work.FirstPublishDate,
		OlKey:            work.OlKey,
		Authors:          toAuthorItems(work.Authors),
		Subjects:         toSubjectItems(work.Subjects),
		CreatedAt:        work.CreatedAt,
		UpdatedAt:        work.UpdatedAt,
	}, nil
}

func (c *workController) WorksUpdate(
	ctx context.Context,
	uid string,
	req openapi.WorkUpdateRequest,
) (*openapi.WorkUpdateResponse, error) {
	cmd := command.WorkUpdateCommand{
		Title:            req.Title,
		Description:      req.Description,
		CoverUids:        derefStringSlice(req.CoverUids),
		FirstPublishDate: req.FirstPublishDate,
		OlKey:            req.OlKey,
		AuthorUids:       derefStringSlice(req.AuthorUids),
		SubjectNames:     derefStringSlice(req.SubjectNames),
	}

	work, err := c.catalogUseCase.UpdateWork(ctx, uid, cmd)
	if err != nil {
		return nil, httpx.MapNotFound(ctx, err)
	}

	dispatchSideEffects(context.WithoutCancel(ctx), "work", uid, "updated", logrus.Fields{})

	return &openapi.WorkUpdateResponse{
		Uid:              work.Uid,
		Title:            work.Title,
		Description:      work.Description,
		CoverUids:        ptrStringSlice(work.CoverUids),
		FirstPublishDate: work.FirstPublishDate,
		OlKey:            work.OlKey,
		Authors:          toAuthorItems(work.Authors),
		Subjects:         toSubjectItems(work.Subjects),
		CreatedAt:        work.CreatedAt,
		UpdatedAt:        work.UpdatedAt,
	}, nil
}

func (c *workController) WorksFullTextSearch(
	ctx context.Context,
	q *string,
	authorUid *string,
	subjectUid *string,
	size *int32,
	page *int32,
) (*openapi.WorkSearchResponse, error) {
	opt := option.NewWorkSearchOption(
		option.WorkSearchOption{}.WithQuery(q),
		option.WorkSearchOption{}.WithAuthorUid(authorUid),
		option.WorkSearchOption{}.WithSubjectUid(subjectUid),
		option.WorkSearchOption{}.WithPage(pagination.NewPageOption(size, page, nil)),
	)

	total, hits, err := c.catalogUseCase.FullTextSearchWorks(ctx, opt)
	if err != nil {
		return nil, err
	}

	items := make([]openapi.WorkSearchItem, len(hits))
	for i, h := range hits {
		items[i] = toWorkSearchItem(h)
	}
	return &openapi.WorkSearchResponse{Items: items, Total: total}, nil
}

func (c *workController) WorksEditions(
	ctx context.Context,
	uid string,
	size *int32,
	page *int32,
	sort *string,
) (*openapi.EditionFindAllResponse, error) {
	opt := option.NewEditionQueryOption(
		option.EditionQueryOption{}.WithWorkUid(&uid),
		option.EditionQueryOption{}.WithPage(pagination.NewPageOption(size, page, sort)),
	)

	total, editions, err := c.catalogUseCase.SearchEditions(ctx, opt)
	if err != nil {
		return nil, err
	}

	items := make([]openapi.EditionItem, len(editions))
	for i, e := range editions {
		items[i] = toEditionItem(e)
	}
	return &openapi.EditionFindAllResponse{Items: items, Total: total}, nil
}

// ─── mapping helpers ────────────────────────────────────────────────────────

func toWorkItem(w *aggregate.WorkAggregate) openapi.WorkItem {
	return openapi.WorkItem{
		Uid:              w.Uid,
		Title:            w.Title,
		Description:      w.Description,
		CoverUids:        ptrStringSlice(w.CoverUids),
		FirstPublishDate: w.FirstPublishDate,
		OlKey:            w.OlKey,
		Authors:          toAuthorItems(w.Authors),
		Subjects:         toSubjectItems(w.Subjects),
		CreatedAt:        w.CreatedAt,
		UpdatedAt:        w.UpdatedAt,
	}
}

func toWorkSearchItem(r *repo.WorkSearchResult) openapi.WorkSearchItem {
	item := openapi.WorkSearchItem{
		Uid:              r.Uid,
		Title:            r.Title,
		Description:      r.Description,
		CoverUids:        ptrStringSlice(r.CoverUids),
		FirstPublishDate: r.FirstPublishDate,
		OlKey:            r.OlKey,
		Authors:          toSearchAuthorItems(r.Authors),
		Subjects:         toSearchSubjectItems(r.Subjects),
		Score:            r.Score,
		CreatedAt:        r.CreatedAt,
		UpdatedAt:        r.UpdatedAt,
	}
	if r.Highlight != nil {
		item.Highlight = &openapi.WorkSearchHighlight{
			Title:       ptrStringSlice(r.Highlight.Title),
			Description: ptrStringSlice(r.Highlight.Description),
		}
	}
	return item
}

func toSearchAuthorItems(src []repo.WorkSearchAuthor) []openapi.AuthorItem {
	out := make([]openapi.AuthorItem, len(src))
	for i, a := range src {
		out[i] = openapi.AuthorItem{
			Uid:            a.Uid,
			Name:           a.Name,
			Bio:            a.Bio,
			BirthDate:      a.BirthDate,
			DeathDate:      a.DeathDate,
			OlKey:          a.OlKey,
			AlternateNames: ptrStringSlice(a.AlternateNames),
		}
	}
	return out
}

func toSearchSubjectItems(src []repo.WorkSearchSubject) []openapi.SubjectItem {
	out := make([]openapi.SubjectItem, len(src))
	for i, s := range src {
		out[i] = openapi.SubjectItem{Uid: s.Uid, Name: s.Name}
	}
	return out
}

func toEditionItem(e *aggregate.EditionAggregate) openapi.EditionItem {
	item := openapi.EditionItem{
		Uid:            e.Uid,
		Title:          e.Title,
		Isbn10:         e.Isbn10,
		Isbn13:         e.Isbn13,
		NumberOfPages:  e.NumberOfPages,
		PublishDate:    e.PublishDate,
		CoverUids:      ptrStringSlice(e.CoverUids),
		Languages:      ptrStringSlice(e.Languages),
		PhysicalFormat: e.PhysicalFormat,
		Description:    e.Description,
		OlKey:          e.OlKey,
		CreatedAt:      e.CreatedAt,
		UpdatedAt:      e.UpdatedAt,
	}
	if e.Work != nil {
		item.Work = toWorkItem(e.Work)
	}
	if e.Publisher != nil {
		pub := toPublisherItem(e.Publisher)
		item.Publisher = &pub
	}
	return item
}

func toAuthorItem(a *aggregate.AuthorAggregate) openapi.AuthorItem {
	return openapi.AuthorItem{
		Uid:            a.Uid,
		Name:           a.Name,
		Bio:            a.Bio,
		BirthDate:      a.BirthDate,
		DeathDate:      a.DeathDate,
		PhotoUids:      ptrStringSlice(a.PhotoUids),
		AlternateNames: ptrStringSlice(a.AlternateNames),
		OlKey:          a.OlKey,
	}
}

func toAuthorItems(authors []*aggregate.AuthorAggregate) []openapi.AuthorItem {
	items := make([]openapi.AuthorItem, len(authors))
	for i, a := range authors {
		items[i] = toAuthorItem(a)
	}
	return items
}

func toSubjectItem(s *aggregate.SubjectAggregate) openapi.SubjectItem {
	return openapi.SubjectItem{
		Uid:  s.Uid,
		Name: s.Name,
	}
}

func toSubjectItems(subjects []*aggregate.SubjectAggregate) []openapi.SubjectItem {
	items := make([]openapi.SubjectItem, len(subjects))
	for i, s := range subjects {
		items[i] = toSubjectItem(s)
	}
	return items
}

func toPublisherItem(p *aggregate.PublisherAggregate) openapi.PublisherItem {
	return openapi.PublisherItem{
		Uid:     p.Uid,
		Name:    p.Name,
		Address: p.Address,
		OlKey:   p.OlKey,
	}
}

// ─── slice helpers ──────────────────────────────────────────────────────────

func ptrStringSlice(s []string) *[]string {
	if len(s) == 0 {
		return nil
	}
	return &s
}

func derefStringSlice(s *[]string) []string {
	if s == nil {
		return nil
	}
	return *s
}
