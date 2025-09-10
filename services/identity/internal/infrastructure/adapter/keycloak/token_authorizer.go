package keycloak

import (
	"context"
	"strings"

	"github.com/codejsha/bookstore-microservices/identity/internal/application/port/security"
)

type AuthzReason int

const (
	AuthzOK AuthzReason = iota
	AuthzInactive
	AuthzIntrospectUnavailable
	AuthzRevoked
	AuthzRevocationUnavailable
	AuthzRiskBlocked
	AuthzRiskRestricted
	AuthzRiskUnavailable
)

type TokenAuthorizer struct {
	introspector Introspector
	revocation   security.RevocationChecker
	risk         security.RiskChecker
}

func NewTokenAuthorizer(
	introspector Introspector,
	revocation security.RevocationChecker,
	risk security.RiskChecker,
) *TokenAuthorizer {
	return &TokenAuthorizer{introspector: introspector, revocation: revocation, risk: risk}
}

func (a *TokenAuthorizer) Authorize(ctx context.Context, token string, write bool) (string, AuthzReason, error) {
	result, err := a.introspector.Introspect(ctx, token)
	if err != nil {
		return "", AuthzIntrospectUnavailable, err
	}
	if !result.Active {
		return result.Sub, AuthzInactive, nil
	}
	revoked, err := a.revocation.IsRevoked(ctx, result.Sub, result.Iat)
	if err != nil {
		return result.Sub, AuthzRevocationUnavailable, err
	}
	if revoked {
		return result.Sub, AuthzRevoked, nil
	}
	entry, err := a.risk.Check(ctx, result.Sub)
	if err != nil {
		return result.Sub, AuthzRiskUnavailable, err
	}
	if entry != nil {
		switch {
		case entry.Level == security.RiskLevelBlock:
			return result.Sub, AuthzRiskBlocked, nil
		case write:
			return result.Sub, AuthzRiskRestricted, nil
		}
	}
	return result.Sub, AuthzOK, nil
}

func BearerToken(header string) string {
	const prefix = "Bearer "
	if len(header) <= len(prefix) || !strings.EqualFold(header[:len(prefix)], prefix) {
		return ""
	}
	return strings.TrimSpace(header[len(prefix):])
}
