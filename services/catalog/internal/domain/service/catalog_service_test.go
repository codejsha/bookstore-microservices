package service

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/codejsha/bookstore-microservices/catalog/internal/application/port/repo"
	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/model/command"
	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/model/option"
)

const (
	cmdAuthorUid    = "0197f7c0-1c1a-7000-8000-000000000001"
	cmdWorkUid      = "0197f7c0-1c1a-7000-8000-0000000000a1"
	cmdPublisherUid = "0197f7c0-1c1a-7000-8000-0000000000b1"
)

// ─── stub repositories ──────────────────────────────────────────────────────

type stubWorkRepo struct {
	findAllFn   func(ctx context.Context, opt option.WorkQueryOption) (int64, []*repo.WorkResult, error)
	findOneFn   func(ctx context.Context, id int64) (*repo.WorkResult, error)
	findByUidFn func(ctx context.Context, uid string) (*repo.WorkResult, error)
	createFn    func(ctx context.Context, cmd command.WorkCreateCommand) (int64, error)
	updateFn    func(ctx context.Context, id int64, cmd command.WorkUpdateCommand) error
}

func (s *stubWorkRepo) FindAll(ctx context.Context, opt option.WorkQueryOption) (int64, []*repo.WorkResult, error) {
	return s.findAllFn(ctx, opt)
}
func (s *stubWorkRepo) FindOne(ctx context.Context, id int64) (*repo.WorkResult, error) {
	return s.findOneFn(ctx, id)
}
func (s *stubWorkRepo) FindByUid(ctx context.Context, uid string) (*repo.WorkResult, error) {
	return s.findByUidFn(ctx, uid)
}
func (s *stubWorkRepo) Create(ctx context.Context, cmd command.WorkCreateCommand) (int64, error) {
	return s.createFn(ctx, cmd)
}
func (s *stubWorkRepo) Update(ctx context.Context, id int64, cmd command.WorkUpdateCommand) error {
	return s.updateFn(ctx, id, cmd)
}

type stubWorkSearchRepo struct {
	findAllFn func(ctx context.Context, opt option.WorkQueryOption) (int64, []*repo.WorkSearchResult, error)
	searchFn  func(ctx context.Context, opt option.WorkSearchOption) (int64, []*repo.WorkSearchResult, error)
}

func (s *stubWorkSearchRepo) FindAll(ctx context.Context, opt option.WorkQueryOption) (int64, []*repo.WorkSearchResult, error) {
	return s.findAllFn(ctx, opt)
}

func (s *stubWorkSearchRepo) Search(ctx context.Context, opt option.WorkSearchOption) (int64, []*repo.WorkSearchResult, error) {
	return s.searchFn(ctx, opt)
}

type stubEditionRepo struct {
	findAllFn   func(ctx context.Context, opt option.EditionQueryOption) (int64, []*repo.EditionResult, error)
	findOneFn   func(ctx context.Context, id int64) (*repo.EditionResult, error)
	findByUidFn func(ctx context.Context, uid string) (*repo.EditionResult, error)
	createFn    func(ctx context.Context, cmd command.EditionCreateCommand, workId int64, publisherId *int64) (int64, error)
	updateFn    func(ctx context.Context, id int64, cmd command.EditionUpdateCommand, workId *int64, publisherId *int64) error
}

func (s *stubEditionRepo) FindAll(ctx context.Context, opt option.EditionQueryOption) (int64, []*repo.EditionResult, error) {
	return s.findAllFn(ctx, opt)
}
func (s *stubEditionRepo) FindOne(ctx context.Context, id int64) (*repo.EditionResult, error) {
	return s.findOneFn(ctx, id)
}
func (s *stubEditionRepo) FindByUid(ctx context.Context, uid string) (*repo.EditionResult, error) {
	return s.findByUidFn(ctx, uid)
}
func (s *stubEditionRepo) Create(ctx context.Context, cmd command.EditionCreateCommand, workId int64, publisherId *int64) (int64, error) {
	return s.createFn(ctx, cmd, workId, publisherId)
}
func (s *stubEditionRepo) Update(ctx context.Context, id int64, cmd command.EditionUpdateCommand, workId *int64, publisherId *int64) error {
	return s.updateFn(ctx, id, cmd, workId, publisherId)
}

type stubAuthorRepo struct {
	findAllFn   func(ctx context.Context, opt option.AuthorQueryOption) (int64, []*repo.AuthorResult, error)
	findOneFn   func(ctx context.Context, id int64) (*repo.AuthorResult, error)
	findByUidFn func(ctx context.Context, uid string) (*repo.AuthorResult, error)
	createFn    func(ctx context.Context, cmd command.AuthorCreateCommand) (int64, error)
	updateFn    func(ctx context.Context, id int64, cmd command.AuthorUpdateCommand) error
}

func (s *stubAuthorRepo) FindAll(ctx context.Context, opt option.AuthorQueryOption) (int64, []*repo.AuthorResult, error) {
	return s.findAllFn(ctx, opt)
}
func (s *stubAuthorRepo) FindOne(ctx context.Context, id int64) (*repo.AuthorResult, error) {
	return s.findOneFn(ctx, id)
}
func (s *stubAuthorRepo) FindByUid(ctx context.Context, uid string) (*repo.AuthorResult, error) {
	return s.findByUidFn(ctx, uid)
}
func (s *stubAuthorRepo) Create(ctx context.Context, cmd command.AuthorCreateCommand) (int64, error) {
	return s.createFn(ctx, cmd)
}
func (s *stubAuthorRepo) Update(ctx context.Context, id int64, cmd command.AuthorUpdateCommand) error {
	return s.updateFn(ctx, id, cmd)
}

type stubPublisherRepo struct {
	findAllFn   func(ctx context.Context, opt option.PublisherQueryOption) (int64, []*repo.PublisherResult, error)
	findOneFn   func(ctx context.Context, id int64) (*repo.PublisherResult, error)
	findByUidFn func(ctx context.Context, uid string) (*repo.PublisherResult, error)
	createFn    func(ctx context.Context, cmd command.PublisherCreateCommand) (int64, error)
	updateFn    func(ctx context.Context, id int64, cmd command.PublisherUpdateCommand) error
}

func (s *stubPublisherRepo) FindAll(ctx context.Context, opt option.PublisherQueryOption) (int64, []*repo.PublisherResult, error) {
	return s.findAllFn(ctx, opt)
}
func (s *stubPublisherRepo) FindOne(ctx context.Context, id int64) (*repo.PublisherResult, error) {
	return s.findOneFn(ctx, id)
}
func (s *stubPublisherRepo) FindByUid(ctx context.Context, uid string) (*repo.PublisherResult, error) {
	return s.findByUidFn(ctx, uid)
}
func (s *stubPublisherRepo) Create(ctx context.Context, cmd command.PublisherCreateCommand) (int64, error) {
	return s.createFn(ctx, cmd)
}
func (s *stubPublisherRepo) Update(ctx context.Context, id int64, cmd command.PublisherUpdateCommand) error {
	return s.updateFn(ctx, id, cmd)
}

type stubSubjectRepo struct {
	findAllFn            func(ctx context.Context, opt option.SubjectQueryOption) (int64, []*repo.SubjectResult, error)
	findOneFn            func(ctx context.Context, id int64) (*repo.SubjectResult, error)
	findByUidFn          func(ctx context.Context, uid string) (*repo.SubjectResult, error)
	findByNameFn         func(ctx context.Context, name string) (*repo.SubjectResult, error)
	findOrCreateByNameFn func(ctx context.Context, name string) (int64, error)
	createFn             func(ctx context.Context, cmd command.SubjectCreateCommand) (int64, error)
	updateFn             func(ctx context.Context, id int64, cmd command.SubjectUpdateCommand) error
}

func (s *stubSubjectRepo) FindAll(ctx context.Context, opt option.SubjectQueryOption) (int64, []*repo.SubjectResult, error) {
	return s.findAllFn(ctx, opt)
}
func (s *stubSubjectRepo) FindOne(ctx context.Context, id int64) (*repo.SubjectResult, error) {
	return s.findOneFn(ctx, id)
}
func (s *stubSubjectRepo) FindByUid(ctx context.Context, uid string) (*repo.SubjectResult, error) {
	return s.findByUidFn(ctx, uid)
}
func (s *stubSubjectRepo) FindByName(ctx context.Context, name string) (*repo.SubjectResult, error) {
	if s.findByNameFn == nil {
		return nil, nil
	}
	return s.findByNameFn(ctx, name)
}
func (s *stubSubjectRepo) FindOrCreateByName(ctx context.Context, name string) (int64, error) {
	return s.findOrCreateByNameFn(ctx, name)
}
func (s *stubSubjectRepo) Create(ctx context.Context, cmd command.SubjectCreateCommand) (int64, error) {
	return s.createFn(ctx, cmd)
}
func (s *stubSubjectRepo) Update(ctx context.Context, id int64, cmd command.SubjectUpdateCommand) error {
	return s.updateFn(ctx, id, cmd)
}

func ptrStr(s string) *string { return &s }
func ptrInt64(v int64) *int64 { return &v }

func newService(
	work *stubWorkRepo,
	workSearch *stubWorkSearchRepo,
	edition *stubEditionRepo,
	author *stubAuthorRepo,
	publisher *stubPublisherRepo,
	subject *stubSubjectRepo,
) *catalogService {
	if work == nil {
		work = &stubWorkRepo{}
	}
	if workSearch == nil {
		workSearch = &stubWorkSearchRepo{}
	}
	if edition == nil {
		edition = &stubEditionRepo{}
	}
	if author == nil {
		author = &stubAuthorRepo{}
	}
	if publisher == nil {
		publisher = &stubPublisherRepo{}
	}
	if subject == nil {
		subject = &stubSubjectRepo{}
	}
	return &catalogService{
		workRepo:       work,
		workSearchRepo: workSearch,
		editionRepo:    edition,
		authorRepo:     author,
		publisherRepo:  publisher,
		subjectRepo:    subject,
	}
}

// ─── Work lifecycle ─────────────────────────────────────────────────────────

func TestSearchWorks_WhenRepoReturnsHits_ReturnsAggregates(t *testing.T) {
	desc := "desc"
	first := "1990-01-01"
	work := &repo.WorkSearchResult{
		Uid:              "u-1",
		Title:            "The Hobbit",
		Description:      &desc,
		FirstPublishDate: &first,
		Authors:          []repo.WorkSearchAuthor{{Uid: "a-1", Name: "Tolkien"}},
		Subjects:         []repo.WorkSearchSubject{{Uid: "s-1", Name: "Fantasy"}},
	}
	work2 := &repo.WorkSearchResult{Uid: "u-2", Title: "LOTR"}

	searchStub := &stubWorkSearchRepo{
		findAllFn: func(ctx context.Context, _ option.WorkQueryOption) (int64, []*repo.WorkSearchResult, error) {
			return 2, []*repo.WorkSearchResult{work, work2}, nil
		},
	}
	svc := newService(nil, searchStub, nil, nil, nil, nil)

	total, aggs, err := svc.SearchWorks(context.Background(), option.NewWorkQueryOption())
	if err != nil {
		t.Fatalf("SearchWorks returned error: %v", err)
	}
	if total != 2 {
		t.Errorf("total = %d, want 2", total)
	}
	if len(aggs) != 2 {
		t.Fatalf("len(aggs) = %d, want 2", len(aggs))
	}
	if aggs[0].Uid != "u-1" || aggs[0].Title != "The Hobbit" {
		t.Errorf("aggs[0] = %+v, want uid=u-1 title=The Hobbit", aggs[0])
	}
	if len(aggs[0].Authors) != 1 || aggs[0].Authors[0].Uid != "a-1" || aggs[0].Authors[0].Name != "Tolkien" {
		t.Errorf("aggs[0].Authors = %+v, want one Tolkien", aggs[0].Authors)
	}
	if len(aggs[0].Subjects) != 1 || aggs[0].Subjects[0].Name != "Fantasy" {
		t.Errorf("aggs[0].Subjects = %+v, want one Fantasy", aggs[0].Subjects)
	}
	if len(aggs[1].Authors) != 0 {
		t.Errorf("aggs[1] should have no authors, got %d", len(aggs[1].Authors))
	}
}

func TestSearchWorks_WhenRepoFails_ReturnsRepoError(t *testing.T) {
	wantErr := errors.New("boom")
	searchStub := &stubWorkSearchRepo{
		findAllFn: func(context.Context, option.WorkQueryOption) (int64, []*repo.WorkSearchResult, error) {
			return 0, nil, wantErr
		},
	}
	svc := newService(nil, searchStub, nil, nil, nil, nil)
	total, aggs, err := svc.SearchWorks(context.Background(), option.NewWorkQueryOption())
	if !errors.Is(err, wantErr) {
		t.Errorf("err = %v, want %v", err, wantErr)
	}
	if total != 0 || aggs != nil {
		t.Errorf("got total=%d aggs=%v, want zero values", total, aggs)
	}
}

func TestFindWork_WhenWorkExists_ReturnsAggregate(t *testing.T) {
	repoStub := &stubWorkRepo{
		findByUidFn: func(_ context.Context, uid string) (*repo.WorkResult, error) {
			if uid != "u-42" {
				t.Errorf("uid = %q, want u-42", uid)
			}
			return &repo.WorkResult{Id: 42, Uid: uid, Title: "Found"}, nil
		},
	}
	svc := newService(repoStub, nil, nil, nil, nil, nil)
	agg, err := svc.FindWork(context.Background(), "u-42")
	if err != nil {
		t.Fatalf("FindWork err: %v", err)
	}
	if agg.Title != "Found" {
		t.Errorf("agg.Title = %q, want Found", agg.Title)
	}
}

func TestCreateWork_WhenCommandValid_ReturnsCreatedAggregate(t *testing.T) {
	createCalled := false
	repoStub := &stubWorkRepo{
		createFn: func(_ context.Context, cmd command.WorkCreateCommand) (int64, error) {
			createCalled = true
			if cmd.Title != "New Book" {
				t.Errorf("create cmd title = %q, want New Book", cmd.Title)
			}
			return 7, nil
		},
		findOneFn: func(_ context.Context, id int64) (*repo.WorkResult, error) {
			if id != 7 {
				t.Errorf("findOne id = %d, want 7", id)
			}
			return &repo.WorkResult{Id: id, Uid: "u-new", Title: "New Book"}, nil
		},
	}
	svc := newService(repoStub, nil, nil, nil, nil, nil)
	agg, err := svc.CreateWork(context.Background(), command.WorkCreateCommand{Title: "New Book", AuthorUids: []string{cmdAuthorUid}})
	if err != nil {
		t.Fatalf("CreateWork err: %v", err)
	}
	if !createCalled {
		t.Error("expected workRepo.Create to be called")
	}
	if agg.Uid != "u-new" {
		t.Errorf("agg.Uid = %q, want u-new", agg.Uid)
	}
}

func TestCreateWork_WhenRepoCreateFails_ReturnsRepoError(t *testing.T) {
	wantErr := errors.New("create fail")
	repoStub := &stubWorkRepo{
		createFn: func(context.Context, command.WorkCreateCommand) (int64, error) { return 0, wantErr },
		findOneFn: func(context.Context, int64) (*repo.WorkResult, error) {
			t.Error("findOne should not be called when create fails")
			return nil, nil
		},
	}
	svc := newService(repoStub, nil, nil, nil, nil, nil)
	_, err := svc.CreateWork(context.Background(), command.WorkCreateCommand{Title: "x", AuthorUids: []string{cmdAuthorUid}})
	if !errors.Is(err, wantErr) {
		t.Errorf("err = %v, want %v", err, wantErr)
	}
}

func TestCreateWork_WhenCommandInvalid_ReturnsErrInvalidCommand(t *testing.T) {
	repoStub := &stubWorkRepo{
		createFn: func(context.Context, command.WorkCreateCommand) (int64, error) {
			t.Error("create should not be called for an invalid command")
			return 0, nil
		},
	}
	svc := newService(repoStub, nil, nil, nil, nil, nil)
	_, err := svc.CreateWork(context.Background(), command.WorkCreateCommand{Title: " "})
	if !errors.Is(err, command.ErrInvalidCommand) {
		t.Errorf("err = %v, want ErrInvalidCommand", err)
	}
}

func TestUpdateWork_WhenWorkExists_ReturnsUpdatedAggregate(t *testing.T) {
	updateCalled := false
	repoStub := &stubWorkRepo{
		findByUidFn: func(_ context.Context, uid string) (*repo.WorkResult, error) {
			return &repo.WorkResult{Id: 11, Uid: uid, Title: "Old"}, nil
		},
		updateFn: func(_ context.Context, id int64, cmd command.WorkUpdateCommand) error {
			updateCalled = true
			if id != 11 {
				t.Errorf("update id = %d, want 11", id)
			}
			if cmd.Title == nil || *cmd.Title != "Renamed" {
				t.Errorf("update cmd.Title = %v, want Renamed", cmd.Title)
			}
			return nil
		},
		findOneFn: func(_ context.Context, _ int64) (*repo.WorkResult, error) {
			return &repo.WorkResult{Id: 11, Uid: "u-1", Title: "Renamed"}, nil
		},
	}
	svc := newService(repoStub, nil, nil, nil, nil, nil)
	title := "Renamed"
	agg, err := svc.UpdateWork(context.Background(), "u-1", command.WorkUpdateCommand{Title: &title})
	if err != nil {
		t.Fatalf("UpdateWork err: %v", err)
	}
	if !updateCalled {
		t.Error("expected workRepo.Update to be called")
	}
	if agg.Title != "Renamed" {
		t.Errorf("agg.Title = %q, want Renamed", agg.Title)
	}
}

func TestUpdateWork_WhenWorkMissing_ReturnsLookupError(t *testing.T) {
	wantErr := errors.New("missing")
	repoStub := &stubWorkRepo{
		findByUidFn: func(context.Context, string) (*repo.WorkResult, error) { return nil, wantErr },
		updateFn: func(context.Context, int64, command.WorkUpdateCommand) error {
			t.Error("Update should not run if FindByUid fails")
			return nil
		},
	}
	svc := newService(repoStub, nil, nil, nil, nil, nil)
	_, err := svc.UpdateWork(context.Background(), "missing", command.WorkUpdateCommand{})
	if !errors.Is(err, wantErr) {
		t.Errorf("err = %v, want %v", err, wantErr)
	}
}

func TestFindWorkWithRelated_WhenWorkHasAuthors_ReturnsRelatedWorks(t *testing.T) {
	workStub := &stubWorkRepo{
		findByUidFn: func(_ context.Context, uid string) (*repo.WorkResult, error) {
			if uid != "w-1" {
				t.Errorf("workUid = %q, want w-1", uid)
			}
			return &repo.WorkResult{
				Id: 1, Uid: "w-1", Title: "Main",
				AuthorUids:  []string{"a-1"},
				AuthorNames: []string{"Author One"},
			}, nil
		},
	}
	editionStub := &stubEditionRepo{
		findAllFn: func(_ context.Context, opt option.EditionQueryOption) (int64, []*repo.EditionResult, error) {
			if opt.WorkUid() == nil || *opt.WorkUid() != "w-1" {
				t.Errorf("edition filter workUid = %v, want w-1", opt.WorkUid())
			}
			return 2, []*repo.EditionResult{
				{Id: 10, Uid: "e-10", Title: "Ed1", WorkUid: "w-1", WorkTitle: "Main"},
				{Id: 11, Uid: "e-11", Title: "Ed2", WorkUid: "w-1", WorkTitle: "Main"},
			}, nil
		},
	}
	searchStub := &stubWorkSearchRepo{
		searchFn: func(_ context.Context, opt option.WorkSearchOption) (int64, []*repo.WorkSearchResult, error) {
			if opt.AuthorUid() == nil || *opt.AuthorUid() != "a-1" {
				t.Errorf("search authorUid = %v, want a-1", opt.AuthorUid())
			}
			return 2, []*repo.WorkSearchResult{
				{Uid: "w-1", Title: "Main"},
				{Uid: "w-2", Title: "Related"},
			}, nil
		},
	}
	svc := newService(workStub, searchStub, editionStub, nil, nil, nil)
	agg, related, err := svc.FindWorkWithRelated(context.Background(), "w-1")
	if err != nil {
		t.Fatalf("FindWorkWithRelated err: %v", err)
	}
	if agg.Uid != "w-1" {
		t.Errorf("agg.Uid = %q, want w-1", agg.Uid)
	}
	if len(agg.Editions) != 2 {
		t.Errorf("editions = %d, want 2", len(agg.Editions))
	}
	if len(related) != 1 || related[0].Uid != "w-2" {
		t.Errorf("related = %+v, want only w-2 (self excluded)", related)
	}
}

func TestFindWorkWithRelated_WhenWorkHasNoAuthors_ReturnsEmptyRelatedWithoutSearch(t *testing.T) {
	workStub := &stubWorkRepo{
		findByUidFn: func(context.Context, string) (*repo.WorkResult, error) {
			return &repo.WorkResult{Id: 1, Uid: "w-1", Title: "Solo"}, nil
		},
	}
	editionStub := &stubEditionRepo{
		findAllFn: func(context.Context, option.EditionQueryOption) (int64, []*repo.EditionResult, error) {
			return 0, nil, nil
		},
	}
	searchStub := &stubWorkSearchRepo{
		searchFn: func(context.Context, option.WorkSearchOption) (int64, []*repo.WorkSearchResult, error) {
			t.Error("Search must not be called when the work has no authors")
			return 0, nil, nil
		},
	}
	svc := newService(workStub, searchStub, editionStub, nil, nil, nil)
	_, related, err := svc.FindWorkWithRelated(context.Background(), "w-1")
	if err != nil {
		t.Fatalf("FindWorkWithRelated err: %v", err)
	}
	if related == nil {
		t.Error("related = nil, want non-nil empty slice")
	}
	if len(related) != 0 {
		t.Errorf("related = %+v, want empty", related)
	}
}

func TestFindWorkWithRelated_WhenWorkMissing_ReturnsLookupError(t *testing.T) {
	wantErr := errors.New("missing work")
	workStub := &stubWorkRepo{
		findByUidFn: func(context.Context, string) (*repo.WorkResult, error) { return nil, wantErr },
	}
	searchStub := &stubWorkSearchRepo{
		searchFn: func(context.Context, option.WorkSearchOption) (int64, []*repo.WorkSearchResult, error) {
			t.Error("Search must not run when the work lookup fails")
			return 0, nil, nil
		},
	}
	svc := newService(workStub, searchStub, nil, nil, nil, nil)
	if _, _, err := svc.FindWorkWithRelated(context.Background(), "missing"); !errors.Is(err, wantErr) {
		t.Errorf("err = %v, want %v", err, wantErr)
	}
}

func TestFindWorkWithRelated_WhenSearchFails_ReturnsSearchError(t *testing.T) {
	wantErr := errors.New("opensearch down")
	workStub := &stubWorkRepo{
		findByUidFn: func(context.Context, string) (*repo.WorkResult, error) {
			return &repo.WorkResult{
				Id: 1, Uid: "w-1",
				AuthorUids:  []string{"a-1"},
				AuthorNames: []string{"A"},
			}, nil
		},
	}
	editionStub := &stubEditionRepo{
		findAllFn: func(context.Context, option.EditionQueryOption) (int64, []*repo.EditionResult, error) {
			return 0, nil, nil
		},
	}
	searchStub := &stubWorkSearchRepo{
		searchFn: func(context.Context, option.WorkSearchOption) (int64, []*repo.WorkSearchResult, error) {
			return 0, nil, wantErr
		},
	}
	svc := newService(workStub, searchStub, editionStub, nil, nil, nil)
	if _, _, err := svc.FindWorkWithRelated(context.Background(), "w-1"); !errors.Is(err, wantErr) {
		t.Errorf("err = %v, want %v", err, wantErr)
	}
}

// ─── FullText search ────────────────────────────────────────────────────────

func TestFullTextSearchWorks_WhenSearchReturnsHits_ReturnsResults(t *testing.T) {
	want := []*repo.WorkSearchResult{{Uid: "u-1", Title: "Hit", Score: 1.5}}
	searchStub := &stubWorkSearchRepo{
		searchFn: func(context.Context, option.WorkSearchOption) (int64, []*repo.WorkSearchResult, error) {
			return 1, want, nil
		},
	}
	svc := newService(nil, searchStub, nil, nil, nil, nil)
	total, got, err := svc.FullTextSearchWorks(context.Background(), option.WorkSearchOption{})
	if err != nil {
		t.Fatalf("FullTextSearchWorks err: %v", err)
	}
	if total != 1 || !reflect.DeepEqual(got, want) {
		t.Errorf("total=%d got=%v, want total=1 want=%v", total, got, want)
	}
}

// ─── Edition lifecycle ──────────────────────────────────────────────────────

func TestSearchEditions_WhenPublishersMixed_ReturnsAggregates(t *testing.T) {
	pubUid := "p-1"
	pubName := "Acme"
	editions := []*repo.EditionResult{
		{Id: 1, Uid: "e-1", Title: "First", WorkUid: "w-1", WorkTitle: "W1", PublisherUid: &pubUid, PublisherName: &pubName},
		{Id: 2, Uid: "e-2", Title: "Second", WorkUid: "w-2", WorkTitle: "W2"},
	}
	editionStub := &stubEditionRepo{
		findAllFn: func(context.Context, option.EditionQueryOption) (int64, []*repo.EditionResult, error) {
			return 2, editions, nil
		},
	}
	svc := newService(nil, nil, editionStub, nil, nil, nil)
	total, aggs, err := svc.SearchEditions(context.Background(), option.NewEditionQueryOption())
	if err != nil {
		t.Fatalf("SearchEditions err: %v", err)
	}
	if total != 2 || len(aggs) != 2 {
		t.Fatalf("total=%d len=%d, want 2/2", total, len(aggs))
	}
	if aggs[0].Publisher == nil || aggs[0].Publisher.Uid != "p-1" {
		t.Errorf("aggs[0].Publisher = %+v, want uid=p-1", aggs[0].Publisher)
	}
	if aggs[1].Publisher != nil {
		t.Errorf("aggs[1].Publisher = %+v, want nil", aggs[1].Publisher)
	}
	if aggs[0].Work == nil || aggs[0].Work.Uid != "w-1" {
		t.Errorf("aggs[0].Work = %+v, want uid=w-1", aggs[0].Work)
	}
}

func TestFindEdition_WhenEditionExists_ReturnsAggregateWithWork(t *testing.T) {
	editionStub := &stubEditionRepo{
		findByUidFn: func(_ context.Context, uid string) (*repo.EditionResult, error) {
			if uid != "e-7" {
				t.Errorf("uid = %q, want e-7", uid)
			}
			return &repo.EditionResult{Id: 7, Uid: uid, Title: "An Edition", WorkUid: "w-1", WorkTitle: "W"}, nil
		},
	}
	svc := newService(nil, nil, editionStub, nil, nil, nil)
	agg, err := svc.FindEdition(context.Background(), "e-7")
	if err != nil {
		t.Fatalf("FindEdition err: %v", err)
	}
	if agg.Uid != "e-7" || agg.Title != "An Edition" {
		t.Errorf("agg = %+v, want uid=e-7 title=An Edition", agg)
	}
	if agg.Work == nil || agg.Work.Uid != "w-1" {
		t.Errorf("agg.Work = %+v, want uid=w-1", agg.Work)
	}
}

func TestFindEdition_WhenRepoFails_ReturnsRepoError(t *testing.T) {
	wantErr := errors.New("edition gone")
	editionStub := &stubEditionRepo{
		findByUidFn: func(context.Context, string) (*repo.EditionResult, error) { return nil, wantErr },
	}
	svc := newService(nil, nil, editionStub, nil, nil, nil)
	if _, err := svc.FindEdition(context.Background(), "e-x"); !errors.Is(err, wantErr) {
		t.Errorf("err = %v, want %v", err, wantErr)
	}
}

func TestCreateEdition_WhenPublisherUidOmitted_CreatesWithNilPublisherId(t *testing.T) {
	workStub := &stubWorkRepo{
		findByUidFn: func(_ context.Context, uid string) (*repo.WorkResult, error) {
			if uid != cmdWorkUid {
				t.Errorf("workUid = %q, want %q", uid, cmdWorkUid)
			}
			return &repo.WorkResult{Id: 50, Uid: uid}, nil
		},
	}
	editionStub := &stubEditionRepo{
		createFn: func(_ context.Context, cmd command.EditionCreateCommand, workId int64, publisherId *int64) (int64, error) {
			if workId != 50 {
				t.Errorf("createFn workId = %d, want 50", workId)
			}
			if publisherId != nil {
				t.Errorf("createFn publisherId = %v, want nil", publisherId)
			}
			if cmd.Title != "1st" {
				t.Errorf("create cmd Title = %q, want 1st", cmd.Title)
			}
			return 99, nil
		},
		findOneFn: func(_ context.Context, id int64) (*repo.EditionResult, error) {
			return &repo.EditionResult{Id: id, Uid: "e-99", Title: "1st", WorkUid: "w-1", WorkTitle: "W"}, nil
		},
	}
	svc := newService(workStub, nil, editionStub, nil, nil, nil)
	agg, err := svc.CreateEdition(context.Background(), command.EditionCreateCommand{Title: "1st", WorkUid: cmdWorkUid})
	if err != nil {
		t.Fatalf("CreateEdition err: %v", err)
	}
	if agg.Uid != "e-99" {
		t.Errorf("agg.Uid = %q, want e-99", agg.Uid)
	}
}

func TestCreateEdition_WhenPublisherUidGiven_CreatesWithResolvedPublisherId(t *testing.T) {
	workStub := &stubWorkRepo{
		findByUidFn: func(context.Context, string) (*repo.WorkResult, error) {
			return &repo.WorkResult{Id: 1, Uid: "w-1"}, nil
		},
	}
	publisherStub := &stubPublisherRepo{
		findByUidFn: func(context.Context, string) (*repo.PublisherResult, error) {
			return &repo.PublisherResult{Id: 7, Uid: "p-1", Name: "Acme"}, nil
		},
	}
	var capturedPublisherId *int64
	editionStub := &stubEditionRepo{
		createFn: func(_ context.Context, _ command.EditionCreateCommand, _ int64, publisherId *int64) (int64, error) {
			capturedPublisherId = publisherId
			return 100, nil
		},
		findOneFn: func(context.Context, int64) (*repo.EditionResult, error) {
			return &repo.EditionResult{Id: 100, Uid: "e-100"}, nil
		},
	}
	svc := newService(workStub, nil, editionStub, nil, publisherStub, nil)
	pubUid := cmdPublisherUid
	_, err := svc.CreateEdition(context.Background(), command.EditionCreateCommand{Title: "x", WorkUid: cmdWorkUid, PublisherUid: &pubUid})
	if err != nil {
		t.Fatalf("CreateEdition err: %v", err)
	}
	if capturedPublisherId == nil || *capturedPublisherId != 7 {
		t.Errorf("capturedPublisherId = %v, want 7", capturedPublisherId)
	}
}

func TestCreateEdition_WhenPublisherLookupFails_ReturnsLookupError(t *testing.T) {
	wantErr := errors.New("no pub")
	workStub := &stubWorkRepo{
		findByUidFn: func(context.Context, string) (*repo.WorkResult, error) {
			return &repo.WorkResult{Id: 1, Uid: "w-1"}, nil
		},
	}
	publisherStub := &stubPublisherRepo{
		findByUidFn: func(context.Context, string) (*repo.PublisherResult, error) { return nil, wantErr },
	}
	editionStub := &stubEditionRepo{
		createFn: func(context.Context, command.EditionCreateCommand, int64, *int64) (int64, error) {
			t.Error("Create should not be called when publisher lookup fails")
			return 0, nil
		},
	}
	svc := newService(workStub, nil, editionStub, nil, publisherStub, nil)
	pub := cmdPublisherUid
	_, err := svc.CreateEdition(context.Background(), command.EditionCreateCommand{Title: "x", WorkUid: cmdWorkUid, PublisherUid: &pub})
	if !errors.Is(err, wantErr) {
		t.Errorf("err = %v, want %v", err, wantErr)
	}
}

func TestUpdateEdition_WhenWorkAndPublisherGiven_UpdatesWithResolvedIds(t *testing.T) {
	editionStub := &stubEditionRepo{
		findByUidFn: func(context.Context, string) (*repo.EditionResult, error) {
			return &repo.EditionResult{Id: 5, Uid: "e-5", WorkUid: "w-old", WorkTitle: "Old"}, nil
		},
		updateFn: func(_ context.Context, id int64, _ command.EditionUpdateCommand, workId *int64, publisherId *int64) error {
			if id != 5 {
				t.Errorf("update id = %d, want 5", id)
			}
			if workId == nil || *workId != 22 {
				t.Errorf("update workId = %v, want 22", workId)
			}
			if publisherId == nil || *publisherId != 33 {
				t.Errorf("update publisherId = %v, want 33", publisherId)
			}
			return nil
		},
		findOneFn: func(context.Context, int64) (*repo.EditionResult, error) {
			return &repo.EditionResult{Id: 5, Uid: "e-5"}, nil
		},
	}
	workStub := &stubWorkRepo{
		findByUidFn: func(context.Context, string) (*repo.WorkResult, error) {
			return &repo.WorkResult{Id: 22, Uid: "w-new"}, nil
		},
	}
	publisherStub := &stubPublisherRepo{
		findByUidFn: func(context.Context, string) (*repo.PublisherResult, error) {
			return &repo.PublisherResult{Id: 33, Uid: "p-new"}, nil
		},
	}
	svc := newService(workStub, nil, editionStub, nil, publisherStub, nil)
	w := cmdWorkUid
	p := cmdPublisherUid
	_, err := svc.UpdateEdition(context.Background(), "e-5", command.EditionUpdateCommand{WorkUid: &w, PublisherUid: &p})
	if err != nil {
		t.Fatalf("UpdateEdition err: %v", err)
	}
}

func TestUpdateEdition_WhenWorkLookupFails_ReturnsLookupError(t *testing.T) {
	wantErr := errors.New("no such work")
	editionStub := &stubEditionRepo{
		findByUidFn: func(context.Context, string) (*repo.EditionResult, error) {
			return &repo.EditionResult{Id: 5, Uid: "e-5"}, nil
		},
		updateFn: func(context.Context, int64, command.EditionUpdateCommand, *int64, *int64) error {
			t.Error("Update must not run when a parallel ref lookup fails")
			return nil
		},
	}
	workStub := &stubWorkRepo{
		findByUidFn: func(context.Context, string) (*repo.WorkResult, error) { return nil, wantErr },
	}
	svc := newService(workStub, nil, editionStub, nil, nil, nil)
	w := cmdWorkUid
	if _, err := svc.UpdateEdition(
		context.Background(), "e-5", command.EditionUpdateCommand{WorkUid: &w},
	); !errors.Is(err, wantErr) {
		t.Errorf("err = %v, want %v", err, wantErr)
	}
}

// ─── Author management ─────────────────────────────────────────────────────

func TestFindAllAuthors_WhenRepoReturnsRows_ReturnsAggregates(t *testing.T) {
	authorStub := &stubAuthorRepo{
		findAllFn: func(context.Context, option.AuthorQueryOption) (int64, []*repo.AuthorResult, error) {
			return 1, []*repo.AuthorResult{{Id: 1, Uid: "a-1", Name: "Asimov"}}, nil
		},
	}
	svc := newService(nil, nil, nil, authorStub, nil, nil)
	total, aggs, err := svc.FindAllAuthors(context.Background(), option.NewAuthorQueryOption())
	if err != nil {
		t.Fatalf("FindAllAuthors err: %v", err)
	}
	if total != 1 || len(aggs) != 1 || aggs[0].Name != "Asimov" {
		t.Errorf("got total=%d aggs=%+v", total, aggs)
	}
}

func TestFindAuthor_WhenAuthorExists_ReturnsAggregate(t *testing.T) {
	authorStub := &stubAuthorRepo{
		findByUidFn: func(_ context.Context, uid string) (*repo.AuthorResult, error) {
			if uid != "a-5" {
				t.Errorf("uid = %q, want a-5", uid)
			}
			return &repo.AuthorResult{Id: 5, Uid: uid, Name: "Herbert", Bio: ptrStr("Dune author")}, nil
		},
	}
	svc := newService(nil, nil, nil, authorStub, nil, nil)
	agg, err := svc.FindAuthor(context.Background(), "a-5")
	if err != nil {
		t.Fatalf("FindAuthor err: %v", err)
	}
	if agg.Uid != "a-5" || agg.Name != "Herbert" {
		t.Errorf("agg = %+v, want uid=a-5 name=Herbert", agg)
	}
	if agg.Bio == nil || *agg.Bio != "Dune author" {
		t.Errorf("agg.Bio = %v, want Dune author", agg.Bio)
	}
}

func TestFindAuthor_WhenRepoFails_ReturnsRepoError(t *testing.T) {
	wantErr := errors.New("author gone")
	authorStub := &stubAuthorRepo{
		findByUidFn: func(context.Context, string) (*repo.AuthorResult, error) { return nil, wantErr },
	}
	svc := newService(nil, nil, nil, authorStub, nil, nil)
	if _, err := svc.FindAuthor(context.Background(), "a-x"); !errors.Is(err, wantErr) {
		t.Errorf("err = %v, want %v", err, wantErr)
	}
}

func TestCreateAuthor_WhenCommandValid_ReturnsCreatedAggregate(t *testing.T) {
	authorStub := &stubAuthorRepo{
		createFn: func(_ context.Context, cmd command.AuthorCreateCommand) (int64, error) {
			if cmd.Name != "Le Guin" {
				t.Errorf("name = %q, want Le Guin", cmd.Name)
			}
			return 3, nil
		},
		findOneFn: func(_ context.Context, id int64) (*repo.AuthorResult, error) {
			return &repo.AuthorResult{Id: id, Uid: "a-3", Name: "Le Guin"}, nil
		},
	}
	svc := newService(nil, nil, nil, authorStub, nil, nil)
	agg, err := svc.CreateAuthor(context.Background(), command.AuthorCreateCommand{Name: "Le Guin"})
	if err != nil {
		t.Fatalf("CreateAuthor err: %v", err)
	}
	if agg.Uid != "a-3" {
		t.Errorf("agg.Uid = %q, want a-3", agg.Uid)
	}
}

func TestUpdateAuthor_WhenAuthorExists_ReturnsUpdatedAggregate(t *testing.T) {
	authorStub := &stubAuthorRepo{
		findByUidFn: func(context.Context, string) (*repo.AuthorResult, error) {
			return &repo.AuthorResult{Id: 9, Uid: "a-9", Name: "Old"}, nil
		},
		updateFn: func(_ context.Context, id int64, cmd command.AuthorUpdateCommand) error {
			if id != 9 {
				t.Errorf("id = %d, want 9", id)
			}
			if cmd.Name == nil || *cmd.Name != "New" {
				t.Errorf("cmd.Name = %v, want New", cmd.Name)
			}
			return nil
		},
		findOneFn: func(context.Context, int64) (*repo.AuthorResult, error) {
			return &repo.AuthorResult{Id: 9, Uid: "a-9", Name: "New"}, nil
		},
	}
	svc := newService(nil, nil, nil, authorStub, nil, nil)
	name := "New"
	agg, err := svc.UpdateAuthor(context.Background(), "a-9", command.AuthorUpdateCommand{Name: &name})
	if err != nil {
		t.Fatalf("UpdateAuthor err: %v", err)
	}
	if agg.Name != "New" {
		t.Errorf("agg.Name = %q, want New", agg.Name)
	}
}

// ─── Publisher management ──────────────────────────────────────────────────

// One pass over the publisher read and write paths: FindAllPublishers, FindPublisher,
// CreatePublisher and UpdatePublisher share the same repo stub and mapping.
func TestPublisherCommands_WhenRepoSucceeds_ReturnAggregates(t *testing.T) {
	store := map[string]*repo.PublisherResult{
		"p-1": {Id: 1, Uid: "p-1", Name: "Acme", Address: ptrStr("Mars")},
	}
	publisherStub := &stubPublisherRepo{
		findAllFn: func(context.Context, option.PublisherQueryOption) (int64, []*repo.PublisherResult, error) {
			out := make([]*repo.PublisherResult, 0, len(store))
			for _, p := range store {
				out = append(out, p)
			}
			return int64(len(out)), out, nil
		},
		findByUidFn: func(_ context.Context, uid string) (*repo.PublisherResult, error) {
			p, ok := store[uid]
			if !ok {
				return nil, errors.New("not found")
			}
			return p, nil
		},
		createFn: func(_ context.Context, cmd command.PublisherCreateCommand) (int64, error) {
			store["p-2"] = &repo.PublisherResult{Id: 2, Uid: "p-2", Name: cmd.Name}
			return 2, nil
		},
		findOneFn: func(_ context.Context, id int64) (*repo.PublisherResult, error) {
			for _, p := range store {
				if p.Id == id {
					return p, nil
				}
			}
			return nil, errors.New("not found")
		},
		updateFn: func(_ context.Context, id int64, cmd command.PublisherUpdateCommand) error {
			for _, p := range store {
				if p.Id == id && cmd.Name != nil {
					p.Name = *cmd.Name
				}
			}
			return nil
		},
	}
	svc := newService(nil, nil, nil, nil, publisherStub, nil)

	total, aggs, err := svc.FindAllPublishers(context.Background(), option.NewPublisherQueryOption())
	if err != nil || total != 1 || len(aggs) != 1 || aggs[0].Name != "Acme" {
		t.Fatalf("FindAllPublishers got total=%d aggs=%+v err=%v", total, aggs, err)
	}

	got, err := svc.FindPublisher(context.Background(), "p-1")
	if err != nil || got.Address == nil || *got.Address != "Mars" {
		t.Fatalf("FindPublisher got=%+v err=%v", got, err)
	}

	created, err := svc.CreatePublisher(context.Background(), command.PublisherCreateCommand{Name: "Beta"})
	if err != nil || created.Uid != "p-2" || created.Name != "Beta" {
		t.Fatalf("CreatePublisher got=%+v err=%v", created, err)
	}

	newName := "Gamma"
	upd, err := svc.UpdatePublisher(context.Background(), "p-2", command.PublisherUpdateCommand{Name: &newName})
	if err != nil || upd.Name != "Gamma" {
		t.Fatalf("UpdatePublisher got=%+v err=%v", upd, err)
	}
}

// ─── Subject management ────────────────────────────────────────────────────

// One pass over the subject read and write paths: FindAllSubjects, FindSubject,
// CreateSubject and UpdateSubject share the same repo stub and mapping.
func TestSubjectCommands_WhenRepoSucceeds_ReturnAggregates(t *testing.T) {
	subjectStub := &stubSubjectRepo{
		findAllFn: func(context.Context, option.SubjectQueryOption) (int64, []*repo.SubjectResult, error) {
			return 1, []*repo.SubjectResult{{Id: 1, Uid: "s-1", Name: "Sci-Fi"}}, nil
		},
		findByUidFn: func(_ context.Context, uid string) (*repo.SubjectResult, error) {
			return &repo.SubjectResult{Id: 1, Uid: uid, Name: "Sci-Fi"}, nil
		},
		createFn: func(_ context.Context, cmd command.SubjectCreateCommand) (int64, error) {
			if cmd.Name != "Drama" {
				t.Errorf("create name = %q, want Drama", cmd.Name)
			}
			return 7, nil
		},
		findOneFn: func(_ context.Context, id int64) (*repo.SubjectResult, error) {
			return &repo.SubjectResult{Id: id, Uid: "s-7", Name: "Drama"}, nil
		},
		updateFn: func(_ context.Context, id int64, cmd command.SubjectUpdateCommand) error {
			if id != 1 {
				t.Errorf("update id = %d, want 1", id)
			}
			if cmd.Name == nil {
				t.Error("update cmd.Name nil, want set")
			}
			return nil
		},
	}
	svc := newService(nil, nil, nil, nil, nil, subjectStub)

	total, aggs, err := svc.FindAllSubjects(context.Background(), option.NewSubjectQueryOption())
	if err != nil || total != 1 || len(aggs) != 1 || aggs[0].Name != "Sci-Fi" {
		t.Fatalf("FindAllSubjects got total=%d aggs=%+v err=%v", total, aggs, err)
	}

	one, err := svc.FindSubject(context.Background(), "s-1")
	if err != nil || one.Name != "Sci-Fi" {
		t.Fatalf("FindSubject got=%+v err=%v", one, err)
	}

	created, err := svc.CreateSubject(context.Background(), command.SubjectCreateCommand{Name: "Drama"})
	if err != nil || created.Uid != "s-7" {
		t.Fatalf("CreateSubject got=%+v err=%v", created, err)
	}

	newName := "Comedy"
	upd, err := svc.UpdateSubject(context.Background(), "s-1", command.SubjectUpdateCommand{Name: &newName})
	if err != nil || upd == nil {
		t.Fatalf("UpdateSubject got=%+v err=%v", upd, err)
	}
}

func TestCreateSubject_WhenNameTaken_ReturnsErrAlreadyExists(t *testing.T) {
	subjectStub := &stubSubjectRepo{
		findByNameFn: func(_ context.Context, name string) (*repo.SubjectResult, error) {
			return &repo.SubjectResult{Id: 3, Uid: "s-3", Name: name}, nil
		},
		createFn: func(context.Context, command.SubjectCreateCommand) (int64, error) {
			t.Error("Create must not run when the name is taken")
			return 0, nil
		},
	}
	svc := newService(nil, nil, nil, nil, nil, subjectStub)
	_, err := svc.CreateSubject(context.Background(), command.SubjectCreateCommand{Name: "Sci-Fi"})
	if !errors.Is(err, repo.ErrAlreadyExists) {
		t.Errorf("err = %v, want ErrAlreadyExists", err)
	}
}

func TestUpdateSubject_WhenNameTakenByAnother_ReturnsErrAlreadyExists(t *testing.T) {
	subjectStub := &stubSubjectRepo{
		findByUidFn: func(_ context.Context, uid string) (*repo.SubjectResult, error) {
			return &repo.SubjectResult{Id: 1, Uid: uid, Name: "Sci-Fi"}, nil
		},
		findByNameFn: func(_ context.Context, name string) (*repo.SubjectResult, error) {
			return &repo.SubjectResult{Id: 9, Uid: "s-9", Name: name}, nil
		},
		updateFn: func(context.Context, int64, command.SubjectUpdateCommand) error {
			t.Error("Update must not run when the name is taken by another subject")
			return nil
		},
	}
	svc := newService(nil, nil, nil, nil, nil, subjectStub)
	name := "Drama"
	_, err := svc.UpdateSubject(context.Background(), "s-1", command.SubjectUpdateCommand{Name: &name})
	if !errors.Is(err, repo.ErrAlreadyExists) {
		t.Errorf("err = %v, want ErrAlreadyExists", err)
	}
}

func TestUpdateSubject_WhenNameBelongsToSameSubject_CallsUpdate(t *testing.T) {
	updated := false
	subjectStub := &stubSubjectRepo{
		findByUidFn: func(_ context.Context, uid string) (*repo.SubjectResult, error) {
			return &repo.SubjectResult{Id: 1, Uid: uid, Name: "Sci-Fi"}, nil
		},
		findByNameFn: func(_ context.Context, name string) (*repo.SubjectResult, error) {
			return &repo.SubjectResult{Id: 1, Uid: "s-1", Name: name}, nil
		},
		updateFn: func(context.Context, int64, command.SubjectUpdateCommand) error {
			updated = true
			return nil
		},
		findOneFn: func(_ context.Context, id int64) (*repo.SubjectResult, error) {
			return &repo.SubjectResult{Id: id, Uid: "s-1", Name: "Sci-Fi"}, nil
		},
	}
	svc := newService(nil, nil, nil, nil, nil, subjectStub)
	name := "Sci-Fi"
	if _, err := svc.UpdateSubject(context.Background(), "s-1", command.SubjectUpdateCommand{Name: &name}); err != nil {
		t.Fatalf("UpdateSubject err: %v", err)
	}
	if !updated {
		t.Error("Update must run when the name belongs to the same subject")
	}
}
