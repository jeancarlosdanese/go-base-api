// internal/utils/converters.go

package utils

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// TryParseUUID tenta converter uma entrada interface{} para uuid.UUID.
func TryParseUUID(c *gin.Context, input interface{}) (uuid.UUID, error) {
	switch v := input.(type) {
	case string:
		parsedUUID, err := uuid.Parse(v)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao converter string para uuid.UUID"})
			return uuid.UUID{}, err
		}
		return parsedUUID, nil
	case uuid.UUID:
		return v, nil
	default:
		err := errors.New("entrada não é do tipo esperado (string ou uuid.UUID)")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return uuid.UUID{}, err
	}
}
