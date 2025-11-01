// @file: internal/middleware/rate_limit_redis.go

package middleware

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// RateLimitRedisMiddleware implementa controle de taxa de requisições por IP usando Redis
func RateLimitRedisMiddleware(redisClient *redis.Client, requests int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Pula rate limiting para health check
		if c.Request.URL.Path == "/health" {
			c.Next()
			return
		}

		ip := getClientIP(c)
		key := fmt.Sprintf("rate_limit:%s", ip)

		count, err := redisClient.Incr(c.Request.Context(), key).Result()
		if err != nil {
			// Em caso de erro, permite a requisição (fail-open)
			c.Next()
			return
		}

		// Define expiração na primeira requisição
		if count == 1 {
			redisClient.Expire(c.Request.Context(), key, window)
		}

		// Verifica se excedeu o limite
		if count > int64(requests) {
			ttl, _ := redisClient.TTL(c.Request.Context(), key).Result()
			resetTime := time.Now().Add(ttl).Unix()

			c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", requests))
			c.Header("X-RateLimit-Remaining", "0")
			c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", resetTime))
			c.Header("Retry-After", fmt.Sprintf("%d", int(ttl.Seconds())))
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "Rate limit exceeded",
				"message":     "Too many requests. Please try again later.",
				"retry_after": int(ttl.Seconds()),
			})
			c.Abort()
			return
		}

		// Headers informativos
		remaining := requests - int(count)
		if remaining < 0 {
			remaining = 0
		}

		ttl, _ := redisClient.TTL(c.Request.Context(), key).Result()
		resetTime := time.Now().Add(ttl).Unix()

		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", requests))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", resetTime))

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
