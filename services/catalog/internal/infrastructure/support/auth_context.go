package support

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const (
	HeaderUserID     = "X-User-Id"
	HeaderUserEmail  = "X-User-Email"
	HeaderUserName   = "X-User-Name"
	HeaderUserRoles  = "X-User-Roles"
	HeaderUserScopes = "X-User-Scopes"
	HeaderJWTPayload = "X-Jwt-Payload"

	ContextKeyPrincipal = "auth.principal"
)

type Principal struct {
	Sub    string
	Email  string
	Name   string
	Roles  []string
	Scopes []string
	Raw    map[string]any
}

func (p *Principal) HasRole(role string) bool {
	for _, r := range p.Roles {
		if r == role {
			return true
		}
	}
	return false
}

func (p *Principal) HasScope(scope string) bool {
	for _, s := range p.Scopes {
		if s == scope {
			return true
		}
	}
	return false
}

func ParsePrincipal(c *gin.Context) (*Principal, error) {
	sub := c.GetHeader(HeaderUserID)
	if sub == "" {
		return nil, nil
	}
	p := &Principal{
		Sub:   sub,
		Email: c.GetHeader(HeaderUserEmail),
		Name:  c.GetHeader(HeaderUserName),
	}
	if v := c.GetHeader(HeaderUserRoles); v != "" {
		p.Roles = parseClaimList(v, ",")
	}
	if v := c.GetHeader(HeaderUserScopes); v != "" {
		p.Scopes = parseClaimList(v, " ")
	}
	if v := c.GetHeader(HeaderJWTPayload); v != "" {
		raw, err := decodePayload(v)
		if err != nil {
			return nil, err
		}
		p.Raw = raw
	}
	return p, nil
}

func decodePayload(v string) (map[string]any, error) {
	b, err := base64.RawURLEncoding.DecodeString(v)
	if err != nil {
		b2, err2 := base64.URLEncoding.DecodeString(v)
		if err2 != nil {
			return nil, fmt.Errorf("invalid x-jwt-payload: %w", err)
		}
		b = b2
	}
	var raw map[string]any
	if err := json.Unmarshal(b, &raw); err != nil {
		return nil, fmt.Errorf("invalid x-jwt-payload json: %w", err)
	}
	return raw, nil
}

func parseClaimList(v string, sep string) []string {
	if out, ok := decodeClaimArray(v); ok {
		return out
	}
	return splitTrimmed(v, sep)
}

func decodeClaimArray(v string) ([]string, bool) {
	b, err := base64.RawURLEncoding.DecodeString(v)
	if err != nil {
		b2, err2 := base64.URLEncoding.DecodeString(v)
		if err2 != nil {
			return nil, false
		}
		b = b2
	}
	var items []string
	if err := json.Unmarshal(b, &items); err != nil {
		return nil, false
	}
	return trimmed(items), true
}

func splitTrimmed(v string, sep string) []string {
	return trimmed(strings.Split(v, sep))
}

func trimmed(parts []string) []string {
	out := make([]string, 0, len(parts))
	for _, s := range parts {
		s = strings.TrimSpace(s)
		if s != "" {
			out = append(out, s)
		}
	}
	return out
}

func GinPrincipalMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		p, err := ParsePrincipal(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if p != nil {
			c.Set(ContextKeyPrincipal, p)
			trace.SpanFromContext(c.Request.Context()).SetAttributes(attribute.String("enduser.id", p.Sub))
		}
		c.Next()
	}
}

func PrincipalFromContext(c *gin.Context) *Principal {
	if v, ok := c.Get(ContextKeyPrincipal); ok {
		if p, ok := v.(*Principal); ok {
			return p
		}
	}
	return nil
}
