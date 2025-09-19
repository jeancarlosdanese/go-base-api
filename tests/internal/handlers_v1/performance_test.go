package handlers_v1_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jeancarlosdanese/go-base-api/internal/routes"
	"github.com/stretchr/testify/assert"
)

func TestMiddlewarePerformance(t *testing.T) {
	// Criar router simples para teste de performance
	r := gin.New()

	// Aplicar middlewares
	r.Use(routes.SecurityHeadersMiddleware())
	r.Use(routes.RateLimitMiddleware(1000, time.Minute))

	r.GET("/test", func(c *gin.Context) {
		c.String(200, "OK")
	})

	// Testar múltiplas requisições
	start := time.Now()
	requestCount := 100

	for i := 0; i < requestCount; i++ {
		req, _ := http.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	}

	duration := time.Since(start)
	avgDuration := duration / time.Duration(requestCount)

	// Middleware deve ser rápido (< 10ms em média)
	assert.True(t, avgDuration < 10*time.Millisecond,
		"Middleware too slow: %v", avgDuration)
}

func TestConcurrentRequests(t *testing.T) {
	// Criar router simples para teste concorrente
	r := gin.New()
	r.Use(routes.SecurityHeadersMiddleware())

	r.GET("/test", func(c *gin.Context) {
		c.String(200, "OK")
	})

	concurrentRequests := 50

	// Channel para coletar resultados
	results := make(chan bool, concurrentRequests)

	// Função para fazer requisições
	makeRequest := func() {
		req, _ := http.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		results <- (w.Code == http.StatusOK)
	}

	// Executar requisições concorrentes
	for i := 0; i < concurrentRequests; i++ {
		go makeRequest()
	}

	// Aguardar resultados
	successCount := 0
	for i := 0; i < concurrentRequests; i++ {
		if <-results {
			successCount++
		}
	}

	// Todas as requisições devem ser bem-sucedidas
	assert.Equal(t, concurrentRequests, successCount,
		"Expected %d successful requests, got %d", concurrentRequests, successCount)
}

func BenchmarkSecurityMiddleware(b *testing.B) {
	r := gin.New()
	r.Use(routes.SecurityHeadersMiddleware())

	r.GET("/test", func(c *gin.Context) {
		c.String(200, "OK")
	})

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req, _ := http.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
		}
	})
}

func BenchmarkRateLimitingMiddleware(b *testing.B) {
	r := gin.New()
	r.Use(routes.RateLimitMiddleware(10000, time.Minute))

	r.GET("/test", func(c *gin.Context) {
		c.String(200, "OK")
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}
}
