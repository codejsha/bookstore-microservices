package errorcode

import (
	"github.com/gin-gonic/gin"
)

func RenderResponse(c *gin.Context, code ErrorCode, err error) {
	c.JSON(code.StatusCode, gin.H{
		"code":    code.Code,
		"message": code.Message,
		"details": err.Error(),
	})
	return
}
