package restcontroller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/codejsha/bookstore-microservices/identity/internal/application/port/security"
	"github.com/codejsha/bookstore-microservices/identity/internal/infrastructure/adapter/keycloak"
)

const AuthzPathPrefix = "/internal/authz"

type AuthzController struct {
	authorizer *keycloak.TokenAuthorizer
}

func NewAuthzController(
	introspector keycloak.Introspector,
	revocation security.RevocationChecker,
	risk security.RiskChecker,
) *AuthzController {
	return &AuthzController{authorizer: keycloak.NewTokenAuthorizer(introspector, revocation, risk)}
}

func (c *AuthzController) RegisterRoutes(engine *gin.Engine) {
	engine.Any(AuthzPathPrefix, c.Check)
	engine.Any(AuthzPathPrefix+"/*checked", c.Check)
}

func (c *AuthzController) Check(ctx *gin.Context) {
	token := bearerToken(ctx.GetHeader("Authorization"))
	if token == "" {
		ctx.Status(http.StatusUnauthorized)
		return
	}
	switch sub, reason, err := c.authorizer.Authorize(ctx.Request.Context(), token, isWrite(ctx.Request.Method)); reason {
	case keycloak.AuthzOK:
		ctx.Status(http.StatusOK)
	case keycloak.AuthzIntrospectUnavailable:
		logrus.WithContext(ctx.Request.Context()).WithError(err).
			Warn("authz denied: introspection unavailable")
		ctx.Status(http.StatusForbidden)
	case keycloak.AuthzRevocationUnavailable:
		logrus.WithContext(ctx.Request.Context()).WithError(err).
			Warn("authz denied: revocation denylist unavailable")
		ctx.Status(http.StatusForbidden)
	case keycloak.AuthzRiskUnavailable:
		logrus.WithContext(ctx.Request.Context()).WithError(err).
			Warn("authz denied: risk store unavailable")
		ctx.Status(http.StatusForbidden)
	case keycloak.AuthzRiskBlocked, keycloak.AuthzRiskRestricted:
		logrus.WithContext(ctx.Request.Context()).WithField("user_uid", sub).
			Warn("authz denied: principal flagged in risk store")
		ctx.Status(http.StatusForbidden)
	default:
		ctx.Status(http.StatusForbidden)
	}
}

func bearerToken(header string) string {
	return keycloak.BearerToken(header)
}

func isWrite(method string) bool {
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
		return false
	default:
		return true
	}
}
