package keycloak

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/codejsha/bookstore-microservices/commonlib-go/pkg/client"
	"github.com/codejsha/bookstore-microservices/identity/internal/application/port/idp"
	"github.com/codejsha/bookstore-microservices/identity/internal/config"
)

func NewAdminTokenHelper(
	cfg *config.Config,
	restyClient *client.RestyClient,
) *AdminTokenHelper {
	return &AdminTokenHelper{
		cfg:         cfg,
		restyClient: restyClient.Client,
	}
}

type AdminTokenHelper struct {
	cfg              *config.Config
	restyClient      *resty.Client
	accessToken      string
	refreshToken     string
	expiresAt        time.Time
	refreshExpiresAt time.Time
	mu               sync.Mutex
}

func (h *AdminTokenHelper) GetTokens() (string, string, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if time.Now().Before(h.expiresAt) {
		return h.accessToken, h.refreshToken, nil
	}

	var (
		accessToken, refreshToken   string
		expiresAt, refreshExpiresAt time.Time
		err                         error
	)

	if time.Now().Before(h.refreshExpiresAt) {
		accessToken, refreshToken, expiresAt, refreshExpiresAt, err = h.exchangeTokens()
	} else {
		accessToken, refreshToken, expiresAt, refreshExpiresAt, err = h.fetchTokens()
	}

	if err != nil {
		return "", "", fmt.Errorf("failed to refresh tokens: %v", err)
	}

	h.accessToken = accessToken
	h.refreshToken = refreshToken
	h.expiresAt = expiresAt
	h.refreshExpiresAt = refreshExpiresAt
	return h.accessToken, h.refreshToken, nil
}

func (h *AdminTokenHelper) fetchTokens() (string, string, time.Time, time.Time, error) {
	data := map[string]string{
		"client_id":     h.cfg.Keycloak.ClientId,
		"client_secret": h.cfg.Keycloak.ClientSecret,
		"grant_type":    "password",
		"username":      h.cfg.Keycloak.RealmAdminUsername,
		"password":      h.cfg.Keycloak.RealmAdminPassword,
		"scope":         "openid profile email",
	}
	restyResp, err := h.restyClient.R().
		SetDebug(h.cfg.App.Logging.IsDebug).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(data).
		Post(h.cfg.Keycloak.Url + fmt.Sprintf("/realms/%s/protocol/openid-connect/token", h.cfg.Keycloak.Realm))
	if err != nil {
		return "", "", time.Time{}, time.Time{}, fmt.Errorf("failed to fetch tokens: %v", err)
	}

	var response idp.TokenExchangeWebResp
	err = json.Unmarshal(restyResp.Body(), &response)
	if err != nil {
		return "", "", time.Time{}, time.Time{}, fmt.Errorf("failed to unmarshal token exchange response: %v", err)
	}

	now := time.Now()
	accessToken := response.AccessToken
	refreshToken := response.RefreshToken
	expiresAt := now.Add(time.Duration(response.ExpiresIn) * time.Second)
	refreshExpiresAt := now.Add(time.Duration(response.RefreshExpiresIn) * time.Second)
	return accessToken, refreshToken, expiresAt, refreshExpiresAt, nil
}

func (h *AdminTokenHelper) exchangeTokens() (string, string, time.Time, time.Time, error) {
	data := map[string]string{
		"client_id":     h.cfg.Keycloak.ClientId,
		"client_secret": h.cfg.Keycloak.ClientSecret,
		"grant_type":    "refresh_token",
		"refresh_token": h.refreshToken,
		"scope":         "openid profile email",
	}
	restyResp, err := h.restyClient.R().
		SetDebug(h.cfg.App.Logging.IsDebug).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(data).
		Post(h.cfg.Keycloak.Url + fmt.Sprintf("/realms/%s/protocol/openid-connect/token", h.cfg.Keycloak.Realm))
	if err != nil {
		return "", "", time.Time{}, time.Time{}, fmt.Errorf("failed to exchange tokens: %v", err)
	}

	var response idp.TokenExchangeWebResp
	err = json.Unmarshal(restyResp.Body(), &response)
	if err != nil {
		return "", "", time.Time{}, time.Time{}, fmt.Errorf("failed to unmarshal token exchange response: %v", err)
	}

	now := time.Now()
	accessToken := response.AccessToken
	refreshToken := response.RefreshToken
	expiresAt := now.Add(time.Duration(response.ExpiresIn) * time.Second)
	refreshExpiresAt := now.Add(time.Duration(response.RefreshExpiresIn) * time.Second)
	return accessToken, refreshToken, expiresAt, refreshExpiresAt, nil
}
