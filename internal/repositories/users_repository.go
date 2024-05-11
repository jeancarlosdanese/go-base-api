// internal/repositories/repository_users.go

package repositories

import (
	"context"

	"github.com/jeancarlosdanese/go-base-api/internal/domain/models"
	"gorm.io/gorm"
)

// UserRepository é uma interface que estende a interface Repository para operações específicas do User.
type UserRepository interface {
	Repository[models.User]
	FindByEmail(ctx context.Context, email string) (*models.User, error)
}

// NewUserRepository cria uma nova instância de um repositório que implementa UserRepository.
func NewUserRepository(db *gorm.DB) UserRepository {
	repo := NewGormRepository[models.User](db).(UserRepository) // Asserção de tipo
	return repo
}

func (r *GormRepository[Entity]) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var entity models.User
	err := r.DB.WithContext(ctx).
		Preload("Roles").
		Preload("Roles.Permissions").
		Preload("Roles.Permissions.Entry").
		Preload("SpecialPermissions").
		Preload("SpecialPermissions.Entry").
		Where("email = ?", email).
		First(&entity).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Usuário não encontrado
		}
		return nil, err // Outro erro
	}
	return &entity, nil
}
