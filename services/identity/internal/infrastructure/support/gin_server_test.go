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

	"github.com/codejsha/bookstore-microservices/identity/generated/application/port/openapi"
	"github.com/codejsha/bookstore-microservices/identity/internal/infrastructure/adapter/restcontroller"
)

type stubRiskApi struct{}

func (stubRiskApi) RiskGetAll(context.Context) (*openapi.RiskEntryListResponse, error) {
	return &openapi.RiskEntryListResponse{Entries: []openapi.RiskEntryResponse{}}, nil
}

func (stubRiskApi) RiskFlagPrincipal(context.Context, string, openapi.RiskFlagRequest) (*openapi.RiskEntryResponse, error) {
	return nil, nil
}

func (stubRiskApi) RiskUnflagPrincipal(context.Context, string) error {
	return nil
}

func TestNewGinServer_riskRoute_servesInjectedRiskApi(t *testing.T) {
	appCfg := &config.AppConfig{Logging: config.LoggingConfig{Level: "error"}}
	serverCfg := &config.ServerConfig{Host: "127.0.0.1", Port: "0", Mode: gin.TestMode}
	s := NewGinServer(
		fxtest.NewLifecycle(t),
		serverCfg,
		logging.NewLogHelper(appCfg),
		nil,
		nil,
		restcontroller.NewAuthzController(nil, nil, nil),
		stubRiskApi{},
	)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/risk", nil)
	req.Header.Set(HeaderUserID, "manager-1")
	req.Header.Set(HeaderUserRoles, "MANAGE,STAFF,USER")
	rec := httptest.NewRecorder()
	s.engine.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("GET /api/v1/risk = %d, want %d: %s", rec.Code, http.StatusOK, rec.Body.String())
	}
}

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
		restcontroller.NewAuthzController(nil, nil, nil),
		stubRiskApi{},
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
