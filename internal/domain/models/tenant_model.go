// internal/domain/models/tenant_model.go

package models

import (
	"github.com/jeancarlosdanese/go-base-api/internal/domain/enums"
	"gorm.io/datatypes"
)

// Tenant representa os dados para criação de um Tenant.
type Tenant struct {
	BaseModel
	Type           enums.PersonType `gorm:"type:person_type;not null" validate:"required,personType" json:"type"`
	Name           string           `gorm:"type:varchar(100);not null" validate:"required" json:"name"`
	CpfCnpj        *string          `gorm:"type:varchar(18);unique" validate:"omitempty,cpfcnpj" json:"cpf_cnpj"`
	Ie             *string          `gorm:"type:varchar(20)" validate:"omitempty" json:"ie"`
	Cep            *string          `gorm:"type:varchar(9)" validate:"omitempty,len=9" json:"cep"`
	Street         *string          `gorm:"type:varchar(100)" validate:"omitempty" json:"street"`
	Number         *string          `gorm:"type:varchar(10)" validate:"omitempty,numeric" json:"number"`
	Neighborhood   *string          `gorm:"type:varchar(100)" validate:"omitempty" json:"neighborhood"`
	City           *string          `gorm:"type:varchar(100)" validate:"omitempty" json:"city"`
	State          *string          `gorm:"type:varchar(2)" validate:"omitempty,len=2" json:"state"`
	Complement     *string          `gorm:"type:varchar(100)" validate:"omitempty" json:"complement"`
	Email          *string          `gorm:"type:varchar(100)" validate:"omitempty,email" json:"email"`
	Phone          *string          `gorm:"type:varchar(15)" validate:"omitempty" json:"phone"`
	CellPhone      *string          `gorm:"type:varchar(15)" validate:"omitempty" json:"cell_phone"`
	ApiKey         *string          `gorm:"type:varchar(100)" validate:"omitempty" json:"api_key"`
	AllowedOrigins *datatypes.JSON  `gorm:"type:jsonb;unique" json:"allowed_origins"`
	Status         enums.StatusType `gorm:"type:status_type;not null;default:'ATIVO'" validate:"required,statusType" json:"status"`
}

type TenantRedis struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	CpfCnpj string `json:"cpfcnpj"`
	Email   string `json:"email"`
}
