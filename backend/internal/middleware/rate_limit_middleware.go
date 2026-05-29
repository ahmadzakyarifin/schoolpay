package middleware

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"
	"golang.org/x/time/rate"
)

type localRateLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var localFallback = struct {
	sync.Mutex
	clients map[string]*localRateLimiter
}{clients: make(map[string]*localRateLimiter)}

func allowLocalFallback(key string, r rate.Limit, burst int) bool {
	localFallback.Lock()
	defer localFallback.Unlock()

	now := time.Now()
	for clientKey, client := range localFallback.clients {
		if now.Sub(client.lastSeen) > 10*time.Minute {
			delete(localFallback.clients, clientKey)
		}
	}

	client, ok := localFallback.clients[key]
	if !ok {
		client = &localRateLimiter{limiter: rate.NewLimiter(r, burst)}
		localFallback.clients[key] = client
	}
	client.lastSeen = now
	return client.limiter.Allow()
}

// RateLimitMiddleware uses Redis for cluster-safe limiting and falls back to
// a local token bucket when Redis is temporarily unavailable. The fallback keeps
// auth/payment endpoints protected instead of silently allowing unlimited traffic.
func RateLimitMiddleware(rdb *redis.Client, prefix string, r rate.Limit, b int) gin.HandlerFunc {
	limiter := redis_rate.NewLimiter(rdb)

	var redisLimit redis_rate.Limit
	if r < 1 {
		reqPerMin := int(r * 60)
		if reqPerMin < 1 {
			reqPerMin = 1
		}
		redisLimit = redis_rate.PerMinute(reqPerMin)
	} else {
		redisLimit = redis_rate.PerSecond(int(r))
	}
	redisLimit.Burst = b

	return func(c *gin.Context) {
		ip := c.ClientIP()
		key := "rate_limit:" + prefix + ":" + ip
		ctx := context.Background()

		res, err := limiter.Allow(ctx, key, redisLimit)
		if err != nil {
			if !allowLocalFallback(key, r, b) {
				c.JSON(http.StatusTooManyRequests, gin.H{
					"status":  "error",
					"message": "terlalu banyak permintaan, silakan coba lagi nanti",
				})
				c.Abort()
				return
			}
			c.Next()
			return
		}

		if res.Allowed == 0 {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"status":  "error",
				"message": "terlalu banyak permintaan, silakan coba lagi nanti",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
