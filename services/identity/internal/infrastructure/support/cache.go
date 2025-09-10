package support

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"

	sharedconfig "github.com/codejsha/shared-library-go/pkg/config"
)

const lockKeyPrefix = "lock:"

var unlockScript = redis.NewScript(
	`if redis.call('GET', KEYS[1]) == ARGV[1] then return redis.call('DEL', KEYS[1]) else return 0 end`,
)

type CacheClient struct {
	client *redis.Client
}

func NewCacheClient(lc fx.Lifecycle, cacheCfg *sharedconfig.CacheConfig) *CacheClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cacheCfg.Host, cacheCfg.Port),
		Password: cacheCfg.Password,
	})
	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			return rdb.Close()
		},
	})
	return &CacheClient{client: rdb}
}

func (c *CacheClient) Get(ctx context.Context, key string, dest interface{}) error {
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}

func (c *CacheClient) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, data, ttl).Err()
}

func (c *CacheClient) Del(ctx context.Context, keys ...string) error {
	return c.client.Del(ctx, keys...).Err()
}

func (c *CacheClient) Exists(ctx context.Context, key string) (bool, error) {
	n, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

func (c *CacheClient) SAdd(ctx context.Context, key string, members ...string) error {
	vals := make([]interface{}, len(members))
	for i, m := range members {
		vals[i] = m
	}
	return c.client.SAdd(ctx, key, vals...).Err()
}

func (c *CacheClient) SRem(ctx context.Context, key string, members ...string) error {
	vals := make([]interface{}, len(members))
	for i, m := range members {
		vals[i] = m
	}
	return c.client.SRem(ctx, key, vals...).Err()
}

func (c *CacheClient) SMembers(ctx context.Context, key string) ([]string, error) {
	return c.client.SMembers(ctx, key).Result()
}

func (c *CacheClient) TryLock(ctx context.Context, key string, ttl time.Duration) (string, error) {
	token := uuid.NewString()
	acquired, err := c.client.SetNX(ctx, lockKeyPrefix+key, token, ttl).Result()
	if err != nil {
		return "", err
	}
	if !acquired {
		return "", nil
	}
	return token, nil
}

func (c *CacheClient) Unlock(ctx context.Context, key, token string) error {
	return unlockScript.Run(ctx, c.client, []string{lockKeyPrefix + key}, token).Err()
}
