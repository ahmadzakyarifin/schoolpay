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
	"github.com/redis/go-redis/v9"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	redisstore "github.com/ulule/limiter/v3/drivers/store/redis"
)

const (
	ByIP     = "ip"
	ByDevice = "device"
	ByEmail  = "email"
	ByUser   = "user"
)

type Rule struct {
	Kind   string
	Limit  int64
	Period time.Duration
}

type RateLimiter struct {
	store limiter.Store
}

type emailPayload struct {
	Email string `json:"email"`
}

// Development / lokal / 1 instance.
func NewMemoryRateLimiter() *RateLimiter {
	store := memory.NewStoreWithOptions(limiter.StoreOptions{
		Prefix:          "schoolpay_rate_limit",
		CleanUpInterval: limiter.DefaultCleanUpInterval,
	})

	return &RateLimiter{store: store}
}

// Production / Redis / multi-instance ready.
func NewRedisRateLimiter(redisClient *redis.Client) (*RateLimiter, error) {
	store, err := redisstore.NewStoreWithOptions(redisClient, limiter.StoreOptions{
		Prefix: "schoolpay_rate_limit",
	})
	if err != nil {
		return nil, err
	}

	return &RateLimiter{store: store}, nil
}

// Rule builder.
func IP(limit int64, period time.Duration) Rule {
	return Rule{Kind: ByIP, Limit: limit, Period: period}
}

func Device(limit int64, period time.Duration) Rule {
	return Rule{Kind: ByDevice, Limit: limit, Period: period}
}

func Email(limit int64, period time.Duration) Rule {
	return Rule{Kind: ByEmail, Limit: limit, Period: period}
}

func User(limit int64, period time.Duration) Rule {
	return Rule{Kind: ByUser, Limit: limit, Period: period}
}

// Use adalah middleware utama.
// Semua rule dicek dulu. Kalau salah satu sudah limit, request ditolak.
func (r *RateLimiter) Use(scope string, rules ...Rule) gin.HandlerFunc {
	if scope == "" {
		scope = "global"
	}

	limiters := make([]*limiter.Limiter, len(rules))

	for i, rule := range rules {
		if rule.Limit < 1 {
			rule.Limit = 1
		}
		if rule.Period <= 0 {
			rule.Period = time.Minute
		}

		limiters[i] = limiter.New(r.store, limiter.Rate{
			Limit:  rule.Limit,
			Period: rule.Period,
		})
	}

	return func(c *gin.Context) {
		path := strings.ReplaceAll(c.FullPath(), "/", ":")
		if path == "" {
			path = strings.ReplaceAll(c.Request.URL.Path, "/", ":")
		}
		if path == "" || path == ":" {
			path = "unknown_path"
		}

		email := ""
		if needsEmail(rules) {
			var payload emailPayload
			if err := c.ShouldBindBodyWithJSON(&payload); err == nil {
				email = strings.ToLower(strings.TrimSpace(payload.Email))
			}
			restoreBody(c)
		}

		var blockedReason string
		var resetAt int64

		for i, rule := range rules {
			key, ok, unauthorized := makeKey(c, scope, path, rule.Kind, email)
			if unauthorized {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"status":  false,
					"message": "Akses ditolak. Sesi login tidak valid.",
					"data":    nil,
				})
				return
			}
			if !ok {
				continue
			}

			info, err := limiters[i].Get(c.Request.Context(), key)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"status":  false,
					"message": "Gagal memproses rate limit.",
					"data":    nil,
				})
				return
			}

			c.Header("X-RateLimit-"+rule.Kind+"-Limit", strconv.FormatInt(info.Limit, 10))
			c.Header("X-RateLimit-"+rule.Kind+"-Remaining", strconv.FormatInt(info.Remaining, 10))
			c.Header("X-RateLimit-"+rule.Kind+"-Reset", strconv.FormatInt(info.Reset, 10))

			if info.Reached {
				blockedReason = rule.Kind
				if info.Reset > resetAt {
					resetAt = info.Reset
				}
			}
		}

		if blockedReason != "" {
			retryAfter := int(resetAt - time.Now().Unix())
			if retryAfter < 1 {
				retryAfter = 1
			}

			c.Header("Retry-After", strconv.Itoa(retryAfter))
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"status":  false,
				"message": "Terlalu banyak request. Mohon tunggu beberapa saat.",
				"data": gin.H{
					"reason":              blockedReason,
					"retry_after_seconds": retryAfter,
				},
			})
			return
		}

		c.Next()
	}
}

func makeKey(c *gin.Context, scope, path, kind, email string) (string, bool, bool) {
	switch kind {
	case ByIP:
		return scope + ":" + path + ":ip:" + c.ClientIP(), true, false

	case ByDevice:
		deviceID := strings.TrimSpace(c.GetHeader("X-Device-ID"))
		if deviceID == "" {
			deviceID = "missing:" + c.GetHeader("User-Agent") + ":" + c.ClientIP()
		}
		if len(deviceID) > 120 {
			deviceID = deviceID[:120]
		}
		return scope + ":" + path + ":device:" + hash(deviceID), true, false

	case ByEmail:
		if email == "" {
			return "", false, false
		}
		return scope + ":" + path + ":email:" + hash(email), true, false

	case ByUser:
		userID, ok := c.Get("user_id")
		if !ok || strings.TrimSpace(fmt.Sprint(userID)) == "" {
			return "", false, true
		}
		return scope + ":" + path + ":user:" + fmt.Sprint(userID), true, false

	default:
		return "", false, false
	}
}

func needsEmail(rules []Rule) bool {
	for _, rule := range rules {
		if rule.Kind == ByEmail {
			return true
		}
	}
	return false
}

func restoreBody(c *gin.Context) {
	body, ok := c.Get(gin.BodyBytesKey)
	if !ok {
		return
	}

	bodyBytes, ok := body.([]byte)
	if !ok {
		return
	}

	c.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))
}

func hash(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])[:16]
}

// Compatibility layer for global/singleton usage to avoid refactoring all routers.
var defaultRateLimiter = NewMemoryRateLimiter()

func SetDefaultRateLimiter(rl *RateLimiter) {
	if rl != nil {
		defaultRateLimiter = rl
	}
}

func RateLimitAuthSaringan(scope string, target string, maxRequests int64) gin.HandlerFunc {
	var rule Rule
	switch target {
	case "ip":
		rule = IP(maxRequests, time.Minute)
	case "device":
		rule = Device(maxRequests, time.Minute)
	case "email":
		rule = Email(maxRequests, time.Minute)
	default:
		rule = IP(maxRequests, time.Minute)
	}
	return defaultRateLimiter.Use(scope, rule)
}

func RateLimitPerUser(scope string, maxRequests int64) gin.HandlerFunc {
	return defaultRateLimiter.Use(scope, User(maxRequests, time.Minute))
}
