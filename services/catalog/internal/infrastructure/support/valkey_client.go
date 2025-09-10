package support

import (
	"context"
	"encoding/json"
	"time"

	"github.com/valkey-io/valkey-go"
)

type ValkeyClient struct {
	client valkey.Client
}

func NewValkeyClient(addr string, password string, db int) (*ValkeyClient, error) {
	c, err := valkey.NewClient(valkey.ClientOption{
		InitAddress: []string{addr},
		Password:    password,
		SelectDB:    db,
	})
	if err != nil {
		return nil, err
	}
	return &ValkeyClient{client: c}, nil
}

func (c *ValkeyClient) Get(ctx context.Context, key string, dest interface{}) error {
	val, err := c.client.Do(ctx, c.client.B().Get().Key(key).Build()).ToString()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}

func (c *ValkeyClient) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	cmd := c.client.B().Set().Key(key).Value(string(data))
	if ttl > 0 {
		return c.client.Do(ctx, cmd.Ex(ttl).Build()).Error()
	}
	return c.client.Do(ctx, cmd.Build()).Error()
}

func (c *ValkeyClient) Incr(ctx context.Context, key string) (int64, error) {
	return c.client.Do(ctx, c.client.B().Incr().Key(key).Build()).ToInt64()
}

func (c *ValkeyClient) Del(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}
	return c.client.Do(ctx, c.client.B().Del().Key(keys...).Build()).Error()
}

func (c *ValkeyClient) Close() {
	c.client.Close()
}
