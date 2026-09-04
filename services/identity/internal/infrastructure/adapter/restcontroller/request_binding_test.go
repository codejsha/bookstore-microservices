package restcontroller

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/codejsha/bookstore-microservices/identity/generated/application/port/openapi"
	"github.com/codejsha/bookstore-microservices/identity/internal/domain/aggregate"
	"github.com/codejsha/bookstore-microservices/identity/internal/domain/model/command"
)

func postUsersRegister(t *testing.T, body string) (*httptest.ResponseRecorder, *bool) {
	t.Helper()

	reached := false
	use := &stubUseCase{
		registerUser: func(_ context.Context, cmd command.UserRegisterCommand) (*aggregate.UserAggregate, error) {
			reached = true
			return newAggregate("idp-1", cmd.Email, cmd.FirstName, cmd.LastName, "", "ACTIVE"), nil
		},
	}

	gin.SetMode(gin.TestMode)
	handler := openapi.NewUserApiHandler(NewUserController(use))
	r := gin.New()
	r.POST("/api/v1/users", handler.UsersRegister)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	return w, &reached
}

func TestUsersRegisterBinding_WhenPasswordMissing_Returns400WithoutCallingUsecase(t *testing.T) {
	w, reached := postUsersRegister(t, `{"email":"u@x.com","first_name":"F","last_name":"L"}`)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", w.Code)
	}
	if *reached {
		t.Error("use case was invoked despite a missing required field")
	}
}

func TestUsersRegisterBinding_WhenEmailOverLength_Returns400WithoutCallingUsecase(t *testing.T) {
	email := strings.Repeat("a", 96) + "@x.com"
	w, reached := postUsersRegister(t, `{"email":"`+email+`","password":"secret","first_name":"F","last_name":"L"}`)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", w.Code)
	}
	if *reached {
		t.Error("use case was invoked despite an over-length field")
	}
}

func TestUsersRegisterBinding_WhenRequestValid_Returns200AndCallsUsecase(t *testing.T) {
	w, reached := postUsersRegister(t, `{"email":"u@x.com","password":"secret","first_name":"F","last_name":"L"}`)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	if !*reached {
		t.Error("use case was not invoked for a valid request")
	}
}
