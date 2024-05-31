// internal/utils/response.go

package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jeancarlosdanese/go-base-api/internal/logging"
)

func HandleAuthenticationError(c *gin.Context, err error) {
	logging.ErrorLogger.Printf("Erro ao autenticar usuário: %v", err)
	var httpStatus int
	var errorMsg string

	switch err.Error() {
	case "usuário ou origem não encontrado", "senha inválida", "invalid credentials":
		httpStatus = http.StatusUnauthorized
		errorMsg = "Credenciais inválidas"
	default:
		httpStatus = http.StatusInternalServerError
		errorMsg = "Erro interno do servidor"
	}

	c.JSON(httpStatus, gin.H{"error": errorMsg})
}
