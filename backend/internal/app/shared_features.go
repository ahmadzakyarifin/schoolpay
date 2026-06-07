package app

import (
	"github.com/ahmadzakyarifin/schoolpay/config"
	academicrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/academic/repository"
	auditrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/audit/repository"
	auditusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/audit/usecase"
	financerepo "github.com/ahmadzakyarifin/schoolpay/internal/module/finance/repository"
	financeusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/finance/usecase"
	notificationrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/notification/repository"
	notificationusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/notification/usecase"
	userauthrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/user_auth/repository"
	"github.com/ahmadzakyarifin/schoolpay/internal/websocket"
	"github.com/ahmadzakyarifin/schoolpay/pkg/utils"
	"github.com/uptrace/bun"
)

type sharedFeatures struct {
	UserRepository         userauthrepo.UserRepo
	AuthRepository         userauthrepo.AuthRepo
	StudentRepository      academicrepo.StudentRepo
	NotificationRepository notificationrepo.NotificationRepo

	AuditLogService            auditusecase.AuditLogService
	FinanceNotificationService notificationusecase.FinanceNotificationService
	PaymentService             financeusecase.PaymentService
	StudentBillService         financeusecase.StudentBillService
}

func buildSharedFeatures(
	database *bun.DB,
	appConfig *config.Config,
	messenger utils.Messenger,
	websocketHub *websocket.Hub,
) sharedFeatures {
	// Repository dasar yang dipakai lintas router, middleware, dan background job.
	userRepository := userauthrepo.NewUserRepo(database)
	authRepository := userauthrepo.NewAuthRepo(database)
	studentRepository := academicrepo.NewStudentRepo(database)
	notificationRepository := notificationrepo.NewNotificationRepo(database)

	// Service audit dan notifikasi finance dipakai oleh banyak fitur.
	auditRepository := auditrepo.NewAuditLogRepo(database)
	auditLogService := auditusecase.NewAuditLogService(auditRepository)
	financeNotificationService := notificationusecase.NewFinanceNotificationService(database, studentRepository, userRepository, notificationRepository, messenger, auditLogService)

	// Service pembayaran dan tagihan dipakai admin, finance, webhook, dan scheduler.
	paymentRepository := financerepo.NewPaymentRepo(database)
	studentBillRepository := financerepo.NewStudentBillRepo(database)
	billingRuleRepository := financerepo.NewBillingRuleRepo(database)
	paymentService := financeusecase.NewPaymentService(database, paymentRepository, studentBillRepository, studentRepository, financeNotificationService, appConfig, websocketHub, auditLogService)
	studentBillService := financeusecase.NewStudentBillService(database, studentBillRepository, billingRuleRepository, studentRepository, financeNotificationService, auditLogService)

	return sharedFeatures{
		UserRepository:             userRepository,
		AuthRepository:             authRepository,
		StudentRepository:          studentRepository,
		NotificationRepository:     notificationRepository,
		AuditLogService:            auditLogService,
		FinanceNotificationService: financeNotificationService,
		PaymentService:             paymentService,
		StudentBillService:         studentBillService,
	}
}
