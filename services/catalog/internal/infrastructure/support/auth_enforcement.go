package support

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const roleAdmin = "admin"

var writeRoutes = map[string]bool{
	"POST /api/v1/works":          true,
	"PUT /api/v1/works/:uid":      true,
	"POST /api/v1/editions":       true,
	"PUT /api/v1/editions/:uid":   true,
	"POST /api/v1/authors":        true,
	"PUT /api/v1/authors/:uid":    true,
	"POST /api/v1/publishers":     true,
	"PUT /api/v1/publishers/:uid": true,
	"POST /api/v1/subjects":       true,
	"PUT /api/v1/subjects/:uid":   true,
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
			if c.Request.Method == http.MethodGet || c.Request.Method == http.MethodHead {
				c.Next()
				return
			}
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
