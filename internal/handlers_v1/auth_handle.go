// internal/handlers_v1/auth_handle.go

package handlers_v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jeancarlosdanese/go-base-api/internal/domain/models"
	"github.com/jeancarlosdanese/go-base-api/internal/logging"
	"github.com/jeancarlosdanese/go-base-api/internal/services"
	"github.com/jeancarlosdanese/go-base-api/internal/utils"
)

// AuthHandler struct para segurar os serviços necessários.
type AuthHandler struct {
	userService       services.UserServiceInterface
	tokenService      services.TokenServiceInterface
	tokenRedisService services.TokenRedisServiceInterface
}

// NewAuthHandler cria uma nova instância de AuthHandler.
func NewAuthHandler(
	userService services.UserServiceInterface,
	tokenService services.TokenServiceInterface,
	tokenRedisService services.TokenRedisServiceInterface) *AuthHandler {
	return &AuthHandler{
		userService:       userService,
		tokenService:      tokenService,
		tokenRedisService: tokenRedisService,
	}
}

// RegisterRoutes registra as rotas para autenticação.
func (h *AuthHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/login", h.Login)
	router.POST("/refresh", h.Refresh)
}

// login realiza o login do usuário e retorna um JWT.
// @Summary Loga um usuário
// @Description Loga um usuário usando email e senha
// @Tags Auth
// @Accept x-www-form-urlencoded
// @Produce json
// @Param email formData string true "Email do Usuário"
// @Param password formData string true "Senha do Usuário"
// @Success 200 {object} map[string]interface{} "Token gerado com sucesso"
// @Failure 400 {object} map[string]string "Erro de autenticação"
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	origin := c.GetString("Origin")
	if origin == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Origem não fornecida"})
		return
	}

	var loginForm models.LoginForm
	if err := c.ShouldBind(&loginForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parâmetros de entrada inválidos"})
		return
	}

	user, err := h.userService.Authenticate(c, loginForm.Email, loginForm.Password, origin)
	if err != nil {
		utils.HandleAuthenticationError(c, err)
		return
	}

	h.generateAndSaveTokens(c, user)
}

// Refresh renova o token usando o refreshToken.
// @Summary Renova o token
// @Description Renova o token usando o refreshToken
// @Tags Auth
// @Accept x-www-form-urlencoded
// @Produce json
// @Param refreshToken formData string true "Refresh Token"
// @Success 200 {object} map[string]interface{} "Token renovado com sucesso"
// @Failure 400 {object} map[string]string "Erro de autenticação"
// @Router /api/v1/auth/refresh [post]
// Refresh renova o token usando o refreshToken.
func (h *AuthHandler) Refresh(c *gin.Context) {
	var refreshTokenRequest models.RefreshTokenRequest
	if err := c.ShouldBind(&refreshTokenRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parâmetros de entrada inválidos"})
		return
	}

	userID, err := h.tokenService.RefreshTokens(refreshTokenRequest.RefreshToken)
	if err != nil {
		logging.WarnLogger.Printf("Erro ao renovar token: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido ou expirado"})
		return
	}

	user, err := h.userService.GetOnlyByID(c, userID)
	if err != nil {
		logging.WarnLogger.Printf("Erro ao buscar usuário: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não encontrado"})
		return
	}

	h.generateAndSaveTokens(c, user)
}

// generateAndSaveTokens gera e salva tokens para o usuário.
func (h *AuthHandler) generateAndSaveTokens(c *gin.Context, user *models.User) {
	roles := user.ExtractRoles()
	policies := user.ExtractPolicies()

	accessToken, refreshToken, err := h.tokenService.CreateTokens(user.ID, roles, policies)
	if err != nil {
		logging.ErrorLogger.Printf("Falha ao gerar tokens para o usuário: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao gerar tokens"})
		return
	}

	if err := h.tokenRedisService.SaveTokenDataRedis(user, accessToken, refreshToken, h.tokenService.GetAccessDuration()); err != nil {
		logging.ErrorLogger.Printf("Falha ao salvar informações do usuário no Redis: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao salvar informações do usuário no Redis"})
		return
	}

	response := prepareTokenResponse(user, accessToken, refreshToken)
	c.JSON(http.StatusOK, response)
}

// prepareTokenResponse monta a resposta completa do token.
func prepareTokenResponse(user *models.User, token, refreshToken string) models.Token {
	return models.Token{
		Type:         "bearer",
		Token:        token,
		RefreshToken: &refreshToken,
		User: models.TokenUser{
			ID:        user.ID.String(),
			Name:      user.Name,
			Username:  user.Username,
			Email:     user.Email,
			Thumbnail: user.Thumbnail,
		},
		Roles:    user.ExtractRoles(),
		Policies: user.ExtractPolicies(),
	}
}
