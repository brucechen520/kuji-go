package redis

import (
	"time"

	"github.com/brucechen520/kuji-go/internal/pkg/core"

	"github.com/redis/go-redis/v9"
)

var _ TokenStore = (*tokenStore)(nil)

type TokenStore interface {
	Save(ctx core.Context, key string, value string, ttl time.Duration) error
	Get(ctx core.Context, key string) (string, error)
	Delete(ctx core.Context, keys ...string) error
}

type tokenStore struct {
	client *redis.Client
}

func NewTokenStore(client *redis.Client) TokenStore {
	return &tokenStore{client: client}
}

func (r *tokenStore) Save(ctx core.Context, key string, value string, ttl time.Duration) error {
	return r.client.Set(ctx.StdContext(), key, value, ttl).Err()
}

func (r *tokenStore) Get(ctx core.Context, key string) (string, error) {
	value, err := r.client.Get(ctx.StdContext(), key).Result()

	return value, err
}

func (r *tokenStore) Delete(ctx core.Context, keys ...string) error {
	return r.client.Del(ctx.StdContext(), keys...).Err()
}
