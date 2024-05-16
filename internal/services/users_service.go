// internal/services/users_service.go

package services

import (
	"context"
	"errors"

	"github.com/jeancarlosdanese/go-base-api/internal/domain/models"
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

// func (s *UserService) FindByEmail(ctx context.Context, email string) (*models.User, error) {
// 	return s.Repo.FindByEmail(ctx, email)
// }

// Authenticate verifica as credenciais de um usuário.
func (s *UserService) Authenticate(ctx context.Context, email, password, origin string) (*models.User, error) {
	user, err := s.Repo.FindByEmail(ctx, email, origin) // Certifique-se que FindByEmail está implementado em UserRepository
	if err != nil {
		return nil, errors.New("usuário não encontrado")
	}

	// Comparar a senha criptografada com a senha fornecida
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("senha inválida")
	}

	return user, nil
}
