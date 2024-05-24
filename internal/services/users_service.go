// internal/services/users_service.go

package services

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/jeancarlosdanese/go-base-api/internal/domain/models"
	"github.com/jeancarlosdanese/go-base-api/internal/logging"
	"github.com/jeancarlosdanese/go-base-api/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	*BaseService[models.User, repositories.UserRepository]
}

func NewUserService(repo repositories.UserRepository) *UserService {
	baseService := NewBaseService[models.User, repositories.UserRepository](repo) // Tipos especificados aqui
	return &UserService{BaseService: baseService}
}

// Create sobrescreve o método Create para adicionar hashing de senha.
func (s *UserService) Create(c *gin.Context, userCreate models.UserCreate) (*models.User, error) {
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

	if err := s.Repo.Create(c, &user); err != nil {
		return nil, err
	}

	return &user, nil
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
