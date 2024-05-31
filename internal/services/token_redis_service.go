// internal/services/token_redis_service.go

package services

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/jeancarlosdanese/go-base-api/internal/domain/models"
)

type TokenRedisServiceInterface interface {
	SaveTokenDataRedis(user *models.User, token, refreshToken string, accessDuration time.Duration) error
	ValidateRefreshToken(refreshToken string) (*models.TokenDataRedis, error)
	GetTokenDataRedisFromToken(token string) (*models.TokenDataRedis, error)
}

type TokenRedisService struct {
	RedisService RedisServiceInterface
}

func NewTokenRedisService(redisService RedisServiceInterface) *TokenRedisService {
	return &TokenRedisService{
		RedisService: redisService,
	}
}

func (s *TokenRedisService) SaveTokenDataRedis(user *models.User, token, refreshToken string, accessDuration time.Duration) error {
	tokenDataRedis := prepareTokenDataRedis(user, token, refreshToken)
	tokenData, err := json.Marshal(tokenDataRedis)
	if err != nil {
		return err
	}
	if err := s.RedisService.Set("token:"+token, tokenData, accessDuration); err != nil {
		return err
	}
	return nil
}

func (s *TokenRedisService) ValidateRefreshToken(refreshToken string) (*models.TokenDataRedis, error) {
	result, err := s.RedisService.Get("refresh_token:" + refreshToken)
	if err != nil {
		return nil, err
	}
	var tokenDataRedis models.TokenDataRedis
	if err := json.Unmarshal([]byte(result), &tokenDataRedis); err != nil {
		return nil, err
	}
	return &tokenDataRedis, nil
}

func (s *TokenRedisService) GetTokenDataRedisFromToken(token string) (*models.TokenDataRedis, error) {
	result, err := s.RedisService.Get("token:" + token)
	if err != nil {
		return nil, err
	}
	var tokenDataRedis models.TokenDataRedis
	if err := json.Unmarshal([]byte(result), &tokenDataRedis); err != nil {
		return nil, err
	}
	return &tokenDataRedis, nil
}

// prepareTokenDataRedis prepares user data to be stored in Redis.
func prepareTokenDataRedis(user *models.User, token, refreshToken string) models.TokenDataRedis {
	// Map para evitar duplicatas e coletar todas as permissões
	policiesMap := make(map[string]bool)

	// Extrair permissões das roles
	for _, role := range user.Roles {
		for _, policy := range role.Policies {
			permKey := policyRoleToString(policy)
			policiesMap[permKey] = true
		}
	}

	// Adicionar permissões especiais
	for _, specialPermission := range user.SpecialPolicies {
		permKey := policyUserToString(specialPermission)
		policiesMap[permKey] = true
	}

	// Converter mapa para slice
	policiesSlice := make([]string, 0, len(policiesMap))
	for policy := range policiesMap {
		policiesSlice = append(policiesSlice, policy)
	}

	return models.TokenDataRedis{
		Token:        token,
		RefreshToken: &refreshToken,
		User: models.UserRedis{
			ID:       user.ID.String(),
			TenantID: user.TenantID.String(),
			Name:     user.Name,
			Username: user.Username,
			Email:    user.Email,
			Roles:    user.ExtractRoles(),
			Policies: policiesSlice,
		},
	}
}

func policyRoleToString(policy *models.PolicyRole) string {
	return fmt.Sprintf("%s:%s", policy.Endpoint.Name, policy.Actions)
}

func policyUserToString(policy *models.PolicyUser) string {
	return fmt.Sprintf("%s:%s", policy.Endpoint.Name, policy.Actions)
}
