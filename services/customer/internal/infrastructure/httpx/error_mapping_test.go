package httpx

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	"github.com/codejsha/bookstore-microservices/customer/internal/application/port/repo"
	"github.com/codejsha/bookstore-microservices/customer/internal/application/usecase"
)

func newTestEngine(controller func(c *gin.Context) (int, any, error)) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(GinResponseMapping())
	r.GET("/x", func(c *gin.Context) {
		status, body, err := controller(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(status, body)
	})
	return r
}

func do(r *gin.Engine) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	r.ServeHTTP(w, req)
	return w
}

func TestGinResponseMapping_NotFoundBecomes404(t *testing.T) {
	r := newTestEngine(func(c *gin.Context) (int, any, error) {
		return 0, nil, MapNotFound(c.Request.Context(), ErrNotFound)
	})
	w := do(r)
	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", w.Code)
	}
	if body := w.Body.String(); body != `{"error":"resource not found"}` {
		t.Errorf("body = %q, want resource-not-found json", body)
	}
}

func TestGinResponseMapping_GormNotFoundBecomes404(t *testing.T) {
	r := newTestEngine(func(c *gin.Context) (int, any, error) {
		return 0, nil, MapNotFound(c.Request.Context(), fmt.Errorf("load work: %w", gorm.ErrRecordNotFound))
	})
	w := do(r)
	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", w.Code)
	}
	if body := w.Body.String(); body != `{"error":"resource not found"}` {
		t.Errorf("body = %q, want resource-not-found json", body)
	}
}

func TestGinResponseMapping_NotImplementedBecomes501(t *testing.T) {
	r := newTestEngine(func(c *gin.Context) (int, any, error) {
		return 0, nil, MapNotImplemented(c.Request.Context(),
			fmt.Errorf("update customer: %w", usecase.ErrNotImplemented))
	})
	w := do(r)
	if w.Code != http.StatusNotImplemented {
		t.Fatalf("status = %d, want 501", w.Code)
	}
	if body := w.Body.String(); body != `{"error":"not implemented"}` {
		t.Errorf("body = %q, want not-implemented json", body)
	}
}

func TestGinResponseMapping_InsufficientPointsBecomes400(t *testing.T) {
	r := newTestEngine(func(c *gin.Context) (int, any, error) {
		return 0, nil, MapInsufficientPoints(c.Request.Context(),
			fmt.Errorf("change points: %w (no balance for user)", repo.ErrInsufficientPoints))
	})
	w := do(r)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", w.Code)
	}
	if body := w.Body.String(); body != `{"error":"insufficient points"}` {
		t.Errorf("body = %q, want insufficient-points json", body)
	}
}

func TestGinResponseMapping_ReviewExistsBecomes409(t *testing.T) {
	r := newTestEngine(func(c *gin.Context) (int, any, error) {
		return 0, nil, MapReviewExists(c.Request.Context(),
			fmt.Errorf("write review: %w", repo.ErrReviewExists))
	})
	w := do(r)
	if w.Code != http.StatusConflict {
		t.Fatalf("status = %d, want 409", w.Code)
	}
	if body := w.Body.String(); body != `{"error":"review already exists for this book"}` {
		t.Errorf("body = %q, want review-exists json", body)
	}
}

func TestGinResponseMapping_GrpcStatusBecomesHTTP(t *testing.T) {
	cases := []struct {
		name     string
		code     codes.Code
		wantCode int
		wantBody string
	}{
		{"not found", codes.NotFound, http.StatusNotFound, `{"error":"resource not found"}`},
		{"unavailable", codes.Unavailable, http.StatusServiceUnavailable, `{"error":"upstream service unavailable"}`},
		{"deadline exceeded", codes.DeadlineExceeded, http.StatusGatewayTimeout, `{"error":"upstream service timed out"}`},
		{"permission denied", codes.PermissionDenied, http.StatusForbidden, `{"error":"forbidden"}`},
		{"unauthenticated", codes.Unauthenticated, http.StatusUnauthorized, `{"error":"unauthenticated"}`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := newTestEngine(func(c *gin.Context) (int, any, error) {
				wrapped := fmt.Errorf("find user via grpc: %w", status.Error(tc.code, "boom"))
				return 0, nil, MapGrpcStatus(c.Request.Context(), wrapped)
			})
			w := do(r)
			if w.Code != tc.wantCode {
				t.Fatalf("status = %d, want %d", w.Code, tc.wantCode)
			}
			if body := w.Body.String(); body != tc.wantBody {
				t.Errorf("body = %q, want %q", body, tc.wantBody)
			}
		})
	}
}

func TestGinResponseMapping_GrpcInternalFallsThroughTo500(t *testing.T) {
	r := newTestEngine(func(c *gin.Context) (int, any, error) {
		wrapped := fmt.Errorf("list orders via grpc: %w", status.Error(codes.Internal, "kaboom"))
		return 0, nil, MapGrpcStatus(c.Request.Context(), wrapped)
	})
	w := do(r)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", w.Code)
	}
	if body := w.Body.String(); body != `{"error":"internal server error"}` {
		t.Errorf("body = %q, want sanitized json", body)
	}
}

func TestGinResponseMapping_NonStatusErrorFallsThroughTo500(t *testing.T) {
	r := newTestEngine(func(c *gin.Context) (int, any, error) {
		return 0, nil, MapGrpcStatus(c.Request.Context(), errors.New("some local failure"))
	})
	w := do(r)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", w.Code)
	}
}

func TestGinResponseMapping_InternalErrorIsSanitized(t *testing.T) {
	r := newTestEngine(func(c *gin.Context) (int, any, error) {
		return 0, nil, errors.New(`pq: duplicate key value violates unique constraint "users_email"`)
	})
	w := do(r)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", w.Code)
	}
	if body := w.Body.String(); body != `{"error":"internal server error"}` {
		t.Errorf("body = %q, want sanitized json", body)
	}
}

func TestGinResponseMapping_SuccessPassthrough(t *testing.T) {
	r := newTestEngine(func(c *gin.Context) (int, any, error) {
		return http.StatusOK, gin.H{"uid": "abc"}, nil
	})
	w := do(r)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	if body := w.Body.String(); body != `{"uid":"abc"}` {
		t.Errorf("body = %q, want passthrough json", body)
	}
}

func TestGinResponseMapping_ClientErrorPassthrough(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(GinResponseMapping())
	r.GET("/x", func(c *gin.Context) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid size"})
	})
	w := do(r)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", w.Code)
	}
	if body := w.Body.String(); body != `{"error":"invalid size"}` {
		t.Errorf("body = %q, want untouched 400 body", body)
	}
}
