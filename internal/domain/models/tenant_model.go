// internal/domain/models/tenant_model.go

package models

import (
	"github.com/jeancarlosdanese/go-base-api/internal/domain/enums"
)

type Tenant struct {
	BaseModel
	Type         enums.PersonType `gorm:"type:person_type;not null" validate:"required,personType" json:"type"`
	Name         string           `gorm:"type:varchar(100);not null" validate:"required" json:"name"`
	CpfCnpj      *string          `gorm:"type:varchar(18);unique" validate:"omitempty,cpfcnpj" json:"cpf_cnpj"`
	Ie           *string          `gorm:"type:varchar(20)" validate:"omitempty" json:"ie"`
	Cep          *string          `gorm:"type:varchar(9)" validate:"omitempty,len=9" json:"cep"`
	Street       *string          `gorm:"type:varchar(100)" validate:"omitempty" json:"street"`
	Number       *string          `gorm:"type:varchar(10)" validate:"omitempty,numeric" json:"number"`
	Neighborhood *string          `gorm:"type:varchar(100)" validate:"omitempty" json:"neighborhood"`
	City         *string          `gorm:"type:varchar(100)" validate:"omitempty" json:"city"`
	State        *string          `gorm:"type:varchar(2)" validate:"omitempty,len=2" json:"state"`
	Complement   *string          `gorm:"type:varchar(100)" validate:"omitempty" json:"complement"`
	Email        *string          `gorm:"type:varchar(100)" validate:"omitempty,email" json:"email"`
	Phone        *string          `gorm:"type:varchar(15)" validate:"omitempty" json:"phone"`
	CellPhone    *string          `gorm:"type:varchar(15)" validate:"omitempty" json:"cell_phone"`
	Subdomain    string           `gorm:"type:varchar(20);not null;index:uni_tenants_subdomain_domain,unique" validate:"required" json:"subdomain"`
	Domain       *string          `gorm:"type:varchar(100);index:uni_tenants_subdomain_domain,unique" validate:"omitempty,url" json:"domain"`
	Status       enums.StatusType `gorm:"type:status_type;not null;default:'ATIVO'" validate:"required,statusType" json:"status"`
}
