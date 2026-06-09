package middleware

import (
	"net"
	"net/http"
	"strings"

	"github.com/ahmadzakyarifin/schoolpay/config"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

const (
	securityRedisPrefix = "schoolpay_security"
	rateLimitBypassKey  = "rate_limit_bypass"
)

func SecurityAccessMiddleware(cfg *config.Config, redisClient *redis.Client) gin.HandlerFunc {
	allowRules := nilToEmpty(cfg).RateLimitAllowIPs
	denyRules := nilToEmpty(cfg).RateLimitDenyIPs

	return func(c *gin.Context) {
		ip := strings.TrimSpace(c.ClientIP())
		if ip == "" {
			c.Next()
			return
		}

		if ipMatchesAny(ip, denyRules) || redisIPExists(c, redisClient, "deny_ip", ip) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"status":  false,
				"message": "Akses dari IP ini diblokir.",
				"data":    nil,
			})
			return
		}

		if ipMatchesAny(ip, allowRules) || redisIPExists(c, redisClient, "allow_ip", ip) {
			c.Set(rateLimitBypassKey, true)
		}

		c.Next()
	}
}

func nilToEmpty(cfg *config.Config) *config.Config {
	if cfg != nil {
		return cfg
	}
	return &config.Config{}
}

func redisIPExists(c *gin.Context, redisClient *redis.Client, listName, ip string) bool {
	if redisClient == nil {
		return false
	}
	key := securityRedisPrefix + ":" + listName + ":" + ip
	count, err := redisClient.Exists(c.Request.Context(), key).Result()
	return err == nil && count > 0
}

func ipMatchesAny(ip string, rules []string) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}

	for _, rule := range rules {
		rule = strings.TrimSpace(rule)
		if rule == "" {
			continue
		}

		if strings.Contains(rule, "/") {
			_, network, err := net.ParseCIDR(rule)
			if err == nil && network.Contains(parsedIP) {
				return true
			}
			continue
		}

		if strings.EqualFold(rule, ip) {
			return true
		}
	}

	return false
}
