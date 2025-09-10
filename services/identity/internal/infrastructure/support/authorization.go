package support

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/codejsha/bookstore-microservices/identity/generated/application/port/openapi"
	"github.com/codejsha/bookstore-microservices/identity/internal/domain/constant"
)

func isAdmin(p *Principal) bool {
	if p == nil {
		return false
	}
	for _, r := range p.Roles {
		switch constant.AuthRoleFromString(r) {
		case constant.AUTHROLE_MANAGE, constant.AUTHROLE_SYSTEM:
			return true
		}
	}
	return false
}

func writeUnauthorized(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, openapi.UnauthorizedError{
		Code:    http.StatusUnauthorized,
		Message: "authentication required",
	})
}

func writeForbidden(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusForbidden, openapi.ForbiddenError{
		Code:    http.StatusForbidden,
		Message: "insufficient privileges",
	})
}

func requireAdmin(c *gin.Context) bool {
	p := PrincipalFromContext(c)
	if p == nil {
		writeUnauthorized(c)
		return false
	}
	if !isAdmin(p) {
		writeForbidden(c)
		return false
	}
	return true
}

func requireSelfOrAdmin(c *gin.Context, uid string) bool {
	p := PrincipalFromContext(c)
	if p == nil {
		writeUnauthorized(c)
		return false
	}
	if p.Sub == uid || isAdmin(p) {
		return true
	}
	writeForbidden(c)
	return false
}

func GinAuthorizationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		path := c.FullPath()

		ok := true
		switch {
		case method == http.MethodGet && path == "/api/v1/users":
			ok = requireAdmin(c)
		case method == http.MethodGet && path == "/api/v1/users/by-email/:email":
			ok = requireAdmin(c)
		case path == "/api/v1/risk" || path == "/api/v1/risk/:uid":
			ok = requireAdmin(c)
		case method == http.MethodPut && path == "/api/v1/users/:uid/roles":
			ok = requireAdmin(c)
		case method == http.MethodPost && path == "/api/v1/users/:uid/suspend":
			ok = requireAdmin(c)
		case method == http.MethodPost && path == "/api/v1/users/:uid/reactivate":
			ok = requireAdmin(c)
		case method == http.MethodPost && path == "/api/v1/users/:uid/deactivate":
			ok = requireAdmin(c)
		case method == http.MethodGet && path == "/api/v1/users/:uid":
			ok = requireSelfOrAdmin(c, c.Param("uid"))
		case method == http.MethodPut && path == "/api/v1/users/:uid":
			ok = requireSelfOrAdmin(c, c.Param("uid"))
		case method == http.MethodPost && path == "/api/v1/users/sync/:idp_uid":
			ok = requireSelfOrAdmin(c, c.Param("idp_uid"))
		}

		if !ok {
			return
		}
		c.Next()
	}
}
