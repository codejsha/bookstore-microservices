package keycloak

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"

	sharedconfig "github.com/codejsha/shared-library-go/pkg/config"

	"github.com/codejsha/bookstore-microservices/identity/internal/config"
)

type fakeTokenCache struct {
	entry    *cachedAdminToken
	getErr   error
	setKey   string
	setVal   cachedAdminToken
	setTTL   time.Duration
	sets     int
	lockHeld bool
	tryLocks int
	unlocks  int
}

func (f *fakeTokenCache) Get(_ context.Context, _ string, dest interface{}) error {
	if f.getErr != nil {
		return f.getErr
	}
	*dest.(*cachedAdminToken) = *f.entry
	return nil
}

func (f *fakeTokenCache) Set(_ context.Context, key string, value interface{}, ttl time.Duration) error {
	f.sets++
	f.setKey = key
	f.setVal = value.(cachedAdminToken)
	f.setTTL = ttl
	return nil
}

func (f *fakeTokenCache) TryLock(_ context.Context, _ string, _ time.Duration) (string, error) {
	f.tryLocks++
	if f.lockHeld {
		return "", nil
	}
	return "test-token", nil
}

func (f *fakeTokenCache) Unlock(_ context.Context, _, _ string) error {
	f.unlocks++
	return nil
}

func TestGetTokens_WhenCachedSessionFresh_AdoptsIt(t *testing.T) {
	now := time.Now()
	cache := &fakeTokenCache{entry: &cachedAdminToken{
		AccessToken:      "cached-access",
		RefreshToken:     "cached-refresh",
		ExpiresAt:        now.Add(1 * time.Minute),
		RefreshExpiresAt: now.Add(10 * time.Minute),
	}}
	h := &AdminTokenHelper{tokenCache: cache}

	access, refresh, err := h.GetTokens()
	if err != nil {
		t.Fatalf("GetTokens: %v", err)
	}
	if access != "cached-access" || refresh != "cached-refresh" {
		t.Fatalf("got (%q, %q), want cached tokens", access, refresh)
	}
}

func TestGetTokens_WhenCacheFails_FallsThroughToKeycloak(t *testing.T) {
	h := &AdminTokenHelper{
		tokenCache:  &fakeTokenCache{getErr: errors.New("cache miss")},
		restyClient: resty.New().SetTimeout(200 * time.Millisecond),
		cfg: &config.Config{
			App:      &sharedconfig.AppConfig{},
			Keycloak: &sharedconfig.KeycloakConfig{Url: "http://127.0.0.1:1", Realm: "test"},
		},
	}
	if _, _, err := h.GetTokens(); err == nil {
		t.Fatal("expected the Keycloak path (error), got a cached session")
	}
}

func TestGetTokens_WhenLockAcquired_ReleasesItAfterTheRoundTrip(t *testing.T) {
	cache := &fakeTokenCache{getErr: errors.New("cache miss")}
	h := &AdminTokenHelper{
		tokenCache:  cache,
		restyClient: resty.New().SetTimeout(200 * time.Millisecond),
		cfg: &config.Config{
			App:      &sharedconfig.AppConfig{},
			Keycloak: &sharedconfig.KeycloakConfig{Url: "http://127.0.0.1:1", Realm: "test"},
		},
	}

	if _, _, err := h.GetTokens(); err == nil {
		t.Fatal("expected the Keycloak path (error)")
	}
	if cache.tryLocks != 1 {
		t.Fatalf("tryLocks = %d, want 1", cache.tryLocks)
	}
	if cache.unlocks != 1 {
		t.Fatalf("unlocks = %d, want 1 (lock must be released)", cache.unlocks)
	}
}

type flakyCache struct {
	fresh *cachedAdminToken
	stale *cachedAdminToken
	gets  int
}

func (f *flakyCache) Get(_ context.Context, _ string, dest interface{}) error {
	f.gets++
	if f.gets == 1 {
		*dest.(*cachedAdminToken) = *f.stale
	} else {
		*dest.(*cachedAdminToken) = *f.fresh
	}
	return nil
}
func (f *flakyCache) Set(context.Context, string, interface{}, time.Duration) error { return nil }
func (f *flakyCache) TryLock(context.Context, string, time.Duration) (string, error) {
	return "", nil
}
func (f *flakyCache) Unlock(context.Context, string, string) error { return nil }

func TestGetTokens_WhenLockHeldByAnotherReplica_AdoptsThePublishedSession(t *testing.T) {
	now := time.Now()
	cache := &flakyCache{
		stale: &cachedAdminToken{ExpiresAt: now.Add(-1 * time.Minute), RefreshExpiresAt: now.Add(-1 * time.Minute)},
		fresh: &cachedAdminToken{
			AccessToken:      "winner-access",
			RefreshToken:     "winner-refresh",
			ExpiresAt:        now.Add(1 * time.Minute),
			RefreshExpiresAt: now.Add(10 * time.Minute),
		},
	}
	h := &AdminTokenHelper{tokenCache: cache}

	access, refresh, err := h.GetTokens()
	if err != nil {
		t.Fatalf("GetTokens: %v", err)
	}
	if access != "winner-access" || refresh != "winner-refresh" {
		t.Fatalf("got (%q, %q), want the winner's published session", access, refresh)
	}
}

func TestStoreToCache_WhenSessionFresh_WritesItWithRefreshTtl(t *testing.T) {
	now := time.Now()
	cache := &fakeTokenCache{}
	h := &AdminTokenHelper{
		tokenCache:       cache,
		accessToken:      "a",
		refreshToken:     "r",
		expiresAt:        now.Add(1 * time.Minute),
		refreshExpiresAt: now.Add(10 * time.Minute),
	}

	h.storeToCache()

	if cache.sets != 1 {
		t.Fatalf("sets = %d, want 1", cache.sets)
	}
	if cache.setKey != adminTokenCacheKey {
		t.Fatalf("key = %q, want %q", cache.setKey, adminTokenCacheKey)
	}
	if cache.setVal.AccessToken != "a" || cache.setVal.RefreshToken != "r" {
		t.Fatalf("cached value = %+v, want in-memory session", cache.setVal)
	}
	if cache.setTTL <= 9*time.Minute || cache.setTTL > 10*time.Minute {
		t.Fatalf("ttl = %v, want ~refresh expiry (10m)", cache.setTTL)
	}
}

func TestFresh_WhenExpiryWithinSkew_ReturnsFalse(t *testing.T) {
	now := time.Now()
	cases := map[string]struct {
		expiresAt time.Time
		want      bool
	}{
		"whenExpiryComfortablyAhead_returnsTrue": {now.Add(2 * expirySkew), true},
		"whenExpiryWithinSkew_returnsFalse":      {now.Add(expirySkew / 2), false},
		"whenAlreadyExpired_returnsFalse":        {now.Add(-time.Second), false},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			if got := fresh(tc.expiresAt); got != tc.want {
				t.Errorf("fresh(%v) = %v, want %v", tc.expiresAt, got, tc.want)
			}
		})
	}
}

func TestGetTokens_WhenCachedTokenNearExpiry_FallsThroughToKeycloak(t *testing.T) {
	now := time.Now()
	cache := &fakeTokenCache{entry: &cachedAdminToken{
		AccessToken:      "near-expiry-access",
		RefreshToken:     "",
		ExpiresAt:        now.Add(expirySkew / 2),
		RefreshExpiresAt: now.Add(-time.Minute),
	}}
	h := &AdminTokenHelper{
		tokenCache:  cache,
		restyClient: resty.New().SetTimeout(200 * time.Millisecond),
		cfg: &config.Config{
			App:      &sharedconfig.AppConfig{},
			Keycloak: &sharedconfig.KeycloakConfig{Url: "http://127.0.0.1:1", Realm: "test"},
		},
	}
	if _, _, err := h.GetTokens(); err == nil {
		t.Fatal("expected a Keycloak round-trip (error); a near-expiry cached token must not be adopted")
	}
}

func TestStoreToCache_WhenSessionExpired_SkipsTheWrite(t *testing.T) {
	cache := &fakeTokenCache{}
	h := &AdminTokenHelper{
		tokenCache:       cache,
		refreshExpiresAt: time.Now().Add(-1 * time.Minute),
	}

	h.storeToCache()

	if cache.sets != 0 {
		t.Fatalf("sets = %d, want 0 for an already-expired session", cache.sets)
	}
}
