// tests/internal/services/tenant_service_test.go

package services_test

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jeancarlosdanese/go-base-api/internal/domain/models"
	"github.com/jeancarlosdanese/go-base-api/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTenantRepository é um repositório mock para testes
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

func (m *MockTenantRepository) FindByApiKey(apiKey, origin string) (*models.Tenant, error) {
	args := m.Called(apiKey, origin)
	return args.Get(0).(*models.Tenant), args.Error(1)
}

func TestTenantService_Create(t *testing.T) {
	repo := new(MockTenantRepository)
	service := services.NewTenantService(repo)

	c := &gin.Context{}

	email := "tenant@example.com"
	tenantCreate := models.Tenant{
		Type:  "FISICA",
		Name:  "Test Tenant",
		Email: &email,
	}

	repo.On("Create", c, mock.AnythingOfType("*models.Tenant")).Return(&tenantCreate, nil)

	tenant, err := service.Create(c, &tenantCreate)

	assert.NoError(t, err)
	assert.NotNil(t, tenant)
	assert.Equal(t, tenantCreate.Type, tenant.Type)
	assert.Equal(t, tenantCreate.Name, tenant.Name)
	assert.Equal(t, tenantCreate.Email, tenant.Email)

	repo.AssertCalled(t, "Create", c, mock.AnythingOfType("*models.Tenant"))
}

func TestTenantService_Update(t *testing.T) {
	repo := new(MockTenantRepository)
	service := services.NewTenantService(repo)

	c := &gin.Context{}
	tenantID := uuid.New()
	email := "tenant@example.com"
	tenantUpdate := models.Tenant{
		BaseModel: models.BaseModel{ID: tenantID},
		Type:      "FISICA",
		Name:      "Updated Tenant",
		Email:     &email,
	}

	// Configuração correta do mock para incluir o ID do tenant como argumento
	repo.On("Update", c, tenantID, &tenantUpdate).Return(&tenantUpdate, nil)

	tenant, err := service.Update(c, tenantID, &tenantUpdate)

	assert.NoError(t, err)
	assert.NotNil(t, tenant)
	assert.Equal(t, tenantUpdate.Type, tenant.Type)
	assert.Equal(t, tenantUpdate.Name, tenant.Name)
	assert.Equal(t, tenantUpdate.Email, tenant.Email)

	repo.AssertExpectations(t)
}

func TestTenantService_UpdatePartial(t *testing.T) {
	repo := new(MockTenantRepository)
	service := services.NewTenantService(repo)

	c := &gin.Context{}

	tenantID := uuid.New()
	updateData := map[string]interface{}{
		"name": "Partially Updated Tenant",
	}

	tenant := models.Tenant{
		BaseModel: models.BaseModel{ID: tenantID},
		Name:      "Partially Updated Tenant",
	}
	repo.On("UpdatePartial", c, tenantID, updateData).Return(&tenant, nil)

	updatedTenant, err := service.UpdatePartial(c, tenantID, updateData)

	assert.NoError(t, err)
	assert.NotNil(t, updatedTenant)
	assert.Equal(t, updateData["name"], updatedTenant.Name)

	repo.AssertCalled(t, "UpdatePartial", c, tenantID, updateData)
}

func TestTenantService_Delete(t *testing.T) {
	repo := new(MockTenantRepository)
	service := services.NewTenantService(repo)

	c := &gin.Context{}

	tenantID := uuid.New()
	repo.On("Delete", c, tenantID).Return(nil)

	err := service.Delete(c, tenantID)

	assert.NoError(t, err)

	repo.AssertCalled(t, "Delete", c, tenantID)
}

func TestTenantService_GetAll(t *testing.T) {
	repo := new(MockTenantRepository)
	service := services.NewTenantService(repo)

	c := &gin.Context{}

	tenants := []models.Tenant{
		{
			BaseModel: models.BaseModel{ID: uuid.New()},
			Name:      "Tenant 1",
		},
		{
			BaseModel: models.BaseModel{ID: uuid.New()},
			Name:      "Tenant 2",
		},
	}
	repo.On("GetAll", c).Return(tenants, nil)

	allTenants, err := service.GetAll(c)

	assert.NoError(t, err)
	assert.Len(t, allTenants, 2)
	assert.Equal(t, tenants, allTenants)

	repo.AssertCalled(t, "GetAll", c)
}

func TestTenantService_GetByID(t *testing.T) {
	repo := new(MockTenantRepository)
	service := services.NewTenantService(repo)

	c := &gin.Context{}

	tenantID := uuid.New()
	tenant := models.Tenant{
		BaseModel: models.BaseModel{ID: tenantID},
		Name:      "Tenant 1",
	}
	repo.On("GetByID", c, tenantID).Return(&tenant, nil)

	retrievedTenant, err := service.GetByID(c, tenantID)

	assert.NoError(t, err)
	assert.NotNil(t, retrievedTenant)
	assert.Equal(t, tenantID, retrievedTenant.ID)
	assert.Equal(t, tenant.Name, retrievedTenant.Name)

	repo.AssertCalled(t, "GetByID", c, tenantID)
}
