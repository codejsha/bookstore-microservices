package keycloak

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/codejsha/shared-library-go/pkg/rest/client"

	"github.com/codejsha/bookstore-microservices/identity/internal/domain/constant"
)

const (
	introspectPath = "/realms/test/protocol/openid-connect/token/introspect"
	opaqueToken    = "opaque-user-token"
)

func newStubIntrospector(baseURL string) Introspector {
	return NewIntrospector(stubConfig(baseURL), &client.RestyClient{Client: newStubRestyClient()})
}

func TestIntrospector_activeToken_postsClientCredentialsAndParsesClaims(t *testing.T) {
	stub := newKeycloakStub(t).on(http.MethodPost, introspectPath, http.StatusOK,
		`{"active":true,"sub":"7c9e6679-7425-40de-944b-e07fc1f90ae7","exp":1900000000,"iat":1899999700,"scope":"openid"}`)

	got, err := newStubIntrospector(stub.URL).Introspect(context.Background(), opaqueToken)
	if err != nil {
		t.Fatalf("introspect: %v", err)
	}

	req := stub.only(t)
	if req.Method != http.MethodPost || req.Path != introspectPath {
		t.Errorf("request = %s %s, want POST %s", req.Method, req.Path, introspectPath)
	}
	if !strings.HasPrefix(req.ContentType, "application/x-www-form-urlencoded") {
		t.Errorf("content type = %q, want a form post", req.ContentType)
	}
	form := req.form(t)
	if form.Get("token") != opaqueToken || form.Get("client_id") != stubClientID || form.Get("client_secret") != stubClientSecret {
		t.Errorf("form = %v, want the token with the client credentials", form)
	}
	want := Introspection{Active: true, Sub: "7c9e6679-7425-40de-944b-e07fc1f90ae7", Exp: 1900000000, Iat: 1899999700}
	if got != want {
		t.Fatalf("introspection = %+v, want %+v", got, want)
	}
}

func TestIntrospector_inactiveToken_reportsInactiveWithoutError(t *testing.T) {
	stub := newKeycloakStub(t).on(http.MethodPost, introspectPath, http.StatusOK, `{"active":false}`)

	got, err := newStubIntrospector(stub.URL).Introspect(context.Background(), opaqueToken)
	if err != nil {
		t.Fatalf("introspect: %v", err)
	}

	if got != (Introspection{}) {
		t.Fatalf("introspection = %+v, want inactive with no claims", got)
	}
}

func TestIntrospector_keycloak401_reportsStatusError(t *testing.T) {
	stub := newKeycloakStub(t).on(http.MethodPost, introspectPath, http.StatusUnauthorized,
		`{"error":"invalid_client","active":true}`)

	got, err := newStubIntrospector(stub.URL).Introspect(context.Background(), opaqueToken)

	if err == nil || !strings.Contains(err.Error(), "401") {
		t.Fatalf("err = %v, want the 401 status reported", err)
	}
	if got.Active {
		t.Errorf("introspection = %+v, want a rejected call never read as active", got)
	}
}

func TestIntrospector_malformedBody_reportsDecodeError(t *testing.T) {
	stub := newKeycloakStub(t).on(http.MethodPost, introspectPath, http.StatusOK, `<html>gateway</html>`)

	got, err := newStubIntrospector(stub.URL).Introspect(context.Background(), opaqueToken)

	if err == nil {
		t.Fatal("err = nil, want a decode error")
	}
	if got.Active {
		t.Errorf("introspection = %+v, want inactive", got)
	}
}

func TestIntrospector_keycloakUnreachable_reportsTransportError(t *testing.T) {
	got, err := newStubIntrospector(unreachableKeycloakURL(t)).Introspect(context.Background(), opaqueToken)

	if err == nil {
		t.Fatal("err = nil, want a transport error")
	}
	if got.Active {
		t.Errorf("introspection = %+v, want inactive", got)
	}
}

type fakeIntrospectionCache struct {
	entries map[string]Introspection
	ttls    map[string]time.Duration
	setErr  error
}

func newFakeIntrospectionCache() *fakeIntrospectionCache {
	return &fakeIntrospectionCache{entries: map[string]Introspection{}, ttls: map[string]time.Duration{}}
}

func (f *fakeIntrospectionCache) Get(_ context.Context, key string, dest interface{}) error {
	entry, ok := f.entries[key]
	if !ok {
		return errors.New("cache miss")
	}
	*dest.(*Introspection) = entry
	return nil
}

func (f *fakeIntrospectionCache) Set(_ context.Context, key string, value interface{}, ttl time.Duration) error {
	if f.setErr != nil {
		return f.setErr
	}
	f.entries[key] = value.(Introspection)
	f.ttls[key] = ttl
	return nil
}

func (f *fakeIntrospectionCache) TryLock(context.Context, string, time.Duration) (string, error) {
	return "", nil
}

func (f *fakeIntrospectionCache) Unlock(context.Context, string, string) error { return nil }

type countingIntrospector struct {
	result Introspection
	err    error
	tokens []string
}

func (c *countingIntrospector) Introspect(_ context.Context, token string) (Introspection, error) {
	c.tokens = append(c.tokens, token)
	return c.result, c.err
}

func TestCachingIntrospector_cacheMiss_asksKeycloakAndStoresUnderHashedKey(t *testing.T) {
	inner := &countingIntrospector{result: Introspection{Active: true, Sub: "sub-1"}}
	cache := newFakeIntrospectionCache()

	got, err := NewCachingIntrospector(inner, cache).Introspect(context.Background(), opaqueToken)
	if err != nil {
		t.Fatalf("introspect: %v", err)
	}

	if got != inner.result {
		t.Fatalf("introspection = %+v, want %+v", got, inner.result)
	}
	if len(inner.tokens) != 1 || inner.tokens[0] != opaqueToken {
		t.Fatalf("inner calls = %v, want one for the token", inner.tokens)
	}
	if len(cache.entries) != 1 {
		t.Fatalf("cache entries = %v, want 1", cache.entries)
	}
	for key, entry := range cache.entries {
		if strings.Contains(key, opaqueToken) {
			t.Errorf("cache key = %q, want the raw token kept out of it", key)
		}
		if !strings.HasPrefix(key, "identity:introspect:") {
			t.Errorf("cache key = %q, want the identity:introspect: namespace", key)
		}
		if entry != inner.result {
			t.Errorf("cached = %+v, want %+v", entry, inner.result)
		}
		if cache.ttls[key] != constant.IntrospectCacheTTL {
			t.Errorf("ttl = %v, want %v", cache.ttls[key], constant.IntrospectCacheTTL)
		}
	}
}

func TestCachingIntrospector_sameTokenTwice_asksKeycloakOnce(t *testing.T) {
	inner := &countingIntrospector{result: Introspection{Active: true, Sub: "sub-1"}}
	c := NewCachingIntrospector(inner, newFakeIntrospectionCache())

	if _, err := c.Introspect(context.Background(), opaqueToken); err != nil {
		t.Fatalf("first introspect: %v", err)
	}
	got, err := c.Introspect(context.Background(), opaqueToken)
	if err != nil {
		t.Fatalf("second introspect: %v", err)
	}

	if got != inner.result {
		t.Fatalf("introspection = %+v, want the cached %+v", got, inner.result)
	}
	if len(inner.tokens) != 1 {
		t.Fatalf("inner calls = %d, want 1", len(inner.tokens))
	}
}

func TestCachingIntrospector_differentTokens_cachedSeparately(t *testing.T) {
	inner := &countingIntrospector{result: Introspection{Active: true, Sub: "sub-1"}}
	cache := newFakeIntrospectionCache()
	c := NewCachingIntrospector(inner, cache)

	for _, token := range []string{"token-a", "token-b"} {
		if _, err := c.Introspect(context.Background(), token); err != nil {
			t.Fatalf("introspect %s: %v", token, err)
		}
	}

	if len(inner.tokens) != 2 || len(cache.entries) != 2 {
		t.Fatalf("inner calls = %d, cache entries = %d, want 2 and 2", len(inner.tokens), len(cache.entries))
	}
}

func TestCachingIntrospector_inactiveToken_cachedSoRevokedTokenIsNotReasked(t *testing.T) {
	inner := &countingIntrospector{result: Introspection{Active: false}}
	c := NewCachingIntrospector(inner, newFakeIntrospectionCache())

	for i := 0; i < 2; i++ {
		got, err := c.Introspect(context.Background(), opaqueToken)
		if err != nil {
			t.Fatalf("introspect: %v", err)
		}
		if got.Active {
			t.Fatalf("introspection = %+v, want inactive", got)
		}
	}

	if len(inner.tokens) != 1 {
		t.Fatalf("inner calls = %d, want 1", len(inner.tokens))
	}
}

func TestCachingIntrospector_keycloakFails_propagatesErrorAndCachesNothing(t *testing.T) {
	failure := errors.New("keycloak down")
	inner := &countingIntrospector{err: failure}
	cache := newFakeIntrospectionCache()

	_, err := NewCachingIntrospector(inner, cache).Introspect(context.Background(), opaqueToken)

	if !errors.Is(err, failure) {
		t.Fatalf("err = %v, want %v", err, failure)
	}
	if len(cache.entries) != 0 {
		t.Fatalf("cache entries = %v, want none", cache.entries)
	}
}

func TestCachingIntrospector_cacheWriteFails_stillAnswers(t *testing.T) {
	inner := &countingIntrospector{result: Introspection{Active: true, Sub: "sub-1"}}
	cache := newFakeIntrospectionCache()
	cache.setErr = errors.New("valkey down")

	got, err := NewCachingIntrospector(inner, cache).Introspect(context.Background(), opaqueToken)
	if err != nil {
		t.Fatalf("introspect: %v", err)
	}

	if got != inner.result {
		t.Fatalf("introspection = %+v, want %+v", got, inner.result)
	}
}

func TestCacheTTLFor_tokenLifetime_capsTtlAtRemainingValidity(t *testing.T) {
	now := time.Now()
	cases := map[string]struct {
		result  Introspection
		wantMin time.Duration
		wantMax time.Duration
	}{
		"inactiveToken_fullTtl": {
			Introspection{Active: false, Exp: now.Add(2 * time.Second).Unix()},
			constant.IntrospectCacheTTL, constant.IntrospectCacheTTL,
		},
		"activeWithoutExpiry_fullTtl": {
			Introspection{Active: true},
			constant.IntrospectCacheTTL, constant.IntrospectCacheTTL,
		},
		"activeExpiringLater_fullTtl": {
			Introspection{Active: true, Exp: now.Add(time.Hour).Unix()},
			constant.IntrospectCacheTTL, constant.IntrospectCacheTTL,
		},
		"activeExpiringSoon_remainingValidity": {
			Introspection{Active: true, Exp: now.Add(4 * time.Second).Unix()},
			time.Second, 4 * time.Second,
		},
		"activeAlreadyExpired_oneSecondFloor": {
			Introspection{Active: true, Exp: now.Add(-time.Minute).Unix()},
			time.Second, time.Second,
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := cacheTTLFor(tc.result)

			if got < tc.wantMin || got > tc.wantMax {
				t.Fatalf("ttl = %v, want within [%v, %v]", got, tc.wantMin, tc.wantMax)
			}
		})
	}
}
