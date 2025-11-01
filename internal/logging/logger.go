// @file: internal/logging/logger.go

package logging

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	// Logger é a instância global do logger zerolog
	Logger zerolog.Logger
)

func init() {
	// Configura o formato baseado na variável de ambiente
	logFormat := os.Getenv("LOG_FORMAT")
	logLevel := os.Getenv("LOG_LEVEL")

	// Configura o nível de log
	level := parseLogLevel(logLevel)
	zerolog.SetGlobalLevel(level)

	// Configura o formato de output
	if logFormat == "console" || logFormat == "" {
		// Formato console com cores (desenvolvimento)
		Logger = log.Output(zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: time.RFC3339,
		})
	} else {
		// Formato JSON (produção)
		Logger = log.With().
			Timestamp().
			Logger()
	}

	// Adiciona contexto global (service name, version, etc.)
	Logger = Logger.With().
		Str("service", "go-base-api").
		Str("version", "1.0.0").
		Logger()

	// Substitui o logger padrão do Go para usar zerolog
	log.Logger = Logger
}

// parseLogLevel converte string para zerolog.Level
func parseLogLevel(level string) zerolog.Level {
	switch level {
	case "trace":
		return zerolog.TraceLevel
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	case "panic":
		return zerolog.PanicLevel
	default:
		// Padrão: info em produção, debug em desenvolvimento
		if os.Getenv("GO_ENV") == "development" || os.Getenv("GO_ENV") == "dev" {
			return zerolog.DebugLevel
		}
		return zerolog.InfoLevel
	}
}

// Wrappers compatíveis com a API antiga (para transição)
var (
	InfoLogger  = &loggerWrapper{level: "info"}
	WarnLogger  = &loggerWrapper{level: "warn"}
	ErrorLogger = &loggerWrapper{level: "error"}
)

// loggerWrapper mantém compatibilidade com a API antiga
type loggerWrapper struct {
	level string
}

func (l *loggerWrapper) Printf(format string, v ...interface{}) {
	switch l.level {
	case "info":
		Logger.Info().Msgf(format, v...)
	case "warn":
		Logger.Warn().Msgf(format, v...)
	case "error":
		Logger.Error().Msgf(format, v...)
	}
}

// SetGlobalLogger permite configurar o logger globalmente
func SetGlobalLogger(logger zerolog.Logger) {
	Logger = logger
	log.Logger = logger
}

// WithContext retorna um logger com contexto adicional
func WithContext() zerolog.Context {
	return Logger.With()
}
