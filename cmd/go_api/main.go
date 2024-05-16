// cmd/go_api/main.go

package main

import (
	"log"

	"github.com/jeancarlosdanese/go-base-api/internal/config" // Importa o pacote onde InitializeServicesContainer está definido
	"github.com/jeancarlosdanese/go-base-api/internal/domain/enums"
	"github.com/jeancarlosdanese/go-base-api/internal/routes" // Importa o pacote de rotas

	"github.com/gin-gonic/gin"
)

// @title 							Swagger Go Base API
// @version 						0.0.2
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

	// Inicializa o servidor
	if err := r.Run("0.0.0.0:5001"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
