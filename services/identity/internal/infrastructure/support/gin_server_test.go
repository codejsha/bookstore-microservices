package support

import (
	"context"
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
