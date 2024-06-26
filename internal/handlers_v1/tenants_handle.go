// internal/handlers_v1/tenants_handle.go

package handlers_v1

import (
	"log"
	"net/http"

	"github.com/jeancarlosdanese/go-base-api/internal/domain/models"
	"github.com/jeancarlosdanese/go-base-api/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// TenantsHandler struct holds the services that are needed.
type TenantsHandler struct {
	tenantService services.TenantServiceInterface
}

func NewTenantsHandler(tenantService services.TenantServiceInterface) *TenantsHandler {
	return &TenantsHandler{tenantService: tenantService}
}

// RegisterRoutes registra as rotas para tenants.
func (h *TenantsHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("", h.GetAll)
	router.GET("/:id", h.GetById)
	router.POST("", h.Create)
	router.PUT("/:id", h.Update)
	router.PATCH("/:id", h.UpdatePatch)
	router.DELETE("/:id", h.Delete)
}

// getAllTenants busca todos os Tenants
// @Summary Busca todos os Tenants
// @Description Busca todos os Tenants
// @Tags Tenants
// @Accept  json
// @Produce  json
// @Success 200 {array} models.Tenant "Lista de Tenants"
// @Failure 500 {object} models.HTTPError "Erro Interno do Servidor"
// @Router /api/v1/tenants [get]
func (h *TenantsHandler) GetAll(c *gin.Context) {
	// ctx := context.Background()
	tenants, err := h.tenantService.GetAll(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tenants)
}

// createTenant cria um novo Tenant
// @Summary Cria um novo Tenant
// @Description Adiciona um novo Tenant ao sistema
// @Tags Tenants
// @Accept json
// @Produce json
// @Param tenant body models.Tenant true "Informações do Tenant"
// @Success 201 {object} models.Tenant "Tenant Criado"
// @Failure 400 {object} models.HTTPError "Erro de Formato de Solicitação"
// @Failure 500 {object} models.HTTPError "Erro Interno do Servidor"
// @Router /api/v1/tenants [post]
func (h *TenantsHandler) Create(c *gin.Context) {
	var tenantCreate models.Tenant
	if err := c.ShouldBindJSON(&tenantCreate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenant, err := h.tenantService.CreateTenantWithApiKey(c, &tenantCreate)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tenant not found"})
		return
	}
	c.JSON(http.StatusCreated, tenant)
}

// getTenantById busca um tenant pelo ID.
// @Summary Busca Tenant por ID
// @Description Busca Tenant por ID
// @Tags Tenants
// @Accept  json
// @Produce  json
// @Param   id     path    string     true        "Tenant ID"
// @Success 200 {object} models.Tenant "Tenant"
// @Failure 404 {object} models.HTTPError "Tenant not found"
// @Failure 400 {object} models.HTTPError "Invalid UUID format"
// @Router /api/v1/tenants/{id} [get]
func (h *TenantsHandler) GetById(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	tenant, err := h.tenantService.GetByID(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tenant not found"})
		return
	}
	c.JSON(http.StatusOK, tenant)
}

// updateTenant atualiza um tenant existente usando PUT.
// @Summary Atualiza um Tenant existente
// @Description Atualiza um Tenant existente com base no ID fornecido
// @Tags Tenants
// @Accept  json
// @Produce  json
// @Param   id     path    string     true        "Tenant ID"
// @Param   tenant body    models.Tenant true "Dados do Tenant"
// @Success 200 {object} models.Tenant "Tenant Atualizado"
// @Failure 400 {object} models.HTTPError "ID Inválido ou Erro de Formato de Solicitação"
// @Failure 500 {object} models.HTTPError "Erro Interno do Servidor"
// @Router /api/v1/tenants/{id} [put]
func (h *TenantsHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id")) // Extrair o ID do recurso da URL
	if err != nil {
		log.Fatalf("Invalid UUID: %v", err)
	}

	var tenant models.Tenant

	if err := c.ShouldBindJSON(&tenant); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Opcional: Definir o ID do tenant com o valor extraído da URL, garantindo que o recurso correto seja atualizado.
	tenant.ID = id

	tenantUpdated, err := h.tenantService.Update(c, id, &tenant)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tenant not found"})
		return
	}

	c.JSON(http.StatusOK, tenantUpdated)
}

// updateTenantPatch atualiza parcialmente um tenant existente usando PATCH.
// @Summary Atualiza parcialmente um Tenant existente
// @Description Atualiza parcialmente um Tenant existente com base no ID fornecido
// @Tags Tenants
// @Accept  json
// @Produce  json
// @Param   id     path    string     true        "Tenant ID"
// @Param   tenant body    models.Tenant true "Dados atualizáveis do Tenant"
// @Success 200 {object} gin.H "Mensagem de sucesso"
// @Failure 400 {object} models.HTTPError "ID Inválido ou Erro de Formato de Solicitação"
// @Failure 500 {object} models.HTTPError "Erro Interno do Servidor"
// @Router /api/v1/tenants/{id} [patch]
func (h *TenantsHandler) UpdatePatch(c *gin.Context) {
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
	// delete(updateData, "cpf_cnpj")

	tenantPatched, err := h.tenantService.UpdatePartial(c, id, updateData)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tenant not found"})
		return
	}

	c.JSON(http.StatusOK, tenantPatched)
}

// deleteTenant exclui um tenant.
// @Summary Exclui um Tenant
// @Description Exclui um Tenant com base no ID fornecido
// @Tags Tenants
// @Accept  json
// @Produce  json
// @Param   id     path    string     true        "Tenant ID"
// @Success 200 {object} gin.H "Mensagem de sucesso"
// @Failure 400 {object} models.HTTPError "ID Inválido"
// @Failure 500 {object} models.HTTPError "Erro Interno do Servidor"
// @Router /api/v1/tenants/{id} [delete]
func (h *TenantsHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		log.Fatalf("Invalid UUID: %v", err)
	}

	// ctx := context.Background()
	if err := h.tenantService.Delete(c, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tenant deleted successfully"})
}
