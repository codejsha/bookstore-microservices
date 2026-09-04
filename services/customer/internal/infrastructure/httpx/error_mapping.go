package httpx

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	"github.com/codejsha/bookstore-microservices/customer/internal/application/port/repo"
	"github.com/codejsha/bookstore-microservices/customer/internal/application/usecase"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/model/command"
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
	if errors.Is(err, command.ErrInvalidCommand) {
		if es, ok := ctx.Value(errorStatusKey{}).(*errorStatus); ok {
			es.status = http.StatusBadRequest
			es.message = err.Error()
		}
	}
	return err
}

func MapNotImplemented(ctx context.Context, err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, usecase.ErrNotImplemented) {
		if es, ok := ctx.Value(errorStatusKey{}).(*errorStatus); ok {
			es.status = http.StatusNotImplemented
			es.message = "not implemented"
		}
	}
	return err
}

func MapInsufficientPoints(ctx context.Context, err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, repo.ErrInsufficientPoints) {
		if es, ok := ctx.Value(errorStatusKey{}).(*errorStatus); ok {
			es.status = http.StatusBadRequest
			es.message = "insufficient points"
		}
	}
	return err
}

func MapPointBalanceOverflow(ctx context.Context, err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, repo.ErrPointBalanceOverflow) {
		if es, ok := ctx.Value(errorStatusKey{}).(*errorStatus); ok {
			es.status = http.StatusBadRequest
			es.message = "point balance limit exceeded"
		}
	}
	return err
}

func MapReviewExists(ctx context.Context, err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, repo.ErrReviewExists) {
		if es, ok := ctx.Value(errorStatusKey{}).(*errorStatus); ok {
			es.status = http.StatusConflict
			es.message = "review already exists for this book"
		}
	}
	return err
}

var grpcStatusToHTTP = map[codes.Code]struct {
	status  int
	message string
}{
	codes.NotFound:         {http.StatusNotFound, "resource not found"},
	codes.Unavailable:      {http.StatusServiceUnavailable, "upstream service unavailable"},
	codes.DeadlineExceeded: {http.StatusGatewayTimeout, "upstream service timed out"},
	codes.PermissionDenied: {http.StatusForbidden, "forbidden"},
	codes.Unauthenticated:  {http.StatusUnauthorized, "unauthenticated"},
}

func grpcCodeFromError(err error) (codes.Code, bool) {
	for err != nil {
		if se, ok := err.(interface{ GRPCStatus() *status.Status }); ok {
			return se.GRPCStatus().Code(), true
		}
		err = errors.Unwrap(err)
	}
	return codes.OK, false
}

func MapGrpcStatus(ctx context.Context, err error) error {
	if err == nil {
		return nil
	}
	code, ok := grpcCodeFromError(err)
	if !ok {
		return err
	}
	mapped, ok := grpcStatusToHTTP[code]
	if !ok {
		return err
	}
	if es, ok := ctx.Value(errorStatusKey{}).(*errorStatus); ok {
		es.status = mapped.status
		es.message = mapped.message
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
