package keycloak

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"

	"github.com/codejsha/bookstore-microservices/commonlib-go/pkg/client"
	"github.com/codejsha/bookstore-microservices/identity/internal/application/port/idp"
	"github.com/codejsha/bookstore-microservices/identity/internal/config"
)

var _ idp.OpenIDConnectAPI = (*openIDConnectClient)(nil)

func NewOpenIDConnectClient(
	cfg *config.Config,
	restyClient *client.RestyClient,
) idp.OpenIDConnectAPI {
	return &openIDConnectClient{
		cfg:         cfg,
		restyClient: restyClient.Client,
	}
}

type openIDConnectClient struct {
	cfg         *config.Config
	restyClient *resty.Client
}

func (c openIDConnectClient) RealmsRealmProtocolOpenidConnectLogoutPost(ctx context.Context, params idp.LogoutWebParam) error {
	data := map[string]string{
		"client_id":     params.ClientId,
		"client_secret": params.ClientSecret,
		"refresh_token": params.RefreshToken,
	}
	restyResp, err := c.restyClient.R().
		SetContext(ctx).
		SetDebug(c.cfg.App.Logging.IsDebug).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(data).
		Post(c.cfg.Keycloak.Url + fmt.Sprintf("/realms/%s/protocol/openid-connect/logout", params.Realm))
	if err != nil {
		return fmt.Errorf("failed to logout: %v", err)
	}
	if restyResp.StatusCode() != 204 {
		return fmt.Errorf("failed to logout: %s", restyResp.String())
	}

	return nil
}

func (c openIDConnectClient) RealmsRealmProtocolOpenidConnectTokenPost(ctx context.Context, params idp.TokenExchangeWebParam) (idp.TokenExchangeWebResp, error) {
	data := map[string]string{
		"client_id":     params.ClientId,
		"client_secret": params.ClientSecret,
		"grant_type":    string(params.GrantType),
		"username":      params.Username,
		"password":      params.Password,
		"refresh_token": params.RefreshToken,
		"scope":         params.Scope,
	}
	restyResp, err := c.restyClient.R().
		SetContext(ctx).
		SetDebug(c.cfg.App.Logging.IsDebug).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(data).
		Post(c.cfg.Keycloak.Url + fmt.Sprintf("/realms/%s/protocol/openid-connect/token", params.Realm))
	if err != nil {
		return idp.TokenExchangeWebResp{}, fmt.Errorf("failed to exchange token: %v", err)
	}

	var response idp.TokenExchangeWebResp
	err = json.Unmarshal(restyResp.Body(), &response)
	if err != nil {
		return idp.TokenExchangeWebResp{}, fmt.Errorf("failed to unmarshal token exchange response: %v", err)
	}
	return response, nil
}

func (c openIDConnectClient) RealmsRealmProtocolOpenidConnectTokenIntrospectPost(ctx context.Context, params idp.TokenIntrospectWebParam) (idp.TokenIntrospectWebResp, error) {
	data := map[string]string{
		"client_id":     params.ClientId,
		"client_secret": params.ClientSecret,
		"token":         params.Token,
	}
	restyResp, err := c.restyClient.R().
		SetContext(ctx).
		SetDebug(c.cfg.App.Logging.IsDebug).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(data).
		Post(c.cfg.Keycloak.Url + fmt.Sprintf("/realms/%s/protocol/openid-connect/token/introspect", params.Realm))
	if err != nil {
		return idp.TokenIntrospectWebResp{}, fmt.Errorf("failed to introspect token: %v", err)
	}

	var response idp.TokenIntrospectWebResp
	err = json.Unmarshal(restyResp.Body(), &response)
	if err != nil {
		return idp.TokenIntrospectWebResp{}, fmt.Errorf("failed to unmarshal token introspect response: %v", err)
	}
	return response, nil
}

func (c openIDConnectClient) RealmsRealmProtocolOpenidConnectUserinfoPost(ctx context.Context, params idp.UserinfoWebParam) (idp.UserinfoWebResp, error) {
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", params.AccessToken),
	}
	restyResp, err := c.restyClient.R().
		SetContext(ctx).
		SetDebug(c.cfg.App.Logging.IsDebug).
		SetHeaders(headers).
		Post(c.cfg.Keycloak.Url + fmt.Sprintf("/realms/%s/protocol/openid-connect/userinfo", params.Realm))
	if err != nil {
		return idp.UserinfoWebResp{}, fmt.Errorf("failed to get userinfo: %v", err)
	}

	var response idp.UserinfoWebResp
	err = json.Unmarshal(restyResp.Body(), &response)
	if err != nil {
		return idp.UserinfoWebResp{}, fmt.Errorf("failed to unmarshal userinfo response: %v", err)
	}
	return response, nil
}
