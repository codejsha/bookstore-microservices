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

	"github.com/codejsha/bookstore-microservices/identity/internal/domain/service"
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

func TestGinResponseMapping_WhenHandlerWrapsUserAlreadyExists_Returns409ProblemDetails(t *testing.T) {
	r := newTestEngine(func(c *gin.Context) (int, any, error) {
		return 0, nil, MapError(c.Request.Context(),
			fmt.Errorf("%w with email a@b.c", service.ErrUserAlreadyExists))
	})
	w := do(r)
	if w.Code != http.StatusConflict {
		t.Fatalf("status = %d, want 409", w.Code)
	}
	if body := w.Body.String(); body != `{"title":"Conflict","status":409,"detail":"user already exists"}` {
		t.Errorf("body = %q, want conflict json", body)
	}
}

func TestGinResponseMapping_WhenHandlerWrapsElevatedRole_Returns400ProblemDetails(t *testing.T) {
	r := newTestEngine(func(c *gin.Context) (int, any, error) {
		return 0, nil, MapError(c.Request.Context(),
			fmt.Errorf("%w: %q", service.ErrElevatedRoleOnRegister, "MANAGE"))
	})
	w := do(r)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", w.Code)
	}
	if body := w.Body.String(); body != `{"title":"Bad Request","status":400,"detail":"registration may not grant elevated roles"}` {
		t.Errorf("body = %q, want bad-request json", body)
	}
}

func TestGinResponseMapping_WhenHandlerWrapsUniqueViolation_Returns409ProblemDetails(t *testing.T) {
	r := newTestEngine(func(c *gin.Context) (int, any, error) {
		return 0, nil, MapError(c.Request.Context(),
			fmt.Errorf("failed to sync user to local DB: %w",
				&pgconn.PgError{Code: "23505", ConstraintName: "uk_users_email"}))
	})
	w := do(r)
	if w.Code != http.StatusConflict {
		t.Fatalf("status = %d, want 409", w.Code)
	}
	if body := w.Body.String(); body != `{"title":"Conflict","status":409,"detail":"resource already exists"}` {
		t.Errorf("body = %q, want conflict json", body)
	}
}

func TestGinResponseMapping_WhenHandlerWrapsDuplicatedKey_Returns409ProblemDetails(t *testing.T) {
	r := newTestEngine(func(c *gin.Context) (int, any, error) {
		return 0, nil, MapError(c.Request.Context(),
			fmt.Errorf("upsert user: %w", gorm.ErrDuplicatedKey))
	})
	w := do(r)
	if w.Code != http.StatusConflict {
		t.Fatalf("status = %d, want 409", w.Code)
	}
	if body := w.Body.String(); body != `{"title":"Conflict","status":409,"detail":"resource already exists"}` {
		t.Errorf("body = %q, want conflict json", body)
	}
}

func TestGinResponseMapping_WhenHandlerWrapsGormNotFound_Returns404ProblemDetails(t *testing.T) {
	r := newTestEngine(func(c *gin.Context) (int, any, error) {
		return 0, nil, MapError(c.Request.Context(),
			fmt.Errorf("find user: %w", gorm.ErrRecordNotFound))
	})
	w := do(r)
	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", w.Code)
	}
	if body := w.Body.String(); body != `{"title":"Not Found","status":404,"detail":"resource not found"}` {
		t.Errorf("body = %q, want resource-not-found json", body)
	}
}

func TestGinResponseMapping_WhenHandlerFailsUnexpectedly_Returns500WithoutDetail(t *testing.T) {
	r := newTestEngine(func(c *gin.Context) (int, any, error) {
		return 0, nil, MapError(c.Request.Context(),
			errors.New(`failed to create user: 500 {"error":"internal keycloak detail"}`))
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
		c.JSON(http.StatusForbidden, gin.H{"message": "insufficient privileges"})
	})
	w := do(r)
	if w.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403", w.Code)
	}
	if body := w.Body.String(); body != `{"message":"insufficient privileges"}` {
		t.Errorf("body = %q, want untouched 403 body", body)
	}
}
