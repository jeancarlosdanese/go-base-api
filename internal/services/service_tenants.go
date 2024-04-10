// internal/services/service_tenants.go

package services

import (
	"context"

	"github.com/google/uuid"
	"hyberica.io/go/go-api/internal/domain/models"
	"hyberica.io/go/go-api/internal/repositories"
)

type TenantService interface {
	CreateTenant(ctx context.Context, tenant *models.Tenant) error
	GetAllTenants(ctx context.Context) ([]models.Tenant, error)
	GetTenantByID(ctx context.Context, id uuid.UUID) (*models.Tenant, error)
	UpdateTenant(ctx context.Context, tenant *models.Tenant) error
	DeleteTenant(ctx context.Context, id uuid.UUID) error
	UpdateTenantPartial(ctx context.Context, id uuid.UUID, updateData map[string]interface{}) error // Novo método para atualizações parciais
}

// TenantServiceImpl é a implementação concreta de TenantService.
type TenantServiceImpl struct {
	repo repositories.TenantRepository
}

// NewTenantService cria e retorna uma nova instância de TenantService.
func NewTenantService(repo repositories.TenantRepository) TenantService {
	return &TenantServiceImpl{repo: repo}
}

// CreateTenant cria um novo tenant.
func (s *TenantServiceImpl) CreateTenant(ctx context.Context, tenant *models.Tenant) error {
	return s.repo.CreateTenant(ctx, tenant)
}

// GetAllTenants retorna todos os tenants.
func (s *TenantServiceImpl) GetAllTenants(ctx context.Context) ([]models.Tenant, error) {
	return s.repo.GetAllTenants(ctx)
}

func (s *TenantServiceImpl) GetTenantByID(ctx context.Context, id uuid.UUID) (*models.Tenant, error) {
	return s.repo.GetTenantByID(ctx, id)
}

func (s *TenantServiceImpl) UpdateTenant(ctx context.Context, tenant *models.Tenant) error {
	return s.repo.UpdateTenant(ctx, tenant)
}

// UpdateTenantPatch atualiza parcialmente um tenant existente.
func (s *TenantServiceImpl) UpdateTenantPartial(ctx context.Context, id uuid.UUID, updateData map[string]interface{}) error {
	return s.repo.UpdateTenantPartial(ctx, id, updateData)
}

func (s *TenantServiceImpl) DeleteTenant(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteTenant(ctx, id)
}
