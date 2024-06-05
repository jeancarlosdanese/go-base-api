// internal/repositories/repository_tenants.go

package repositories

import (
	"errors"
	"fmt"

	"github.com/jeancarlosdanese/go-base-api/internal/domain/models"
	"github.com/jeancarlosdanese/go-base-api/internal/logging"
	"gorm.io/gorm"
)

// TenantRepository é uma interface que estende a interface Repository para operações específicas do Tenant.
type TenantRepository interface {
	GormRepositoryInterface[models.Tenant]
	FindByApiKey(apiKey, origin string) (*models.Tenant, error)
}

// NewTenantRepository cria uma nova instância de um repositório que implementa TenantRepository.
func NewTenantRepository(db *gorm.DB) TenantRepository {
	return NewGormRepository[models.Tenant](db).(TenantRepository) // Retornando diretamente a interface
}

func (r *GormRepository[Entity]) FindByApiKey(apiKey, origin string) (*models.Tenant, error) {
	var tenant *models.Tenant
	formattedOrigin := fmt.Sprintf(`["%s"]`, origin)
	err := r.DB.
		Where("api_key = ? AND allowed_origins @> ?", apiKey, formattedOrigin).
		Take(&tenant).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logging.InfoLogger.Printf("Tenant ou origem não encontrado: %s, %s", apiKey, origin)
			return nil, errors.New("tenant ou origem não encontrado")
		}
		logging.ErrorLogger.Printf("Erro ao buscar Tenant por apiKey e origem: %v", err)
		return nil, err
	}

	return tenant, nil
}
