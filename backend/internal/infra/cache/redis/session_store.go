package redis

import (
	"context"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

type SessionStore struct {
	Client *goredis.Client
}

func (s SessionStore) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return s.Client.Set(ctx, key, value, ttl).Err()
}

func (s SessionStore) Get(ctx context.Context, key string) (string, bool, error) {
	v, err := s.Client.Get(ctx, key).Result()
	if err == nil {
		return v, true, nil
	}
	if err == goredis.Nil {
		return "", false, nil
	}
	return "", false, err
}

func (s SessionStore) Delete(ctx context.Context, key string) error {
	return s.Client.Del(ctx, key).Err()
}

