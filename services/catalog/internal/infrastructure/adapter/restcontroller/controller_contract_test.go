package restcontroller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/codejsha/bookstore-microservices/catalog/generated/application/port/openapi"
	"github.com/codejsha/bookstore-microservices/catalog/internal/application/port/repo"
	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/aggregate"
	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/model/command"
	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/model/option"
	"github.com/codejsha/bookstore-microservices/catalog/internal/infrastructure/httpx"
)

type problemDetails struct {
	Title  string `json:"title"`
	Status int    `json:"status"`
	Detail string `json:"detail"`
}

func serveThroughErrorMapping(t *testing.T, call func(ctx context.Context) error) (*httptest.ResponseRecorder, problemDetails) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(httpx.GinResponseMapping())
	r.GET("/probe", func(c *gin.Context) {
		if err := call(c.Request.Context()); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{})
	})
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/probe", nil))
	var problem problemDetails
	if rec.Code >= http.StatusBadRequest {
		if err := json.Unmarshal(rec.Body.Bytes(), &problem); err != nil {
			t.Fatalf("problem body %q: %v", rec.Body.String(), err)
		}
	}
	return rec, problem
}

func failingUseCase(err error) *stubUseCase {
	return &stubUseCase{
		findWork: func(context.Context, string) (*aggregate.WorkAggregate, error) { return nil, err },
		createWork: func(context.Context, command.WorkCreateCommand) (*aggregate.WorkAggregate, error) {
			return nil, err
		},
		updateWork: func(context.Context, string, command.WorkUpdateCommand) (*aggregate.WorkAggregate, error) {
			return nil, err
		},
		findEdition: func(context.Context, string) (*aggregate.EditionAggregate, error) { return nil, err },
		createEdition: func(context.Context, command.EditionCreateCommand) (*aggregate.EditionAggregate, error) {
			return nil, err
		},
		updateEdition: func(context.Context, string, command.EditionUpdateCommand) (*aggregate.EditionAggregate, error) {
			return nil, err
		},
		findAuthor: func(context.Context, string) (*aggregate.AuthorAggregate, error) { return nil, err },
		createAuthor: func(context.Context, command.AuthorCreateCommand) (*aggregate.AuthorAggregate, error) {
			return nil, err
		},
		updateAuthor: func(context.Context, string, command.AuthorUpdateCommand) (*aggregate.AuthorAggregate, error) {
			return nil, err
		},
		findPublisher: func(context.Context, string) (*aggregate.PublisherAggregate, error) { return nil, err },
		createPublisher: func(context.Context, command.PublisherCreateCommand) (*aggregate.PublisherAggregate, error) {
			return nil, err
		},
		updatePublisher: func(context.Context, string, command.PublisherUpdateCommand) (*aggregate.PublisherAggregate, error) {
			return nil, err
		},
		findSubject: func(context.Context, string) (*aggregate.SubjectAggregate, error) { return nil, err },
		createSubject: func(context.Context, command.SubjectCreateCommand) (*aggregate.SubjectAggregate, error) {
			return nil, err
		},
		updateSubject: func(context.Context, string, command.SubjectUpdateCommand) (*aggregate.SubjectAggregate, error) {
			return nil, err
		},
	}
}

type controllerCall struct {
	name string
	run  func(ctx context.Context, use *stubUseCase) error
}

var readCalls = []controllerCall{
	{"WorksRead", func(ctx context.Context, use *stubUseCase) error {
		_, err := NewWorkController(use).WorksRead(ctx, openapi.WorkApiCredentials{}, "w-1")
		return err
	}},
	{"EditionsRead", func(ctx context.Context, use *stubUseCase) error {
		_, err := NewEditionController(use).EditionsRead(ctx, openapi.EditionApiCredentials{}, "e-1")
		return err
	}},
	{"AuthorsRead", func(ctx context.Context, use *stubUseCase) error {
		_, err := NewAuthorController(use).AuthorsRead(ctx, openapi.AuthorApiCredentials{}, "a-1")
		return err
	}},
	{"PublishersRead", func(ctx context.Context, use *stubUseCase) error {
		_, err := NewPublisherController(use).PublishersRead(ctx, openapi.PublisherApiCredentials{}, "p-1")
		return err
	}},
	{"SubjectsRead", func(ctx context.Context, use *stubUseCase) error {
		_, err := NewSubjectController(use).SubjectsRead(ctx, openapi.SubjectApiCredentials{}, "s-1")
		return err
	}},
}

var createCalls = []controllerCall{
	{"WorksCreate", func(ctx context.Context, use *stubUseCase) error {
		_, err := NewWorkController(use).WorksCreate(ctx, openapi.WorkApiCredentials{}, openapi.WorkCreateRequest{Title: "x"})
		return err
	}},
	{"EditionsCreate", func(ctx context.Context, use *stubUseCase) error {
		_, err := NewEditionController(use).EditionsCreate(ctx, openapi.EditionApiCredentials{}, openapi.EditionCreateRequest{Title: "x"})
		return err
	}},
	{"AuthorsCreate", func(ctx context.Context, use *stubUseCase) error {
		_, err := NewAuthorController(use).AuthorsCreate(ctx, openapi.AuthorApiCredentials{}, openapi.AuthorCreateRequest{Name: "x"})
		return err
	}},
	{"PublishersCreate", func(ctx context.Context, use *stubUseCase) error {
		_, err := NewPublisherController(use).PublishersCreate(ctx, openapi.PublisherApiCredentials{}, openapi.PublisherCreateRequest{Name: "x"})
		return err
	}},
	{"SubjectsCreate", func(ctx context.Context, use *stubUseCase) error {
		_, err := NewSubjectController(use).SubjectsCreate(ctx, openapi.SubjectApiCredentials{}, openapi.SubjectCreateRequest{Name: "x"})
		return err
	}},
}

var updateCalls = []controllerCall{
	{"WorksUpdate", func(ctx context.Context, use *stubUseCase) error {
		_, err := NewWorkController(use).WorksUpdate(ctx, openapi.WorkApiCredentials{}, "w-1", openapi.WorkUpdateRequest{})
		return err
	}},
	{"EditionsUpdate", func(ctx context.Context, use *stubUseCase) error {
		_, err := NewEditionController(use).EditionsUpdate(ctx, openapi.EditionApiCredentials{}, "e-1", openapi.EditionUpdateRequest{})
		return err
	}},
	{"AuthorsUpdate", func(ctx context.Context, use *stubUseCase) error {
		_, err := NewAuthorController(use).AuthorsUpdate(ctx, openapi.AuthorApiCredentials{}, "a-1", openapi.AuthorUpdateRequest{})
		return err
	}},
	{"PublishersUpdate", func(ctx context.Context, use *stubUseCase) error {
		_, err := NewPublisherController(use).PublishersUpdate(ctx, openapi.PublisherApiCredentials{}, "p-1", openapi.PublisherUpdateRequest{})
		return err
	}},
	{"SubjectsUpdate", func(ctx context.Context, use *stubUseCase) error {
		_, err := NewSubjectController(use).SubjectsUpdate(ctx, openapi.SubjectApiCredentials{}, "s-1", openapi.SubjectUpdateRequest{})
		return err
	}},
}

func joinCalls(groups ...[]controllerCall) []controllerCall {
	var out []controllerCall
	for _, g := range groups {
		out = append(out, g...)
	}
	return out
}

func assertProblem(t *testing.T, rec *httptest.ResponseRecorder, got problemDetails, status int, detail string) {
	t.Helper()
	if rec.Code != status {
		t.Fatalf("status = %d, want %d (body %s)", rec.Code, status, rec.Body.String())
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/problem+json" {
		t.Errorf("content-type = %q, want application/problem+json", ct)
	}
	want := problemDetails{Title: http.StatusText(status), Status: status, Detail: detail}
	if got != want {
		t.Errorf("problem = %+v, want %+v", got, want)
	}
}

func TestControllers_resourceMissing_mapsToNotFoundProblem(t *testing.T) {
	captureSideEffects(t)
	missing := fmt.Errorf("load: %w", gorm.ErrRecordNotFound)
	for _, call := range joinCalls(readCalls, updateCalls, createCalls[:2]) {
		t.Run(call.name, func(t *testing.T) {
			rec, problem := serveThroughErrorMapping(t, func(ctx context.Context) error {
				return call.run(ctx, failingUseCase(missing))
			})
			assertProblem(t, rec, problem, http.StatusNotFound, "resource not found")
		})
	}
}

func TestControllers_commandInvalid_mapsToBadRequestProblemWithReason(t *testing.T) {
	captureSideEffects(t)
	invalid := command.WorkCreateCommand{}.Validate()
	if !errors.Is(invalid, command.ErrInvalidCommand) {
		t.Fatalf("precondition: %v is not an invalid-command error", invalid)
	}
	for _, call := range joinCalls(createCalls, updateCalls) {
		t.Run(call.name, func(t *testing.T) {
			rec, problem := serveThroughErrorMapping(t, func(ctx context.Context) error {
				return call.run(ctx, failingUseCase(invalid))
			})
			assertProblem(t, rec, problem, http.StatusBadRequest, invalid.Error())
		})
	}
}

func TestControllers_resourceAlreadyExists_mapsToConflictProblem(t *testing.T) {
	captureSideEffects(t)
	duplicate := fmt.Errorf("create: %w", repo.ErrAlreadyExists)
	for _, call := range joinCalls(createCalls, updateCalls) {
		t.Run(call.name, func(t *testing.T) {
			rec, problem := serveThroughErrorMapping(t, func(ctx context.Context) error {
				return call.run(ctx, failingUseCase(duplicate))
			})
			assertProblem(t, rec, problem, http.StatusConflict, "resource already exists")
		})
	}
}

func TestControllers_unexpectedFailure_mapsToInternalErrorWithoutLeakingCause(t *testing.T) {
	captureSideEffects(t)
	boom := errors.New("pq: connection refused to 10.0.0.5")
	for _, call := range joinCalls(readCalls, createCalls, updateCalls) {
		t.Run(call.name, func(t *testing.T) {
			rec, problem := serveThroughErrorMapping(t, func(ctx context.Context) error {
				return call.run(ctx, failingUseCase(boom))
			})
			assertProblem(t, rec, problem, http.StatusInternalServerError, "internal server error")
		})
	}
}

func TestControllers_writeFails_emitsNoSideEffects(t *testing.T) {
	pub, inv := captureSideEffects(t)
	use := failingUseCase(errors.New("write failed"))

	for _, call := range joinCalls(createCalls, updateCalls) {
		if err := call.run(context.Background(), use); err == nil {
			t.Errorf("%s err = nil, want the usecase failure", call.name)
		}
	}
	drainSideEffects(t)

	if events := pub.published(); len(events) != 0 {
		t.Errorf("published events = %+v, want none", events)
	}
	if got := inv.calls.Load(); got != 0 {
		t.Errorf("cache invalidations = %d, want 0", got)
	}
}

func TestControllers_writeSucceeds_publishesChangeEventAndInvalidatesCache(t *testing.T) {
	use := &stubUseCase{
		createWork: func(context.Context, command.WorkCreateCommand) (*aggregate.WorkAggregate, error) {
			return &aggregate.WorkAggregate{Uid: "w-new"}, nil
		},
		updateWork: func(_ context.Context, uid string, _ command.WorkUpdateCommand) (*aggregate.WorkAggregate, error) {
			return &aggregate.WorkAggregate{Uid: uid}, nil
		},
		createEdition: func(context.Context, command.EditionCreateCommand) (*aggregate.EditionAggregate, error) {
			return &aggregate.EditionAggregate{Uid: "e-new"}, nil
		},
		updateEdition: func(_ context.Context, uid string, _ command.EditionUpdateCommand) (*aggregate.EditionAggregate, error) {
			return &aggregate.EditionAggregate{Uid: uid}, nil
		},
		createAuthor: func(context.Context, command.AuthorCreateCommand) (*aggregate.AuthorAggregate, error) {
			return &aggregate.AuthorAggregate{Uid: "a-new"}, nil
		},
		updateAuthor: func(_ context.Context, uid string, _ command.AuthorUpdateCommand) (*aggregate.AuthorAggregate, error) {
			return &aggregate.AuthorAggregate{Uid: uid}, nil
		},
		createPublisher: func(context.Context, command.PublisherCreateCommand) (*aggregate.PublisherAggregate, error) {
			return &aggregate.PublisherAggregate{Uid: "p-new"}, nil
		},
		updatePublisher: func(_ context.Context, uid string, _ command.PublisherUpdateCommand) (*aggregate.PublisherAggregate, error) {
			return &aggregate.PublisherAggregate{Uid: uid}, nil
		},
		createSubject: func(context.Context, command.SubjectCreateCommand) (*aggregate.SubjectAggregate, error) {
			return &aggregate.SubjectAggregate{Uid: "s-new"}, nil
		},
		updateSubject: func(_ context.Context, uid string, _ command.SubjectUpdateCommand) (*aggregate.SubjectAggregate, error) {
			return &aggregate.SubjectAggregate{Uid: uid}, nil
		},
	}
	cases := []struct {
		call         controllerCall
		topic        string
		id           string
		kind         string
		metadataKeys []string
	}{
		{createCalls[0], "work.changed.v1", "w-new", "work.created", []string{"action", "title"}},
		{createCalls[1], "edition.changed.v1", "e-new", "edition.created", []string{"action", "title"}},
		{createCalls[2], "author.changed.v1", "a-new", "author.created", []string{"action", "name"}},
		{createCalls[3], "publisher.changed.v1", "p-new", "publisher.created", []string{"action", "name"}},
		{createCalls[4], "subject.changed.v1", "s-new", "subject.created", []string{"action", "name"}},
		{updateCalls[0], "work.changed.v1", "w-1", "work.updated", []string{"action"}},
		{updateCalls[1], "edition.changed.v1", "e-1", "edition.updated", []string{"action"}},
		{updateCalls[2], "author.changed.v1", "a-1", "author.updated", []string{"action"}},
		{updateCalls[3], "publisher.changed.v1", "p-1", "publisher.updated", []string{"action"}},
		{updateCalls[4], "subject.changed.v1", "s-1", "subject.updated", []string{"action"}},
	}
	for _, c := range cases {
		t.Run(c.call.name, func(t *testing.T) {
			pub, inv := captureSideEffects(t)

			if err := c.call.run(context.Background(), use); err != nil {
				t.Fatalf("err: %v", err)
			}
			drainSideEffects(t)

			events := pub.published()
			if len(events) != 1 {
				t.Fatalf("published events = %d, want 1", len(events))
			}
			if events[0].topic != c.topic || events[0].id != c.id || events[0].kind != c.kind {
				t.Errorf("event = %+v, want %s / %s / %s", events[0], c.topic, c.id, c.kind)
			}
			for _, key := range c.metadataKeys {
				if _, ok := events[0].metadata[key]; !ok {
					t.Errorf("metadata = %v, missing %q", events[0].metadata, key)
				}
			}
			if len(events[0].metadata) != len(c.metadataKeys) {
				t.Errorf("metadata = %v, want exactly %v", events[0].metadata, c.metadataKeys)
			}
			if got := inv.calls.Load(); got != 1 {
				t.Errorf("cache invalidations = %d, want 1", got)
			}
		})
	}
}

func TestWorkController_requestContextCancelledAfterCreate_stillRunsSideEffects(t *testing.T) {
	pub, inv := captureSideEffects(t)
	ctx, cancel := context.WithCancel(context.Background())
	use := &stubUseCase{
		createWork: func(context.Context, command.WorkCreateCommand) (*aggregate.WorkAggregate, error) {
			cancel()
			return &aggregate.WorkAggregate{Uid: "w-new"}, nil
		},
	}

	if _, err := NewWorkController(use).WorksCreate(ctx, openapi.WorkApiCredentials{}, openapi.WorkCreateRequest{Title: "Dune"}); err != nil {
		t.Fatalf("err: %v", err)
	}
	drainSideEffects(t)

	if events := pub.published(); len(events) != 1 || events[0].id != "w-new" {
		t.Errorf("published events = %+v, want the work.created event despite the cancelled request", events)
	}
	if got := inv.calls.Load(); got != 1 {
		t.Errorf("cache invalidations = %d, want 1", got)
	}
}

func TestWorkController_searchPagingGiven_passesFiltersAndPageToUsecase(t *testing.T) {
	var captured option.WorkQueryOption
	use := &stubUseCase{
		searchWorks: func(_ context.Context, opt option.WorkQueryOption) (int64, []*aggregate.WorkAggregate, error) {
			captured = opt
			return 0, nil, nil
		},
	}
	size, page := int32(5), int32(3)

	res, err := NewWorkController(use).WorksSearch(context.Background(), openapi.WorkApiCredentials{},
		ptrStr("dune"), ptrStr("a-1"), ptrStr("s-1"), ptrStr("OL1W"), &size, &page, ptrStr("title,desc"))
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if *captured.Title() != "dune" || *captured.AuthorUid() != "a-1" || *captured.SubjectUid() != "s-1" || *captured.OlKey() != "OL1W" {
		t.Errorf("filters = %+v", captured)
	}
	if captured.Page().GetSize() != 5 || captured.Page().GetPage() != 3 || captured.Page().GetSort() != "title,desc" {
		t.Errorf("page = %+v", captured.Page())
	}
	got, ok := res.(openapi.WorksSearchOk)
	if !ok {
		t.Fatalf("response = %T, want openapi.WorksSearchOk", res)
	}
	if got.Body.Items == nil || len(got.Body.Items) != 0 || got.Body.Total != 0 {
		t.Errorf("body = %+v, want an empty items array and zero total", got.Body)
	}
}

func TestWorkController_fullTextSearchHits_mapsScoreHighlightAndRelations(t *testing.T) {
	var captured option.WorkSearchOption
	createdAt := time.Date(2026, 10, 1, 9, 0, 0, 0, time.UTC)
	use := &stubUseCase{
		fullTextSearchWorks: func(_ context.Context, opt option.WorkSearchOption) (int64, []*repo.WorkSearchResult, error) {
			captured = opt
			return 7, []*repo.WorkSearchResult{
				{
					Uid:       "w-1",
					Title:     "The Hobbit",
					CoverUids: []string{"c-1"},
					Authors:   []repo.WorkSearchAuthor{{Uid: "a-1", Name: "Tolkien", AlternateNames: []string{"JRRT"}}},
					Subjects:  []repo.WorkSearchSubject{{Uid: "s-1", Name: "Fantasy"}},
					Score:     3.5,
					Highlight: &repo.WorkSearchHighlight{Title: []string{"The <em>Hobbit</em>"}},
					CreatedAt: createdAt,
				},
				{Uid: "w-2", Title: "Dune", Score: 1},
			}, nil
		},
	}
	size, page := int32(10), int32(1)

	res, err := NewWorkController(use).WorksFullTextSearch(context.Background(), openapi.WorkApiCredentials{},
		ptrStr("hobbit"), ptrStr("a-1"), ptrStr("s-1"), &size, &page)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if *captured.Query() != "hobbit" || *captured.AuthorUid() != "a-1" || *captured.SubjectUid() != "s-1" ||
		captured.Page().GetSize() != 10 || captured.Page().GetPage() != 1 || captured.Page().GetSort() != "" {
		t.Errorf("search option = %+v", captured)
	}
	got, ok := res.(openapi.WorksFullTextSearchOk)
	if !ok {
		t.Fatalf("response = %T, want openapi.WorksFullTextSearchOk", res)
	}
	if got.Body.Total != 7 || len(got.Body.Items) != 2 {
		t.Fatalf("body = %+v", got.Body)
	}
	first := got.Body.Items[0]
	if first.Uid != "w-1" || first.Score != 3.5 || !first.CreatedAt.Equal(createdAt) || first.CoverUids == nil || (*first.CoverUids)[0] != "c-1" {
		t.Errorf("first = %+v", first)
	}
	if len(first.Authors) != 1 || first.Authors[0].Uid != "a-1" || first.Authors[0].AlternateNames == nil || (*first.Authors[0].AlternateNames)[0] != "JRRT" {
		t.Errorf("authors = %+v", first.Authors)
	}
	if !reflect.DeepEqual(first.Subjects, []openapi.SubjectItem{{Uid: "s-1", Name: "Fantasy"}}) {
		t.Errorf("subjects = %+v", first.Subjects)
	}
	if first.Highlight == nil || first.Highlight.Title == nil || (*first.Highlight.Title)[0] != "The <em>Hobbit</em>" || first.Highlight.Description != nil {
		t.Errorf("highlight = %+v", first.Highlight)
	}
	second := got.Body.Items[1]
	if second.Highlight != nil || second.CoverUids != nil || second.Authors == nil || second.Subjects == nil {
		t.Errorf("second = %+v, want no highlight, no covers and empty relation arrays", second)
	}
}

func TestWorkController_fullTextSearchUsecaseFails_propagatesError(t *testing.T) {
	wantErr := errors.New("opensearch down")
	use := &stubUseCase{
		fullTextSearchWorks: func(context.Context, option.WorkSearchOption) (int64, []*repo.WorkSearchResult, error) {
			return 0, nil, wantErr
		},
	}

	res, err := NewWorkController(use).WorksFullTextSearch(context.Background(), openapi.WorkApiCredentials{}, nil, nil, nil, nil, nil)

	if !errors.Is(err, wantErr) || res != nil {
		t.Errorf("res = %v, err = %v, want nil and %v", res, err, wantErr)
	}
}

func TestWorkController_editionsOfWorkRequested_scopesSearchToThatWork(t *testing.T) {
	var captured option.EditionQueryOption
	use := &stubUseCase{
		searchEditions: func(_ context.Context, opt option.EditionQueryOption) (int64, []*aggregate.EditionAggregate, error) {
			captured = opt
			return 2, []*aggregate.EditionAggregate{
				{Uid: "e-1", Title: "First", Languages: []string{"eng"}, Work: &aggregate.WorkAggregate{Uid: "w-1", Title: "The Hobbit"},
					Publisher: &aggregate.PublisherAggregate{Uid: "p-1", Name: "Allen & Unwin"}},
				{Uid: "e-2", Title: "Second"},
			}, nil
		},
	}
	size, page := int32(2), int32(0)

	res, err := NewWorkController(use).WorksEditions(context.Background(), openapi.WorkApiCredentials{}, "w-1", &size, &page, ptrStr("title"))
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if captured.WorkUid() == nil || *captured.WorkUid() != "w-1" || captured.Title() != nil || captured.PublisherUid() != nil {
		t.Errorf("edition option = %+v, want only the work uid filter", captured)
	}
	if captured.Page().GetSize() != 2 || captured.Page().GetSort() != "title" {
		t.Errorf("page = %+v", captured.Page())
	}
	got, ok := res.(openapi.WorksEditionsOk)
	if !ok {
		t.Fatalf("response = %T, want openapi.WorksEditionsOk", res)
	}
	if got.Body.Total != 2 || len(got.Body.Items) != 2 {
		t.Fatalf("body = %+v", got.Body)
	}
	first := got.Body.Items[0]
	if first.Work != (openapi.EditionWorkRef{Uid: "w-1", Title: "The Hobbit"}) || first.Publisher == nil || first.Publisher.Name != "Allen & Unwin" ||
		first.Languages == nil || (*first.Languages)[0] != "eng" {
		t.Errorf("first = %+v", first)
	}
	if got.Body.Items[1].Publisher != nil || got.Body.Items[1].Languages != nil {
		t.Errorf("second = %+v, want no publisher and no languages", got.Body.Items[1])
	}
}

func TestWorkController_editionsOfWorkUsecaseFails_propagatesError(t *testing.T) {
	wantErr := errors.New("db down")
	use := &stubUseCase{
		searchEditions: func(context.Context, option.EditionQueryOption) (int64, []*aggregate.EditionAggregate, error) {
			return 0, nil, wantErr
		},
	}

	_, err := NewWorkController(use).WorksEditions(context.Background(), openapi.WorkApiCredentials{}, "w-1", nil, nil, nil)

	if !errors.Is(err, wantErr) {
		t.Errorf("err = %v, want %v", err, wantErr)
	}
}

func TestWorkController_updateRelationsGiven_passesEveryFieldAndMapsBody(t *testing.T) {
	var captured command.WorkUpdateCommand
	updatedAt := time.Date(2026, 10, 2, 0, 0, 0, 0, time.UTC)
	use := &stubUseCase{
		updateWork: func(_ context.Context, uid string, cmd command.WorkUpdateCommand) (*aggregate.WorkAggregate, error) {
			captured = cmd
			return &aggregate.WorkAggregate{
				Uid: uid, Title: "Good Omens", OlKey: cmd.OlKey, CoverUids: cmd.CoverUids, UpdatedAt: &updatedAt,
				Authors:  []*aggregate.AuthorAggregate{{Uid: "a-1", Name: "Pratchett"}, {Uid: "a-2", Name: "Gaiman"}},
				Subjects: []*aggregate.SubjectAggregate{{Uid: "s-1", Name: "Humor"}},
			}, nil
		},
	}
	captureSideEffects(t)
	covers, authors, subjects := []string{"c-1"}, []string{"a-1", "a-2"}, []string{"Humor"}

	res, err := NewWorkController(use).WorksUpdate(context.Background(), openapi.WorkApiCredentials{}, "w-1", openapi.WorkUpdateRequest{
		Description:      ptrStr("the nice and accurate prophecies"),
		CoverUids:        &covers,
		FirstPublishDate: ptrStr("1990"),
		OlKey:            ptrStr("OL5W"),
		AuthorUids:       &authors,
		SubjectNames:     &subjects,
	})
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	want := command.WorkUpdateCommand{
		Description:      ptrStr("the nice and accurate prophecies"),
		CoverUids:        covers,
		FirstPublishDate: ptrStr("1990"),
		OlKey:            ptrStr("OL5W"),
		AuthorUids:       authors,
		SubjectNames:     subjects,
	}
	if !reflect.DeepEqual(captured, want) {
		t.Errorf("command = %+v, want %+v", captured, want)
	}
	got, ok := res.(openapi.WorksUpdateOk)
	if !ok {
		t.Fatalf("response = %T, want openapi.WorksUpdateOk", res)
	}
	if got.Body.Uid != "w-1" || got.Body.OlKey == nil || *got.Body.OlKey != "OL5W" || got.Body.UpdatedAt == nil || !got.Body.UpdatedAt.Equal(updatedAt) {
		t.Errorf("body = %+v", got.Body)
	}
	if len(got.Body.Authors) != 2 || got.Body.Authors[1].Name != "Gaiman" || len(got.Body.Subjects) != 1 || got.Body.Subjects[0].Uid != "s-1" {
		t.Errorf("authors = %+v, subjects = %+v", got.Body.Authors, got.Body.Subjects)
	}
}

func TestWorkController_updateOmitsRelations_leavesRelationFieldsNil(t *testing.T) {
	var captured command.WorkUpdateCommand
	use := &stubUseCase{
		updateWork: func(_ context.Context, uid string, cmd command.WorkUpdateCommand) (*aggregate.WorkAggregate, error) {
			captured = cmd
			return &aggregate.WorkAggregate{Uid: uid}, nil
		},
	}
	captureSideEffects(t)

	if _, err := NewWorkController(use).WorksUpdate(context.Background(), openapi.WorkApiCredentials{}, "w-1", openapi.WorkUpdateRequest{Title: ptrStr("Renamed")}); err != nil {
		t.Fatalf("err: %v", err)
	}

	if captured.AuthorUids != nil || captured.SubjectNames != nil || captured.CoverUids != nil {
		t.Errorf("command = %+v, want nil relations so the repository keeps them", captured)
	}
}

func TestEditionController_searchFiltersGiven_passesEveryFilterAndMapsItems(t *testing.T) {
	var captured option.EditionQueryOption
	use := &stubUseCase{
		searchEditions: func(_ context.Context, opt option.EditionQueryOption) (int64, []*aggregate.EditionAggregate, error) {
			captured = opt
			return 1, []*aggregate.EditionAggregate{{Uid: "e-1", Title: "First edition", Isbn13: ptrStr("9780261102217")}}, nil
		},
	}
	size, page := int32(50), int32(2)

	res, err := NewEditionController(use).EditionsSearch(context.Background(), openapi.EditionApiCredentials{},
		ptrStr("first"), ptrStr("978"), ptrStr("w-1"), ptrStr("p-1"), ptrStr("eng"), ptrStr("OL1M"), &size, &page, ptrStr("created_at,desc"))
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if *captured.Title() != "first" || *captured.Isbn() != "978" || *captured.WorkUid() != "w-1" ||
		*captured.PublisherUid() != "p-1" || *captured.Language() != "eng" || *captured.OlKey() != "OL1M" {
		t.Errorf("filters = %+v", captured)
	}
	if captured.Page().GetSize() != 50 || captured.Page().GetPage() != 2 || captured.Page().GetSort() != "created_at,desc" {
		t.Errorf("page = %+v", captured.Page())
	}
	got, ok := res.(openapi.EditionsSearchOk)
	if !ok {
		t.Fatalf("response = %T, want openapi.EditionsSearchOk", res)
	}
	if got.Body.Total != 1 || len(got.Body.Items) != 1 || got.Body.Items[0].Uid != "e-1" ||
		got.Body.Items[0].Isbn13 == nil || *got.Body.Items[0].Isbn13 != "9780261102217" {
		t.Errorf("body = %+v", got.Body)
	}
}

func TestEditionController_searchUsecaseFails_propagatesError(t *testing.T) {
	wantErr := errors.New("db down")
	use := &stubUseCase{
		searchEditions: func(context.Context, option.EditionQueryOption) (int64, []*aggregate.EditionAggregate, error) {
			return 0, nil, wantErr
		},
	}

	res, err := NewEditionController(use).EditionsSearch(context.Background(), openapi.EditionApiCredentials{}, nil, nil, nil, nil, nil, nil, nil, nil, nil)

	if !errors.Is(err, wantErr) || res != nil {
		t.Errorf("res = %v, err = %v, want nil and %v", res, err, wantErr)
	}
}

func TestEditionController_createRequestValid_passesCommandAndPointsLocationAtNewEdition(t *testing.T) {
	var captured command.EditionCreateCommand
	use := &stubUseCase{
		createEdition: func(_ context.Context, cmd command.EditionCreateCommand) (*aggregate.EditionAggregate, error) {
			captured = cmd
			return &aggregate.EditionAggregate{Uid: "e-new"}, nil
		},
	}
	captureSideEffects(t)
	pages := int32(310)
	covers, languages := []string{"c-1"}, []string{"eng", "kor"}

	res, err := NewEditionController(use).EditionsCreate(context.Background(), openapi.EditionApiCredentials{}, openapi.EditionCreateRequest{
		Title:          "First edition",
		Isbn10:         ptrStr("0261102214"),
		Isbn13:         ptrStr("9780261102217"),
		NumberOfPages:  &pages,
		PublishDate:    ptrStr("1937-09-21"),
		CoverUids:      &covers,
		Languages:      &languages,
		PhysicalFormat: ptrStr("Hardcover"),
		Description:    ptrStr("first printing"),
		WorkUid:        "w-1",
		PublisherUid:   ptrStr("p-1"),
		OlKey:          ptrStr("OL1M"),
	})
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	want := command.EditionCreateCommand{
		Title:          "First edition",
		Isbn10:         ptrStr("0261102214"),
		Isbn13:         ptrStr("9780261102217"),
		NumberOfPages:  &pages,
		PublishDate:    ptrStr("1937-09-21"),
		CoverUids:      covers,
		Languages:      languages,
		PhysicalFormat: ptrStr("Hardcover"),
		Description:    ptrStr("first printing"),
		WorkUid:        "w-1",
		PublisherUid:   ptrStr("p-1"),
		OlKey:          ptrStr("OL1M"),
	}
	if !reflect.DeepEqual(captured, want) {
		t.Errorf("command = %+v, want %+v", captured, want)
	}
	got, ok := res.(openapi.EditionsCreateCreated)
	if !ok {
		t.Fatalf("response = %T, want openapi.EditionsCreateCreated", res)
	}
	if got.Location != "/api/v1/editions/e-new" {
		t.Errorf("Location = %q, want /api/v1/editions/e-new", got.Location)
	}
}

func TestEditionController_updateMovesWorkAndPublisher_mapsRefsInBody(t *testing.T) {
	var captured command.EditionUpdateCommand
	use := &stubUseCase{
		updateEdition: func(_ context.Context, uid string, cmd command.EditionUpdateCommand) (*aggregate.EditionAggregate, error) {
			captured = cmd
			return &aggregate.EditionAggregate{
				Uid: uid, Title: "Reissue", CoverUids: []string{"c-2"},
				Work:      &aggregate.WorkAggregate{Uid: "w-2", Title: "Dune"},
				Publisher: &aggregate.PublisherAggregate{Uid: "p-2", Name: "Ace", Address: ptrStr("New York")},
			}, nil
		},
	}
	captureSideEffects(t)
	pages := int32(412)

	res, err := NewEditionController(use).EditionsUpdate(context.Background(), openapi.EditionApiCredentials{}, "e-1", openapi.EditionUpdateRequest{
		Title:         ptrStr("Reissue"),
		NumberOfPages: &pages,
		WorkUid:       ptrStr("w-2"),
		PublisherUid:  ptrStr("p-2"),
	})
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if captured.Title == nil || *captured.Title != "Reissue" || captured.NumberOfPages == nil || *captured.NumberOfPages != 412 ||
		captured.WorkUid == nil || *captured.WorkUid != "w-2" || captured.PublisherUid == nil || *captured.PublisherUid != "p-2" ||
		captured.CoverUids != nil || captured.Languages != nil {
		t.Errorf("command = %+v", captured)
	}
	got, ok := res.(openapi.EditionsUpdateOk)
	if !ok {
		t.Fatalf("response = %T, want openapi.EditionsUpdateOk", res)
	}
	if got.Body.Uid != "e-1" || got.Body.Work != (openapi.EditionWorkRef{Uid: "w-2", Title: "Dune"}) {
		t.Errorf("body = %+v", got.Body)
	}
	if got.Body.Publisher == nil || got.Body.Publisher.Uid != "p-2" || got.Body.Publisher.Address == nil || *got.Body.Publisher.Address != "New York" {
		t.Errorf("publisher = %+v", got.Body.Publisher)
	}
	if got.Body.CoverUids == nil || (*got.Body.CoverUids)[0] != "c-2" || got.Body.Languages != nil {
		t.Errorf("cover_uids = %v, languages = %v", got.Body.CoverUids, got.Body.Languages)
	}
}

func TestEditionController_updatedEditionHasNoPublisher_omitsPublisherFromBody(t *testing.T) {
	use := &stubUseCase{
		updateEdition: func(_ context.Context, uid string, _ command.EditionUpdateCommand) (*aggregate.EditionAggregate, error) {
			return &aggregate.EditionAggregate{Uid: uid, Work: &aggregate.WorkAggregate{Uid: "w-1", Title: "The Hobbit"}}, nil
		},
	}
	captureSideEffects(t)

	res, err := NewEditionController(use).EditionsUpdate(context.Background(), openapi.EditionApiCredentials{}, "e-1", openapi.EditionUpdateRequest{})
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	got, ok := res.(openapi.EditionsUpdateOk)
	if !ok {
		t.Fatalf("response = %T, want openapi.EditionsUpdateOk", res)
	}
	if got.Body.Publisher != nil || got.Body.Work.Uid != "w-1" {
		t.Errorf("body = %+v, want the work ref and no publisher", got.Body)
	}
}

func TestPublisherController_createRequestValid_passesCommandAndPointsLocationAtNewPublisher(t *testing.T) {
	var captured command.PublisherCreateCommand
	use := &stubUseCase{
		createPublisher: func(_ context.Context, cmd command.PublisherCreateCommand) (*aggregate.PublisherAggregate, error) {
			captured = cmd
			return &aggregate.PublisherAggregate{Uid: "p-new"}, nil
		},
	}
	captureSideEffects(t)

	res, err := NewPublisherController(use).PublishersCreate(context.Background(), openapi.PublisherApiCredentials{}, openapi.PublisherCreateRequest{
		Name: "Allen & Unwin", Address: ptrStr("London"), OlKey: ptrStr("OL1P"),
	})
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	want := command.PublisherCreateCommand{Name: "Allen & Unwin", Address: ptrStr("London"), OlKey: ptrStr("OL1P")}
	if !reflect.DeepEqual(captured, want) {
		t.Errorf("command = %+v, want %+v", captured, want)
	}
	got, ok := res.(openapi.PublishersCreateCreated)
	if !ok {
		t.Fatalf("response = %T, want openapi.PublishersCreateCreated", res)
	}
	if got.Location != "/api/v1/publishers/p-new" {
		t.Errorf("Location = %q, want /api/v1/publishers/p-new", got.Location)
	}
}

func TestPublisherController_publisherExists_mapsEveryFieldOnRead(t *testing.T) {
	use := &stubUseCase{
		findPublisher: func(_ context.Context, uid string) (*aggregate.PublisherAggregate, error) {
			return &aggregate.PublisherAggregate{Uid: uid, Name: "Ace", Address: ptrStr("New York"), OlKey: ptrStr("OL2P")}, nil
		},
	}

	res, err := NewPublisherController(use).PublishersRead(context.Background(), openapi.PublisherApiCredentials{}, "p-2")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	got, ok := res.(openapi.PublishersReadOk)
	if !ok {
		t.Fatalf("response = %T, want openapi.PublishersReadOk", res)
	}
	want := openapi.PublisherFindResponse{Uid: "p-2", Name: "Ace", Address: ptrStr("New York"), OlKey: ptrStr("OL2P")}
	if !reflect.DeepEqual(got.Body, want) {
		t.Errorf("body = %+v, want %+v", got.Body, want)
	}
}

func TestPublisherController_updateRequestValid_passesCommandAndMapsBody(t *testing.T) {
	var capturedUid string
	var captured command.PublisherUpdateCommand
	use := &stubUseCase{
		updatePublisher: func(_ context.Context, uid string, cmd command.PublisherUpdateCommand) (*aggregate.PublisherAggregate, error) {
			capturedUid, captured = uid, cmd
			return &aggregate.PublisherAggregate{Uid: uid, Name: *cmd.Name, Address: cmd.Address}, nil
		},
	}
	captureSideEffects(t)

	res, err := NewPublisherController(use).PublishersUpdate(context.Background(), openapi.PublisherApiCredentials{}, "p-1", openapi.PublisherUpdateRequest{
		Name: ptrStr("HarperCollins"), Address: ptrStr("Glasgow"),
	})
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if capturedUid != "p-1" || !reflect.DeepEqual(captured, command.PublisherUpdateCommand{Name: ptrStr("HarperCollins"), Address: ptrStr("Glasgow")}) {
		t.Errorf("uid = %q, command = %+v", capturedUid, captured)
	}
	got, ok := res.(openapi.PublishersUpdateOk)
	if !ok {
		t.Fatalf("response = %T, want openapi.PublishersUpdateOk", res)
	}
	want := openapi.PublisherUpdateResponse{Uid: "p-1", Name: "HarperCollins", Address: ptrStr("Glasgow")}
	if !reflect.DeepEqual(got.Body, want) {
		t.Errorf("body = %+v, want %+v", got.Body, want)
	}
}

func TestPublisherController_listUsecaseFails_propagatesError(t *testing.T) {
	wantErr := errors.New("db down")
	use := &stubUseCase{
		findAllPublishers: func(context.Context, option.PublisherQueryOption) (int64, []*aggregate.PublisherAggregate, error) {
			return 0, nil, wantErr
		},
	}

	res, err := NewPublisherController(use).PublishersGetAll(context.Background(), openapi.PublisherApiCredentials{}, nil, nil, nil, nil)

	if !errors.Is(err, wantErr) || res != nil {
		t.Errorf("res = %v, err = %v, want nil and %v", res, err, wantErr)
	}
}

func TestAuthorController_listPagingGiven_passesNameAndPageAndMapsItems(t *testing.T) {
	var captured option.AuthorQueryOption
	use := &stubUseCase{
		findAllAuthors: func(_ context.Context, opt option.AuthorQueryOption) (int64, []*aggregate.AuthorAggregate, error) {
			captured = opt
			return 12, []*aggregate.AuthorAggregate{{Uid: "a-1", Name: "Asimov", PhotoUids: []string{"ph-1"}}}, nil
		},
	}
	size, page := int32(1), int32(4)

	res, err := NewAuthorController(use).AuthorsGetAll(context.Background(), openapi.AuthorApiCredentials{}, ptrStr("asi"), &size, &page, ptrStr("name"))
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if captured.Name() == nil || *captured.Name() != "asi" || captured.Page().GetSize() != 1 || captured.Page().GetPage() != 4 || captured.Page().GetSort() != "name" {
		t.Errorf("option = %+v", captured)
	}
	got, ok := res.(openapi.AuthorsGetAllOk)
	if !ok {
		t.Fatalf("response = %T, want openapi.AuthorsGetAllOk", res)
	}
	if got.Body.Total != 12 || len(got.Body.Items) != 1 || got.Body.Items[0].Uid != "a-1" ||
		got.Body.Items[0].PhotoUids == nil || (*got.Body.Items[0].PhotoUids)[0] != "ph-1" {
		t.Errorf("body = %+v", got.Body)
	}
}

func TestAuthorController_updateRequestValid_passesCommandAndMapsBody(t *testing.T) {
	var captured command.AuthorUpdateCommand
	use := &stubUseCase{
		updateAuthor: func(_ context.Context, uid string, cmd command.AuthorUpdateCommand) (*aggregate.AuthorAggregate, error) {
			captured = cmd
			return &aggregate.AuthorAggregate{Uid: uid, Name: "Asimov", Bio: cmd.Bio, AlternateNames: cmd.AlternateNames}, nil
		},
	}
	captureSideEffects(t)
	alternates := []string{"Paul French"}

	res, err := NewAuthorController(use).AuthorsUpdate(context.Background(), openapi.AuthorApiCredentials{}, "a-1", openapi.AuthorUpdateRequest{
		Bio: ptrStr("biochemist"), DeathDate: ptrStr("1992-04-06"), AlternateNames: &alternates,
	})
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	want := command.AuthorUpdateCommand{Bio: ptrStr("biochemist"), DeathDate: ptrStr("1992-04-06"), AlternateNames: alternates}
	if !reflect.DeepEqual(captured, want) {
		t.Errorf("command = %+v, want %+v", captured, want)
	}
	got, ok := res.(openapi.AuthorsUpdateOk)
	if !ok {
		t.Fatalf("response = %T, want openapi.AuthorsUpdateOk", res)
	}
	if got.Body.Uid != "a-1" || got.Body.Bio == nil || *got.Body.Bio != "biochemist" ||
		got.Body.AlternateNames == nil || (*got.Body.AlternateNames)[0] != "Paul French" || got.Body.PhotoUids != nil {
		t.Errorf("body = %+v", got.Body)
	}
}

func TestSubjectController_listPagingGiven_passesNameAndPageAndMapsItems(t *testing.T) {
	var captured option.SubjectQueryOption
	use := &stubUseCase{
		findAllSubjects: func(_ context.Context, opt option.SubjectQueryOption) (int64, []*aggregate.SubjectAggregate, error) {
			captured = opt
			return 2, []*aggregate.SubjectAggregate{{Uid: "s-1", Name: "Fantasy"}, {Uid: "s-2", Name: "Fantastic Fiction"}}, nil
		},
	}
	size := int32(2)

	res, err := NewSubjectController(use).SubjectsGetAll(context.Background(), openapi.SubjectApiCredentials{}, ptrStr("fanta"), &size, nil, nil)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if captured.Name() == nil || *captured.Name() != "fanta" || captured.Page().GetSize() != 2 || captured.Page().Page() != nil {
		t.Errorf("option = %+v", captured)
	}
	got, ok := res.(openapi.SubjectsGetAllOk)
	if !ok {
		t.Fatalf("response = %T, want openapi.SubjectsGetAllOk", res)
	}
	if got.Body.Total != 2 || !reflect.DeepEqual(got.Body.Items, []openapi.SubjectItem{{Uid: "s-1", Name: "Fantasy"}, {Uid: "s-2", Name: "Fantastic Fiction"}}) {
		t.Errorf("body = %+v", got.Body)
	}
}

func TestSubjectController_updateRequestValid_passesCommandAndMapsBody(t *testing.T) {
	use := &stubUseCase{
		updateSubject: func(_ context.Context, uid string, cmd command.SubjectUpdateCommand) (*aggregate.SubjectAggregate, error) {
			return &aggregate.SubjectAggregate{Uid: uid, Name: *cmd.Name}, nil
		},
	}
	captureSideEffects(t)

	res, err := NewSubjectController(use).SubjectsUpdate(context.Background(), openapi.SubjectApiCredentials{}, "s-1", openapi.SubjectUpdateRequest{Name: ptrStr("High Fantasy")})
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	got, ok := res.(openapi.SubjectsUpdateOk)
	if !ok {
		t.Fatalf("response = %T, want openapi.SubjectsUpdateOk", res)
	}
	if got.Body != (openapi.SubjectUpdateResponse{Uid: "s-1", Name: "High Fantasy"}) {
		t.Errorf("body = %+v", got.Body)
	}
}

func TestControllers_listUsecaseFails_propagatesError(t *testing.T) {
	wantErr := errors.New("db down")
	use := &stubUseCase{
		findAllAuthors: func(context.Context, option.AuthorQueryOption) (int64, []*aggregate.AuthorAggregate, error) {
			return 0, nil, wantErr
		},
		findAllSubjects: func(context.Context, option.SubjectQueryOption) (int64, []*aggregate.SubjectAggregate, error) {
			return 0, nil, wantErr
		},
	}

	if _, err := NewAuthorController(use).AuthorsGetAll(context.Background(), openapi.AuthorApiCredentials{}, nil, nil, nil, nil); !errors.Is(err, wantErr) {
		t.Errorf("AuthorsGetAll err = %v, want %v", err, wantErr)
	}
	if _, err := NewSubjectController(use).SubjectsGetAll(context.Background(), openapi.SubjectApiCredentials{}, nil, nil, nil, nil); !errors.Is(err, wantErr) {
		t.Errorf("SubjectsGetAll err = %v, want %v", err, wantErr)
	}
}
