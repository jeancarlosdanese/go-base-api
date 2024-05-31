// tests/integration/integration_auth_test.go

package integration_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jeancarlosdanese/go-base-api/tests/testutils"
	"github.com/stretchr/testify/assert"
)

var testAuthRouter *gin.Engine

func TestAuthHandler_Login(t *testing.T) {
	testAuthRouter = testutils.SetupRouter()
	t.Run("successful login", func(t *testing.T) {
		loginData := map[string]string{
			"email":    "master@domain.local",
			"password": "master123",
		}
		loginBody, _ := json.Marshal(loginData)
		loginReq, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(loginBody))
		loginReq.Header.Set("Content-Type", "application/json")
		loginReq.Header.Set("Origin", "http://localhost")
		loginResp := httptest.NewRecorder()
		testAuthRouter.ServeHTTP(loginResp, loginReq)

		if loginResp.Code != http.StatusOK {
			log.Fatalf("Expected status code 200 for login request, got %v", loginResp.Code)
		}

		var loginResponse map[string]interface{}
		err := json.Unmarshal(loginResp.Body.Bytes(), &loginResponse)
		if err != nil {
			log.Fatalf("Failed to parse login response: %v", err)
		}
		_, ok := loginResponse["token"].(string)
		if !ok {
			log.Fatalf("Login response should contain token")
		}
	})

	t.Run("invalid login", func(t *testing.T) {
		loginPayload := map[string]string{
			"email":    "invalid@example.com",
			"password": "wrongPassword",
		}
		body, _ := json.Marshal(loginPayload)
		req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Origin", "http://localhost")
		w := httptest.NewRecorder()

		testAuthRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		log.Printf("RESPONSE: %v", response["error"])
		assert.Equal(t, "Credenciais inválidas", response["error"])
	})
}

func TestAuthHandler_RefreshToken(t *testing.T) {
	_, refreshToken := testutils.LoginAndGetToken(t, testAuthRouter, "master@domain.local", "master123")
	t.Run("successful token refresh", func(t *testing.T) {
		loginData := map[string]string{
			"refreshToken": refreshToken,
		}
		body, err := json.Marshal(loginData)
		if err != nil {
			t.Fatalf("Error marshaling login data: %v", err)
		}
		req := httptest.NewRequest("POST", "/api/v1/auth/refresh", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		testAuthRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.NotEmpty(t, response["token"], "Expected new access token")
		assert.NotEmpty(t, response["refreshToken"], "Expected new refresh token")
	})

	t.Run("invalid refresh token", func(t *testing.T) {
		loginData := map[string]string{
			"refreshToken": "invalidToken",
		}
		body, err := json.Marshal(loginData)
		if err != nil {
			t.Fatalf("Error marshaling login data: %v", err)
		}
		req := httptest.NewRequest("POST", "/api/v1/auth/refresh", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		testAuthRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "Token inválido ou expirado", response["error"])
	})
}
