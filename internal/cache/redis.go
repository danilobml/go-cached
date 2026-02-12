package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	rdb *redis.Client
}

func NewRedisClient(addr, password string, db int) *RedisClient {
	return &RedisClient{
		rdb: redis.NewClient(&redis.Options{
			Addr:         addr,
			Password:     password,
			DB:           db,
			ReadTimeout:  500 * time.Millisecond,
			WriteTimeout: 500 * time.Millisecond,
		}),
	}
}

func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	return r.rdb.Get(ctx, key).Result()
}

func (r *RedisClient) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return r.rdb.Set(ctx, key, value, ttl).Err()
}

func (r *RedisClient) Del(ctx context.Context, keys ...string) error {
	return r.rdb.Del(ctx, keys...).Err()
}

func (r *RedisClient) Ping(ctx context.Context) error {
	return r.rdb.Ping(ctx).Err()
}

func (r *RedisClient) Close() error {
	return r.rdb.Close()
}
