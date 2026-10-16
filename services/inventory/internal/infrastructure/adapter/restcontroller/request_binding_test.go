package restcontroller

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/codejsha/bookstore-microservices/inventory/generated/application/port/openapi"
)

type recordingWarehouseApi struct {
	called   bool
	capacity int32
}

func (a *recordingWarehouseApi) WarehousesGetAll(
	_ context.Context,
	_ openapi.WarehouseApiCredentials,
	_ *string,
	_ *int32,
	_ *int32,
	_ *string,
) (openapi.WarehousesGetAllResponse, error) {
	a.called = true
	return openapi.WarehousesGetAllOk{Body: openapi.WarehouseFindAllResponse{}}, nil
}

func (a *recordingWarehouseApi) WarehousesCreate(
	_ context.Context,
	_ openapi.WarehouseApiCredentials,
	req openapi.WarehouseCreateRequest,
) (openapi.WarehousesCreateResponse, error) {
	a.called = true
	a.capacity = req.Capacity
	return openapi.WarehousesCreateCreated{Location: "/api/v1/warehouses/" + uidWarehouse1}, nil
}

func (a *recordingWarehouseApi) WarehousesRead(
	_ context.Context,
	_ openapi.WarehouseApiCredentials,
	uid string,
) (openapi.WarehousesReadResponse, error) {
	a.called = true
	return openapi.WarehousesReadOk{Body: openapi.WarehouseFindResponse{Uid: uid}}, nil
}

func (a *recordingWarehouseApi) WarehousesUpdate(
	_ context.Context,
	_ openapi.WarehouseApiCredentials,
	uid string,
	_ openapi.WarehouseUpdateRequest,
) (openapi.WarehousesUpdateResponse, error) {
	a.called = true
	return openapi.WarehousesUpdateOk{Body: openapi.WarehouseUpdateResponse{Uid: uid}}, nil
}

type recordingAuditApi struct {
	called         bool
	actualQuantity int32
}

func (a *recordingAuditApi) AuditsGetAll(
	_ context.Context,
	_ openapi.AuditApiCredentials,
	_ *string,
	_ *openapi.AuditStatus,
	_ *int32,
	_ *int32,
	_ *string,
) (openapi.AuditsGetAllResponse, error) {
	a.called = true
	return openapi.AuditsGetAllOk{Body: openapi.AuditFindAllResponse{}}, nil
}

func (a *recordingAuditApi) AuditsCreate(
	_ context.Context,
	_ openapi.AuditApiCredentials,
	req openapi.AuditCreateRequest,
) (openapi.AuditsCreateResponse, error) {
	a.called = true
	if len(req.Items) > 0 {
		a.actualQuantity = req.Items[0].ActualQuantity
	}
	return openapi.AuditsCreateOk{Body: openapi.AuditFindResponse{}}, nil
}

func (a *recordingAuditApi) AuditsRead(
	_ context.Context,
	_ openapi.AuditApiCredentials,
	uid string,
) (openapi.AuditsReadResponse, error) {
	a.called = true
	return openapi.AuditsReadOk{Body: openapi.AuditFindResponse{Uid: uid}}, nil
}

func (a *recordingAuditApi) AuditsComplete(
	_ context.Context,
	_ openapi.AuditApiCredentials,
	uid string,
) (openapi.AuditsCompleteResponse, error) {
	a.called = true
	return openapi.AuditsCompleteOk{Body: openapi.AuditFindResponse{Uid: uid}}, nil
}

type recordingClosingApi struct {
	called bool
}

func (a *recordingClosingApi) ClosingsGetAll(
	_ context.Context,
	_ openapi.ClosingApiCredentials,
	_ *string,
	_ *int32,
	_ *int32,
	_ *openapi.ClosingStatus,
	_ *int32,
	_ *int32,
	_ *string,
) (openapi.ClosingsGetAllResponse, error) {
	a.called = true
	return openapi.ClosingsGetAllOk{Body: openapi.ClosingFindAllResponse{}}, nil
}

func (a *recordingClosingApi) ClosingsCreate(
	_ context.Context,
	_ openapi.ClosingApiCredentials,
	_ openapi.ClosingCreateRequest,
) (openapi.ClosingsCreateResponse, error) {
	a.called = true
	return openapi.ClosingsCreateOk{Body: openapi.ClosingFindResponse{}}, nil
}

func (a *recordingClosingApi) ClosingsRead(
	_ context.Context,
	_ openapi.ClosingApiCredentials,
	uid string,
) (openapi.ClosingsReadResponse, error) {
	a.called = true
	return openapi.ClosingsReadOk{Body: openapi.ClosingFindResponse{Uid: uid}}, nil
}

func postJSON(r http.Handler, path, body string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-token")
	r.ServeHTTP(w, req)
	return w
}

func assertProblemDetail(t *testing.T, w *httptest.ResponseRecorder, detail string) {
	t.Helper()
	var problem openapi.Problem
	if err := json.Unmarshal(w.Body.Bytes(), &problem); err != nil {
		t.Fatalf("body is not a problem: %v: %s", err, w.Body.String())
	}
	if problem.Detail != detail {
		t.Errorf("detail = %q, want %q", problem.Detail, detail)
	}
}

func postWarehouseCreate(t *testing.T, body string) (*httptest.ResponseRecorder, *recordingWarehouseApi) {
	t.Helper()

	gin.SetMode(gin.TestMode)
	api := &recordingWarehouseApi{}
	r := gin.New()
	openapi.RegisterWarehouseApiGinRoutes(r, api)

	return postJSON(r, "/api/v1/warehouses", body), api
}

func postAuditCreate(t *testing.T, body string) (*httptest.ResponseRecorder, *recordingAuditApi) {
	t.Helper()

	gin.SetMode(gin.TestMode)
	api := &recordingAuditApi{}
	r := gin.New()
	openapi.RegisterAuditApiGinRoutes(r, api)

	return postJSON(r, "/api/v1/audits", body), api
}

func TestWarehouseCreateBinding_capacityZero_created201(t *testing.T) {
	w, api := postWarehouseCreate(t, `{"name":"central","capacity":0}`)

	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201: %s", w.Code, w.Body.String())
	}
	if loc := w.Header().Get("Location"); loc != "/api/v1/warehouses/"+uidWarehouse1 {
		t.Errorf("Location = %q", loc)
	}
	if !api.called {
		t.Error("controller was not invoked for a zero capacity, which the contract allows")
	}
	if api.capacity != 0 {
		t.Errorf("capacity = %d, want 0", api.capacity)
	}
}

func TestWarehouseCreateBinding_WhenCapacityNegative_Returns400(t *testing.T) {
	w, api := postWarehouseCreate(t, `{"name":"central","capacity":-1}`)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", w.Code)
	}
	assertProblemDetail(t, w, "request validation failed")
	if api.called {
		t.Error("controller was invoked despite a below-minimum field")
	}
}

func TestWarehouseCreateBinding_WhenNameMissing_Returns400(t *testing.T) {
	w, api := postWarehouseCreate(t, `{"capacity":10}`)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", w.Code)
	}
	assertProblemDetail(t, w, "malformed request body")
	if api.called {
		t.Error("controller was invoked despite a missing required field")
	}
}

func TestWarehouseCreateBinding_WhenNameOverLength_Returns400(t *testing.T) {
	w, api := postWarehouseCreate(t, `{"name":"`+strings.Repeat("x", 256)+`","capacity":10}`)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", w.Code)
	}
	assertProblemDetail(t, w, "request validation failed")
	if !strings.Contains(w.Body.String(), "WarehouseCreateRequest.Name: failed max") {
		t.Errorf("body = %s, want a Name max violation", w.Body.String())
	}
	if api.called {
		t.Error("controller was invoked despite an over-length field")
	}
}

func TestWarehouseCreateBinding_bearerMissing_unauthorized401(t *testing.T) {
	gin.SetMode(gin.TestMode)
	api := &recordingWarehouseApi{}
	r := gin.New()
	openapi.RegisterWarehouseApiGinRoutes(r, api)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/warehouses", strings.NewReader(`{"name":"central","capacity":1}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", w.Code)
	}
	if api.called {
		t.Error("controller was invoked without credentials")
	}
}

func TestAuditCreateBinding_WhenActualQuantityZero_Returns200(t *testing.T) {
	body := `{"warehouse_uid":"` + uidWarehouse1 + `","items":[{"edition_uid":"` + uidEdition1 + `","actual_quantity":0}]}`
	w, api := postAuditCreate(t, body)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", w.Code, w.Body.String())
	}
	if !api.called {
		t.Error("controller was not invoked for a zero actual quantity, which a stock count may legitimately report")
	}
	if api.actualQuantity != 0 {
		t.Errorf("actual_quantity = %d, want 0", api.actualQuantity)
	}
}

func TestAuditCreateBinding_WhenItemEditionUidMissing_Returns400(t *testing.T) {
	body := `{"warehouse_uid":"` + uidWarehouse1 + `","items":[{"actual_quantity":3}]}`
	w, api := postAuditCreate(t, body)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", w.Code)
	}
	assertProblemDetail(t, w, "malformed request body")
	if api.called {
		t.Error("controller was invoked despite an item missing a required field")
	}
}

func TestClosingCreateBinding_WhenPeriodOutOfRange_Returns400(t *testing.T) {
	cases := map[string]struct {
		body   string
		detail string
	}{
		"whenMonthZero_returns400":        {`{"warehouse_uid":"` + uidWarehouse1 + `","year":2026,"month":0}`, "request validation failed"},
		"whenMonthAboveTwelve_returns400": {`{"warehouse_uid":"` + uidWarehouse1 + `","year":2026,"month":13}`, "request validation failed"},
		"whenYearBelowMin_returns400":     {`{"warehouse_uid":"` + uidWarehouse1 + `","year":1999,"month":6}`, "request validation failed"},
		"whenPeriodMissing_returns400":    {`{"warehouse_uid":"` + uidWarehouse1 + `"}`, "malformed request body"},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			api := &recordingClosingApi{}
			r := gin.New()
			openapi.RegisterClosingApiGinRoutes(r, api)

			w := postJSON(r, "/api/v1/closings", tc.body)

			if w.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want 400", w.Code)
			}
			assertProblemDetail(t, w, tc.detail)
			if api.called {
				t.Error("controller was invoked despite an out-of-range period")
			}
		})
	}
}
