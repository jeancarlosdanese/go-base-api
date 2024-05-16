// internal/domain/enums/action_name.go

package enums

import (
	"log"

	"github.com/go-playground/validator/v10"
)

// ActionType define os tipos possíveis de ações em uma entidade.
type ActionType string

const (
	GET    ActionType = "GET"
	POST   ActionType = "POST"
	PUT    ActionType = "PUT"
	PATCH  ActionType = "PATCH"
	DELETE ActionType = "DELETE"
)

// ValidActionTypes mapeia as ações válidas para validação rápida.
var ValidActionTypes = map[ActionType]bool{
	GET:    true,
	POST:   true,
	PUT:    true,
	PATCH:  true,
	DELETE: true,
}

// validateActionType verifica se o valor do ActionType é um dos definidos como válidos.
func validateActionType(fl validator.FieldLevel) bool {
	action, ok := fl.Field().Interface().(ActionType)
	if !ok {
		log.Printf("Error: invalid data type for ActionType field")
		return false
	}
	if _, exists := ValidActionTypes[action]; exists {
		return true
	}
	log.Printf("Validation failed: %s is not a valid ActionType", action)
	return false
}

// IsValid verifica se o valor de ActionType é válido.
func (a ActionType) IsValid() bool {
	return ValidActionTypes[a]
}
