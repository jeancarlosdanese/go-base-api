// internal/routes/router.go

package routes

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	_ "github.com/jeancarlosdanese/go-base-api/docs"

	"net"
	"sync"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"

	"github.com/jeancarlosdanese/go-base-api/internal/app"
	contextkeys "github.com/jeancarlosdanese/go-base-api/internal/domain/context_keys"
	"github.com/jeancarlosdanese/go-base-api/internal/domain/models"
	handlers_v1 "github.com/jeancarlosdanese/go-base-api/internal/handlers_v1"
	"github.com/jeancarlosdanese/go-base-api/internal/services"
)

// SetupRouter agora aceita ServicesContainer como argumento.
func SetupRouter(r *gin.Engine, sc *app.ServicesContainer) {
	// Middleware global de segurança
	r.Use(SecurityHeadersMiddleware())

	// Timeout e limite de tamanho do request body
	r.Use(RequestSizeLimitMiddleware(10 << 20)) // 10 MB

	// Compressão gzip para respostas
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	// Middleware de métricas
	r.Use(MetricsMiddleware())

	// Rate limiting simples - 50 requests por minuto por IP
	r.Use(RateLimitMiddleware(50, time.Minute))

	// Servir arquivos estáticos (favicon, etc.)
	r.Static("/static", "./static")

	// Favicon específico com cache usando StaticFile
	r.StaticFile("/favicon.ico", "./static/favicon.ico")

	// Servir arquivos de documentação (Postman, Insomnia, etc.)
	r.Static("/docs", "./docs")

	// Health check endpoint (não limitado por rate limiting)
	r.GET("/health", HealthCheckHandler(sc.DB))

	// Metrics endpoint
	r.GET("/metrics", MetricsHandler)

	// Setup da rota do Swagger com redirecionamento automático
	swaggerHandler := ginSwagger.WrapHandler(swaggerFiles.Handler)
	r.GET("/swagger/*any", func(c *gin.Context) {
		path := c.Param("any")
		if path == "/" {
			c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
			return
		}
		swaggerHandler(c)
	})

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

		// Se não houver Origin (ex: curl, Postman sem Origin definido),
		// usar localhost como padrão para desenvolvimento
		if origin == "" {
			origin = "localhost"
		}

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

		// Se não houver Origin (ex: curl, Postman sem Origin definido),
		// usar localhost como padrão para desenvolvimento
		if origin == "" {
			origin = "localhost"
		}

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

// SecurityHeadersMiddleware adiciona headers de segurança HTTP a todas as respostas
func SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Previne MIME type sniffing
		c.Header("X-Content-Type-Options", "nosniff")

		// Previne clickjacking
		c.Header("X-Frame-Options", "DENY")

		// Ativa proteção XSS
		c.Header("X-XSS-Protection", "1; mode=block")

		// HTTP Strict Transport Security (HSTS) - apenas para HTTPS
		if c.Request.TLS != nil {
			c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}

		// Referrer Policy
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		// Content Security Policy básica
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data:")

		// Permissions Policy (anteriormente Feature Policy)
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		c.Next()
	}
}

// HealthCheckHandler verifica o status de saúde da aplicação e suas dependências
func HealthCheckHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		health := gin.H{
			"status":    "ok",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
			"services": gin.H{
				"database": checkDatabaseHealth(db),
				"redis":    checkRedisHealth(),
			},
		}

		// Verifica se todos os serviços estão saudáveis
		allHealthy := true
		for _, service := range health["services"].(gin.H) {
			if serviceStatus, ok := service.(gin.H); ok {
				if status, exists := serviceStatus["status"]; exists && status != "ok" {
					allHealthy = false
					break
				}
			}
		}

		statusCode := http.StatusOK
		if !allHealthy {
			statusCode = http.StatusServiceUnavailable
			health["status"] = "degraded"
		}

		c.JSON(statusCode, health)
	}
}

// checkDatabaseHealth verifica a saúde do banco de dados
func checkDatabaseHealth(db *gorm.DB) gin.H {
	if db == nil {
		return gin.H{
			"status": "error",
			"error":  "Database connection is nil",
		}
	}

	// Testa a conexão executando uma query simples
	sqlDB, err := db.DB()
	if err != nil {
		return gin.H{
			"status": "error",
			"error":  err.Error(),
		}
	}

	// Testa ping no banco
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return gin.H{
			"status": "error",
			"error":  err.Error(),
		}
	}

	// Verifica estatísticas básicas
	stats := sqlDB.Stats()
	return gin.H{
		"status":              "ok",
		"open_connections":    stats.OpenConnections,
		"in_use":              stats.InUse,
		"idle":                stats.Idle,
		"wait_count":          stats.WaitCount,
		"wait_duration":       stats.WaitDuration.String(),
		"max_idle_closed":     stats.MaxIdleClosed,
		"max_lifetime_closed": stats.MaxLifetimeClosed,
	}
}

// checkRedisHealth verifica a saúde do Redis
func checkRedisHealth() gin.H {
	// Para uma implementação completa, seria necessário injetar o cliente Redis
	// Por enquanto, retornamos um status básico
	return gin.H{
		"status": "ok",
		"note":   "Redis health check requires Redis client injection",
	}
}

// Rate limiting store usando map com mutex para thread safety
var (
	rateLimitStore = make(map[string][]time.Time)
	rateLimitMutex sync.RWMutex
)

// Métricas básicas
var (
	metricsData = struct {
		requestsTotal    int64
		requestsByPath   map[string]int64
		requestsByMethod map[string]int64
		requestsByStatus map[int]int64
		responseTime     []time.Duration
		mutex            sync.RWMutex
	}{
		requestsByPath:   make(map[string]int64),
		requestsByMethod: make(map[string]int64),
		requestsByStatus: make(map[int]int64),
		responseTime:     make([]time.Duration, 0),
	}
)

// RateLimitMiddleware implementa controle de taxa de requisições por IP
func RateLimitMiddleware(requests int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Pula rate limiting para health check
		if c.Request.URL.Path == "/health" {
			c.Next()
			return
		}

		// Obtém IP do cliente
		ip := getClientIP(c)

		rateLimitMutex.Lock()

		// Limpa timestamps antigos
		now := time.Now()
		windowStart := now.Add(-window)

		if timestamps, exists := rateLimitStore[ip]; exists {
			// Remove timestamps fora da janela
			validTimestamps := make([]time.Time, 0)
			for _, timestamp := range timestamps {
				if timestamp.After(windowStart) {
					validTimestamps = append(validTimestamps, timestamp)
				}
			}
			rateLimitStore[ip] = validTimestamps
		}

		// Verifica se excedeu o limite
		if len(rateLimitStore[ip]) >= requests {
			rateLimitMutex.Unlock()
			c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", requests))
			c.Header("X-RateLimit-Remaining", "0")
			c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", now.Add(window).Unix()))
			c.Header("Retry-After", fmt.Sprintf("%d", int(window.Seconds())))
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "Rate limit exceeded",
				"message":     "Too many requests. Please try again later.",
				"retry_after": int(window.Seconds()),
			})
			c.Abort()
			return
		}

		// Adiciona timestamp atual
		rateLimitStore[ip] = append(rateLimitStore[ip], now)

		rateLimitMutex.Unlock()

		// Headers informativos
		remaining := requests - len(rateLimitStore[ip])
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", requests))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", now.Add(window).Unix()))

		c.Next()
	}
}

// getClientIP obtém o IP real do cliente considerando proxies
func getClientIP(c *gin.Context) string {
	// Tenta headers de proxy primeiro
	if ip := c.GetHeader("X-Real-IP"); ip != "" {
		return ip
	}
	if ip := c.GetHeader("X-Forwarded-For"); ip != "" {
		// Pega o primeiro IP da lista
		if idx := strings.Index(ip, ","); idx > 0 {
			return strings.TrimSpace(ip[:idx])
		}
		return ip
	}

	// Fallback para RemoteAddr
	ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
	if err != nil {
		return c.Request.RemoteAddr
	}
	return ip
}

// MetricsMiddleware coleta métricas básicas das requisições
func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		// Coleta métricas após a requisição
		duration := time.Since(start)
		statusCode := c.Writer.Status()
		path := c.Request.URL.Path
		method := c.Request.Method

		metricsData.mutex.Lock()
		metricsData.requestsTotal++
		metricsData.requestsByPath[path]++
		metricsData.requestsByMethod[method]++
		metricsData.requestsByStatus[statusCode]++

		// Mantém apenas as últimas 1000 medições de tempo de resposta
		metricsData.responseTime = append(metricsData.responseTime, duration)
		if len(metricsData.responseTime) > 1000 {
			metricsData.responseTime = metricsData.responseTime[1:]
		}
		metricsData.mutex.Unlock()
	}
}

// MetricsHandler retorna métricas em formato JSON
func MetricsHandler(c *gin.Context) {
	metricsData.mutex.RLock()
	defer metricsData.mutex.RUnlock()

	// Calcula tempo de resposta médio
	var avgResponseTime time.Duration
	if len(metricsData.responseTime) > 0 {
		var total time.Duration
		for _, duration := range metricsData.responseTime {
			total += duration
		}
		avgResponseTime = total / time.Duration(len(metricsData.responseTime))
	}

	// Calcula distribuição de status codes
	statusDistribution := make(map[string]int64)
	for status, count := range metricsData.requestsByStatus {
		var category string
		switch {
		case status >= 200 && status < 300:
			category = "2xx"
		case status >= 300 && status < 400:
			category = "3xx"
		case status >= 400 && status < 500:
			category = "4xx"
		case status >= 500:
			category = "5xx"
		default:
			category = "other"
		}
		statusDistribution[category] += count
	}

	metrics := gin.H{
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"uptime":    "N/A", // Seria necessário armazenar o tempo de inicialização
		"requests": gin.H{
			"total":               metricsData.requestsTotal,
			"by_method":           metricsData.requestsByMethod,
			"by_path":             metricsData.requestsByPath,
			"by_status":           metricsData.requestsByStatus,
			"status_distribution": statusDistribution,
		},
		"performance": gin.H{
			"avg_response_time":    avgResponseTime.String(),
			"avg_response_time_ms": avgResponseTime.Milliseconds(),
			"samples_count":        len(metricsData.responseTime),
		},
		"system": gin.H{
			"goroutines": "N/A", // Seria necessário usar runtime.NumGoroutine()
			"memory":     "N/A", // Seria necessário usar runtime.MemStats
		},
	}

	c.JSON(http.StatusOK, metrics)
}

// RequestSizeLimitMiddleware limita o tamanho do corpo da requisição
func RequestSizeLimitMiddleware(maxSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Para métodos que têm body
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
			if c.Request.ContentLength > maxSize {
				c.JSON(http.StatusRequestEntityTooLarge, gin.H{
					"error":    "Request too large",
					"message":  "Request body exceeds maximum allowed size",
					"max_size": fmt.Sprintf("%d bytes", maxSize),
				})
				c.Abort()
				return
			}

			// Para multipart forms, verifica também
			if c.Request.Header.Get("Content-Type") != "" &&
				strings.Contains(c.Request.Header.Get("Content-Type"), "multipart/form-data") {
				c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxSize)
			}
		}

		c.Next()
	}
}
