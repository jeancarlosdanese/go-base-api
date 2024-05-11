// internal/domain/enums/status_type.go

package enums

import (
	"log"

	"github.com/go-playground/validator/v10"
)

// StatusType define os status possíveis para uma entidade.
type StatusType string

const (
	Ativo   StatusType = "ATIVO"
	Inativo StatusType = "INATIVO"
)

// ValidStatusTypes fornece uma maneira de validar o StatusType.
var ValidStatusTypes = map[StatusType]bool{
	Ativo:   true,
	Inativo: true,
}

// validateStatusType verifica se o valor do StatusType é um dos definidos como válidos.
func validateStatusType(fl validator.FieldLevel) bool {
	personType, ok := fl.Field().Interface().(StatusType)
	if !ok {
		log.Printf("Error: invalid data type for StatusType field")
		return false
	}
	if _, exists := ValidStatusTypes[personType]; exists {
		return true
	}
	log.Printf("Validation failed: %s is not a valid StatusType", personType)
	return false
}

// IsValid verifica se o valor de StatusType é válido.
func (p StatusType) IsValid() bool {
	_, ok := ValidStatusTypes[p]
	return ok
}
