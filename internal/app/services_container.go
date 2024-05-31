// internal/app/services_container.go

package app

import (
	"log"
	"os"
	"time"

	// Import correto
	"github.com/jeancarlosdanese/go-base-api/internal/db"
	"github.com/jeancarlosdanese/go-base-api/internal/repositories"
	"github.com/jeancarlosdanese/go-base-api/internal/services"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type ServicesContainer struct {
	CasbinService     services.CasbinServiceInterface
	TokenService      services.TokenServiceInterface
	TenantService     services.TenantServiceInterface
	UserService       services.UserServiceInterface
	RedisService      services.RedisServiceInterface
	TokenRedisService services.TokenRedisServiceInterface
	DB                *gorm.DB
}

func NewServicesContainer() (*ServicesContainer, error) {
	// Carrega .env ou .env.test com base na variável GO_ENV
	envFile := ".env"
	if os.Getenv("GO_ENV") == "test" {
		envFile = ".env.test"
	}
	if err := godotenv.Load(envFile); err != nil {
		log.Printf("Warning: %s file not found", envFile)
	}

	gormDB, err := db.NewDatabaseConnection()
	if err != nil {
		return nil, err
	}

	// Inicializa o Redis
	db.InitializeRedis()
	redisService := services.NewRedisService()
	tokenRedisService := services.NewTokenRedisService(redisService)

	casbinService, err := services.NewCasbinService(gormDB)
	if err != nil {
		return nil, err
	}

	tokenService := services.NewTokenService(os.Getenv("JWT_SECRET_KEY"), time.Hour*24, time.Hour*24*90)

	tenantsRepo := repositories.NewTenantRepository(gormDB)
	tenantsService := services.NewTenantService(tenantsRepo)

	usersRepo := repositories.NewUserRepository(gormDB)
	usersService := services.NewUserService(usersRepo)

	return &ServicesContainer{
		CasbinService:     casbinService,
		TokenService:      tokenService,
		TenantService:     tenantsService,
		UserService:       usersService,
		RedisService:      redisService,
		TokenRedisService: tokenRedisService,
		DB:                gormDB,
	}, nil
}
