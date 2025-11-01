// @file: internal/db/redis.go

package db

import (
	"context"
	"fmt"
	"os"
	"strconv" // Importe strconv para usar Atoi

	"github.com/jeancarlosdanese/go-base-api/internal/logging"
	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

// InitializeRedis configura e inicializa a conexão Redis.
func InitializeRedis() {
	addr := fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	redisDBStr := os.Getenv("REDIS_DB")      // Obtém a variável de ambiente como uma string
	redisDB, err := strconv.Atoi(redisDBStr) // Converte de string para int
	if err != nil {
		logging.Logger.Fatal().
			Err(err).
			Str("redis_db", redisDBStr).
			Str("operation", "init_redis").
			Msg("Não foi possível converter REDIS_DB para int")
	}

	redisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: os.Getenv("REDIS_PASSWORD"), // Senha, se necessário. Deixe vazio se não houver senha.
		DB:       redisDB,                     // Usa a conversão para definir o DB.
		PoolSize: 10,                          // Tamanho do pool de conexões.
	})

	ctx := context.Background()
	_, err = redisClient.Ping(ctx).Result()
	if err != nil {
		logging.Logger.Fatal().
			Err(err).
			Str("addr", addr).
			Int("db", redisDB).
			Str("operation", "init_redis").
			Msg("Não foi possível conectar ao Redis")
	}

	logging.Logger.Info().
		Str("addr", addr).
		Int("db", redisDB).
		Str("operation", "init_redis").
		Msg("Conexão com Redis estabelecida com sucesso")
}

// GetRedisClient retorna uma instância do cliente Redis.
func GetRedisClient() *redis.Client {
	return redisClient
}
