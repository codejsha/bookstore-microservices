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
	AccessToken      string `json:"access_token"`
	RefreshToken     string `json:"refresh_token"`
	ExpiresIn        int64  `json:"expires_in"`
	RefreshExpiresIn int64  `json:"refresh_expires_in"`
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
	AccessToken      string    `json:"access_token"`
	RefreshToken     string    `json:"refresh_token"`
	ExpiresAt        time.Time `json:"expires_at"`
	RefreshExpiresAt time.Time `json:"refresh_expires_at"`
}

type AdminTokenHelper struct {
	cfg              *config.Config
	restyClient      *resty.Client
	tokenCache       TokenCache
	accessToken      string
	refreshToken     string
	expiresAt        time.Time
	refreshExpiresAt time.Time
	mu               sync.Mutex
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

func (h *AdminTokenHelper) GetTokens() (string, string, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if fresh(h.expiresAt) {
		return h.accessToken, h.refreshToken, nil
	}
	h.loadFromCache()
	if fresh(h.expiresAt) {
		return h.accessToken, h.refreshToken, nil
	}
	lockToken := h.acquireRefreshLock()
	if lockToken != "" {
		defer h.releaseRefreshLock(lockToken)
		h.loadFromCache()
	}
	if fresh(h.expiresAt) {
		return h.accessToken, h.refreshToken, nil
	}

	var (
		accessToken, refreshToken   string
		expiresAt, refreshExpiresAt time.Time
		err                         error
	)
	if time.Now().Before(h.refreshExpiresAt) {
		accessToken, refreshToken, expiresAt, refreshExpiresAt, err = h.exchangeTokens()
	} else {
		accessToken, refreshToken, expiresAt, refreshExpiresAt, err = h.fetchTokens()
	}

	if err != nil {
		return "", "", fmt.Errorf("failed to refresh tokens: %v", err)
	}

	h.accessToken = accessToken
	h.refreshToken = refreshToken
	h.expiresAt = expiresAt
	h.refreshExpiresAt = refreshExpiresAt
	h.storeToCache()
	return h.accessToken, h.refreshToken, nil
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
	if time.Now().Before(cached.RefreshExpiresAt) {
		h.refreshToken = cached.RefreshToken
		h.refreshExpiresAt = cached.RefreshExpiresAt
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
	ttl := time.Until(h.refreshExpiresAt)
	if ttl <= 0 {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), cacheOpTimeout)
	defer cancel()

	_ = h.tokenCache.Set(ctx, adminTokenCacheKey, cachedAdminToken{
		AccessToken:      h.accessToken,
		RefreshToken:     h.refreshToken,
		ExpiresAt:        h.expiresAt,
		RefreshExpiresAt: h.refreshExpiresAt,
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

func (h *AdminTokenHelper) fetchTokens() (string, string, time.Time, time.Time, error) {
	data := map[string]string{
		"client_id":     h.cfg.Keycloak.ClientId,
		"client_secret": h.cfg.Keycloak.ClientSecret,
		"grant_type":    "password",
		"username":      h.cfg.Keycloak.RealmAdminUsername,
		"password":      h.cfg.Keycloak.RealmAdminPassword,
		"scope":         "openid profile email",
	}
	ctx, cancel := context.WithTimeout(context.Background(), tokenRequestTimeout)
	defer cancel()
	restyResp, err := h.restyClient.R().
		SetContext(ctx).
		SetDebug(h.cfg.App.Logging.IsDebugEnabled).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(data).
		Post(h.cfg.Keycloak.Url + fmt.Sprintf("/realms/%s/protocol/openid-connect/token", h.cfg.Keycloak.Realm))
	if err != nil {
		return "", "", time.Time{}, time.Time{}, fmt.Errorf("failed to fetch tokens: %v", err)
	}

	var response tokenResponse
	err = json.Unmarshal(restyResp.Body(), &response)
	if err != nil {
		return "", "", time.Time{}, time.Time{}, fmt.Errorf("failed to unmarshal token exchange response: %v", err)
	}

	now := time.Now()
	accessToken := response.AccessToken
	refreshToken := response.RefreshToken
	expiresAt := now.Add(time.Duration(response.ExpiresIn) * time.Second)
	refreshExpiresAt := now.Add(time.Duration(response.RefreshExpiresIn) * time.Second)
	return accessToken, refreshToken, expiresAt, refreshExpiresAt, nil
}

func (h *AdminTokenHelper) exchangeTokens() (string, string, time.Time, time.Time, error) {
	data := map[string]string{
		"client_id":     h.cfg.Keycloak.ClientId,
		"client_secret": h.cfg.Keycloak.ClientSecret,
		"grant_type":    "refresh_token",
		"refresh_token": h.refreshToken,
		"scope":         "openid profile email",
	}
	ctx, cancel := context.WithTimeout(context.Background(), tokenRequestTimeout)
	defer cancel()
	restyResp, err := h.restyClient.R().
		SetContext(ctx).
		SetDebug(h.cfg.App.Logging.IsDebugEnabled).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(data).
		Post(h.cfg.Keycloak.Url + fmt.Sprintf("/realms/%s/protocol/openid-connect/token", h.cfg.Keycloak.Realm))
	if err != nil {
		return "", "", time.Time{}, time.Time{}, fmt.Errorf("failed to exchange tokens: %v", err)
	}

	var response tokenResponse
	err = json.Unmarshal(restyResp.Body(), &response)
	if err != nil {
		return "", "", time.Time{}, time.Time{}, fmt.Errorf("failed to unmarshal token exchange response: %v", err)
	}

	now := time.Now()
	accessToken := response.AccessToken
	refreshToken := response.RefreshToken
	expiresAt := now.Add(time.Duration(response.ExpiresIn) * time.Second)
	refreshExpiresAt := now.Add(time.Duration(response.RefreshExpiresIn) * time.Second)
	return accessToken, refreshToken, expiresAt, refreshExpiresAt, nil
}
