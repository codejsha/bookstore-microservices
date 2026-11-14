package keycloak

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/codejsha/shared-library-go/pkg/rest/client"
)

const (
	tokenPath      = "/realms/test/protocol/openid-connect/token"
	tokenGrantBody = `{"access_token":"new-access","expires_in":300,"token_type":"Bearer"}`
)

type scriptedTokenCache struct {
	entry   *cachedAdminToken
	lockErr error
	stored  []cachedAdminToken
	ttls    []time.Duration
	locks   int
	unlocks []string
}

func (s *scriptedTokenCache) Get(_ context.Context, _ string, dest interface{}) error {
	if s.entry == nil {
		return errors.New("cache miss")
	}
	*dest.(*cachedAdminToken) = *s.entry
	return nil
}

func (s *scriptedTokenCache) Set(_ context.Context, _ string, value interface{}, ttl time.Duration) error {
	s.stored = append(s.stored, value.(cachedAdminToken))
	s.ttls = append(s.ttls, ttl)
	return nil
}

func (s *scriptedTokenCache) TryLock(context.Context, string, time.Duration) (string, error) {
	s.locks++
	if s.lockErr != nil {
		return "", s.lockErr
	}
	return "lock-owner", nil
}

func (s *scriptedTokenCache) Unlock(_ context.Context, _, token string) error {
	s.unlocks = append(s.unlocks, token)
	return nil
}

func TestAdminTokenHelperAccessToken_noSession_requestsClientCredentialsGrantAndPublishesIt(t *testing.T) {
	stub := newKeycloakStub(t).on(http.MethodPost, tokenPath, http.StatusOK, tokenGrantBody)
	cache := &scriptedTokenCache{}
	h := NewAdminTokenHelper(stubConfig(stub.URL), &client.RestyClient{Client: newStubRestyClient()}, cache)

	access, err := h.AccessToken()
	if err != nil {
		t.Fatalf("access token: %v", err)
	}

	if access != "new-access" {
		t.Fatalf("access = %q, want the granted token", access)
	}
	req := stub.only(t)
	if req.Method != http.MethodPost || req.Path != tokenPath {
		t.Errorf("request = %s %s, want POST %s", req.Method, req.Path, tokenPath)
	}
	if !strings.HasPrefix(req.ContentType, "application/x-www-form-urlencoded") {
		t.Errorf("content type = %q, want a form post", req.ContentType)
	}
	form := req.form(t)
	want := map[string]string{
		"grant_type":    "client_credentials",
		"client_id":     stubClientID,
		"client_secret": stubClientSecret,
	}
	for key, value := range want {
		if form.Get(key) != value {
			t.Errorf("form[%s] = %q, want %q", key, form.Get(key), value)
		}
	}
	for _, key := range []string{"username", "password", "refresh_token"} {
		if form.Has(key) {
			t.Errorf("form = %v, want no %s on a client credentials grant", form, key)
		}
	}
	if len(cache.stored) != 1 || cache.stored[0].AccessToken != "new-access" {
		t.Fatalf("published tokens = %+v, want the granted token once", cache.stored)
	}
	if cache.ttls[0] <= 4*time.Minute || cache.ttls[0] > 5*time.Minute {
		t.Errorf("published ttl = %v, want ~expires_in (5m)", cache.ttls[0])
	}
	if lifetime := time.Until(cache.stored[0].ExpiresAt); lifetime <= 4*time.Minute || lifetime > 5*time.Minute {
		t.Errorf("published expiry in %v, want ~expires_in (5m)", lifetime)
	}
	if len(cache.unlocks) != 1 || cache.unlocks[0] != "lock-owner" {
		t.Errorf("unlocks = %v, want the acquired lock released by its owner token", cache.unlocks)
	}
}

func TestAdminTokenHelperAccessToken_tokenStillFresh_reusesItWithoutAnotherGrant(t *testing.T) {
	stub := newKeycloakStub(t).on(http.MethodPost, tokenPath, http.StatusOK, tokenGrantBody)
	h := NewAdminTokenHelper(stubConfig(stub.URL), &client.RestyClient{Client: newStubRestyClient()}, nil)

	if _, err := h.AccessToken(); err != nil {
		t.Fatalf("first access token: %v", err)
	}
	access, err := h.AccessToken()
	if err != nil {
		t.Fatalf("second access token: %v", err)
	}

	if access != "new-access" {
		t.Fatalf("access = %q, want the granted token", access)
	}
	if requests := stub.recorded(); len(requests) != 1 {
		t.Fatalf("token requests = %d, want 1", len(requests))
	}
}

func TestAdminTokenHelperAccessToken_tokenExpired_requestsANewClientCredentialsGrant(t *testing.T) {
	stub := newKeycloakStub(t).on(http.MethodPost, tokenPath, http.StatusOK, tokenGrantBody)
	h := &AdminTokenHelper{
		cfg:         stubConfig(stub.URL),
		restyClient: newStubRestyClient(),
		accessToken: "old-access",
		expiresAt:   time.Now().Add(-time.Minute),
	}

	access, err := h.AccessToken()
	if err != nil {
		t.Fatalf("access token: %v", err)
	}

	if access != "new-access" {
		t.Fatalf("access = %q, want the new token", access)
	}
	if got := stub.only(t).form(t).Get("grant_type"); got != "client_credentials" {
		t.Errorf("grant_type = %q, want client_credentials", got)
	}
}

func TestAdminTokenHelperAccessToken_cachedTokenExpired_requestsANewGrant(t *testing.T) {
	stub := newKeycloakStub(t).on(http.MethodPost, tokenPath, http.StatusOK, tokenGrantBody)
	cache := &scriptedTokenCache{entry: &cachedAdminToken{
		AccessToken: "peer-access",
		ExpiresAt:   time.Now().Add(-time.Minute),
	}}
	h := &AdminTokenHelper{cfg: stubConfig(stub.URL), restyClient: newStubRestyClient(), tokenCache: cache}

	access, err := h.AccessToken()
	if err != nil {
		t.Fatalf("access token: %v", err)
	}

	if access != "new-access" {
		t.Fatalf("access = %q, want the new token instead of the expired cached one", access)
	}
	if got := len(stub.recorded()); got != 1 {
		t.Errorf("token requests = %d, want 1", got)
	}
}

func TestAdminTokenHelperAccessToken_lockBackendFails_requestsWithoutHoldingALock(t *testing.T) {
	stub := newKeycloakStub(t).on(http.MethodPost, tokenPath, http.StatusOK, tokenGrantBody)
	cache := &scriptedTokenCache{lockErr: errors.New("valkey down")}
	h := &AdminTokenHelper{cfg: stubConfig(stub.URL), restyClient: newStubRestyClient(), tokenCache: cache}

	access, err := h.AccessToken()
	if err != nil {
		t.Fatalf("access token: %v", err)
	}

	if access != "new-access" {
		t.Fatalf("access = %q, want new-access", access)
	}
	if cache.locks != 1 {
		t.Errorf("lock attempts = %d, want 1", cache.locks)
	}
	if len(cache.unlocks) != 0 {
		t.Errorf("unlocks = %v, want none for a lock that was never acquired", cache.unlocks)
	}
}

func TestAdminTokenHelperAccessToken_malformedGrantBody_reportsErrorAndKeepsOldToken(t *testing.T) {
	stub := newKeycloakStub(t).on(http.MethodPost, tokenPath, http.StatusOK, `<html>gateway</html>`)
	expired := time.Now().Add(-time.Minute)
	h := &AdminTokenHelper{
		cfg:         stubConfig(stub.URL),
		restyClient: newStubRestyClient(),
		accessToken: "old-access",
		expiresAt:   expired,
	}

	access, err := h.AccessToken()

	if err == nil {
		t.Fatal("err = nil, want a decode error")
	}
	if access != "" {
		t.Errorf("access = %q, want none", access)
	}
	if h.accessToken != "old-access" || !h.expiresAt.Equal(expired) {
		t.Errorf("token = (%q, %v), want the previous one untouched", h.accessToken, h.expiresAt)
	}
}

func TestAdminTokenHelperAccessToken_keycloakUnreachable_reportsError(t *testing.T) {
	h := &AdminTokenHelper{cfg: stubConfig(unreachableKeycloakURL(t)), restyClient: newStubRestyClient()}

	_, err := h.AccessToken()

	if err == nil || !strings.Contains(err.Error(), "failed to refresh token") {
		t.Fatalf("err = %v, want the refresh failure reported", err)
	}
}

func TestAdminTokenHelperAccessToken_grantRejected_reportsStatusAndOAuthErrorCode(t *testing.T) {
	stub := newKeycloakStub(t).on(http.MethodPost, tokenPath, http.StatusUnauthorized,
		`{"error":"unauthorized_client","error_description":"Invalid client or Invalid client credentials"}`)
	h := &AdminTokenHelper{cfg: stubConfig(stub.URL), restyClient: newStubRestyClient()}

	access, err := h.AccessToken()

	if err == nil {
		t.Fatal("err = nil, want the rejected grant reported")
	}
	for _, want := range []string{"401", "unauthorized_client", "Invalid client credentials"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("err = %v, want it to mention %q", err, want)
		}
	}
	if access != "" {
		t.Errorf("access = %q, want none", access)
	}
	if h.accessToken != "" || !h.expiresAt.IsZero() {
		t.Errorf("token = (%q, %v), want it left empty", h.accessToken, h.expiresAt)
	}
	if got := len(stub.recorded()); got != 1 {
		t.Errorf("token requests = %d, want a single grant attempt", got)
	}
}

func TestAdminTokenHelperAccessToken_successBodyWithoutAccessToken_reportsError(t *testing.T) {
	stub := newKeycloakStub(t).on(http.MethodPost, tokenPath, http.StatusOK, `{}`)
	h := &AdminTokenHelper{cfg: stubConfig(stub.URL), restyClient: newStubRestyClient()}

	_, err := h.AccessToken()

	if err == nil || !strings.Contains(err.Error(), "no access token") {
		t.Fatalf("err = %v, want the missing access token reported", err)
	}
}

type capturingRestyLogger struct {
	mu    sync.Mutex
	lines []string
}

func (l *capturingRestyLogger) add(format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.lines = append(l.lines, fmt.Sprintf(format, v...))
}

func (l *capturingRestyLogger) Errorf(format string, v ...interface{}) { l.add(format, v...) }
func (l *capturingRestyLogger) Warnf(format string, v ...interface{})  { l.add(format, v...) }
func (l *capturingRestyLogger) Debugf(format string, v ...interface{}) { l.add(format, v...) }

func (l *capturingRestyLogger) output() string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return strings.Join(l.lines, "\n")
}

func TestAdminTokenHelper_debugLoggingEnabled_keepsCredentialsAndTokensOutOfTheLog(t *testing.T) {
	stub := newKeycloakStub(t).
		on(http.MethodPost, tokenPath, http.StatusOK, tokenGrantBody).
		on(http.MethodGet, userPath, http.StatusOK, `{"id":"idp-1"}`)
	logs := &capturingRestyLogger{}
	shared := &client.RestyClient{Client: resty.New().SetTimeout(2 * time.Second).SetLogger(logs)}
	cfg := stubConfig(stub.URL)
	cfg.App.Logging.IsDebugEnabled = true
	h := NewAdminTokenHelper(cfg, shared, nil)
	users := NewUsersClient(cfg, shared, h)

	if _, err := users.GetUser(context.Background(), stubRealm, "idp-1"); err != nil {
		t.Fatalf("get user: %v", err)
	}

	out := logs.output()
	if !strings.Contains(out, tokenPath) || !strings.Contains(out, userPath) {
		t.Fatalf("debug log = %q, want both round trips logged", out)
	}
	for _, secret := range []string{stubClientSecret, "new-access"} {
		if strings.Contains(out, secret) {
			t.Errorf("debug log leaks %q", secret)
		}
	}
	if !strings.Contains(out, redactedValue) {
		t.Errorf("debug log = %q, want redaction markers", out)
	}
}

func TestRedactBody_secretBearingPayload_masksOnlySecretValues(t *testing.T) {
	cases := map[string]struct {
		body string
		want string
	}{
		"emptyBody_unchanged": {"", ""},
		"formGrant_masksCredentials": {
			"client_id=identity&client_secret=abc&grant_type=password&password=p%40ss&username=admin",
			"client_id=identity&client_secret=[REDACTED]&grant_type=password&password=[REDACTED]&username=admin",
		},
		"formRefresh_masksRefreshToken": {
			"grant_type=refresh_token&refresh_token=eyJ.a.b",
			"grant_type=refresh_token&refresh_token=[REDACTED]",
		},
		"jsonTokenResponse_masksTokens": {
			`{"access_token":"eyJ.a.b","expires_in":300,"refresh_token": "eyJ.c.d"}`,
			`{"access_token":"[REDACTED]","expires_in":300,"refresh_token":"[REDACTED]"}`,
		},
		"jsonCredential_masksValue": {
			`{"username":"reader","credentials":[{"type":"password","value":"hunter2"}]}`,
			`{"username":"reader","credentials":[{"type":"password","value":"[REDACTED]"}]}`,
		},
		"noSecrets_unchanged": {
			`{"username":"reader","email":"reader@example.com"}`,
			`{"username":"reader","email":"reader@example.com"}`,
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			if got := redactBody(tc.body); got != tc.want {
				t.Errorf("redactBody(%q) = %q, want %q", tc.body, got, tc.want)
			}
		})
	}
}

func TestRedactHeaders_credentialHeadersPresent_masksThemAndKeepsTheRest(t *testing.T) {
	h := http.Header{}
	h.Set("Authorization", "Bearer eyJ.a.b")
	h.Add("Set-Cookie", "KEYCLOAK_SESSION=abc")
	h.Add("Set-Cookie", "AUTH_SESSION_ID=def")
	h.Set("Content-Type", "application/json")

	redactHeaders(h)

	if got := h.Values("Authorization"); len(got) != 1 || got[0] != redactedValue {
		t.Errorf("authorization = %v, want redacted", got)
	}
	if got := h.Values("Set-Cookie"); len(got) != 1 || got[0] != redactedValue {
		t.Errorf("set-cookie = %v, want every cookie collapsed into one redacted value", got)
	}
	if got := h.Get("Content-Type"); got != "application/json" {
		t.Errorf("content type = %q, want it untouched", got)
	}
	if _, present := h["Cookie"]; present {
		t.Errorf("headers = %v, want no Cookie header invented", h)
	}
}
