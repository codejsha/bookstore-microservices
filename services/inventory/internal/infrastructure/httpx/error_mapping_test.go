package httpx

import (
	"context"
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

func TestGinResponseMapping_WhenHandlerReturnsNotFound_Returns404ProblemDetails(t *testing.T) {
	r := newTestEngine(func(c *gin.Context) (int, any, error) {
		return 0, nil, MapNotFound(c.Request.Context(), ErrNotFound)
	})
	w := do(r)
	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", w.Code)
	}
	if body := w.Body.String(); body != `{"title":"Not Found","status":404,"detail":"resource not found"}` {
		t.Errorf("body = %q, want resource-not-found json", body)
	}
}

func TestGinResponseMapping_WhenHandlerWrapsGormNotFound_Returns404ProblemDetails(t *testing.T) {
	r := newTestEngine(func(c *gin.Context) (int, any, error) {
		return 0, nil, MapNotFound(c.Request.Context(), fmt.Errorf("load work: %w", gorm.ErrRecordNotFound))
	})
	w := do(r)
	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", w.Code)
	}
	if body := w.Body.String(); body != `{"title":"Not Found","status":404,"detail":"resource not found"}` {
		t.Errorf("body = %q, want resource-not-found json", body)
	}
}

func TestGinResponseMapping_WhenHandlerWrapsBusinessError_Returns4xxProblemDetails(t *testing.T) {
	cases := []struct {
		name       string
		err        error
		wantStatus int
		wantBody   string
	}{
		{
			name:       "whenQuantityInvalid_returns400",
			err:        fmt.Errorf("reserve quantity must be positive: got %d: %w", -1, repo.ErrInvalidQuantity),
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"title":"Bad Request","status":400,"detail":"invalid quantity"}`,
		},
		{
			name:       "whenStockInsufficient_returns400",
			err:        fmt.Errorf("insufficient stock for edition x: %w", repo.ErrInsufficientStock),
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"title":"Bad Request","status":400,"detail":"insufficient stock"}`,
		},
		{
			name:       "whenTransferWarehousesMatch_returns400",
			err:        fmt.Errorf("transfer must differ: %w", repo.ErrSameWarehouse),
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"title":"Bad Request","status":400,"detail":"source and target warehouse must differ"}`,
		},
		{
			name:       "whenStockNotFound_returns404",
			err:        fmt.Errorf("stock not found for edition e in warehouse w: %w", repo.ErrStockNotFound),
			wantStatus: http.StatusNotFound,
			wantBody:   `{"title":"Not Found","status":404,"detail":"stock not found"}`,
		},
		{
			name:       "whenWarehouseNotFound_returns404",
			err:        fmt.Errorf("warehouse w not found: %w", repo.ErrWarehouseNotFound),
			wantStatus: http.StatusNotFound,
			wantBody:   `{"title":"Not Found","status":404,"detail":"warehouse not found"}`,
		},
		{
			name:       "whenClosingDuplicated_returns409",
			err:        fmt.Errorf("closing already exists for warehouse x 2026-07: %w", repo.ErrDuplicateClosing),
			wantStatus: http.StatusConflict,
			wantBody:   `{"title":"Conflict","status":409,"detail":"monthly closing already exists for this warehouse and period"}`,
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

func TestGinResponseMapping_WhenBusinessErrorUnknown_Returns500WithoutDetail(t *testing.T) {
	r := newTestEngine(func(c *gin.Context) (int, any, error) {
		return 0, nil, MapBusinessError(c.Request.Context(), errors.New("some unexpected failure"))
	})
	w := do(r)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", w.Code)
	}
	if body := w.Body.String(); body != `{"title":"Internal Server Error","status":500,"detail":"internal server error"}` {
		t.Errorf("body = %q, want sanitized json", body)
	}
}

func TestGinResponseMapping_WhenHandlerFailsUnexpectedly_Returns500WithoutDetail(t *testing.T) {
	r := newTestEngine(func(c *gin.Context) (int, any, error) {
		return 0, nil, errors.New(`pq: duplicate key value violates unique constraint "users_email"`)
	})
	w := do(r)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", w.Code)
	}
	if body := w.Body.String(); body != `{"title":"Internal Server Error","status":500,"detail":"internal server error"}` {
		t.Errorf("body = %q, want sanitized json", body)
	}
}

func TestGinResponseMapping_WhenHandlerSucceeds_PassesResponseThrough(t *testing.T) {
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

func TestGinResponseMapping_WhenHandlerWrote4xx_PassesResponseThrough(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(GinResponseMapping())
	r.GET("/x", func(c *gin.Context) {
		c.JSON(http.StatusBadRequest, gin.H{"title": "Bad Request", "status": 400, "detail": "invalid size"})
	})
	w := do(r)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", w.Code)
	}
	if body := w.Body.String(); body != `{"detail":"invalid size","status":400,"title":"Bad Request"}` {
		t.Errorf("body = %q, want untouched 400 body", body)
	}
}

func TestGinResponseMapping_WhenStockNotFoundWrappedInGenericError_Returns404(t *testing.T) {
	r := newTestEngine(func(c *gin.Context) (int, any, error) {
		err := fmt.Errorf("complete transfer %s: %w", "t-1",
			fmt.Errorf("stock not found for edition e in warehouse w: %w", repo.ErrStockNotFound))
		return 0, nil, MapBusinessError(c.Request.Context(), err)
	})
	w := do(r)
	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", w.Code)
	}
	if body := w.Body.String(); body != `{"title":"Not Found","status":404,"detail":"stock not found"}` {
		t.Errorf("body = %q, want stock-not-found json", body)
	}
}

func TestMapBusinessError_WhenStockNotFoundWrapped_PreservesSentinel(t *testing.T) {
	in := fmt.Errorf("no stock for edition e in source warehouse w: %w", repo.ErrStockNotFound)
	out := MapBusinessError(context.Background(), in)
	if !errors.Is(out, repo.ErrStockNotFound) {
		t.Errorf("out = %v, want wrapped ErrStockNotFound", out)
	}
}

func TestMapBusinessError_WhenWarehouseNotFoundWrapped_PreservesSentinel(t *testing.T) {
	in := fmt.Errorf("target warehouse w not found: %w", repo.ErrWarehouseNotFound)
	out := MapBusinessError(context.Background(), in)
	if !errors.Is(out, repo.ErrWarehouseNotFound) {
		t.Errorf("out = %v, want wrapped ErrWarehouseNotFound", out)
	}
}
