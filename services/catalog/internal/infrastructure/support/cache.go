package support

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	sharedconfig "github.com/codejsha/shared-library-go/pkg/config"

	"github.com/codejsha/bookstore-microservices/catalog/internal/infrastructure/adapter/cache"
	"github.com/codejsha/bookstore-microservices/catalog/internal/infrastructure/adapter/restcontroller"
)

type CacheClientCloser struct {
	client *CacheClient
}

func (c *CacheClientCloser) Close() error {
	if c == nil || c.client == nil {
		return nil
	}
	return c.client.Close()
}

func NewWorkListCache(cacheCfg *sharedconfig.CacheConfig) (*cache.WorkListCache, *CacheClientCloser) {
	if cacheCfg == nil || cacheCfg.Host == "" {
		logrus.Info("valkey cache disabled (no cache.host configured); work-list caching is off")
		return cache.NewWorkListCache(nil), &CacheClientCloser{}
	}

	addr := fmt.Sprintf("%s:%d", cacheCfg.Host, cacheCfg.Port)
	pingCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	client, err := NewCacheClient(pingCtx, addr, cacheCfg.Password, 0)
	if err != nil {
		logrus.WithError(err).Warn("valkey unreachable; work-list caching disabled")
		return cache.NewWorkListCache(nil), &CacheClientCloser{}
	}
	logrus.WithField("addr", addr).Info("valkey cache enabled for work-list read path")
	return cache.NewWorkListCache(client), &CacheClientCloser{client: client}
}

func RegisterCacheInvalidator(c *cache.WorkListCache) {
	restcontroller.SetCacheInvalidator(c)
}
