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
	tokenPath        = "/realms/test/protocol/openid-connect/token"
	tokenGrantBody   = `{"access_token":"new-access","refresh_token":"new-refresh","expires_in":300,"refresh_expires_in":1800}`
	wantGrantedScope = "openid profile email"
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

func TestAdminTokenHelperGetTokens_noSession_requestsPasswordGrantAndPublishesSession(t *testing.T) {
	stub := newKeycloakStub(t).on(http.MethodPost, tokenPath, http.StatusOK, tokenGrantBody)
	cache := &scriptedTokenCache{}
	h := NewAdminTokenHelper(stubConfig(stub.URL), &client.RestyClient{Client: newStubRestyClient()}, cache)

	access, refresh, err := h.GetTokens()
	if err != nil {
		t.Fatalf("get tokens: %v", err)
	}

	if access != "new-access" || refresh != "new-refresh" {
		t.Fatalf("tokens = (%q, %q), want the granted pair", access, refresh)
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
		"grant_type":    "password",
		"client_id":     stubClientID,
		"client_secret": stubClientSecret,
		"username":      stubAdminUser,
		"password":      stubAdminPass,
		"scope":         wantGrantedScope,
	}
	for key, value := range want {
		if form.Get(key) != value {
			t.Errorf("form[%s] = %q, want %q", key, form.Get(key), value)
		}
	}
	if form.Has("refresh_token") {
		t.Errorf("form = %v, want no refresh_token on a password grant", form)
	}
	if len(cache.stored) != 1 {
		t.Fatalf("published sessions = %+v, want 1", cache.stored)
	}
	if cache.stored[0].AccessToken != "new-access" || cache.stored[0].RefreshToken != "new-refresh" {
		t.Errorf("published session = %+v, want the granted pair", cache.stored[0])
	}
	if cache.ttls[0] <= 29*time.Minute || cache.ttls[0] > 30*time.Minute {
		t.Errorf("published ttl = %v, want ~refresh_expires_in (30m)", cache.ttls[0])
	}
	if lifetime := time.Until(cache.stored[0].ExpiresAt); lifetime <= 4*time.Minute || lifetime > 5*time.Minute {
		t.Errorf("published access expiry in %v, want ~expires_in (5m)", lifetime)
	}
	if len(cache.unlocks) != 1 || cache.unlocks[0] != "lock-owner" {
		t.Errorf("unlocks = %v, want the acquired lock released by its owner token", cache.unlocks)
	}
}

func TestAdminTokenHelperGetTokens_sessionStillFresh_reusesItWithoutAnotherGrant(t *testing.T) {
	stub := newKeycloakStub(t).on(http.MethodPost, tokenPath, http.StatusOK, tokenGrantBody)
	h := NewAdminTokenHelper(stubConfig(stub.URL), &client.RestyClient{Client: newStubRestyClient()}, nil)

	if _, _, err := h.GetTokens(); err != nil {
		t.Fatalf("first get tokens: %v", err)
	}
	access, refresh, err := h.GetTokens()
	if err != nil {
		t.Fatalf("second get tokens: %v", err)
	}

	if access != "new-access" || refresh != "new-refresh" {
		t.Fatalf("tokens = (%q, %q), want the granted pair", access, refresh)
	}
	if requests := stub.recorded(); len(requests) != 1 {
		t.Fatalf("token requests = %d, want 1", len(requests))
	}
}

func TestAdminTokenHelperGetTokens_accessExpiredRefreshValid_requestsRefreshGrant(t *testing.T) {
	stub := newKeycloakStub(t).on(http.MethodPost, tokenPath, http.StatusOK, tokenGrantBody)
	h := &AdminTokenHelper{
		cfg:              stubConfig(stub.URL),
		restyClient:      newStubRestyClient(),
		accessToken:      "old-access",
		refreshToken:     "old-refresh",
		expiresAt:        time.Now().Add(-time.Minute),
		refreshExpiresAt: time.Now().Add(10 * time.Minute),
	}

	access, refresh, err := h.GetTokens()
	if err != nil {
		t.Fatalf("get tokens: %v", err)
	}

	if access != "new-access" || refresh != "new-refresh" {
		t.Fatalf("tokens = (%q, %q), want the rotated pair", access, refresh)
	}
	req := stub.only(t)
	if req.Path != tokenPath {
		t.Errorf("path = %s, want %s", req.Path, tokenPath)
	}
	form := req.form(t)
	want := map[string]string{
		"grant_type":    "refresh_token",
		"refresh_token": "old-refresh",
		"client_id":     stubClientID,
		"client_secret": stubClientSecret,
		"scope":         wantGrantedScope,
	}
	for key, value := range want {
		if form.Get(key) != value {
			t.Errorf("form[%s] = %q, want %q", key, form.Get(key), value)
		}
	}
	if form.Has("username") || form.Has("password") {
		t.Errorf("form = %v, want no admin credentials on a refresh grant", form)
	}
}

func TestAdminTokenHelperGetTokens_cachedAccessExpiredRefreshValid_exchangesTheCachedRefreshToken(t *testing.T) {
	stub := newKeycloakStub(t).on(http.MethodPost, tokenPath, http.StatusOK, tokenGrantBody)
	cache := &scriptedTokenCache{entry: &cachedAdminToken{
		AccessToken:      "peer-access",
		RefreshToken:     "peer-refresh",
		ExpiresAt:        time.Now().Add(-time.Minute),
		RefreshExpiresAt: time.Now().Add(10 * time.Minute),
	}}
	h := &AdminTokenHelper{cfg: stubConfig(stub.URL), restyClient: newStubRestyClient(), tokenCache: cache}

	access, _, err := h.GetTokens()
	if err != nil {
		t.Fatalf("get tokens: %v", err)
	}

	if access != "new-access" {
		t.Fatalf("access = %q, want the rotated token instead of the expired cached one", access)
	}
	form := stub.only(t).form(t)
	if form.Get("grant_type") != "refresh_token" || form.Get("refresh_token") != "peer-refresh" {
		t.Errorf("form = %v, want a refresh grant with the cached refresh token", form)
	}
}

func TestAdminTokenHelperGetTokens_lockBackendFails_refreshesWithoutHoldingALock(t *testing.T) {
	stub := newKeycloakStub(t).on(http.MethodPost, tokenPath, http.StatusOK, tokenGrantBody)
	cache := &scriptedTokenCache{lockErr: errors.New("valkey down")}
	h := &AdminTokenHelper{cfg: stubConfig(stub.URL), restyClient: newStubRestyClient(), tokenCache: cache}

	access, _, err := h.GetTokens()
	if err != nil {
		t.Fatalf("get tokens: %v", err)
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

func TestAdminTokenHelperGetTokens_malformedPasswordGrantBody_reportsError(t *testing.T) {
	stub := newKeycloakStub(t).on(http.MethodPost, tokenPath, http.StatusOK, `<html>gateway</html>`)
	h := &AdminTokenHelper{cfg: stubConfig(stub.URL), restyClient: newStubRestyClient()}

	access, refresh, err := h.GetTokens()

	if err == nil {
		t.Fatal("err = nil, want a decode error")
	}
	if access != "" || refresh != "" {
		t.Errorf("tokens = (%q, %q), want none", access, refresh)
	}
}

func TestAdminTokenHelperGetTokens_malformedRefreshGrantBody_reportsErrorAndKeepsOldSession(t *testing.T) {
	stub := newKeycloakStub(t).on(http.MethodPost, tokenPath, http.StatusOK, `<html>gateway</html>`)
	h := &AdminTokenHelper{
		cfg:              stubConfig(stub.URL),
		restyClient:      newStubRestyClient(),
		refreshToken:     "old-refresh",
		refreshExpiresAt: time.Now().Add(10 * time.Minute),
	}

	_, _, err := h.GetTokens()

	if err == nil {
		t.Fatal("err = nil, want a decode error")
	}
	if h.refreshToken != "old-refresh" {
		t.Errorf("refresh token = %q, want the previous one kept for the next attempt", h.refreshToken)
	}
}

func TestAdminTokenHelperGetTokens_keycloakUnreachableOnRefreshGrant_reportsError(t *testing.T) {
	h := &AdminTokenHelper{
		cfg:              stubConfig(unreachableKeycloakURL(t)),
		restyClient:      newStubRestyClient(),
		refreshToken:     "old-refresh",
		refreshExpiresAt: time.Now().Add(10 * time.Minute),
	}

	_, _, err := h.GetTokens()

	if err == nil || !strings.Contains(err.Error(), "failed to refresh tokens") {
		t.Fatalf("err = %v, want the refresh failure reported", err)
	}
}

func TestAdminTokenHelperGetTokens_passwordGrantRejected_reportsStatusAndOAuthErrorCode(t *testing.T) {
	stub := newKeycloakStub(t).on(http.MethodPost, tokenPath, http.StatusUnauthorized,
		`{"error":"invalid_grant","error_description":"Invalid user credentials"}`)
	h := &AdminTokenHelper{cfg: stubConfig(stub.URL), restyClient: newStubRestyClient()}

	access, refresh, err := h.GetTokens()

	if err == nil {
		t.Fatal("err = nil, want the rejected grant reported")
	}
	for _, want := range []string{"401", "invalid_grant", "Invalid user credentials"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("err = %v, want it to mention %q", err, want)
		}
	}
	if access != "" || refresh != "" {
		t.Errorf("tokens = (%q, %q), want none", access, refresh)
	}
	if h.accessToken != "" || !h.expiresAt.IsZero() {
		t.Errorf("session = (%q, %v), want it left empty", h.accessToken, h.expiresAt)
	}
}

func TestAdminTokenHelperGetTokens_refreshGrantRejected_fallsBackToPasswordGrant(t *testing.T) {
	stub := newKeycloakStub(t)
	stub.respond = func(r stubRequest) *stubReply {
		if strings.Contains(r.Body, "grant_type=refresh_token") {
			return &stubReply{Status: http.StatusBadRequest, Body: `{"error":"invalid_grant","error_description":"Token is not active"}`}
		}
		return &stubReply{Status: http.StatusOK, Body: tokenGrantBody}
	}
	h := &AdminTokenHelper{
		cfg:              stubConfig(stub.URL),
		restyClient:      newStubRestyClient(),
		refreshToken:     "revoked-refresh",
		refreshExpiresAt: time.Now().Add(10 * time.Minute),
	}

	access, refresh, err := h.GetTokens()
	if err != nil {
		t.Fatalf("get tokens: %v", err)
	}

	if access != "new-access" || refresh != "new-refresh" {
		t.Fatalf("tokens = (%q, %q), want the pair from the password grant", access, refresh)
	}
	requests := stub.recorded()
	if len(requests) != 2 {
		t.Fatalf("requests = %d, want the refresh grant then the password grant", len(requests))
	}
	if got := requests[1].form(t).Get("grant_type"); got != "password" {
		t.Errorf("second grant_type = %q, want password", got)
	}
}

func TestAdminTokenHelperGetTokens_refreshAndPasswordGrantsRejected_reportsErrorAndKeepsOldSession(t *testing.T) {
	stub := newKeycloakStub(t).on(http.MethodPost, tokenPath, http.StatusUnauthorized, `{"error":"invalid_grant"}`)
	h := &AdminTokenHelper{
		cfg:              stubConfig(stub.URL),
		restyClient:      newStubRestyClient(),
		refreshToken:     "old-refresh",
		refreshExpiresAt: time.Now().Add(10 * time.Minute),
	}

	_, _, err := h.GetTokens()

	if err == nil || !strings.Contains(err.Error(), "fetch tokens") {
		t.Fatalf("err = %v, want the password grant failure reported", err)
	}
	if h.refreshToken != "old-refresh" {
		t.Errorf("refresh token = %q, want the previous one kept", h.refreshToken)
	}
}

func TestAdminTokenHelperGetTokens_refreshGrantServerError_reportsErrorWithoutPasswordGrant(t *testing.T) {
	stub := newKeycloakStub(t).on(http.MethodPost, tokenPath, http.StatusServiceUnavailable, `unavailable`)
	h := &AdminTokenHelper{
		cfg:              stubConfig(stub.URL),
		restyClient:      newStubRestyClient(),
		refreshToken:     "old-refresh",
		refreshExpiresAt: time.Now().Add(10 * time.Minute),
	}

	_, _, err := h.GetTokens()

	if err == nil || !strings.Contains(err.Error(), "503") {
		t.Fatalf("err = %v, want the 503 reported", err)
	}
	if got := len(stub.recorded()); got != 1 {
		t.Errorf("requests = %d, want only the refresh grant", got)
	}
}

func TestAdminTokenHelperGetTokens_successBodyWithoutAccessToken_reportsError(t *testing.T) {
	stub := newKeycloakStub(t).on(http.MethodPost, tokenPath, http.StatusOK, `{}`)
	h := &AdminTokenHelper{cfg: stubConfig(stub.URL), restyClient: newStubRestyClient()}

	_, _, err := h.GetTokens()

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
	for _, secret := range []string{stubClientSecret, stubAdminPass, "new-access", "new-refresh"} {
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
