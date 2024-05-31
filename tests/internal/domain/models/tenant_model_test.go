// tests/internal/domain/models/tenant_model_test.go

package models_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/jeancarlosdanese/go-base-api/internal/domain/enums"
	"github.com/jeancarlosdanese/go-base-api/internal/domain/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/datatypes"
)

func TestTenant_Validation(t *testing.T) {
	allowedOrigins := datatypes.JSON([]byte(`["http://example.com"]`))
	tenant := models.Tenant{
		BaseModel:      models.BaseModel{ID: uuid.New()},
		Type:           enums.Fisica,
		Name:           "Test Tenant",
		AllowedOrigins: &allowedOrigins,
		Status:         enums.Ativo,
	}

	// Verificar se os campos obrigatórios estão preenchidos
	assert.NotEmpty(t, tenant.ID, "ID não pode ser vazio")
	assert.NotEmpty(t, tenant.Type, "Type não pode ser vazio")
	assert.NotEmpty(t, tenant.Name, "Name não pode ser vazio")
	assert.NotEmpty(t, tenant.Status, "Status não pode ser vazio")

	// Verificar se os campos opcionais podem ser vazios
	assert.Nil(t, tenant.CpfCnpj, "CpfCnpj deve ser vazio inicialmente")
	assert.Nil(t, tenant.Ie, "Ie deve ser vazio inicialmente")
	assert.Nil(t, tenant.Cep, "Cep deve ser vazio inicialmente")
	assert.Nil(t, tenant.Street, "Street deve ser vazio inicialmente")
	assert.Nil(t, tenant.Number, "Number deve ser vazio inicialmente")
	assert.Nil(t, tenant.Neighborhood, "Neighborhood deve ser vazio inicialmente")
	assert.Nil(t, tenant.City, "City deve ser vazio inicialmente")
	assert.Nil(t, tenant.State, "State deve ser vazio inicialmente")
	assert.Nil(t, tenant.Complement, "Complement deve ser vazio inicialmente")
	assert.Nil(t, tenant.Email, "Email deve ser vazio inicialmente")
	assert.Nil(t, tenant.Phone, "Phone deve ser vazio inicialmente")
	assert.Nil(t, tenant.CellPhone, "CellPhone deve ser vazio inicialmente")
}
