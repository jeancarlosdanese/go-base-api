// internal/repositories/repository_tenants.go

package repositories

import (
	"context"

	"github.com/google/uuid"
	"hyberica.io/go/go-api/internal/domain/models"

	"gorm.io/gorm"
)

// TenantRepository define a interface para interação com os tenants no banco de dados.
type TenantRepository interface {
	CreateTenant(ctx context.Context, tenant *models.Tenant) error
	GetTenantByID(ctx context.Context, id uuid.UUID) (*models.Tenant, error)
	GetAllTenants(ctx context.Context) ([]models.Tenant, error) // Adicionado
	UpdateTenant(ctx context.Context, tenant *models.Tenant) error
	UpdateTenantPartial(ctx context.Context, id uuid.UUID, updateData map[string]interface{}) error
	DeleteTenant(ctx context.Context, id uuid.UUID) error
}

// GormTenantRepository é a implementação de TenantRepository usando GORM.
type GormTenantRepository struct {
	DB *gorm.DB
}

// NewGormTenantRepository cria uma nova instância de GormTenantRepository.
func NewGormTenantRepository(db *gorm.DB) *GormTenantRepository {
	return &GormTenantRepository{DB: db}
}

// Aqui vão as implementações dos métodos definidos na interface TenantRepository.
// Exemplo: CreateTenant

func (r *GormTenantRepository) CreateTenant(ctx context.Context, tenant *models.Tenant) error {
	result := r.DB.WithContext(ctx).Create(tenant)
	return result.Error
}

func (r *GormTenantRepository) GetAllTenants(ctx context.Context) ([]models.Tenant, error) {
	var tenants []models.Tenant
	result := r.DB.WithContext(ctx).Find(&tenants)
	if result.Error != nil {
		return nil, result.Error
	}
	return tenants, nil
}

func (r *GormTenantRepository) GetTenantByID(ctx context.Context, id uuid.UUID) (*models.Tenant, error) {
	var tenant models.Tenant
	result := r.DB.WithContext(ctx).Where("id = ?", id).First(&tenant)
	if result.Error != nil {
		return nil, result.Error
	}
	return &tenant, nil
}

func (r *GormTenantRepository) UpdateTenant(ctx context.Context, tenant *models.Tenant) error {
	result := r.DB.WithContext(ctx).Save(tenant)
	return result.Error
}

func (r *GormTenantRepository) UpdateTenantPartial(ctx context.Context, id uuid.UUID, updateData map[string]interface{}) error {
	result := r.DB.WithContext(ctx).Model(&models.Tenant{}).Where("id = ?", id).Updates(updateData)
	return result.Error
}

func (r *GormTenantRepository) DeleteTenant(ctx context.Context, id uuid.UUID) error {
	result := r.DB.WithContext(ctx).Delete(&models.Tenant{}, id)
	return result.Error
}
