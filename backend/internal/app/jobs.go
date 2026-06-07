package app

import (
	"context"
	"time"

	"github.com/ahmadzakyarifin/schoolpay/config"
	"github.com/ahmadzakyarifin/schoolpay/internal/middleware"
	academicrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/academic/repository"
	auditusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/audit/usecase"
	financeusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/finance/usecase"
	notificationrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/notification/repository"
	userauthrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/user_auth/repository"
	"github.com/ahmadzakyarifin/schoolpay/internal/worker"
	"github.com/ahmadzakyarifin/schoolpay/pkg/utils"
	"github.com/uptrace/bun"
)

func startBackgroundJobs(
	appContext context.Context,
	database *bun.DB,
	appConfig *config.Config,
	messenger utils.Messenger,
	studentRepository academicrepo.StudentRepo,
	userRepository userauthrepo.UserRepo,
	notificationRepository notificationrepo.NotificationRepo,
	authRepository userauthrepo.AuthRepo,
	auditLogService auditusecase.AuditLogService,
	studentBillService financeusecase.StudentBillService,
) {
	// Scheduler tagihan berjalan berkala untuk proses otomatis terkait billing.
	studentBillService.RunScheduler()

	// Worker database memproses background_jobs, misalnya email dan notifikasi.
	go worker.StartDatabaseWorker(appContext, database, studentRepository, userRepository, notificationRepository, authRepository, messenger, auditLogService, appConfig)

	// Cleanup idempotency menjaga tabel idempotency_keys tidak membesar tanpa batas.
	middleware.StartIdempotencyCleanupJob(appContext, database, time.Hour, 24*time.Hour, 15*time.Minute)
}
