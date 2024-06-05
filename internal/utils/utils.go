// internal/utils/utils.go

package utils

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jeancarlosdanese/go-base-api/internal/domain/models"
)

func GenerateApiKey(bytesNumber int) (string, error) {
	b := make([]byte, bytesNumber) // Gera 32, 64... bytes aleatórios.
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func GetTenantFromContext(c *gin.Context, key string) (*models.TenantRedis, bool) {
	tenantData, exists := c.Get(key)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant não encontrado"})
		c.Abort()
		return nil, false
	}

	tenant, ok := tenantData.(*models.TenantRedis)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Erro ao converter Tenant"})
		c.Abort()
		return nil, false
	}

	return tenant, true
}
