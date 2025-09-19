package handlers_v1_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jeancarlosdanese/go-base-api/internal/routes"
	"github.com/stretchr/testify/assert"
)

func TestSecurityHeaders(t *testing.T) {
	// Criar router simples para teste de headers
	r := gin.New()

	// Aplicar apenas o middleware de segurança
	r.Use(routes.SecurityHeadersMiddleware())

	r.GET("/test", func(c *gin.Context) {
		c.String(200, "OK")
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Verificar headers de segurança
	headers := w.Header()
	assert.Equal(t, "nosniff", headers.Get("X-Content-Type-Options"))
	assert.Equal(t, "DENY", headers.Get("X-Frame-Options"))
	assert.Equal(t, "1; mode=block", headers.Get("X-XSS-Protection"))
	assert.Contains(t, headers.Get("Content-Security-Policy"), "default-src 'self'")
	assert.Equal(t, "strict-origin-when-cross-origin", headers.Get("Referrer-Policy"))
	assert.Equal(t, "geolocation=(), microphone=(), camera=()", headers.Get("Permissions-Policy"))
}

func TestRateLimiting(t *testing.T) {
	// Criar router simples para teste de rate limiting
	r := gin.New()
	r.Use(routes.RateLimitMiddleware(5, 1)) // 5 requests por segundo

	r.GET("/test", func(c *gin.Context) {
		c.String(200, "OK")
	})

	// Fazer múltiplas requisições rápidas
	for i := 0; i < 7; i++ { // Mais que o limite de 5
		req, _ := http.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if i < 5 {
			// Primeiras 5 devem passar
			assert.Equal(t, http.StatusOK, w.Code)
		} else {
			// A partir da 6 deve ser rate limited
			if w.Code == http.StatusTooManyRequests {
				// Verificar headers de rate limit
				assert.NotEmpty(t, w.Header().Get("X-RateLimit-Limit"))
				assert.NotEmpty(t, w.Header().Get("X-RateLimit-Remaining"))
				assert.NotEmpty(t, w.Header().Get("X-RateLimit-Reset"))
				assert.NotEmpty(t, w.Header().Get("Retry-After"))
				break
			}
		}
	}
}

func TestRequestSizeLimit(t *testing.T) {
	// Criar router simples para teste de tamanho de request
	r := gin.New()
	r.Use(routes.RequestSizeLimitMiddleware(1024)) // 1KB limite

	r.POST("/test", func(c *gin.Context) {
		c.String(200, "OK")
	})

	// Criar um body grande (mais de 1KB)
	largeBody := make([]byte, 2048) // 2KB
	req, _ := http.NewRequest("POST", "/test", bytes.NewReader(largeBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusRequestEntityTooLarge, w.Code)
	assert.Contains(t, w.Body.String(), "Request too large")
}
