// internal/domain/models/security_model.go

package models

import (
	"github.com/google/uuid"
)

// Endpoint representa um recurso no sistema que pode ser acessado por usuários.
type Endpoint struct {
	ID   uint   `gorm:"primarykey" validate:"required" json:"id"`
	Name string `gorm:"type:varchar(254);not null;unique" validate:"required" json:"name"` // Nome do recurso, único e não nulo
}

// Role representa um papel no sistema.
type Role struct {
	ID   uint   `gorm:"primarykey" validate:"required" json:"id"`
	Name string `gorm:"type:varchar(36);not null;unique" validate:"required" json:"name"` // Nome do recurso, único e não nulo

	Policies []*PolicyRole `gorm:"many2many:policies_roles;" json:"policies,omitempty"`
}

type PolicyRole struct {
	RoleID     uint   `gorm:"not null;primarykey;" validate:"required" json:"role_id"`
	EndpointID uint   `gorm:"not null;primarykey;" validate:"required" json:"endpoint_id"`
	Actions    string `gorm:"type:varchar(36);not null;" validate:"required" json:"actions"`

	// // constraints
	Role     *Role     `gorm:"foreignKey:RoleID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
	Endpoint *Endpoint `gorm:"foreignKey:EndpointID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
}

// TableName especifica o nome da tabela para GORM.
func (PolicyRole) TableName() string {
	return "policies_roles"
}

// PolicyUser representa permissões especiais atribuídas a um usuário para um recurso específico.
type PolicyUser struct {
	UserID     uuid.UUID `gorm:"not null;primarykey;" validate:"required" json:"user_id"`
	EndpointID uint      `gorm:"not null;primarykey;" validate:"required" json:"endpoint_id"`
	Actions    string    `gorm:"type:varchar(36);not null;" validate:"required" json:"actions"`

	// // constraints
	User     *User     `gorm:"foreignKey:UserID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
	Endpoint *Endpoint `gorm:"foreignKey:EndpointID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
}

// TableName especifica o nome da tabela para GORM.
func (PolicyUser) TableName() string {
	return "policies_users"
}

// UserRole representa permissões especiais atribuídas a um usuário para um recurso específico.
type UserRole struct {
	UserID uuid.UUID `gorm:"not null;primarykey;" validate:"required" json:"user_id"`
	RoleID uint      `gorm:"not null;primarykey;" validate:"required" json:"role_id"`

	// constraints
	User *User `gorm:"foreignKey:UserID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
	Role *Role `gorm:"foreignKey:RoleID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
}

// TableName especifica o nome da tabela para GORM.
func (UserRole) TableName() string {
	return "users_roles"
}
