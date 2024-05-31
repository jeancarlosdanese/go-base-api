// internal/repositories/auth_repository.go

package repositories

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	contextkeys "github.com/jeancarlosdanese/go-base-api/internal/domain/context_keys"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// AuthRepository define as operações básicas de um repositório com tipo genérico para entidade.

type AuthRepositoryInterface[Entity any] interface {
	GormRepositoryInterface[Entity]
}

// NewGormAuthRepository cria uma nova instância de GormAuthRepository.
func NewGormAuthRepository[Entity any](db *gorm.DB) AuthRepositoryInterface[Entity] {
	return &GormAuthRepository[Entity]{DB: db}
}

type GormAuthRepository[Entity any] struct {
	DB *gorm.DB
}

func (r *GormAuthRepository[Entity]) Create(c *gin.Context, entity *Entity) (*Entity, error) {
	err := r.DB.WithContext(c).Create(entity).Error
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (r *GormAuthRepository[Entity]) Update(c *gin.Context, id uuid.UUID, entity *Entity) (*Entity, error) {
	tenantID, exists := c.Get(string(contextkeys.TenantIDKey))
	if !exists {
		return nil, fmt.Errorf("tenant não encontrado")
	}

	fmt.Printf("\nGormAuthRepository: %v\n", entity)

	// Inclui a cláusula de tenantID na consulta de atualização
	err := r.DB.WithContext(c).Model(entity).Where("tenant_id = ?", tenantID).Where("id = ?", id).Save(entity).Error
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (r *GormAuthRepository[Entity]) UpdatePartial(c *gin.Context, id uuid.UUID, updateData map[string]interface{}) (*Entity, error) {
	tenantID, exists := c.Get(string(contextkeys.TenantIDKey))
	if !exists {
		return nil, fmt.Errorf("tenant não encontado")
	}

	var entity Entity // Cria uma referência para o tipo Entity

	err := r.DB.WithContext(c).
		Model(&entity).
		Clauses(clause.Returning{}).
		Where("tenant_id = ?", tenantID).
		Where("id = ?", id).
		Updates(updateData).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
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
