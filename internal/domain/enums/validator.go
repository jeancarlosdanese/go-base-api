// internal/domain/enums/validator.go

package enums

import "github.com/go-playground/validator/v10"

// Validator é a instância compartilhada do validador usada em todo o pacote enums.
var Validator *validator.Validate

// Validatable define a capacidade de um tipo validar seus próprios valores.
type Validatable interface {
	IsValid() bool
}

func Initialize() {
	Validator = validator.New() // Inicializa o validador uma única vez para todo o pacote
	Validator.RegisterValidation("roleType", validateRoleType)
	Validator.RegisterValidation("actionName", validateActionType)
	Validator.RegisterValidation("personType", validatePersonType)
	Validator.RegisterValidation("statusType", validateStatusType)
}
