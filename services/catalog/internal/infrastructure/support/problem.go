package support

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

const problemMediaType = "application/problem+json"

// abortWithProblem stops the request and writes an RFC 9457 problem body.
func abortWithProblem(c *gin.Context, status int, detail string) {
	body := gin.H{"title": http.StatusText(status), "status": status, "detail": detail}
	payload, err := json.Marshal(body)
	if err != nil {
		c.AbortWithStatusJSON(status, body)
		return
	}
	c.Abort()
	c.Data(status, problemMediaType, payload)
}

// RegisterBindingTagNames makes the binding validator report wire (json) field
// names in violation messages instead of Go struct field names.
func RegisterBindingTagNames() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(field reflect.StructField) string {
			name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
	}
}
