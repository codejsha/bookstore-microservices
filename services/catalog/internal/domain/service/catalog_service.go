package service

import (
	"context"

	"golang.org/x/sync/errgroup"

	"github.com/codejsha/bookstore-microservices/catalog/internal/application/port/repo"
	"github.com/codejsha/bookstore-microservices/catalog/internal/application/usecase"
	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/aggregate"
	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/model/command"
	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/model/option"
)

var _ usecase.CatalogUseCase = (*catalogService)(nil)

type catalogService struct {
	workRepo       repo.WorkRepo
	workSearchRepo repo.WorkSearchRepo
	editionRepo    repo.EditionRepo
	authorRepo     repo.AuthorRepo
	publisherRepo  repo.PublisherRepo
	subjectRepo    repo.SubjectRepo
}

func NewCatalogService(
	workRepo repo.WorkRepo,
	workSearchRepo repo.WorkSearchRepo,
	editionRepo repo.EditionRepo,
	authorRepo repo.AuthorRepo,
	publisherRepo repo.PublisherRepo,
	subjectRepo repo.SubjectRepo,
) usecase.CatalogUseCase {
	return &catalogService{
		workRepo:       workRepo,
		workSearchRepo: workSearchRepo,
		editionRepo:    editionRepo,
		authorRepo:     authorRepo,
		publisherRepo:  publisherRepo,
		subjectRepo:    subjectRepo,
	}
}

// ─── Work lifecycle ─────────────────────────────────────────────────────────

func (s catalogService) SearchWorks(ctx context.Context, opt option.WorkQueryOption) (int64, []*aggregate.WorkAggregate, error) {
	total, works, err := s.workSearchRepo.FindAll(ctx, opt)
	if err != nil {
		return 0, nil, err
	}
	aggs := make([]*aggregate.WorkAggregate, len(works))
	for i, w := range works {
		aggs[i] = s.toWorkAggregateFromSearch(w)
	}
	return total, aggs, nil
}

func (s catalogService) FullTextSearchWorks(ctx context.Context, opt option.WorkSearchOption) (int64, []*repo.WorkSearchResult, error) {
	return s.workSearchRepo.Search(ctx, opt)
}

func (s catalogService) FindWork(ctx context.Context, uid string) (*aggregate.WorkAggregate, error) {
	w, err := s.workRepo.FindByUid(ctx, uid)
	if err != nil {
		return nil, err
	}
	return s.toWorkAggregate(w), nil
}

func (s catalogService) FindWorkWithRelated(
	ctx context.Context, uid string,
) (*aggregate.WorkAggregate, []*repo.WorkSearchResult, error) {
	work, err := s.workRepo.FindByUid(ctx, uid)
	if err != nil {
		return nil, nil, err
	}

	var (
		editions []*repo.EditionResult
		related  []*repo.WorkSearchResult
	)
	g, gctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		workUid := work.Uid
		opt := option.NewEditionQueryOption(option.EditionQueryOption{}.WithWorkUid(&workUid))
		_, eds, err := s.editionRepo.FindAll(gctx, opt)
		if err != nil {
			return err
		}
		editions = eds
		return nil
	})

	if len(work.AuthorUids) > 0 {
		authorUid := work.AuthorUids[0]
		g.Go(func() error {
			searchOpt := option.NewWorkSearchOption(
				option.WorkSearchOption{}.WithAuthorUid(&authorUid),
			)
			_, results, err := s.workSearchRepo.Search(gctx, searchOpt)
			if err != nil {
				return err
			}
			out := make([]*repo.WorkSearchResult, 0, len(results))
			for _, r := range results {
				if r.Uid == work.Uid {
					continue
				}
				out = append(out, r)
			}
			related = out
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, nil, err
	}

	agg := s.toWorkAggregate(work)
	if len(editions) > 0 {
		agg.Editions = make([]*aggregate.EditionAggregate, len(editions))
		for i, e := range editions {
			agg.Editions[i] = s.toEditionAggregate(e)
		}
	}
	if related == nil {
		related = []*repo.WorkSearchResult{}
	}
	return agg, related, nil
}

func (s catalogService) CreateWork(ctx context.Context, cmd command.WorkCreateCommand) (*aggregate.WorkAggregate, error) {
	if err := cmd.Validate(); err != nil {
		return nil, err
	}
	id, err := s.workRepo.Create(ctx, cmd)
	if err != nil {
		return nil, err
	}
	w, err := s.workRepo.FindOne(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.toWorkAggregate(w), nil
}

func (s catalogService) UpdateWork(ctx context.Context, uid string, cmd command.WorkUpdateCommand) (*aggregate.WorkAggregate, error) {
	if err := cmd.Validate(); err != nil {
		return nil, err
	}
	existing, err := s.workRepo.FindByUid(ctx, uid)
	if err != nil {
		return nil, err
	}
	if err := s.workRepo.Update(ctx, existing.Id, cmd); err != nil {
		return nil, err
	}
	w, err := s.workRepo.FindOne(ctx, existing.Id)
	if err != nil {
		return nil, err
	}
	return s.toWorkAggregate(w), nil
}

// ─── Edition lifecycle ──────────────────────────────────────────────────────

func (s catalogService) SearchEditions(ctx context.Context, opt option.EditionQueryOption) (int64, []*aggregate.EditionAggregate, error) {
	total, editions, err := s.editionRepo.FindAll(ctx, opt)
	if err != nil {
		return 0, nil, err
	}
	aggs := make([]*aggregate.EditionAggregate, len(editions))
	for i, e := range editions {
		aggs[i] = s.toEditionAggregate(e)
	}
	return total, aggs, nil
}

func (s catalogService) FindEdition(ctx context.Context, uid string) (*aggregate.EditionAggregate, error) {
	e, err := s.editionRepo.FindByUid(ctx, uid)
	if err != nil {
		return nil, err
	}
	return s.toEditionAggregate(e), nil
}

func (s catalogService) CreateEdition(ctx context.Context, cmd command.EditionCreateCommand) (*aggregate.EditionAggregate, error) {
	if err := cmd.Validate(); err != nil {
		return nil, err
	}
	work, publisherId, err := s.resolveEditionRefs(ctx, cmd.WorkUid, cmd.PublisherUid)
	if err != nil {
		return nil, err
	}
	id, err := s.editionRepo.Create(ctx, cmd, work.Id, publisherId)
	if err != nil {
		return nil, err
	}
	e, err := s.editionRepo.FindOne(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.toEditionAggregate(e), nil
}

func (s catalogService) UpdateEdition(ctx context.Context, uid string, cmd command.EditionUpdateCommand) (*aggregate.EditionAggregate, error) {
	if err := cmd.Validate(); err != nil {
		return nil, err
	}
	var (
		existing    *repo.EditionResult
		workId      *int64
		publisherId *int64
	)
	g, gctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		e, err := s.editionRepo.FindByUid(gctx, uid)
		if err != nil {
			return err
		}
		existing = e
		return nil
	})
	if cmd.WorkUid != nil {
		workUid := *cmd.WorkUid
		g.Go(func() error {
			w, err := s.workRepo.FindByUid(gctx, workUid)
			if err != nil {
				return err
			}
			workId = &w.Id
			return nil
		})
	}
	if cmd.PublisherUid != nil {
		publisherUid := *cmd.PublisherUid
		g.Go(func() error {
			p, err := s.publisherRepo.FindByUid(gctx, publisherUid)
			if err != nil {
				return err
			}
			publisherId = &p.Id
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}
	if err := s.editionRepo.Update(ctx, existing.Id, cmd, workId, publisherId); err != nil {
		return nil, err
	}
	e, err := s.editionRepo.FindOne(ctx, existing.Id)
	if err != nil {
		return nil, err
	}
	return s.toEditionAggregate(e), nil
}

func (s catalogService) resolveEditionRefs(
	ctx context.Context, workUid string, publisherUid *string,
) (*repo.WorkResult, *int64, error) {
	var (
		work        *repo.WorkResult
		publisherId *int64
	)
	g, gctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		w, err := s.workRepo.FindByUid(gctx, workUid)
		if err != nil {
			return err
		}
		work = w
		return nil
	})
	if publisherUid != nil {
		pubUid := *publisherUid
		g.Go(func() error {
			p, err := s.publisherRepo.FindByUid(gctx, pubUid)
			if err != nil {
				return err
			}
			publisherId = &p.Id
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return nil, nil, err
	}
	return work, publisherId, nil
}

// ─── Author management ──────────────────────────────────────────────────────

func (s catalogService) FindAllAuthors(ctx context.Context, opt option.AuthorQueryOption) (int64, []*aggregate.AuthorAggregate, error) {
	total, authors, err := s.authorRepo.FindAll(ctx, opt)
	if err != nil {
		return 0, nil, err
	}
	aggs := make([]*aggregate.AuthorAggregate, len(authors))
	for i, a := range authors {
		aggs[i] = s.toAuthorAggregate(a)
	}
	return total, aggs, nil
}

func (s catalogService) FindAuthor(ctx context.Context, uid string) (*aggregate.AuthorAggregate, error) {
	a, err := s.authorRepo.FindByUid(ctx, uid)
	if err != nil {
		return nil, err
	}
	return s.toAuthorAggregate(a), nil
}

func (s catalogService) CreateAuthor(ctx context.Context, cmd command.AuthorCreateCommand) (*aggregate.AuthorAggregate, error) {
	if err := cmd.Validate(); err != nil {
		return nil, err
	}
	id, err := s.authorRepo.Create(ctx, cmd)
	if err != nil {
		return nil, err
	}
	a, err := s.authorRepo.FindOne(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.toAuthorAggregate(a), nil
}

func (s catalogService) UpdateAuthor(ctx context.Context, uid string, cmd command.AuthorUpdateCommand) (*aggregate.AuthorAggregate, error) {
	if err := cmd.Validate(); err != nil {
		return nil, err
	}
	existing, err := s.authorRepo.FindByUid(ctx, uid)
	if err != nil {
		return nil, err
	}
	if err := s.authorRepo.Update(ctx, existing.Id, cmd); err != nil {
		return nil, err
	}
	a, err := s.authorRepo.FindOne(ctx, existing.Id)
	if err != nil {
		return nil, err
	}
	return s.toAuthorAggregate(a), nil
}

// ─── Publisher management ───────────────────────────────────────────────────

func (s catalogService) FindAllPublishers(ctx context.Context, opt option.PublisherQueryOption) (int64, []*aggregate.PublisherAggregate, error) {
	total, publishers, err := s.publisherRepo.FindAll(ctx, opt)
	if err != nil {
		return 0, nil, err
	}
	aggs := make([]*aggregate.PublisherAggregate, len(publishers))
	for i, p := range publishers {
		aggs[i] = &aggregate.PublisherAggregate{
			Uid:     p.Uid,
			Name:    p.Name,
			Address: p.Address,
			OlKey:   p.OlKey,
		}
	}
	return total, aggs, nil
}

func (s catalogService) FindPublisher(ctx context.Context, uid string) (*aggregate.PublisherAggregate, error) {
	p, err := s.publisherRepo.FindByUid(ctx, uid)
	if err != nil {
		return nil, err
	}
	return &aggregate.PublisherAggregate{
		Uid:     p.Uid,
		Name:    p.Name,
		Address: p.Address,
		OlKey:   p.OlKey,
	}, nil
}

func (s catalogService) CreatePublisher(ctx context.Context, cmd command.PublisherCreateCommand) (*aggregate.PublisherAggregate, error) {
	if err := cmd.Validate(); err != nil {
		return nil, err
	}
	id, err := s.publisherRepo.Create(ctx, cmd)
	if err != nil {
		return nil, err
	}
	p, err := s.publisherRepo.FindOne(ctx, id)
	if err != nil {
		return nil, err
	}
	return &aggregate.PublisherAggregate{
		Uid:     p.Uid,
		Name:    p.Name,
		Address: p.Address,
		OlKey:   p.OlKey,
	}, nil
}

func (s catalogService) UpdatePublisher(ctx context.Context, uid string, cmd command.PublisherUpdateCommand) (*aggregate.PublisherAggregate, error) {
	if err := cmd.Validate(); err != nil {
		return nil, err
	}
	existing, err := s.publisherRepo.FindByUid(ctx, uid)
	if err != nil {
		return nil, err
	}
	if err := s.publisherRepo.Update(ctx, existing.Id, cmd); err != nil {
		return nil, err
	}
	p, err := s.publisherRepo.FindOne(ctx, existing.Id)
	if err != nil {
		return nil, err
	}
	return &aggregate.PublisherAggregate{
		Uid:     p.Uid,
		Name:    p.Name,
		Address: p.Address,
		OlKey:   p.OlKey,
	}, nil
}

// ─── Subject management ─────────────────────────────────────────────────────

func (s catalogService) FindAllSubjects(ctx context.Context, opt option.SubjectQueryOption) (int64, []*aggregate.SubjectAggregate, error) {
	total, subjects, err := s.subjectRepo.FindAll(ctx, opt)
	if err != nil {
		return 0, nil, err
	}
	aggs := make([]*aggregate.SubjectAggregate, len(subjects))
	for i, sub := range subjects {
		aggs[i] = &aggregate.SubjectAggregate{
			Uid:  sub.Uid,
			Name: sub.Name,
		}
	}
	return total, aggs, nil
}

func (s catalogService) FindSubject(ctx context.Context, uid string) (*aggregate.SubjectAggregate, error) {
	sub, err := s.subjectRepo.FindByUid(ctx, uid)
	if err != nil {
		return nil, err
	}
	return &aggregate.SubjectAggregate{
		Uid:  sub.Uid,
		Name: sub.Name,
	}, nil
}

func (s catalogService) CreateSubject(ctx context.Context, cmd command.SubjectCreateCommand) (*aggregate.SubjectAggregate, error) {
	if err := cmd.Validate(); err != nil {
		return nil, err
	}
	if err := s.ensureSubjectNameFree(ctx, cmd.Name, nil); err != nil {
		return nil, err
	}
	id, err := s.subjectRepo.Create(ctx, cmd)
	if err != nil {
		return nil, err
	}
	sub, err := s.subjectRepo.FindOne(ctx, id)
	if err != nil {
		return nil, err
	}
	return &aggregate.SubjectAggregate{
		Uid:  sub.Uid,
		Name: sub.Name,
	}, nil
}

func (s catalogService) UpdateSubject(ctx context.Context, uid string, cmd command.SubjectUpdateCommand) (*aggregate.SubjectAggregate, error) {
	if err := cmd.Validate(); err != nil {
		return nil, err
	}
	existing, err := s.subjectRepo.FindByUid(ctx, uid)
	if err != nil {
		return nil, err
	}
	if cmd.Name != nil {
		if err := s.ensureSubjectNameFree(ctx, *cmd.Name, &existing.Id); err != nil {
			return nil, err
		}
	}
	if err := s.subjectRepo.Update(ctx, existing.Id, cmd); err != nil {
		return nil, err
	}
	sub, err := s.subjectRepo.FindOne(ctx, existing.Id)
	if err != nil {
		return nil, err
	}
	return &aggregate.SubjectAggregate{
		Uid:  sub.Uid,
		Name: sub.Name,
	}, nil
}

func (s catalogService) ensureSubjectNameFree(ctx context.Context, name string, selfId *int64) error {
	found, err := s.subjectRepo.FindByName(ctx, name)
	if err != nil {
		return err
	}
	if found == nil {
		return nil
	}
	if selfId != nil && found.Id == *selfId {
		return nil
	}
	return repo.ErrAlreadyExists
}

// ─── mapping helpers ────────────────────────────────────────────────

func (s catalogService) toWorkAggregate(w *repo.WorkResult) *aggregate.WorkAggregate {
	agg := &aggregate.WorkAggregate{
		Uid:              w.Uid,
		Title:            w.Title,
		Description:      w.Description,
		CoverUids:        w.CoverUids,
		FirstPublishDate: w.FirstPublishDate,
		OlKey:            w.OlKey,
		CreatedAt:        w.CreatedAt,
		UpdatedAt:        w.UpdatedAt,
	}
	if len(w.AuthorUids) > 0 {
		agg.Authors = make([]*aggregate.AuthorAggregate, len(w.AuthorUids))
		for i := range w.AuthorUids {
			agg.Authors[i] = &aggregate.AuthorAggregate{
				Uid:  w.AuthorUids[i],
				Name: w.AuthorNames[i],
			}
		}
	}
	if len(w.SubjectUids) > 0 {
		agg.Subjects = make([]*aggregate.SubjectAggregate, len(w.SubjectUids))
		for i := range w.SubjectUids {
			agg.Subjects[i] = &aggregate.SubjectAggregate{
				Uid:  w.SubjectUids[i],
				Name: w.SubjectNames[i],
			}
		}
	}
	return agg
}

func (s catalogService) toWorkAggregateFromSearch(w *repo.WorkSearchResult) *aggregate.WorkAggregate {
	agg := &aggregate.WorkAggregate{
		Uid:              w.Uid,
		Title:            w.Title,
		Description:      w.Description,
		CoverUids:        w.CoverUids,
		FirstPublishDate: w.FirstPublishDate,
		OlKey:            w.OlKey,
		CreatedAt:        w.CreatedAt,
		UpdatedAt:        w.UpdatedAt,
	}
	if len(w.Authors) > 0 {
		agg.Authors = make([]*aggregate.AuthorAggregate, len(w.Authors))
		for i, a := range w.Authors {
			agg.Authors[i] = &aggregate.AuthorAggregate{
				Uid:  a.Uid,
				Name: a.Name,
			}
		}
	}
	if len(w.Subjects) > 0 {
		agg.Subjects = make([]*aggregate.SubjectAggregate, len(w.Subjects))
		for i, sub := range w.Subjects {
			agg.Subjects[i] = &aggregate.SubjectAggregate{
				Uid:  sub.Uid,
				Name: sub.Name,
			}
		}
	}
	return agg
}

func (s catalogService) toEditionAggregate(e *repo.EditionResult) *aggregate.EditionAggregate {
	agg := &aggregate.EditionAggregate{
		Uid:            e.Uid,
		Title:          e.Title,
		Isbn10:         e.Isbn10,
		Isbn13:         e.Isbn13,
		NumberOfPages:  e.NumberOfPages,
		PublishDate:    e.PublishDate,
		CoverUids:      e.CoverUids,
		Languages:      e.Languages,
		PhysicalFormat: e.PhysicalFormat,
		Description:    e.Description,
		OlKey:          e.OlKey,
		CreatedAt:      e.CreatedAt,
		UpdatedAt:      e.UpdatedAt,
		Work: &aggregate.WorkAggregate{
			Uid:   e.WorkUid,
			Title: e.WorkTitle,
		},
	}

	if e.PublisherUid != nil {
		name := ""
		if e.PublisherName != nil {
			name = *e.PublisherName
		}
		agg.Publisher = &aggregate.PublisherAggregate{
			Uid:  *e.PublisherUid,
			Name: name,
		}
	}
	return agg
}

func (s catalogService) toAuthorAggregate(a *repo.AuthorResult) *aggregate.AuthorAggregate {
	return &aggregate.AuthorAggregate{
		Uid:            a.Uid,
		Name:           a.Name,
		Bio:            a.Bio,
		BirthDate:      a.BirthDate,
		DeathDate:      a.DeathDate,
		PhotoUids:      a.PhotoUids,
		AlternateNames: a.AlternateNames,
		OlKey:          a.OlKey,
	}
}
