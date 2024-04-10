// cmd/go_api/main.go

package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"hyberica.io/go/go-api/internal/config" // Importa o pacote onde InitializeServicesContainer está definido
	"hyberica.io/go/go-api/internal/routes" // Importa o pacote de rotas
)

func main() {
	r := gin.Default()

	// Inicializa o container de serviços usando o Google Wire.
	// Essa chamada irá configurar todos os serviços necessários, incluindo o pool de conexões do banco de dados.
	sc, err := config.InitializeServicesContainer() // Chama a nova função de inicialização
	if err != nil {
		log.Fatalf("Failed to initialize services container: %v", err)
	}

	// Configura as rotas com o container de serviços
	routes.SetupRouter(r, sc)

	// Inicializa o servidor
	if err := r.Run(); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
