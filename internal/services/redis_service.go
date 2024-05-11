// internal/services/redis_service.go

package services

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jeancarlosdanese/go-base-api/internal/db"
)

type RedisService struct {
	client *redis.Client
}

func NewRedisService() *RedisService {
	return &RedisService{
		client: db.GetRedisClient(),
	}
}

func (s *RedisService) Set(key string, value interface{}, expiration time.Duration) error {
	ctx := context.Background()
	return s.client.Set(ctx, key, value, expiration).Err()
}

func (s *RedisService) Get(key string) (string, error) {
	ctx := context.Background()
	return s.client.Get(ctx, key).Result()
}
