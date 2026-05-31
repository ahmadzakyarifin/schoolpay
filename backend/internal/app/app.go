package app

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ahmadzakyarifin/schoolpay/config"
	_ "github.com/ahmadzakyarifin/schoolpay/docs"
	"github.com/ahmadzakyarifin/schoolpay/internal/middleware"
	academicrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/academic/repository"
	auditrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/audit/repository"
	auditusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/audit/usecase"
	financehandler "github.com/ahmadzakyarifin/schoolpay/internal/module/finance/delivery"
	financerepo "github.com/ahmadzakyarifin/schoolpay/internal/module/finance/repository"
	financeusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/finance/usecase"
	notificationrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/notification/repository"
	notificationusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/notification/usecase"
	userauthrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/user_auth/repository"
	"github.com/ahmadzakyarifin/schoolpay/internal/router/admin"
	"github.com/ahmadzakyarifin/schoolpay/internal/router/auth"
	"github.com/ahmadzakyarifin/schoolpay/internal/router/parent"
	"github.com/ahmadzakyarifin/schoolpay/internal/router/webhook"
	"github.com/ahmadzakyarifin/schoolpay/internal/websocket"
	"github.com/ahmadzakyarifin/schoolpay/internal/worker"
	"github.com/ahmadzakyarifin/schoolpay/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/uptrace/bun"
	"golang.org/x/time/rate"
)

type App struct {
	Server    *gin.Engine
	DB        *bun.DB
	Cfg       config.Config
	Messenger utils.Messenger
	Hub       *websocket.Hub
	RedisDB   *redis.Client
}

func NewApp(db *bun.DB, cfg *config.Config) *App {
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	g := gin.New()
	g.Use(middleware.LoggerMiddleware())
	g.Use(gin.Recovery())

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("custom_phone", utils.ValidatePhoneStruct)
		v.RegisterValidation("custom_nik", utils.ValidateNIKStruct)
	}

	// CORS Middleware
	g.Use(func(c *gin.Context) {
		allowedOrigins := map[string]bool{
			"http://localhost:5173": true,
			"http://localhost:5174": true,
			cfg.FrontendURL:         true,
		}

		origin := c.GetHeader("Origin")
		allowOrigin := allowedOrigins[origin]
		if !allowOrigin && cfg.AppEnv != "production" && origin != "" {
			if parsedOrigin, err := url.Parse(origin); err == nil {
				hostname := strings.ToLower(parsedOrigin.Hostname())
				allowOrigin = parsedOrigin.Scheme == "http" && (hostname == "localhost" || hostname == "127.0.0.1")
			}
		}
		if allowOrigin {
			c.Header("Access-Control-Allow-Origin", origin)
		}

		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, Accept, Origin, X-Idempotency-Key")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
		c.Header("Vary", "Origin")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	g.Static("/uploads", "./public/uploads")

	msg := utils.NewMessenger(cfg.WAHAURL, cfg.WAHAApiKey, cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPEmail, cfg.SMTPPass)
	hub := websocket.NewHub()
	go hub.Run()

	redisClient := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})

	asynqClient := asynq.NewClient(asynq.RedisClientOpt{Addr: cfg.RedisAddr})
	asynqSrv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: cfg.RedisAddr},
		asynq.Config{
			Concurrency: 5,
			Queues: map[string]int{
				"default": 10,
			},
		},
	)

	appInstance := &App{
		Server:    g,
		DB:        db,
		Cfg:       *cfg,
		Messenger: msg,
		Hub:       hub,
		RedisDB:   redisClient,
	}

	api := appInstance.Server.Group("/api")
	api.Use(middleware.RateLimitMiddleware(redisClient, "global", rate.Limit(50), 100))
	api.Use(middleware.IdempotencyMiddleware(redisClient))
	api.GET("/health", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()

		dbStatus := "ok"
		if err := db.PingContext(ctx); err != nil {
			dbStatus = err.Error()
		}

		redisStatus := "ok"
		if err := redisClient.Ping(ctx).Err(); err != nil {
			redisStatus = err.Error()
		}

		statusCode := http.StatusOK
		overall := "ok"
		if dbStatus != "ok" || redisStatus != "ok" {
			statusCode = http.StatusServiceUnavailable
			overall = "degraded"
		}

		c.JSON(statusCode, gin.H{
			"status": overall,
			"dependencies": gin.H{
				"database": dbStatus,
				"redis":    redisStatus,
			},
		})
	})

	// Initialize UserRepo early for AuthMiddleware
	userRepo := userauthrepo.NewUserRepo(db)

	wsGroup := api.Group("")
	wsGroup.Use(middleware.AuthMiddleware(cfg.JWTSecret, userRepo))
	wsGroup.Use(middleware.RoleMiddleware("admin", "parent"))
	wsGroup.GET("/ws", func(c *gin.Context) {
		websocket.ServeWs(hub, c.Writer, c.Request, cfg.FrontendURL)
	})

	// Singletons for services with background workers/schedulers
	payRepo := financerepo.NewPaymentRepo(db)
	sbRepo := financerepo.NewStudentBillRepo(db)
	brRepo := financerepo.NewBillingRuleRepo(db)
	stuRepo := academicrepo.NewStudentRepo(db)
	notiRepo := notificationrepo.NewNotificationRepo(db)

	auditRepo := auditrepo.NewAuditLogRepo(db)
	auditSvc := auditusecase.NewAuditLogService(auditRepo)
	finNotifSvc := notificationusecase.NewFinanceNotificationService(db, stuRepo, userRepo, notiRepo, msg, auditSvc, asynqClient)
	paySvc := financeusecase.NewPaymentService(db, payRepo, sbRepo, stuRepo, finNotifSvc, cfg, appInstance.Hub, auditSvc)
	sbSvc := financeusecase.NewStudentBillService(db, sbRepo, brRepo, stuRepo, finNotifSvc, auditSvc)

	// Start Schedulers
	sbSvc.RunScheduler()

	// Start Asynq Worker
	go worker.StartAsynqWorker(asynqSrv, db, stuRepo, userRepo, notiRepo, msg, auditSvc)

	//  Auth Router (Public)
	auth.RouterAuthSetup(api, appInstance.DB, &appInstance.Cfg, msg, redisClient, asynqClient, userRepo)

	//  Admin Router (Auth + Role Admin)
	adminGroup := api.Group("")
	adminGroup.Use(middleware.AuthMiddleware(cfg.JWTSecret, userRepo))
	adminGroup.Use(middleware.RoleMiddleware("admin"))
	admin.SetupAdminRoutes(adminGroup, api, appInstance.DB, &appInstance.Cfg, msg, appInstance.Hub, paySvc, sbSvc, finNotifSvc, auditSvc, redisClient, asynqClient)

	// Parent Router (Auth + Role Parent)
	parentGroup := api.Group("/parent")
	parentGroup.Use(middleware.AuthMiddleware(cfg.JWTSecret, userRepo))
	parentGroup.Use(middleware.RoleMiddleware("parent"))
	parent.SetupParentRoutes(parentGroup, appInstance.DB, &appInstance.Cfg, msg, redisClient, appInstance.Hub)

	//  Finance Feature (Cross-role)
	finGroup := api.Group("/finance")
	finGroup.Use(middleware.AuthMiddleware(cfg.JWTSecret, userRepo))
	finGroup.Use(middleware.RoleMiddleware("admin", "parent"))
	financePaymentLimit := middleware.RateLimitMiddleware(redisClient, "finance_payment", rate.Limit(10.0/60.0), 10)

	payHdl := financehandler.NewPaymentHandler(paySvc)
	sbHdl := financehandler.NewStudentBillHandler(sbSvc, paySvc)

	finGroup.POST("/payments", financePaymentLimit, payHdl.Process)
	finGroup.GET("/my-payments", payHdl.GetHistory)
	finGroup.POST("/payment-intent", financePaymentLimit, payHdl.CreateIntent)
	finGroup.GET("/my-bills", sbHdl.GetMyBills)
	finGroup.GET("/payments/:id/receipt", payHdl.GetReceipt)

	// 5. Webhook Router
	webhook.RouterWebhookSetup(appInstance.Server.Group(""), appInstance.DB, &appInstance.Cfg, msg, appInstance.Hub, paySvc, sbSvc, finNotifSvc, auditSvc, redisClient)

	// 6. Swagger API Documentation
	g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return appInstance
}
