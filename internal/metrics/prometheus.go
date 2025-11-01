// @file: internal/metrics/prometheus.go

package metrics

import (
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// httpRequestsTotal é o contador total de requisições HTTP
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	// httpRequestDurationSeconds é o histograma de duração das requisições HTTP
	httpRequestDurationSeconds = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets, // Usa buckets padrão: .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10
		},
		[]string{"method", "path", "status"},
	)

	// httpRequestSizeBytes é o histograma do tamanho das requisições HTTP
	httpRequestSizeBytes = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_size_bytes",
			Help:    "Size of HTTP requests in bytes",
			Buckets: []float64{100, 500, 1000, 5000, 10000, 50000, 100000, 500000, 1000000, 5000000},
		},
		[]string{"method", "path"},
	)

	// httpResponseSizeBytes é o histograma do tamanho das respostas HTTP
	httpResponseSizeBytes = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_size_bytes",
			Help:    "Size of HTTP responses in bytes",
			Buckets: []float64{100, 500, 1000, 5000, 10000, 50000, 100000, 500000, 1000000, 5000000},
		},
		[]string{"method", "path", "status"},
	)

	// httpInFlightRequests é o gauge de requisições em andamento
	httpInFlightRequests = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_in_flight_requests",
			Help: "Number of HTTP requests currently being processed",
		},
	)
)

// Middleware retorna um middleware Gin que coleta métricas Prometheus
func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ignora endpoints de métricas e health check para evitar poluição de métricas
		if c.Request.URL.Path == "/metrics" || c.Request.URL.Path == "/health" {
			c.Next()
			return
		}

		start := time.Now()

		// Incrementa contador de requisições em andamento
		httpInFlightRequests.Inc()
		defer httpInFlightRequests.Dec()

		// Captura tamanho da requisição
		requestSize := float64(c.Request.ContentLength)
		if requestSize < 0 {
			requestSize = 0
		}

		// Processa a requisição
		c.Next()

		// Calcula duração e tamanho da resposta
		duration := time.Since(start).Seconds()
		statusCode := strconv.Itoa(c.Writer.Status())
		method := c.Request.Method
		path := normalizePath(c.Request.URL.Path)

		// Registra tamanho da resposta (aproximado pelo Content-Length)
		responseSize := float64(c.Writer.Size())
		if responseSize < 0 {
			responseSize = 0
		}

		// Registra métricas
		httpRequestsTotal.WithLabelValues(method, path, statusCode).Inc()
		httpRequestDurationSeconds.WithLabelValues(method, path, statusCode).Observe(duration)
		httpResponseSizeBytes.WithLabelValues(method, path, statusCode).Observe(responseSize)

		// Registra tamanho da requisição apenas para métodos que têm body
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
			httpRequestSizeBytes.WithLabelValues(method, path).Observe(requestSize)
		}
	}
}

// normalizePath normaliza o caminho da URL para reduzir cardinalidade das métricas
// Remove IDs e parâmetros dinâmicos, substituindo por :id
func normalizePath(path string) string {
	// Primeiro, tenta remover UUIDs (formato: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx)
	if len(path) > 36 {
		for i := 0; i <= len(path)-36; i++ {
			if i+36 <= len(path) {
				substr := path[i : i+36]
				if isValidUUID(substr) {
					return path[:i] + ":id" + normalizePath(path[i+36:])
				}
			}
		}
	}

	// Remove IDs numéricos simples (ex: /api/v1/users/123 -> /api/v1/users/:id)
	// Usa uma abordagem simples: se o último segmento após / é numérico, substitui
	parts := strings.Split(path, "/")
	for i := len(parts) - 1; i >= 0; i-- {
		if parts[i] != "" {
			if isNumeric(parts[i]) {
				parts[i] = ":id"
			}
			break
		}
	}

	return strings.Join(parts, "/")
}

// isNumeric verifica se uma string contém apenas dígitos
func isNumeric(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

// isValidUUID verifica se uma string é um UUID válido
func isValidUUID(s string) bool {
	if len(s) != 36 {
		return false
	}
	for i, c := range s {
		if i == 8 || i == 13 || i == 18 || i == 23 {
			if c != '-' {
				return false
			}
		} else {
			if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
				return false
			}
		}
	}
	return true
}
