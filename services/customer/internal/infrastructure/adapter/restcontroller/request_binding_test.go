package restcontroller

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/codejsha/bookstore-microservices/customer/generated/application/port/openapi"
)

type recordingPointApi struct {
	called bool
}

func (a *recordingPointApi) PointsGetBalance(
	_ context.Context,
	_ openapi.PointApiCredentials,
	uid string,
) (openapi.PointsGetBalanceResponse, error) {
	a.called = true
	return openapi.PointsGetBalanceOk{Body: openapi.PointBalanceResponse{UserUid: uid}}, nil
}

func (a *recordingPointApi) PointsEarn(
	_ context.Context,
	_ openapi.PointApiCredentials,
	uid string,
	_ openapi.PointEarnRequest,
) (openapi.PointsEarnResponse, error) {
	a.called = true
	return openapi.PointsEarnOk{Body: openapi.PointBalanceResponse{UserUid: uid}}, nil
}

func (a *recordingPointApi) PointsHistory(
	_ context.Context,
	_ openapi.PointApiCredentials,
	_ string,
	_ *int32,
	_ *int32,
	_ *string,
) (openapi.PointsHistoryResponse, error) {
	a.called = true
	return openapi.PointsHistoryOk{Body: openapi.PointHistoryFindAllResponse{}}, nil
}

func (a *recordingPointApi) PointsSpend(
	_ context.Context,
	_ openapi.PointApiCredentials,
	uid string,
	_ openapi.PointSpendRequest,
) (openapi.PointsSpendResponse, error) {
	a.called = true
	return openapi.PointsSpendOk{Body: openapi.PointBalanceResponse{UserUid: uid}}, nil
}

func postPointsEarn(t *testing.T, body string) (*httptest.ResponseRecorder, *recordingPointApi) {
	t.Helper()

	gin.SetMode(gin.TestMode)
	api := &recordingPointApi{}

	r := gin.New()
	openapi.RegisterPointApiGinRoutes(r, api)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/customers/6f1e3a2c-0a1b-4c2d-8e3f-9a0b1c2d3e4f/points/earn",
		strings.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-token")
	r.ServeHTTP(w, req)

	return w, api
}

func TestPointsEarnBinding_WhenAmountMissing_Returns400WithoutCallingController(t *testing.T) {
	w, api := postPointsEarn(t, `{"reason":"promotion"}`)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", w.Code)
	}
	if api.called {
		t.Error("controller was invoked despite a missing required field")
	}
}

func TestPointsEarnBinding_WhenReasonOverLength_Returns400WithoutCallingController(t *testing.T) {
	w, api := postPointsEarn(t, `{"amount":10,"reason":"`+strings.Repeat("x", 256)+`"}`)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", w.Code)
	}
	if api.called {
		t.Error("controller was invoked despite an over-length field")
	}
}

func TestPointsEarnBinding_WhenRequestValid_Returns200AndCallsController(t *testing.T) {
	w, api := postPointsEarn(t, `{"amount":10,"reason":"promotion"}`)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	if !api.called {
		t.Error("controller was not invoked for a valid request")
	}
}
