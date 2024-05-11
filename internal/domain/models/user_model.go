// internal/domain/models/user_model.go

package models

import (
	"fmt"

	"github.com/google/uuid"
)

// User representa um usuário no sistema.
type User struct {
	BaseModel
	TenantID           uuid.UUID     `gorm:"type:uuid;not null;uniqueIndex:uni_users_tenant_id_email" json:"tenant_id"`
	Username           string        `gorm:"type:varchar(80);not null" json:"username"`
	Name               string        `gorm:"type:varchar(254);not null" json:"name"`
	Email              string        `gorm:"type:varchar(100);not null;uniqueIndex:uni_users_tenant_id_email" json:"email"`
	Password           string        `gorm:"type:varchar(60);not null" json:"-"` // Omitindo senha no JSON
	Thumbnail          string        `gorm:"type:varchar(70);" validate:"omitempty" json:"thumbnail"`
	Roles              []*Role       `gorm:"many2many:users_roles;" json:"roles,omitempty"`
	SpecialPermissions []*Permission `gorm:"many2many:permissions_users;" json:"special_permissions,omitempty"`

	// constraints
	Tenant *Tenant `gorm:"foreignKey:TenantID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
}

// ExtractRoles extrai e retorna os nomes dos roles do usuário.
func (u *User) ExtractRoles() []string {
	var roles []string
	for _, role := range u.Roles {
		roles = append(roles, role.Name)
	}
	return roles
}

// ExtractPermissions extrai e retorna as permissões especiais formatadas do usuário.
func (u *User) ExtractPermissions() []string {
	var permissions []string
	for _, permission := range u.SpecialPermissions {
		permissions = append(permissions, fmt.Sprintf("%s:%s", permission.Entry.Name, permission.Action))
	}
	return permissions
}

type UserDataRedis struct {
	Token        string    `json:"token"`
	RefreshToken *string   `json:"refreshToken"`
	User         UserRedis `json:"user"`
	Roles        []string  `json:"roles"`
	Permissions  []string  `json:"permissions"`
}

type UserRedis struct {
	ID       string `json:"id"`
	TenantID string `json:"tenant_id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
}
