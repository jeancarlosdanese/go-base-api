package models_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/jeancarlosdanese/go-base-api/internal/domain/models"
	"github.com/stretchr/testify/assert"
)

func TestUser_Validation(t *testing.T) {
	user := models.User{
		TenantID: uuid.New(),
		Username: "testuser",
		Name:     "Test User",
		Email:    "testuser@example.com",
		Password: "password123",
	}

	// Verificar se os campos obrigatórios estão preenchidos
	assert.NotEmpty(t, user.TenantID, "TenantID não pode ser vazio")
	assert.NotEmpty(t, user.Username, "Username não pode ser vazio")
	assert.NotEmpty(t, user.Name, "Name não pode ser vazio")
	assert.NotEmpty(t, user.Email, "Email não pode ser vazio")
	assert.NotEmpty(t, user.Password, "Password não pode ser vazio")

	// Verificar se os campos opcionais podem ser vazios
	assert.Empty(t, user.Thumbnail, "Thumbnail deve ser vazio inicialmente")
}

func TestUser_ExtractRoles(t *testing.T) {
	roles := []*models.Role{
		{Name: "admin"},
		{Name: "user"},
	}
	user := models.User{Roles: roles}

	extractedRoles := user.ExtractRoles()

	assert.Contains(t, extractedRoles, "admin")
	assert.Contains(t, extractedRoles, "user")
	assert.Len(t, extractedRoles, 2)
}

func TestUser_ExtractPolicies(t *testing.T) {
	policies := []*models.PolicyUser{
		{Endpoint: &models.Endpoint{Name: "/api/v1/resource1"}, Actions: "GET"},
		{Endpoint: &models.Endpoint{Name: "/api/v1/resource2"}, Actions: "POST"},
	}
	user := models.User{SpecialPolicies: policies}

	extractedPolicies := user.ExtractPolicies()

	assert.Contains(t, extractedPolicies, "/api/v1/resource1:GET")
	assert.Contains(t, extractedPolicies, "/api/v1/resource2:POST")
	assert.Len(t, extractedPolicies, 2)
}
