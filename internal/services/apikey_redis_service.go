// @file: internal/services/apikey_redis_service.go

package services

import (
	"encoding/json"
	"time"

	"github.com/jeancarlosdanese/go-base-api/internal/domain/models"
	"github.com/jeancarlosdanese/go-base-api/internal/logging"
	"github.com/redis/go-redis/v9"
)

type ApiKeyRedisServiceInterface interface {
	SaveApiKeyDataRedis(tenant *models.Tenant, apiKey string, accessDuration time.Duration) error
	GetTenantRedisFromApiKey(apiKey, origin string) (*models.TenantRedis, error)
}

type ApiKeyRedisService struct {
	TenantService  TenantServiceInterface
	RedisService   RedisServiceInterface
	AccessDuration time.Duration
}

func NewApiKeyRedisService(tenantsService TenantServiceInterface, redisService RedisServiceInterface, accessDuration time.Duration) *ApiKeyRedisService {
	return &ApiKeyRedisService{
		TenantService:  tenantsService,
		RedisService:   redisService,
		AccessDuration: accessDuration,
	}
}

func (s *ApiKeyRedisService) SaveApiKeyDataRedis(tenant *models.Tenant, apiKey string, accessDuration time.Duration) error {
	apiKeyDataRedis := prepareApiKeyDataRedis(tenant)
	apiKeyData, err := json.Marshal(apiKeyDataRedis)
	if err != nil {
		return err
	}
	if err := s.RedisService.Set("apiKey:"+apiKey, apiKeyData, accessDuration); err != nil {
		logging.Logger.Error().
			Err(err).
			Str("operation", "save_apikey_redis").
			Msg("Erro ao salvar dados da API Key no Redis")
		return err
	}
	return nil
}

func (s *ApiKeyRedisService) GetTenantRedisFromApiKey(apiKey, origin string) (*models.TenantRedis, error) {
	result, err := s.RedisService.Get("apiKey:" + apiKey)
	if err != nil && err != redis.Nil {
		logging.Logger.Error().
			Err(err).
			Str("operation", "get_tenant_from_apikey").
			Msg("Erro ao recuperar dados do Redis")
		return nil, err
	}

	if result == "" {
		tenant, err := s.TenantService.ApiKeyAuthenticate(apiKey, origin)
		if err != nil {
			logging.Logger.Error().
				Err(err).
				Str("origin", origin).
				Str("operation", "authenticate_apikey").
				Msg("Erro ao autenticar API Key")
			return nil, err
		}

		if err := s.SaveApiKeyDataRedis(tenant, apiKey, s.AccessDuration); err != nil {
			logging.Logger.Error().
				Err(err).
				Str("operation", "save_apikey_data_redis").
				Msg("Erro ao salvar dados da API Key no Redis")
			return nil, err
		}

		// Retrieve again to confirm saving was successful
		result, err = s.RedisService.Get("apiKey:" + apiKey)
		if err != nil {
			logging.Logger.Error().
				Err(err).
				Str("operation", "verify_stored_apikey").
				Msg("Erro ao verificar chave armazenada no Redis")
			return nil, err
		}
	}

	var apiKeyDataRedis models.TenantRedis
	if err := json.Unmarshal([]byte(result), &apiKeyDataRedis); err != nil {
		logging.Logger.Error().
			Err(err).
			Str("operation", "unmarshal_apikey_data").
			Msg("Erro ao fazer unmarshal dos dados da API Key")
		return nil, err
	}

	return &apiKeyDataRedis, nil
}

// prepareApiKeyDataRedis prepares user data to be stored in Redis.
func prepareApiKeyDataRedis(tenant *models.Tenant) models.TenantRedis {
	var cpfCnpj, email string

	if tenant.CpfCnpj != nil {
		cpfCnpj = *tenant.CpfCnpj
	}
	if tenant.Email != nil {
		email = *tenant.Email
	}

	return models.TenantRedis{
		ID:      tenant.ID.String(),
		Name:    tenant.Name,
		CpfCnpj: cpfCnpj,
		Email:   email,
	}
}
