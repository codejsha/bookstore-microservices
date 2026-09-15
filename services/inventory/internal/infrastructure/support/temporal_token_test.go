package support

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	restclient "github.com/codejsha/shared-library-go/pkg/rest/client"

	"github.com/codejsha/bookstore-microservices/inventory/internal/config"
)

type fakeClock struct {
	mu  sync.Mutex
	now time.Time
}

func (c *fakeClock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.now
}

func (c *fakeClock) Advance(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.now = c.now.Add(d)
}

type scriptedFetcher struct {
	mu    sync.Mutex
	steps []func() (temporalAccessToken, error)
	calls int
}

func (f *scriptedFetcher) fetch(context.Context) (temporalAccessToken, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.calls++
	step := f.steps[0]
	if len(f.steps) > 1 {
		f.steps = f.steps[1:]
	}
	return step()
}

func (f *scriptedFetcher) callCount() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.calls
}

func tokenStep(value string, lifetime time.Duration) func() (temporalAccessToken, error) {
	return func() (temporalAccessToken, error) {
		return temporalAccessToken{value: value, expiresIn: lifetime}, nil
	}
}

func failStep() (temporalAccessToken, error) {
	return temporalAccessToken{}, errors.New("token endpoint unavailable")
}

func newTestTokenSource(clock *fakeClock, steps ...func() (temporalAccessToken, error)) (*TemporalTokenSource, *scriptedFetcher) {
	fetcher := &scriptedFetcher{steps: steps}
	source := newTemporalTokenSource(fetcher.fetch, clock.Now)
	source.minBackoff = time.Second
	source.maxBackoff = 8 * time.Second
	source.backoff = source.minBackoff
	return source, fetcher
}

func authorizationHeader(t *testing.T, source *TemporalTokenSource) string {
	t.Helper()
	headers, err := source.GetHeaders(context.Background())
	if err != nil {
		t.Fatalf("GetHeaders error = %v", err)
	}
	return headers[temporalAuthHeader]
}

func TestTemporalTokenSourceInit_300sLifetime_schedulesRefreshAt75Percent(t *testing.T) {
	clock := &fakeClock{now: time.Unix(0, 0)}
	source, _ := newTestTokenSource(clock, tokenStep("first", 300*time.Second))

	if err := source.Init(context.Background()); err != nil {
		t.Fatalf("Init error = %v", err)
	}

	if got := source.nextRefreshDelay(); got != 225*time.Second {
		t.Fatalf("nextRefreshDelay = %v, want 225s", got)
	}
}

func TestTemporalTokenSourceInit_fetchFails_propagatesError(t *testing.T) {
	clock := &fakeClock{now: time.Unix(0, 0)}
	source, _ := newTestTokenSource(clock, failStep)

	if err := source.Init(context.Background()); err == nil {
		t.Fatal("Init error = nil, want error")
	}
}

func TestTemporalTokenSourceGetHeaders_beforeInit_fails(t *testing.T) {
	clock := &fakeClock{now: time.Unix(0, 0)}
	source, _ := newTestTokenSource(clock, failStep)

	if _, err := source.GetHeaders(context.Background()); err == nil {
		t.Fatal("GetHeaders error = nil, want error")
	}
}

func TestTemporalTokenSourceGetHeaders_repeatedCalls_reusesCachedToken(t *testing.T) {
	clock := &fakeClock{now: time.Unix(0, 0)}
	source, fetcher := newTestTokenSource(clock, tokenStep("first", 300*time.Second))
	if err := source.Init(context.Background()); err != nil {
		t.Fatalf("Init error = %v", err)
	}

	for range 5 {
		if got := authorizationHeader(t, source); got != "Bearer first" {
			t.Fatalf("authorization = %q, want %q", got, "Bearer first")
		}
	}

	if got := fetcher.callCount(); got != 1 {
		t.Fatalf("fetch calls = %d, want 1", got)
	}
}

func TestTemporalTokenSourceRefresh_fetchSucceeds_replacesTokenAndReschedules(t *testing.T) {
	clock := &fakeClock{now: time.Unix(0, 0)}
	source, _ := newTestTokenSource(clock, tokenStep("first", 300*time.Second), tokenStep("second", 120*time.Second))
	if err := source.Init(context.Background()); err != nil {
		t.Fatalf("Init error = %v", err)
	}
	clock.Advance(225 * time.Second)

	delay := source.refresh(context.Background())

	if got := authorizationHeader(t, source); got != "Bearer second" {
		t.Fatalf("authorization = %q, want %q", got, "Bearer second")
	}
	if delay != 90*time.Second {
		t.Fatalf("delay = %v, want 90s", delay)
	}
}

func TestTemporalTokenSourceRefresh_fetchFails_keepsPreviousToken(t *testing.T) {
	clock := &fakeClock{now: time.Unix(0, 0)}
	source, _ := newTestTokenSource(clock, tokenStep("first", 300*time.Second), failStep)
	if err := source.Init(context.Background()); err != nil {
		t.Fatalf("Init error = %v", err)
	}
	clock.Advance(225 * time.Second)

	delay := source.refresh(context.Background())

	if got := authorizationHeader(t, source); got != "Bearer first" {
		t.Fatalf("authorization = %q, want %q", got, "Bearer first")
	}
	if delay != time.Second {
		t.Fatalf("delay = %v, want 1s", delay)
	}
}

func TestTemporalTokenSourceRefresh_consecutiveFailures_backsOffUpToMax(t *testing.T) {
	clock := &fakeClock{now: time.Unix(0, 0)}
	source, _ := newTestTokenSource(clock, tokenStep("first", 300*time.Second), failStep)
	if err := source.Init(context.Background()); err != nil {
		t.Fatalf("Init error = %v", err)
	}

	want := []time.Duration{time.Second, 2 * time.Second, 4 * time.Second, 8 * time.Second, 8 * time.Second}
	for i, w := range want {
		if got := source.refresh(context.Background()); got != w {
			t.Fatalf("refresh #%d delay = %v, want %v", i+1, got, w)
		}
	}
}

func TestTemporalTokenSourceRefresh_successAfterFailures_resetsBackoff(t *testing.T) {
	clock := &fakeClock{now: time.Unix(0, 0)}
	source, _ := newTestTokenSource(clock,
		tokenStep("first", 300*time.Second), failStep, failStep, tokenStep("second", 300*time.Second), failStep)
	if err := source.Init(context.Background()); err != nil {
		t.Fatalf("Init error = %v", err)
	}
	for range 3 {
		source.refresh(context.Background())
	}

	if got := source.refresh(context.Background()); got != time.Second {
		t.Fatalf("delay = %v, want 1s", got)
	}
	if got := authorizationHeader(t, source); got != "Bearer second" {
		t.Fatalf("authorization = %q, want %q", got, "Bearer second")
	}
}

func TestTemporalTokenSourceStart_refreshDue_refreshesInBackgroundUntilStopped(t *testing.T) {
	var calls atomic.Int32
	source := newTemporalTokenSource(func(context.Context) (temporalAccessToken, error) {
		calls.Add(1)
		return temporalAccessToken{value: "token", expiresIn: 40 * time.Millisecond}, nil
	}, time.Now)
	if err := source.Init(context.Background()); err != nil {
		t.Fatalf("Init error = %v", err)
	}

	source.Start()
	deadline := time.Now().Add(2 * time.Second)
	for calls.Load() < 3 && time.Now().Before(deadline) {
		time.Sleep(5 * time.Millisecond)
	}
	source.Stop()

	if got := calls.Load(); got < 3 {
		t.Fatalf("fetch calls = %d, want at least 3", got)
	}
	stoppedAt := calls.Load()
	time.Sleep(100 * time.Millisecond)
	if got := calls.Load(); got != stoppedAt {
		t.Fatalf("fetch calls after Stop = %d, want %d", got, stoppedAt)
	}
}

func TestNewTemporalTokenSource_missingCredentials_fails(t *testing.T) {
	_, err := NewTemporalTokenSource(&config.TemporalAuthConfig{Enabled: true, TokenURL: "http://keycloak"}, restclient.NewRestyClient())

	if err == nil {
		t.Fatal("NewTemporalTokenSource error = nil, want error")
	}
}

func TestKeycloakTokenFetcher_clientCredentialsGrant_parsesToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("ParseForm error = %v", err)
		}
		if r.Method != http.MethodPost ||
			r.PostForm.Get("grant_type") != "client_credentials" ||
			r.PostForm.Get("client_id") != "temporal-worker-inventory" ||
			r.PostForm.Get("client_secret") != "s3cret" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"access_token":"abc.def.ghi","expires_in":300,"token_type":"Bearer"}`))
	}))
	t.Cleanup(server.Close)

	fetch := newKeycloakTokenFetcher(restclient.NewRestyClient(), &config.TemporalAuthConfig{
		TokenURL:     server.URL,
		ClientID:     "temporal-worker-inventory",
		ClientSecret: "s3cret",
	})
	token, err := fetch(context.Background())

	if err != nil {
		t.Fatalf("fetch error = %v", err)
	}
	if token.value != "abc.def.ghi" || token.expiresIn != 300*time.Second {
		t.Fatalf("token = %+v, want abc.def.ghi/300s", token)
	}
}

func TestKeycloakTokenFetcher_errorStatus_fails(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	t.Cleanup(server.Close)

	fetch := newKeycloakTokenFetcher(restclient.NewRestyClient(), &config.TemporalAuthConfig{
		TokenURL:     server.URL,
		ClientID:     "temporal-worker-inventory",
		ClientSecret: "wrong",
	})

	if _, err := fetch(context.Background()); err == nil {
		t.Fatal("fetch error = nil, want error")
	}
}
