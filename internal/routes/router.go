// internal/routes/router.go

package routes

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	_ "github.com/jeancarlosdanese/go-base-api/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/jeancarlosdanese/go-base-api/internal/app"
	contextkeys "github.com/jeancarlosdanese/go-base-api/internal/domain/context_keys"
	"github.com/jeancarlosdanese/go-base-api/internal/domain/models"
	handlers_v1 "github.com/jeancarlosdanese/go-base-api/internal/handlers_v1"
	"github.com/jeancarlosdanese/go-base-api/internal/services"
)

// SetupRouter agora aceita ServicesContainer como argumento.
func SetupRouter(r *gin.Engine, sc *app.ServicesContainer) {
	// Setup da rota do Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Definindo o grupo de rotas para a versão 1 da API
	v1 := r.Group("/api/v1")

	// Teste de rota com autenticação vi X-API-Key
	authApiKeyGroup := v1.Group("/auth-apikey")
	authApiKeyGroup.Use(XApiKeyMiddleware(sc.ApiKeyRedisService))
	{
		authApiKeyHandler := handlers_v1.NewAuthApiKeyHandler()
		authApiKeyHandler.RegisterRoutes(authApiKeyGroup)
	}

	// Configuração de rotas não autenticadas
	authGroup := v1.Group("/auth")
	authGroup.Use(OriginMiddleware())
	{
		authHandler := handlers_v1.NewAuthHandler(sc.UserService, sc.TokenService, sc.TokenRedisService)
		// auth.POST("/login", authHandler.Login) // Registra diretamente a rota POST /login no grupo /auth
		authHandler.RegisterRoutes(authGroup)
	}

	// Middleware de autenticação que é aplicado a todas as rotas que necessitam autenticação
	secured := v1.Group("/")
	secured.Use(AuthMiddleware(sc.TokenService, sc.TokenRedisService))
	{
		// Grupo para gestão de tenants
		tenantsGroup := secured.Group("/tenants")
		// tenantsGroup.Use(RoleMiddleware("administration")) // Apenas usuários com role "administration"
		tenantsGroup.Use(PolicyMiddleware(sc.CasbinService))
		{
			tenantsHandler := handlers_v1.NewTenantsHandler(sc.TenantService)
			tenantsHandler.RegisterRoutes(tenantsGroup)
			// tenantsGroup.GET("", tenantsHandler.GetAll)
			// tenantsGroup.POST("", tenantsHandler.Create)
		}

		{
			usersHandler := handlers_v1.NewUsersHandler(sc.UserService)
			usersGroup := secured.Group("/users")
			usersGroup.Use(PolicyMiddleware(sc.CasbinService))
			// Aqui você pode adicionar middlewares específicos para /users se necessário
			usersHandler.RegisterRoutes(usersGroup)
		}
	}
}

func extractToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	// Divide o cabeçalho em partes e verifica se é um token tipo Bearer
	parts := strings.Split(authHeader, " ")
	if len(parts) == 2 && parts[0] == "Bearer" {
		return parts[1]
	}
	return ""
}

func OriginMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		// Remover o protocolo (http:// ou https://)
		origin = strings.TrimPrefix(origin, "http://")
		origin = strings.TrimPrefix(origin, "https://")

		c.Set("Origin", origin) // Armazenar no contexto

		c.Next() // continuar com a cadeia de middlewares/handlers
	}
}

// XApiKeyMiddleware é um Middleware para verificar a API Key
func XApiKeyMiddleware(apiKeyRedisService services.ApiKeyRedisServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		// Remover o protocolo (http:// ou https://)
		origin = strings.TrimPrefix(origin, "http://")
		origin = strings.TrimPrefix(origin, "https://")

		c.Set("Origin", origin) // Armazenar no contexto

		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "API Key inválida"})
			return
		}

		// Tenta recuperar as informações do Tenant do Redis
		tenantRedis, err := apiKeyRedisService.GetTenantRedisFromApiKey(apiKey, origin)
		if err != nil || tenantRedis == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Falha ao recuperar informações do Tenant"})
			c.Abort()
			return
		}

		// Configura o contexto com o usuário para uso posterior
		c.Set(string(contextkeys.TenantDataKey), tenantRedis)
		c.Set(string(contextkeys.TenantIDKey), tenantRedis.ID)

		c.Next() // continuar com a cadeia de middlewares/handlers
	}
}

// AuthMiddleware é um Midleware para verificar o Bearer Token
func AuthMiddleware(tokenService services.TokenServiceInterface, tokenRedisService services.TokenRedisServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Capturar o tempo de início
		start := time.Now()
		tokenString := extractToken(c) // Função auxiliar para extrair o token
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token não fornecido"})
			c.Abort()
			return
		}

		_, err := tokenService.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido: " + err.Error()})
			c.Abort()
			return
		}

		// Tenta recuperar as informações do usuário do Redis
		userRedis, err := tokenRedisService.GetUserRedisFromToken(tokenString)
		if err != nil || userRedis == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Falha ao recuperar informações do usuário"})
			c.Abort()
			return
		}

		// Configura o contexto com o usuário para uso posterior
		c.Set(string(contextkeys.UserDataKey), userRedis)
		c.Set(string(contextkeys.TenantIDKey), userRedis.TenantID)
		c.Next() // Prosseguir com a próxima função no pipeline
		// Capturar o tempo de término
		end := time.Now()

		// Calcular a diferença
		duration := end.Sub(start)

		fmt.Printf("O processo demorou %v\n", duration)
	}
}

// RoleMiddleware verifica se o usuário possui as roles necessárias.
func RoleMiddleware(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenData, exists := c.Get(string(contextkeys.UserDataKey))
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
			c.Abort()
			return
		}

		// Asumindo que a estrutura User tem um campo Roles que é um slice de strings.
		userRedis := tokenData.(*models.UserRedis)
		isValidRole := false
		for _, role := range userRedis.Roles {

			for _, requiredRole := range requiredRoles {
				if role == requiredRole {
					isValidRole = true
					break
				}
			}
			if isValidRole {
				break
			}
		}

		if !isValidRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "Acesso negado - role inválida"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// PolicyMiddleware verifica permissões específicas para ações ou recursos usando Casbin.
func PolicyMiddleware(casbinService services.CasbinServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenData, exists := c.Get(string(contextkeys.UserDataKey))
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
			c.Abort()
			return
		}

		userRedis := tokenData.(*models.UserRedis) // Certifique-se de que este cast está correto conforme sua implementação
		obj := c.Request.URL.Path
		act := c.Request.Method

		// Tenta verificar permissões usando ID do usuário e roles
		if !checkPermissions(userRedis, casbinService, obj, act) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Acesso negado - permissão insuficiente"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// checkPermissions tries to verify permissions using the user's ID and their roles
func checkPermissions(userRedis *models.UserRedis, casbinService services.CasbinServiceInterface, obj, act string) bool {
	userID := fmt.Sprintf("%v", userRedis.ID) // Ensure user ID is converted to string

	// Verify special permissions using the user ID
	if casbinService.CheckPermission(userID, obj, act) {
		return true
	}

	// Check permissions based on the roles
	for _, role := range userRedis.Roles {
		if casbinService.CheckPermission(role, obj, act) {
			return true
		}
	}

	// No permission found
	return false
}
