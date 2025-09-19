// internal/services/redis_service.go

package services

import (
	"context"
	"log"
	"time"

	"github.com/jeancarlosdanese/go-base-api/internal/db"
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
	// log.Printf("INFO: Setting key in Redis: %s", key)
	err := r.Client.Set(context.Background(), key, value, expiration).Err()
	if err != nil {
		log.Printf("ERROR: Error setting key in Redis: %v", err)
	}
	return err
}

func (r *RedisService) Get(key string) (string, error) {
	// log.Printf("INFO: Getting key from Redis: %s", key)
	result, err := r.Client.Get(context.Background(), key).Result()
	if err != nil {
		log.Printf("ERROR: Error getting key from Redis: %v", err)
	}
	return result, err
}
