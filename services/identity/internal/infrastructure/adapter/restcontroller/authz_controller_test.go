package restcontroller

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/codejsha/bookstore-microservices/identity/internal/application/port/security"
	"github.com/codejsha/bookstore-microservices/identity/internal/infrastructure/adapter/keycloak"
)

type fakeRiskChecker struct {
	entry *security.RiskEntry
	err   error
}

func (f *fakeRiskChecker) Check(context.Context, string) (*security.RiskEntry, error) {
	return f.entry, f.err
}

type fakeIntrospector struct {
	result keycloak.Introspection
	err    error
	calls  int
}

func (f *fakeIntrospector) Introspect(context.Context, string) (keycloak.Introspection, error) {
	f.calls++
	return f.result, f.err
}

type fakeRevocationChecker struct {
	revoked bool
	err     error
	subs    []string
	iats    []int64
}

func (f *fakeRevocationChecker) IsRevoked(_ context.Context, sub string, tokenIssuedAt int64) (bool, error) {
	f.subs = append(f.subs, sub)
	f.iats = append(f.iats, tokenIssuedAt)
	return f.revoked, f.err
}

func check(t *testing.T, in *fakeIntrospector, rev *fakeRevocationChecker, authHeader string) int {
	t.Helper()
	return checkRisk(t, in, rev, &fakeRiskChecker{}, http.MethodGet, authHeader)
}

func checkRisk(
	t *testing.T,
	in *fakeIntrospector,
	rev *fakeRevocationChecker,
	risk *fakeRiskChecker,
	method string,
	authHeader string,
) int {
	t.Helper()
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	NewAuthzController(in, rev, risk).RegisterRoutes(engine)

	req := httptest.NewRequest(method, AuthzPathPrefix+"/api/v1/books", nil)
	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)
	return rec.Code
}

func TestCheck_ActiveTokenAllowed(t *testing.T) {
	in := &fakeIntrospector{result: keycloak.Introspection{Active: true, Sub: "user-1", Iat: 1700}}
	rev := &fakeRevocationChecker{}

	if code := check(t, in, rev, "Bearer good-token"); code != http.StatusOK {
		t.Fatalf("status = %d, want 200", code)
	}
	if len(rev.subs) != 1 || rev.subs[0] != "user-1" {
		t.Fatalf("denylist checked %v, want [user-1]", rev.subs)
	}
	if len(rev.iats) != 1 || rev.iats[0] != 1700 {
		t.Fatalf("denylist iats %v, want [1700]", rev.iats)
	}
}

func TestCheck_InactiveTokenDenied(t *testing.T) {
	in := &fakeIntrospector{result: keycloak.Introspection{Active: false}}
	if code := check(t, in, &fakeRevocationChecker{}, "Bearer dead-token"); code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403", code)
	}
}

func TestCheck_RevokedSubjectDeniedDespiteActiveToken(t *testing.T) {
	in := &fakeIntrospector{result: keycloak.Introspection{Active: true, Sub: "user-1"}}
	rev := &fakeRevocationChecker{revoked: true}

	if code := check(t, in, rev, "Bearer good-token"); code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403 for a revoked subject", code)
	}
}

func TestCheck_IntrospectionFailureFailsClosed(t *testing.T) {
	in := &fakeIntrospector{err: errors.New("keycloak unreachable")}
	if code := check(t, in, &fakeRevocationChecker{}, "Bearer any"); code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403 when the verdict cannot be established", code)
	}
}

func TestCheck_DenylistFailureFailsClosed(t *testing.T) {
	in := &fakeIntrospector{result: keycloak.Introspection{Active: true, Sub: "user-1"}}
	rev := &fakeRevocationChecker{err: errors.New("valkey unreachable")}

	if code := check(t, in, rev, "Bearer good"); code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403 when revocation cannot be ruled out", code)
	}
}

func TestCheck_MissingOrMalformedBearerIsUnauthorized(t *testing.T) {
	for name, header := range map[string]string{
		"absent":      "",
		"no scheme":   "abc.def.ghi",
		"wrong kind":  "Basic dXNlcjpwYXNz",
		"empty token": "Bearer ",
	} {
		t.Run(name, func(t *testing.T) {
			in := &fakeIntrospector{}
			if code := check(t, in, &fakeRevocationChecker{}, header); code != http.StatusUnauthorized {
				t.Fatalf("status = %d, want 401", code)
			}
			if in.calls != 0 {
				t.Errorf("introspected %d times, want 0 — nothing to introspect", in.calls)
			}
		})
	}
}

func TestBearerToken_IsSchemeCaseInsensitive(t *testing.T) {
	if got := bearerToken("bearer abc"); got != "abc" {
		t.Fatalf("bearerToken(%q) = %q, want %q", "bearer abc", got, "abc")
	}
}

func TestCheck_BlockedPrincipalDenied(t *testing.T) {
	in := &fakeIntrospector{result: keycloak.Introspection{Active: true, Sub: "user-1", Iat: 1700}}
	risk := &fakeRiskChecker{entry: &security.RiskEntry{Sub: "user-1", Level: security.RiskLevelBlock}}

	if code := checkRisk(t, in, &fakeRevocationChecker{}, risk, http.MethodGet, "Bearer t"); code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403 — blocked principals must be denied even on reads", code)
	}
}

func TestCheck_RestrictedPrincipalReadAllowedWriteDenied(t *testing.T) {
	entry := &security.RiskEntry{Sub: "user-1", Level: security.RiskLevelRestrict}

	in := &fakeIntrospector{result: keycloak.Introspection{Active: true, Sub: "user-1", Iat: 1700}}
	if code := checkRisk(t, in, &fakeRevocationChecker{}, &fakeRiskChecker{entry: entry}, http.MethodGet, "Bearer t"); code != http.StatusOK {
		t.Fatalf("GET status = %d, want 200 — restricted principals may still read", code)
	}

	in = &fakeIntrospector{result: keycloak.Introspection{Active: true, Sub: "user-1", Iat: 1700}}
	if code := checkRisk(t, in, &fakeRevocationChecker{}, &fakeRiskChecker{entry: entry}, http.MethodPost, "Bearer t"); code != http.StatusForbidden {
		t.Fatalf("POST status = %d, want 403 — restricted principals must not write", code)
	}
}

func TestCheck_RiskStoreUnavailableFailsClosed(t *testing.T) {
	in := &fakeIntrospector{result: keycloak.Introspection{Active: true, Sub: "user-1", Iat: 1700}}
	risk := &fakeRiskChecker{err: errors.New("valkey down")}

	if code := checkRisk(t, in, &fakeRevocationChecker{}, risk, http.MethodGet, "Bearer t"); code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403 — an unreachable risk store must fail closed", code)
	}
}
