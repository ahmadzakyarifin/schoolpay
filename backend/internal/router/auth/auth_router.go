package auth

import (
	"github.com/ahmadzakyarifin/schoolpay/config"
	"github.com/ahmadzakyarifin/schoolpay/internal/middleware"
	academicrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/academic/repository"
	auditrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/audit/repository"
	auditusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/audit/usecase"
	authhandler "github.com/ahmadzakyarifin/schoolpay/internal/module/user_auth/delivery"
	authrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/user_auth/repository"
	authusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/user_auth/usecase"
	"github.com/ahmadzakyarifin/schoolpay/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/uptrace/bun"
)

func RouterAuthSetup(g *gin.RouterGroup, db *bun.DB, cfg *config.Config, msg utils.Messenger, userRepo authrepo.UserRepo, redisClient *redis.Client) {
	authRepo := authrepo.NewAuthRepo(db)

	auditRepo := auditrepo.NewAuditLogRepo(db)
	auditService := auditusecase.NewAuditLogService(auditRepo)

	authService := authusecase.NewAuthService(authRepo, msg, auditService)
	authHandler := authhandler.NewAuthHandler(authService, cfg, redisClient)

	studentRepo := academicrepo.NewStudentRepo(db)
	profileHandler := authhandler.NewProfileHandler(db, userRepo, studentRepo, auditService)

	loginIPLimit := middleware.RateLimitAuthSaringan("auth_login", "ip", 20)
	loginDeviceLimit := middleware.RateLimitAuthSaringan("auth_login", "device", 10)
	loginEmailLimit := middleware.RateLimitAuthSaringan("auth_login", "email", 5)
	authRefreshLimit := middleware.RateLimitAuthSaringan("auth_refresh", "ip", 60)
	authLogoutLimit := middleware.RateLimitAuthSaringan("auth_logout", "ip", 60)
	forgotPasswordIPLimit := middleware.RateLimitAuthSaringan("auth_forgot_password", "ip", 10)
	forgotPasswordEmailLimit := middleware.RateLimitAuthSaringan("auth_forgot_password", "email", 3)
	resetPasswordLimit := middleware.RateLimitAuthSaringan("auth_reset_password", "ip", 10)
	changePasswordLimit := middleware.RateLimitPerUser("auth_change_password", 5)

	auth := g.Group("/auth")
	{
		auth.POST("/login", loginIPLimit, loginDeviceLimit, loginEmailLimit, authHandler.Login)
		auth.POST("/refresh", authRefreshLimit, authHandler.RefreshToken)
		auth.POST("/logout", authLogoutLimit, authHandler.Logout)
		auth.POST("/forgot-password", forgotPasswordIPLimit, forgotPasswordEmailLimit, authHandler.ForgotPassword)
		auth.POST("/reset-password", resetPasswordLimit, authHandler.ResetPassword)
	}

	// Secured auth endpoints
	authSecured := g.Group("/auth")
	authSecured.Use(middleware.AuthMiddleware(cfg.JWTSecret, userRepo, redisClient))
	authSecured.Use(middleware.RateLimitPerUser("auth_private", 300))
	{
		authSecured.GET("/me", profileHandler.Me)
		authSecured.PUT("/profile", profileHandler.UpdateProfile)
		authSecured.POST("/profile/photo", profileHandler.UploadPhoto)
		authSecured.POST("/change-password", changePasswordLimit, authHandler.ChangePassword)
	}
}
