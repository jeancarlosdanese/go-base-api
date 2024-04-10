// internal/domain/enums/person_type.go

package enums

// PersonType define os tipos possíveis para uma pessoa.
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

// IsValid verifica se o valor de PersonType é válido.
func (t PersonType) IsValid() bool {
	_, ok := ValidPersonTypes[t]
	return ok
}
