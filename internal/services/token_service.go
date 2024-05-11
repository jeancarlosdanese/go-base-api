// internal/services/token_service.go

package services

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenService struct {
	secretKey       []byte
	accessDuration  time.Duration
	refreshDuration time.Duration
}

func NewTokenService(secretKey string, accessDuration, refreshDuration time.Duration) *TokenService {
	return &TokenService{
		secretKey:       []byte(secretKey),
		accessDuration:  accessDuration,
		refreshDuration: refreshDuration,
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

func (t *TokenService) createAccessToken(userID uuid.UUID, roles, permissions []string) (string, error) {
	claims := jwt.MapClaims{
		"sub":         userID.String(),
		"roles":       roles,
		"permissions": permissions,
		"exp":         time.Now().Add(t.accessDuration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(t.secretKey)
}

func (t *TokenService) createRefreshToken(userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID.String(),
		"exp": time.Now().Add(t.refreshDuration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(t.secretKey)
}

func (t *TokenService) ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("algoritmo de assinatura inesperado: %v", token.Header["alg"])
		}

		return t.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("token inválido")
	}

	return token, nil
}
