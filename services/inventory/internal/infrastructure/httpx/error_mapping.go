package httpx

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/codejsha/bookstore-microservices/inventory/internal/application/port/repo"
	"github.com/codejsha/bookstore-microservices/inventory/internal/domain/model/command"
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

func MapBusinessError(ctx context.Context, err error) error {
	if err == nil {
		return nil
	}
	es, ok := ctx.Value(errorStatusKey{}).(*errorStatus)
	if !ok {
		return err
	}
	switch {
	case errors.Is(err, command.ErrInvalidCommand):
		es.status, es.message = http.StatusBadRequest, err.Error()
	case errors.Is(err, repo.ErrInvalidQuantity):
		es.status, es.message = http.StatusBadRequest, repo.ErrInvalidQuantity.Error()
	case errors.Is(err, repo.ErrInsufficientStock):
		es.status, es.message = http.StatusBadRequest, repo.ErrInsufficientStock.Error()
	case errors.Is(err, repo.ErrSameWarehouse):
		es.status, es.message = http.StatusBadRequest, repo.ErrSameWarehouse.Error()
	case errors.Is(err, repo.ErrStockNotFound):
		es.status, es.message = http.StatusNotFound, repo.ErrStockNotFound.Error()
	case errors.Is(err, repo.ErrWarehouseNotFound):
		es.status, es.message = http.StatusNotFound, repo.ErrWarehouseNotFound.Error()
	case errors.Is(err, repo.ErrDuplicateClosing):
		es.status, es.message = http.StatusConflict, repo.ErrDuplicateClosing.Error()
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
				body = []byte(fmt.Sprintf(`{"title":%q,"status":%d,"detail":%q}`, http.StatusText(es.status), es.status, es.message))
				orig.Header().Set("Content-Type", "application/problem+json")
			case status >= http.StatusInternalServerError:
				body = []byte(fmt.Sprintf(`{"title":%q,"status":%d,"detail":"internal server error"}`, http.StatusText(status), status))
				orig.Header().Set("Content-Type", "application/problem+json")
			}
			orig.WriteHeader(status)
			_, _ = orig.Write(body)
		}()

		c.Next()
	}
}
