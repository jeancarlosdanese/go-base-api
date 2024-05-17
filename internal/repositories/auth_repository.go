// internal/repositories/auth_repository.go

package repositories

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	contextkeys "github.com/jeancarlosdanese/go-base-api/internal/domain/context_keys"
	"gorm.io/gorm"
)

// AuthRepository define as operações básicas de um repositório com tipo genérico para entidade.
type AuthRepository[Entity any] interface {
	Repository[Entity]
}

// NewGormAuthRepository cria uma nova instância de GormAuthRepository.
func NewGormAuthRepository[Entity any](db *gorm.DB) AuthRepository[Entity] {
	return &GormAuthRepository[Entity]{DB: db}
}

type GormAuthRepository[Entity any] struct {
	DB *gorm.DB
}

func (r *GormAuthRepository[Entity]) Create(c *gin.Context, entity *Entity) error {
	return r.DB.WithContext(c.Request.Context()).Create(entity).Error
}

func (r *GormAuthRepository[Entity]) Update(c *gin.Context, entity *Entity) error {
	tenantID, exists := c.Get(string(contextkeys.TenantIDKey))
	if !exists {
		return fmt.Errorf("tenant não encontado")
	}

	return r.DB.WithContext(c).Where("tenant_id = ?", tenantID).Save(entity).Error
}

func (r *GormAuthRepository[Entity]) UpdatePartial(c *gin.Context, id uuid.UUID, updateData map[string]interface{}) error {
	tenantID, exists := c.Get(string(contextkeys.TenantIDKey))
	if !exists {
		return fmt.Errorf("tenant não encontado")
	}

	entity := new(Entity)
	return r.DB.WithContext(c).Model(&entity).Where("tenant_id = ?", tenantID).Where("id = ?", id).Updates(updateData).Error
}

func (r *GormAuthRepository[Entity]) Delete(c *gin.Context, id uuid.UUID) error {
	tenantID, exists := c.Get(string(contextkeys.TenantIDKey))
	if !exists {
		return fmt.Errorf("tenant não encontado")
	}

	entity := new(Entity)
	return r.DB.WithContext(c).Where("tenant_id = ?", tenantID).Where("id = ?", id).Delete(entity).Error
}

func (r *GormAuthRepository[Entity]) GetAll(c *gin.Context) ([]Entity, error) {
	tenantID, exists := c.Get(string(contextkeys.TenantIDKey))
	if !exists {
		return nil, fmt.Errorf("tenant não encontado")
	}

	var entities []Entity
	err := r.DB.WithContext(c).Where("tenant_id = ?", tenantID).Find(&entities).Error
	if err != nil {
		return nil, err
	}
	return entities, nil
}

func (r *GormAuthRepository[Entity]) GetByID(c *gin.Context, id uuid.UUID) (*Entity, error) {
	tenantID, exists := c.Get(string(contextkeys.TenantIDKey))
	if !exists {
		return nil, fmt.Errorf("tenant não encontado")
	}

	var entity Entity
	err := r.DB.WithContext(c).Where("tenant_id = ?", tenantID).First(&entity, id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}
