// internal/domain/enums/action_name.go

package enums

import (
	"log"

	"github.com/go-playground/validator/v10"
)

// RoleType define os tipos possíveis de ações em uma entidade.
type RoleType string

const (
	Admin       RoleType = "admin"
	Coordinator RoleType = "coordinator"
	Secretary   RoleType = "secretary"
	Student     RoleType = "student"
)

// ValidRoleTypes mapeia as ações válidas para validação rápida.
var ValidRoleTypes = map[RoleType]bool{
	Admin:       true,
	Coordinator: true,
	Secretary:   true,
	Student:     true,
}

// validateRoleType verifica se o valor do RoleType é um dos definidos como válidos.
func validateRoleType(fl validator.FieldLevel) bool {
	action, ok := fl.Field().Interface().(RoleType)
	if !ok {
		log.Printf("Error: invalid data type for RoleType field")
		return false
	}
	if _, exists := ValidRoleTypes[action]; exists {
		return true
	}
	log.Printf("Validation failed: %s is not a valid RoleType", action)
	return false
}

// IsValid verifica se o valor de RoleType é válido.
func (a RoleType) IsValid() bool {
	return ValidRoleTypes[a]
}
