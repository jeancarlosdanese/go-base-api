// internal/services/token_redis_service.go

package services

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/jeancarlosdanese/go-base-api/internal/domain/models"
)

type TokenRedisService struct {
	*RedisService
}

func NewTokenRedisService(redisService *RedisService) *TokenRedisService {
	return &TokenRedisService{
		RedisService: redisService,
	}
}

func (s *TokenRedisService) SaveUserTokenInfo(user *models.User, token, refreshToken string) error {
	userDataRedis := prepareUserDataRedis(user, token, refreshToken)
	userData, err := json.Marshal(userDataRedis)
	if err != nil {
		return err
	}
	return s.Set("token:"+token, userData, time.Hour*1)
}

func (s *TokenRedisService) GetUserFromToken(token string) (*models.UserDataRedis, error) {
	result, err := s.Get("token:" + token)
	if err != nil {
		return nil, err
	}
	var userDataRedis models.UserDataRedis
	if err := json.Unmarshal([]byte(result), &userDataRedis); err != nil {
		return nil, err
	}
	return &userDataRedis, nil
}

// prepareUserDataRedis prepares user data to be stored in Redis.
func prepareUserDataRedis(user *models.User, token, refreshToken string) models.UserDataRedis {
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

	return models.UserDataRedis{
		Token:        token,
		RefreshToken: &refreshToken,
		User: models.UserRedis{
			ID:       user.ID.String(),
			TenantID: user.TenantID.String(),
			Name:     user.Name,
			Username: user.Username,
			Email:    user.Email,
		},
		Roles:    user.ExtractRoles(),
		Policies: policiesSlice,
	}
}

func policyRoleToString(policy *models.PolicyRole) string {
	fmt.Printf("%s:%s:%s", fmt.Sprintf("%d", policy.RoleID), policy.Endpoint.Name, policy.Actions)
	return fmt.Sprintf("%s:%s", policy.Endpoint.Name, policy.Actions)
}

func policyUserToString(policy *models.PolicyUser) string {
	return fmt.Sprintf("%s:%s", policy.Endpoint.Name, policy.Actions)
}
