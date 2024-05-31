package mocks

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jeancarlosdanese/go-base-api/internal/domain/models"
	"github.com/stretchr/testify/mock"
)

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
