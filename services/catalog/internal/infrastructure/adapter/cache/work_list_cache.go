package cache

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/codejsha/bookstore-microservices/catalog/internal/application/port/repo"
	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/model/option"
)

const (
	epochKey      = "catalog:works:epoch"
	listKeyPrefix = "catalog:works:list:"
	listTTL       = 60 * time.Second
)

type Store interface {
	Get(ctx context.Context, key string, dest interface{}) error
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Incr(ctx context.Context, key string) (int64, error)
}

type cachedWorkList struct {
	Total   int64              `json:"total"`
	Results []*repo.WorkResult `json:"results"`
}

type WorkListCache struct {
	store Store
}

func NewWorkListCache(store Store) *WorkListCache {
	return &WorkListCache{store: store}
}

func (c *WorkListCache) enabled() bool {
	return c != nil && c.store != nil
}

func (c *WorkListCache) Lookup(ctx context.Context, opt option.WorkQueryOption) (total int64, results []*repo.WorkResult, ok bool) {
	if !c.enabled() {
		return 0, nil, false
	}
	var out cachedWorkList
	if err := c.store.Get(ctx, c.listKey(ctx, opt), &out); err != nil {
		return 0, nil, false
	}
	return out.Total, out.Results, true
}

func (c *WorkListCache) Save(ctx context.Context, opt option.WorkQueryOption, total int64, results []*repo.WorkResult) {
	if !c.enabled() {
		return
	}
	_ = c.store.Set(ctx, c.listKey(ctx, opt), cachedWorkList{Total: total, Results: results}, listTTL)
}

func (c *WorkListCache) Invalidate(ctx context.Context) error {
	if !c.enabled() {
		return nil
	}
	_, err := c.store.Incr(ctx, epochKey)
	return err
}

func (c *WorkListCache) listKey(ctx context.Context, opt option.WorkQueryOption) string {
	return fmt.Sprintf("%s%d:%s", listKeyPrefix, c.currentEpoch(ctx), hashWorkQuery(opt))
}

func (c *WorkListCache) currentEpoch(ctx context.Context) int64 {
	var epoch int64
	if err := c.store.Get(ctx, epochKey, &epoch); err != nil {
		return 0
	}
	return epoch
}

func hashWorkQuery(opt option.WorkQueryOption) string {
	p := opt.Page()
	parts := []string{
		"size=" + strconv.FormatInt(int64(p.GetSize()), 10),
		"page=" + strconv.FormatInt(int64(p.GetPage()), 10),
		"sort=" + p.GetSort(),
		"title=" + derefString(opt.Title()),
		"author=" + derefString(opt.AuthorUid()),
		"subject=" + derefString(opt.SubjectUid()),
		"olkey=" + derefString(opt.OlKey()),
	}
	sum := sha256.Sum256([]byte(strings.Join(parts, "&")))
	return hex.EncodeToString(sum[:])
}

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
