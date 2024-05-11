// internal/handlers_v1/auth_handle.go

package handlers_v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jeancarlosdanese/go-base-api/internal/domain/models"
	"github.com/jeancarlosdanese/go-base-api/internal/services"
)

// AuthHandler struct para segurar os serviços necessários.
type AuthHandler struct {
	userService       *services.UserService
	tokenService      *services.TokenService
	tokenRedisService *services.TokenRedisService
}

// NewAuthHandler cria uma nova instância de AuthHandler.
func NewAuthHandler(
	userService *services.UserService,
	tokenService *services.TokenService,
	tokenRedisService *services.TokenRedisService) *AuthHandler {
	return &AuthHandler{
		userService:       userService,
		tokenService:      tokenService,
		tokenRedisService: tokenRedisService,
	}
}

// RegisterRoutes registra as rotas para autenticação.
func (h *AuthHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/login", h.Login)
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
	var loginForm models.LoginForm
	if err := c.ShouldBind(&loginForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parâmetros de entrada inválidos"})
		return
	}

	user, err := h.userService.Authenticate(c.Request.Context(), loginForm.Email, loginForm.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciais inválidas"})
		return
	}

	// Extrai roles e permissions especiais usando métodos de receptor
	roles := user.ExtractRoles()
	permissions := user.ExtractPermissions()

	accessToken, refreshToken, err := h.tokenService.CreateTokens(user.ID, roles, permissions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao gerar tokens"})
		return
	}

	// Salva as informações do usuário no Redis
	if err := h.tokenRedisService.SaveUserTokenInfo(user, accessToken, refreshToken); err != nil {
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
			Thumbnail: &user.Thumbnail,
		},
		Roles:       user.ExtractRoles(),
		Permissions: user.ExtractPermissions(),
	}
}
