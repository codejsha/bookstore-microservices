package support

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/codejsha/bookstore-microservices/identity/internal/domain/constant"
)

func hasManage(p *Principal) bool {
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
	abortWithProblem(c, http.StatusUnauthorized, "authentication required")
}

func writeForbidden(c *gin.Context) {
	abortWithProblem(c, http.StatusForbidden, "insufficient privileges")
}

func requireManage(c *gin.Context) bool {
	p := PrincipalFromContext(c)
	if p == nil {
		writeUnauthorized(c)
		return false
	}
	if !hasManage(p) {
		writeForbidden(c)
		return false
	}
	return true
}

func requireSelfOrManage(c *gin.Context, uid string) bool {
	p := PrincipalFromContext(c)
	if p == nil {
		writeUnauthorized(c)
		return false
	}
	if p.Sub == uid || hasManage(p) {
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
			ok = requireManage(c)
		case method == http.MethodGet && path == "/api/v1/users/by-email/:email":
			ok = requireManage(c)
		case path == "/api/v1/risk" || path == "/api/v1/risk/:uid":
			ok = requireManage(c)
		case method == http.MethodPut && path == "/api/v1/users/:uid/roles":
			ok = requireManage(c)
		case method == http.MethodPost && path == "/api/v1/users/:uid/suspend":
			ok = requireManage(c)
		case method == http.MethodPost && path == "/api/v1/users/:uid/reactivate":
			ok = requireManage(c)
		case method == http.MethodPost && path == "/api/v1/users/:uid/deactivate":
			ok = requireManage(c)
		case method == http.MethodGet && path == "/api/v1/users/:uid":
			ok = requireSelfOrManage(c, c.Param("uid"))
		case method == http.MethodPut && path == "/api/v1/users/:uid":
			ok = requireSelfOrManage(c, c.Param("uid"))
		case method == http.MethodPost && path == "/api/v1/users/sync/:idp_uid":
			ok = requireSelfOrManage(c, c.Param("idp_uid"))
		}

		if !ok {
			return
		}
		c.Next()
	}
}
