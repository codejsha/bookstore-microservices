package keycloak

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/codejsha/shared-library-go/pkg/rest/client"

	idp "github.com/codejsha/bookstore-microservices/identity/generated/application/port/idpapi"
	"github.com/codejsha/bookstore-microservices/identity/internal/application/port/security"
)

const (
	usersPath        = "/admin/realms/test/users"
	userPath         = "/admin/realms/test/users/idp-1"
	userLogoutPath   = "/admin/realms/test/users/idp-1/logout"
	realmRolesPath   = "/admin/realms/test/roles"
	userRolesPath    = "/admin/realms/test/users/idp-1/role-mappings/realm"
	wantBearerHeader = "Bearer " + stubAccessToken
)

var idpSentinels = []error{security.ErrIdpUserNotFound, security.ErrIdpUserConflict, security.ErrIdpBadRequest}

func newRecordingUsersClient(stub *keycloakStub) *usersClient {
	return &usersClient{
		cfg:         stubConfig(stub.URL),
		restyClient: newStubRestyClient(),
		tokenHelper: freshTokenHelper(),
	}
}

func assertNoIdpSentinel(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("err = nil, want an error")
	}
	for _, sentinel := range idpSentinels {
		if errors.Is(err, sentinel) {
			t.Fatalf("err = %v, want no %v", err, sentinel)
		}
	}
}

func roleNames(t *testing.T, body string) []string {
	t.Helper()
	var roles []idp.RoleRepresentation
	if err := json.Unmarshal([]byte(body), &roles); err != nil {
		t.Fatalf("decode role body %q: %v", body, err)
	}
	names := make([]string, 0, len(roles))
	for _, role := range roles {
		if role.Name == nil {
			t.Fatalf("role without a name in %q", body)
		}
		names = append(names, *role.Name)
	}
	sort.Strings(names)
	return names
}

func TestNewUsersClient_sharedRestyClient_callsKeycloakWithHelperToken(t *testing.T) {
	stub := newKeycloakStub(t).on(http.MethodDelete, userPath, http.StatusNoContent, "")
	c := NewUsersClient(stubConfig(stub.URL), &client.RestyClient{Client: newStubRestyClient()}, freshTokenHelper())

	if err := c.DeleteUser(context.Background(), stubRealm, "idp-1"); err != nil {
		t.Fatalf("delete user: %v", err)
	}

	if got := stub.only(t); got.Auth != wantBearerHeader {
		t.Errorf("authorization = %q, want %q", got.Auth, wantBearerHeader)
	}
}

func TestUsersClientListUsers_emailGiven_sendsExactMatchQueryAndDropsOtherEmails(t *testing.T) {
	stub := newKeycloakStub(t).on(http.MethodGet, usersPath, http.StatusOK,
		`[{"id":"idp-1","email":"Reader@Example.com"},{"id":"idp-2","email":"reader@example.com.au"},{"id":"idp-3"}]`)
	c := newRecordingUsersClient(stub)

	users, err := c.ListUsers(context.Background(), stubRealm, "reader@example.com")
	if err != nil {
		t.Fatalf("list users: %v", err)
	}

	got := stub.only(t)
	if got.Method != http.MethodGet || got.Path != usersPath {
		t.Errorf("request = %s %s, want GET %s", got.Method, got.Path, usersPath)
	}
	if got.Query.Get("email") != "reader@example.com" || got.Query.Get("exact") != "true" {
		t.Errorf("query = %v, want email=reader@example.com and exact=true", got.Query)
	}
	if got.Auth != wantBearerHeader {
		t.Errorf("authorization = %q, want %q", got.Auth, wantBearerHeader)
	}
	if len(users) != 1 || users[0].Id == nil || *users[0].Id != "idp-1" {
		t.Fatalf("users = %+v, want only the case-insensitive exact match idp-1", users)
	}
}

func TestUsersClientListUsers_emailOmitted_sendsNoFilterAndKeepsEveryUser(t *testing.T) {
	stub := newKeycloakStub(t).on(http.MethodGet, usersPath, http.StatusOK,
		`[{"id":"idp-1","email":"a@example.com"},{"id":"idp-2"}]`)
	c := newRecordingUsersClient(stub)

	users, err := c.ListUsers(context.Background(), stubRealm, "")
	if err != nil {
		t.Fatalf("list users: %v", err)
	}

	if got := stub.only(t); len(got.Query) != 0 {
		t.Errorf("query = %v, want none", got.Query)
	}
	if len(users) != 2 {
		t.Fatalf("users = %+v, want both", users)
	}
}

func TestUsersClientListUsers_keycloak500_reportsStatusWithoutSentinel(t *testing.T) {
	stub := newKeycloakStub(t).on(http.MethodGet, usersPath, http.StatusInternalServerError, keycloakErrorBody)
	c := newRecordingUsersClient(stub)

	_, err := c.ListUsers(context.Background(), stubRealm, "")

	assertNoIdpSentinel(t, err)
	if !strings.Contains(err.Error(), "status 500") {
		t.Errorf("err = %q, want the upstream status", err.Error())
	}
}

func TestUsersClientListUsers_malformedBody_reportsDecodeError(t *testing.T) {
	stub := newKeycloakStub(t).on(http.MethodGet, usersPath, http.StatusOK, `{"not":"a list"}`)
	c := newRecordingUsersClient(stub)

	users, err := c.ListUsers(context.Background(), stubRealm, "")

	assertNoIdpSentinel(t, err)
	if users != nil {
		t.Errorf("users = %+v, want nil", users)
	}
}

func TestUsersClientCreateUser_keycloak201_postsJsonRepresentation(t *testing.T) {
	stub := newKeycloakStub(t).on(http.MethodPost, usersPath, http.StatusCreated, "")
	c := newRecordingUsersClient(stub)
	username, email, enabled := "reader", "reader@example.com", true

	err := c.CreateUser(context.Background(), stubRealm, idp.UserRepresentation{
		Username: &username, Email: &email, Enabled: &enabled,
	})
	if err != nil {
		t.Fatalf("create user: %v", err)
	}

	got := stub.only(t)
	if got.Method != http.MethodPost || got.Path != usersPath {
		t.Errorf("request = %s %s, want POST %s", got.Method, got.Path, usersPath)
	}
	if got.ContentType != "application/json" {
		t.Errorf("content type = %q, want application/json", got.ContentType)
	}
	if got.Auth != wantBearerHeader {
		t.Errorf("authorization = %q, want %q", got.Auth, wantBearerHeader)
	}
	var sent map[string]any
	if err := json.Unmarshal([]byte(got.Body), &sent); err != nil {
		t.Fatalf("decode body %q: %v", got.Body, err)
	}
	want := map[string]any{"username": "reader", "email": "reader@example.com", "enabled": true}
	if !reflect.DeepEqual(sent, want) {
		t.Errorf("body = %v, want %v", sent, want)
	}
}

func TestUsersClientCreateUser_keycloak200_treatedAsFailure(t *testing.T) {
	stub := newKeycloakStub(t).on(http.MethodPost, usersPath, http.StatusOK, "")
	c := newRecordingUsersClient(stub)

	err := c.CreateUser(context.Background(), stubRealm, idp.UserRepresentation{})

	assertNoIdpSentinel(t, err)
}

func TestUsersClientGetUser_keycloak200_parsesRepresentation(t *testing.T) {
	stub := newKeycloakStub(t).on(http.MethodGet, userPath, http.StatusOK,
		`{"id":"idp-1","username":"reader","email":"reader@example.com","firstName":"Ada","enabled":false}`)
	c := newRecordingUsersClient(stub)

	user, err := c.GetUser(context.Background(), stubRealm, "idp-1")
	if err != nil {
		t.Fatalf("get user: %v", err)
	}

	got := stub.only(t)
	if got.Method != http.MethodGet || got.Path != userPath {
		t.Errorf("request = %s %s, want GET %s", got.Method, got.Path, userPath)
	}
	if got.Auth != wantBearerHeader {
		t.Errorf("authorization = %q, want %q", got.Auth, wantBearerHeader)
	}
	if user.Id == nil || *user.Id != "idp-1" ||
		user.Email == nil || *user.Email != "reader@example.com" ||
		user.FirstName == nil || *user.FirstName != "Ada" ||
		user.Enabled == nil || *user.Enabled {
		t.Fatalf("user = %+v, want the Keycloak representation", user)
	}
}

func TestUsersClientGetUser_malformedBody_reportsDecodeError(t *testing.T) {
	stub := newKeycloakStub(t).on(http.MethodGet, userPath, http.StatusOK, `[`)
	c := newRecordingUsersClient(stub)

	_, err := c.GetUser(context.Background(), stubRealm, "idp-1")

	assertNoIdpSentinel(t, err)
}

func TestUsersClientUpdateUser_keycloak204_putsJsonRepresentation(t *testing.T) {
	stub := newKeycloakStub(t).on(http.MethodPut, userPath, http.StatusNoContent, "")
	c := newRecordingUsersClient(stub)
	firstName, enabled := "Grace", false

	err := c.UpdateUser(context.Background(), stubRealm, "idp-1", idp.UserRepresentation{
		FirstName: &firstName, Enabled: &enabled,
	})
	if err != nil {
		t.Fatalf("update user: %v", err)
	}

	got := stub.only(t)
	if got.Method != http.MethodPut || got.Path != userPath {
		t.Errorf("request = %s %s, want PUT %s", got.Method, got.Path, userPath)
	}
	if got.ContentType != "application/json" {
		t.Errorf("content type = %q, want application/json", got.ContentType)
	}
	var sent map[string]any
	if err := json.Unmarshal([]byte(got.Body), &sent); err != nil {
		t.Fatalf("decode body %q: %v", got.Body, err)
	}
	want := map[string]any{"firstName": "Grace", "enabled": false}
	if !reflect.DeepEqual(sent, want) {
		t.Errorf("body = %v, want %v", sent, want)
	}
}

func TestUsersClientUpdateUser_keycloak409_wrapsErrIdpUserConflict(t *testing.T) {
	stub := newKeycloakStub(t).on(http.MethodPut, userPath, http.StatusConflict, keycloakErrorBody)
	c := newRecordingUsersClient(stub)

	err := c.UpdateUser(context.Background(), stubRealm, "idp-1", idp.UserRepresentation{})

	if !errors.Is(err, security.ErrIdpUserConflict) {
		t.Fatalf("err = %v, want ErrIdpUserConflict", err)
	}
}

func TestUsersClientLogoutUser_keycloak204_postsToLogoutEndpoint(t *testing.T) {
	stub := newKeycloakStub(t).on(http.MethodPost, userLogoutPath, http.StatusNoContent, "")
	c := newRecordingUsersClient(stub)

	if err := c.LogoutUser(context.Background(), stubRealm, "idp-1"); err != nil {
		t.Fatalf("logout user: %v", err)
	}

	got := stub.only(t)
	if got.Method != http.MethodPost || got.Path != userLogoutPath {
		t.Errorf("request = %s %s, want POST %s", got.Method, got.Path, userLogoutPath)
	}
	if got.Auth != wantBearerHeader {
		t.Errorf("authorization = %q, want %q", got.Auth, wantBearerHeader)
	}
}

func TestUsersClientLogoutUser_keycloak404_wrapsErrIdpUserNotFound(t *testing.T) {
	stub := newKeycloakStub(t).on(http.MethodPost, userLogoutPath, http.StatusNotFound, keycloakErrorBody)
	c := newRecordingUsersClient(stub)

	err := c.LogoutUser(context.Background(), stubRealm, "idp-1")

	if !errors.Is(err, security.ErrIdpUserNotFound) {
		t.Fatalf("err = %v, want ErrIdpUserNotFound", err)
	}
}

func TestUsersClientDeleteUser_keycloak204_sendsDeleteForThatUser(t *testing.T) {
	stub := newKeycloakStub(t).on(http.MethodDelete, userPath, http.StatusNoContent, "")
	c := newRecordingUsersClient(stub)

	if err := c.DeleteUser(context.Background(), stubRealm, "idp-1"); err != nil {
		t.Fatalf("delete user: %v", err)
	}

	got := stub.only(t)
	if got.Method != http.MethodDelete || got.Path != userPath {
		t.Errorf("request = %s %s, want DELETE %s", got.Method, got.Path, userPath)
	}
}

func TestUsersClientGetUserRealmRoles_defaultRoleAssigned_omitsDefaultAndUnnamedRoles(t *testing.T) {
	stub := newKeycloakStub(t).on(http.MethodGet, userRolesPath, http.StatusOK,
		`[{"id":"r1","name":"STAFF"},{"id":"r2","name":"default-roles-test"},{"id":"r3"},{"id":"r4","name":"MANAGER"}]`)
	c := newRecordingUsersClient(stub)

	roles, err := c.GetUserRealmRoles(context.Background(), stubRealm, "idp-1")
	if err != nil {
		t.Fatalf("get user realm roles: %v", err)
	}

	got := stub.only(t)
	if got.Method != http.MethodGet || got.Path != userRolesPath {
		t.Errorf("request = %s %s, want GET %s", got.Method, got.Path, userRolesPath)
	}
	if want := []string{"STAFF", "MANAGER"}; !reflect.DeepEqual(roles, want) {
		t.Fatalf("roles = %v, want %v", roles, want)
	}
}

func TestUsersClientGetUserRealmRoles_noRolesAssigned_emptyList(t *testing.T) {
	stub := newKeycloakStub(t).on(http.MethodGet, userRolesPath, http.StatusOK, `[]`)
	c := newRecordingUsersClient(stub)

	roles, err := c.GetUserRealmRoles(context.Background(), stubRealm, "idp-1")
	if err != nil {
		t.Fatalf("get user realm roles: %v", err)
	}

	if roles == nil || len(roles) != 0 {
		t.Fatalf("roles = %#v, want an empty non-nil list", roles)
	}
}

func TestUsersClientGetUserRealmRoles_malformedBody_reportsDecodeError(t *testing.T) {
	stub := newKeycloakStub(t).on(http.MethodGet, userRolesPath, http.StatusOK, `{}`)
	c := newRecordingUsersClient(stub)

	_, err := c.GetUserRealmRoles(context.Background(), stubRealm, "idp-1")

	assertNoIdpSentinel(t, err)
}

const realmRolesBody = `[{"id":"r-user","name":"USER"},{"id":"r-staff","name":"STAFF"},` +
	`{"id":"r-manager","name":"MANAGER"},{"id":"r-default","name":"default-roles-test"},{"id":"r-unnamed"}]`

func TestUsersClientSetUserRealmRoles_assignedRolesDiffer_grantsMissingAndRevokesSurplus(t *testing.T) {
	stub := newKeycloakStub(t).
		on(http.MethodGet, realmRolesPath, http.StatusOK, realmRolesBody).
		on(http.MethodGet, userRolesPath, http.StatusOK,
			`[{"id":"r-user","name":"USER"},{"id":"r-staff","name":"STAFF"},{"id":"r-default","name":"default-roles-test"},{"id":"r-unnamed"}]`).
		on(http.MethodPost, userRolesPath, http.StatusNoContent, "").
		on(http.MethodDelete, userRolesPath, http.StatusNoContent, "")
	c := newRecordingUsersClient(stub)

	err := c.SetUserRealmRoles(context.Background(), stubRealm, "idp-1", []string{"USER", "MANAGER"})
	if err != nil {
		t.Fatalf("set user realm roles: %v", err)
	}

	grants := stub.recordedFor(http.MethodPost, userRolesPath)
	if len(grants) != 1 {
		t.Fatalf("grant requests = %+v, want 1", grants)
	}
	if grants[0].ContentType != "application/json" || grants[0].Auth != wantBearerHeader {
		t.Errorf("grant headers = (%q, %q), want json with the bearer token", grants[0].ContentType, grants[0].Auth)
	}
	if got, want := roleNames(t, grants[0].Body), []string{"MANAGER"}; !reflect.DeepEqual(got, want) {
		t.Errorf("granted = %v, want %v", got, want)
	}
	if !strings.Contains(grants[0].Body, `"id":"r-manager"`) {
		t.Errorf("grant body = %s, want the realm role id Keycloak requires", grants[0].Body)
	}
	revokes := stub.recordedFor(http.MethodDelete, userRolesPath)
	if len(revokes) != 1 {
		t.Fatalf("revoke requests = %+v, want 1", revokes)
	}
	if got, want := roleNames(t, revokes[0].Body), []string{"STAFF"}; !reflect.DeepEqual(got, want) {
		t.Errorf("revoked = %v, want %v with the default role left alone", got, want)
	}
}

func TestUsersClientSetUserRealmRoles_assignedRolesAlreadyMatch_sendsNoMappingChange(t *testing.T) {
	stub := newKeycloakStub(t).
		on(http.MethodGet, realmRolesPath, http.StatusOK, realmRolesBody).
		on(http.MethodGet, userRolesPath, http.StatusOK,
			`[{"id":"r-staff","name":"STAFF"},{"id":"r-default","name":"default-roles-test"}]`)
	c := newRecordingUsersClient(stub)

	if err := c.SetUserRealmRoles(context.Background(), stubRealm, "idp-1", []string{"STAFF"}); err != nil {
		t.Fatalf("set user realm roles: %v", err)
	}

	for _, r := range stub.recorded() {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected %s %s, want read-only calls", r.Method, r.Path)
		}
	}
}

func TestUsersClientSetUserRealmRoles_emptyRoles_revokesEveryNonDefaultRole(t *testing.T) {
	stub := newKeycloakStub(t).
		on(http.MethodGet, realmRolesPath, http.StatusOK, realmRolesBody).
		on(http.MethodGet, userRolesPath, http.StatusOK,
			`[{"id":"r-user","name":"USER"},{"id":"r-staff","name":"STAFF"},{"id":"r-default","name":"default-roles-test"}]`).
		on(http.MethodDelete, userRolesPath, http.StatusNoContent, "")
	c := newRecordingUsersClient(stub)

	if err := c.SetUserRealmRoles(context.Background(), stubRealm, "idp-1", nil); err != nil {
		t.Fatalf("set user realm roles: %v", err)
	}

	if grants := stub.recordedFor(http.MethodPost, userRolesPath); len(grants) != 0 {
		t.Errorf("grant requests = %+v, want none", grants)
	}
	revokes := stub.recordedFor(http.MethodDelete, userRolesPath)
	if len(revokes) != 1 {
		t.Fatalf("revoke requests = %+v, want 1", revokes)
	}
	if got, want := roleNames(t, revokes[0].Body), []string{"STAFF", "USER"}; !reflect.DeepEqual(got, want) {
		t.Errorf("revoked = %v, want %v", got, want)
	}
}

func TestUsersClientSetUserRealmRoles_unknownRealmRole_wrapsErrIdpBadRequestBeforeAnyMapping(t *testing.T) {
	stub := newKeycloakStub(t).on(http.MethodGet, realmRolesPath, http.StatusOK, realmRolesBody)
	c := newRecordingUsersClient(stub)

	err := c.SetUserRealmRoles(context.Background(), stubRealm, "idp-1", []string{"STAFF", "SUPERUSER"})

	if !errors.Is(err, security.ErrIdpBadRequest) {
		t.Fatalf("err = %v, want ErrIdpBadRequest", err)
	}
	if !strings.Contains(err.Error(), "SUPERUSER") {
		t.Errorf("err = %q, want the missing role named", err.Error())
	}
	if got := stub.only(t); got.Path != realmRolesPath {
		t.Errorf("request = %s %s, want only the realm role listing", got.Method, got.Path)
	}
}

func TestUsersClientSetUserRealmRoles_userMissingInKeycloak_wrapsErrIdpUserNotFound(t *testing.T) {
	stub := newKeycloakStub(t).
		on(http.MethodGet, realmRolesPath, http.StatusOK, realmRolesBody).
		on(http.MethodGet, userRolesPath, http.StatusNotFound, keycloakErrorBody)
	c := newRecordingUsersClient(stub)

	err := c.SetUserRealmRoles(context.Background(), stubRealm, "idp-1", []string{"STAFF"})

	if !errors.Is(err, security.ErrIdpUserNotFound) {
		t.Fatalf("err = %v, want ErrIdpUserNotFound", err)
	}
}

func TestUsersClientSetUserRealmRoles_realmRoleListingForbidden_reportsStatusWithoutSentinel(t *testing.T) {
	stub := newKeycloakStub(t).on(http.MethodGet, realmRolesPath, http.StatusForbidden, keycloakErrorBody)
	c := newRecordingUsersClient(stub)

	err := c.SetUserRealmRoles(context.Background(), stubRealm, "idp-1", []string{"STAFF"})

	assertNoIdpSentinel(t, err)
	if !strings.Contains(err.Error(), "status 403") {
		t.Errorf("err = %q, want the upstream status", err.Error())
	}
}

func TestUsersClientSetUserRealmRoles_realmRoleListingMalformed_reportsDecodeError(t *testing.T) {
	stub := newKeycloakStub(t).on(http.MethodGet, realmRolesPath, http.StatusOK, `{}`)
	c := newRecordingUsersClient(stub)

	err := c.SetUserRealmRoles(context.Background(), stubRealm, "idp-1", []string{"STAFF"})

	assertNoIdpSentinel(t, err)
	if got := stub.only(t); got.Path != realmRolesPath {
		t.Errorf("request = %s %s, want only the realm role listing", got.Method, got.Path)
	}
}

func TestUsersClientSetUserRealmRoles_grantRejected_reportsGrantFailureAndSkipsRevoke(t *testing.T) {
	stub := newKeycloakStub(t).
		on(http.MethodGet, realmRolesPath, http.StatusOK, realmRolesBody).
		on(http.MethodGet, userRolesPath, http.StatusOK, `[{"id":"r-staff","name":"STAFF"}]`).
		on(http.MethodPost, userRolesPath, http.StatusForbidden, keycloakErrorBody).
		on(http.MethodDelete, userRolesPath, http.StatusNoContent, "")
	c := newRecordingUsersClient(stub)

	err := c.SetUserRealmRoles(context.Background(), stubRealm, "idp-1", []string{"MANAGER"})

	assertNoIdpSentinel(t, err)
	if !strings.Contains(err.Error(), "failed to grant realm roles") {
		t.Errorf("err = %q, want the grant step named", err.Error())
	}
	if revokes := stub.recordedFor(http.MethodDelete, userRolesPath); len(revokes) != 0 {
		t.Errorf("revoke requests = %+v, want none after a failed grant", revokes)
	}
}

func TestUsersClientSetUserRealmRoles_revokeRejected_reportsRevokeFailure(t *testing.T) {
	stub := newKeycloakStub(t).
		on(http.MethodGet, realmRolesPath, http.StatusOK, realmRolesBody).
		on(http.MethodGet, userRolesPath, http.StatusOK, `[{"id":"r-staff","name":"STAFF"}]`).
		on(http.MethodDelete, userRolesPath, http.StatusInternalServerError, keycloakErrorBody)
	c := newRecordingUsersClient(stub)

	err := c.SetUserRealmRoles(context.Background(), stubRealm, "idp-1", nil)

	assertNoIdpSentinel(t, err)
	if !strings.Contains(err.Error(), "failed to revoke realm roles") {
		t.Errorf("err = %q, want the revoke step named", err.Error())
	}
}

func usersClientCalls() map[string]func(c *usersClient) error {
	ctx := context.Background()
	return map[string]func(c *usersClient) error{
		"listUsers": func(c *usersClient) error {
			_, err := c.ListUsers(ctx, stubRealm, "reader@example.com")
			return err
		},
		"createUser": func(c *usersClient) error {
			return c.CreateUser(ctx, stubRealm, idp.UserRepresentation{})
		},
		"getUser": func(c *usersClient) error {
			_, err := c.GetUser(ctx, stubRealm, "idp-1")
			return err
		},
		"updateUser": func(c *usersClient) error {
			return c.UpdateUser(ctx, stubRealm, "idp-1", idp.UserRepresentation{})
		},
		"logoutUser": func(c *usersClient) error {
			return c.LogoutUser(ctx, stubRealm, "idp-1")
		},
		"deleteUser": func(c *usersClient) error {
			return c.DeleteUser(ctx, stubRealm, "idp-1")
		},
		"getUserRealmRoles": func(c *usersClient) error {
			_, err := c.GetUserRealmRoles(ctx, stubRealm, "idp-1")
			return err
		},
		"setUserRealmRoles": func(c *usersClient) error {
			return c.SetUserRealmRoles(ctx, stubRealm, "idp-1", []string{"STAFF"})
		},
	}
}

func TestUsersClient_keycloakUnreachable_reportsTransportErrorWithoutSentinel(t *testing.T) {
	for name, call := range usersClientCalls() {
		t.Run(name, func(t *testing.T) {
			c := &usersClient{
				cfg:         stubConfig(unreachableKeycloakURL(t)),
				restyClient: newStubRestyClient(),
				tokenHelper: freshTokenHelper(),
			}

			assertNoIdpSentinel(t, call(c))
		})
	}
}

func TestUsersClient_adminTokenUnavailable_failsBeforeCallingTheAdminApi(t *testing.T) {
	for name, call := range usersClientCalls() {
		t.Run(name, func(t *testing.T) {
			stub := newKeycloakStub(t)
			c := &usersClient{
				cfg:         stubConfig(stub.URL),
				restyClient: newStubRestyClient(),
				tokenHelper: &AdminTokenHelper{
					cfg:         stubConfig(unreachableKeycloakURL(t)),
					restyClient: newStubRestyClient(),
				},
			}

			err := call(c)

			assertNoIdpSentinel(t, err)
			if !strings.Contains(err.Error(), "failed to get token") {
				t.Errorf("err = %q, want the token failure named", err.Error())
			}
			if requests := stub.recorded(); len(requests) != 0 {
				t.Errorf("admin api requests = %+v, want none without a token", requests)
			}
		})
	}
}

func TestUsersClientMapRealmRoles_adminTokenUnavailable_failsBeforeCallingTheAdminApi(t *testing.T) {
	stub := newKeycloakStub(t)
	c := &usersClient{
		cfg:         stubConfig(stub.URL),
		restyClient: newStubRestyClient(),
		tokenHelper: &AdminTokenHelper{
			cfg:         stubConfig(unreachableKeycloakURL(t)),
			restyClient: newStubRestyClient(),
		},
	}

	err := c.mapRealmRoles(context.Background(), stubRealm, "idp-1", nil, http.MethodPost)

	assertNoIdpSentinel(t, err)
	if requests := stub.recorded(); len(requests) != 0 {
		t.Errorf("admin api requests = %+v, want none without a token", requests)
	}
}

func TestUsersClientListUsers_realmMissingInKeycloak_reportsStatusWithoutSentinel(t *testing.T) {
	stub := newKeycloakStub(t).on(http.MethodGet, usersPath, http.StatusNotFound, keycloakErrorBody)
	c := newRecordingUsersClient(stub)

	_, err := c.ListUsers(context.Background(), stubRealm, "")

	assertNoIdpSentinel(t, err)
	if !strings.Contains(err.Error(), "status 404") {
		t.Errorf("err = %q, want the upstream status", err.Error())
	}
}

func TestUsersClientCreateUser_realmMissingInKeycloak_reportsStatusWithoutSentinel(t *testing.T) {
	stub := newKeycloakStub(t).on(http.MethodPost, usersPath, http.StatusNotFound, keycloakErrorBody)
	c := newRecordingUsersClient(stub)

	err := c.CreateUser(context.Background(), stubRealm, idp.UserRepresentation{})

	assertNoIdpSentinel(t, err)
	if !strings.Contains(err.Error(), "status 404") {
		t.Errorf("err = %q, want the upstream status", err.Error())
	}
}

func TestUsersClientSetUserRealmRoles_realmMissingInKeycloak_reportsStatusWithoutSentinel(t *testing.T) {
	stub := newKeycloakStub(t).on(http.MethodGet, realmRolesPath, http.StatusNotFound, keycloakErrorBody)
	c := newRecordingUsersClient(stub)

	err := c.SetUserRealmRoles(context.Background(), stubRealm, "idp-1", []string{"STAFF"})

	assertNoIdpSentinel(t, err)
	if !strings.Contains(err.Error(), "status 404") {
		t.Errorf("err = %q, want the upstream status", err.Error())
	}
}
