package support

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"go.uber.org/fx"

	sharedconfig "github.com/codejsha/shared-library-go/pkg/config"

	"github.com/codejsha/bookstore-microservices/catalog/internal/infrastructure/adapter/cache"
	"github.com/codejsha/bookstore-microservices/catalog/internal/infrastructure/adapter/restcontroller"
)

func NewWorkListCache(lc fx.Lifecycle, cacheCfg *sharedconfig.CacheConfig) *cache.WorkListCache {
	if cacheCfg == nil || cacheCfg.Host == "" {
		logrus.Info("valkey cache disabled (no cache.host configured); work-list caching is off")
		return cache.NewWorkListCache(nil)
	}

	addr := fmt.Sprintf("%s:%d", cacheCfg.Host, cacheCfg.Port)
	pingCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	client, err := NewCacheClient(pingCtx, addr, cacheCfg.Password, 0)
	if err != nil {
		logrus.WithError(err).Warn("valkey unreachable; work-list caching disabled")
		return cache.NewWorkListCache(nil)
	}
	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			return client.Close()
		},
	})
	logrus.WithField("addr", addr).Info("valkey cache enabled for work-list read path")
	return cache.NewWorkListCache(client)
}

func RegisterCacheInvalidator(c *cache.WorkListCache) {
	restcontroller.SetCacheInvalidator(c)
}
