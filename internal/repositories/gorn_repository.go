// internal/repositories/gorm_repository.go

package repositories

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Repository define as operações básicas de um repositório com tipo genérico para entidade.
type GormRepositoryInterface[Entity any] interface {
	Create(c *gin.Context, entity *Entity) (*Entity, error)
	Update(c *gin.Context, id uuid.UUID, entity *Entity) (*Entity, error)
	UpdatePartial(c *gin.Context, id uuid.UUID, updateData map[string]interface{}) (*Entity, error)
	Delete(c *gin.Context, id uuid.UUID) error
	GetAll(c *gin.Context) ([]Entity, error)
	GetByID(c *gin.Context, id uuid.UUID) (*Entity, error)
}

// NewGormRepository cria uma nova instância de GormRepository.
func NewGormRepository[Entity any](db *gorm.DB) GormRepositoryInterface[Entity] {
	return &GormRepository[Entity]{DB: db}
}

type GormRepository[Entity any] struct {
	DB *gorm.DB
}

func (r *GormRepository[Entity]) Create(c *gin.Context, entity *Entity) (*Entity, error) {
	err := r.DB.WithContext(c).Create(entity).Error
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (r *GormRepository[Entity]) Update(c *gin.Context, id uuid.UUID, entity *Entity) (*Entity, error) {
	// Inicia uma transação
	tx := r.DB.WithContext(c)

	// Primeiro, busca o tenant atual para assegurar que ele existe
	var existing Entity
	if err := tx.Where("id = ?", id).First(&existing).Error; err != nil {
		return nil, err // Retorna erro se o tenant não for encontrado
	}

	// Atualiza apenas os campos que foram realmente passados na requisição
	if err := tx.Model(&existing).Updates(entity).Error; err != nil {
		return nil, err
	}

	return &existing, nil
}

func (r *GormRepository[Entity]) UpdatePartial(c *gin.Context, id uuid.UUID, updateData map[string]interface{}) (*Entity, error) {
	var entity Entity // Cria uma referência para o tipo Entity

	err := r.DB.WithContext(c).
		Model(&entity).
		Clauses(clause.Returning{}).
		Where("id = ?", id).
		Updates(updateData).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *GormRepository[Entity]) Delete(c *gin.Context, id uuid.UUID) error {
	entity := new(Entity) // Cria uma referência para o tipo Entity
	return r.DB.WithContext(c).Where("id = ?", id).Delete(entity).Error
}

func (r *GormRepository[Entity]) GetAll(c *gin.Context) ([]Entity, error) {
	var entities []Entity
	err := r.DB.WithContext(c).Find(&entities).Error
	if err != nil {
		return nil, err
	}
	return entities, nil
}

func (r *GormRepository[Entity]) GetByID(c *gin.Context, id uuid.UUID) (*Entity, error) {
	var entity Entity
	err := r.DB.WithContext(c).First(&entity, id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}
