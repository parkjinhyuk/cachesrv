package repository

import (
	"cachesrv/internal/model"
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"time"
)

type CacheRepository interface {
	Set(ctx context.Context, key string, cache *model.Cache, ttl time.Duration) error
	Get(ctx context.Context, key string) (*model.Cache, error)
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

func (r *redisRepository) Set(ctx context.Context, key string, cache *model.Cache, ttl time.Duration) error {
	data, err := json.Marshal(cache)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, data, ttl).Err()
}

func (r *redisRepository) Get(ctx context.Context, key string) (*model.Cache, error) {
	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}

	var cache model.Cache
	if err := json.Unmarshal([]byte(data), &cache); err != nil {
		return nil, err
	}

	return &cache, nil
}
