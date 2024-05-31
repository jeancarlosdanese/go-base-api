// internal/services/token_service.go

package services

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenServiceInterface interface {
	CreateTokens(userID uuid.UUID, roles []string, permissions []string) (string, string, error)
	RefreshTokens(refreshToken string) (uuid.UUID, error)
	ValidateToken(tokenString string) (*jwt.Token, error)
	GetAccessDuration() time.Duration
}
type TokenService struct {
	SecretKey       []byte
	AccessDuration  time.Duration
	RefreshDuration time.Duration
}

func NewTokenService(secretKey string, accessDuration, refreshDuration time.Duration) *TokenService {
	return &TokenService{
		SecretKey:       []byte(secretKey),
		AccessDuration:  accessDuration,
		RefreshDuration: refreshDuration,
	}
}

// CreateTokens cria tanto o token de acesso quanto o refresh token
func (t *TokenService) CreateTokens(userID uuid.UUID, roles []string, permissions []string) (string, string, error) {
	accessToken, err := t.createAccessToken(userID, roles, permissions)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := t.createRefreshToken(userID)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func (t *TokenService) RefreshTokens(refreshToken string) (uuid.UUID, error) {
	// Validar e parsear o refreshToken
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("algoritmo de assinatura inesperado: %v", token.Header["alg"])
		}
		return t.SecretKey, nil
	})

	if err != nil || !token.Valid {
		return uuid.Nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return uuid.Nil, err
	}

	userID, err := uuid.Parse(claims["sub"].(string))
	if err != nil {
		return uuid.Nil, err
	}

	return userID, nil
}

func (t *TokenService) createAccessToken(userID uuid.UUID, roles, permissions []string) (string, error) {
	claims := jwt.MapClaims{
		"sub":         userID.String(),
		"roles":       roles,
		"permissions": permissions,
		"exp":         time.Now().Add(t.AccessDuration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(t.SecretKey)
}

func (t *TokenService) createRefreshToken(userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID.String(),
		"exp": time.Now().Add(t.RefreshDuration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(t.SecretKey)
}

func (t *TokenService) ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("algoritmo de assinatura inesperado: %v", token.Header["alg"])
		}

		return t.SecretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("token inválido")
	}

	return token, nil
}

func (t *TokenService) GetAccessDuration() time.Duration {
	return t.AccessDuration
}
