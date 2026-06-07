package middleware

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

var globalStore = memory.NewStoreWithOptions(limiter.StoreOptions{
	Prefix:          "rate_limit",
	CleanUpInterval: limiter.DefaultCleanUpInterval,
})

type emailPayload struct {
	Email string `json:"email"`
}

// Fase sebelum login / publik: saringan bertingkat IP -> device -> email.
func RateLimitAuthSaringan(scope string, target string, maxRequests int64) gin.HandlerFunc {
	if maxRequests < 1 {
		maxRequests = 1
	}
	instance := limiter.New(globalStore, limiter.Rate{Period: time.Minute, Limit: maxRequests})

	return func(c *gin.Context) {
		var key string

		switch target {
		case "ip":
			key = "ip:" + c.ClientIP()

		case "device":
			deviceID := strings.TrimSpace(c.GetHeader("X-Device-ID"))
			if deviceID == "" {
				c.Next()
				return
			}
			if len(deviceID) > 80 {
				deviceID = deviceID[:80]
			}
			key = "device:" + shortHash(deviceID)

		case "email":
			var payload emailPayload
			if err := c.ShouldBindBodyWithJSON(&payload); err != nil {
				c.Next()
				return
			}
			if body, exists := c.Get(gin.BodyBytesKey); exists {
				if bodyBytes, ok := body.([]byte); ok {
					c.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))
				}
			}

			email := strings.ToLower(strings.TrimSpace(payload.Email))
			if email == "" {
				c.Next()
				return
			}
			key = "email:" + shortHash(email)

		default:
			c.Next()
			return
		}

		executeLimit(c, instance, scope+":"+key)
	}
}

// Fase setelah login: selalu kunci berdasarkan user_id dari AuthMiddleware.
func RateLimitPerUser(scope string, maxRequests int64) gin.HandlerFunc {
	if maxRequests < 1 {
		maxRequests = 1
	}
	instance := limiter.New(globalStore, limiter.Rate{Period: time.Minute, Limit: maxRequests})

	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists || strings.TrimSpace(fmt.Sprint(userID)) == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  false,
				"message": "Akses ditolak. Sesi login tidak valid.",
				"data":    nil,
			})
			return
		}

		executeLimit(c, instance, scope+":user:"+fmt.Sprint(userID))
	}
}

func executeLimit(c *gin.Context, instance *limiter.Limiter, fullKey string) {
	info, err := instance.Get(c.Request.Context(), fullKey)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Gagal memproses sistem keamanan rate limit.",
			"data":    nil,
		})
		return
	}

	c.Header("X-RateLimit-Limit", strconv.FormatInt(info.Limit, 10))
	c.Header("X-RateLimit-Remaining", strconv.FormatInt(info.Remaining, 10))
	c.Header("X-RateLimit-Reset", strconv.FormatInt(info.Reset, 10))

	if info.Reached {
		retryAfter := int(info.Reset - time.Now().Unix())
		if retryAfter < 1 {
			retryAfter = 1
		}

		c.Header("Retryk-After", strconv.Itoa(retryAfter))
		c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
			"status":  false,
			"message": "Terlalu banyak request. Mohon tunggu beberapa saat.",
			"data":    gin.H{"retry_after_seconds": retryAfter},
		})
		return
	}

	c.Next()
}

func shortHash(value string) string {
	hash := sha256.Sum256([]byte(value))

	return hex.EncodeToString(hash[:])[:16]
}
