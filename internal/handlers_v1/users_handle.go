// internal/handlers_v1/users_handle.go

package handlers_v1

import (
	"fmt"
	"log"
	"net/http"

	contextkeys "github.com/jeancarlosdanese/go-base-api/internal/domain/context_keys"
	"github.com/jeancarlosdanese/go-base-api/internal/domain/models"
	"github.com/jeancarlosdanese/go-base-api/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UsersHandler struct holds the services that are needed.
type UsersHandler struct {
	service *services.UserService
}

func NewUsersHandler(service *services.UserService) *UsersHandler {
	return &UsersHandler{service: service}
}

// RegisterRoutes registra as rotas para users.
func (h *UsersHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("", h.getAll)
	router.GET("/:id", h.getById)
	router.POST("", h.create)
	router.PUT("/:id", h.update)
	router.PATCH("/:id", h.updatePatch)
	router.DELETE("/:id", h.delete)
}

// getAllUsers busca todos os Users
// @Summary Busca todos os Users
// @Description Busca todos os Users
// @Tags Users
// @Accept  json
// @Produce  json
// @Success 200 {array} models.User "Lista de Users"
// @Failure 500 {object} models.HTTPError "Erro Interno do Servidor"
// @Router /api/v1/users [get]
func (h *UsersHandler) getAll(c *gin.Context) {
	users, err := h.service.GetAll(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

// createUser cria um novo User
// @Summary Cria um novo User
// @Description Adiciona um novo User ao sistema
// @Tags Users
// @Accept json
// @Produce json
// @Param user body models.User true "Informações do User"
// @Success 201 {object} models.User "User Criado"
// @Failure 400 {object} models.HTTPError "Erro de Formato de Solicitação"
// @Failure 500 {object} models.HTTPError "Erro Interno do Servidor"
// @Router /api/v1/users [post]
func (h *UsersHandler) create(c *gin.Context) {
	tenantID, exists := c.Get(string(contextkeys.TenantIDKey))
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tenant não encontrado"})
	}

	// Verificar o tipo do tenantID
	fmt.Printf("Tipo de tenantID: %T\n", tenantID)

	var tenantUUID uuid.UUID
	switch v := tenantID.(type) {
	case string:
		uuidParsed, err := uuid.Parse(v)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao converter tenantID para uuid.UUID"})
		}
		tenantUUID = uuidParsed
	case uuid.UUID:
		tenantUUID = v
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "TenantID não é do tipo esperado (string ou uuid.UUID"})
	}

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user.TenantID = tenantUUID
	if err := h.service.Create(c, &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, user)
}

// getUserById busca um user pelo ID.
// @Summary Busca User por ID
// @Description Busca User por ID
// @Tags Users
// @Accept  json
// @Produce  json
// @Param   id     path    string     true        "User ID"
// @Success 200 {object} models.User "User"
// @Failure 404 {object} models.HTTPError "User not found"
// @Failure 400 {object} models.HTTPError "Invalid UUID format"
// @Router /api/v1/users/{id} [get]
func (h *UsersHandler) getById(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	user, err := h.service.GetByID(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// updateUser atualiza um user existente usando PUT.
// @Summary Atualiza um User existente
// @Description Atualiza um User existente com base no ID fornecido
// @Tags Users
// @Accept  json
// @Produce  json
// @Param   id     path    string     true        "User ID"
// @Param   user body    models.User true "Dados do User"
// @Success 200 {object} models.User "User Atualizado"
// @Failure 400 {object} models.HTTPError "ID Inválido ou Erro de Formato de Solicitação"
// @Failure 500 {object} models.HTTPError "Erro Interno do Servidor"
// @Router /api/v1/users/{id} [put]
func (h *UsersHandler) update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id")) // Extrair o ID do recurso da URL
	if err != nil {
		log.Fatalf("Invalid UUID: %v", err)
	}

	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Opcional: Definir o ID do user com o valor extraído da URL, garantindo que o recurso correto seja atualizado.
	user.ID = id

	if err := h.service.Update(c, &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

// updateUserPatch atualiza parcialmente um user existente usando PATCH.
// @Summary Atualiza parcialmente um User existente
// @Description Atualiza parcialmente um User existente com base no ID fornecido
// @Tags Users
// @Accept  json
// @Produce  json
// @Param   id     path    string     true        "User ID"
// @Param   user body    models.User true "Dados atualizáveis do User"
// @Success 200 {object} gin.H "Mensagem de sucesso"
// @Failure 400 {object} models.HTTPError "ID Inválido ou Erro de Formato de Solicitação"
// @Failure 500 {object} models.HTTPError "Erro Interno do Servidor"
// @Router /api/v1/users/{id} [patch]
func (h *UsersHandler) updatePatch(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id")) // Extrair o ID do recurso da URL
	if err != nil {
		log.Fatalf("Invalid UUID: %v", err)
	}

	var updateData map[string]interface{}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Remover campos que não devem ser atualizáveis
	delete(updateData, "cpf_cnpj")

	if err := h.service.UpdatePartial(c, id, updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

// deleteUser exclui um user.
// @Summary Exclui um User
// @Description Exclui um User com base no ID fornecido
// @Tags Users
// @Accept  json
// @Produce  json
// @Param   id     path    string     true        "User ID"
// @Success 200 {object} gin.H "Mensagem de sucesso"
// @Failure 400 {object} models.HTTPError "ID Inválido"
// @Failure 500 {object} models.HTTPError "Erro Interno do Servidor"
// @Router /api/v1/users/{id} [delete]
func (h *UsersHandler) delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		log.Fatalf("Invalid UUID: %v", err)
	}

	if err := h.service.Delete(c, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
