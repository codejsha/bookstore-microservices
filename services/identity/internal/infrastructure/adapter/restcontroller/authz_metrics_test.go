package restcontroller

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"

	"github.com/codejsha/bookstore-microservices/identity/internal/domain/constant"
	"github.com/codejsha/bookstore-microservices/identity/internal/infrastructure/adapter/keycloak"
)

func installManualReader(t *testing.T) *sdkmetric.ManualReader {
	t.Helper()
	reader := sdkmetric.NewManualReader()
	provider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader))
	previous := otel.GetMeterProvider()
	otel.SetMeterProvider(provider)
	t.Cleanup(func() {
		otel.SetMeterProvider(previous)
		_ = provider.Shutdown(context.Background())
	})
	return reader
}

func decisionCount(t *testing.T, reader *sdkmetric.ManualReader, outcome, reason string) int64 {
	t.Helper()
	var rm metricdata.ResourceMetrics
	if err := reader.Collect(context.Background(), &rm); err != nil {
		t.Fatalf("collect metrics: %v", err)
	}
	want := attribute.NewSet(
		attribute.String("outcome", outcome),
		attribute.String("reason", reason),
	)
	for _, scope := range rm.ScopeMetrics {
		for _, m := range scope.Metrics {
			if m.Name != constant.MetricAuthzDecisions {
				continue
			}
			sum, ok := m.Data.(metricdata.Sum[int64])
			if !ok {
				t.Fatalf("metric %s is %T, want Sum[int64]", m.Name, m.Data)
			}
			if !sum.IsMonotonic {
				t.Fatalf("metric %s is not monotonic", m.Name)
			}
			for _, dp := range sum.DataPoints {
				if dp.Attributes.Equals(&want) {
					return dp.Value
				}
			}
		}
	}
	return 0
}

func TestCheck_ActiveToken_CountsAllowDecision(t *testing.T) {
	reader := installManualReader(t)
	in := &fakeIntrospector{result: keycloak.Introspection{Active: true, Sub: "user-1", Iat: 1700}}

	if code := check(t, in, &fakeRevocationChecker{}, "Bearer good"); code != http.StatusOK {
		t.Fatalf("status = %d, want 200", code)
	}
	if got := decisionCount(t, reader, "allow", "ok"); got != 1 {
		t.Fatalf("allow/ok decisions = %d, want 1", got)
	}
}

func TestCheck_RevocationUnavailable_CountsDenyWithReason(t *testing.T) {
	reader := installManualReader(t)
	in := &fakeIntrospector{result: keycloak.Introspection{Active: true, Sub: "user-1"}}
	rev := &fakeRevocationChecker{err: errors.New("valkey unreachable")}

	if code := check(t, in, rev, "Bearer good"); code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403", code)
	}
	if got := decisionCount(t, reader, "deny", "revocation_unavailable"); got != 1 {
		t.Fatalf("deny/revocation_unavailable decisions = %d, want 1", got)
	}
	if got := decisionCount(t, reader, "allow", "ok"); got != 0 {
		t.Fatalf("allow/ok decisions = %d, want 0", got)
	}
}

func TestCheck_RiskStoreUnavailable_CountsDenyWithReason(t *testing.T) {
	reader := installManualReader(t)
	in := &fakeIntrospector{result: keycloak.Introspection{Active: true, Sub: "user-1", Iat: 1700}}
	risk := &fakeRiskChecker{err: errors.New("valkey down")}

	if code := checkRisk(t, in, &fakeRevocationChecker{}, risk, http.MethodGet, "Bearer t"); code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403", code)
	}
	if got := decisionCount(t, reader, "deny", "risk_unavailable"); got != 1 {
		t.Fatalf("deny/risk_unavailable decisions = %d, want 1", got)
	}
}

func TestCheck_MissingBearer_CountsDenyWithoutIntrospecting(t *testing.T) {
	reader := installManualReader(t)
	in := &fakeIntrospector{}

	if code := check(t, in, &fakeRevocationChecker{}, ""); code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", code)
	}
	if got := decisionCount(t, reader, "deny", "missing_bearer"); got != 1 {
		t.Fatalf("deny/missing_bearer decisions = %d, want 1", got)
	}
	if in.calls != 0 {
		t.Fatalf("introspected %d times, want 0", in.calls)
	}
}

func TestAuthzReasonString_EveryReason_MapsToStableLabel(t *testing.T) {
	want := map[keycloak.AuthzReason]string{
		keycloak.AuthzOK:                    "ok",
		keycloak.AuthzInactive:              "inactive",
		keycloak.AuthzIntrospectUnavailable: "introspect_unavailable",
		keycloak.AuthzRevoked:               "revoked",
		keycloak.AuthzRevocationUnavailable: "revocation_unavailable",
		keycloak.AuthzRiskBlocked:           "risk_blocked",
		keycloak.AuthzRiskRestricted:        "risk_restricted",
		keycloak.AuthzRiskUnavailable:       "risk_unavailable",
	}
	for reason, label := range want {
		if got := reason.String(); got != label {
			t.Errorf("AuthzReason(%d).String() = %q, want %q", int(reason), got, label)
		}
	}
	if got := keycloak.AuthzReason(99).String(); got != "unknown" {
		t.Errorf("unknown reason label = %q, want %q", got, "unknown")
	}
}
