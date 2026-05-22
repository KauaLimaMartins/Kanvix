package redis

import (
	"context"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

type Cache struct {
	Client *goredis.Client
}

func (c Cache) GetBytes(ctx context.Context, key string) ([]byte, bool, error) {
	b, err := c.Client.Get(ctx, key).Bytes()
	if err == nil {
		return b, true, nil
	}
	if err == goredis.Nil {
		return nil, false, nil
	}
	return nil, false, err
}

func (c Cache) SetBytes(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return c.Client.Set(ctx, key, value, ttl).Err()
}

