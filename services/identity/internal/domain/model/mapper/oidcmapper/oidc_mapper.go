package oidcmapper

import (
	"fmt"

	"github.com/codejsha/bookstore-microservices/commonlib-go/pkg/config"
	"github.com/codejsha/bookstore-microservices/commonlib-go/pkg/convert"
	"github.com/codejsha/bookstore-microservices/identity/internal/application/port/idp"
	"github.com/codejsha/bookstore-microservices/identity/internal/application/port/openapi"
)

func ToTokenExchangeParam(kcCfg *config.KeycloakConfig, value interface{}) (idp.TokenExchangeWebParam, error) {
	switch value.(type) {
	case openapi.SignInWebReq:
		r := value.(openapi.SignInWebReq)
		request := idp.TokenExchangeWebParam{
			Realm:        kcCfg.Realm,
			ClientId:     kcCfg.ClientId,
			ClientSecret: kcCfg.ClientSecret,
			GrantType:    idp.GRANTTYPE_PASSWORD,
			Username:     r.Email,
			Password:     r.Password,
			Scope:        "openid profile email",
		}
		return request, nil
	case openapi.TokenExchangeWebParam:
		r := value.(openapi.TokenExchangeWebParam)
		request := idp.TokenExchangeWebParam{
			Realm:        kcCfg.Realm,
			ClientId:     kcCfg.ClientId,
			ClientSecret: kcCfg.ClientSecret,
			GrantType:    idp.GRANTTYPE_REFRESH_TOKEN,
			RefreshToken: r.XRefreshToken,
			Scope:        "openid profile email",
		}
		return request, nil
	default:
		return idp.TokenExchangeWebParam{}, fmt.Errorf("failed to convert request to TokenExchangeParams: %v", value)
	}
}

func ToLogoutParam(kcCfg *config.KeycloakConfig, value interface{}) (idp.LogoutWebParam, error) {
	switch value.(type) {
	case openapi.SignOutWebParam:
		r := value.(openapi.SignOutWebParam)
		request := idp.LogoutWebParam{
			Realm:        kcCfg.Realm,
			ClientId:     kcCfg.ClientId,
			ClientSecret: kcCfg.ClientSecret,
			RefreshToken: r.XRefreshToken,
		}
		return request, nil
	default:
		return idp.LogoutWebParam{}, fmt.Errorf("failed to convert request to LogoutParams: %v", value)
	}
}

func ToUserRepresentation(value interface{}) (idp.UserRepresentation, error) {
	switch value.(type) {
	case openapi.SignUpWebReq:
		r := value.(openapi.SignUpWebReq)
		request := idp.UserRepresentation{
			Username:  convert.ToPtr(r.Email),
			FirstName: convert.ToPtr(r.FirstName),
			LastName:  convert.ToPtr(r.LastName),
			Email:     convert.ToPtr(r.Email),
			Enabled:   convert.ToPtr(true),
			Credentials: []idp.CredentialRepresentation{
				{
					Type:      convert.ToPtr("password"),
					Value:     convert.ToPtr(r.Password),
					Temporary: convert.ToPtr(false),
				},
			},
			RealmRoles: []string{
				string(openapi.AUTHROLE_PROFILE),
				string(openapi.AUTHROLE_ORDER),
				string(openapi.AUTHROLE_VIEW),
			},
		}
		return request, nil
	default:
		return idp.UserRepresentation{}, fmt.Errorf("failed to convert request to UserRepresentation: %v", value)
	}
}
