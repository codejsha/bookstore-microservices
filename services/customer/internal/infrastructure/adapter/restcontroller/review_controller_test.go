package restcontroller

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/codejsha/bookstore-microservices/customer/generated/application/port/openapi"
	"github.com/codejsha/bookstore-microservices/customer/internal/application/usecase"
	"github.com/codejsha/bookstore-microservices/customer/internal/infrastructure/httpx"
)

type removeReviewUseCase struct {
	usecase.CustomerUseCase
	err error
}

func (u *removeReviewUseCase) RemoveReview(context.Context, string, string) error {
	return u.err
}

func deleteReview(t *testing.T, uc usecase.CustomerUseCase) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(httpx.GinResponseMapping())
	openapi.RegisterReviewApiGinRoutes(r, NewReviewController(uc))

	w := httptest.NewRecorder()
	req := httptest.NewRequest(
		http.MethodDelete,
		"/api/v1/customers/6f1e3a2c-0a1b-4c2d-8e3f-9a0b1c2d3e4f/reviews/7a2b3c4d-0a1b-4c2d-8e3f-9a0b1c2d3e4f",
		nil,
	)
	req.Header.Set("Authorization", "Bearer test-token")
	r.ServeHTTP(w, req)
	return w
}

func TestCustomerReviewsRemove_ReviewMissingOrForeign_Returns404ProblemDetails(t *testing.T) {
	w := deleteReview(t, &removeReviewUseCase{err: fmt.Errorf("remove review r-1: %w", usecase.ErrNotFound)})
	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", w.Code)
	}
	if body := w.Body.String(); body != `{"title":"Not Found","status":404,"detail":"resource not found"}` {
		t.Errorf("body = %q, want resource-not-found problem details", body)
	}
}

func TestCustomerReviewsRemove_ReviewOwnedByCaller_Returns204(t *testing.T) {
	w := deleteReview(t, &removeReviewUseCase{})
	WaitForSideEffects(context.Background())
	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204", w.Code)
	}
}
