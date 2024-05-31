package mocks

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jeancarlosdanese/go-base-api/internal/domain/models"
	"github.com/stretchr/testify/mock"
)

type MockTenantRepository struct {
	mock.Mock
}

func (m *MockTenantRepository) Create(c *gin.Context, entity *models.Tenant) (*models.Tenant, error) {
	args := m.Called(c, entity)
	return args.Get(0).(*models.Tenant), args.Error(1)
}

func (m *MockTenantRepository) Update(c *gin.Context, id uuid.UUID, entity *models.Tenant) (*models.Tenant, error) {
	args := m.Called(c, id, entity)
	return args.Get(0).(*models.Tenant), args.Error(1)
}

func (m *MockTenantRepository) UpdatePartial(c *gin.Context, id uuid.UUID, updateData map[string]interface{}) (*models.Tenant, error) {
	args := m.Called(c, id, updateData)
	return args.Get(0).(*models.Tenant), args.Error(1)
}

func (m *MockTenantRepository) Delete(c *gin.Context, id uuid.UUID) error {
	args := m.Called(c, id)
	return args.Error(0)
}

func (m *MockTenantRepository) GetAll(c *gin.Context) ([]models.Tenant, error) {
	args := m.Called(c)
	return args.Get(0).([]models.Tenant), args.Error(1)
}

func (m *MockTenantRepository) GetByID(c *gin.Context, id uuid.UUID) (*models.Tenant, error) {
	args := m.Called(c, id)
	return args.Get(0).(*models.Tenant), args.Error(1)
}
