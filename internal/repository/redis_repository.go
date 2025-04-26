package repository

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type CacheRepository interface {
	Set(ctx context.Context, key string, value string, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, error)
}

type redisRepository struct {
	client *redis.Client
}

func NewRedisRepository(addr string, password string, db int) CacheRepository {
	c := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &redisRepository{
		client: c,
	}
}

func (r *redisRepository) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

func (r *redisRepository) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}
