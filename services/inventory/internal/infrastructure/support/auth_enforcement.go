package support

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const roleAdmin = "admin"

var writeRoutes = map[string]bool{
	"POST /api/v1/audits":                  true,
	"POST /api/v1/audits/:uid/complete":    true,
	"POST /api/v1/closings":                true,
	"POST /api/v1/stocks/adjust":           true,
	"POST /api/v1/stocks/receive":          true,
	"POST /api/v1/stocks/release":          true,
	"POST /api/v1/stocks/reserve":          true,
	"POST /api/v1/transfers":               true,
	"POST /api/v1/transfers/:uid/cancel":   true,
	"POST /api/v1/transfers/:uid/complete": true,
	"POST /api/v1/warehouses":              true,
	"PUT /api/v1/warehouses/:uid":          true,
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
		if writeRoutes[c.Request.Method+" "+route] && !p.HasRole(roleAdmin) {
			abortWithProblem(c, http.StatusForbidden, "forbidden")
			return
		}
		c.Next()
	}
}
