// cmd/go_api/main.go

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jeancarlosdanese/go-base-api/internal/config" // Importa o pacote onde InitializeServicesContainer está definido
	"github.com/jeancarlosdanese/go-base-api/internal/domain/enums"
	"github.com/jeancarlosdanese/go-base-api/internal/routes" // Importa o pacote de rotas

	"github.com/gin-gonic/gin"
)

// @title 							Swagger Go Base API
// @version 						0.0.3
// @description 					This is a Go Base API.
// @termsOfService 					github.com/jeancarlosdanese/go-base-api/blob/main/LICENSE
// @contact.name 					Go Base API Support
// @contact.url 					github.com/jeancarlosdanese/go-base-api
// @license.name 					MIT
// @license.url 					github.com/jeancarlosdanese/go-base-api/blob/main/LICENSE
// @host 							localhost:5001
// @SecurityDefinitions.apiKey 		Bearer
// @in header
// @name 							Authorization
// @BasePath 						/api/v1/
func main() {
	enums.Initialize() // Garante que tudo está configurado antes de usar.

	r := gin.Default()

	// Inicializa o container de serviços usando o Google Wire.
	// Essa chamada irá configurar todos os serviços necessários, incluindo o pool de conexões do banco de dados.
	sc, err := config.InitializeServicesContainer() // Chama a nova função de inicialização
	if err != nil {
		log.Fatalf("Failed to initialize services container: %v", err)
	}

	// Configura as rotas com o container de serviços
	routes.SetupRouter(r, sc)

	// Cria o servidor HTTP
	server := &http.Server{
		Addr:    "0.0.0.0:5001",
		Handler: r,
	}

	// Canal para escutar sinais do sistema operacional
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	// Goroutine para iniciar o servidor
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Erro ao iniciar o servidor: %v", err)
		}
	}()

	log.Println("Servidor iniciado...")

	// Esperar por um sinal de término
	<-stop
	log.Println("Desligando o servidor...")

	// Contexto com timeout para o desligamento gracioso
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Iniciar o desligamento gracioso
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Erro ao desligar o servidor: %v", err)
	}

	log.Println("Servidor desligado.")
}
