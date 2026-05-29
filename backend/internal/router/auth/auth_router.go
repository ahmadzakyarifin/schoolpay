package auth

import (
	"github.com/ahmadzakyarifin/schoolpay/config"
	"github.com/ahmadzakyarifin/schoolpay/internal/middleware"
	auditrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/audit/repository"
	auditusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/audit/usecase"
	authhandler "github.com/ahmadzakyarifin/schoolpay/internal/module/user_auth/delivery"
	authrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/user_auth/repository"
	authusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/user_auth/usecase"
	"github.com/ahmadzakyarifin/schoolpay/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"github.com/uptrace/bun"
	"golang.org/x/time/rate"
)

func RouterAuthSetup(g *gin.RouterGroup, db *bun.DB, cfg *config.Config, msg utils.Messenger, redisClient *redis.Client, asynqClient *asynq.Client, userRepo authrepo.UserRepo) {
	repo := authrepo.NewAuthRepo(db, redisClient)
	auditRepo := auditrepo.NewAuditLogRepo(db)
	auditSvc := auditusecase.NewAuditLogService(auditRepo)
	svc := authusecase.NewAuthService(repo, msg, auditSvc, asynqClient)
	hdl := authhandler.NewAuthHandler(svc, cfg)

	loginLimit := middleware.RateLimitMiddleware(redisClient, "auth_login", rate.Limit(5.0/60.0), 5)
	forgotPasswordLimit := middleware.RateLimitMiddleware(redisClient, "auth_forgot_password", rate.Limit(3.0/60.0), 3)
	resetPasswordLimit := middleware.RateLimitMiddleware(redisClient, "auth_reset_password", rate.Limit(5.0/60.0), 5)

	auth := g.Group("/auth")
	{
		auth.POST("/login", loginLimit, hdl.Login)
		auth.POST("/refresh", hdl.RefreshToken)
		auth.POST("/logout", hdl.Logout)
		auth.POST("/forgot-password", forgotPasswordLimit, hdl.ForgotPassword)
		auth.POST("/reset-password", resetPasswordLimit, hdl.ResetPassword)
	}

	// Secured auth endpoints
	authSecured := g.Group("/auth")
	authSecured.Use(middleware.AuthMiddleware(cfg.JWTSecret, userRepo))
	{
		authSecured.POST("/change-password", hdl.ChangePassword)
	}
}
