// @file: cmd/go_api/main.go

package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jeancarlosdanese/go-base-api/internal/config" // Importa o pacote onde InitializeServicesContainer está definido
	"github.com/jeancarlosdanese/go-base-api/internal/domain/enums"
	"github.com/jeancarlosdanese/go-base-api/internal/logging"
	"github.com/jeancarlosdanese/go-base-api/internal/middleware"
	"github.com/jeancarlosdanese/go-base-api/internal/routes" // Importa o pacote de rotas

	"github.com/gin-gonic/gin"
	"github.com/jeancarlosdanese/go-base-api/docs"
)

// @title 							Go Base API
// @version 						1.0.0
// @description 					API REST base com autenticação JWT, RBAC (Casbin), multi-tenancy, rate limiting (Redis) e métricas Prometheus. Template completo para desenvolvimento de APIs Go com arquitetura limpa.
// @termsOfService 					https://github.com/jeancarlosdanese/go-base-api/blob/main/LICENSE
// @contact.name 					Go Base API Support
// @contact.url 					https://github.com/jeancarlosdanese/go-base-api/go-api
// @license.name 					MIT
// @license.url 					https://github.com/jeancarlosdanese/go-base-api/blob/main/LICENSE
// @host 							http://localhost:5001
// @SecurityDefinitions.apiKey 		Bearer
// @in header
// @name 							Authorization
// @BasePath 						/api/v1/
func main() {
	enums.Initialize() // Garante que tudo está configurado antes de usar.

	// Configura o logger do Gin para usar zerolog
	middleware.SetupGinLogger()

	// Configura versão do Swagger em runtime
	version := os.Getenv("SERVICE_VERSION")
	if version == "" {
		version = "1.0.0"
	}
	docs.SwaggerInfo.Version = version

	r := gin.Default()

	// Inicializa o container de serviços usando o Google Wire.
	// Essa chamada irá configurar todos os serviços necessários, incluindo o pool de conexões do banco de dados.
	sc, err := config.InitializeServicesContainer() // Chama a nova função de inicialização
	if err != nil {
		logging.Logger.Fatal().
			Err(err).
			Str("operation", "init_services_container").
			Msg("Falha ao inicializar container de serviços")
	}

	// Configura as rotas com o container de serviços
	routes.SetupRouter(r, sc)

	// Obtém porta do ambiente ou usa padrão
	port := os.Getenv("PORT")
	if port == "" {
		port = "5001"
	}

	// Cria o servidor HTTP
	server := &http.Server{
		Addr:    "0.0.0.0:" + port,
		Handler: r,
	}

	// Canal para escutar sinais do sistema operacional
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	// Goroutine para iniciar o servidor
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logging.Logger.Fatal().
				Err(err).
				Str("port", port).
				Str("operation", "start_server").
				Msg("Erro ao iniciar o servidor")
		}
	}()

	logging.Logger.Info().
		Str("port", port).
		Str("operation", "start_server").
		Msg("Servidor iniciado")

	// Esperar por um sinal de término
	<-stop
	logging.Logger.Info().
		Str("operation", "shutdown_server").
		Msg("Desligando o servidor")

	// Contexto com timeout para o desligamento gracioso
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Iniciar o desligamento gracioso
	if err := server.Shutdown(ctx); err != nil {
		logging.Logger.Fatal().
			Err(err).
			Str("operation", "shutdown_server").
			Msg("Erro ao desligar o servidor")
	}

	logging.Logger.Info().
		Str("operation", "shutdown_server").
		Msg("Servidor desligado")
}
