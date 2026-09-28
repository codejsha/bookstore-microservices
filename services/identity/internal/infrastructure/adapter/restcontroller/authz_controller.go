package restcontroller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"

	"github.com/codejsha/bookstore-microservices/identity/internal/application/port/security"
	"github.com/codejsha/bookstore-microservices/identity/internal/infrastructure/adapter/keycloak"
)

const AuthzPathPrefix = "/internal/authz"

type AuthzController struct {
	authorizer *keycloak.TokenAuthorizer
	metrics    *authzMetrics
}

func NewAuthzController(
	introspector keycloak.Introspector,
	revocation security.RevocationChecker,
	risk security.RiskChecker,
) *AuthzController {
	return &AuthzController{
		authorizer: keycloak.NewTokenAuthorizer(introspector, revocation, risk),
		metrics:    newAuthzMetrics(otel.GetMeterProvider()),
	}
}

func (c *AuthzController) RegisterRoutes(engine *gin.Engine) {
	engine.Any(AuthzPathPrefix, c.Check)
	engine.Any(AuthzPathPrefix+"/*checked", c.Check)
}

func (c *AuthzController) Check(ctx *gin.Context) {
	reqCtx := ctx.Request.Context()
	token := bearerToken(ctx.GetHeader("Authorization"))
	if token == "" {
		c.metrics.deny(reqCtx, authzReasonMissingBearer)
		ctx.Status(http.StatusUnauthorized)
		return
	}
	sub, reason, err := c.authorizer.Authorize(reqCtx, token, isWrite(ctx.Request.Method))
	if reason == keycloak.AuthzOK {
		c.metrics.allow(reqCtx, reason.String())
		ctx.Status(http.StatusOK)
		return
	}
	c.metrics.deny(reqCtx, reason.String())
	switch reason {
	case keycloak.AuthzIntrospectUnavailable:
		logrus.WithContext(reqCtx).WithError(err).
			Warn("authz denied: introspection unavailable")
	case keycloak.AuthzRevocationUnavailable:
		logrus.WithContext(reqCtx).WithError(err).
			Warn("authz denied: revocation denylist unavailable")
	case keycloak.AuthzRiskUnavailable:
		logrus.WithContext(reqCtx).WithError(err).
			Warn("authz denied: risk store unavailable")
	case keycloak.AuthzRiskBlocked, keycloak.AuthzRiskRestricted:
		logrus.WithContext(reqCtx).WithField("user_uid", sub).
			Warn("authz denied: principal flagged in risk store")
	}
	ctx.Status(http.StatusForbidden)
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
