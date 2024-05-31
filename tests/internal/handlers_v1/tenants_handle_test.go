// tests/internal/handlers_v1/tenants_handle_test.go

package handlers_v1_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jeancarlosdanese/go-base-api/internal/domain/models"
	"github.com/jeancarlosdanese/go-base-api/internal/handlers_v1"
	"github.com/jeancarlosdanese/go-base-api/internal/services"
	"github.com/jeancarlosdanese/go-base-api/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTenantsHandler_GetAll(t *testing.T) {
	mockRepo := new(mocks.MockTenantRepository)
	service := services.NewTenantService(mockRepo)
	handler := handlers_v1.NewTenantsHandler(service)

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
	mockRepo.On("GetAll", mock.Anything).Return(tenants, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	handler.GetAll(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var response []models.Tenant
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 2)
	mockRepo.AssertExpectations(t)
}

func TestTenantsHandler_Create(t *testing.T) {
	mockRepo := new(mocks.MockTenantRepository)
	service := services.NewTenantService(mockRepo)
	handler := handlers_v1.NewTenantsHandler(service)

	tenant := models.Tenant{
		BaseModel: models.BaseModel{ID: uuid.New()},
		Name:      "New Tenant",
	}
	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.Tenant")).Return(&tenant, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	tenantData, _ := json.Marshal(tenant)
	c.Request = httptest.NewRequest("POST", "/tenants", bytes.NewBuffer(tenantData))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Create(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	var response models.Tenant
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, tenant.Name, response.Name)
	mockRepo.AssertExpectations(t)
}

func TestTenantsHandler_GetByID(t *testing.T) {
	mockRepo := new(mocks.MockTenantRepository)
	service := services.NewTenantService(mockRepo)
	handler := handlers_v1.NewTenantsHandler(service)

	tenantID := uuid.New()
	tenant := models.Tenant{
		BaseModel: models.BaseModel{ID: tenantID},
		Name:      "Tenant 1",
	}
	mockRepo.On("GetByID", mock.Anything, tenantID).Return(&tenant, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: tenantID.String()}}

	handler.GetById(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var response models.Tenant
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, tenantID, response.ID)
	mockRepo.AssertExpectations(t)
}

func TestTenantsHandler_Update(t *testing.T) {
	mockRepo := new(mocks.MockTenantRepository)
	service := services.NewTenantService(mockRepo)
	handler := handlers_v1.NewTenantsHandler(service)

	tenantID := uuid.New()
	tenant := models.Tenant{
		BaseModel: models.BaseModel{ID: tenantID},
		Name:      "Updated Tenant",
	}

	// Corrige a configuração do mock para incluir o UUID
	mockRepo.On("Update", mock.Anything, tenantID, &tenant).Return(&tenant, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: tenantID.String()}}

	tenantData, _ := json.Marshal(tenant)
	c.Request = httptest.NewRequest("PUT", "/tenants/"+tenantID.String(), bytes.NewBuffer(tenantData))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Update(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var response models.Tenant
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, tenant.Name, response.Name)
	mockRepo.AssertExpectations(t)
}

func TestTenantsHandler_Delete(t *testing.T) {
	mockRepo := new(mocks.MockTenantRepository)
	service := services.NewTenantService(mockRepo)
	handler := handlers_v1.NewTenantsHandler(service)

	tenantID := uuid.New()
	mockRepo.On("Delete", mock.Anything, tenantID).Return(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: tenantID.String()}}

	c.Request = httptest.NewRequest("DELETE", "/tenants/"+tenantID.String(), nil)

	handler.Delete(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockRepo.AssertExpectations(t)
}
