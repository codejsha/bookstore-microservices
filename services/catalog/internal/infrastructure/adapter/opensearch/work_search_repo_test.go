package opensearch

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/codejsha/shared-library-go/pkg/pagination"

	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/model/option"
)

func listOpt(opts ...option.WorkQueryOptionFunc) option.WorkQueryOption {
	return option.NewWorkQueryOption(opts...)
}

func pageOpt(size, page int32, sort string) pagination.PageOption {
	return pagination.NewPageOption(&size, &page, &sort)
}

func decodeListQuery(t *testing.T, opt option.WorkQueryOption) map[string]any {
	t.Helper()
	body, err := buildListQuery(opt)
	if err != nil {
		t.Fatalf("buildListQuery: %v", err)
	}
	var got map[string]any
	if err := json.Unmarshal(body, &got); err != nil {
		t.Fatalf("unmarshal body: %v", err)
	}
	return got
}

func TestBuildListQuery_WhenNoFilters_ReturnsMatchAll(t *testing.T) {
	got := decodeListQuery(t, listOpt())

	query, ok := got["query"].(map[string]any)
	if !ok {
		t.Fatalf("query = %v, want object", got["query"])
	}
	if _, ok := query["match_all"]; !ok {
		t.Errorf("query = %v, want match_all", query)
	}
	if got["track_total_hits"] != true {
		t.Errorf("track_total_hits = %v, want true", got["track_total_hits"])
	}
	if _, ok := got["highlight"]; ok {
		t.Error("list query must not request highlighting")
	}
}

func TestBuildListQuery_WhenFiltersGiven_ReturnsBoolMustClauses(t *testing.T) {
	title := "hobbit"
	authorUid := "a-1"
	subjectUid := "s-1"
	olKey := "/works/OL1W"

	got := decodeListQuery(t, listOpt(
		option.WorkQueryOption{}.WithTitle(&title),
		option.WorkQueryOption{}.WithAuthorUid(&authorUid),
		option.WorkQueryOption{}.WithSubjectUid(&subjectUid),
		option.WorkQueryOption{}.WithOlKey(&olKey),
	))

	boolQuery, ok := got["query"].(map[string]any)["bool"].(map[string]any)
	if !ok {
		t.Fatalf("query = %v, want bool query", got["query"])
	}

	must, _ := boolQuery["must"].([]any)
	if len(must) != 1 {
		t.Fatalf("must = %v, want one clause", must)
	}
	wantMust := map[string]any{"match": map[string]any{"title": map[string]any{"query": title}}}
	if !reflect.DeepEqual(must[0], wantMust) {
		t.Errorf("must[0] = %v, want %v", must[0], wantMust)
	}

	filter, _ := boolQuery["filter"].([]any)
	if len(filter) != 3 {
		t.Fatalf("filter = %v, want three clauses", filter)
	}
	wantFilter := []any{
		map[string]any{"nested": map[string]any{
			"path":  "authors",
			"query": map[string]any{"term": map[string]any{"authors.uid": authorUid}},
		}},
		map[string]any{"nested": map[string]any{
			"path":  "subjects",
			"query": map[string]any{"term": map[string]any{"subjects.uid": subjectUid}},
		}},
		map[string]any{"term": map[string]any{"ol_key": olKey}},
	}
	if !reflect.DeepEqual(filter, wantFilter) {
		t.Errorf("filter = %v, want %v", filter, wantFilter)
	}
}

func TestBuildListQuery_WhenFiltersBlank_ReturnsMatchAll(t *testing.T) {
	blank := ""
	got := decodeListQuery(t, listOpt(
		option.WorkQueryOption{}.WithTitle(&blank),
		option.WorkQueryOption{}.WithAuthorUid(&blank),
		option.WorkQueryOption{}.WithSubjectUid(&blank),
		option.WorkQueryOption{}.WithOlKey(&blank),
	))

	if _, ok := got["query"].(map[string]any)["match_all"]; !ok {
		t.Errorf("query = %v, want match_all", got["query"])
	}
}

func TestBuildListQuery_WhenPageAndSizeGiven_ReturnsFromAndSize(t *testing.T) {
	got := decodeListQuery(t, listOpt(
		option.WorkQueryOption{}.WithPage(pageOpt(20, 3, "")),
	))

	if got["from"] != float64(40) {
		t.Errorf("from = %v, want 40", got["from"])
	}
	if got["size"] != float64(20) {
		t.Errorf("size = %v, want 20", got["size"])
	}
}

func TestBuildSort(t *testing.T) {
	uidTiebreaker := map[string]any{"uid": map[string]any{"order": "asc"}}
	cases := []struct {
		name string
		sort string
		want []map[string]any
	}{
		{
			name: "whenSortEmpty_returnsTiebreakerOnly",
			sort: "",
			want: []map[string]any{uidTiebreaker},
		},
		{
			name: "whenFieldUnknown_returnsTiebreakerOnly",
			sort: "bogus,desc",
			want: []map[string]any{uidTiebreaker},
		},
		{
			name: "whenSortingByTitle_returnsRawSubfield",
			sort: "title,desc",
			want: []map[string]any{
				{"title.raw": map[string]any{"order": "desc"}},
				uidTiebreaker,
			},
		},
		{
			name: "whenColonSeparatorWithoutDirection_returnsAscending",
			sort: "title:whatever",
			want: []map[string]any{
				{"title.raw": map[string]any{"order": "asc"}},
				uidTiebreaker,
			},
		},
		{
			name: "whenSortingByDateField_returnsMissingLast",
			sort: "updated_at,desc",
			want: []map[string]any{
				{"updated_at": map[string]any{"order": "desc", "missing": "_last"}},
				uidTiebreaker,
			},
		},
		{
			name: "whenMultipleTerms_preservesOrder",
			sort: " title ; created_at,DESC ",
			want: []map[string]any{
				{"title.raw": map[string]any{"order": "asc"}},
				{"created_at": map[string]any{"order": "desc", "missing": "_last"}},
				uidTiebreaker,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := buildSort(c.sort)
			if !reflect.DeepEqual(got, c.want) {
				t.Errorf("buildSort(%q) = %v, want %v", c.sort, got, c.want)
			}
		})
	}
}
