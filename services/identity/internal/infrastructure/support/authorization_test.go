package support

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/codejsha/bookstore-microservices/identity/generated/application/port/openapi"
)

type syncUserApiStub struct {
	openapi.UserApi
	syncCalled bool
}

func (s *syncUserApiStub) UsersSyncFromIdp(_ context.Context, _ openapi.UserApiCredentials, idpUid openapi.Uuid) (openapi.UsersSyncFromIdpResponse, error) {
	s.syncCalled = true
	return openapi.UsersSyncFromIdpOk{Body: openapi.UserFindResponse{Uid: idpUid, Roles: []openapi.AuthRole{}}}, nil
}

func TestAuthorization_UserSync(t *testing.T) {
	const caller = "22222222-2222-2222-2222-222222222222"
	const other = "11111111-1111-1111-1111-111111111111"

	tests := []struct {
		name       string
		roles      string
		target     string
		wantStatus int
		wantCalled bool
	}{
		{"userSync_otherUid_forbidden", "USER", other, http.StatusForbidden, false},
		{"userSync_selfUid_allowed", "USER", caller, http.StatusOK, true},
		{"userSync_managerOnOtherUid_allowed", "MANAGER", other, http.StatusOK, true},
	}

	gin.SetMode(gin.TestMode)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine := gin.New()
			engine.Use(GinPrincipalMiddleware())
			engine.Use(GinAuthorizationMiddleware())
			stub := &syncUserApiStub{}
			openapi.RegisterUserApiGinRoutes(engine, stub)

			req := httptest.NewRequest(http.MethodPost, "/api/v1/users/sync/"+tt.target, nil)
			req.Header.Set("Authorization", "Bearer token")
			req.Header.Set(HeaderUserID, caller)
			req.Header.Set(HeaderUserRoles, tt.roles)
			rec := httptest.NewRecorder()
			engine.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d (body=%s)", rec.Code, tt.wantStatus, rec.Body.String())
			}
			if stub.syncCalled != tt.wantCalled {
				t.Fatalf("handler called = %v, want %v", stub.syncCalled, tt.wantCalled)
			}
		})
	}
}
