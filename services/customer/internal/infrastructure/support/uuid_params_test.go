package support

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func uuidParamStatus(t *testing.T, target string) int {
	t.Helper()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(GinUuidParamMiddleware())
	ok := func(c *gin.Context) { c.Status(http.StatusOK) }
	r.GET("/api/v1/tickets", ok)
	r.GET("/api/v1/tickets/:uid/comments", ok)
	r.GET("/api/v1/customers/:uid/orders/:orderUid", ok)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, target, nil))
	return w.Code
}

func TestGinUuidParamMiddleware_MalformedUid_Returns400(t *testing.T) {
	const valid = "018f0000-0000-7000-8000-0000000000a1"
	for name, target := range map[string]string{
		"path":         "/api/v1/tickets/not-a-uuid/comments",
		"second path":  "/api/v1/customers/" + valid + "/orders/x",
		"query":        "/api/v1/tickets?assignee_uid=x",
		"urn form":     "/api/v1/tickets/urn:uuid:" + valid + "/comments",
		"repeat query": "/api/v1/tickets?category_uid=" + valid + "&category_uid=bad",
	} {
		t.Run(name, func(t *testing.T) {
			if got := uuidParamStatus(t, target); got != http.StatusBadRequest {
				t.Errorf("GET %s = %d, want 400", target, got)
			}
		})
	}
}

func TestGinUuidParamMiddleware_ValidOrNonUidParams_PassThrough(t *testing.T) {
	const valid = "018f0000-0000-7000-8000-0000000000a1"
	for _, target := range []string{
		"/api/v1/tickets/" + valid + "/comments",
		"/api/v1/customers/" + valid + "/orders/" + valid,
		"/api/v1/tickets?customer_uid=" + valid + "&status=OPEN&sort=created_at,desc",
	} {
		if got := uuidParamStatus(t, target); got != http.StatusOK {
			t.Errorf("GET %s = %d, want 200", target, got)
		}
	}
}
