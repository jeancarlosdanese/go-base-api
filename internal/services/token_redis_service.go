// internal/services/token_redis_service.go

package services

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/jeancarlosdanese/go-base-api/internal/domain/models"
)

type TokenRedisServiceInterface interface {
	SaveUserRedis(user *models.User, token, refreshToken string, accessDuration time.Duration) error
	ValidateRefreshToken(refreshToken string) (*models.UserRedis, error)
	GetUserRedisFromToken(token string) (*models.UserRedis, error)
}

type TokenRedisService struct {
	RedisService RedisServiceInterface
}

func NewTokenRedisService(redisService RedisServiceInterface) *TokenRedisService {
	return &TokenRedisService{
		RedisService: redisService,
	}
}

func (s *TokenRedisService) SaveUserRedis(user *models.User, token, refreshToken string, accessDuration time.Duration) error {
	tokenDataRedis := prepareUserRedis(user)
	tokenData, err := json.Marshal(tokenDataRedis)
	if err != nil {
		return err
	}
	if err := s.RedisService.Set("token:"+token, tokenData, accessDuration); err != nil {
		return err
	}
	return nil
}

func (s *TokenRedisService) ValidateRefreshToken(refreshToken string) (*models.UserRedis, error) {
	result, err := s.RedisService.Get("refresh_token:" + refreshToken)
	if err != nil {
		return nil, err
	}
	var tokenDataRedis models.UserRedis
	if err := json.Unmarshal([]byte(result), &tokenDataRedis); err != nil {
		return nil, err
	}
	return &tokenDataRedis, nil
}

func (s *TokenRedisService) GetUserRedisFromToken(token string) (*models.UserRedis, error) {
	result, err := s.RedisService.Get("token:" + token)
	if err != nil {
		return nil, err
	}
	var tokenDataRedis models.UserRedis
	if err := json.Unmarshal([]byte(result), &tokenDataRedis); err != nil {
		return nil, err
	}
	return &tokenDataRedis, nil
}

// prepareUserRedis prepares user data to be stored in Redis.
func prepareUserRedis(user *models.User) models.UserRedis {
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

	return models.UserRedis{
		ID:       user.ID.String(),
		TenantID: user.TenantID.String(),
		Name:     user.Name,
		Username: user.Username,
		Email:    user.Email,
		Roles:    user.ExtractRoles(),
		Policies: policiesSlice,
	}
}

func policyRoleToString(policy *models.PolicyRole) string {
	return fmt.Sprintf("%s:%s", policy.Endpoint.Name, policy.Actions)
}

func policyUserToString(policy *models.PolicyUser) string {
	return fmt.Sprintf("%s:%s", policy.Endpoint.Name, policy.Actions)
}
