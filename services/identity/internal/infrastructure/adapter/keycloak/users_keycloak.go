package keycloak

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-resty/resty/v2"

	"github.com/codejsha/shared-library-go/pkg/rest/client"

	idp "github.com/codejsha/bookstore-microservices/identity/generated/application/port/idpapi"
	"github.com/codejsha/bookstore-microservices/identity/internal/application/port/security"
	"github.com/codejsha/bookstore-microservices/identity/internal/config"
)

const defaultRolesPrefix = "default-roles-"

var _ security.UsersClient = (*usersClient)(nil)

type usersClient struct {
	cfg         *config.Config
	restyClient *resty.Client
	tokenHelper *AdminTokenHelper
}

func NewUsersClient(
	cfg *config.Config,
	restyClient *client.RestyClient,
	tokenHelper *AdminTokenHelper,
) security.UsersClient {
	return &usersClient{
		cfg:         cfg,
		tokenHelper: tokenHelper,
		restyClient: restyClient.Client,
	}
}

func (c *usersClient) authHeaders() (map[string]string, error) {
	token, _, err := c.tokenHelper.GetTokens()
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}
	return map[string]string{"Authorization": "Bearer " + token}, nil
}

func (c *usersClient) ListUsers(ctx context.Context, realm, email string) ([]idp.UserRepresentation, error) {
	headers, err := c.authHeaders()
	if err != nil {
		return nil, err
	}

	req := c.restyClient.R().
		SetContext(ctx).
		SetDebug(c.cfg.App.Logging.IsDebugEnabled).
		SetHeaders(headers)
	if email != "" {
		req.SetQueryParam("email", email)
		req.SetQueryParam("exact", "true")
	}

	resp, err := req.Get(c.cfg.Keycloak.Url + fmt.Sprintf("/admin/realms/%s/users", realm))
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to list users: %s", resp.String())
	}

	var users []idp.UserRepresentation
	if err := json.Unmarshal(resp.Body(), &users); err != nil {
		return nil, fmt.Errorf("failed to unmarshal users: %w", err)
	}
	if email != "" {
		filtered := users[:0]
		for _, u := range users {
			if u.Email != nil && strings.EqualFold(*u.Email, email) {
				filtered = append(filtered, u)
			}
		}
		users = filtered
	}
	return users, nil
}

func (c *usersClient) CreateUser(ctx context.Context, realm string, req idp.UserRepresentation) error {
	headers, err := c.authHeaders()
	if err != nil {
		return err
	}
	headers["Content-Type"] = "application/json"

	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal user: %w", err)
	}

	resp, err := c.restyClient.R().
		SetContext(ctx).
		SetDebug(c.cfg.App.Logging.IsDebugEnabled).
		SetHeaders(headers).
		SetBody(body).
		Post(c.cfg.Keycloak.Url + fmt.Sprintf("/admin/realms/%s/users", realm))
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	if resp.StatusCode() != 201 {
		return fmt.Errorf("failed to create user: %s", resp.String())
	}
	return nil
}

func (c *usersClient) GetUser(ctx context.Context, realm, userId string) (idp.UserRepresentation, error) {
	headers, err := c.authHeaders()
	if err != nil {
		return idp.UserRepresentation{}, err
	}

	resp, err := c.restyClient.R().
		SetContext(ctx).
		SetDebug(c.cfg.App.Logging.IsDebugEnabled).
		SetHeaders(headers).
		Get(c.cfg.Keycloak.Url + fmt.Sprintf("/admin/realms/%s/users/%s", realm, userId))
	if err != nil {
		return idp.UserRepresentation{}, fmt.Errorf("failed to get user: %w", err)
	}
	if resp.StatusCode() != 200 {
		return idp.UserRepresentation{}, fmt.Errorf("failed to get user: %s", resp.String())
	}

	var user idp.UserRepresentation
	if err := json.Unmarshal(resp.Body(), &user); err != nil {
		return idp.UserRepresentation{}, fmt.Errorf("failed to unmarshal user: %w", err)
	}
	return user, nil
}

func (c *usersClient) UpdateUser(ctx context.Context, realm, userId string, req idp.UserRepresentation) error {
	headers, err := c.authHeaders()
	if err != nil {
		return err
	}
	headers["Content-Type"] = "application/json"

	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal user: %w", err)
	}

	resp, err := c.restyClient.R().
		SetContext(ctx).
		SetDebug(c.cfg.App.Logging.IsDebugEnabled).
		SetHeaders(headers).
		SetBody(body).
		Put(c.cfg.Keycloak.Url + fmt.Sprintf("/admin/realms/%s/users/%s", realm, userId))
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	if resp.StatusCode() != 204 {
		return fmt.Errorf("failed to update user: %s", resp.String())
	}
	return nil
}

func (c *usersClient) LogoutUser(ctx context.Context, realm, userId string) error {
	headers, err := c.authHeaders()
	if err != nil {
		return err
	}

	resp, err := c.restyClient.R().
		SetContext(ctx).
		SetDebug(c.cfg.App.Logging.IsDebugEnabled).
		SetHeaders(headers).
		Post(c.cfg.Keycloak.Url + fmt.Sprintf("/admin/realms/%s/users/%s/logout", realm, userId))
	if err != nil {
		return fmt.Errorf("failed to logout user: %w", err)
	}
	if resp.StatusCode() != 204 {
		return fmt.Errorf("failed to logout user: %s", resp.String())
	}
	return nil
}

func (c *usersClient) DeleteUser(ctx context.Context, realm, userId string) error {
	headers, err := c.authHeaders()
	if err != nil {
		return err
	}

	resp, err := c.restyClient.R().
		SetContext(ctx).
		SetDebug(c.cfg.App.Logging.IsDebugEnabled).
		SetHeaders(headers).
		Delete(c.cfg.Keycloak.Url + fmt.Sprintf("/admin/realms/%s/users/%s", realm, userId))
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	if resp.StatusCode() != 204 {
		return fmt.Errorf("failed to delete user: %s", resp.String())
	}
	return nil
}

// ─── Realm role mappings ────────────────────────────────────────────────────

func (c *usersClient) GetUserRealmRoles(ctx context.Context, realm, userId string) ([]string, error) {
	assigned, err := c.getUserRealmRoles(ctx, realm, userId)
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(assigned))
	for _, role := range assigned {
		if role.Name == nil || strings.HasPrefix(*role.Name, defaultRolesPrefix) {
			continue
		}
		names = append(names, *role.Name)
	}
	return names, nil
}

func (c *usersClient) SetUserRealmRoles(ctx context.Context, realm, userId string, roles []string) error {
	available, err := c.listRealmRoles(ctx, realm)
	if err != nil {
		return err
	}

	desired := make([]idp.RoleRepresentation, 0, len(roles))
	for _, name := range roles {
		role, ok := available[name]
		if !ok {
			return fmt.Errorf("realm role %q does not exist in realm %s", name, realm)
		}
		desired = append(desired, role)
	}

	assigned, err := c.getUserRealmRoles(ctx, realm, userId)
	if err != nil {
		return err
	}

	wanted := make(map[string]bool, len(roles))
	for _, name := range roles {
		wanted[name] = true
	}

	var toAdd, toRemove []idp.RoleRepresentation
	for _, role := range desired {
		if !containsRole(assigned, *role.Name) {
			toAdd = append(toAdd, role)
		}
	}
	for _, role := range assigned {
		if role.Name == nil || strings.HasPrefix(*role.Name, defaultRolesPrefix) {
			continue
		}
		if !wanted[*role.Name] {
			toRemove = append(toRemove, role)
		}
	}
	if len(toAdd) > 0 {
		if err := c.mapRealmRoles(ctx, realm, userId, toAdd, http.MethodPost); err != nil {
			return fmt.Errorf("failed to grant realm roles: %w", err)
		}
	}
	if len(toRemove) > 0 {
		if err := c.mapRealmRoles(ctx, realm, userId, toRemove, http.MethodDelete); err != nil {
			return fmt.Errorf("failed to revoke realm roles: %w", err)
		}
	}
	return nil
}

func containsRole(roles []idp.RoleRepresentation, name string) bool {
	for _, role := range roles {
		if role.Name != nil && *role.Name == name {
			return true
		}
	}
	return false
}

func (c *usersClient) listRealmRoles(ctx context.Context, realm string) (map[string]idp.RoleRepresentation, error) {
	headers, err := c.authHeaders()
	if err != nil {
		return nil, err
	}

	resp, err := c.restyClient.R().
		SetContext(ctx).
		SetDebug(c.cfg.App.Logging.IsDebugEnabled).
		SetHeaders(headers).
		Get(c.cfg.Keycloak.Url + fmt.Sprintf("/admin/realms/%s/roles", realm))
	if err != nil {
		return nil, fmt.Errorf("failed to list realm roles: %w", err)
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to list realm roles: %s", resp.String())
	}

	var roles []idp.RoleRepresentation
	if err := json.Unmarshal(resp.Body(), &roles); err != nil {
		return nil, fmt.Errorf("failed to unmarshal realm roles: %w", err)
	}

	byName := make(map[string]idp.RoleRepresentation, len(roles))
	for _, role := range roles {
		if role.Name != nil {
			byName[*role.Name] = role
		}
	}
	return byName, nil
}

func (c *usersClient) getUserRealmRoles(ctx context.Context, realm, userId string) ([]idp.RoleRepresentation, error) {
	headers, err := c.authHeaders()
	if err != nil {
		return nil, err
	}

	resp, err := c.restyClient.R().
		SetContext(ctx).
		SetDebug(c.cfg.App.Logging.IsDebugEnabled).
		SetHeaders(headers).
		Get(c.cfg.Keycloak.Url + fmt.Sprintf("/admin/realms/%s/users/%s/role-mappings/realm", realm, userId))
	if err != nil {
		return nil, fmt.Errorf("failed to get user realm roles: %w", err)
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to get user realm roles: %s", resp.String())
	}

	var roles []idp.RoleRepresentation
	if err := json.Unmarshal(resp.Body(), &roles); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user realm roles: %w", err)
	}
	return roles, nil
}

func (c *usersClient) mapRealmRoles(
	ctx context.Context,
	realm, userId string,
	roles []idp.RoleRepresentation,
	method string,
) error {
	headers, err := c.authHeaders()
	if err != nil {
		return err
	}
	headers["Content-Type"] = "application/json"

	body, err := json.Marshal(roles)
	if err != nil {
		return fmt.Errorf("failed to marshal realm roles: %w", err)
	}

	resp, err := c.restyClient.R().
		SetContext(ctx).
		SetDebug(c.cfg.App.Logging.IsDebugEnabled).
		SetHeaders(headers).
		SetBody(body).
		Execute(method, c.cfg.Keycloak.Url+fmt.Sprintf("/admin/realms/%s/users/%s/role-mappings/realm", realm, userId))
	if err != nil {
		return fmt.Errorf("failed to map realm roles: %w", err)
	}
	if resp.StatusCode() != 204 {
		return fmt.Errorf("failed to map realm roles: %s", resp.String())
	}
	return nil
}
