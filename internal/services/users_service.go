// internal/services/users_service.go

package services

import (
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jeancarlosdanese/go-base-api/internal/domain/models"
	"github.com/jeancarlosdanese/go-base-api/internal/logging"
	"github.com/jeancarlosdanese/go-base-api/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

// UserServiceInterface define as operações adicionais do UserService além das operações CRUD básicas.
type UserServiceInterface interface {
	BaseServiceInterface[models.User]
	CreateUserWithPassword(c *gin.Context, entity *models.UserCreate) (*models.User, error)
	Authenticate(c *gin.Context, email, password, origin string) (*models.User, error)
	GetOnlyByID(c *gin.Context, id uuid.UUID) (*models.User, error)
}
type UserService struct {
	*BaseService[models.User, repositories.UserRepository]
}

func NewUserService(repo repositories.UserRepository) *UserService {
	baseService := NewBaseService[models.User, repositories.UserRepository](repo) // Tipos especificados aqui
	return &UserService{BaseService: baseService}
}

// // Create sobrescreve o método Create para retornar um erro, alertando para o uso do méetodo CreateUserWithPassword.
// func (s *UserService) Create(c *gin.Context, user *models.User) (*models.User, error) {
// 	return nil, errors.New("use CreateUserWithPassword to create users")
// }

// CreateUserWithPassword é o método indicado para adicionar usuários, para fazer o hashing de senha.
func (s *UserService) CreateUserWithPassword(c *gin.Context, userCreate *models.UserCreate) (*models.User, error) {
	// Gera um hash para a senha do usuário
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userCreate.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Cria a entidade User com os dados do userCreateulário
	user := models.User{
		TenantID: userCreate.TenantID,
		Username: userCreate.Username,
		Name:     userCreate.Name,
		Email:    userCreate.Email,
		Password: string(hashedPassword),
	}

	userCreated, err := s.Repo.Create(c, &user)
	if err != nil {
		return nil, err
	}

	return userCreated, nil
}

// Authenticate verifica as credenciais de um usuário.
func (s *UserService) Authenticate(c *gin.Context, email, password, origin string) (*models.User, error) {
	user, err := s.Repo.FindByEmail(c, email, origin)
	if err != nil {
		logging.InfoLogger.Printf("Erro ao buscar usuário por email e origem: %v", err)
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		logging.InfoLogger.Printf("Senha inválida para o usuário: %s", email)
		return nil, errors.New("senha inválida")
	}

	return user, nil
}

// GetOnlyByID busca usuário apenas pelo seu ID
func (s *UserService) GetOnlyByID(c *gin.Context, id uuid.UUID) (*models.User, error) {
	user, err := s.Repo.GetOnlyByID(c, id)
	if err != nil {
		log.Println(err)
		logging.InfoLogger.Printf("Usuário não encontrado pelo ID: %v", id)
		return nil, errors.New("not found user by id")
	}

	return user, nil
}
