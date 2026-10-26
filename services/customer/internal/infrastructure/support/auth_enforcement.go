package support

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	RoleStaff   = "STAFF"
	RoleManager = "MANAGER"
	RoleSystem  = "SYSTEM"
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

var managerOnlyRoutes = map[string]bool{
	"/api/v1/customers/:uid/points/earn": true,
}

var staffOnlyMethodRoutes = map[string]bool{
	"PATCH /api/v1/tickets/:uid/status": true,
}

var managerOnlyMethodRoutes = map[string]bool{
	"POST /api/v1/categories":        true,
	"PATCH /api/v1/categories/:uid":  true,
	"DELETE /api/v1/categories/:uid": true,
	"POST /api/v1/faqs":              true,
	"PATCH /api/v1/faqs/:uid":        true,
	"DELETE /api/v1/faqs/:uid":       true,
}

func GinOwnershipMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		methodRoute := c.Request.Method + " " + c.FullPath()
		if staffOnlyRoutes[c.FullPath()] || staffOnlyMethodRoutes[methodRoute] {
			if !requireStaff(c) {
				return
			}
			c.Next()
			return
		}
		if managerOnlyRoutes[c.FullPath()] || managerOnlyMethodRoutes[methodRoute] {
			if !requireManager(c) {
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

func requireManager(c *gin.Context) bool {
	p := PrincipalFromContext(c)
	if p == nil {
		abortWithProblem(c, http.StatusUnauthorized, "authentication required")
		return false
	}
	if hasAnyRole(p, RoleManager, RoleSystem) {
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
