// internal/handlers_v1/auth_apikey_handle.go

package handlers_v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	contextkeys "github.com/jeancarlosdanese/go-base-api/internal/domain/context_keys"
	_ "github.com/jeancarlosdanese/go-base-api/internal/domain/models"
	"github.com/jeancarlosdanese/go-base-api/internal/utils"
)

// AuthApiKeyHandler struct holds the services that are needed.
type AuthApiKeyHandler struct{}

func NewAuthApiKeyHandler() *AuthApiKeyHandler {
	return &AuthApiKeyHandler{}
}

// RegisterRoutes registra as rotas para users.
func (h *AuthApiKeyHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/tenant-by-apikey", h.GetTenantByApiKey)
}

// getTenantByApiKey busca um tenant pelo X-API-Key.
// @Summary Busca Tenant por X-API-Key
// @Description Busca Tenant por X-API-Key
// @Tags auth-apikey
// @Accept  json
// @Produce  json
// @Success 200 {object} models.TenantRedis "Tenant"
// @Failure 404 {object} models.HTTPError "Tenant not found"
// @Failure 400 {object} models.HTTPError "Invalid X-API-Key format"
// @Router /api/v1/auth-apikey/tenant-by-apikey [get]
func (h *AuthApiKeyHandler) GetTenantByApiKey(c *gin.Context) {
	origin := c.GetString("Origin")
	if origin == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Origem não fornecida",
			"hint":  "Certifique-se de que o cliente HTTP está enviando o header Origin ou configure localhost como padrão",
		})
		return
	}

	tenant, ok := utils.GetTenantFromContext(c, string(contextkeys.TenantDataKey))
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"error":  "Tenant não encontrado",
			"hint":   "Verifique se a API Key é válida e está associada a um tenant",
			"origin": origin,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"origin": origin,
		"tenant": tenant,
		"status": "success",
	})
}
