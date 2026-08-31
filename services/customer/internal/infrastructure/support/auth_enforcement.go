package support

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const RoleAdmin = "admin"

var ownerParamByRoute = map[string]string{
	"/api/v1/customers/:uid":                       "uid",
	"/api/v1/customers/:uid/orders":                "uid",
	"/api/v1/customers/:uid/orders/:order_uid":     "uid",
	"/api/v1/customers/:uid/payments":              "uid",
	"/api/v1/customers/:uid/payments/:payment_uid": "uid",
	"/api/v1/customers/:uid/points":                "uid",
	"/api/v1/customers/:uid/points/history":        "uid",
	"/api/v1/customers/:uid/points/spend":          "uid",
	"/api/v1/customers/:uid/reviews":               "uid",
	"/api/v1/customers/:uid/reviews/:review_uid":   "uid",
	"/api/v1/customers/:uid/wishlist":              "uid",
	"/api/v1/customers/:uid/wishlist/add":          "uid",
	"/api/v1/customers/:uid/wishlist/remove":       "uid",
}

var adminOnlyRoutes = map[string]bool{
	"/api/v1/customers":                  true,
	"/api/v1/customers/:uid/points/earn": true,
}

func GinOwnershipMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if adminOnlyRoutes[c.FullPath()] {
			if !requireAdmin(c) {
				return
			}
			c.Next()
			return
		}
		param, gated := ownerParamByRoute[c.FullPath()]
		if !gated {
			c.Next()
			return
		}
		if !requireOwnerOrAdmin(c, c.Param(param)) {
			return
		}
		c.Next()
	}
}

func requireAdmin(c *gin.Context) bool {
	p := PrincipalFromContext(c)
	if p == nil {
		abortWithProblem(c, http.StatusUnauthorized, "authentication required")
		return false
	}
	if p.HasRole(RoleAdmin) {
		return true
	}
	abortWithProblem(c, http.StatusForbidden, "forbidden")
	return false
}

func requireOwnerOrAdmin(c *gin.Context, targetUid string) bool {
	p := PrincipalFromContext(c)
	if p == nil {
		abortWithProblem(c, http.StatusUnauthorized, "authentication required")
		return false
	}
	if p.Sub == targetUid || p.HasRole(RoleAdmin) {
		return true
	}
	abortWithProblem(c, http.StatusForbidden, "forbidden")
	return false
}
