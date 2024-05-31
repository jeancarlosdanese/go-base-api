// internal/services/base_service.go

package services

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jeancarlosdanese/go-base-api/internal/repositories"
)

// Service define as operações básicas de um serviço com tipo genérico para entidade.
type BaseServiceInterface[Entity any] interface {
	Create(c *gin.Context, entity *Entity) (*Entity, error)
	Update(c *gin.Context, id uuid.UUID, entity *Entity) (*Entity, error)
	UpdatePartial(c *gin.Context, id uuid.UUID, updateData map[string]interface{}) (*Entity, error)
	Delete(c *gin.Context, id uuid.UUID) error
	GetAll(c *gin.Context) ([]Entity, error)
	GetByID(c *gin.Context, id uuid.UUID) (*Entity, error)
}

// BaseService implementa operações CRUD genéricas para qualquer entidade.
type BaseService[Entity any, Repo repositories.GormRepositoryInterface[Entity]] struct {
	Repo Repo
}

func NewBaseService[Entity any, Repo repositories.GormRepositoryInterface[Entity]](repo Repo) *BaseService[Entity, Repo] {
	return &BaseService[Entity, Repo]{Repo: repo}
}

func (s *BaseService[Entity, Repo]) Create(c *gin.Context, entity *Entity) (*Entity, error) {
	return s.Repo.Create(c, entity)
}

func (s *BaseService[Entity, Repo]) Update(c *gin.Context, id uuid.UUID, entity *Entity) (*Entity, error) {
	return s.Repo.Update(c, id, entity)
}

func (s *BaseService[Entity, Repo]) UpdatePartial(c *gin.Context, id uuid.UUID, updateData map[string]interface{}) (*Entity, error) {
	return s.Repo.UpdatePartial(c, id, updateData)
}

func (s *BaseService[Entity, Repo]) Delete(c *gin.Context, id uuid.UUID) error {
	return s.Repo.Delete(c, id)
}

func (s *BaseService[Entity, Repo]) GetAll(c *gin.Context) ([]Entity, error) {
	return s.Repo.GetAll(c)
}

func (s *BaseService[Entity, Repo]) GetByID(c *gin.Context, id uuid.UUID) (*Entity, error) {
	return s.Repo.GetByID(c, id)
}
