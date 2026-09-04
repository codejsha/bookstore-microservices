package support

import (
	"context"
	"errors"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

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

func (f *fakeRevocationChecker) IsRevoked(_ context.Context, sub string, iat int64) (bool, error) {
	f.subs = append(f.subs, sub)
	f.iats = append(f.iats, iat)
	return f.revoked, f.err
}

func invoke(in *fakeIntrospector, rev *fakeRevocationChecker, md metadata.MD) (codes.Code, bool) {
	return invokeRisk(in, rev, &fakeRiskChecker{}, md)
}

func invokeRisk(in *fakeIntrospector, rev *fakeRevocationChecker, risk *fakeRiskChecker, md metadata.MD) (codes.Code, bool) {
	authorizer := keycloak.NewTokenAuthorizer(in, rev, risk)
	interceptor := newAuthUnaryInterceptor(authorizer)

	ctx := context.Background()
	if md != nil {
		ctx = metadata.NewIncomingContext(ctx, md)
	}
	handlerRan := false
	_, err := interceptor(ctx, nil, &grpc.UnaryServerInfo{}, func(context.Context, any) (any, error) {
		handlerRan = true
		return nil, nil
	})
	return status.Code(err), handlerRan
}

func bearerMD(token string) metadata.MD {
	return metadata.Pairs("authorization", "Bearer "+token)
}

func TestInterceptor_WhenTokenActive_CallsHandler(t *testing.T) {
	in := &fakeIntrospector{result: keycloak.Introspection{Active: true, Sub: "user-1", Iat: 1700}}
	rev := &fakeRevocationChecker{}

	code, ran := invoke(in, rev, bearerMD("good"))
	if code != codes.OK || !ran {
		t.Fatalf("code=%v ran=%v, want OK+handler ran", code, ran)
	}
	if len(rev.subs) != 1 || rev.subs[0] != "user-1" || rev.iats[0] != 1700 {
		t.Fatalf("denylist checked subs=%v iats=%v, want [user-1]/[1700]", rev.subs, rev.iats)
	}
}

func TestInterceptor_WhenMetadataMissing_ReturnsUnauthenticatedWithoutIntrospecting(t *testing.T) {
	in := &fakeIntrospector{}
	code, ran := invoke(in, &fakeRevocationChecker{}, nil)
	if code != codes.Unauthenticated || ran {
		t.Fatalf("code=%v ran=%v, want Unauthenticated, handler not run", code, ran)
	}
	if in.calls != 0 {
		t.Errorf("introspected %d times, want 0 — nothing to introspect", in.calls)
	}
}

func TestInterceptor_WhenTokenMissingOrMalformed_ReturnsUnauthenticatedWithoutIntrospecting(t *testing.T) {
	for name, md := range map[string]metadata.MD{
		"no authorization key": metadata.Pairs("other", "x"),
		"wrong scheme":         metadata.Pairs("authorization", "Basic dXNlcjpwYXNz"),
		"empty bearer":         metadata.Pairs("authorization", "Bearer "),
	} {
		t.Run(name, func(t *testing.T) {
			in := &fakeIntrospector{}
			code, ran := invoke(in, &fakeRevocationChecker{}, md)
			if code != codes.Unauthenticated || ran {
				t.Fatalf("code=%v ran=%v, want Unauthenticated, handler not run", code, ran)
			}
			if in.calls != 0 {
				t.Errorf("introspected %d times, want 0", in.calls)
			}
		})
	}
}

func TestInterceptor_WhenTokenInactive_ReturnsUnauthenticated(t *testing.T) {
	in := &fakeIntrospector{result: keycloak.Introspection{Active: false}}
	code, ran := invoke(in, &fakeRevocationChecker{}, bearerMD("dead"))
	if code != codes.Unauthenticated || ran {
		t.Fatalf("code=%v ran=%v, want Unauthenticated for an inactive token", code, ran)
	}
}

func TestInterceptor_WhenSubjectRevoked_ReturnsPermissionDenied(t *testing.T) {
	in := &fakeIntrospector{result: keycloak.Introspection{Active: true, Sub: "user-1"}}
	rev := &fakeRevocationChecker{revoked: true}
	code, ran := invoke(in, rev, bearerMD("good"))
	if code != codes.PermissionDenied || ran {
		t.Fatalf("code=%v ran=%v, want PermissionDenied for a revoked subject", code, ran)
	}
}

func TestInterceptor_WhenIntrospectionUnavailable_ReturnsPermissionDenied(t *testing.T) {
	in := &fakeIntrospector{err: errors.New("keycloak unreachable")}
	code, ran := invoke(in, &fakeRevocationChecker{}, bearerMD("any"))
	if code != codes.PermissionDenied || ran {
		t.Fatalf("code=%v ran=%v, want PermissionDenied when the verdict cannot be established", code, ran)
	}
}

func TestInterceptor_WhenDenylistUnavailable_ReturnsPermissionDenied(t *testing.T) {
	in := &fakeIntrospector{result: keycloak.Introspection{Active: true, Sub: "user-1"}}
	rev := &fakeRevocationChecker{err: errors.New("valkey unreachable")}
	code, ran := invoke(in, rev, bearerMD("good"))
	if code != codes.PermissionDenied || ran {
		t.Fatalf("code=%v ran=%v, want PermissionDenied when revocation cannot be ruled out", code, ran)
	}
}

func TestGrpcAuth_WhenPrincipalBlocked_ReturnsPermissionDeniedWithoutHandler(t *testing.T) {
	in := &fakeIntrospector{result: keycloak.Introspection{Active: true, Sub: "user-1", Iat: 1700}}
	risk := &fakeRiskChecker{entry: &security.RiskEntry{Sub: "user-1", Level: security.RiskLevelBlock}}

	code, ran := invokeRisk(in, &fakeRevocationChecker{}, risk, metadata.Pairs("authorization", "Bearer t"))
	if ran || code != codes.PermissionDenied {
		t.Fatalf("code = %v ran = %v, want PermissionDenied without handler", code, ran)
	}
}

func TestGrpcAuth_WhenPrincipalRestricted_StillCallsReadHandler(t *testing.T) {
	in := &fakeIntrospector{result: keycloak.Introspection{Active: true, Sub: "user-1", Iat: 1700}}
	risk := &fakeRiskChecker{entry: &security.RiskEntry{Sub: "user-1", Level: security.RiskLevelRestrict}}

	code, ran := invokeRisk(in, &fakeRevocationChecker{}, risk, metadata.Pairs("authorization", "Bearer t"))
	if !ran || code != codes.OK {
		t.Fatalf("code = %v ran = %v, want handler to run — gRPC surface is read-only", code, ran)
	}
}
