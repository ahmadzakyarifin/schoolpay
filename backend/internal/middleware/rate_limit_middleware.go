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
	redisstore "github.com/ulule/limiter/v3/drivers/store/redis"
)

const (
	ByIP          = "ip"
	ByDevice      = "device"
	ByEmail       = "email"
	ByUser        = "user"
	ByIPEmail     = "ip_email"
	ByDeviceEmail = "device_email"
)

type Rule struct {
	Name   string
	Kind   string
	Limit  int64
	Period time.Duration
}

type RateLimiter struct {
	store limiter.Store
}

type preparedRule struct {
	rule    Rule
	limiter *limiter.Limiter
}

type emailPayload struct {
	Email string `json:"email"`
}

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
	return namedRule(ByIP, ByIP, limit, period)
}

func Device(limit int64, period time.Duration) Rule {
	return namedRule(ByDevice, ByDevice, limit, period)
}

func Email(limit int64, period time.Duration) Rule {
	return namedRule(ByEmail, ByEmail, limit, period)
}

func EmailMinute(limit int64, period time.Duration) Rule {
	return namedRule("email_minute", ByEmail, limit, period)
}

func EmailHour(limit int64, period time.Duration) Rule {
	return namedRule("email_hour", ByEmail, limit, period)
}

func EmailDay(limit int64, period time.Duration) Rule {
	return namedRule("email_day", ByEmail, limit, period)
}

func User(limit int64, period time.Duration) Rule {
	return namedRule(ByUser, ByUser, limit, period)
}

func IPEmail(limit int64, period time.Duration) Rule {
	return namedRule(ByIPEmail, ByIPEmail, limit, period)
}

func DeviceEmail(limit int64, period time.Duration) Rule {
	return namedRule(ByDeviceEmail, ByDeviceEmail, limit, period)
}

func namedRule(name, kind string, limit int64, period time.Duration) Rule {
	return Rule{Name: name, Kind: kind, Limit: limit, Period: period}
}

// Use adalah middleware utama.
// Alur:
//  1. Cek setiap rule Redis.
//  2. Jika hard limit kena, balas 429.
//  3. Jika limit yang bisa diselesaikan CAPTCHA kena, minta CAPTCHA.
//  4. Jika aman, lanjut ke handler.
func (r *RateLimiter) Use(scope string, rules ...Rule) gin.HandlerFunc {
	scope = normalizeScope(scope)
	preparedRules := r.prepareRules(rules)

	return func(c *gin.Context) {
		if c.GetBool(rateLimitBypassKey) {
			c.Next()
			return
		}

		request := requestIdentity{
			scope: scope,
			path:  rateLimitPath(c),
			email: emailFromBodyIfNeeded(c, rules),
		}

		result, ok := checkRateLimitRules(c, preparedRules, request)
		if !ok {
			return
		}

		if result.hardBlockedReason != "" {
			abortTooManyRequests(c, result.hardBlockedReason, result.hardResetAt)
			return
		}

		if result.captchaReason != "" && !c.GetBool(captchaVerifiedKey) {
			abortCaptchaRequired(c, result.captchaReason, result.captchaResetAt)
			return
		}

		c.Next()
	}
}

type requestIdentity struct {
	scope string
	path  string
	email string
}

type rateLimitResult struct {
	captchaReason     string
	captchaResetAt    int64
	hardBlockedReason string
	hardResetAt       int64
}

func (r *RateLimiter) prepareRules(rules []Rule) []preparedRule {
	prepared := make([]preparedRule, len(rules))
	for i, rule := range rules {
		rule = normalizeRule(rule)
		prepared[i] = preparedRule{
			rule: rule,
			limiter: limiter.New(r.store, limiter.Rate{
				Limit:  rule.Limit,
				Period: rule.Period,
			}),
		}
	}
	return prepared
}

func normalizeRule(rule Rule) Rule {
	if rule.Limit < 1 {
		rule.Limit = 1
	}
	if strings.TrimSpace(rule.Name) == "" {
		rule.Name = rule.Kind
	}
	if rule.Period <= 0 {
		rule.Period = time.Minute
	}
	return rule
}

func normalizeScope(scope string) string {
	if scope == "" {
		return "global"
	}
	return scope
}

func checkRateLimitRules(c *gin.Context, rules []preparedRule, request requestIdentity) (rateLimitResult, bool) {
	var result rateLimitResult

	for _, item := range rules {
		key, ok, unauthorized := makeKey(c, request.scope, request.path, item.rule, request.email)
		if unauthorized {
			abortUnauthorizedRateLimit(c)
			return result, false
		}
		if !ok {
			continue
		}

		info, err := item.limiter.Get(c.Request.Context(), key)
		if err != nil {
			abortRateLimitError(c)
			return result, false
		}

		writeRateLimitHeaders(c, item.rule.Name, info)
		if !info.Reached {
			continue
		}

		if captchaCanSatisfy(item.rule.Kind) {
			result.captchaReason = item.rule.Name
			if info.Reset > result.captchaResetAt {
				result.captchaResetAt = info.Reset
			}
			continue
		}

		result.hardBlockedReason = item.rule.Name
		if info.Reset > result.hardResetAt {
			result.hardResetAt = info.Reset
		}
	}

	return result, true
}

func writeRateLimitHeaders(c *gin.Context, name string, info limiter.Context) {
	c.Header("X-RateLimit-"+name+"-Limit", strconv.FormatInt(info.Limit, 10))
	c.Header("X-RateLimit-"+name+"-Remaining", strconv.FormatInt(info.Remaining, 10))
	c.Header("X-RateLimit-"+name+"-Reset", strconv.FormatInt(info.Reset, 10))
}

func abortUnauthorizedRateLimit(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"status":  false,
		"message": "Akses ditolak. Sesi login tidak valid.",
		"data":    nil,
	})
}

func abortRateLimitError(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
		"status":  false,
		"message": "Gagal memproses rate limit.",
		"data":    nil,
	})
}

func abortTooManyRequests(c *gin.Context, reason string, resetAt int64) {
	retryAfter := retryAfterSeconds(resetAt)
	applyRateLimitDelay(retryAfter)

	c.Header("Retry-After", strconv.Itoa(retryAfter))
	c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
		"status":  false,
		"message": "Terlalu banyak request. Mohon tunggu beberapa saat.",
		"data": gin.H{
			"reason":              reason,
			"retry_after_seconds": retryAfter,
		},
	})
}

func abortCaptchaRequired(c *gin.Context, reason string, resetAt int64) {
	c.AbortWithStatusJSON(http.StatusOK, gin.H{
		"status":  false,
		"message": "Verifikasi tambahan diperlukan.",
		"data": gin.H{
			"captcha_required":    true,
			"reason":              reason,
			"retry_after_seconds": retryAfterSeconds(resetAt),
		},
	})
}

func retryAfterSeconds(resetAt int64) int {
	return maxInt(1, int(resetAt-time.Now().Unix()))
}

func rateLimitPath(c *gin.Context) string {
	path := strings.ReplaceAll(c.FullPath(), "/", ":")
	if path == "" {
		path = strings.ReplaceAll(c.Request.URL.Path, "/", ":")
	}
	if path == "" || path == ":" {
		return "unknown_path"
	}
	return path
}

func emailFromBodyIfNeeded(c *gin.Context, rules []Rule) string {
	if !needsEmail(rules) {
		return ""
	}

	var payload emailPayload
	if err := c.ShouldBindBodyWithJSON(&payload); err != nil {
		restoreBody(c)
		return ""
	}

	restoreBody(c)
	return strings.ToLower(strings.TrimSpace(payload.Email))
}

func makeKey(c *gin.Context, scope, path string, rule Rule, email string) (string, bool, bool) {
	prefix := scope + ":" + path + ":" + rule.Name

	switch rule.Kind {
	case ByIP:
		return prefix + ":ip:" + c.ClientIP(), true, false

	case ByDevice:
		return prefix + ":device:" + hash(resolveDeviceID(c)), true, false

	case ByEmail:
		if email == "" {
			return "", false, false
		}
		return prefix + ":email:" + hash(email), true, false

	case ByUser:
		userID, ok := c.Get("user_id")
		if !ok || strings.TrimSpace(fmt.Sprint(userID)) == "" {
			return "", false, true
		}
		return prefix + ":user:" + fmt.Sprint(userID), true, false

	case ByIPEmail:
		if email == "" {
			return "", false, false
		}
		return prefix + ":ip_email:" + hash(c.ClientIP()+":"+email), true, false

	case ByDeviceEmail:
		if email == "" {
			return "", false, false
		}
		return prefix + ":device_email:" + hash(resolveDeviceID(c)+":"+email), true, false

	default:
		return "", false, false
	}
}

func needsEmail(rules []Rule) bool {
	for _, rule := range rules {
		if rule.Kind == ByEmail || rule.Kind == ByIPEmail || rule.Kind == ByDeviceEmail {
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

func resolveDeviceID(c *gin.Context) string {
	deviceID := strings.TrimSpace(c.GetHeader("X-Device-ID"))
	if deviceID != "" {
		if len(deviceID) > 120 {
			deviceID = deviceID[:120]
		}
		return deviceID
	}

	userAgent := c.GetHeader("User-Agent")
	if len(userAgent) > 100 {
		userAgent = userAgent[:100]
	}

	ip := c.ClientIP()
	parts := strings.Split(ip, ".")
	subnet := ip
	if len(parts) == 4 {
		subnet = parts[0] + "." + parts[1] + "." + parts[2] + ".0"
	}

	return "fallback:" + userAgent + ":" + subnet
}

func captchaCanSatisfy(kind string) bool {
	return captchaChallengeEnabled && (kind == ByDevice || kind == ByEmail || kind == ByIPEmail || kind == ByDeviceEmail)
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Compatibility layer for global/singleton usage to avoid refactoring all routers.
var defaultRateLimiter *RateLimiter
var defaultRateLimitBlockDelay time.Duration

func SetDefaultRateLimiter(rl *RateLimiter) {
	if rl != nil {
		defaultRateLimiter = rl
	}
}

func SetRateLimitBlockDelay(delay time.Duration) {
	if delay < 0 {
		delay = 0
	}
	if delay > 5*time.Second {
		delay = 5 * time.Second
	}
	defaultRateLimitBlockDelay = delay
}

func RateLimitAuthSaringan(scope string, target string, maxRequests int64) gin.HandlerFunc {
	if defaultRateLimiter == nil {
		panic("rate limiter Redis store is not initialized")
	}

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

func RateLimitRules(scope string, rules ...Rule) gin.HandlerFunc {
	if defaultRateLimiter == nil {
		panic("rate limiter Redis store is not initialized")
	}
	return defaultRateLimiter.Use(scope, rules...)
}

func RateLimitPerUser(scope string, maxRequests int64) gin.HandlerFunc {
	if defaultRateLimiter == nil {
		panic("rate limiter Redis store is not initialized")
	}
	return defaultRateLimiter.Use(scope, User(maxRequests, time.Minute))
}

func applyRateLimitDelay(retryAfterSeconds int) {
	if defaultRateLimitBlockDelay <= 0 {
		return
	}

	multiplier := retryAfterSeconds / 10
	if multiplier < 1 {
		multiplier = 1
	}
	if multiplier > 5 {
		multiplier = 5
	}

	delay := time.Duration(multiplier) * defaultRateLimitBlockDelay
	if delay > 5*time.Second {
		delay = 5 * time.Second
	}
	time.Sleep(delay)
}
