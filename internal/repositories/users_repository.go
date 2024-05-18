// internal/repositories/repository_users.go

package repositories

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jeancarlosdanese/go-base-api/internal/domain/models"
	"github.com/jeancarlosdanese/go-base-api/internal/logging"
	"gorm.io/gorm"
)

// UserRepository é uma interface que estende a interface Repository para operações específicas do User.
type UserRepository interface {
	Repository[models.User]
	FindByEmail(c *gin.Context, email, origin string) (*models.User, error)
	GetOnlyByID(c *gin.Context, id uuid.UUID) (*models.User, error)
}

// NewUserRepository cria uma nova instância de um repositório que implementa UserRepository.
func NewUserRepository(db *gorm.DB) UserRepository {
	repo := NewGormAuthRepository[models.User](db).(UserRepository) // Asserção de tipo
	return repo
}

func (r *GormAuthRepository[Entity]) FindByEmail(c *gin.Context, email, origin string) (*models.User, error) {
	var user models.User
	formattedOrigin := fmt.Sprintf(`["%s"]`, origin)
	err := r.DB.WithContext(c).
		Preload("Roles.Policies.Endpoint").
		Where("email = ? AND EXISTS (SELECT 1 FROM tenants WHERE tenants.id = users.tenant_id AND allowed_origins @> ?)", email, formattedOrigin).
		Take(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logging.InfoLogger.Printf("Usuário ou origem não encontrado: %s, %s", email, origin)
			return nil, errors.New("usuário ou origem não encontrado")
		}
		logging.ErrorLogger.Printf("Erro ao buscar usuário por email e origem: %v", err)
		return nil, err
	}

	if err := r.DB.WithContext(c).
		Preload("Endpoint").
		Where("user_id = ?", user.ID).
		Find(&user.SpecialPolicies).Error; err != nil {
		logging.ErrorLogger.Printf("Erro ao carregar special policies para o usuário: %v", err)
		return nil, err
	}

	return &user, nil
}

func (r *GormAuthRepository[Entity]) GetOnlyByID(c *gin.Context, id uuid.UUID) (*Entity, error) {
	fmt.Printf("UserID: %v", id)

	var entity Entity
	err := r.DB.WithContext(c).First(&entity, id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}
