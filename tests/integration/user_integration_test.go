// tests/integration/user_integration_test.go

package integration_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jeancarlosdanese/go-base-api/internal/domain/models"
	"github.com/jeancarlosdanese/go-base-api/tests/testutils"
	"github.com/stretchr/testify/assert"
)

var testUserRouter *gin.Engine

func setupUserRouter() {
	if testUserRouter == nil {
		testUserRouter = testutils.SetupRouter() // Configura o router e os serviços uma única vez
	}
}

func getTokenForUser(t *testing.T) string {
	t.Helper()
	email := "master@domain.local"
	password := "master123"
	token, _ := testutils.LoginAndGetToken(t, testUserRouter, email, password)
	if token == "" {
		t.Fatalf("Failed to obtain token by %v and %v", email, password)
	}
	return token
}

func createMockUser() *models.User {
	// Criação de um novo User
	return &models.User{
		Name:     "New User",
		Email:    "newuser@example.com",
		Password: "securepassword",
	}
}
func TestUserHandlers(t *testing.T) {
	setupUserRouter()
	token := getTokenForUser(t) // Supõe que getToken já realiza o login e retorna um token válido
	userData := createMockUser()

	t.Run("Create User", func(t *testing.T) {

		body, _ := json.Marshal(userData)
		req, _ := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		resp := httptest.NewRecorder()

		testUserRouter.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusCreated, resp.Code)
		var createdUser models.User
		err := json.Unmarshal(resp.Body.Bytes(), &createdUser)
		assert.NoError(t, err)
		assert.Equal(t, "New User", createdUser.Name)
		userData.ID = createdUser.ID
		log.Printf("USER_ID CREATED: %v", userData.ID)
	})

	t.Run("Get All Users", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/users", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp := httptest.NewRecorder()

		testUserRouter.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		var users []models.User
		err := json.Unmarshal(resp.Body.Bytes(), &users)
		assert.NoError(t, err)
		assert.NotEmpty(t, users)
	})

	t.Run("Update User", func(t *testing.T) {
		// Supõe que já existe um usuário que pode ser atualizado
		userUpdate := map[string]interface{}{
			"name": "Updated Name",
		}
		log.Printf("USER_ID UPDATE: %v", userData.ID)
		body, _ := json.Marshal(userUpdate)
		req, _ := http.NewRequest("PUT", "/api/v1/users/"+userData.ID.String(), bytes.NewBuffer(body)) // Substituir {id} pelo ID real
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		resp := httptest.NewRecorder()

		testUserRouter.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		var updatedUser models.User
		err := json.Unmarshal(resp.Body.Bytes(), &updatedUser)
		assert.NoError(t, err)
		assert.Equal(t, "Updated Name", updatedUser.Name)
	})

	t.Run("Delete User", func(t *testing.T) {
		// Supõe que já existe um usuário que pode ser deletado
		req, _ := http.NewRequest("DELETE", "/api/v1/users/"+userData.ID.String(), nil) // Substituir {id} pelo ID real
		req.Header.Set("Authorization", "Bearer "+token)
		resp := httptest.NewRecorder()

		testUserRouter.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
	})
}
