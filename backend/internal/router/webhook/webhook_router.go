package webhook

import (
	"github.com/ahmadzakyarifin/schoolpay/config"
	"github.com/ahmadzakyarifin/schoolpay/internal/middleware"
	academicrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/academic/repository"
	auditusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/audit/usecase"
	financerepo "github.com/ahmadzakyarifin/schoolpay/internal/module/finance/repository"
	financeusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/finance/usecase"
	notificationrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/notification/repository"
	notificationusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/notification/usecase"
	supportrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/support/repository"
	supportusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/support/usecase"
	userauthrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/user_auth/repository"
	webhookhandler "github.com/ahmadzakyarifin/schoolpay/internal/module/webhook/delivery"
	webhookrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/webhook/repository"
	webhookusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/webhook/usecase"
	"github.com/ahmadzakyarifin/schoolpay/internal/websocket"
	"github.com/ahmadzakyarifin/schoolpay/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
)

func RouterWebhookSetup(
	g *gin.RouterGroup,
	db *bun.DB,
	cfg *config.Config,
	msg utils.Messenger,
	hub *websocket.Hub,
	paySvc financeusecase.PaymentService,
	sbSvc financeusecase.StudentBillService,
	finNotifSvc notificationusecase.FinanceNotificationService,
	auditSvc auditusecase.AuditLogService,
) {
	repo := webhookrepo.NewWebhookRepo(db)
	payRepo := financerepo.NewPaymentRepo(db)
	sbRepo := financerepo.NewStudentBillRepo(db)
	stuRepo := academicrepo.NewStudentRepo(db)
	userRepo := userauthrepo.NewUserRepo(db)
	notiRepo := notificationrepo.NewNotificationRepo(db)
	supportRepo := supportrepo.NewSupportRepo(db)

	waSvc := notificationusecase.NewWhatsAppService()
	supportSvc := supportusecase.NewSupportService(db, supportRepo, userRepo, waSvc, auditSvc)
	svc := webhookusecase.NewWebhookService(repo, waSvc, notiRepo, sbRepo, payRepo, stuRepo, userRepo, hub, supportSvc, cfg)
	pgSvc := financeusecase.NewPaymentGatewayService(cfg)

	hdl := webhookhandler.NewWebhookHandler(svc, paySvc, pgSvc, cfg)

	// Register Webhook to WAHA automatically on startup
	go waSvc.RegisterWebhook()

	wahaWebhookLimit := middleware.RateLimitAuthSaringan("webhook_waha", "ip", 600)
	paymentWebhookLimit := middleware.RateLimitAuthSaringan("webhook_payment", "ip", 120)

	g.POST("/wa-webhook", wahaWebhookLimit, hdl.HandleWAHA)
	g.POST("/payments/callback", paymentWebhookLimit, hdl.HandlePayment)
}
