// internal/services/tenants_service.go

package services

import (
	"github.com/jeancarlosdanese/go-base-api/internal/domain/models"
	"github.com/jeancarlosdanese/go-base-api/internal/repositories"
)

type TenantService struct {
	*BaseService[models.Tenant, repositories.TenantRepository]
}

func NewTenantService(repo repositories.TenantRepository) *TenantService {
	baseService := NewBaseService(repo)
	return &TenantService{BaseService: baseService}
}
