package delivery

import (
	"context"
	"net/http"
	"time"

	"github.com/ahmadzakyarifin/schoolpay/config"
	"github.com/ahmadzakyarifin/schoolpay/internal/dto"
	"github.com/ahmadzakyarifin/schoolpay/internal/module/user_auth/usecase"
	"github.com/ahmadzakyarifin/schoolpay/pkg/utils"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	s   usecase.AuthService
	cfg *config.Config
}

func NewAuthHandler(service usecase.AuthService, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		s:   service,
		cfg: cfg,
	}
}

// Login godoc
// @Summary User Login
// @Description Authenticate a user and return an access token and refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login Credentials"
// @Success 200 {object} utils.Response{data=dto.LoginResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorValidationResponse(c, http.StatusBadRequest, "validasi gagal", utils.GetValidationErrors(err))
		return
	}

	ctx := c.Request.Context()
	ctx = context.WithValue(ctx, "ip_address", c.ClientIP())
	ctx = context.WithValue(ctx, "user_agent", c.Request.UserAgent())
	c.Request = c.Request.WithContext(ctx)

	res, err := h.s.Login(ctx, req, h.cfg.JWTSecret)
	if err != nil {
		utils.ErrorResponseRaw(c, http.StatusUnauthorized, err)
		return
	}

	// Set Refresh Token in HttpOnly Cookie
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("refresh_token", res.RefreshToken, int(time.Until(res.RefreshTokenExpiry).Seconds()), "/", "", false, true)

	utils.SuccessResponse(c, http.StatusOK, "login berhasil", res)
}

// RefreshToken godoc
// @Summary Refresh Access Token
// @Description Get a new access token using a refresh token from cookies
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response{data=dto.LoginResponse}
// @Failure 401 {object} utils.Response
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "session expired, please login again")
		return
	}

	ctx := c.Request.Context()
	ctx = context.WithValue(ctx, "ip_address", c.ClientIP())
	ctx = context.WithValue(ctx, "user_agent", c.Request.UserAgent())
	c.Request = c.Request.WithContext(ctx)

	req := dto.RefreshTokenRequest{RefreshToken: refreshToken}
	res, err := h.s.RefreshToken(ctx, req, h.cfg.JWTSecret)
	if err != nil {
		utils.ErrorResponseRaw(c, http.StatusUnauthorized, err)
		return
	}

	// Rotate Refresh Token in Cookie
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("refresh_token", res.RefreshToken, int(time.Until(res.RefreshTokenExpiry).Seconds()), "/", "", false, true)

	utils.SuccessResponse(c, http.StatusOK, "token berhasil diperbarui", res)
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req dto.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorValidationResponse(c, http.StatusBadRequest, "validasi gagal", utils.GetValidationErrors(err))
		return
	}

	ctx := c.Request.Context()
	ctx = context.WithValue(ctx, "ip_address", c.ClientIP())
	ctx = context.WithValue(ctx, "user_agent", c.Request.UserAgent())
	c.Request = c.Request.WithContext(ctx)

	if err := h.s.ForgotPassword(ctx, req, h.cfg); err != nil {
		utils.ErrorResponseRaw(c, http.StatusBadRequest, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "link reset password telah dikirim ke email", nil)
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req dto.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorValidationResponse(c, http.StatusBadRequest, "validasi gagal", utils.GetValidationErrors(err))
		return
	}

	ctx := c.Request.Context()
	ctx = context.WithValue(ctx, "ip_address", c.ClientIP())
	ctx = context.WithValue(ctx, "user_agent", c.Request.UserAgent())
	c.Request = c.Request.WithContext(ctx)

	if err := h.s.ResetPassword(ctx, req); err != nil {
		utils.ErrorResponseRaw(c, http.StatusBadRequest, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "password berhasil diubah", nil)
}

// Logout godoc
// @Summary User Logout
// @Description Invalidate the refresh token and clear cookies
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	if refreshToken, err := c.Cookie("refresh_token"); err == nil && refreshToken != "" {
		_ = h.s.Logout(c.Request.Context(), refreshToken)
	}
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)
	utils.SuccessResponse(c, http.StatusOK, "logout berhasil", nil)
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req struct {
		CurrentPassword string `json:"current_password" binding:"required"`
		NewPassword     string `json:"new_password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorValidationResponse(c, http.StatusBadRequest, "validasi gagal", utils.GetValidationErrors(err))
		return
	}

	userID, _ := c.Get("user_id")
	if err := h.s.ChangePassword(c.Request.Context(), userID.(uint), req.CurrentPassword, req.NewPassword); err != nil {
		utils.ErrorResponseRaw(c, http.StatusBadRequest, err)
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)
	utils.SuccessResponse(c, http.StatusOK, "password berhasil diperbarui, silakan login ulang", nil)
}
