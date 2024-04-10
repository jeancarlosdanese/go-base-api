// internal/handlers_v1/tenants_handle.go

package handlers_v1

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"hyberica.io/go/go-api/internal/domain/models"
	"hyberica.io/go/go-api/internal/services"
)

// TenantsHandler struct holds the services that are needed.
type TenantsHandler struct {
	service services.TenantService
}

// NewTenantsHandler creates a new TenantsHandler.
func NewTenantsHandler(service services.TenantService) *TenantsHandler {
	return &TenantsHandler{service: service}
}

// RegisterRoutes registra as rotas para tenants.
func (h *TenantsHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("", h.getAllTenants)
	router.GET("/:id", h.getTenantById)
	router.POST("", h.createTenant)
	router.PUT("/:id", h.updateTenant)
	router.PATCH("/:id", h.updateTenantPatch)
	router.DELETE("/:id", h.deleteTenant)
}

// getAllTenants busca todos os Tenans
func (h *TenantsHandler) getAllTenants(c *gin.Context) {
	ctx := context.Background()
	tenants, err := h.service.GetAllTenants(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tenants)
}

// createTenant cria um novo Tenant
func (h *TenantsHandler) createTenant(c *gin.Context) {
	var tenant models.Tenant
	if err := c.ShouldBindJSON(&tenant); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx := context.Background()
	if err := h.service.CreateTenant(ctx, &tenant); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, tenant)
}

// getTenantById busca um tenant pelo ID.
func (h *TenantsHandler) getTenantById(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		log.Fatalf("Invalid UUID: %v", err)
	}

	ctx := context.Background()
	tenant, err := h.service.GetTenantByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tenant not found"})
		return
	}
	c.JSON(http.StatusOK, tenant)
}

// updateTenant atualiza um tenant existente usando PUT.
func (h *TenantsHandler) updateTenant(c *gin.Context) {
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

	ctx := context.Background()
	if err := h.service.UpdateTenant(ctx, &tenant); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tenant)
}

// updateTenantPatch atualiza parcialmente um tenant existente usando PATCH.
func (h *TenantsHandler) updateTenantPatch(c *gin.Context) {
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

	ctx := context.Background()
	if err := h.service.UpdateTenantPartial(ctx, id, updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tenant updated successfully"})
}

// deleteTenant exclui um tenant.
func (h *TenantsHandler) deleteTenant(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		log.Fatalf("Invalid UUID: %v", err)
	}

	ctx := context.Background()
	if err := h.service.DeleteTenant(ctx, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tenant deleted successfully"})
}
