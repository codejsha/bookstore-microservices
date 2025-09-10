package support

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/codejsha/bookstore-microservices/identity/internal/application/port/security"

	"github.com/codejsha/bookstore-microservices/identity/internal/domain/constant"
)

func revokedKey(sub string) string { return "identity:revoked:" + sub }

type RevocationCache struct {
	*CacheClient
	now func() time.Time
}

func NewRevocationCache(c *CacheClient) *RevocationCache {
	return &RevocationCache{CacheClient: c, now: time.Now}
}

func (r *RevocationCache) Revoke(ctx context.Context, sub string) error {
	return r.Set(ctx, revokedKey(sub), r.now().Unix(), constant.RevocationTTL)
}

func (r *RevocationCache) IsRevoked(ctx context.Context, sub string, tokenIssuedAt int64) (bool, error) {
	var revokedAt int64
	err := r.Get(ctx, revokedKey(sub), &revokedAt)
	if errors.Is(err, redis.Nil) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return tokenRevoked(revokedAt, tokenIssuedAt), nil
}

func tokenRevoked(revokedAt, tokenIssuedAt int64) bool {
	if tokenIssuedAt == 0 {
		return true
	}
	return tokenIssuedAt < revokedAt+int64(constant.RevocationSkew.Seconds())
}

func ProvideSessionRevoker(r *RevocationCache) security.SessionRevoker       { return r }
func ProvideRevocationChecker(r *RevocationCache) security.RevocationChecker { return r }
