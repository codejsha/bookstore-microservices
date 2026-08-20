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
	e.GET("/api/v1/works", ok)
	e.GET("/api/v1/works/:uid", ok)
	e.POST("/api/v1/works", ok)
	e.PUT("/api/v1/works/:uid", ok)
	e.POST("/api/v1/authors", ok)
	e.PUT("/api/v1/editions/:uid", ok)
	e.POST("/api/v1/publishers", ok)
	e.POST("/api/v1/subjects", ok)
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
	shopper   = map[string]string{HeaderUserID: "u-1", HeaderUserRoles: "VIEW,PROFILE"}
	curator   = map[string]string{HeaderUserID: "u-2", HeaderUserRoles: "admin"}
)

func TestGinAuthorizationMiddleware_WritesRequireAdmin(t *testing.T) {
	writes := []struct{ method, path string }{
		{http.MethodPost, "/api/v1/works"},
		{http.MethodPut, "/api/v1/works/w-1"},
		{http.MethodPost, "/api/v1/authors"},
		{http.MethodPut, "/api/v1/editions/e-1"},
		{http.MethodPost, "/api/v1/publishers"},
		{http.MethodPost, "/api/v1/subjects"},
	}
	e := newAuthedEngine()

	for _, w := range writes {
		t.Run(w.method+" "+w.path, func(t *testing.T) {
			if got := do(t, e, w.method, w.path, anonymous); got != http.StatusUnauthorized {
				t.Errorf("anonymous = %d, want 401", got)
			}
			if got := do(t, e, w.method, w.path, shopper); got != http.StatusForbidden {
				t.Errorf("non-admin = %d, want 403", got)
			}
			if got := do(t, e, w.method, w.path, curator); got != http.StatusOK {
				t.Errorf("admin = %d, want 200", got)
			}
		})
	}
}

func TestGinAuthorizationMiddleware_ReadsStayAnonymous(t *testing.T) {
	e := newAuthedEngine()

	for _, path := range []string{"/api/v1/works", "/api/v1/works/w-1"} {
		t.Run(path, func(t *testing.T) {
			if got := do(t, e, http.MethodGet, path, anonymous); got != http.StatusOK {
				t.Errorf("anonymous = %d, want 200", got)
			}
			if got := do(t, e, http.MethodGet, path, shopper); got != http.StatusOK {
				t.Errorf("shopper = %d, want 200", got)
			}
		})
	}
}

func TestGinAuthorizationMiddleware_HealthStaysAnonymous(t *testing.T) {
	if got := do(t, newAuthedEngine(), http.MethodGet, "/health", anonymous); got != http.StatusOK {
		t.Errorf("GET /health = %d, want 200", got)
	}
}
