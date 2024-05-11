package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/jeancarlosdanese/go-base-api/internal/repositories"
)

// Service define as operações básicas de um serviço com tipo genérico para entidade.
type Service[Entity any] interface {
	Create(ctx context.Context, entity *Entity) error
	Update(ctx context.Context, entity *Entity) error
	UpdatePartial(ctx context.Context, id uuid.UUID, updateData map[string]interface{}) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetAll(ctx context.Context) ([]Entity, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Entity, error)
}

// BaseService implementa operações CRUD genéricas para qualquer entidade.
type BaseService[Entity any, Repo repositories.Repository[Entity]] struct {
	Repo Repo
}

func NewBaseService[Entity any, Repo repositories.Repository[Entity]](repo Repo) *BaseService[Entity, Repo] {
	return &BaseService[Entity, Repo]{Repo: repo}
}

func (s *BaseService[Entity, Repo]) Create(ctx context.Context, entity *Entity) error {
	return s.Repo.Create(ctx, entity)
}

func (s *BaseService[Entity, Repo]) Update(ctx context.Context, entity *Entity) error {
	return s.Repo.Update(ctx, entity)
}

func (s *BaseService[Entity, Repo]) UpdatePartial(ctx context.Context, id uuid.UUID, updateData map[string]interface{}) error {
	return s.Repo.UpdatePartial(ctx, id, updateData)
}

func (s *BaseService[Entity, Repo]) Delete(ctx context.Context, id uuid.UUID) error {
	return s.Repo.Delete(ctx, id)
}

func (s *BaseService[Entity, Repo]) GetAll(ctx context.Context) ([]Entity, error) {
	return s.Repo.GetAll(ctx)
}

func (s *BaseService[Entity, Repo]) GetByID(ctx context.Context, id uuid.UUID) (*Entity, error) {
	return s.Repo.GetByID(ctx, id)
}
