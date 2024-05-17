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

	// Configuração de rotas não autenticadas
	auth := v1.Group("/auth")
	{
		authHandler := handlers_v1.NewAuthHandler(sc.UserService, sc.TokenService, sc.TokenRedisService)
		auth.POST("/login", authHandler.Login) // Registra diretamente a rota POST /login no grupo /auth
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
		c.Set("origin", origin) // Armazenar no contexto

		c.Next() // continuar com a cadeia de middlewares/handlers
	}
}

func AuthMiddleware(tokenService *services.TokenService, tokenRedisService *services.TokenRedisService) gin.HandlerFunc {
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
		tokenDataRedis, err := tokenRedisService.GetTokenDataRedisFromToken(tokenString)
		if err != nil || tokenDataRedis == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Falha ao recuperar informações do usuário"})
			c.Abort()
			return
		}

		// Configura o contexto com o usuário para uso posterior
		c.Set(string(contextkeys.TokenDataKey), tokenDataRedis)
		c.Set(string(contextkeys.TenantIDKey), tokenDataRedis.User.TenantID)
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
		tokenData, exists := c.Get(string(contextkeys.TokenDataKey))
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
			c.Abort()
			return
		}

		// Asumindo que a estrutura User tem um campo Roles que é um slice de strings.
		tokenDataRedis := tokenData.(*models.TokenDataRedis)
		isValidRole := false
		for _, role := range tokenDataRedis.User.Roles {

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
func PolicyMiddleware(casbinService *services.CasbinService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenData, exists := c.Get(string(contextkeys.TokenDataKey))
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
			c.Abort()
			return
		}

		tokenDataRedis := tokenData.(*models.TokenDataRedis) // Certifique-se de que este cast está correto conforme sua implementação
		obj := c.Request.URL.Path
		act := c.Request.Method

		// Tenta verificar permissões usando ID do usuário e roles
		if !checkPermissions(tokenDataRedis, casbinService, obj, act) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Acesso negado - permissão insuficiente"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// checkPermissions tenta verificar as permissões com o ID do usuário e suas roles
func checkPermissions(tokenDataRedis *models.TokenDataRedis, casbinService *services.CasbinService, obj, act string) bool {
	// Verifica permissões especiais usando o ID do usuário
	if casbinService.CheckPermission(tokenDataRedis.User.ID, obj, act) {
		return true
	}

	// Verifica permissões baseadas nas roles
	for _, role := range tokenDataRedis.User.Roles {
		if casbinService.CheckPermission(role, obj, act) {
			return true
		}
	}

	// Nenhuma permissão encontrada
	return false
}
