package support

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GinUuidParamMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, p := range c.Params {
			if isUidName(p.Key) && !isCanonicalUuid(p.Value) {
				abortWithProblem(c, http.StatusBadRequest, p.Key+" must be a valid uuid")
				return
			}
		}
		for key, values := range c.Request.URL.Query() {
			if !isUidName(key) {
				continue
			}
			for _, value := range values {
				if !isCanonicalUuid(value) {
					abortWithProblem(c, http.StatusBadRequest, key+" must be a valid uuid")
					return
				}
			}
		}
		c.Next()
	}
}

func isUidName(name string) bool {
	return strings.HasSuffix(strings.ToLower(name), "uid")
}

func isCanonicalUuid(value string) bool {
	if len(value) != 36 {
		return false
	}
	_, err := uuid.Parse(value)
	return err == nil
}
