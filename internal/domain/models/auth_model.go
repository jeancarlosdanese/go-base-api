// internal/domain/models/login_model.go

package models

// LoginForm representa os dados de entrada para o login do usuário.
type LoginForm struct {
	Email    string `form:"email" binding:"required,email"`
	Password string `form:"password" binding:"required"`
}

type Token struct {
	Type         string    `json:"type"`
	Token        string    `json:"token"`
	RefreshToken *string   `json:"refreshToken"`
	User         TokenUser `json:"user"`
	Roles        []string  `json:"roles"`
	Policies     []string  `json:"policies"`
}

type TokenUser struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Username  string  `json:"username"`
	Email     string  `json:"email"`
	Thumbnail *string `json:"thumbnail"`
}
