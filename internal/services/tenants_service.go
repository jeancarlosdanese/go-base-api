// internal/services/tenants_service.go

package services

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jeancarlosdanese/go-base-api/internal/domain/models"
	"github.com/jeancarlosdanese/go-base-api/internal/logging"
	"github.com/jeancarlosdanese/go-base-api/internal/repositories"
	"github.com/jeancarlosdanese/go-base-api/internal/utils"
)

// TenantServiceInterface define as operações adicionais do TenantService além das operações CRUD básicas.
type TenantServiceInterface interface {
	BaseServiceInterface[models.Tenant]
	CreateTenantWithApiKey(c *gin.Context, entity *models.Tenant) (*models.Tenant, error)
	ApiKeyAuthenticate(apiKey, origin string) (*models.Tenant, error)
}
type TenantService struct {
	*BaseService[models.Tenant, repositories.TenantRepository]
}

func NewTenantService(repo repositories.TenantRepository) *TenantService {
	baseService := NewBaseService[models.Tenant, repositories.TenantRepository](repo)
	return &TenantService{BaseService: baseService}
}

// CreateTenantWithApiKey cria um tenant gerando e adicionando uma ApiKey para o tenant
func (s *TenantService) CreateTenantWithApiKey(c *gin.Context, tenant *models.Tenant) (*models.Tenant, error) {
	// Gera uma ApiKey para a o Tenant
	apikey, err := utils.GenerateApiKey(64)
	if err != nil {
		return nil, err
	}

	fmt.Printf("API_KEY: %v", apikey)

	tenant.ApiKey = &apikey

	tenantCreated, err := s.Repo.Create(c, tenant)
	if err != nil {
		return nil, err
	}

	return tenantCreated, nil
}

// ApiKeyAuthenticate verifica a apiKey do Tenant.
func (s *TenantService) ApiKeyAuthenticate(apiKey, origin string) (*models.Tenant, error) {
	user, err := s.Repo.FindByApiKey(apiKey, origin)
	if err != nil {
		logging.InfoLogger.Printf("Erro ao buscar Tenant por apiKey: %v", err)
		return nil, err
	}

	return user, nil
}
