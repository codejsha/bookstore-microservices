package support

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	restclient "github.com/codejsha/shared-library-go/pkg/rest/client"
	"github.com/sirupsen/logrus"

	"github.com/codejsha/bookstore-microservices/inventory/internal/config"
)

const (
	temporalTokenRefreshRatio = 0.75
	temporalTokenMinBackoff   = time.Second
	temporalTokenMaxBackoff   = 30 * time.Second
	temporalTokenInitTimeout  = 15 * time.Second
	temporalAuthHeader        = "authorization"
	temporalBearerPrefix      = "Bearer "
)

type temporalAccessToken struct {
	value     string
	expiresIn time.Duration
}

type temporalTokenFetcher func(ctx context.Context) (temporalAccessToken, error)

type TemporalTokenSource struct {
	fetch        temporalTokenFetcher
	now          func() time.Time
	refreshRatio float64
	minBackoff   time.Duration
	maxBackoff   time.Duration

	mu        sync.RWMutex
	token     string
	expiresAt time.Time
	refreshAt time.Time

	refreshMu sync.Mutex
	backoff   time.Duration

	lifecycleMu sync.Mutex
	cancel      context.CancelFunc
	done        chan struct{}
}

func newTemporalTokenSource(fetch temporalTokenFetcher, now func() time.Time) *TemporalTokenSource {
	return &TemporalTokenSource{
		fetch:        fetch,
		now:          now,
		refreshRatio: temporalTokenRefreshRatio,
		minBackoff:   temporalTokenMinBackoff,
		maxBackoff:   temporalTokenMaxBackoff,
		backoff:      temporalTokenMinBackoff,
	}
}

func NewTemporalTokenSource(
	cfg *config.TemporalAuthConfig,
	restyClient *restclient.RestyClient,
) (*TemporalTokenSource, error) {
	if err := validateTemporalAuthConfig(cfg); err != nil {
		return nil, err
	}
	return newTemporalTokenSource(newKeycloakTokenFetcher(restyClient, cfg), time.Now), nil
}

func validateTemporalAuthConfig(cfg *config.TemporalAuthConfig) error {
	var missing []error
	if cfg.TokenURL == "" {
		missing = append(missing, errors.New("temporal.auth.tokenUrl is required"))
	}
	if cfg.ClientID == "" {
		missing = append(missing, errors.New("temporal.auth.clientId is required"))
	}
	if cfg.ClientSecret == "" {
		missing = append(missing, errors.New("temporal.auth.clientSecret is required"))
	}
	return errors.Join(missing...)
}

func (s *TemporalTokenSource) Init(ctx context.Context) error {
	token, err := s.fetch(ctx)
	if err != nil {
		return fmt.Errorf("fetch temporal access token: %w", err)
	}
	return s.store(token)
}

func (s *TemporalTokenSource) GetHeaders(context.Context) (map[string]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.token == "" {
		return nil, errors.New("temporal access token is not initialized")
	}
	return map[string]string{temporalAuthHeader: temporalBearerPrefix + s.token}, nil
}

func (s *TemporalTokenSource) Start() {
	s.lifecycleMu.Lock()
	defer s.lifecycleMu.Unlock()
	if s.cancel != nil {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel
	s.done = make(chan struct{})
	go s.loop(ctx, s.done)
}

func (s *TemporalTokenSource) Stop() {
	s.lifecycleMu.Lock()
	defer s.lifecycleMu.Unlock()
	if s.cancel == nil {
		return
	}
	s.cancel()
	<-s.done
	s.cancel = nil
	s.done = nil
}

func (s *TemporalTokenSource) loop(ctx context.Context, done chan struct{}) {
	defer close(done)
	timer := time.NewTimer(s.nextRefreshDelay())
	defer timer.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			timer.Reset(s.refresh(ctx))
		}
	}
}

func (s *TemporalTokenSource) refresh(ctx context.Context) time.Duration {
	s.refreshMu.Lock()
	defer s.refreshMu.Unlock()

	token, err := s.fetch(ctx)
	if err == nil {
		err = s.store(token)
	}
	if err == nil {
		s.backoff = s.minBackoff
		return s.nextRefreshDelay()
	}

	delay := s.backoff
	s.backoff = min(s.backoff*2, s.maxBackoff)
	if ctx.Err() == nil {
		s.mu.RLock()
		expiresAt := s.expiresAt
		s.mu.RUnlock()
		logrus.WithError(err).WithFields(logrus.Fields{
			"retry_in_ms": delay.Milliseconds(),
			"expires_at":  expiresAt,
		}).Warn("temporal access token refresh failed")
	}
	return delay
}

func (s *TemporalTokenSource) nextRefreshDelay() time.Duration {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return max(s.refreshAt.Sub(s.now()), 0)
}

func (s *TemporalTokenSource) store(token temporalAccessToken) error {
	if token.value == "" {
		return errors.New("temporal access token is empty")
	}
	if token.expiresIn <= 0 {
		return errors.New("temporal access token lifetime must be positive")
	}
	now := s.now()
	refreshAfter := time.Duration(float64(token.expiresIn) * s.refreshRatio)

	s.mu.Lock()
	defer s.mu.Unlock()
	s.token = token.value
	s.expiresAt = now.Add(token.expiresIn)
	s.refreshAt = now.Add(refreshAfter)
	return nil
}

type keycloakTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

func newKeycloakTokenFetcher(
	restyClient *restclient.RestyClient,
	cfg *config.TemporalAuthConfig,
) temporalTokenFetcher {
	return func(ctx context.Context) (temporalAccessToken, error) {
		var body keycloakTokenResponse
		resp, err := restyClient.Client.R().
			SetContext(ctx).
			SetHeader("Accept", "application/json").
			SetFormData(map[string]string{
				"grant_type":    "client_credentials",
				"client_id":     cfg.ClientID,
				"client_secret": cfg.ClientSecret,
			}).
			SetResult(&body).
			Post(cfg.TokenURL)
		if err != nil {
			return temporalAccessToken{}, fmt.Errorf("request token: %w", err)
		}
		if resp.IsError() {
			return temporalAccessToken{}, fmt.Errorf("token endpoint returned status %d", resp.StatusCode())
		}
		if body.AccessToken == "" {
			return temporalAccessToken{}, errors.New("token response has no access_token")
		}
		if body.ExpiresIn <= 0 {
			return temporalAccessToken{}, errors.New("token response has no positive expires_in")
		}
		return temporalAccessToken{
			value:     body.AccessToken,
			expiresIn: time.Duration(body.ExpiresIn) * time.Second,
		}, nil
	}
}
