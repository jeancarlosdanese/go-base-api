// internal/repositories/repository_tenants.go

package repositories

import (
	"github.com/jeancarlosdanese/go-base-api/internal/domain/models"
	"gorm.io/gorm"
)

// TenantRepository é uma interface que estende a interface Repository para operações específicas do Tenant.
type TenantRepository interface {
	GormRepositoryInterface[models.Tenant]
}

// NewTenantRepository cria uma nova instância de um repositório que implementa TenantRepository.
func NewTenantRepository(db *gorm.DB) TenantRepository {
	return NewGormRepository[models.Tenant](db) // Retornando diretamente a interface
}
