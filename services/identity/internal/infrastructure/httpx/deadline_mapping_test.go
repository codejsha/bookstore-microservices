package httpx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func newDeadlineEngine(timeout time.Duration, handler gin.HandlerFunc) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	})
	r.Use(GinResponseMapping())
	r.GET("/probe", handler)
	return r
}

func TestGinResponseMapping_serverErrorAfterRequestDeadline_mapsTo504ProblemDetails(t *testing.T) {
	r := newDeadlineEngine(time.Millisecond, func(c *gin.Context) {
		<-c.Request.Context().Done()
		c.JSON(http.StatusInternalServerError, gin.H{"error": c.Request.Context().Err().Error()})
	})

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/probe", nil))

	if rec.Code != http.StatusGatewayTimeout {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusGatewayTimeout)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/problem+json" {
		t.Fatalf("content-type = %q, want application/problem+json", ct)
	}
	if !strings.Contains(rec.Body.String(), `"status":504`) {
		t.Fatalf("body = %s, want problem details with status 504", rec.Body.String())
	}
}

func TestGinResponseMapping_serverErrorWithinRequestDeadline_staysInternalServerError(t *testing.T) {
	r := newDeadlineEngine(time.Minute, func(c *gin.Context) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "boom"})
	})

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/probe", nil))

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
	}
}

func TestGinResponseMapping_clientErrorAfterRequestDeadline_passesThrough(t *testing.T) {
	r := newDeadlineEngine(time.Millisecond, func(c *gin.Context) {
		<-c.Request.Context().Done()
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad"})
	})

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/probe", nil))

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}
