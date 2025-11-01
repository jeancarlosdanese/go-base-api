// @file: internal/middleware/logger_gin.go

package middleware

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jeancarlosdanese/go-base-api/internal/logging"
	"github.com/rs/zerolog"
)

// GinLoggerMiddleware retorna um middleware Gin que registra logs estruturados
func GinLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ignora endpoints de métricas e health check para evitar poluição de logs
		if c.Request.URL.Path == "/metrics" || c.Request.URL.Path == "/health" {
			c.Next()
			return
		}

		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Processa a requisição
		c.Next()

		// Calcula duração
		duration := time.Since(start)

		// Cria evento de log
		event := logging.Logger.Info()
		if c.Writer.Status() >= 500 {
			event = logging.Logger.Error()
		} else if c.Writer.Status() >= 400 {
			event = logging.Logger.Warn()
		}

		// Adiciona campos estruturados
		event.
			Int("status", c.Writer.Status()).
			Str("method", c.Request.Method).
			Str("path", path).
			Int("size", c.Writer.Size()).
			Dur("latency", duration).
			Str("client_ip", c.ClientIP()).
			Str("user_agent", c.Request.UserAgent())

		if raw != "" {
			event.Str("query", raw)
		}

		// Adiciona erro se houver
		if len(c.Errors) > 0 {
			errors := make([]error, len(c.Errors))
			for i, err := range c.Errors {
				errors[i] = err
			}
			event.Errs("errors", errors)
		}

		// Mensagem final
		event.Msg("HTTP Request")
	}
}

// GinRecoveryMiddleware retorna um middleware Gin que captura panics e loga
func GinRecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log do panic
				logging.Logger.Error().
					Interface("panic", err).
					Str("method", c.Request.Method).
					Str("path", c.Request.URL.Path).
					Str("client_ip", c.ClientIP()).
					Stack().
					Msg("Panic recovered")

				// Retorna erro 500
				c.JSON(500, gin.H{
					"error":   "Internal server error",
					"message": "An unexpected error occurred",
				})
				c.Abort()
			}
		}()

		c.Next()
	}
}

// SetupGinLogger configura o logger do Gin para usar zerolog
func SetupGinLogger() {
	// Desabilita o logger padrão do Gin
	gin.SetMode(gin.ReleaseMode)

	// Configura Gin para usar zerolog através do output
	if os.Getenv("GIN_MODE") != "release" {
		gin.DefaultWriter = &zerologWriter{level: zerolog.InfoLevel}
		gin.DefaultErrorWriter = &zerologWriter{level: zerolog.ErrorLevel}
	} else {
		// Em produção, não loga para stdout
		gin.DefaultWriter = os.Stdout
		gin.DefaultErrorWriter = os.Stderr
	}
}

// zerologWriter é um writer que redireciona output do Gin para zerolog
type zerologWriter struct {
	level zerolog.Level
}

func (w *zerologWriter) Write(p []byte) (n int, err error) {
	log := logging.Logger.WithLevel(w.level)
	log.Msg(string(p[:len(p)-1])) // Remove \n do final
	return len(p), nil
}
