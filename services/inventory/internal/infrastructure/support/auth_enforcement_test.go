package support

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func newAuthedEngine() *gin.Engine {
	gin.SetMode(gin.TestMode)
	e := gin.New()
	e.GET("/health", func(c *gin.Context) { c.Status(http.StatusOK) })
	e.Use(GinPrincipalMiddleware())
	e.Use(GinAuthorizationMiddleware())

	ok := func(c *gin.Context) { c.Status(http.StatusOK) }
	e.GET("/api/v1/stocks", ok)
	e.POST("/api/v1/stocks/adjust", ok)
	e.POST("/api/v1/stocks/receive", ok)
	e.POST("/api/v1/stocks/release", ok)
	e.POST("/api/v1/stocks/reserve", ok)
	e.POST("/api/v1/audits", ok)
	e.POST("/api/v1/audits/:uid/complete", ok)
	e.POST("/api/v1/transfers", ok)
	e.POST("/api/v1/transfers/:uid/cancel", ok)
	e.POST("/api/v1/transfers/:uid/complete", ok)
	e.POST("/api/v1/closings", ok)
	e.POST("/api/v1/warehouses", ok)
	e.PUT("/api/v1/warehouses/:uid", ok)
	return e
}

func do(t *testing.T, e *gin.Engine, method, path string, headers map[string]string) int {
	t.Helper()
	req := httptest.NewRequest(method, path, nil)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	return rec.Code
}

var (
	anonymous = map[string]string{}
	shopper   = map[string]string{HeaderUserID: "u-1", HeaderUserRoles: "USER"}
	clerk     = map[string]string{HeaderUserID: "u-2", HeaderUserRoles: "STAFF,USER"}
	manager   = map[string]string{HeaderUserID: "u-3", HeaderUserRoles: "MANAGE,STAFF,USER"}
	automaton = map[string]string{HeaderUserID: "u-4", HeaderUserRoles: "SYSTEM"}
)

var staffRequests = []struct{ method, path string }{
	{http.MethodPost, "/api/v1/stocks/adjust"},
	{http.MethodPost, "/api/v1/stocks/receive"},
	{http.MethodPost, "/api/v1/stocks/release"},
	{http.MethodPost, "/api/v1/stocks/reserve"},
	{http.MethodPost, "/api/v1/audits"},
	{http.MethodPost, "/api/v1/audits/a-1/complete"},
	{http.MethodPost, "/api/v1/transfers"},
	{http.MethodPost, "/api/v1/transfers/t-1/cancel"},
	{http.MethodPost, "/api/v1/transfers/t-1/complete"},
}

var manageRequests = []struct{ method, path string }{
	{http.MethodPost, "/api/v1/closings"},
	{http.MethodPost, "/api/v1/warehouses"},
	{http.MethodPut, "/api/v1/warehouses/w-1"},
}

func TestGinAuthorizationMiddleware_RequestWithoutPrincipal_Unauthorized(t *testing.T) {
	e := newAuthedEngine()

	for _, r := range append(append([]struct{ method, path string }{{http.MethodGet, "/api/v1/stocks"}}, staffRequests...), manageRequests...) {
		t.Run(r.method+" "+r.path, func(t *testing.T) {
			if got := do(t, e, r.method, r.path, anonymous); got != http.StatusUnauthorized {
				t.Errorf("anonymous = %d, want 401", got)
			}
		})
	}
}

func TestGinAuthorizationMiddleware_StaffRouteAsUser_Forbidden(t *testing.T) {
	e := newAuthedEngine()

	for _, r := range staffRequests {
		t.Run(r.method+" "+r.path, func(t *testing.T) {
			if got := do(t, e, r.method, r.path, shopper); got != http.StatusForbidden {
				t.Errorf("user = %d, want 403", got)
			}
		})
	}
}

func TestGinAuthorizationMiddleware_StaffRouteAsStaff_Allowed(t *testing.T) {
	e := newAuthedEngine()

	for _, r := range staffRequests {
		t.Run(r.method+" "+r.path, func(t *testing.T) {
			if got := do(t, e, r.method, r.path, clerk); got != http.StatusOK {
				t.Errorf("staff = %d, want 200", got)
			}
		})
	}
}

func TestGinAuthorizationMiddleware_ManageRouteAsStaff_Forbidden(t *testing.T) {
	e := newAuthedEngine()

	for _, r := range manageRequests {
		t.Run(r.method+" "+r.path, func(t *testing.T) {
			if got := do(t, e, r.method, r.path, clerk); got != http.StatusForbidden {
				t.Errorf("staff = %d, want 403", got)
			}
		})
	}
}

func TestGinAuthorizationMiddleware_GatedRouteAsCompositeManager_Allowed(t *testing.T) {
	e := newAuthedEngine()

	for _, r := range append(append([]struct{ method, path string }{}, staffRequests...), manageRequests...) {
		t.Run(r.method+" "+r.path, func(t *testing.T) {
			if got := do(t, e, r.method, r.path, manager); got != http.StatusOK {
				t.Errorf("manager = %d, want 200", got)
			}
		})
	}
}

func TestGinAuthorizationMiddleware_GatedRouteAsSystem_Allowed(t *testing.T) {
	e := newAuthedEngine()

	for _, r := range append(append([]struct{ method, path string }{}, staffRequests...), manageRequests...) {
		t.Run(r.method+" "+r.path, func(t *testing.T) {
			if got := do(t, e, r.method, r.path, automaton); got != http.StatusOK {
				t.Errorf("system = %d, want 200", got)
			}
		})
	}
}

func TestGinAuthorizationMiddleware_ReadAsUser_Allowed(t *testing.T) {
	if got := do(t, newAuthedEngine(), http.MethodGet, "/api/v1/stocks", shopper); got != http.StatusOK {
		t.Errorf("user = %d, want 200", got)
	}
}

func TestGinAuthorizationMiddleware_HealthWithoutPrincipal_Allowed(t *testing.T) {
	if got := do(t, newAuthedEngine(), http.MethodGet, "/health", anonymous); got != http.StatusOK {
		t.Errorf("GET /health = %d, want 200", got)
	}
}
