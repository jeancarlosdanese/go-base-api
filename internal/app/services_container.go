// @file: internal/app/services_container.go

package app

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	// Import correto
	"github.com/jeancarlosdanese/go-base-api/internal/db"
	"github.com/jeancarlosdanese/go-base-api/internal/logging"
	"github.com/jeancarlosdanese/go-base-api/internal/repositories"
	"github.com/jeancarlosdanese/go-base-api/internal/services"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ServicesContainer struct {
	CasbinService      services.CasbinServiceInterface
	TokenService       services.TokenServiceInterface
	TenantService      services.TenantServiceInterface
	UserService        services.UserServiceInterface
	RedisService       services.RedisServiceInterface
	TokenRedisService  services.TokenRedisServiceInterface
	ApiKeyRedisService services.ApiKeyRedisServiceInterface
	RedisClient        *redis.Client
	DB                 *gorm.DB
}

func NewServicesContainer() (*ServicesContainer, error) {
	// Carrega .env ou .env.test com base na variável GO_ENV
	envFile := ".env"
	if os.Getenv("GO_ENV") == "test" {
		envFile = ".env.test"
	}
	if err := godotenv.Load(envFile); err != nil {
		logging.Logger.Warn().
			Err(err).
			Str("env_file", envFile).
			Msg("Arquivo .env não encontrado, usando valores padrão")
	}

	gormDB, err := db.NewDatabaseConnection()
	if err != nil {
		return nil, err
	}

	// Inicializa o Redis
	db.InitializeRedis()
	redisClient := db.GetRedisClient()
	redisService := services.NewRedisService()
	tokenRedisService := services.NewTokenRedisService(redisService)

	casbinService, err := services.NewCasbinService(gormDB)
	if err != nil {
		return nil, err
	}

	// JWT_ACCESS_DURATION=24h
	accessDuration, err := parseDurationWithDays(os.Getenv("JWT_ACCESS_DURATION"))
	if err != nil {
		return nil, err
	}
	// JWT_REFRESH_DURATION=90d
	refreshDuration, err := parseDurationWithDays(os.Getenv("JWT_REFRESH_DURATION"))
	if err != nil {
		return nil, err
	}

	if accessDuration == 0 || refreshDuration == 0 {
		return nil, errors.New("JWT_ACCESS_DURATION or JWT_REFRESH_DURATION is invalid")
	}

	tokenService := services.NewTokenService(os.Getenv("JWT_SECRET_KEY"), accessDuration, refreshDuration)

	tenantsRepo := repositories.NewTenantRepository(gormDB)
	tenantService := services.NewTenantService(tenantsRepo)

	apiKeyRedisService := services.NewApiKeyRedisService(tenantService, redisService, time.Hour*24)

	usersRepo := repositories.NewUserRepository(gormDB)
	userService := services.NewUserService(usersRepo)

	return &ServicesContainer{
		CasbinService:      casbinService,
		TokenService:       tokenService,
		TenantService:      tenantService,
		UserService:        userService,
		RedisService:       redisService,
		TokenRedisService:  tokenRedisService,
		ApiKeyRedisService: apiKeyRedisService,
		RedisClient:        redisClient,
		DB:                 gormDB,
	}, nil
}

// parseDurationWithDays converte uma string de duração, suportando "d" (dias) que não é nativo do Go.
// Exemplos: "90d" -> 2160h, "24h" -> 24h, "1h30m" -> 1h30m
func parseDurationWithDays(s string) (time.Duration, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, errors.New("duration string is empty")
	}

	// Verifica se termina com "d" (dias)
	if strings.HasSuffix(s, "d") {
		// Remove o "d" final
		daysStr := strings.TrimSuffix(s, "d")
		days, err := strconv.Atoi(daysStr)
		if err != nil {
			return 0, errors.New("invalid days format: " + err.Error())
		}
		// Converte dias para horas (1 dia = 24 horas)
		return time.Duration(days) * 24 * time.Hour, nil
	}

	// Se não é dias, usa o parser padrão do Go
	return time.ParseDuration(s)
}
