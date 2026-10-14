package keycloak

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"

	sharedconfig "github.com/codejsha/shared-library-go/pkg/config"

	idp "github.com/codejsha/bookstore-microservices/identity/generated/application/port/idpapi"
	"github.com/codejsha/bookstore-microservices/identity/internal/application/port/security"
	"github.com/codejsha/bookstore-microservices/identity/internal/config"
)

const keycloakErrorBody = `{"errorMessage":"internal keycloak detail"}`

func TestIdpStatusError_mappedStatus_wrapsSentinelWithoutBody(t *testing.T) {
	cases := map[string]struct {
		status int
		want   error
	}{
		"notFound_userNotFound": {http.StatusNotFound, security.ErrIdpUserNotFound},
		"conflict_userConflict": {http.StatusConflict, security.ErrIdpUserConflict},
		"badRequest_badRequest": {http.StatusBadRequest, security.ErrIdpBadRequest},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			err := idpStatusError("get user", tc.status, keycloakErrorBody)

			if !errors.Is(err, tc.want) {
				t.Fatalf("err = %v, want %v", err, tc.want)
			}
			if strings.Contains(err.Error(), "internal keycloak detail") {
				t.Errorf("err = %q, want the Keycloak body left out", err.Error())
			}
		})
	}
}

func TestIdpStatusError_unmappedStatus_wrapsNoSentinel(t *testing.T) {
	err := idpStatusError("get user", http.StatusInternalServerError, keycloakErrorBody)

	for _, sentinel := range []error{security.ErrIdpUserNotFound, security.ErrIdpUserConflict, security.ErrIdpBadRequest} {
		if errors.Is(err, sentinel) {
			t.Errorf("err = %v, want no %v", err, sentinel)
		}
	}
}

func newStubUsersClient(t *testing.T, status int) *usersClient {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(status)
		_, _ = w.Write([]byte(keycloakErrorBody))
	}))
	t.Cleanup(srv.Close)
	return &usersClient{
		cfg: &config.Config{
			App:      &sharedconfig.AppConfig{},
			Keycloak: &sharedconfig.KeycloakConfig{Url: srv.URL, Realm: "test"},
		},
		restyClient: resty.New().SetTimeout(2 * time.Second),
		tokenHelper: &AdminTokenHelper{accessToken: "token", expiresAt: time.Now().Add(time.Hour)},
	}
}

func TestUsersClientGetUser_keycloak404_wrapsErrIdpUserNotFound(t *testing.T) {
	c := newStubUsersClient(t, http.StatusNotFound)

	_, err := c.GetUser(context.Background(), "test", "idp-1")

	if !errors.Is(err, security.ErrIdpUserNotFound) {
		t.Fatalf("err = %v, want ErrIdpUserNotFound", err)
	}
}

func TestUsersClientCreateUser_keycloak409_wrapsErrIdpUserConflict(t *testing.T) {
	c := newStubUsersClient(t, http.StatusConflict)

	err := c.CreateUser(context.Background(), "test", idp.UserRepresentation{})

	if !errors.Is(err, security.ErrIdpUserConflict) {
		t.Fatalf("err = %v, want ErrIdpUserConflict", err)
	}
}

func TestUsersClientUpdateUser_keycloak400_wrapsErrIdpBadRequest(t *testing.T) {
	c := newStubUsersClient(t, http.StatusBadRequest)

	err := c.UpdateUser(context.Background(), "test", "idp-1", idp.UserRepresentation{})

	if !errors.Is(err, security.ErrIdpBadRequest) {
		t.Fatalf("err = %v, want ErrIdpBadRequest", err)
	}
}

func TestUsersClientDeleteUser_keycloak404_wrapsErrIdpUserNotFound(t *testing.T) {
	c := newStubUsersClient(t, http.StatusNotFound)

	err := c.DeleteUser(context.Background(), "test", "idp-1")

	if !errors.Is(err, security.ErrIdpUserNotFound) {
		t.Fatalf("err = %v, want ErrIdpUserNotFound", err)
	}
}

func TestUsersClientGetUserRealmRoles_keycloak404_wrapsErrIdpUserNotFound(t *testing.T) {
	c := newStubUsersClient(t, http.StatusNotFound)

	_, err := c.GetUserRealmRoles(context.Background(), "test", "idp-1")

	if !errors.Is(err, security.ErrIdpUserNotFound) {
		t.Fatalf("err = %v, want ErrIdpUserNotFound", err)
	}
}
