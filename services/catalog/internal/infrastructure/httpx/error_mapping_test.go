package httpx

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"

	"github.com/codejsha/bookstore-microservices/catalog/internal/application/port/repo"
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

func TestGinResponseMapping_WhenHandlerWrapsUniqueViolation_Returns409ProblemDetails(t *testing.T) {
	r := newTestEngine(func(c *gin.Context) (int, any, error) {
		pgErr := &pgconn.PgError{Code: "23505", ConstraintName: "uk_subject_name"}
		return 0, nil, MapConflict(c.Request.Context(), fmt.Errorf("create subject: %w", pgErr))
	})
	w := do(r)
	if w.Code != http.StatusConflict {
		t.Fatalf("status = %d, want 409", w.Code)
	}
	if body := w.Body.String(); body != `{"title":"Conflict","status":409,"detail":"resource already exists"}` {
		t.Errorf("body = %q, want already-exists json", body)
	}
}

func TestGinResponseMapping_WhenHandlerReturnsAlreadyExists_Returns409ProblemDetails(t *testing.T) {
	r := newTestEngine(func(c *gin.Context) (int, any, error) {
		return 0, nil, MapConflict(c.Request.Context(), repo.ErrAlreadyExists)
	})
	w := do(r)
	if w.Code != http.StatusConflict {
		t.Fatalf("status = %d, want 409", w.Code)
	}
}

func TestGinResponseMapping_WhenHandlerWrapsDuplicatedKey_Returns409ProblemDetails(t *testing.T) {
	r := newTestEngine(func(c *gin.Context) (int, any, error) {
		return 0, nil, MapConflict(c.Request.Context(), fmt.Errorf("create work: %w", gorm.ErrDuplicatedKey))
	})
	w := do(r)
	if w.Code != http.StatusConflict {
		t.Fatalf("status = %d, want 409", w.Code)
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
