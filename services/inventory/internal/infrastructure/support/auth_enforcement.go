package support

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	roleStaff  = "STAFF"
	roleManage = "MANAGE"
	roleSystem = "SYSTEM"
)

var writeRoutes = map[string]string{
	"POST /api/v1/audits":                  roleStaff,
	"POST /api/v1/audits/:uid/complete":    roleStaff,
	"POST /api/v1/closings":                roleManage,
	"POST /api/v1/stocks/adjust":           roleStaff,
	"POST /api/v1/stocks/receive":          roleStaff,
	"POST /api/v1/stocks/release":          roleStaff,
	"POST /api/v1/stocks/reserve":          roleStaff,
	"POST /api/v1/transfers":               roleStaff,
	"POST /api/v1/transfers/:uid/cancel":   roleStaff,
	"POST /api/v1/transfers/:uid/complete": roleStaff,
	"POST /api/v1/warehouses":              roleManage,
	"PUT /api/v1/warehouses/:uid":          roleManage,
}

func GinAuthorizationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		route := c.FullPath()
		if route == "" {
			c.Next()
			return
		}
		p := PrincipalFromContext(c)
		if p == nil {
			abortWithProblem(c, http.StatusUnauthorized, "authentication required")
			return
		}
		if required, gated := writeRoutes[c.Request.Method+" "+route]; gated && !hasAnyRole(p, required, roleSystem) {
			abortWithProblem(c, http.StatusForbidden, "forbidden")
			return
		}
		c.Next()
	}
}
