package restcontroller

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	"gorm.io/gorm"

	"github.com/codejsha/bookstore-microservices/catalog/generated/application/port/openapi"
	"github.com/codejsha/bookstore-microservices/catalog/internal/application/port/repo"
	"github.com/codejsha/bookstore-microservices/catalog/internal/application/usecase"
	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/aggregate"
	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/model/command"
	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/model/option"
)

type stubUseCase struct {
	searchWorks         func(context.Context, option.WorkQueryOption) (int64, []*aggregate.WorkAggregate, error)
	fullTextSearchWorks func(context.Context, option.WorkSearchOption) (int64, []*repo.WorkSearchResult, error)
	findWork            func(context.Context, string) (*aggregate.WorkAggregate, error)
	createWork          func(context.Context, command.WorkCreateCommand) (*aggregate.WorkAggregate, error)
	updateWork          func(context.Context, string, command.WorkUpdateCommand) (*aggregate.WorkAggregate, error)
	searchEditions      func(context.Context, option.EditionQueryOption) (int64, []*aggregate.EditionAggregate, error)
	findEdition         func(context.Context, string) (*aggregate.EditionAggregate, error)
	createEdition       func(context.Context, command.EditionCreateCommand) (*aggregate.EditionAggregate, error)
	updateEdition       func(context.Context, string, command.EditionUpdateCommand) (*aggregate.EditionAggregate, error)
	findAllAuthors      func(context.Context, option.AuthorQueryOption) (int64, []*aggregate.AuthorAggregate, error)
	findAuthor          func(context.Context, string) (*aggregate.AuthorAggregate, error)
	createAuthor        func(context.Context, command.AuthorCreateCommand) (*aggregate.AuthorAggregate, error)
	updateAuthor        func(context.Context, string, command.AuthorUpdateCommand) (*aggregate.AuthorAggregate, error)
	findAllPublishers   func(context.Context, option.PublisherQueryOption) (int64, []*aggregate.PublisherAggregate, error)
	findPublisher       func(context.Context, string) (*aggregate.PublisherAggregate, error)
	createPublisher     func(context.Context, command.PublisherCreateCommand) (*aggregate.PublisherAggregate, error)
	updatePublisher     func(context.Context, string, command.PublisherUpdateCommand) (*aggregate.PublisherAggregate, error)
	findAllSubjects     func(context.Context, option.SubjectQueryOption) (int64, []*aggregate.SubjectAggregate, error)
	findSubject         func(context.Context, string) (*aggregate.SubjectAggregate, error)
	createSubject       func(context.Context, command.SubjectCreateCommand) (*aggregate.SubjectAggregate, error)
	updateSubject       func(context.Context, string, command.SubjectUpdateCommand) (*aggregate.SubjectAggregate, error)
}

var _ usecase.CatalogUseCase = (*stubUseCase)(nil)

func (s *stubUseCase) SearchWorks(ctx context.Context, opt option.WorkQueryOption) (int64, []*aggregate.WorkAggregate, error) {
	return s.searchWorks(ctx, opt)
}
func (s *stubUseCase) FullTextSearchWorks(ctx context.Context, opt option.WorkSearchOption) (int64, []*repo.WorkSearchResult, error) {
	return s.fullTextSearchWorks(ctx, opt)
}
func (s *stubUseCase) FindWork(ctx context.Context, uid string) (*aggregate.WorkAggregate, error) {
	return s.findWork(ctx, uid)
}
func (s *stubUseCase) FindWorkWithRelated(context.Context, string) (*aggregate.WorkAggregate, []*repo.WorkSearchResult, error) {
	return nil, nil, nil
}
func (s *stubUseCase) CreateWork(ctx context.Context, cmd command.WorkCreateCommand) (*aggregate.WorkAggregate, error) {
	return s.createWork(ctx, cmd)
}
func (s *stubUseCase) UpdateWork(ctx context.Context, uid string, cmd command.WorkUpdateCommand) (*aggregate.WorkAggregate, error) {
	return s.updateWork(ctx, uid, cmd)
}
func (s *stubUseCase) SearchEditions(ctx context.Context, opt option.EditionQueryOption) (int64, []*aggregate.EditionAggregate, error) {
	return s.searchEditions(ctx, opt)
}
func (s *stubUseCase) FindEdition(ctx context.Context, uid string) (*aggregate.EditionAggregate, error) {
	return s.findEdition(ctx, uid)
}
func (s *stubUseCase) CreateEdition(ctx context.Context, cmd command.EditionCreateCommand) (*aggregate.EditionAggregate, error) {
	return s.createEdition(ctx, cmd)
}
func (s *stubUseCase) UpdateEdition(ctx context.Context, uid string, cmd command.EditionUpdateCommand) (*aggregate.EditionAggregate, error) {
	return s.updateEdition(ctx, uid, cmd)
}
func (s *stubUseCase) FindAllAuthors(ctx context.Context, opt option.AuthorQueryOption) (int64, []*aggregate.AuthorAggregate, error) {
	return s.findAllAuthors(ctx, opt)
}
func (s *stubUseCase) FindAuthor(ctx context.Context, uid string) (*aggregate.AuthorAggregate, error) {
	return s.findAuthor(ctx, uid)
}
func (s *stubUseCase) CreateAuthor(ctx context.Context, cmd command.AuthorCreateCommand) (*aggregate.AuthorAggregate, error) {
	return s.createAuthor(ctx, cmd)
}
func (s *stubUseCase) UpdateAuthor(ctx context.Context, uid string, cmd command.AuthorUpdateCommand) (*aggregate.AuthorAggregate, error) {
	return s.updateAuthor(ctx, uid, cmd)
}
func (s *stubUseCase) FindAllPublishers(ctx context.Context, opt option.PublisherQueryOption) (int64, []*aggregate.PublisherAggregate, error) {
	return s.findAllPublishers(ctx, opt)
}
func (s *stubUseCase) FindPublisher(ctx context.Context, uid string) (*aggregate.PublisherAggregate, error) {
	return s.findPublisher(ctx, uid)
}
func (s *stubUseCase) CreatePublisher(ctx context.Context, cmd command.PublisherCreateCommand) (*aggregate.PublisherAggregate, error) {
	return s.createPublisher(ctx, cmd)
}
func (s *stubUseCase) UpdatePublisher(ctx context.Context, uid string, cmd command.PublisherUpdateCommand) (*aggregate.PublisherAggregate, error) {
	return s.updatePublisher(ctx, uid, cmd)
}
func (s *stubUseCase) FindAllSubjects(ctx context.Context, opt option.SubjectQueryOption) (int64, []*aggregate.SubjectAggregate, error) {
	return s.findAllSubjects(ctx, opt)
}
func (s *stubUseCase) FindSubject(ctx context.Context, uid string) (*aggregate.SubjectAggregate, error) {
	return s.findSubject(ctx, uid)
}
func (s *stubUseCase) CreateSubject(ctx context.Context, cmd command.SubjectCreateCommand) (*aggregate.SubjectAggregate, error) {
	return s.createSubject(ctx, cmd)
}
func (s *stubUseCase) UpdateSubject(ctx context.Context, uid string, cmd command.SubjectUpdateCommand) (*aggregate.SubjectAggregate, error) {
	return s.updateSubject(ctx, uid, cmd)
}

func ptrStr(s string) *string { return &s }

// ─── Work controller ───────────────────────────────────────────────────────

func TestWorkController_WhenSearchFiltersGiven_ReturnsPagedWorks(t *testing.T) {
	use := &stubUseCase{
		searchWorks: func(_ context.Context, opt option.WorkQueryOption) (int64, []*aggregate.WorkAggregate, error) {
			if opt.Title() == nil || *opt.Title() != "hobbit" {
				t.Errorf("title = %v, want hobbit", opt.Title())
			}
			return 1, []*aggregate.WorkAggregate{{Uid: "w-1", Title: "Hobbit"}}, nil
		},
	}
	ctrl := NewWorkController(use)
	title := "hobbit"
	res, err := ctrl.WorksSearch(context.Background(), openapi.WorkApiCredentials{}, &title, nil, nil, nil, nil, nil, nil)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	got, ok := res.(openapi.WorksSearchOk)
	if !ok {
		t.Fatalf("response = %T, want openapi.WorksSearchOk", res)
	}
	if got.Body.Total != 1 || len(got.Body.Items) != 1 || got.Body.Items[0].Uid != "w-1" {
		t.Errorf("resp = %+v", got.Body)
	}
}

func TestWorkController_WhenSearchUsecaseFails_PropagatesError(t *testing.T) {
	wantErr := errors.New("usecase down")
	use := &stubUseCase{
		searchWorks: func(context.Context, option.WorkQueryOption) (int64, []*aggregate.WorkAggregate, error) {
			return 0, nil, wantErr
		},
	}
	ctrl := NewWorkController(use)
	_, err := ctrl.WorksSearch(context.Background(), openapi.WorkApiCredentials{}, nil, nil, nil, nil, nil, nil, nil)
	if !errors.Is(err, wantErr) {
		t.Errorf("err = %v, want %v", err, wantErr)
	}
}

func TestWorkController_WhenWorkExists_ReturnsWorkResponse(t *testing.T) {
	use := &stubUseCase{
		findWork: func(_ context.Context, uid string) (*aggregate.WorkAggregate, error) {
			if uid != "w-1" {
				t.Errorf("uid = %q, want w-1", uid)
			}
			return &aggregate.WorkAggregate{
				Uid:       "w-1",
				Title:     "T",
				CoverUids: []string{"c1"},
				Authors:   []*aggregate.AuthorAggregate{{Uid: "a-1", Name: "AA"}},
				Subjects:  []*aggregate.SubjectAggregate{{Uid: "s-1", Name: "SS"}},
			}, nil
		},
	}
	ctrl := NewWorkController(use)
	res, err := ctrl.WorksRead(context.Background(), openapi.WorkApiCredentials{}, "w-1")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	got, ok := res.(openapi.WorksReadOk)
	if !ok {
		t.Fatalf("response = %T, want openapi.WorksReadOk", res)
	}
	if got.Body.Uid != "w-1" || got.Body.Title != "T" {
		t.Errorf("resp = %+v", got.Body)
	}
	if got.Body.CoverUids == nil || len(*got.Body.CoverUids) != 1 {
		t.Errorf("CoverUids = %+v, want [c1]", got.Body.CoverUids)
	}
	if len(got.Body.Authors) != 1 || got.Body.Authors[0].Uid != "a-1" {
		t.Errorf("Authors = %+v", got.Body.Authors)
	}
	if len(got.Body.Subjects) != 1 || got.Body.Subjects[0].Name != "SS" {
		t.Errorf("Subjects = %+v", got.Body.Subjects)
	}
}

func TestWorkController_WhenCreateRequestValid_PassesCommandToUsecase(t *testing.T) {
	captured := command.WorkCreateCommand{}
	use := &stubUseCase{
		createWork: func(_ context.Context, cmd command.WorkCreateCommand) (*aggregate.WorkAggregate, error) {
			captured = cmd
			return &aggregate.WorkAggregate{Uid: "w-new"}, nil
		},
	}
	ctrl := NewWorkController(use)
	covers := []string{"c1", "c2"}
	subjects := []string{"sci-fi"}
	desc := "blurb"
	req := openapi.WorkCreateRequest{
		Title:        "New",
		Description:  &desc,
		CoverUids:    &covers,
		AuthorUids:   []string{"a-1"},
		SubjectNames: &subjects,
	}
	res, err := ctrl.WorksCreate(context.Background(), openapi.WorkApiCredentials{}, req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	got, ok := res.(openapi.WorksCreateCreated)
	if !ok {
		t.Fatalf("response = %T, want openapi.WorksCreateCreated", res)
	}
	if !strings.HasSuffix(got.Location, "/w-new") {
		t.Errorf("Location = %q, want suffix /w-new", got.Location)
	}
	if captured.Title != "New" || !reflect.DeepEqual(captured.CoverUids, covers) {
		t.Errorf("captured = %+v", captured)
	}
}

func TestWorkController_WhenUpdateRequestValid_ReturnsUpdatedWork(t *testing.T) {
	use := &stubUseCase{
		updateWork: func(_ context.Context, uid string, cmd command.WorkUpdateCommand) (*aggregate.WorkAggregate, error) {
			if uid != "w-9" {
				t.Errorf("uid = %q, want w-9", uid)
			}
			if cmd.Title == nil || *cmd.Title != "Renamed" {
				t.Errorf("cmd.Title = %v", cmd.Title)
			}
			return &aggregate.WorkAggregate{Uid: uid, Title: *cmd.Title}, nil
		},
	}
	ctrl := NewWorkController(use)
	title := "Renamed"
	res, err := ctrl.WorksUpdate(context.Background(), openapi.WorkApiCredentials{}, "w-9", openapi.WorkUpdateRequest{Title: &title})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	got, ok := res.(openapi.WorksUpdateOk)
	if !ok {
		t.Fatalf("response = %T, want openapi.WorksUpdateOk", res)
	}
	if got.Body.Title != "Renamed" {
		t.Errorf("Title = %q", got.Body.Title)
	}
}

func TestWorkController_WhenCreateHitsMissingReference_PreservesNotFoundError(t *testing.T) {
	use := &stubUseCase{
		createWork: func(context.Context, command.WorkCreateCommand) (*aggregate.WorkAggregate, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}
	ctrl := NewWorkController(use)
	_, err := ctrl.WorksCreate(context.Background(), openapi.WorkApiCredentials{}, openapi.WorkCreateRequest{Title: "x", AuthorUids: []string{"missing"}})
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Errorf("err = %v, want gorm.ErrRecordNotFound preserved", err)
	}
}

func TestEditionController_WhenCreateHitsMissingReference_PreservesNotFoundError(t *testing.T) {
	use := &stubUseCase{
		createEdition: func(context.Context, command.EditionCreateCommand) (*aggregate.EditionAggregate, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}
	ctrl := NewEditionController(use)
	_, err := ctrl.EditionsCreate(context.Background(), openapi.EditionApiCredentials{}, openapi.EditionCreateRequest{Title: "x", WorkUid: "missing"})
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Errorf("err = %v, want gorm.ErrRecordNotFound preserved", err)
	}
}

// ─── Author controller ────────────────────────────────────────────────────

func TestAuthorController_WhenAuthorExists_ReturnsAuthorResponse(t *testing.T) {
	use := &stubUseCase{
		findAuthor: func(context.Context, string) (*aggregate.AuthorAggregate, error) {
			return &aggregate.AuthorAggregate{Uid: "a-1", Name: "Asimov", AlternateNames: []string{"Paul French"}}, nil
		},
	}
	ctrl := NewAuthorController(use)
	res, err := ctrl.AuthorsRead(context.Background(), openapi.AuthorApiCredentials{}, "a-1")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	got, ok := res.(openapi.AuthorsReadOk)
	if !ok {
		t.Fatalf("response = %T, want openapi.AuthorsReadOk", res)
	}
	if got.Body.Name != "Asimov" {
		t.Errorf("Name = %q", got.Body.Name)
	}
	if got.Body.AlternateNames == nil || (*got.Body.AlternateNames)[0] != "Paul French" {
		t.Errorf("AlternateNames = %+v", got.Body.AlternateNames)
	}
}

func TestAuthorController_WhenCreateRequestValid_PassesCommandToUsecase(t *testing.T) {
	use := &stubUseCase{
		createAuthor: func(_ context.Context, cmd command.AuthorCreateCommand) (*aggregate.AuthorAggregate, error) {
			if cmd.Name != "Le Guin" {
				t.Errorf("cmd.Name = %q", cmd.Name)
			}
			return &aggregate.AuthorAggregate{Uid: "a-new", Name: cmd.Name}, nil
		},
	}
	ctrl := NewAuthorController(use)
	res, err := ctrl.AuthorsCreate(context.Background(), openapi.AuthorApiCredentials{}, openapi.AuthorCreateRequest{Name: "Le Guin"})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	got, ok := res.(openapi.AuthorsCreateCreated)
	if !ok {
		t.Fatalf("response = %T, want openapi.AuthorsCreateCreated", res)
	}
	if !strings.HasSuffix(got.Location, "/a-new") {
		t.Errorf("Location = %q, want suffix /a-new", got.Location)
	}
}

// ─── Publisher controller ─────────────────────────────────────────────────

func TestPublisherController_WhenNameFilterGiven_ReturnsPagedPublishers(t *testing.T) {
	use := &stubUseCase{
		findAllPublishers: func(_ context.Context, opt option.PublisherQueryOption) (int64, []*aggregate.PublisherAggregate, error) {
			if opt.Name() == nil || *opt.Name() != "Acme" {
				t.Errorf("name = %v", opt.Name())
			}
			return 1, []*aggregate.PublisherAggregate{{Uid: "p-1", Name: "Acme", Address: ptrStr("Mars")}}, nil
		},
	}
	ctrl := NewPublisherController(use)
	name := "Acme"
	res, err := ctrl.PublishersGetAll(context.Background(), openapi.PublisherApiCredentials{}, &name, nil, nil, nil)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	got, ok := res.(openapi.PublishersGetAllOk)
	if !ok {
		t.Fatalf("response = %T, want openapi.PublishersGetAllOk", res)
	}
	if got.Body.Total != 1 || len(got.Body.Items) != 1 {
		t.Fatalf("resp = %+v", got.Body)
	}
	if got.Body.Items[0].Address == nil || *got.Body.Items[0].Address != "Mars" {
		t.Errorf("Address = %v", got.Body.Items[0].Address)
	}
}

// ─── Subject controller ───────────────────────────────────────────────────

func TestSubjectController_WhenCreateRequestValid_PassesCommandToUsecase(t *testing.T) {
	use := &stubUseCase{
		createSubject: func(_ context.Context, cmd command.SubjectCreateCommand) (*aggregate.SubjectAggregate, error) {
			if cmd.Name != "Drama" {
				t.Errorf("cmd.Name = %q", cmd.Name)
			}
			return &aggregate.SubjectAggregate{Uid: "s-new", Name: cmd.Name}, nil
		},
	}
	ctrl := NewSubjectController(use)
	res, err := ctrl.SubjectsCreate(context.Background(), openapi.SubjectApiCredentials{}, openapi.SubjectCreateRequest{Name: "Drama"})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	got, ok := res.(openapi.SubjectsCreateCreated)
	if !ok {
		t.Fatalf("response = %T, want openapi.SubjectsCreateCreated", res)
	}
	if !strings.HasSuffix(got.Location, "/s-new") {
		t.Errorf("Location = %q, want suffix /s-new", got.Location)
	}
}

func TestSubjectController_WhenSubjectExists_ReturnsSubjectResponse(t *testing.T) {
	use := &stubUseCase{
		findSubject: func(context.Context, string) (*aggregate.SubjectAggregate, error) {
			return &aggregate.SubjectAggregate{Uid: "s-1", Name: "Sci-Fi"}, nil
		},
	}
	ctrl := NewSubjectController(use)
	res, err := ctrl.SubjectsRead(context.Background(), openapi.SubjectApiCredentials{}, "s-1")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	got, ok := res.(openapi.SubjectsReadOk)
	if !ok {
		t.Fatalf("response = %T, want openapi.SubjectsReadOk", res)
	}
	if got.Body.Name != "Sci-Fi" {
		t.Errorf("Name = %q", got.Body.Name)
	}
}

// ─── Edition controller ───────────────────────────────────────────────────

func TestEditionController_WhenEditionExists_ReturnsEditionResponse(t *testing.T) {
	use := &stubUseCase{
		findEdition: func(context.Context, string) (*aggregate.EditionAggregate, error) {
			return &aggregate.EditionAggregate{
				Uid:       "e-1",
				Title:     "Edition",
				Work:      &aggregate.WorkAggregate{Uid: "w-1", Title: "W"},
				Publisher: &aggregate.PublisherAggregate{Uid: "p-1", Name: "Acme"},
			}, nil
		},
	}
	ctrl := NewEditionController(use)
	res, err := ctrl.EditionsRead(context.Background(), openapi.EditionApiCredentials{}, "e-1")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	got, ok := res.(openapi.EditionsReadOk)
	if !ok {
		t.Fatalf("response = %T, want openapi.EditionsReadOk", res)
	}
	if got.Body.Work != (openapi.EditionWorkRef{Uid: "w-1", Title: "W"}) || got.Body.Publisher == nil || got.Body.Publisher.Uid != "p-1" {
		t.Errorf("resp = %+v", got.Body)
	}
}

func TestEditionController_WhenEditionHasNoPublisher_OmitsPublisherFromResponse(t *testing.T) {
	use := &stubUseCase{
		findEdition: func(context.Context, string) (*aggregate.EditionAggregate, error) {
			return &aggregate.EditionAggregate{Uid: "e-1", Work: &aggregate.WorkAggregate{Uid: "w-1"}}, nil
		},
	}
	ctrl := NewEditionController(use)
	res, err := ctrl.EditionsRead(context.Background(), openapi.EditionApiCredentials{}, "e-1")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	got, ok := res.(openapi.EditionsReadOk)
	if !ok {
		t.Fatalf("response = %T, want openapi.EditionsReadOk", res)
	}
	if got.Body.Publisher != nil {
		t.Errorf("Publisher = %+v, want nil", got.Body.Publisher)
	}
}

// ─── slice helpers ─────────────────────────────────────────────────────────

// The non-empty case is asserted too: it must come back as a pointer to the same slice.
func TestPtrStringSlice_WhenSliceNilOrEmpty_ReturnsNil(t *testing.T) {
	if got := ptrStringSlice(nil); got != nil {
		t.Errorf("nil -> %v, want nil", got)
	}
	if got := ptrStringSlice([]string{}); got != nil {
		t.Errorf("empty -> %v, want nil", got)
	}
	if got := ptrStringSlice([]string{"a"}); got == nil || (*got)[0] != "a" {
		t.Errorf("[a] -> %v, want pointer", got)
	}
}

// The non-nil case is asserted too: it must come back as the pointed-to slice.
func TestDerefStringSlice_WhenPointerNil_ReturnsNil(t *testing.T) {
	if got := derefStringSlice(nil); got != nil {
		t.Errorf("nil -> %v, want nil", got)
	}
	in := []string{"a", "b"}
	if got := derefStringSlice(&in); !reflect.DeepEqual(got, in) {
		t.Errorf("got %v, want %v", got, in)
	}
}
