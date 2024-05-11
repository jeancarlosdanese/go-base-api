// internal/repositories/gorm_repository.go

package repositories

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository define as operações básicas de um repositório com tipo genérico para entidade.
type Repository[Entity any] interface {
	// Create insere uma nova entidade no banco de dados.
	Create(ctx context.Context, entity *Entity) error

	// Update atualiza uma entidade existente no banco de dados.
	Update(ctx context.Context, entity *Entity) error

	// UpdatePartial atualiza parcialmente uma entidade existente no banco de dados.
	UpdatePartial(ctx context.Context, id uuid.UUID, updateData map[string]interface{}) error

	// Delete remove uma entidade do banco de dados por ID.
	Delete(ctx context.Context, id uuid.UUID) error

	// GetAll recupera todos os itens do repositório.
	GetAll(ctx context.Context) ([]Entity, error)

	// GetByID busca uma entidade específica pelo ID.
	GetByID(ctx context.Context, id uuid.UUID) (*Entity, error)
}

// NewGormRepository cria uma nova instância de GormRepository.
func NewGormRepository[Entity any](db *gorm.DB) Repository[Entity] {
	return &GormRepository[Entity]{DB: db}
}

type GormRepository[Entity any] struct {
	DB *gorm.DB
}

func (r *GormRepository[Entity]) Create(ctx context.Context, entity *Entity) error {
	return r.DB.WithContext(ctx).Create(entity).Error
}

func (r *GormRepository[Entity]) Update(ctx context.Context, entity *Entity) error {
	return r.DB.WithContext(ctx).Save(entity).Error
}

func (r *GormRepository[Entity]) UpdatePartial(ctx context.Context, id uuid.UUID, updateData map[string]interface{}) error {
	entity := new(Entity) // Cria uma referência para o tipo Entity
	return r.DB.WithContext(ctx).Model(&entity).Where("id = ?", id).Updates(updateData).Error
}

func (r *GormRepository[Entity]) Delete(ctx context.Context, id uuid.UUID) error {
	entity := new(Entity) // Cria uma referência para o tipo Entity
	return r.DB.WithContext(ctx).Where("id = ?", id).Delete(entity).Error
}

func (r *GormRepository[Entity]) GetAll(ctx context.Context) ([]Entity, error) {
	var entities []Entity
	err := r.DB.WithContext(ctx).Find(&entities).Error
	if err != nil {
		return nil, err
	}
	return entities, nil
}

func (r *GormRepository[Entity]) GetByID(ctx context.Context, id uuid.UUID) (*Entity, error) {
	var entity Entity
	err := r.DB.WithContext(ctx).First(&entity, id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}
