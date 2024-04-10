// internal/routes/router.go

package routes

import (
	"github.com/gin-gonic/gin"
	"hyberica.io/go/go-api/internal/app" // Ajuste conforme a localização do ServicesContainer
	handlers_v1 "hyberica.io/go/go-api/internal/handlers_v1"
)

// SetupRouter agora aceita ServicesContainer como argumento.
func SetupRouter(r *gin.Engine, sc *app.ServicesContainer) {
	// Definindo o grupo de rotas para a versão 1 da API
	v1 := r.Group("/api/v1")

	// Configurando as rotas de tenants usando o handler e o serviço de tenants
	{
		tenantsHandler := handlers_v1.NewTenantsHandler(sc.TenantService)
		tenantsGroup := v1.Group("/tenants")
		tenantsHandler.RegisterRoutes(tenantsGroup)
	}

	// Adicione mais configurações de rotas para outros handlers conforme necessário
}
