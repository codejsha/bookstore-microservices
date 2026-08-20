package support

import (
	"github.com/codejsha/bookstore-microservices/identity/internal/infrastructure/adapter/keycloak"
)

type AdminTokenCache struct {
	*CacheClient
}

func NewAdminTokenCache(c *CacheClient) keycloak.TokenCache {
	return &AdminTokenCache{CacheClient: c}
}
