// internal/app/services_container.go

package app

import (
	"gorm.io/gorm"
	"hyberica.io/go/go-api/internal/db"
	"hyberica.io/go/go-api/internal/repositories"
	"hyberica.io/go/go-api/internal/services"
)

type ServicesContainer struct {
	TenantService services.TenantService
	DB            *gorm.DB
}

func NewServicesContainer() (*ServicesContainer, error) {
	db, err := db.NewDatabaseConnection()
	if err != nil {
		return nil, err
	}

	tenantRepo := repositories.NewGormTenantRepository(db) // Assume que agora aceita *gorm.DB
	tenantService := services.NewTenantService(tenantRepo)

	return &ServicesContainer{
		TenantService: tenantService,
		DB:            db,
	}, nil
}
