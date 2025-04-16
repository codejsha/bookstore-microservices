package aggregate

import (
	"fmt"

	"github.com/codejsha/bookstore-microservices/identity/internal/application/port/idp"
	"github.com/codejsha/bookstore-microservices/identity/internal/application/port/openapi"
)

type TokenAggregate struct {
	SessionState string
	AccessToken  string
	RefreshToken string
	IdToken      string
	Scope        string
}

func (a *TokenAggregate) ToSignInWebResp() openapi.SignInWebResp {
	return openapi.SignInWebResp{
		AccessToken: a.AccessToken,
	}
}

func (a *TokenAggregate) ToTokenExchangeWebResp() openapi.TokenExchangeWebResp {
	return openapi.TokenExchangeWebResp{
		AccessToken: a.AccessToken,
	}

}

func (a *TokenAggregate) FromValue(value interface{}) (*TokenAggregate, error) {
	switch value.(type) {
	case idp.TokenExchangeWebResp:
		r := value.(idp.TokenExchangeWebResp)
		tokenAgg := &TokenAggregate{
			SessionState: r.SessionState,
			AccessToken:  r.AccessToken,
			RefreshToken: r.RefreshToken,
			IdToken:      *r.IdToken,
			Scope:        r.Scope,
		}
		return tokenAgg, nil
	default:
		return nil, fmt.Errorf("failed to convert request to TokenAggregate: %v", value)
	}
}
