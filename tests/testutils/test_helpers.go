// tests/testutils/test_helpers.go

package testutils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jeancarlosdanese/go-base-api/internal/app"
	"github.com/jeancarlosdanese/go-base-api/internal/domain/models"
	"github.com/jeancarlosdanese/go-base-api/internal/routes"
)

// Setup the Gin router and any necessary services
func SetupRouter() *gin.Engine {
	// Definir a variável de ambiente GO_ENV para "test"
	log.Println("INIT OF TEST")
	os.Setenv("GO_ENV", "test")

	// Verificar se as variáveis de ambiente foram carregadas corretamente
	fmt.Printf("DB_NAME: %s\n", os.Getenv("DB_NAME"))

	// Set up service container and routes
	sc, err := app.NewServicesContainer()
	if err != nil {
		panic("Failed to set up services container: " + err.Error())
	}

	router := gin.Default()
	routes.SetupRouter(router, sc)

	return router
}

// LoginAndGetToken realiza o login usando credenciais válidas e retorna o usuário e o token para ser usado em outros testes de integração.
func LoginAndGetToken(t *testing.T, router *gin.Engine, email, password string) (string, string) {
	t.Helper()
	loginData := map[string]string{
		"email":    email,
		"password": password,
	}
	loginBody, err := json.Marshal(loginData)
	if err != nil {
		t.Fatalf("Error marshaling login data: %v", err)
	}

	loginReq, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(loginBody))
	loginReq.Header.Set("Content-Type", "application/json")
	loginReq.Header.Set("Origin", "http://localhost")
	loginResp := httptest.NewRecorder()
	router.ServeHTTP(loginResp, loginReq)

	if loginResp.Code != http.StatusOK {
		t.Fatalf("Expected status code 200 for login request, got %v", loginResp.Code)
	}

	var token *models.Token

	if err := json.Unmarshal(loginResp.Body.Bytes(), &token); err != nil {
		t.Fatalf("Failed to parse login response: %v", err)
	}

	return token.Token, *token.RefreshToken
}
