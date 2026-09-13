package support

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx/fxtest"

	"github.com/codejsha/shared-library-go/pkg/config"
	"github.com/codejsha/shared-library-go/pkg/logging"
)

func newTestGinServer(t *testing.T) *GinServer {
	t.Helper()
	appCfg := &config.AppConfig{Logging: config.LoggingConfig{Level: "error"}}
	serverCfg := &config.ServerConfig{Host: "127.0.0.1", Port: "0", Mode: gin.TestMode}
	return NewGinServer(
		fxtest.NewLifecycle(t),
		serverCfg,
		logging.NewLogHelper(appCfg),
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
	)
}

func getHealthPath(s *GinServer, path string) int {
	req := httptest.NewRequest(http.MethodGet, path, nil)
	rec := httptest.NewRecorder()
	s.engine.ServeHTTP(rec, req)
	return rec.Code
}

func TestGinServer_healthRoute_returns200(t *testing.T) {
	s := newTestGinServer(t)

	if code := getHealthPath(s, "/health"); code != http.StatusOK {
		t.Fatalf("GET /health = %d, want %d", code, http.StatusOK)
	}
}

func TestGinServer_readyRoute_nilDataSource_returns503(t *testing.T) {
	s := newTestGinServer(t)

	if code := getHealthPath(s, "/health/ready"); code != http.StatusServiceUnavailable {
		t.Fatalf("GET /health/ready = %d, want %d", code, http.StatusServiceUnavailable)
	}
}

func TestGinServer_readyRoute_checkSucceeds_returns200(t *testing.T) {
	s := newTestGinServer(t)
	s.readyCheck = func(context.Context) error { return nil }

	if code := getHealthPath(s, "/health/ready"); code != http.StatusOK {
		t.Fatalf("GET /health/ready = %d, want %d", code, http.StatusOK)
	}
}

func TestGinServer_readyRoute_checkFails_returns503(t *testing.T) {
	s := newTestGinServer(t)
	s.readyCheck = func(context.Context) error { return errors.New("ping failed") }

	if code := getHealthPath(s, "/health/ready"); code != http.StatusServiceUnavailable {
		t.Fatalf("GET /health/ready = %d, want %d", code, http.StatusServiceUnavailable)
	}
}
