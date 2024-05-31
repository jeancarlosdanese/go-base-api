// tests/internal/handlers_v1/users_handle_test.go

package handlers_v1_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	contextkeys "github.com/jeancarlosdanese/go-base-api/internal/domain/context_keys"
	"github.com/jeancarlosdanese/go-base-api/internal/domain/models"
	"github.com/jeancarlosdanese/go-base-api/internal/handlers_v1"
	"github.com/jeancarlosdanese/go-base-api/internal/services"
	"github.com/jeancarlosdanese/go-base-api/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUsersHandler_GetAll(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	userService := services.NewUserService(mockRepo)
	handler := handlers_v1.NewUsersHandler(userService)

	users := []models.User{
		{
			BaseModel: models.BaseModel{ID: uuid.New()},
			Username:  "User1",
			Name:      "User One",
		},
		{
			BaseModel: models.BaseModel{ID: uuid.New()},
			Username:  "User2",
			Name:      "User Two",
		},
	}
	mockRepo.On("GetAll", mock.Anything).Return(users, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	handler.GetAll(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var response []models.User
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 2)
	mockRepo.AssertExpectations(t)
}

func TestUsersHandler_Create(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	service := services.NewUserService(mockRepo)
	handler := handlers_v1.NewUsersHandler(service)

	user := models.User{
		BaseModel: models.BaseModel{ID: uuid.New()},
		Username:  "newuser",
		Name:      "New User",
		Email:     "newuser@example.com",
	}
	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.User")).Return(&user, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	tenantID := uuid.New().String()
	c.Set(string(contextkeys.TenantIDKey), tenantID)

	userCreate := models.UserCreate{
		Username: "newuser",
		Name:     "New User",
		Email:    "newuser@example.com",
		Password: "password123",
	}

	userData, _ := json.Marshal(userCreate)
	c.Request = httptest.NewRequest("POST", "/users", bytes.NewBuffer(userData))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Create(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	var response models.User
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, user.Username, response.Username)
	assert.Equal(t, user.Name, response.Name)
	assert.Equal(t, user.Email, response.Email)
	mockRepo.AssertExpectations(t)
}

func TestUsersHandler_GetByID(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	service := services.NewUserService(mockRepo)
	handler := handlers_v1.NewUsersHandler(service)

	userID := uuid.New()
	user := models.User{
		BaseModel: models.BaseModel{ID: userID},
		Username:  "User1",
		Name:      "User One",
	}
	mockRepo.On("GetByID", mock.Anything, userID).Return(&user, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: userID.String()}}

	handler.GetById(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var response models.User
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, userID, response.ID)
	assert.Equal(t, user.Username, response.Username)
	assert.Equal(t, user.Name, response.Name)
	mockRepo.AssertExpectations(t)
}

func TestUsersHandler_Update(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	service := services.NewUserService(mockRepo)
	handler := handlers_v1.NewUsersHandler(service)

	userID := uuid.New()
	tenantID := uuid.New() // Suponha que esta é a identificação do tenant necessária
	user := models.User{
		BaseModel: models.BaseModel{ID: userID},
		TenantID:  tenantID, // Vinculação com o tenant
		Username:  "updateduser",
		Name:      "Updated User",
		Email:     "updateduser@example.com",
	}

	// Configurando o mock para esperar o contexto, o UUID do usuário e o objeto usuário
	mockRepo.On("Update", mock.Anything, userID, &user).Return(&user, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: userID.String()}}
	c.Set(string(contextkeys.TenantIDKey), tenantID.String()) // Definindo o tenantID no contexto

	userData, _ := json.Marshal(user)
	c.Request = httptest.NewRequest("PUT", "/users/"+userID.String(), bytes.NewBuffer(userData))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Update(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var response models.User
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, user.Username, response.Username)
	assert.Equal(t, user.Name, response.Name)
	assert.Equal(t, user.Email, response.Email)
	mockRepo.AssertExpectations(t)
}

func TestUsersHandler_UpdatePartial(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	service := services.NewUserService(mockRepo)
	handler := handlers_v1.NewUsersHandler(service)

	userID := uuid.New()
	updateData := map[string]interface{}{
		"name": "Partially Updated User",
	}
	user := models.User{
		BaseModel: models.BaseModel{ID: userID},
		Name:      "Partially Updated User",
	}
	mockRepo.On("UpdatePartial", mock.Anything, userID, updateData).Return(&user, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: userID.String()}}

	updateDataBytes, _ := json.Marshal(updateData)
	c.Request = httptest.NewRequest("PATCH", "/users/"+userID.String(), bytes.NewBuffer(updateDataBytes))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.UpdatePartial(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var response models.User
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, updateData["name"], response.Name)
	mockRepo.AssertExpectations(t)
}

func TestUsersHandler_Delete(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	service := services.NewUserService(mockRepo)
	handler := handlers_v1.NewUsersHandler(service)

	userID := uuid.New()
	mockRepo.On("Delete", mock.Anything, userID).Return(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: userID.String()}}

	c.Request = httptest.NewRequest("DELETE", "/users/"+userID.String(), nil)

	handler.Delete(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockRepo.AssertExpectations(t)
}
