package support

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/codejsha/bookstore-microservices/identity/internal/domain/constant"
)

func hasManager(p *Principal) bool {
	if p == nil {
		return false
	}
	for _, r := range p.Roles {
		switch constant.AuthRoleFromString(r) {
		case constant.AUTHROLE_MANAGER, constant.AUTHROLE_SYSTEM:
			return true
		}
	}
	return false
}

func writeUnauthorized(c *gin.Context) {
	abortWithProblem(c, http.StatusUnauthorized, "authentication required")
}

func writeForbidden(c *gin.Context) {
	abortWithProblem(c, http.StatusForbidden, "insufficient privileges")
}

func requireManager(c *gin.Context) bool {
	p := PrincipalFromContext(c)
	if p == nil {
		writeUnauthorized(c)
		return false
	}
	if !hasManager(p) {
		writeForbidden(c)
		return false
	}
	return true
}

func requireSelfOrManager(c *gin.Context, uid string) bool {
	p := PrincipalFromContext(c)
	if p == nil {
		writeUnauthorized(c)
		return false
	}
	if p.Sub == uid || hasManager(p) {
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
			ok = requireManager(c)
		case method == http.MethodGet && path == "/api/v1/users/by-email/:email":
			ok = requireManager(c)
		case path == "/api/v1/risk" || path == "/api/v1/risk/:uid":
			ok = requireManager(c)
		case method == http.MethodPut && path == "/api/v1/users/:uid/roles":
			ok = requireManager(c)
		case method == http.MethodPost && path == "/api/v1/users/:uid/suspend":
			ok = requireManager(c)
		case method == http.MethodPost && path == "/api/v1/users/:uid/reactivate":
			ok = requireManager(c)
		case method == http.MethodPost && path == "/api/v1/users/:uid/deactivate":
			ok = requireManager(c)
		case method == http.MethodGet && path == "/api/v1/users/:uid":
			ok = requireSelfOrManager(c, c.Param("uid"))
		case method == http.MethodPut && path == "/api/v1/users/:uid":
			ok = requireSelfOrManager(c, c.Param("uid"))
		case method == http.MethodPost && path == "/api/v1/users/sync/:idp_uid":
			ok = requireSelfOrManager(c, c.Param("idp_uid"))
		}

		if !ok {
			return
		}
		c.Next()
	}
}
