package httpx

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/codejsha/bookstore-microservices/identity/internal/domain/service"
)

var ErrBadRequest = errors.New("bad request")

type callerUidKey struct{}

func WithCallerUid(ctx context.Context, uid string) context.Context {
	return context.WithValue(ctx, callerUidKey{}, uid)
}

func CallerUidFromContext(ctx context.Context) string {
	uid, _ := ctx.Value(callerUidKey{}).(string)
	return uid
}

type errorStatus struct {
	status  int
	message string
}

type errorStatusKey struct{}

func MapError(ctx context.Context, err error) error {
	if err == nil {
		return nil
	}
	switch {
	case errors.Is(err, ErrBadRequest):
		record(ctx, http.StatusBadRequest, err.Error())
	case errors.Is(err, service.ErrUserAlreadyExists):
		record(ctx, http.StatusConflict, "user already exists")
	case errors.Is(err, service.ErrElevatedRoleOnRegister):
		record(ctx, http.StatusBadRequest, "registration may not grant elevated roles")
	case errors.Is(err, gorm.ErrRecordNotFound):
		record(ctx, http.StatusNotFound, "resource not found")
	}
	return err
}

func record(ctx context.Context, status int, message string) {
	if es, ok := ctx.Value(errorStatusKey{}).(*errorStatus); ok {
		es.status = status
		es.message = message
	}
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
