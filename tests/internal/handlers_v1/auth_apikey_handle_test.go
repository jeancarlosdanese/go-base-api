// tests/internal/handlers_v1/auth_apikey_handle_test.go

package handlers_v1_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	contextkeys "github.com/jeancarlosdanese/go-base-api/internal/domain/context_keys"
	"github.com/jeancarlosdanese/go-base-api/internal/domain/models"
	"github.com/jeancarlosdanese/go-base-api/internal/handlers_v1"
	"github.com/stretchr/testify/assert"
)

func TestGetTenantByApiKey(t *testing.T) {
	// Setup the Gin router and context
	router := gin.Default()
	authHandler := handlers_v1.NewAuthApiKeyHandler()
	api := router.Group("/api/v1/auth-apikey")
	authHandler.RegisterRoutes(api)

	// Mock tenant data
	tenantID := uuid.New().String()
	tenant := &models.TenantRedis{
		ID:   tenantID,
		Name: "Test Tenant",
	}

	// Create a request to send to the handler
	req, _ := http.NewRequest("GET", "/api/v1/auth-apikey/tenant-by-apikey", nil)
	req.Header.Set("Origin", "test.com")
	req.Header.Set("X-API-Key", "test-api-key")

	// Create a response recorder to record the response
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Set the context variables
	c.Set("Origin", "test.com")
	c.Set(string(contextkeys.TenantDataKey), tenant)

	// Call the handler
	authHandler.GetTenantByApiKey(c)

	// Check the response
	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "test.com", response["origin"])
	assert.Equal(t, tenantID, response["tenant"].(map[string]interface{})["id"])
	assert.Equal(t, "Test Tenant", response["tenant"].(map[string]interface{})["name"])
}
