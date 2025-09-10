package keycloak

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/codejsha/shared-library-go/pkg/rest/client"

	"github.com/codejsha/bookstore-microservices/identity/internal/config"

	"github.com/codejsha/bookstore-microservices/identity/internal/domain/constant"
)

type Introspection struct {
	Active bool   `json:"active"`
	Sub    string `json:"sub"`
	Exp    int64  `json:"exp"`
	Iat    int64  `json:"iat"`
}

type Introspector interface {
	Introspect(ctx context.Context, token string) (Introspection, error)
}

type introspector struct {
	cfg         *config.Config
	restyClient *resty.Client
}

func NewIntrospector(cfg *config.Config, restyClient *client.RestyClient) Introspector {
	return &introspector{cfg: cfg, restyClient: restyClient.Client}
}

func (i *introspector) Introspect(ctx context.Context, token string) (Introspection, error) {
	resp, err := i.restyClient.R().
		SetContext(ctx).
		SetDebug(i.cfg.App.Logging.IsDebugEnabled).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(map[string]string{
			"token":         token,
			"client_id":     i.cfg.Keycloak.ClientId,
			"client_secret": i.cfg.Keycloak.ClientSecret,
		}).
		Post(i.cfg.Keycloak.Url + fmt.Sprintf(
			"/realms/%s/protocol/openid-connect/token/introspect", i.cfg.Keycloak.Realm))
	if err != nil {
		return Introspection{}, fmt.Errorf("failed to introspect token: %w", err)
	}
	if resp.StatusCode() != 200 {
		return Introspection{}, fmt.Errorf("failed to introspect token: %s", resp.Status())
	}

	var result Introspection
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return Introspection{}, fmt.Errorf("failed to unmarshal introspection response: %w", err)
	}
	return result, nil
}

type cachingIntrospector struct {
	inner Introspector
	cache TokenCache
}

func NewCachingIntrospector(inner Introspector, cache TokenCache) Introspector {
	return &cachingIntrospector{inner: inner, cache: cache}
}

func introspectCacheKey(token string) string {
	sum := sha256.Sum256([]byte(token))
	return "identity:introspect:" + hex.EncodeToString(sum[:])
}

func (c *cachingIntrospector) Introspect(ctx context.Context, token string) (Introspection, error) {
	key := introspectCacheKey(token)

	var cached Introspection
	if err := c.cache.Get(ctx, key, &cached); err == nil {
		return cached, nil
	}

	result, err := c.inner.Introspect(ctx, token)
	if err != nil {
		return Introspection{}, err
	}
	_ = c.cache.Set(ctx, key, result, cacheTTLFor(result))
	return result, nil
}

func cacheTTLFor(result Introspection) time.Duration {
	ttl := constant.IntrospectCacheTTL
	if result.Active && result.Exp > 0 {
		if remaining := time.Until(time.Unix(result.Exp, 0)); remaining < ttl {
			ttl = remaining
		}
	}
	if ttl <= 0 {
		return time.Second
	}
	return ttl
}
