// internal/services/tenants_service.go

package services

import (
	"github.com/jeancarlosdanese/go-base-api/internal/domain/models"
	"github.com/jeancarlosdanese/go-base-api/internal/repositories"
)

// TenantServiceInterface define as operações adicionais do TenantService além das operações CRUD básicas.
type TenantServiceInterface interface {
	BaseServiceInterface[models.Tenant]
}
type TenantService struct {
	*BaseService[models.Tenant, repositories.TenantRepository]
}

func NewTenantService(repo repositories.TenantRepository) *TenantService {
	baseService := NewBaseService(repo)
	return &TenantService{BaseService: baseService}
}
