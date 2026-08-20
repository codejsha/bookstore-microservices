package opensearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"

	"github.com/codejsha/bookstore-microservices/catalog/internal/application/port/repo"
	"github.com/codejsha/bookstore-microservices/catalog/internal/config"
	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/model/option"
	"github.com/codejsha/bookstore-microservices/catalog/internal/infrastructure/support"
	"github.com/codejsha/bookstore-microservices/catalog/internal/infrastructure/support/utils"
)

type workSearchRepo struct {
	client *support.OpensearchClient
	index  string
}

func NewWorkSearchRepository(client *support.OpensearchClient, cfg *config.OpensearchConfig) repo.WorkSearchRepo {
	return &workSearchRepo{client: client, index: cfg.Index}
}

func (r *workSearchRepo) FindAll(
	ctx context.Context,
	opt option.WorkQueryOption,
) (int64, []*repo.WorkSearchResult, error) {
	body, err := buildListQuery(opt)
	if err != nil {
		return 0, nil, err
	}
	return r.search(ctx, body)
}

func (r *workSearchRepo) Search(
	ctx context.Context,
	opt option.WorkSearchOption,
) (int64, []*repo.WorkSearchResult, error) {
	body, err := buildQuery(opt)
	if err != nil {
		return 0, nil, err
	}
	return r.search(ctx, body)
}

func (r *workSearchRepo) search(ctx context.Context, body []byte) (int64, []*repo.WorkSearchResult, error) {
	resp, err := r.client.Search(ctx, &opensearchapi.SearchReq{
		Indices: []string{r.index},
		Body:    bytes.NewReader(body),
	})
	if err != nil {
		return 0, nil, fmt.Errorf("opensearch search: %w", err)
	}

	var total int64
	if resp.Hits.Total.Value > 0 {
		total = int64(resp.Hits.Total.Value)
	}

	results := make([]*repo.WorkSearchResult, 0, len(resp.Hits.Hits))
	for _, hit := range resp.Hits.Hits {
		r, err := hitToResult(hit.Source, hit.Highlight, float32(hit.Score))
		if err != nil {
			return 0, nil, err
		}
		results = append(results, r)
	}
	return total, results, nil
}

func buildQuery(opt option.WorkSearchOption) ([]byte, error) {
	must := []map[string]any{}
	filter := []map[string]any{}

	if q := opt.Query(); q != nil && *q != "" {
		must = append(must, map[string]any{
			"multi_match": map[string]any{
				"query":     *q,
				"fields":    []string{"title^3", "description", "authors.name", "subjects.name"},
				"type":      "best_fields",
				"fuzziness": "AUTO",
			},
		})
	}

	if uidStr := opt.AuthorUid(); uidStr != nil && *uidStr != "" {
		filter = append(filter, map[string]any{
			"nested": map[string]any{
				"path": "authors",
				"query": map[string]any{
					"term": map[string]any{"authors.uid": *uidStr},
				},
			},
		})
	}

	if uidStr := opt.SubjectUid(); uidStr != nil && *uidStr != "" {
		filter = append(filter, map[string]any{
			"nested": map[string]any{
				"path": "subjects",
				"query": map[string]any{
					"term": map[string]any{"subjects.uid": *uidStr},
				},
			},
		})
	}

	offset, size := utils.PageOffsetLimit(opt.Page())
	from := offset

	query := map[string]any{}
	if len(must) > 0 || len(filter) > 0 {
		query["bool"] = map[string]any{"must": must, "filter": filter}
	} else {
		query["match_all"] = map[string]any{}
	}

	return json.Marshal(map[string]any{
		"from":  from,
		"size":  size,
		"query": query,
		"highlight": map[string]any{
			"fields": map[string]any{
				"title":       map[string]any{},
				"description": map[string]any{},
			},
		},
		"track_total_hits": true,
	})
}

func buildListQuery(opt option.WorkQueryOption) ([]byte, error) {
	must := []map[string]any{}
	filter := []map[string]any{}

	if title := opt.Title(); title != nil && *title != "" {
		must = append(must, map[string]any{
			"match": map[string]any{
				"title": map[string]any{"query": *title},
			},
		})
	}

	if uidStr := opt.AuthorUid(); uidStr != nil && *uidStr != "" {
		filter = append(filter, map[string]any{
			"nested": map[string]any{
				"path": "authors",
				"query": map[string]any{
					"term": map[string]any{"authors.uid": *uidStr},
				},
			},
		})
	}

	if uidStr := opt.SubjectUid(); uidStr != nil && *uidStr != "" {
		filter = append(filter, map[string]any{
			"nested": map[string]any{
				"path": "subjects",
				"query": map[string]any{
					"term": map[string]any{"subjects.uid": *uidStr},
				},
			},
		})
	}

	if olKey := opt.OlKey(); olKey != nil && *olKey != "" {
		filter = append(filter, map[string]any{
			"term": map[string]any{"ol_key": *olKey},
		})
	}

	query := map[string]any{}
	if len(must) > 0 || len(filter) > 0 {
		query["bool"] = map[string]any{"must": must, "filter": filter}
	} else {
		query["match_all"] = map[string]any{}
	}

	from, size := utils.PageOffsetLimit(opt.Page())

	return json.Marshal(map[string]any{
		"from":             from,
		"size":             size,
		"query":            query,
		"sort":             buildSort(opt.Page().GetSort()),
		"track_total_hits": true,
	})
}

var sortFields = map[string]string{
	"title":      "title.raw",
	"created_at": "created_at",
	"updated_at": "updated_at",
}

func buildSort(sort string) []map[string]any {
	out := []map[string]any{}
	for _, term := range strings.Split(sort, ";") {
		name, desc := parseSortTerm(term)
		fieldName, ok := sortFields[name]
		if !ok {
			continue
		}
		order := "asc"
		if desc {
			order = "desc"
		}
		spec := map[string]any{"order": order}
		if fieldName == "created_at" || fieldName == "updated_at" {
			spec["missing"] = "_last"
		}
		out = append(out, map[string]any{fieldName: spec})
	}
	return append(out, map[string]any{"uid": map[string]any{"order": "asc"}})
}

func parseSortTerm(term string) (name string, desc bool) {
	term = strings.TrimSpace(term)
	sep := strings.IndexAny(term, ",:")
	if sep < 0 {
		return term, false
	}
	name = strings.TrimSpace(term[:sep])
	dir := strings.ToLower(strings.TrimSpace(term[sep+1:]))
	return name, dir == "desc"
}

type bookDoc struct {
	Uid              string       `json:"uid"`
	Title            string       `json:"title"`
	Description      *string      `json:"description,omitempty"`
	CoverUids        []string     `json:"cover_uids,omitempty"`
	FirstPublishDate *string      `json:"first_publish_date,omitempty"`
	OlKey            *string      `json:"ol_key,omitempty"`
	Authors          []docAuthor  `json:"authors"`
	Subjects         []docSubject `json:"subjects"`
	Editions         []docEdition `json:"editions"`
	CreatedAt        time.Time    `json:"created_at"`
	UpdatedAt        *time.Time   `json:"updated_at,omitempty"`
}

type docAuthor struct {
	Uid            string   `json:"uid"`
	Name           string   `json:"name"`
	Bio            *string  `json:"bio,omitempty"`
	BirthDate      *string  `json:"birth_date,omitempty"`
	DeathDate      *string  `json:"death_date,omitempty"`
	OlKey          *string  `json:"ol_key,omitempty"`
	AlternateNames []string `json:"alternate_names,omitempty"`
}

type docSubject struct {
	Uid  string `json:"uid"`
	Name string `json:"name"`
}

type docEdition struct {
	Uid            string        `json:"uid"`
	Title          string        `json:"title"`
	Isbn10         *string       `json:"isbn10,omitempty"`
	Isbn13         *string       `json:"isbn13,omitempty"`
	PublishDate    *string       `json:"publish_date,omitempty"`
	Languages      []string      `json:"languages,omitempty"`
	PhysicalFormat *string       `json:"physical_format,omitempty"`
	Description    *string       `json:"description,omitempty"`
	Publisher      *docPublisher `json:"publisher,omitempty"`
}

type docPublisher struct {
	Uid  string `json:"uid"`
	Name string `json:"name"`
}

func hitToResult(
	raw json.RawMessage,
	highlight map[string][]string,
	score float32,
) (*repo.WorkSearchResult, error) {
	var d bookDoc
	if err := json.Unmarshal(raw, &d); err != nil {
		return nil, fmt.Errorf("unmarshal hit: %w", err)
	}

	result := &repo.WorkSearchResult{
		Uid:              d.Uid,
		Title:            d.Title,
		Description:      d.Description,
		CoverUids:        d.CoverUids,
		FirstPublishDate: d.FirstPublishDate,
		OlKey:            d.OlKey,
		Authors:          mapAuthors(d.Authors),
		Subjects:         mapSubjects(d.Subjects),
		Editions:         mapEditions(d.Editions),
		Score:            score,
		CreatedAt:        d.CreatedAt,
		UpdatedAt:        d.UpdatedAt,
	}
	if len(highlight) > 0 {
		result.Highlight = &repo.WorkSearchHighlight{
			Title:       highlight["title"],
			Description: highlight["description"],
		}
	}
	return result, nil
}

func mapAuthors(src []docAuthor) []repo.WorkSearchAuthor {
	out := make([]repo.WorkSearchAuthor, len(src))
	for i, a := range src {
		out[i] = repo.WorkSearchAuthor{
			Uid:            a.Uid,
			Name:           a.Name,
			Bio:            a.Bio,
			BirthDate:      a.BirthDate,
			DeathDate:      a.DeathDate,
			OlKey:          a.OlKey,
			AlternateNames: a.AlternateNames,
		}
	}
	return out
}

func mapSubjects(src []docSubject) []repo.WorkSearchSubject {
	out := make([]repo.WorkSearchSubject, len(src))
	for i, s := range src {
		out[i] = repo.WorkSearchSubject{
			Uid:  s.Uid,
			Name: s.Name,
		}
	}
	return out
}

func mapEditions(src []docEdition) []repo.WorkSearchEdition {
	out := make([]repo.WorkSearchEdition, len(src))
	for i, e := range src {
		out[i] = repo.WorkSearchEdition{
			Uid:            e.Uid,
			Title:          e.Title,
			Isbn10:         e.Isbn10,
			Isbn13:         e.Isbn13,
			PublishDate:    e.PublishDate,
			Languages:      e.Languages,
			PhysicalFormat: e.PhysicalFormat,
			Description:    e.Description,
		}
		if e.Publisher != nil && e.Publisher.Uid != "" {
			out[i].Publisher = &repo.WorkSearchPublisher{
				Uid:  e.Publisher.Uid,
				Name: e.Publisher.Name,
			}
		}
	}
	return out
}
