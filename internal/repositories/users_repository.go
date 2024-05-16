// internal/repositories/repository_users.go

package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/jeancarlosdanese/go-base-api/internal/domain/models"
	"gorm.io/gorm"
)

// UserRepository é uma interface que estende a interface Repository para operações específicas do User.
type UserRepository interface {
	Repository[models.User]
	FindByEmail(ctx context.Context, email, origin string) (*models.User, error)
}

// NewUserRepository cria uma nova instância de um repositório que implementa UserRepository.
func NewUserRepository(db *gorm.DB) UserRepository {
	repo := NewGormRepository[models.User](db).(UserRepository) // Asserção de tipo
	return repo
}

func (r *GormRepository[Entity]) FindByEmail(ctx context.Context, email, origin string) (*models.User, error) {
	var user models.User
	err := r.DB.WithContext(ctx).
		// Preload o Tenant apenas se a origem estiver na lista de allowed_origins do Tenant.
		// Usando JSONB para verificar se a origem está contida na lista allowed_origins.
		Preload("Tenant", "allowed_origins @> ?", fmt.Sprintf(`["%s"]`, origin)).
		Preload("Roles.Policies.Endpoint").
		Where("email = ?", email).
		First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Usuário não encontrado
		}
		return nil, err // Outro erro
	}

	if err := r.DB.WithContext(ctx).
		Preload("Endpoint").
		Where("user_id = ?", user.ID).
		Find(&user.SpecialPolicies).Error; err != nil {
		return nil, err // Erro ao carregar special policies
	}

	// Verificar se o tenant foi carregado (validar origem)
	if user.Tenant == nil {
		return nil, errors.New("origem inválida")
	}

	return &user, nil
}
