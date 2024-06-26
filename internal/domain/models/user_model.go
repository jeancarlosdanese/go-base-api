// internal/domain/models/user_model.go

package models

import (
	"fmt"

	"github.com/google/uuid"
)

// User representa um usuário no sistema.
type User struct {
	BaseModel
	TenantID        uuid.UUID     `gorm:"type:uuid;not null;uniqueIndex:uni_users_tenant_id_email" json:"-"`
	Username        string        `gorm:"type:varchar(80);not null" json:"username"`
	Name            string        `gorm:"type:varchar(254);not null" json:"name"`
	Email           string        `gorm:"type:varchar(100);not null;uniqueIndex:uni_users_tenant_id_email" json:"email"`
	Password        string        `gorm:"type:varchar(60);not null" json:"-"`
	Thumbnail       *string       `gorm:"type:varchar(70)" json:"thumbnail"`
	Roles           []*Role       `gorm:"many2many:users_roles;" json:"roles"`
	SpecialPolicies []*PolicyUser `gorm:"many2many:policies_users;" json:"policies"`

	// constraints
	Tenant *Tenant `gorm:"foreignKey:TenantID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;" json:"tenant,omitempty"`
}

// UserCreate é usado para receber dados do formulário de criação de usuário.
type UserCreate struct {
	TenantID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:uni_users_tenant_id_email" json:"tenant_id"`
	Username string    `gorm:"type:varchar(80);not null" json:"username"`
	Name     string    `gorm:"type:varchar(254);not null" json:"name"`
	Email    string    `gorm:"type:varchar(100);not null;uniqueIndex:uni_users_tenant_id_email" json:"email"`
	Password string    `gorm:"type:varchar(60);not null" json:"password"`
}

// ExtractRoles extrai e retorna os nomes dos roles do usuário.
func (u *User) ExtractRoles() []string {
	var roles []string
	for _, role := range u.Roles {
		roles = append(roles, role.Name)
	}
	return roles
}

// ExtractPolicies extrai e retorna as permissões especiais formatadas do usuário.
func (u *User) ExtractPolicies() []string {
	var policies []string
	for _, policy := range u.SpecialPolicies {
		if policy.Endpoint != nil {
			// Como Actions já é uma string, você pode usá-la diretamente
			policies = append(policies, fmt.Sprintf("%s:%s", policy.Endpoint.Name, policy.Actions))
		}
	}
	return policies
}

type UserRedis struct {
	ID       string   `json:"id"`
	TenantID string   `json:"tenant_id"`
	Name     string   `json:"name"`
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Roles    []string `json:"roles"`
	Policies []string `json:"policies"`
}
