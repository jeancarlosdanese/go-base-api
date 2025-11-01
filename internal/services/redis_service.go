// @file: internal/services/redis_service.go

package services

import (
	"context"
	"time"

	"github.com/jeancarlosdanese/go-base-api/internal/db"
	"github.com/jeancarlosdanese/go-base-api/internal/logging"
	"github.com/redis/go-redis/v9"
)

type RedisServiceInterface interface {
	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string) (string, error)
}

type RedisService struct {
	Client *redis.Client
}

func NewRedisService() *RedisService {
	return &RedisService{
		Client: db.GetRedisClient(),
	}
}

func (r *RedisService) Set(key string, value interface{}, expiration time.Duration) error {
	err := r.Client.Set(context.Background(), key, value, expiration).Err()
	if err != nil {
		logging.Logger.Error().
			Err(err).
			Str("key", key).
			Dur("expiration", expiration).
			Str("operation", "redis_set").
			Msg("Erro ao definir chave no Redis")
	}
	return err
}

func (r *RedisService) Get(key string) (string, error) {
	result, err := r.Client.Get(context.Background(), key).Result()
	if err != nil && err != redis.Nil {
		logging.Logger.Error().
			Err(err).
			Str("key", key).
			Str("operation", "redis_get").
			Msg("Erro ao obter chave do Redis")
	}
	return result, err
}
