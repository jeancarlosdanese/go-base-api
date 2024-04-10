// internal/domain/enums/status_type.go

package enums

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

// IsValid verifica se o valor de StatusType é válido.
func (s StatusType) IsValid() bool {
	_, ok := ValidStatusTypes[s]
	return ok
}
