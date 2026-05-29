package admin

import (
	"github.com/ahmadzakyarifin/schoolpay/config"
	academichandler "github.com/ahmadzakyarifin/schoolpay/internal/module/academic/delivery"
	academicrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/academic/repository"
	academicusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/academic/usecase"
	audithandler "github.com/ahmadzakyarifin/schoolpay/internal/module/audit/delivery"
	auditusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/audit/usecase"
	financehandler "github.com/ahmadzakyarifin/schoolpay/internal/module/finance/delivery"
	financerepo "github.com/ahmadzakyarifin/schoolpay/internal/module/finance/repository"
	financeusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/finance/usecase"
	notificationhandler "github.com/ahmadzakyarifin/schoolpay/internal/module/notification/delivery"
	notificationrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/notification/repository"
	notificationusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/notification/usecase"
	supporthandler "github.com/ahmadzakyarifin/schoolpay/internal/module/support/delivery"
	supportrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/support/repository"
	supportusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/support/usecase"
	userauthhandler "github.com/ahmadzakyarifin/schoolpay/internal/module/user_auth/delivery"
	userauthrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/user_auth/repository"
	userauthusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/user_auth/usecase"
	"github.com/ahmadzakyarifin/schoolpay/internal/middleware"
	"github.com/ahmadzakyarifin/schoolpay/internal/websocket"
	"github.com/ahmadzakyarifin/schoolpay/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"github.com/uptrace/bun"
	"golang.org/x/time/rate"
)

func SetupAdminRoutes(
	g *gin.RouterGroup,
	publicAPI *gin.RouterGroup,
	db *bun.DB,
	cfg *config.Config,
	msg utils.Messenger,
	hub *websocket.Hub,
	paySvc financeusecase.PaymentService,
	sbSvc financeusecase.StudentBillService,
	finNotifSvc notificationusecase.FinanceNotificationService,
	auditSvc auditusecase.AuditLogService,
	redisClient *redis.Client,
	asynqClient *asynq.Client,
) {

	// Repositories
	userRepo := userauthrepo.NewUserRepo(db)
	authRepo := userauthrepo.NewAuthRepo(db, redisClient)
	notiRepo := notificationrepo.NewNotificationRepo(db)
	stuRepo := academicrepo.NewStudentRepo(db)

	// Academic Repositories
	majorRepo := academicrepo.NewMajorRepo(db)
	classRepo := academicrepo.NewClassRepo(db)
	ayRepo := academicrepo.NewAcademicYearRepo(db)

	// Finance Repositories
	btRepo := financerepo.NewBillTypeRepo(db)
	brRepo := financerepo.NewBillingRuleRepo(db)
	sbRepo := financerepo.NewStudentBillRepo(db)
	repRepo := financerepo.NewFinanceReportRepo(db)

	// Shared Services
	waSvc := notificationusecase.NewWhatsAppService()
	supportRepo := supportrepo.NewSupportRepo(db)

	// Core Feature Services
	userSvc := userauthusecase.NewUserService(db, userRepo, authRepo, msg, notiRepo, stuRepo, auditSvc)
	stuSvc := academicusecase.NewStudentService(db, stuRepo, userRepo, authRepo, msg, notiRepo, ayRepo, majorRepo, classRepo, sbRepo, cfg, auditSvc)

	// Academic Services
	majorSvc := academicusecase.NewMajorService(db, majorRepo, auditSvc)
	classSvc := academicusecase.NewClassService(db, classRepo, majorRepo, auditSvc)
	aySvc := academicusecase.NewAcademicYearService(db, ayRepo, auditSvc)

	// Finance Services
	btSvc := financeusecase.NewBillTypeService(db, btRepo, auditSvc)
	brSvc := financeusecase.NewBillingRuleService(db, brRepo, finNotifSvc, auditSvc)
	repSvc := financeusecase.NewFinanceReportService(db, repRepo, auditSvc)
	supportSvc := supportusecase.NewSupportService(db, supportRepo, userRepo, waSvc, auditSvc)

	// Handlers
	userHdl := userauthhandler.NewUserHandler(userSvc, cfg)
	stuHdl := academichandler.NewStudentHandler(stuSvc, cfg)

	// Academic Handlers
	majorHdl := academichandler.NewMajorHandler(majorSvc)
	classHdl := academichandler.NewClassHandler(classSvc)
	ayHdl := academichandler.NewAcademicYearHandler(aySvc)

	// Finance Handlers
	btHdl := financehandler.NewBillTypeHandler(btSvc)
	brHdl := financehandler.NewBillingRuleHandler(brSvc)
	sbHdl := financehandler.NewStudentBillHandler(sbSvc, paySvc)
	payHdl := financehandler.NewPaymentHandler(paySvc)
	repHdl := financehandler.NewFinanceReportHandler(repSvc)

	dashHdl := financehandler.NewDashboardHandler(db, userRepo, stuRepo, repRepo, ayRepo, notiRepo, auditSvc)
	waHdl := notificationhandler.NewWhatsAppHandler(waSvc, notiRepo, msg, db, auditSvc)
	supportHdl := supporthandler.NewSupportHandler(supportSvc)
	upHdl := academichandler.NewUploadHandler()

	financeExportLimit := middleware.RateLimitMiddleware(redisClient, "finance_export", rate.Limit(10.0/60.0), 10)

	// Dashboard Stats
	g.GET("/dashboard/stats", dashHdl.GetStats)
	g.GET("/dashboard/communication-details", dashHdl.GetCommunicationDetails)
	g.GET("/dashboard/export", financeExportLimit, dashHdl.ExportGlobalReport)

	// Password & Profile
	// Change password dipindahkan ke auth_router.go agar bisa diakses parent juga

	// WebSocket Endpoint
	g.GET("/ws", func(c *gin.Context) {
		websocket.ServeWs(hub, c.Writer, c.Request)
	})

	notificationWriteLimit := middleware.RateLimitMiddleware(redisClient, "notification_write", rate.Limit(10.0/60.0), 10)
	whatsappChatLimit := middleware.RateLimitMiddleware(redisClient, "whatsapp_chat", rate.Limit(20.0/60.0), 20)
	supportWriteLimit := middleware.RateLimitMiddleware(redisClient, "support_write", rate.Limit(20.0/60.0), 20)
	masterWriteLimit := middleware.RateLimitMiddleware(redisClient, "master_write", rate.Limit(60.0/60.0), 60)
	financeWriteLimit := middleware.RateLimitMiddleware(redisClient, "finance_write", rate.Limit(30.0/60.0), 30)
	financePaymentLimit := middleware.RateLimitMiddleware(redisClient, "finance_payment_admin", rate.Limit(10.0/60.0), 10)

	// User Management

	// Public user routes
	publicAPI.POST("/users/activate", userHdl.Activate)

	users := g.Group("/users")
	{
		// Static Routes (Harus di atas)
		users.GET("", userHdl.GetAll)
		users.GET("/check-unique", userHdl.CheckUnique)
		users.GET("/notifications", userHdl.GetNotifications)
		users.GET("/export", userHdl.Export)
		users.POST("", masterWriteLimit, userHdl.Create)
		users.POST("/bulk-resend-notification", notificationWriteLimit, userHdl.BulkResendNotification)
		users.POST("/bulk-delete", masterWriteLimit, userHdl.BulkDelete)
		users.PATCH("/bulk-restore", masterWriteLimit, userHdl.BulkRestore)

		// Parameterized Routes (ID di bawah)
		users.GET("/:id", userHdl.GetByID)
		users.GET("/:id/students", stuHdl.GetByParentID)
		users.PUT("/:id", masterWriteLimit, userHdl.Update)
		users.DELETE("/:id", masterWriteLimit, userHdl.Delete)
		users.PATCH("/:id/status", masterWriteLimit, userHdl.ToggleStatus)
		users.POST("/:id/resend-notification", notificationWriteLimit, userHdl.ResendNotification)
		users.PATCH("/:id/restore", masterWriteLimit, userHdl.Restore)
		users.GET("/:id/dependency-info", userHdl.GetDependencyInfo)
	}

	// Student Management
	students := g.Group("/students")
	{
		// Static Routes
		students.GET("", stuHdl.GetAll)
		students.GET("/check-unique", stuHdl.CheckUnique)
		students.GET("/export", stuHdl.Export)
		students.GET("/filters", stuHdl.GetFilters)
		students.POST("", masterWriteLimit, stuHdl.Create)
		students.POST("/bulk-graduate", masterWriteLimit, stuHdl.BulkGraduate)
		students.POST("/bulk-promote", masterWriteLimit, stuHdl.BulkPromote)
		students.POST("/bulk-delete", masterWriteLimit, stuHdl.BulkDelete)
		students.PATCH("/bulk-restore", masterWriteLimit, stuHdl.BulkRestore)
		students.POST("/upload-photo", masterWriteLimit, upHdl.UploadStudentPhoto)

		// Parameterized Routes
		students.GET("/:id", stuHdl.GetByID)
		students.PUT("/:id", masterWriteLimit, stuHdl.Update)
		students.DELETE("/:id", masterWriteLimit, stuHdl.Delete)
		students.PATCH("/:id/status", masterWriteLimit, stuHdl.ToggleStatus)
		students.GET("/:id/parents", stuHdl.GetParents)
		students.GET("/:id/class-history", stuHdl.GetClassHistory)
		students.PATCH("/:id/restore", masterWriteLimit, stuHdl.Restore)
		students.GET("/:id/dependency-info", stuHdl.GetDependencyInfo)
	}

	// Academic Management
	academic := g.Group("/academic")
	{
		academic.GET("/major", majorHdl.GetAll)
		academic.GET("/major/check-unique", majorHdl.CheckUnique)
		academic.POST("/major", masterWriteLimit, majorHdl.Create)
		academic.POST("/major/bulk-delete", masterWriteLimit, majorHdl.BulkDelete)
		academic.PATCH("/major/bulk-restore", masterWriteLimit, majorHdl.BulkRestore)
		academic.PUT("/major/:id", masterWriteLimit, majorHdl.Update)
		academic.DELETE("/major/:id", masterWriteLimit, majorHdl.Delete)
		academic.PATCH("/major/:id/status", masterWriteLimit, majorHdl.ToggleStatus)
		academic.PATCH("/major/:id/restore", masterWriteLimit, majorHdl.Restore)
		academic.GET("/major/:id/dependency-info", majorHdl.GetDependencyInfo)

		academic.GET("/class", classHdl.GetAll)
		academic.GET("/class/check-unique", classHdl.CheckUnique)
		academic.POST("/class", masterWriteLimit, classHdl.Create)
		academic.PUT("/class/:id", masterWriteLimit, classHdl.Update)
		academic.DELETE("/class/:id", masterWriteLimit, classHdl.Delete)
		academic.PATCH("/class/:id/status", masterWriteLimit, classHdl.ToggleStatus)
		academic.PATCH("/class/:id/restore", masterWriteLimit, classHdl.Restore)
		academic.POST("/class/bulk-delete", masterWriteLimit, classHdl.BulkDelete)
		academic.PATCH("/class/bulk-restore", masterWriteLimit, classHdl.BulkRestore)
		academic.GET("/class/suggest-name", classHdl.SuggestNextName)
		academic.GET("/class/:id/dependency-info", classHdl.GetDependencyInfo)

		academic.GET("/years", ayHdl.GetAll)
		academic.GET("/years/active", ayHdl.GetActive)
		academic.GET("/years/check-unique", ayHdl.CheckUnique)
		academic.POST("/years", masterWriteLimit, ayHdl.Create)
		academic.POST("/years/bulk-delete", masterWriteLimit, ayHdl.BulkDelete)
		academic.PATCH("/years/bulk-restore", masterWriteLimit, ayHdl.BulkRestore)
		academic.PUT("/years/:id", masterWriteLimit, ayHdl.Update)
		academic.DELETE("/years/:id", masterWriteLimit, ayHdl.Delete)
		academic.PATCH("/years/:id/restore", masterWriteLimit, ayHdl.Restore)
		academic.POST("/years/:id/majors", masterWriteLimit, ayHdl.AssignMajors)
		academic.GET("/years/:id/majors", ayHdl.GetMajors)
		academic.POST("/years/:id/classes", masterWriteLimit, ayHdl.AssignClasses)
		academic.GET("/years/:id/classes", ayHdl.GetClasses)
		academic.GET("/years/:id/dependency-info", ayHdl.GetDependencyInfo)
	}

	// Finance Management
	finance := g.Group("/finance")
	{
		finance.GET("/bill-types", btHdl.GetAll)
		finance.GET("/bill-types/check-unique", btHdl.CheckUnique)
		finance.POST("/bill-types", financeWriteLimit, btHdl.Create)
		finance.POST("/bill-types/bulk-delete", financeWriteLimit, btHdl.BulkDelete)
		finance.PATCH("/bill-types/bulk-restore", financeWriteLimit, btHdl.BulkRestore)
		finance.PUT("/bill-types/:id", financeWriteLimit, btHdl.Update)
		finance.DELETE("/bill-types/:id", financeWriteLimit, btHdl.Delete)
		finance.PATCH("/bill-types/:id/status", financeWriteLimit, btHdl.ToggleStatus)
		finance.PATCH("/bill-types/:id/restore", financeWriteLimit, btHdl.Restore)
		finance.GET("/bill-types/:id/dependency-info", btHdl.GetDependencyInfo)

		finance.GET("/billing-rules", brHdl.GetAll)
		finance.GET("/billing-rules/check-unique", brHdl.CheckUnique)
		finance.POST("/billing-rules", financeWriteLimit, brHdl.Create)
		finance.POST("/billing-rules/bulk-delete", financeWriteLimit, brHdl.BulkDelete)
		finance.PATCH("/billing-rules/bulk-restore", financeWriteLimit, brHdl.BulkRestore)
		finance.PUT("/billing-rules/:id", financeWriteLimit, brHdl.Update)
		finance.DELETE("/billing-rules/:id", financeWriteLimit, brHdl.Delete)
		finance.PATCH("/billing-rules/:id/status", financeWriteLimit, brHdl.ToggleStatus)
		finance.PATCH("/billing-rules/:id/restore", financeWriteLimit, brHdl.Restore)
		finance.GET("/billing-rules/:id/dependency-info", brHdl.GetDependencyInfo)

		finance.GET("/bills", sbHdl.GetAll)
		finance.POST("/bills", financeWriteLimit, sbHdl.Create)
		finance.PUT("/bills/:id", financeWriteLimit, sbHdl.Update)
		finance.DELETE("/bills/:id", financeWriteLimit, sbHdl.Delete)
		finance.POST("/bills/:id/remind", notificationWriteLimit, sbHdl.Remind)
		finance.POST("/bills/:id/pay-manual", financePaymentLimit, sbHdl.MarkAsPaidManual)
		finance.POST("/generate-bills", financePaymentLimit, sbHdl.Generate)
		finance.POST("/generate-bills/bulk", financePaymentLimit, sbHdl.BulkGenerate)
		finance.POST("/generate-bills/bulk-cancel", financePaymentLimit, sbHdl.BulkCancel)

		finance.GET("/payments", payHdl.GetHistory)

		finance.GET("/arrears", repHdl.GetArrears)
		finance.GET("/export-trend", financeExportLimit, repHdl.ExportTrend)
	}

	// WhatsApp Management
	g.GET("/whatsapp/status", waHdl.GetStatus)
	g.GET("/whatsapp/qr", waHdl.GetQR)
	g.POST("/whatsapp/logout", whatsappChatLimit, waHdl.Logout)
	g.GET("/whatsapp/stats", waHdl.GetStats)
	g.GET("/whatsapp/logs", waHdl.GetLogs)
	g.POST("/whatsapp/notifications/:id/resend", notificationWriteLimit, waHdl.ResendSpecificNotification)
	g.GET("/whatsapp/chat/:phone", waHdl.GetChatHistory)
	g.POST("/whatsapp/chat/:phone", whatsappChatLimit, waHdl.SendChatMessage)

	support := g.Group("/support")
	{
		support.GET("/conversations", supportHdl.List)
		support.GET("/conversations/:id/messages", supportHdl.Messages)
		support.POST("/conversations/:id/reply", supportWriteLimit, supportHdl.Reply)
		support.PATCH("/conversations/:id/assign", supportWriteLimit, supportHdl.Assign)
		support.PATCH("/conversations/:id/close", supportWriteLimit, supportHdl.Close)
	}

	// Audit Log Management
	auditHdl := audithandler.NewAuditLogHandler(auditSvc)

	g.GET("/audit-logs", auditHdl.GetLogs)
	g.GET("/audit-logs/entity/:entityType/:entityID", auditHdl.GetEntityLogs)
}
