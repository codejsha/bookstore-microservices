package support

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/codejsha/bookstore-microservices/identity/internal/application/port/security"
)

const riskIndexKey = "identity:risk:index"

func riskKey(sub string) string { return "identity:risk:" + sub }

func autoRiskKey(sub string) string { return "identity:risk:auto:" + sub }

type riskBatchStore interface {
	SMembers(ctx context.Context, key string) ([]string, error)
	MGet(ctx context.Context, keys ...string) ([]any, error)
	SRem(ctx context.Context, key string, members ...string) error
}

type RiskCache struct {
	*CacheClient
	batch riskBatchStore
}

func NewRiskCache(c *CacheClient) *RiskCache {
	return &RiskCache{CacheClient: c, batch: c}
}

func (r *RiskCache) Flag(ctx context.Context, entry security.RiskEntry, ttl time.Duration) error {
	if err := r.Set(ctx, riskKey(entry.Sub), entry, ttl); err != nil {
		return err
	}
	return r.SAdd(ctx, riskIndexKey, entry.Sub)
}

func (r *RiskCache) Unflag(ctx context.Context, sub string) error {
	if err := r.Del(ctx, riskKey(sub), autoRiskKey(sub)); err != nil {
		return err
	}
	return r.SRem(ctx, riskIndexKey, sub)
}

func (r *RiskCache) Check(ctx context.Context, sub string) (*security.RiskEntry, error) {
	entry, err := r.getEntry(ctx, riskKey(sub))
	if entry != nil || err != nil {
		return entry, err
	}
	return r.getEntry(ctx, autoRiskKey(sub))
}

func (r *RiskCache) getEntry(ctx context.Context, key string) (*security.RiskEntry, error) {
	var entry security.RiskEntry
	err := r.Get(ctx, key, &entry)
	if errors.Is(err, redis.Nil) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

func (r *RiskCache) List(ctx context.Context) ([]security.RiskEntry, error) {
	subs, err := r.batch.SMembers(ctx, riskIndexKey)
	if err != nil {
		return nil, err
	}
	if len(subs) == 0 {
		return []security.RiskEntry{}, nil
	}

	values, err := r.batch.MGet(ctx, riskKeysFor(subs, riskKey)...)
	if err != nil {
		return nil, err
	}
	entries, missing, err := decodeRiskValues(subs, values)
	if err != nil {
		return nil, err
	}
	if len(missing) == 0 {
		return entries, nil
	}

	autoValues, err := r.batch.MGet(ctx, riskKeysFor(missing, autoRiskKey)...)
	if err != nil {
		return nil, err
	}
	autoEntries, stale, err := decodeRiskValues(missing, autoValues)
	if err != nil {
		return nil, err
	}
	entries = append(entries, autoEntries...)
	if len(stale) > 0 {
		_ = r.batch.SRem(ctx, riskIndexKey, stale...)
	}

	return entries, nil
}

func riskKeysFor(subs []string, keyFn func(string) string) []string {
	keys := make([]string, len(subs))
	for i, sub := range subs {
		keys[i] = keyFn(sub)
	}
	return keys
}

func decodeRiskValues(subs []string, values []any) ([]security.RiskEntry, []string, error) {
	if len(values) != len(subs) {
		return nil, nil, fmt.Errorf("risk cache: got %d values for %d subjects", len(values), len(subs))
	}
	entries := make([]security.RiskEntry, 0, len(subs))
	missing := make([]string, 0, len(subs))
	for i, value := range values {
		var raw []byte
		switch v := value.(type) {
		case nil:
			missing = append(missing, subs[i])
			continue
		case string:
			raw = []byte(v)
		case []byte:
			raw = v
		default:
			return nil, nil, fmt.Errorf("risk cache: unexpected value type %T for subject %s", value, subs[i])
		}
		var entry security.RiskEntry
		if err := json.Unmarshal(raw, &entry); err != nil {
			return nil, nil, fmt.Errorf("risk cache: decode entry for subject %s: %w", subs[i], err)
		}
		entries = append(entries, entry)
	}
	return entries, missing, nil
}

func ProvideRiskStore(r *RiskCache) security.RiskStore     { return r }
func ProvideRiskChecker(r *RiskCache) security.RiskChecker { return r }
