package support

import (
	"github.com/codejsha/bookstore-microservices/identity/internal/infrastructure/adapter/cache"
)

type UserCacheStore struct {
	*CacheClient
}

func NewUserCacheStore(c *CacheClient) cache.Store {
	return &UserCacheStore{CacheClient: c}
}
