package restcontroller

import (
	"context"
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
	_ *string,
	_ *int32,
	_ *int32,
	_ *string,
) (*openapi.WarehouseFindAllResponse, error) {
	a.called = true
	return &openapi.WarehouseFindAllResponse{}, nil
}

func (a *recordingWarehouseApi) WarehousesCreate(_ context.Context, req openapi.WarehouseCreateRequest) error {
	a.called = true
	a.capacity = req.Capacity
	return nil
}

func (a *recordingWarehouseApi) WarehousesRead(_ context.Context, uid string) (*openapi.WarehouseFindResponse, error) {
	a.called = true
	return &openapi.WarehouseFindResponse{Uid: uid}, nil
}

func (a *recordingWarehouseApi) WarehousesUpdate(
	_ context.Context,
	uid string,
	_ openapi.WarehouseUpdateRequest,
) (*openapi.WarehouseUpdateResponse, error) {
	a.called = true
	return &openapi.WarehouseUpdateResponse{Uid: uid}, nil
}

type recordingAuditApi struct {
	called         bool
	actualQuantity int32
}

func (a *recordingAuditApi) AuditsGetAll(
	_ context.Context,
	_ *string,
	_ *openapi.AuditStatus,
	_ *int32,
	_ *int32,
	_ *string,
) (*openapi.AuditFindAllResponse, error) {
	a.called = true
	return &openapi.AuditFindAllResponse{}, nil
}

func (a *recordingAuditApi) AuditsCreate(
	_ context.Context,
	req openapi.AuditCreateRequest,
) (*openapi.AuditFindResponse, error) {
	a.called = true
	if len(req.Items) > 0 {
		a.actualQuantity = req.Items[0].ActualQuantity
	}
	return &openapi.AuditFindResponse{}, nil
}

func (a *recordingAuditApi) AuditsRead(_ context.Context, uid string) (*openapi.AuditFindResponse, error) {
	a.called = true
	return &openapi.AuditFindResponse{Uid: uid}, nil
}

func (a *recordingAuditApi) AuditsComplete(_ context.Context, uid string) (*openapi.AuditFindResponse, error) {
	a.called = true
	return &openapi.AuditFindResponse{Uid: uid}, nil
}

func postWarehouseCreate(t *testing.T, body string) (*httptest.ResponseRecorder, *recordingWarehouseApi) {
	t.Helper()

	gin.SetMode(gin.TestMode)
	api := &recordingWarehouseApi{}
	handler := openapi.NewWarehouseApiHandler(api)

	r := gin.New()
	r.POST("/api/v1/warehouses", handler.WarehousesCreate)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/warehouses", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	return w, api
}

func postAuditCreate(t *testing.T, body string) (*httptest.ResponseRecorder, *recordingAuditApi) {
	t.Helper()

	gin.SetMode(gin.TestMode)
	api := &recordingAuditApi{}
	handler := openapi.NewAuditApiHandler(api)

	r := gin.New()
	r.POST("/api/v1/audits", handler.AuditsCreate)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/audits", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	return w, api
}

func TestWarehouseCreateBinding_WhenCapacityZero_Returns200(t *testing.T) {
	w, api := postWarehouseCreate(t, `{"name":"central","capacity":0}`)

	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204: %s", w.Code, w.Body.String())
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
	if api.called {
		t.Error("controller was invoked despite a below-minimum field")
	}
}

func TestWarehouseCreateBinding_WhenNameMissing_Returns400(t *testing.T) {
	w, api := postWarehouseCreate(t, `{"capacity":10}`)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", w.Code)
	}
	if api.called {
		t.Error("controller was invoked despite a missing required field")
	}
}

func TestWarehouseCreateBinding_WhenNameOverLength_Returns400(t *testing.T) {
	w, api := postWarehouseCreate(t, `{"name":"`+strings.Repeat("x", 256)+`","capacity":10}`)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", w.Code)
	}
	if api.called {
		t.Error("controller was invoked despite an over-length field")
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
	if api.called {
		t.Error("controller was invoked despite an item missing a required field")
	}
}

func TestClosingCreateBinding_WhenPeriodOutOfRange_Returns400(t *testing.T) {
	var req openapi.ClosingCreateRequest

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/api/v1/closings", func(c *gin.Context) {
		if err := c.ShouldBindJSON(&req); err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		c.Status(http.StatusNoContent)
	})

	cases := map[string]string{
		"whenMonthZero_returns400":        `{"warehouse_uid":"` + uidWarehouse1 + `","year":2026,"month":0}`,
		"whenMonthAboveTwelve_returns400": `{"warehouse_uid":"` + uidWarehouse1 + `","year":2026,"month":13}`,
		"whenYearBelowMin_returns400":     `{"warehouse_uid":"` + uidWarehouse1 + `","year":1999,"month":6}`,
		"whenPeriodMissing_returns400":    `{"warehouse_uid":"` + uidWarehouse1 + `"}`,
	}
	for name, body := range cases {
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()
			httpReq := httptest.NewRequest(http.MethodPost, "/api/v1/closings", strings.NewReader(body))
			httpReq.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, httpReq)

			if w.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want 400", w.Code)
			}
		})
	}
}
