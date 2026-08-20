package httpx

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/codejsha/bookstore-microservices/inventory/internal/application/port/repo"
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

func TestGinResponseMapping_BusinessErrorsBecome4xx(t *testing.T) {
	cases := []struct {
		name       string
		err        error
		wantStatus int
		wantBody   string
	}{
		{
			name:       "invalid quantity → 400",
			err:        fmt.Errorf("reserve quantity must be positive: got %d: %w", -1, repo.ErrInvalidQuantity),
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"error":"invalid quantity"}`,
		},
		{
			name:       "insufficient stock → 400",
			err:        fmt.Errorf("insufficient stock for edition x: %w", repo.ErrInsufficientStock),
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"error":"insufficient stock"}`,
		},
		{
			name:       "same warehouse transfer → 400",
			err:        fmt.Errorf("transfer must differ: %w", repo.ErrSameWarehouse),
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"error":"source and target warehouse must differ"}`,
		},
		{
			name:       "duplicate closing → 409",
			err:        fmt.Errorf("closing already exists for warehouse x 2026-07: %w", repo.ErrDuplicateClosing),
			wantStatus: http.StatusConflict,
			wantBody:   `{"error":"monthly closing already exists for this warehouse and period"}`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := newTestEngine(func(c *gin.Context) (int, any, error) {
				return 0, nil, MapBusinessError(c.Request.Context(), tc.err)
			})
			w := do(r)
			if w.Code != tc.wantStatus {
				t.Fatalf("status = %d, want %d", w.Code, tc.wantStatus)
			}
			if body := w.Body.String(); body != tc.wantBody {
				t.Errorf("body = %q, want %q", body, tc.wantBody)
			}
		})
	}
}

func TestGinResponseMapping_UnknownBusinessErrorStays500(t *testing.T) {
	r := newTestEngine(func(c *gin.Context) (int, any, error) {
		return 0, nil, MapBusinessError(c.Request.Context(), errors.New("some unexpected failure"))
	})
	w := do(r)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", w.Code)
	}
	if body := w.Body.String(); body != `{"error":"internal server error"}` {
		t.Errorf("body = %q, want sanitized json", body)
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
