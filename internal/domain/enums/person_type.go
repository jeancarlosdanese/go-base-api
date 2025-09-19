// internal/domain/enums/person_type.go

package enums

import (
	"log"

	"github.com/go-playground/validator/v10"
)

// PersonType define os tipos possíveis para uma pessoa.
// @name PersonType
type PersonType string

const (
	Fisica   PersonType = "FISICA"
	Juridica PersonType = "JURIDICA"
)

// ValidPersonTypes oferece uma maneira de validar o PersonType.
var ValidPersonTypes = map[PersonType]bool{
	Fisica:   true,
	Juridica: true,
}

// validatePersonType verifica se o valor do PersonType é um dos definidos como válidos.
func validatePersonType(fl validator.FieldLevel) bool {
	personType, ok := fl.Field().Interface().(PersonType)
	if !ok {
		log.Printf("Error: invalid data type for PersonType field")
		return false
	}
	if _, exists := ValidPersonTypes[personType]; exists {
		return true
	}
	log.Printf("Validation failed: %s is not a valid PersonType", personType)
	return false
}

// IsValid verifica se o valor de PersonType é válido.
func (p PersonType) IsValid() bool {
	_, ok := ValidPersonTypes[p]
	return ok
}
