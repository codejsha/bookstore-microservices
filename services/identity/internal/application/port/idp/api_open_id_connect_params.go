// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.
//
// Keycloak OpenID Connect API
//
// Keycloak OpenID Connect API of Bookstore microservices
//
// API version: 0.1.0
// Contact: admin@example.com

package idp

// Logout
// Logout (operationId: RealmsRealmProtocolOpenidConnectLogoutPost)
type LogoutWebParam struct {
	Realm        string `uri:"realm"`
	ClientId     string `form:"client_id"`
	ClientSecret string `form:"client_secret"`
	RefreshToken string `form:"refresh_token"`
}

// Token introspection
// Token introspection (operationId: RealmsRealmProtocolOpenidConnectTokenIntrospectPost)
type TokenIntrospectWebParam struct {
	Realm        string `uri:"realm"`
	ClientId     string `form:"client_id"`
	ClientSecret string `form:"client_secret"`
	Token        string `form:"token"`
}

// Token exchange
// Token exchange (operationId: RealmsRealmProtocolOpenidConnectTokenPost)
type TokenExchangeWebParam struct {
	Realm        string    `uri:"realm"`
	ClientId     string    `form:"client_id"`
	ClientSecret string    `form:"client_secret"`
	GrantType    GrantType `form:"grant_type"`
	Username     string    `form:"username,omitempty"`
	Password     string    `form:"password,omitempty"`
	RefreshToken string    `form:"refresh_token,omitempty"`
	Scope        string    `form:"scope"`
}

// User info
// User info (operationId: RealmsRealmProtocolOpenidConnectUserinfoPost)
type UserinfoWebParam struct {
	Realm       string `uri:"realm"`
	AccessToken string `form:"access_token"`
}
