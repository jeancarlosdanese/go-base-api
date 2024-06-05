package handlers_v1_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jeancarlosdanese/go-base-api/internal/domain/models"
	"github.com/jeancarlosdanese/go-base-api/internal/handlers_v1"
	"github.com/jeancarlosdanese/go-base-api/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthHandler_Login(t *testing.T) {
	mockUserService := mocks.NewUserService(t)
	mockTokenService := mocks.NewTokenService(t)
	mockTokenRedisService := mocks.NewTokenRedisService(t)
	handler := handlers_v1.NewAuthHandler(mockUserService, mockTokenService, mockTokenRedisService)

	t.Run("successful login", func(t *testing.T) {
		userID := uuid.New()
		user := &models.User{
			BaseModel: models.BaseModel{ID: userID},
			Name:      "John Doe",
			Username:  "johndoe",
			Email:     "john@example.com",
		}

		mockUserService.On("Authenticate", mock.Anything, "john@example.com", "password123", "localhost").Return(user, nil)
		mockTokenService.On("GetAccessDuration").Return(time.Hour * 24) // Adicionando esta linha
		mockTokenService.On("CreateTokens", user.ID, mock.Anything, mock.Anything).Return("access-token", "refresh-token", nil)
		mockTokenRedisService.On("SaveUserRedis", user, "access-token", "refresh-token", mock.AnythingOfType("time.Duration")).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("Origin", "localhost")
		body := `{"email":"john@example.com","password":"password123"}`
		c.Request = httptest.NewRequest("POST", "/login", bytes.NewBufferString(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.Login(c)

		assert.Equal(t, http.StatusOK, w.Code)
		response := make(map[string]interface{})
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "access-token", response["token"])
		assert.Equal(t, "refresh-token", response["refreshToken"])

		// Verifica se todas as expectativas nos mocks foram atendidas
		mockUserService.AssertExpectations(t)
		mockTokenService.AssertExpectations(t)
		mockTokenRedisService.AssertExpectations(t)
	})

	t.Run("invalid login credentials", func(t *testing.T) {
		mockUserService.On("Authenticate", mock.Anything, "john@example.com", "wrongpassword", "localhost").Return(nil, errors.New("invalid credentials"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("Origin", "localhost")
		body := `{"email":"john@example.com","password":"wrongpassword"}`
		c.Request = httptest.NewRequest("POST", "/login", bytes.NewBufferString(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.Login(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		response := make(map[string]interface{})
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "Credenciais inválidas", response["error"])

		// Verifica se todas as expectativas nos mocks foram atendidas
		mockUserService.AssertExpectations(t)
		mockTokenService.AssertExpectations(t)
		mockTokenRedisService.AssertExpectations(t)
	})

	t.Run("successful token refresh", func(t *testing.T) {
		userID := uuid.New()
		user := &models.User{
			BaseModel: models.BaseModel{ID: userID},
			Name:      "John Doe",
			Email:     "john@example.com",
		}

		mockTokenService.On("RefreshTokens", "valid-refresh-token").Return(user.ID, nil)
		mockUserService.On("GetOnlyByID", mock.Anything, user.ID).Return(user, nil)
		mockTokenService.On("GetAccessDuration").Return(time.Hour * 24) // Assume que o token expira em 24 horas
		mockTokenService.On("CreateTokens", user.ID, mock.Anything, mock.Anything).Return("new-access-token", "new-refresh-token", nil)
		mockTokenRedisService.On("SaveUserRedis", user, "new-access-token", "new-refresh-token", mock.AnythingOfType("time.Duration")).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		body := `{"refreshToken":"valid-refresh-token"}`
		c.Request = httptest.NewRequest("POST", "/refresh", bytes.NewBufferString(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.Refresh(c)

		assert.Equal(t, http.StatusOK, w.Code)
		response := make(map[string]interface{})
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "new-access-token", response["token"])
		assert.Equal(t, "new-refresh-token", response["refreshToken"])

		// Verifica se todas as expectativas nos mocks foram atendidas
		mockUserService.AssertExpectations(t)
		mockTokenService.AssertExpectations(t)
		mockTokenRedisService.AssertExpectations(t)
	})
}
