package keycloak

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/codejsha/shared-library-go/pkg/rest/client"

	"github.com/codejsha/bookstore-microservices/identity/internal/config"
)

const redactedValue = "[REDACTED]"

var (
	formSecretRe = regexp.MustCompile(`(?i)(password|client_secret|refresh_token|access_token|code)=[^&]*`)
	jsonSecretRe = regexp.MustCompile(`(?i)"(password|client_secret|access_token|refresh_token|value|secret)"\s*:\s*"[^"]*"`)
)

func installDebugRedaction(c *resty.Client) {
	c.OnRequestLog(func(rl *resty.RequestLog) error {
		redactHeaders(rl.Header)
		rl.Body = redactBody(rl.Body)
		return nil
	})
	c.OnResponseLog(func(rl *resty.ResponseLog) error {
		redactHeaders(rl.Header)
		rl.Body = redactBody(rl.Body)
		return nil
	})
}

func redactHeaders(h http.Header) {
	for _, key := range []string{"Authorization", "Cookie", "Set-Cookie", "Proxy-Authorization"} {
		if len(h.Values(key)) > 0 {
			h.Set(key, redactedValue)
		}
	}
}

func redactBody(body string) string {
	if body == "" {
		return body
	}
	body = formSecretRe.ReplaceAllString(body, `${1}=`+redactedValue)
	body = jsonSecretRe.ReplaceAllString(body, `"${1}":"`+redactedValue+`"`)
	return body
}

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

type tokenErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

type tokenGrantError struct {
	op          string
	status      int
	code        string
	description string
}

func (e *tokenGrantError) Error() string {
	if e.code == "" {
		return fmt.Sprintf("failed to %s: status %d", e.op, e.status)
	}
	return fmt.Sprintf("failed to %s: status %d: %s: %s", e.op, e.status, e.code, e.description)
}

const (
	adminTokenCacheKey = "identity:keycloak:admin-token"
	adminTokenLockKey  = "identity:keycloak:admin-token-refresh"
)

const (
	cacheOpTimeout      = 2 * time.Second
	refreshLockTTL      = 10 * time.Second
	refreshLockWait     = 5 * time.Second
	refreshLockPoll     = 100 * time.Millisecond
	tokenRequestTimeout = 8 * time.Second
	expirySkew          = 30 * time.Second
)

func fresh(expiresAt time.Time) bool {
	return time.Now().Add(expirySkew).Before(expiresAt)
}

type TokenCache interface {
	Get(ctx context.Context, key string, dest interface{}) error
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	TryLock(ctx context.Context, key string, ttl time.Duration) (string, error)
	Unlock(ctx context.Context, key, token string) error
}

type cachedAdminToken struct {
	AccessToken string    `json:"access_token"`
	ExpiresAt   time.Time `json:"expires_at"`
}

type AdminTokenHelper struct {
	cfg         *config.Config
	restyClient *resty.Client
	tokenCache  TokenCache
	accessToken string
	expiresAt   time.Time
	mu          sync.Mutex
}

func NewAdminTokenHelper(
	cfg *config.Config,
	restyClient *client.RestyClient,
	tokenCache TokenCache,
) *AdminTokenHelper {
	installDebugRedaction(restyClient.Client)
	return &AdminTokenHelper{
		cfg:         cfg,
		restyClient: restyClient.Client,
		tokenCache:  tokenCache,
	}
}

func (h *AdminTokenHelper) AccessToken() (string, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if fresh(h.expiresAt) {
		return h.accessToken, nil
	}
	h.loadFromCache()
	if fresh(h.expiresAt) {
		return h.accessToken, nil
	}
	lockToken := h.acquireRefreshLock()
	if lockToken != "" {
		defer h.releaseRefreshLock(lockToken)
		h.loadFromCache()
	}
	if fresh(h.expiresAt) {
		return h.accessToken, nil
	}

	accessToken, expiresAt, err := h.requestToken()
	if err != nil {
		return "", fmt.Errorf("failed to refresh token: %w", err)
	}

	h.accessToken = accessToken
	h.expiresAt = expiresAt
	h.storeToCache()
	return h.accessToken, nil
}

func (h *AdminTokenHelper) loadFromCache() {
	if h.tokenCache == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), cacheOpTimeout)
	defer cancel()

	var cached cachedAdminToken
	if err := h.tokenCache.Get(ctx, adminTokenCacheKey, &cached); err != nil {
		return
	}
	if fresh(cached.ExpiresAt) {
		h.accessToken = cached.AccessToken
		h.expiresAt = cached.ExpiresAt
	}
}

func (h *AdminTokenHelper) storeToCache() {
	if h.tokenCache == nil {
		return
	}
	ttl := time.Until(h.expiresAt)
	if ttl <= 0 {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), cacheOpTimeout)
	defer cancel()

	_ = h.tokenCache.Set(ctx, adminTokenCacheKey, cachedAdminToken{
		AccessToken: h.accessToken,
		ExpiresAt:   h.expiresAt,
	}, ttl)
}

func (h *AdminTokenHelper) acquireRefreshLock() string {
	if h.tokenCache == nil {
		return ""
	}
	deadline := time.Now().Add(refreshLockWait)
	for {
		ctx, cancel := context.WithTimeout(context.Background(), cacheOpTimeout)
		token, err := h.tokenCache.TryLock(ctx, adminTokenLockKey, refreshLockTTL)
		cancel()
		if err != nil {
			return ""
		}
		if token != "" {
			return token
		}
		if time.Now().After(deadline) {
			return ""
		}
		time.Sleep(refreshLockPoll)
		h.loadFromCache()
		if fresh(h.expiresAt) {
			return ""
		}
	}
}

func (h *AdminTokenHelper) releaseRefreshLock(token string) {
	if h.tokenCache == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), cacheOpTimeout)
	defer cancel()
	_ = h.tokenCache.Unlock(ctx, adminTokenLockKey, token)
}

func (h *AdminTokenHelper) requestToken() (string, time.Time, error) {
	const op = "request client credentials token"
	ctx, cancel := context.WithTimeout(context.Background(), tokenRequestTimeout)
	defer cancel()
	restyResp, err := h.restyClient.R().
		SetContext(ctx).
		SetDebug(h.cfg.App.Logging.IsDebugEnabled).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(map[string]string{
			"grant_type":    "client_credentials",
			"client_id":     h.cfg.Keycloak.ClientId,
			"client_secret": h.cfg.Keycloak.ClientSecret,
		}).
		Post(h.cfg.Keycloak.Url + fmt.Sprintf("/realms/%s/protocol/openid-connect/token", h.cfg.Keycloak.Realm))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to %s: %v", op, err)
	}
	if restyResp.IsError() {
		var rejection tokenErrorResponse
		_ = json.Unmarshal(restyResp.Body(), &rejection)
		return "", time.Time{}, &tokenGrantError{
			op:          op,
			status:      restyResp.StatusCode(),
			code:        rejection.Error,
			description: rejection.ErrorDescription,
		}
	}

	var response tokenResponse
	if err := json.Unmarshal(restyResp.Body(), &response); err != nil {
		return "", time.Time{}, fmt.Errorf("failed to unmarshal %s response: %v", op, err)
	}
	if response.AccessToken == "" {
		return "", time.Time{}, fmt.Errorf("failed to %s: response carries no access token", op)
	}
	return response.AccessToken, time.Now().Add(time.Duration(response.ExpiresIn) * time.Second), nil
}
