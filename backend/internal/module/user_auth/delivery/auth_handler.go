package delivery

import (
	"net/http"
	"strings"
	"time"

	"github.com/ahmadzakyarifin/schoolpay/config"
	"github.com/ahmadzakyarifin/schoolpay/internal/dto"
	"github.com/ahmadzakyarifin/schoolpay/internal/helper"
	"github.com/ahmadzakyarifin/schoolpay/internal/module/user_auth/usecase"
	"github.com/ahmadzakyarifin/schoolpay/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type AuthHandler struct {
	s           usecase.AuthService
	cfg         *config.Config
	redisClient *redis.Client
}

func NewAuthHandler(service usecase.AuthService, cfg *config.Config, redisClient *redis.Client) *AuthHandler {
	return &AuthHandler{
		s:           service,
		cfg:         cfg,
		redisClient: redisClient,
	}
}

// Login godoc
// @Summary User Login
// @Description Authenticate a user and return an access token and refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login Credentials"
// @Success 200 {object} helper.Response{data=dto.LoginResponse}
// @Failure 400 {object} helper.Response
// @Failure 401 {object} helper.Response
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorValidationResponse(c, http.StatusBadRequest, "validasi gagal", helper.GetValidationErrors(err))
		return
	}

	audit := helper.BuildAuditMeta(c)
	res, err := h.s.Login(c.Request.Context(), req, audit, h.cfg.JWTSecret)
	if err != nil {
		helper.ErrorResponseRaw(c, http.StatusUnauthorized, err)
		return
	}

	// SameSite Lax Mode: Cookie hanya dikirim jika user membuka web kita secara langsung.
	// Ini mencegah website lain "menitip" request jahat menggunakan identitas/cookie user (keamanan CSRF).
	c.SetSameSite(http.SameSiteLaxMode)

	c.SetCookie("refresh_token", res.RefreshToken, int(time.Until(res.RefreshTokenExpiry).Seconds()), "/", "", false, true)

	helper.SuccessResponse(c, http.StatusOK, "login berhasil", res)
}

// RefreshToken godoc
// @Summary Refresh Access Token
// @Description Get a new access token using a refresh token from cookies
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} helper.Response{data=dto.LoginResponse}
// @Failure 401 {object} helper.Response
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		helper.ErrorResponse(c, http.StatusUnauthorized, "session expired, please login again")
		return
	}

	audit := helper.BuildAuditMeta(c)
	res, err := h.s.RefreshToken(c.Request.Context(), refreshToken, audit, h.cfg.JWTSecret)
	if err != nil {
		helper.ErrorResponseRaw(c, http.StatusUnauthorized, err)
		return
	}

	// Rotate Refresh Token in Cookie
	// SameSite Lax Mode: Cookie hanya dikirim jika user membuka web kita secara langsung.
	// Ini mencegah website lain "menitip" request jahat menggunakan identitas/cookie user (keamanan CSRF).
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("refresh_token", res.RefreshToken, int(time.Until(res.RefreshTokenExpiry).Seconds()), "/", "", false, true)

	helper.SuccessResponse(c, http.StatusOK, "token berhasil diperbarui", res)
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req dto.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorValidationResponse(c, http.StatusBadRequest, "validasi gagal", helper.GetValidationErrors(err))
		return
	}

	audit := helper.BuildAuditMeta(c)
	if err := h.s.ForgotPassword(c.Request.Context(), req, audit, h.cfg); err != nil {
		helper.ErrorResponseRaw(c, http.StatusBadRequest, err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "link reset password telah dikirim ke email", nil)
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req dto.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorValidationResponse(c, http.StatusBadRequest, "validasi gagal", helper.GetValidationErrors(err))
		return
	}

	audit := helper.BuildAuditMeta(c)
	if err := h.s.ResetPassword(c.Request.Context(), req, audit); err != nil {
		helper.ErrorResponseRaw(c, http.StatusBadRequest, err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "password berhasil diubah", nil)
}

// Logout godoc
// @Summary User Logout
// @Description Invalidate the refresh token and clear cookies
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} helper.Response
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err == nil && refreshToken != "" {
		_ = h.s.Logout(c.Request.Context(), refreshToken)
	}
	
	// Masukkan access token saat ini ke blacklist Redis
	h.blacklistToken(c)

	// SameSite Lax Mode: Cookie hanya dikirim jika user membuka web kita secara langsung.
	// Ini mencegah website lain "menitip" request jahat menggunakan identitas/cookie user (keamanan CSRF).
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)
	helper.SuccessResponse(c, http.StatusOK, "logout berhasil", nil)
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req dto.ChangePasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorValidationResponse(c, http.StatusBadRequest, "validasi gagal", helper.GetValidationErrors(err))
		return
	}

	userID, _ := c.Get("user_id")
	audit := helper.BuildAuditMeta(c)
	if err := h.s.ChangePassword(c.Request.Context(), userID.(uint), req, audit); err != nil {
		helper.ErrorResponseRaw(c, http.StatusBadRequest, err)
		return
	}

	// Masukkan access token saat ini ke blacklist Redis karena password berubah
	h.blacklistToken(c)

	// SameSite Lax Mode: Cookie hanya dikirim jika user membuka web kita secara langsung.
	// Ini mencegah website lain "menitip" request jahat menggunakan identitas/cookie user (keamanan CSRF).
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)
	helper.SuccessResponse(c, http.StatusOK, "password berhasil diperbarui, silakan login ulang", nil)
}

func (h *AuthHandler) blacklistToken(c *gin.Context) {
	if h.redisClient == nil {
		return
	}

	var tokenStr string
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
		}
	} else {
		tokenStr = c.Query("token")
	}

	if tokenStr == "" {
		return
	}

	claims, err := utils.ValidateToken(tokenStr, h.cfg.JWTSecret)
	if err != nil {
		return
	}

	// Hitung durasi sisa masa aktif token
	expiry := claims.ExpiresAt.Time
	duration := time.Until(expiry)
	if duration > 0 {
		_ = h.redisClient.Set(c.Request.Context(), "blacklist:"+tokenStr, "1", duration).Err()
	}
}
