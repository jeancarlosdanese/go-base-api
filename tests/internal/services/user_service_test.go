// tests/internal/services/user_service_test.go

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

// MockUserRepository é um repositório mock para testes
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(c *gin.Context, entity *models.User) (*models.User, error) {
	args := m.Called(c, entity)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Update(c *gin.Context, id uuid.UUID, entity *models.User) (*models.User, error) {
	args := m.Called(c, id, entity)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) UpdatePartial(c *gin.Context, id uuid.UUID, updateData map[string]interface{}) (*models.User, error) {
	args := m.Called(c, id, updateData)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Delete(c *gin.Context, id uuid.UUID) error {
	args := m.Called(c, id)
	return args.Error(0)
}

func (m *MockUserRepository) GetAll(c *gin.Context) ([]models.User, error) {
	args := m.Called(c)
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserRepository) GetByID(c *gin.Context, id uuid.UUID) (*models.User, error) {
	args := m.Called(c, id)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(c *gin.Context, email, origin string) (*models.User, error) {
	args := m.Called(c, email, origin)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetOnlyByID(c *gin.Context, id uuid.UUID) (*models.User, error) {
	args := m.Called(c, id)
	return args.Get(0).(*models.User), args.Error(1)
}

func TestUserService_Create(t *testing.T) {
	repo := new(MockUserRepository)
	service := services.NewUserService(repo)

	c := &gin.Context{}

	userCreate := models.UserCreate{
		TenantID: uuid.New(),
		Username: "testuser",
		Name:     "Test User",
		Email:    "testuser@example.com",
		Password: "password123",
	}

	createdUser := &models.User{
		TenantID: userCreate.TenantID,
		Username: userCreate.Username,
		Name:     userCreate.Name,
		Email:    userCreate.Email,
		Password: "hashedpassword",
	}

	repo.On("Create", c, mock.AnythingOfType("*models.User")).Return(createdUser, nil)

	user, err := service.CreateUserWithPassword(c, &userCreate)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, userCreate.Username, user.Username)
	assert.Equal(t, userCreate.Name, user.Name)
	assert.Equal(t, userCreate.Email, user.Email)
	assert.NotEqual(t, userCreate.Password, user.Password)

	repo.AssertCalled(t, "Create", c, mock.AnythingOfType("*models.User"))
}

func TestUserService_Update(t *testing.T) {
	repo := new(MockUserRepository)
	service := services.NewUserService(repo)

	c := &gin.Context{}

	userID := uuid.New()
	userUpdate := models.User{
		BaseModel: models.BaseModel{ID: userID},
		Username:  "updateduser",
		Name:      "Updated User",
		Email:     "updateduser@example.com",
	}

	repo.On("Update", c, userID, mock.AnythingOfType("*models.User")).Return(&userUpdate, nil)

	user, err := service.Update(c, userID, &userUpdate)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, userUpdate.Username, user.Username)
	assert.Equal(t, userUpdate.Name, user.Name)
	assert.Equal(t, userUpdate.Email, user.Email)

	repo.AssertExpectations(t)
}

func TestUserService_UpdatePartial(t *testing.T) {
	repo := new(MockUserRepository)
	service := services.NewUserService(repo)

	c := &gin.Context{}

	userID := uuid.New()
	updateData := map[string]interface{}{
		"name": "Partially Updated User",
	}

	user := models.User{
		BaseModel: models.BaseModel{ID: userID},
		Name:      "Partially Updated User",
	}
	repo.On("UpdatePartial", c, userID, updateData).Return(&user, nil)

	updatedUser, err := service.UpdatePartial(c, userID, updateData)

	assert.NoError(t, err)
	assert.NotNil(t, updatedUser)
	assert.Equal(t, updateData["name"], updatedUser.Name)

	repo.AssertCalled(t, "UpdatePartial", c, userID, updateData)
}

func TestUserService_Delete(t *testing.T) {
	repo := new(MockUserRepository)
	service := services.NewUserService(repo)

	c := &gin.Context{}

	userID := uuid.New()
	repo.On("Delete", c, userID).Return(nil)

	err := service.Delete(c, userID)

	assert.NoError(t, err)

	repo.AssertCalled(t, "Delete", c, userID)
}

func TestUserService_GetAll(t *testing.T) {
	repo := new(MockUserRepository)
	service := services.NewUserService(repo)

	c := &gin.Context{}

	users := []models.User{
		{
			BaseModel: models.BaseModel{ID: uuid.New()},
			Username:  "user1",
		},
		{
			BaseModel: models.BaseModel{ID: uuid.New()},
			Username:  "user2",
		},
	}
	repo.On("GetAll", c).Return(users, nil)

	allUsers, err := service.GetAll(c)

	assert.NoError(t, err)
	assert.Len(t, allUsers, 2)
	assert.Equal(t, users, allUsers)

	repo.AssertCalled(t, "GetAll", c)
}

func TestUserService_GetByID(t *testing.T) {
	repo := new(MockUserRepository)
	service := services.NewUserService(repo)

	c := &gin.Context{}

	userID := uuid.New()
	user := models.User{
		BaseModel: models.BaseModel{ID: userID},
		Username:  "user1",
	}
	repo.On("GetByID", c, userID).Return(&user, nil)

	retrievedUser, err := service.GetByID(c, userID)

	assert.NoError(t, err)
	assert.NotNil(t, retrievedUser)
	assert.Equal(t, userID, retrievedUser.ID)
	assert.Equal(t, user.Username, retrievedUser.Username)

	repo.AssertCalled(t, "GetByID", c, userID)
}
