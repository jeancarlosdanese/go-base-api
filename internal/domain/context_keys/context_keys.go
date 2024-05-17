// internal/domain/context_keys/context_keys.go

package contextkeys

import (
	"log"

	"github.com/go-playground/validator/v10"
)

// ContextKey define os tipos possíveis de ações em uma entidade.
type ContextKey string

const (
	TenantIDKey  ContextKey = "TenantID"
	TokenDataKey ContextKey = "TokenData"
)

// ValidContextKeys mapeia as ações válidas para validação rápida.
var ValidContextKeys = map[ContextKey]bool{
	TenantIDKey:  true,
	TokenDataKey: true,
}

// validateContextKey verifica se o valor do ContextKey é um dos definidos como válidos.
func validateContextKey(fl validator.FieldLevel) bool {
	action, ok := fl.Field().Interface().(ContextKey)
	if !ok {
		log.Printf("Error: invalid data type for ContextKey field")
		return false
	}
	if _, exists := ValidContextKeys[action]; exists {
		return true
	}
	log.Printf("Validation failed: %s is not a valid ContextKey", action)
	return false
}

// IsValid verifica se o valor de ContextKey é válido.
func (a ContextKey) IsValid() bool {
	return ValidContextKeys[a]
}
