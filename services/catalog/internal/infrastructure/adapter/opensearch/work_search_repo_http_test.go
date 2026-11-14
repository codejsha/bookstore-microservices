package opensearch

import (
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/codejsha/bookstore-microservices/catalog/internal/application/port/repo"
	"github.com/codejsha/bookstore-microservices/catalog/internal/config"
	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/model/option"
	"github.com/codejsha/bookstore-microservices/catalog/internal/infrastructure/support"
)

type searchRequest struct {
	method      string
	path        string
	contentType string
	body        map[string]any
}

type fakeSearchCluster struct {
	mu       sync.Mutex
	requests []searchRequest
	status   int
	response string
}

func (f *fakeSearchCluster) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	raw, _ := io.ReadAll(r.Body)
	var body map[string]any
	_ = json.Unmarshal(raw, &body)
	f.mu.Lock()
	f.requests = append(f.requests, searchRequest{
		method:      r.Method,
		path:        r.URL.Path,
		contentType: r.Header.Get("Content-Type"),
		body:        body,
	})
	status, response := f.status, f.response
	f.mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = io.WriteString(w, response)
}

func (f *fakeSearchCluster) onlyRequest(t *testing.T) searchRequest {
	t.Helper()
	f.mu.Lock()
	defer f.mu.Unlock()
	if len(f.requests) != 1 {
		t.Fatalf("search requests = %d, want 1", len(f.requests))
	}
	return f.requests[0]
}

func newSearchRepo(t *testing.T, status int, response string) (repo.WorkSearchRepo, *fakeSearchCluster) {
	t.Helper()
	cluster := &fakeSearchCluster{status: status, response: response}
	srv := httptest.NewServer(cluster)
	t.Cleanup(srv.Close)

	u, err := url.Parse(srv.URL)
	if err != nil {
		t.Fatalf("parse server url: %v", err)
	}
	host, portText, err := net.SplitHostPort(u.Host)
	if err != nil {
		t.Fatalf("split host port: %v", err)
	}
	port, err := strconv.Atoi(portText)
	if err != nil {
		t.Fatalf("port: %v", err)
	}
	cfg := &config.OpensearchConfig{Scheme: u.Scheme, Host: host, Port: port, Index: "books-test"}
	client, err := support.NewOpensearchClient(cfg)
	if err != nil {
		t.Fatalf("opensearch client: %v", err)
	}
	return NewWorkSearchRepository(client, cfg), cluster
}

func searchResponse(total int, hits ...string) string {
	return `{"took":3,"timed_out":false,"_shards":{"total":1,"successful":1,"skipped":0,"failed":0},` +
		`"hits":{"total":{"value":` + strconv.Itoa(total) + `,"relation":"eq"},"max_score":1.0,"hits":[` + strings.Join(hits, ",") + `]}}`
}

const hobbitHit = `{
	"_index": "books-test", "_id": "w-1", "_score": 4.5,
	"_source": {
		"uid": "w-1", "title": "The Hobbit", "description": "There and back again",
		"cover_uids": ["c-1"], "first_publish_date": "1937", "ol_key": "OL1W",
		"authors": [{"uid": "a-1", "name": "J.R.R. Tolkien", "bio": "philologist", "birth_date": "1892-01-03",
			"death_date": "1973-09-02", "ol_key": "OL26320A", "alternate_names": ["Tolkien"]}],
		"subjects": [{"uid": "s-1", "name": "Fantasy"}, {"uid": "s-2", "name": "Adventure"}],
		"editions": [
			{"uid": "e-1", "title": "First edition", "isbn10": "0261102214", "isbn13": "9780261102217",
				"publish_date": "1937-09-21", "languages": ["eng"], "physical_format": "Hardcover",
				"description": "first printing", "publisher": {"uid": "p-1", "name": "Allen & Unwin"}},
			{"uid": "e-2", "title": "Unattributed", "publisher": {"uid": "", "name": ""}},
			{"uid": "e-3", "title": "No publisher"}
		],
		"created_at": "2026-10-01T09:00:00Z", "updated_at": "2026-10-02T10:30:00Z"
	},
	"highlight": {"title": ["The <em>Hobbit</em>"], "description": ["There and <em>back</em> again"]}
}`

const duneHit = `{
	"_index": "books-test", "_id": "w-2", "_score": 1.25,
	"_source": {"uid": "w-2", "title": "Dune", "authors": [], "subjects": [], "editions": [],
		"created_at": "2026-10-03T00:00:00Z"}
}`

func TestWorkSearchRepoSearch_queryAndFiltersGiven_postsBoolQueryToConfiguredIndex(t *testing.T) {
	r, cluster := newSearchRepo(t, http.StatusOK, searchResponse(0))
	q, authorUid, subjectUid := "hobbit", "a-1", "s-1"

	_, _, err := r.Search(context.Background(), option.NewWorkSearchOption(
		option.WorkSearchOption{}.WithQuery(&q),
		option.WorkSearchOption{}.WithAuthorUid(&authorUid),
		option.WorkSearchOption{}.WithSubjectUid(&subjectUid),
		option.WorkSearchOption{}.WithPage(pageOpt(5, 2, "")),
	))
	if err != nil {
		t.Fatalf("Search: %v", err)
	}

	req := cluster.onlyRequest(t)
	if req.method != http.MethodPost || req.path != "/books-test/_search" {
		t.Errorf("request = %s %s, want POST /books-test/_search", req.method, req.path)
	}
	if !strings.HasPrefix(req.contentType, "application/json") {
		t.Errorf("content-type = %q, want application/json", req.contentType)
	}
	if req.body["from"] != float64(10) || req.body["size"] != float64(5) || req.body["track_total_hits"] != true {
		t.Errorf("from = %v, size = %v, track_total_hits = %v", req.body["from"], req.body["size"], req.body["track_total_hits"])
	}
	wantQuery := map[string]any{"bool": map[string]any{
		"must": []any{map[string]any{"multi_match": map[string]any{
			"query":     "hobbit",
			"fields":    []any{"title^3", "description", "authors.name", "subjects.name"},
			"type":      "best_fields",
			"fuzziness": "AUTO",
		}}},
		"filter": []any{
			map[string]any{"nested": map[string]any{
				"path":  "authors",
				"query": map[string]any{"term": map[string]any{"authors.uid": "a-1"}},
			}},
			map[string]any{"nested": map[string]any{
				"path":  "subjects",
				"query": map[string]any{"term": map[string]any{"subjects.uid": "s-1"}},
			}},
		},
	}}
	if !reflect.DeepEqual(req.body["query"], wantQuery) {
		t.Errorf("query = %v, want %v", req.body["query"], wantQuery)
	}
	wantHighlight := map[string]any{"fields": map[string]any{"title": map[string]any{}, "description": map[string]any{}}}
	if !reflect.DeepEqual(req.body["highlight"], wantHighlight) {
		t.Errorf("highlight = %v, want %v", req.body["highlight"], wantHighlight)
	}
	if _, ok := req.body["sort"]; ok {
		t.Errorf("sort = %v, want relevance order without an explicit sort", req.body["sort"])
	}
}

func TestWorkSearchRepoSearch_filtersOnlyGiven_sendsBoolQueryWithEmptyMust(t *testing.T) {
	r, cluster := newSearchRepo(t, http.StatusOK, searchResponse(0))
	authorUid := "a-1"

	_, _, err := r.Search(context.Background(), option.NewWorkSearchOption(
		option.WorkSearchOption{}.WithAuthorUid(&authorUid),
	))
	if err != nil {
		t.Fatalf("Search: %v", err)
	}

	boolQuery, ok := cluster.onlyRequest(t).body["query"].(map[string]any)["bool"].(map[string]any)
	if !ok {
		t.Fatalf("query is not a bool query")
	}
	if must, _ := boolQuery["must"].([]any); boolQuery["must"] == nil || len(must) != 0 {
		t.Errorf("must = %v, want an empty array", boolQuery["must"])
	}
	if filter, _ := boolQuery["filter"].([]any); len(filter) != 1 {
		t.Errorf("filter = %v, want one clause", boolQuery["filter"])
	}
}

func TestWorkSearchRepoSearch_noCriteriaGiven_sendsMatchAllWithDefaultPage(t *testing.T) {
	r, cluster := newSearchRepo(t, http.StatusOK, searchResponse(0))
	blank := ""

	_, _, err := r.Search(context.Background(), option.NewWorkSearchOption(
		option.WorkSearchOption{}.WithQuery(&blank),
		option.WorkSearchOption{}.WithAuthorUid(&blank),
		option.WorkSearchOption{}.WithSubjectUid(&blank),
	))
	if err != nil {
		t.Fatalf("Search: %v", err)
	}

	body := cluster.onlyRequest(t).body
	if !reflect.DeepEqual(body["query"], map[string]any{"match_all": map[string]any{}}) {
		t.Errorf("query = %v, want match_all", body["query"])
	}
	if body["from"] != float64(0) || body["size"] != float64(20) {
		t.Errorf("from = %v, size = %v, want 0 and the default size 20", body["from"], body["size"])
	}
}

func TestWorkSearchRepoSearch_hitsWithHighlight_parsesDocumentsScoresAndTotal(t *testing.T) {
	r, _ := newSearchRepo(t, http.StatusOK, searchResponse(42, hobbitHit, duneHit))

	total, got, err := r.Search(context.Background(), option.NewWorkSearchOption())
	if err != nil {
		t.Fatalf("Search: %v", err)
	}

	if total != 42 || len(got) != 2 {
		t.Fatalf("total = %d, len = %d, want 42 and 2", total, len(got))
	}
	updatedAt := time.Date(2026, 10, 2, 10, 30, 0, 0, time.UTC)
	str := func(s string) *string { return &s }
	wantHobbit := &repo.WorkSearchResult{
		Uid:              "w-1",
		Title:            "The Hobbit",
		Description:      str("There and back again"),
		CoverUids:        []string{"c-1"},
		FirstPublishDate: str("1937"),
		OlKey:            str("OL1W"),
		Authors: []repo.WorkSearchAuthor{{
			Uid: "a-1", Name: "J.R.R. Tolkien", Bio: str("philologist"), BirthDate: str("1892-01-03"),
			DeathDate: str("1973-09-02"), OlKey: str("OL26320A"), AlternateNames: []string{"Tolkien"},
		}},
		Subjects: []repo.WorkSearchSubject{{Uid: "s-1", Name: "Fantasy"}, {Uid: "s-2", Name: "Adventure"}},
		Editions: []repo.WorkSearchEdition{
			{
				Uid: "e-1", Title: "First edition", Isbn10: str("0261102214"), Isbn13: str("9780261102217"),
				PublishDate: str("1937-09-21"), Languages: []string{"eng"}, PhysicalFormat: str("Hardcover"),
				Description: str("first printing"), Publisher: &repo.WorkSearchPublisher{Uid: "p-1", Name: "Allen & Unwin"},
			},
			{Uid: "e-2", Title: "Unattributed"},
			{Uid: "e-3", Title: "No publisher"},
		},
		Score: 4.5,
		Highlight: &repo.WorkSearchHighlight{
			Title:       []string{"The <em>Hobbit</em>"},
			Description: []string{"There and <em>back</em> again"},
		},
		CreatedAt: time.Date(2026, 10, 1, 9, 0, 0, 0, time.UTC),
		UpdatedAt: &updatedAt,
	}
	if !reflect.DeepEqual(got[0], wantHobbit) {
		t.Errorf("first hit = %+v, want %+v", got[0], wantHobbit)
	}

	dune := got[1]
	if dune.Uid != "w-2" || dune.Score != 1.25 || dune.Highlight != nil || dune.UpdatedAt != nil || dune.Description != nil {
		t.Errorf("second hit = %+v", dune)
	}
	if dune.Authors == nil || len(dune.Authors) != 0 || len(dune.Subjects) != 0 || len(dune.Editions) != 0 {
		t.Errorf("second hit relations = %v/%v/%v, want empty", dune.Authors, dune.Subjects, dune.Editions)
	}
}

func TestWorkSearchRepoSearch_noHits_yieldsEmptySliceAndZeroTotal(t *testing.T) {
	r, _ := newSearchRepo(t, http.StatusOK, searchResponse(0))

	total, got, err := r.Search(context.Background(), option.NewWorkSearchOption())
	if err != nil {
		t.Fatalf("Search: %v", err)
	}

	if total != 0 || got == nil || len(got) != 0 {
		t.Errorf("total = %d, results = %v, want 0 and an empty slice", total, got)
	}
}

func TestWorkSearchRepoSearch_clusterRejectsQuery_wrapsUpstreamError(t *testing.T) {
	r, _ := newSearchRepo(t, http.StatusBadRequest,
		`{"error":{"root_cause":[{"type":"parsing_exception","reason":"bad query"}],"type":"parsing_exception","reason":"bad query"},"status":400}`)

	total, got, err := r.Search(context.Background(), option.NewWorkSearchOption())

	if err == nil || !strings.HasPrefix(err.Error(), "opensearch search: ") {
		t.Fatalf("err = %v, want it wrapped with the opensearch search prefix", err)
	}
	if !strings.Contains(err.Error(), "bad query") {
		t.Errorf("err = %v, want the upstream reason kept", err)
	}
	if total != 0 || got != nil {
		t.Errorf("total = %d, results = %v, want 0 and nil", total, got)
	}
}

func TestWorkSearchRepoSearch_hitSourceMalformed_failsInsteadOfDroppingTheHit(t *testing.T) {
	r, _ := newSearchRepo(t, http.StatusOK, searchResponse(2, duneHit,
		`{"_index":"books-test","_id":"w-9","_score":1.0,"_source":{"uid":"w-9","title":"Broken","created_at":"not-a-timestamp"}}`))

	total, got, err := r.Search(context.Background(), option.NewWorkSearchOption())

	if err == nil || !strings.HasPrefix(err.Error(), "unmarshal hit: ") {
		t.Fatalf("err = %v, want an unmarshal hit error", err)
	}
	if total != 0 || got != nil {
		t.Errorf("total = %d, results = %v, want 0 and nil", total, got)
	}
}

func TestWorkSearchRepoSearch_contextAlreadyCancelled_failsWithoutResults(t *testing.T) {
	r, _ := newSearchRepo(t, http.StatusOK, searchResponse(1, duneHit))
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, got, err := r.Search(ctx, option.NewWorkSearchOption())

	if err == nil || !strings.Contains(err.Error(), context.Canceled.Error()) {
		t.Errorf("err = %v, want the context cancellation surfaced", err)
	}
	if got != nil {
		t.Errorf("results = %v, want nil", got)
	}
}

func TestWorkSearchRepoFindAll_filtersAndSortGiven_postsListQueryAndParsesHits(t *testing.T) {
	r, cluster := newSearchRepo(t, http.StatusOK, searchResponse(1, duneHit))
	title, olKey := "dune", "OL2W"

	total, got, err := r.FindAll(context.Background(), option.NewWorkQueryOption(
		option.WorkQueryOption{}.WithTitle(&title),
		option.WorkQueryOption{}.WithOlKey(&olKey),
		option.WorkQueryOption{}.WithPage(pageOpt(10, 1, "title,desc;created_at")),
	))
	if err != nil {
		t.Fatalf("FindAll: %v", err)
	}

	if total != 1 || len(got) != 1 || got[0].Uid != "w-2" || got[0].Title != "Dune" || got[0].Highlight != nil {
		t.Errorf("total = %d, results = %+v", total, got)
	}
	req := cluster.onlyRequest(t)
	if req.method != http.MethodPost || req.path != "/books-test/_search" {
		t.Errorf("request = %s %s, want POST /books-test/_search", req.method, req.path)
	}
	if req.body["from"] != float64(10) || req.body["size"] != float64(10) {
		t.Errorf("from = %v, size = %v, want 10 and 10", req.body["from"], req.body["size"])
	}
	wantQuery := map[string]any{"bool": map[string]any{
		"must":   []any{map[string]any{"match": map[string]any{"title": map[string]any{"query": "dune"}}}},
		"filter": []any{map[string]any{"term": map[string]any{"ol_key": "OL2W"}}},
	}}
	if !reflect.DeepEqual(req.body["query"], wantQuery) {
		t.Errorf("query = %v, want %v", req.body["query"], wantQuery)
	}
	wantSort := []any{
		map[string]any{"title.raw": map[string]any{"order": "desc"}},
		map[string]any{"created_at": map[string]any{"order": "asc", "missing": "_last"}},
		map[string]any{"uid": map[string]any{"order": "asc"}},
	}
	if !reflect.DeepEqual(req.body["sort"], wantSort) {
		t.Errorf("sort = %v, want %v", req.body["sort"], wantSort)
	}
	if _, ok := req.body["highlight"]; ok {
		t.Error("list query must not request highlighting")
	}
}

func TestWorkSearchRepoFindAll_clusterFails_wrapsUpstreamError(t *testing.T) {
	r, _ := newSearchRepo(t, http.StatusInternalServerError,
		`{"error":{"root_cause":[{"type":"search_phase_execution_exception","reason":"all shards failed"}],"type":"search_phase_execution_exception","reason":"all shards failed"},"status":500}`)

	_, got, err := r.FindAll(context.Background(), option.NewWorkQueryOption())

	if err == nil || !strings.HasPrefix(err.Error(), "opensearch search: ") {
		t.Fatalf("err = %v, want it wrapped with the opensearch search prefix", err)
	}
	if got != nil {
		t.Errorf("results = %v, want nil", got)
	}
}
