package httpx

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var ErrNotFound = errors.New("resource not found")

type errorStatus struct {
	status  int
	message string
}

type errorStatusKey struct{}

func MapNotFound(ctx context.Context, err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, ErrNotFound) || errors.Is(err, gorm.ErrRecordNotFound) {
		if es, ok := ctx.Value(errorStatusKey{}).(*errorStatus); ok {
			es.status = http.StatusNotFound
			es.message = "resource not found"
		}
	}
	return err
}

type captureWriter struct {
	gin.ResponseWriter
	buf    *bytes.Buffer
	status int
	wrote  bool
}

func (w *captureWriter) WriteHeader(code int)              { w.status = code; w.wrote = true }
func (w *captureWriter) WriteHeaderNow()                   {}
func (w *captureWriter) Write(b []byte) (int, error)       { return w.buf.Write(b) }
func (w *captureWriter) WriteString(s string) (int, error) { return w.buf.WriteString(s) }
func (w *captureWriter) Status() int                       { return w.status }
func (w *captureWriter) Size() int                         { return w.buf.Len() }
func (w *captureWriter) Written() bool                     { return w.wrote }

func GinResponseMapping() gin.HandlerFunc {
	return func(c *gin.Context) {
		es := &errorStatus{}
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), errorStatusKey{}, es))

		orig := c.Writer
		cw := &captureWriter{ResponseWriter: orig, buf: &bytes.Buffer{}, status: http.StatusOK}
		c.Writer = cw

		defer func() {
			c.Writer = orig
			if !cw.wrote {
				return
			}
			status := cw.status
			body := cw.buf.Bytes()
			switch {
			case es.status != 0:
				status = es.status
				body = []byte(fmt.Sprintf(`{"error":%q}`, es.message))
			case status >= http.StatusInternalServerError:
				body = []byte(`{"error":"internal server error"}`)
			}
			orig.WriteHeader(status)
			_, _ = orig.Write(body)
		}()

		c.Next()
	}
}
