package support

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	RoleStaff  = "STAFF"
	RoleManage = "MANAGE"
	RoleSystem = "SYSTEM"
)

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

var staffOnlyRoutes = map[string]bool{
	"/api/v1/customers": true,
}

var manageOnlyRoutes = map[string]bool{
	"/api/v1/customers/:uid/points/earn": true,
}

func GinOwnershipMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if staffOnlyRoutes[c.FullPath()] {
			if !requireStaff(c) {
				return
			}
			c.Next()
			return
		}
		if manageOnlyRoutes[c.FullPath()] {
			if !requireManage(c) {
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
		if !requireOwnerOrStaff(c, c.Param(param)) {
			return
		}
		c.Next()
	}
}

func requireStaff(c *gin.Context) bool {
	p := PrincipalFromContext(c)
	if p == nil {
		abortWithProblem(c, http.StatusUnauthorized, "authentication required")
		return false
	}
	if hasAnyRole(p, RoleStaff, RoleSystem) {
		return true
	}
	abortWithProblem(c, http.StatusForbidden, "forbidden")
	return false
}

func requireManage(c *gin.Context) bool {
	p := PrincipalFromContext(c)
	if p == nil {
		abortWithProblem(c, http.StatusUnauthorized, "authentication required")
		return false
	}
	if hasAnyRole(p, RoleManage, RoleSystem) {
		return true
	}
	abortWithProblem(c, http.StatusForbidden, "forbidden")
	return false
}

func requireOwnerOrStaff(c *gin.Context, targetUid string) bool {
	p := PrincipalFromContext(c)
	if p == nil {
		abortWithProblem(c, http.StatusUnauthorized, "authentication required")
		return false
	}
	if p.Sub == targetUid || hasAnyRole(p, RoleStaff, RoleSystem) {
		return true
	}
	abortWithProblem(c, http.StatusForbidden, "forbidden")
	return false
}
