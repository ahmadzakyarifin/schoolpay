package app

import (
	"context"

	"github.com/ahmadzakyarifin/schoolpay/config"
	_ "github.com/ahmadzakyarifin/schoolpay/docs"
	"github.com/ahmadzakyarifin/schoolpay/internal/middleware"
	"github.com/ahmadzakyarifin/schoolpay/internal/router/admin"
	"github.com/ahmadzakyarifin/schoolpay/internal/router/auth"
	financerouter "github.com/ahmadzakyarifin/schoolpay/internal/router/finance"
	"github.com/ahmadzakyarifin/schoolpay/internal/router/parent"
	"github.com/ahmadzakyarifin/schoolpay/internal/router/webhook"
	"github.com/ahmadzakyarifin/schoolpay/internal/websocket"
	"github.com/ahmadzakyarifin/schoolpay/pkg/utils"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/uptrace/bun"
)

type App struct {
	Server    *gin.Engine
	DB        *bun.DB
	Cfg       config.Config
	Messenger utils.Messenger
	Hub       *websocket.Hub
}

func NewApp(database *bun.DB, appConfig *config.Config) *App {
	if appConfig.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	routerEngine := gin.Default()
	routerEngine.Use(middleware.CORSMiddleware(appConfig.AppEnv, appConfig.FrontendURL))
	routerEngine.Static("/uploads", "./public/uploads")

	// Buat komponen global yang dipakai banyak fitur.
	// Messenger dipakai auth, notification, support, dan worker.
	messenger := utils.NewMessenger(appConfig.WAHAURL, appConfig.WAHAApiKey, appConfig.SMTPHost, appConfig.SMTPPort, appConfig.SMTPEmail, appConfig.SMTPPass)

	// Hub websocket adalah jalur realtime global, misalnya event pembayaran baru.
	websocketHub := websocket.NewHub()
	go websocketHub.Run()

	appInstance := &App{
		Server:    routerEngine,
		DB:        database,
		Cfg:       *appConfig,
		Messenger: messenger,
		Hub:       websocketHub,
	}

	// Buat group API dan pasang middleware yang berlaku untuk semua endpoint API.
	apiGroup := appInstance.Server.Group("/api")
	apiGroup.Use(middleware.RateLimitAuthSaringan("global", "ip", 3000))
	apiGroup.Use(middleware.IdempotencyMiddleware(database))

	appContext := context.Background()
	sharedFeatures := buildSharedFeatures(database, appConfig, messenger, websocketHub)

	// Endpoint realtime global. Route ini sengaja ada di app karena dipakai
	// lintas role, bukan milik modul admin/parent/finance tertentu.
	websocketGroup := apiGroup.Group("")
	websocketGroup.Use(middleware.AuthMiddleware(appConfig.JWTSecret, sharedFeatures.UserRepository))
	websocketGroup.Use(middleware.RoleMiddleware("admin", "parent"))
	websocketGroup.Use(middleware.RateLimitPerUser("websocket", 60))
	// WebSocket dipakai untuk push event realtime dari server ke frontend,
	// misalnya notifikasi pembayaran baru ke dashboard admin tanpa refresh halaman.
	websocketGroup.GET("/ws", func(ginContext *gin.Context) {
		websocket.ServeWs(websocketHub, ginContext.Writer, ginContext.Request, appConfig.FrontendURL)
	})

	// Jalankan proses non-HTTP yang hidup bersama aplikasi.
	// Jalankan pekerjaan aplikasi yang hidup di background, seperti scheduler
	// tagihan, worker database, dan cleanup idempotency key lama.
	startBackgroundJobs(
		appContext,
		database,
		appConfig,
		messenger,
		sharedFeatures.StudentRepository,
		sharedFeatures.UserRepository,
		sharedFeatures.NotificationRepository,
		sharedFeatures.AuthRepository,
		sharedFeatures.AuditLogService,
		sharedFeatures.StudentBillService,
	)

	// Daftarkan router fitur. Detail endpoint masing-masing fitur berada
	// di package router terkait agar app.go tetap mudah dibaca.
	auth.RouterAuthSetup(apiGroup, appInstance.DB, &appInstance.Cfg, messenger, sharedFeatures.UserRepository)

	adminGroup := apiGroup.Group("")
	adminGroup.Use(middleware.AuthMiddleware(appConfig.JWTSecret, sharedFeatures.UserRepository))
	adminGroup.Use(middleware.RoleMiddleware("admin"))
	adminGroup.Use(middleware.RateLimitPerUser("admin_private", 600))
	admin.SetupAdminRoutes(
		adminGroup,
		apiGroup,
		appInstance.DB,
		&appInstance.Cfg,
		messenger,
		appInstance.Hub,
		sharedFeatures.PaymentService,
		sharedFeatures.StudentBillService,
		sharedFeatures.FinanceNotificationService,
		sharedFeatures.AuditLogService,
	)

	parentGroup := apiGroup.Group("/parent")
	parentGroup.Use(middleware.AuthMiddleware(appConfig.JWTSecret, sharedFeatures.UserRepository))
	parentGroup.Use(middleware.RoleMiddleware("parent"))
	parentGroup.Use(middleware.RateLimitPerUser("parent_private", 300))
	parent.SetupParentRoutes(parentGroup, appInstance.DB, &appInstance.Cfg, messenger, appInstance.Hub)

	financerouter.SetupFinanceRoutes(apiGroup, appConfig.JWTSecret, sharedFeatures.UserRepository, sharedFeatures.PaymentService, sharedFeatures.StudentBillService)

	// Webhook memakai root group karena callback gateway biasanya tidak berada
	// di bawah /api dan punya mekanisme autentikasi/verifikasi sendiri.
	webhook.RouterWebhookSetup(
		appInstance.Server.Group(""),
		appInstance.DB,
		&appInstance.Cfg,
		messenger,
		appInstance.Hub,
		sharedFeatures.PaymentService,
		sharedFeatures.StudentBillService,
		sharedFeatures.FinanceNotificationService,
		sharedFeatures.AuditLogService,
	)

	// Dokumentasi API.
	routerEngine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return appInstance
}
