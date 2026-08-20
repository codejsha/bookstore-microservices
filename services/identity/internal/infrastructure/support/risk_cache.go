package support

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/codejsha/bookstore-microservices/identity/internal/application/port/security"
)

const riskIndexKey = "identity:risk:index"

func riskKey(sub string) string { return "identity:risk:" + sub }

func autoRiskKey(sub string) string { return "identity:risk:auto:" + sub }

type RiskCache struct {
	*CacheClient
}

func NewRiskCache(c *CacheClient) *RiskCache {
	return &RiskCache{CacheClient: c}
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
	subs, err := r.SMembers(ctx, riskIndexKey)
	if err != nil {
		return nil, err
	}
	entries := make([]security.RiskEntry, 0, len(subs))
	for _, sub := range subs {
		entry, err := r.Check(ctx, sub)
		if err != nil {
			return nil, err
		}
		if entry == nil {
			_ = r.SRem(ctx, riskIndexKey, sub)
			continue
		}
		entries = append(entries, *entry)
	}
	return entries, nil
}

func ProvideRiskStore(r *RiskCache) security.RiskStore     { return r }
func ProvideRiskChecker(r *RiskCache) security.RiskChecker { return r }
