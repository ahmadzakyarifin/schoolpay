package app

import (
	"context"
	"time"

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
	"github.com/redis/go-redis/v9"
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

func NewApp(database *bun.DB, appConfig *config.Config, redisClient *redis.Client) *App {
	if appConfig.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	if redisClient == nil {
		panic("Redis client is required")
	}

	rl, err := middleware.NewRedisRateLimiter(redisClient)
	if err != nil {
		panic("failed to initialize Redis rate limiter: " + err.Error())
	}
	middleware.SetDefaultRateLimiter(rl)
	middleware.SetRateLimitBlockDelay(time.Duration(appConfig.RateLimitBlockDelayMs) * time.Millisecond)

	routerEngine := gin.Default()
	configureCloudflareClientIP(routerEngine)
	routerEngine.Use(middleware.CORSMiddleware(appConfig.AppEnv, appConfig.FrontendURL))
	routerEngine.Use(middleware.SecurityAccessMiddleware(appConfig, redisClient))
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
	websocketGroup.Use(middleware.AuthMiddleware(appConfig.JWTSecret, sharedFeatures.UserRepository, redisClient))
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
	auth.RouterAuthSetup(apiGroup, appInstance.DB, &appInstance.Cfg, messenger, sharedFeatures.UserRepository, redisClient)

	adminGroup := apiGroup.Group("")
	adminGroup.Use(middleware.AuthMiddleware(appConfig.JWTSecret, sharedFeatures.UserRepository, redisClient))
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
	parentGroup.Use(middleware.AuthMiddleware(appConfig.JWTSecret, sharedFeatures.UserRepository, redisClient))
	parentGroup.Use(middleware.RoleMiddleware("parent"))
	parentGroup.Use(middleware.RateLimitPerUser("parent_private", 300))
	parent.SetupParentRoutes(parentGroup, appInstance.DB, &appInstance.Cfg, messenger, appInstance.Hub)

	// Update pemanggilan SetupFinanceRoutes agar menggunakan AuthMiddleware dengan redisClient secara konsisten (dilakukan di setup router masing-masing)
	// Karena SetupFinanceRoutes menggunakan middleware package, ia akan memanggil AuthMiddleware(jwtSecret, userRepo) yang sekarang butuh redisClient, maka kita perlu memperbarui parameter finance router juga.
	financerouter.SetupFinanceRoutes(apiGroup, appConfig.JWTSecret, sharedFeatures.UserRepository, sharedFeatures.PaymentService, sharedFeatures.StudentBillService, redisClient)

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

func configureCloudflareClientIP(router *gin.Engine) {
	router.RemoteIPHeaders = []string{"CF-Connecting-IP", "X-Forwarded-For", "X-Real-IP"}
	if err := router.SetTrustedProxies(cloudflareTrustedProxies()); err != nil {
		panic("failed to configure Cloudflare trusted proxies: " + err.Error())
	}
}

func cloudflareTrustedProxies() []string {
	return []string{
		"173.245.48.0/20",
		"103.21.244.0/22",
		"103.22.200.0/22",
		"103.31.4.0/22",
		"141.101.64.0/18",
		"108.162.192.0/18",
		"190.93.240.0/20",
		"188.114.96.0/20",
		"197.234.240.0/22",
		"198.41.128.0/17",
		"162.158.0.0/15",
		"104.16.0.0/13",
		"104.24.0.0/14",
		"172.64.0.0/13",
		"131.0.72.0/22",
		"2400:cb00::/32",
		"2606:4700::/32",
		"2803:f800::/32",
		"2405:b500::/32",
		"2405:8100::/32",
		"2a06:98c0::/29",
		"2c0f:f248::/32",
	}
}
