// internal/domain/models/security_model.go

package models

import (
	"github.com/google/uuid"
	"github.com/jeancarlosdanese/go-base-api/internal/domain/enums"
)

// Entry representa um recurso no sistema que pode ser acessado por usuários.
type Entry struct {
	ID   uint   `gorm:"primarykey" validate:"required" json:"id"`
	Name string `gorm:"type:varchar(36);not null;unique" validate:"required" json:"name"` // Nome do recurso, único e não nulo
}

// Permission representa um papel no sistema.
type Permission struct {
	ID      uint             `gorm:"primarykey" validate:"required" json:"id"`
	EntryID uint             `gorm:"not null;uniqueIndex:uni_permissions_entry_id_action" validate:"required" json:"entry_id"`
	Action  enums.ActionType `gorm:"type:action_type;not null;validate:required,actionType;uniqueIndex:uni_permissions_entry_id_action" json:"action"`
	Roles   []*Role          `gorm:"many2many:permissions_roles;" json:"roles,omitempty"`
	Users   []*User          `gorm:"many2many:permissions_users;" json:"users,omitempty"`

	// constraints
	Entry *Entry `gorm:"foreignKey:EntryID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
}

// Role representa um papel no sistema.
type Role struct {
	ID          uint          `gorm:"primarykey" validate:"required" json:"id"`
	Name        string        `gorm:"type:varchar(36);not null;unique" validate:"required" json:"name"` // Nome do recurso, único e não nulo
	Permissions []*Permission `gorm:"many2many:permissions_roles;" json:"permissions,omitempty"`
}

type PermissionRole struct {
	PermissionID uint `gorm:"not null;primarykey" validate:"required" json:"permission_id"`
	RoleID       uint `gorm:"not null;primarykey" validate:"required" json:"role_id"`

	// constraints
	Permission *Permission `gorm:"foreignKey:PermissionID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
	Role       *Role       `gorm:"foreignKey:RoleID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
}

// TableName especifica o nome da tabela para GORM.
func (PermissionRole) TableName() string {
	return "permissions_roles"
}

// PermissionUser representa permissões especiais atribuídas a um usuário para um recurso específico.
type PermissionUser struct {
	PermissionID uint `gorm:"not null;primarykey" validate:"required" json:"permission_id"`
	UserID       uint `gorm:"not null;primarykey" validate:"required" json:"user_id"`

	// constraints
	Permission *Permission `gorm:"foreignKey:PermissionID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
	User       *User       `gorm:"foreignKey:UserID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
}

// TableName especifica o nome da tabela para GORM.
func (PermissionUser) TableName() string {
	return "permissions_users"
}

// UserRole representa permissões especiais atribuídas a um usuário para um recurso específico.
type UserRole struct {
	UserID uuid.UUID `gorm:"not null;primarykey" validate:"required" json:"user_id"`
	RoleID uint      `gorm:"not null;primarykey" validate:"required" json:"role_id"`

	// constraints
	User *User `gorm:"foreignKey:UserID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
	Role *Role `gorm:"foreignKey:RoleID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
}

// TableName especifica o nome da tabela para GORM.
func (UserRole) TableName() string {
	return "users_roles"
}
