// tests/integration/integration_tenant_test.go

package integration_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jeancarlosdanese/go-base-api/internal/domain/models"
	"github.com/stretchr/testify/assert"

	"github.com/jeancarlosdanese/go-base-api/tests/testutils"
)

var testTenantRouter *gin.Engine

func setupTenantRouter() {
	if testTenantRouter == nil {
		testTenantRouter = testutils.SetupRouter() // Configura o router e os serviços uma única vez
	}
}

func getTokenForTenant(t *testing.T) string {
	t.Helper()
	email := "master@domain.local"
	password := "master123"
	token, _ := testutils.LoginAndGetToken(t, testTenantRouter, email, password)
	if token == "" {
		t.Fatalf("Failed to obtain token by %v and %v", email, password)
	}
	return token
}

func createMockTenant() *models.Tenant {
	// Criação de um novo tenant
	email := "tenant@example.com"
	return &models.Tenant{
		Name:  "New Tenant",
		Type:  "FISICA",
		Email: &email,
	}
}

func TestTenantHandleIntegration(t *testing.T) {
	setupTenantRouter()
	testTenantToken := getTokenForTenant(t)
	newTenant := createMockTenant()

	t.Run("Create Tenant", func(t *testing.T) {
		// Setup específico para este teste
		tenantPayload, _ := json.Marshal(&newTenant)
		req, _ := http.NewRequest("POST", "/api/v1/tenants", bytes.NewBuffer(tenantPayload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+testTenantToken)
		req.Header.Set("Origin", "http://localhost")
		resp := httptest.NewRecorder()
		testTenantRouter.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusCreated, resp.Code, "Expected status code 201 for tenant creation request")
		var createdTenant models.Tenant
		err := json.Unmarshal(resp.Body.Bytes(), &createdTenant)
		assert.NoError(t, err)
		assert.Equal(t, "New Tenant", createdTenant.Name, "Expected tenant name to be 'New Tenant'")

		newTenant.ID = createdTenant.ID
	})

	t.Run("Get All Tenants", func(t *testing.T) {
		// Execute get all tenants request
		tenantsReq, _ := http.NewRequest("GET", "/api/v1/tenants", nil)
		tenantsReq.Header.Set("Authorization", "Bearer "+testTenantToken)

		tenantsResp := httptest.NewRecorder()
		testTenantRouter.ServeHTTP(tenantsResp, tenantsReq)

		assert.Equal(t, http.StatusOK, tenantsResp.Code, "Expected status code 200 for get all tenants request, got %v", tenantsResp.Code)

		var tenantsResponse []map[string]interface{}
		err := json.Unmarshal(tenantsResp.Body.Bytes(), &tenantsResponse)
		assert.NoError(t, err, "Failed to parse get all tenants response")
		assert.NotEmpty(t, tenantsResponse, "Expected non-empty list of tenants")
	})

	t.Run("Get Tenant By ID", func(t *testing.T) {
		// Execute get tenant by ID request
		getTenantReq, _ := http.NewRequest("GET", "/api/v1/tenants/"+newTenant.ID.String(), nil)
		getTenantReq.Header.Set("Authorization", "Bearer "+testTenantToken)

		getTenantResp := httptest.NewRecorder()
		testTenantRouter.ServeHTTP(getTenantResp, getTenantReq)

		assert.Equal(t, http.StatusOK, getTenantResp.Code, "Expected status code 200 for get tenant by ID request, got %v", getTenantResp.Code)

		var getTenantResponse map[string]interface{}
		err := json.Unmarshal(getTenantResp.Body.Bytes(), &getTenantResponse)
		assert.NoError(t, err, "Failed to parse get tenant by ID response")
		assert.Equal(t, "New Tenant", getTenantResponse["name"], "Expected tenant name to be 'New Tenant', got %v", getTenantResponse["name"])
	})

	t.Run("Update Tenant", func(t *testing.T) {
		// Update tenant request payload
		updatePayload := map[string]interface{}{
			"type":   "FISICA",
			"name":   "Updated Tenant Name",
			"status": "ATIVO",
		}
		updatePayloadBytes, _ := json.Marshal(updatePayload)

		// Execute update tenant request
		updateReq, _ := http.NewRequest("PUT", "/api/v1/tenants/"+newTenant.ID.String(), bytes.NewBuffer(updatePayloadBytes))
		updateReq.Header.Set("Content-Type", "application/json")
		updateReq.Header.Set("Authorization", "Bearer "+testTenantToken)

		updateResp := httptest.NewRecorder()
		testTenantRouter.ServeHTTP(updateResp, updateReq)

		assert.Equal(t, http.StatusOK, updateResp.Code, "Expected status code 200 for update tenant request, got %v", updateResp.Code)

		var updateResponse map[string]interface{}
		err := json.Unmarshal(updateResp.Body.Bytes(), &updateResponse)
		assert.NoError(t, err, "Failed to parse update tenant response")
		assert.Equal(t, "Updated Tenant Name", updateResponse["name"], "Expected tenant name to be 'Updated Tenant Name', got %v", updateResponse["name"])
	})

	t.Run("Partial Update tenant", func(t *testing.T) {
		// Update tenant request payload
		updatePayload := map[string]interface{}{
			"name": "Partial updated Tenant Name",
		}
		updatePayloadBytes, _ := json.Marshal(updatePayload)

		// Execute update tenant request
		updateReq, _ := http.NewRequest("PATCH", "/api/v1/tenants/"+newTenant.ID.String(), bytes.NewBuffer(updatePayloadBytes))
		updateReq.Header.Set("Content-Type", "application/json")
		updateReq.Header.Set("Authorization", "Bearer "+testTenantToken)

		updateResp := httptest.NewRecorder()
		testTenantRouter.ServeHTTP(updateResp, updateReq)

		assert.Equal(t, http.StatusOK, updateResp.Code, "Expected status code 200 for update tenant request, got %v", updateResp.Code)

		var updateResponse map[string]interface{}
		err := json.Unmarshal(updateResp.Body.Bytes(), &updateResponse)
		assert.NoError(t, err, "Failed to parse update tenant response")
		assert.Equal(t, "Partial updated Tenant Name", updateResponse["name"], "Expected tenant name to be 'Updated Tenant Name', got %v", updateResponse["name"])
	})

	t.Run("Delete Tenant", func(t *testing.T) {
		// Execute delete tenant request
		deleteReq, _ := http.NewRequest("DELETE", "/api/v1/tenants/"+newTenant.ID.String(), nil)
		deleteReq.Header.Set("Authorization", "Bearer "+testTenantToken)

		deleteResp := httptest.NewRecorder()
		testTenantRouter.ServeHTTP(deleteResp, deleteReq)

		assert.Equal(t, http.StatusOK, deleteResp.Code, "Expected status code 200 for delete tenant request, got %v", deleteResp.Code)
	})
}

// func TestCreateTenantIntegration(t *testing.T) {
// 	tenantPayload, _ := json.Marshal(&newTenant)
// 	req, _ := http.NewRequest("POST", "/api/v1/tenants", bytes.NewBuffer(tenantPayload))
// 	req.Header.Set("Content-Type", "application/json")
// 	req.Header.Set("Authorization", "Bearer "+testTenantToken)
// 	req.Header.Set("Origin", "http://localhost")
// 	resp := httptest.NewRecorder()
// 	testTenantRouter.ServeHTTP(resp, req)

// 	assert.Equal(t, http.StatusCreated, resp.Code, "Expected status code 201 for tenant creation request")
// 	var createdTenant models.Tenant
// 	err := json.Unmarshal(resp.Body.Bytes(), &createdTenant)
// 	assert.NoError(t, err)
// 	assert.Equal(t, "New Tenant", createdTenant.Name, "Expected tenant name to be 'New Tenant'")

// 	newTenant.ID = createdTenant.ID
// }

// func TestGetAllTenantsIntegration(t *testing.T) {
// 	// Execute get all tenants request
// 	tenantsReq, _ := http.NewRequest("GET", "/api/v1/tenants", nil)
// 	tenantsReq.Header.Set("Authorization", "Bearer "+testTenantToken)

// 	tenantsResp := httptest.NewRecorder()
// 	testTenantRouter.ServeHTTP(tenantsResp, tenantsReq)

// 	assert.Equal(t, http.StatusOK, tenantsResp.Code, "Expected status code 200 for get all tenants request, got %v", tenantsResp.Code)

// 	var tenantsResponse []map[string]interface{}
// 	err := json.Unmarshal(tenantsResp.Body.Bytes(), &tenantsResponse)
// 	assert.NoError(t, err, "Failed to parse get all tenants response")
// 	assert.NotEmpty(t, tenantsResponse, "Expected non-empty list of tenants")
// }

// func TestGetTenantByIDIntegration(t *testing.T) {
// 	// Execute get tenant by ID request
// 	getTenantReq, _ := http.NewRequest("GET", "/api/v1/tenants/"+newTenant.ID.String(), nil)
// 	getTenantReq.Header.Set("Authorization", "Bearer "+testTenantToken)

// 	getTenantResp := httptest.NewRecorder()
// 	testTenantRouter.ServeHTTP(getTenantResp, getTenantReq)

// 	assert.Equal(t, http.StatusOK, getTenantResp.Code, "Expected status code 200 for get tenant by ID request, got %v", getTenantResp.Code)

// 	var getTenantResponse map[string]interface{}
// 	err := json.Unmarshal(getTenantResp.Body.Bytes(), &getTenantResponse)
// 	assert.NoError(t, err, "Failed to parse get tenant by ID response")
// 	assert.Equal(t, "New Tenant", getTenantResponse["name"], "Expected tenant name to be 'New Tenant', got %v", getTenantResponse["name"])
// }

// func TestUpdateTenantIntegration(t *testing.T) {
// 	// Update tenant request payload
// 	updatePayload := map[string]interface{}{
// 		"type":   "FISICA",
// 		"name":   "Updated Tenant Name",
// 		"status": "ATIVO",
// 	}
// 	updatePayloadBytes, _ := json.Marshal(updatePayload)

// 	// Execute update tenant request
// 	updateReq, _ := http.NewRequest("PUT", "/api/v1/tenants/"+newTenant.ID.String(), bytes.NewBuffer(updatePayloadBytes))
// 	updateReq.Header.Set("Content-Type", "application/json")
// 	updateReq.Header.Set("Authorization", "Bearer "+testTenantToken)

// 	updateResp := httptest.NewRecorder()
// 	testTenantRouter.ServeHTTP(updateResp, updateReq)

// 	assert.Equal(t, http.StatusOK, updateResp.Code, "Expected status code 200 for update tenant request, got %v", updateResp.Code)

// 	var updateResponse map[string]interface{}
// 	err := json.Unmarshal(updateResp.Body.Bytes(), &updateResponse)
// 	assert.NoError(t, err, "Failed to parse update tenant response")
// 	assert.Equal(t, "Updated Tenant Name", updateResponse["name"], "Expected tenant name to be 'Updated Tenant Name', got %v", updateResponse["name"])
// }

// func TestPartialUpdateTenantIntegration(t *testing.T) {
// 	// Update tenant request payload
// 	updatePayload := map[string]interface{}{
// 		"name": "Partial updated Tenant Name",
// 	}
// 	updatePayloadBytes, _ := json.Marshal(updatePayload)

// 	// Execute update tenant request
// 	updateReq, _ := http.NewRequest("PUT", "/api/v1/tenants/"+newTenant.ID.String(), bytes.NewBuffer(updatePayloadBytes))
// 	updateReq.Header.Set("Content-Type", "application/json")
// 	updateReq.Header.Set("Authorization", "Bearer "+testTenantToken)

// 	updateResp := httptest.NewRecorder()
// 	testTenantRouter.ServeHTTP(updateResp, updateReq)

// 	assert.Equal(t, http.StatusOK, updateResp.Code, "Expected status code 200 for update tenant request, got %v", updateResp.Code)

// 	var updateResponse map[string]interface{}
// 	err := json.Unmarshal(updateResp.Body.Bytes(), &updateResponse)
// 	assert.NoError(t, err, "Failed to parse update tenant response")
// 	assert.Equal(t, "Partial updated Tenant Name", updateResponse["name"], "Expected tenant name to be 'Updated Tenant Name', got %v", updateResponse["name"])
// }

// func TestDeleteTenantIntegration(t *testing.T) {
// 	// Execute delete tenant request
// 	deleteReq, _ := http.NewRequest("DELETE", "/api/v1/tenants/"+newTenant.ID.String(), nil)
// 	deleteReq.Header.Set("Authorization", "Bearer "+testTenantToken)

// 	deleteResp := httptest.NewRecorder()
// 	testTenantRouter.ServeHTTP(deleteResp, deleteReq)

// 	assert.Equal(t, http.StatusOK, deleteResp.Code, "Expected status code 200 for delete tenant request, got %v", deleteResp.Code)
// }
