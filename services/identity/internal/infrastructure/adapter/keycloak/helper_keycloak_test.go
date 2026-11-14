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

func TestAdminTokenHelperAccessToken_cachedTokenFresh_adoptsIt(t *testing.T) {
	cache := &fakeTokenCache{entry: &cachedAdminToken{
		AccessToken: "cached-access",
		ExpiresAt:   time.Now().Add(1 * time.Minute),
	}}
	h := &AdminTokenHelper{tokenCache: cache}

	access, err := h.AccessToken()
	if err != nil {
		t.Fatalf("access token: %v", err)
	}
	if access != "cached-access" {
		t.Fatalf("access = %q, want the cached token", access)
	}
}

func unreachableHelper(cache TokenCache) *AdminTokenHelper {
	return &AdminTokenHelper{
		tokenCache:  cache,
		restyClient: resty.New().SetTimeout(200 * time.Millisecond),
		cfg: &config.Config{
			App:      &sharedconfig.AppConfig{},
			Keycloak: &sharedconfig.KeycloakConfig{Url: "http://127.0.0.1:1", Realm: "test"},
		},
	}
}

func TestAdminTokenHelperAccessToken_cacheFails_fallsThroughToKeycloak(t *testing.T) {
	h := unreachableHelper(&fakeTokenCache{getErr: errors.New("cache miss")})
	if _, err := h.AccessToken(); err == nil {
		t.Fatal("expected the Keycloak path (error), got a cached token")
	}
}

func TestAdminTokenHelperAccessToken_lockAcquired_releasesItAfterTheRoundTrip(t *testing.T) {
	cache := &fakeTokenCache{getErr: errors.New("cache miss")}
	h := unreachableHelper(cache)

	if _, err := h.AccessToken(); err == nil {
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

func TestAdminTokenHelperAccessToken_lockHeldByAnotherReplica_adoptsThePublishedToken(t *testing.T) {
	now := time.Now()
	cache := &flakyCache{
		stale: &cachedAdminToken{ExpiresAt: now.Add(-1 * time.Minute)},
		fresh: &cachedAdminToken{AccessToken: "winner-access", ExpiresAt: now.Add(1 * time.Minute)},
	}
	h := &AdminTokenHelper{tokenCache: cache}

	access, err := h.AccessToken()
	if err != nil {
		t.Fatalf("access token: %v", err)
	}
	if access != "winner-access" {
		t.Fatalf("access = %q, want the winner's published token", access)
	}
}

func TestStoreToCache_tokenFresh_writesItWithAccessTtl(t *testing.T) {
	cache := &fakeTokenCache{}
	h := &AdminTokenHelper{
		tokenCache:  cache,
		accessToken: "a",
		expiresAt:   time.Now().Add(5 * time.Minute),
	}

	h.storeToCache()

	if cache.sets != 1 {
		t.Fatalf("sets = %d, want 1", cache.sets)
	}
	if cache.setKey != adminTokenCacheKey {
		t.Fatalf("key = %q, want %q", cache.setKey, adminTokenCacheKey)
	}
	if cache.setVal.AccessToken != "a" {
		t.Fatalf("cached value = %+v, want the in-memory token", cache.setVal)
	}
	if cache.setTTL <= 4*time.Minute || cache.setTTL > 5*time.Minute {
		t.Fatalf("ttl = %v, want ~access expiry (5m)", cache.setTTL)
	}
}

func TestFresh_expiryWithinSkew_isNotFresh(t *testing.T) {
	now := time.Now()
	cases := map[string]struct {
		expiresAt time.Time
		want      bool
	}{
		"expiryComfortablyAhead_isFresh": {now.Add(2 * expirySkew), true},
		"expiryWithinSkew_isNotFresh":    {now.Add(expirySkew / 2), false},
		"alreadyExpired_isNotFresh":      {now.Add(-time.Second), false},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			if got := fresh(tc.expiresAt); got != tc.want {
				t.Errorf("fresh(%v) = %v, want %v", tc.expiresAt, got, tc.want)
			}
		})
	}
}

func TestAdminTokenHelperAccessToken_cachedTokenNearExpiry_fallsThroughToKeycloak(t *testing.T) {
	cache := &fakeTokenCache{entry: &cachedAdminToken{
		AccessToken: "near-expiry-access",
		ExpiresAt:   time.Now().Add(expirySkew / 2),
	}}
	if _, err := unreachableHelper(cache).AccessToken(); err == nil {
		t.Fatal("expected a Keycloak round-trip (error); a near-expiry cached token must not be adopted")
	}
}

func TestStoreToCache_tokenExpired_skipsTheWrite(t *testing.T) {
	cache := &fakeTokenCache{}
	h := &AdminTokenHelper{tokenCache: cache, expiresAt: time.Now().Add(-1 * time.Minute)}

	h.storeToCache()

	if cache.sets != 0 {
		t.Fatalf("sets = %d, want 0 for an already-expired token", cache.sets)
	}
}
