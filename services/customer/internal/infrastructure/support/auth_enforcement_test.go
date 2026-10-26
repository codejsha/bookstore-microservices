package support

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func newEnforcementRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(GinPrincipalMiddleware())
	r.Use(GinOwnershipMiddleware())
	ok := func(c *gin.Context) { c.Status(http.StatusOK) }
	r.GET("/api/v1/categories", ok)
	r.POST("/api/v1/categories", ok)
	r.PATCH("/api/v1/categories/:uid", ok)
	r.DELETE("/api/v1/categories/:uid", ok)
	r.GET("/api/v1/faqs/:uid", ok)
	r.POST("/api/v1/faqs", ok)
	r.PATCH("/api/v1/tickets/:uid", ok)
	r.PATCH("/api/v1/tickets/:uid/status", ok)
	r.POST("/api/v1/customers/:uid/points/earn", ok)
	r.GET("/api/v1/customers/:uid/orders/:orderUid", ok)
	r.GET("/api/v1/customers/:uid/payments/:paymentUid", ok)
	r.GET("/api/v1/customers/:uid/reviews/:reviewUid", ok)
	r.PUT("/api/v1/customers/:uid/reviews/:reviewUid", ok)
	r.DELETE("/api/v1/customers/:uid/reviews/:reviewUid", ok)
	return r
}

func enforce(t *testing.T, method, path, roles string) int {
	t.Helper()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(method, path, nil)
	if roles != "" {
		req.Header.Set(HeaderUserID, "018f0000-0000-7000-8000-0000000000a1")
		req.Header.Set(HeaderUserRoles, roles)
	}
	newEnforcementRouter().ServeHTTP(w, req)
	return w.Code
}

func TestGinOwnershipMiddleware_CategoryReads_OpenToAnyPrincipal(t *testing.T) {
	if got := enforce(t, http.MethodGet, "/api/v1/categories", "USER"); got != http.StatusOK {
		t.Errorf("GET categories as USER = %d, want 200", got)
	}
	if got := enforce(t, http.MethodGet, "/api/v1/faqs/x", "USER"); got != http.StatusOK {
		t.Errorf("GET faq as USER = %d, want 200", got)
	}
}

func TestGinOwnershipMiddleware_CategoryAndFaqWrites_RequireManager(t *testing.T) {
	cases := []struct{ method, path string }{
		{http.MethodPost, "/api/v1/categories"},
		{http.MethodPatch, "/api/v1/categories/x"},
		{http.MethodDelete, "/api/v1/categories/x"},
		{http.MethodPost, "/api/v1/faqs"},
	}
	for _, tc := range cases {
		if got := enforce(t, tc.method, tc.path, "STAFF,USER"); got != http.StatusForbidden {
			t.Errorf("%s %s as STAFF = %d, want 403", tc.method, tc.path, got)
		}
		if got := enforce(t, tc.method, tc.path, "MANAGER"); got != http.StatusOK {
			t.Errorf("%s %s as MANAGER = %d, want 200", tc.method, tc.path, got)
		}
		if got := enforce(t, tc.method, tc.path, "SYSTEM"); got != http.StatusOK {
			t.Errorf("%s %s as SYSTEM = %d, want 200", tc.method, tc.path, got)
		}
		if got := enforce(t, tc.method, tc.path, ""); got != http.StatusUnauthorized {
			t.Errorf("%s %s anonymous = %d, want 401", tc.method, tc.path, got)
		}
	}
}

func TestGinOwnershipMiddleware_TicketStatus_RequiresStaffButTicketPatchDoesNot(t *testing.T) {
	if got := enforce(t, http.MethodPatch, "/api/v1/tickets/x/status", "USER"); got != http.StatusForbidden {
		t.Errorf("PATCH status as USER = %d, want 403", got)
	}
	if got := enforce(t, http.MethodPatch, "/api/v1/tickets/x/status", "STAFF"); got != http.StatusOK {
		t.Errorf("PATCH status as STAFF = %d, want 200", got)
	}
	if got := enforce(t, http.MethodPatch, "/api/v1/tickets/x", "USER"); got != http.StatusOK {
		t.Errorf("PATCH ticket as USER = %d, want 200 (ownership is checked in the controller)", got)
	}
}

func TestGinOwnershipMiddleware_ExistingManagerRoute_StillGated(t *testing.T) {
	if got := enforce(t, http.MethodPost, "/api/v1/customers/x/points/earn", "STAFF"); got != http.StatusForbidden {
		t.Errorf("points earn as STAFF = %d, want 403", got)
	}
	if got := enforce(t, http.MethodPost, "/api/v1/customers/x/points/earn", "MANAGER"); got != http.StatusOK {
		t.Errorf("points earn as MANAGER = %d, want 200", got)
	}
}

func TestGinOwnershipMiddleware_NestedItemRoutes_GateOnCustomerUid(t *testing.T) {
	own := "018f0000-0000-7000-8000-0000000000a1"
	other := "018f0000-0000-7000-8000-0000000000b2"
	cases := []struct{ method, suffix string }{
		{http.MethodGet, "/orders/o1"},
		{http.MethodGet, "/payments/p1"},
		{http.MethodGet, "/reviews/r1"},
		{http.MethodPut, "/reviews/r1"},
		{http.MethodDelete, "/reviews/r1"},
	}
	for _, tc := range cases {
		if got := enforce(t, tc.method, "/api/v1/customers/"+own+tc.suffix, "USER"); got != http.StatusOK {
			t.Errorf("%s %s as owner = %d, want 200", tc.method, tc.suffix, got)
		}
		if got := enforce(t, tc.method, "/api/v1/customers/"+other+tc.suffix, "USER"); got != http.StatusForbidden {
			t.Errorf("%s %s as other customer = %d, want 403", tc.method, tc.suffix, got)
		}
		if got := enforce(t, tc.method, "/api/v1/customers/"+other+tc.suffix, "STAFF"); got != http.StatusOK {
			t.Errorf("%s %s as STAFF = %d, want 200", tc.method, tc.suffix, got)
		}
		if got := enforce(t, tc.method, "/api/v1/customers/"+other+tc.suffix, ""); got != http.StatusUnauthorized {
			t.Errorf("%s %s anonymous = %d, want 401", tc.method, tc.suffix, got)
		}
	}
}

func TestGinOwnershipMiddleware_RouteMaps_MatchRegisteredRoutes(t *testing.T) {
	s := newTestGinServer(t)
	paths := map[string]bool{}
	methodRoutes := map[string]bool{}
	for _, r := range s.engine.Routes() {
		paths[r.Path] = true
		methodRoutes[r.Method+" "+r.Path] = true
	}
	for route, param := range ownerParamByRoute {
		if !paths[route] {
			t.Errorf("ownerParamByRoute key %q is not a registered route", route)
		}
		if !strings.Contains(route+"/", "/:"+param+"/") {
			t.Errorf("ownerParamByRoute key %q has no :%s segment", route, param)
		}
	}
	for route := range staffOnlyRoutes {
		if !paths[route] {
			t.Errorf("staffOnlyRoutes key %q is not a registered route", route)
		}
	}
	for route := range managerOnlyRoutes {
		if !paths[route] {
			t.Errorf("managerOnlyRoutes key %q is not a registered route", route)
		}
	}
	for route := range staffOnlyMethodRoutes {
		if !methodRoutes[route] {
			t.Errorf("staffOnlyMethodRoutes key %q is not a registered route", route)
		}
	}
	for route := range managerOnlyMethodRoutes {
		if !methodRoutes[route] {
			t.Errorf("managerOnlyMethodRoutes key %q is not a registered route", route)
		}
	}
}
