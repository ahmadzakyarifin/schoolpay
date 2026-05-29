package usecase

import (
	"encoding/json"
	"fmt"

	academicrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/academic/repository"
	auditusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/audit/usecase"
	financedomain "github.com/ahmadzakyarifin/schoolpay/internal/module/finance/domain"
	notificationrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/notification/repository"
	userauthrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/user_auth/repository"
	"github.com/ahmadzakyarifin/schoolpay/pkg/utils"
	"github.com/hibiken/asynq"
	"github.com/uptrace/bun"
)

const TaskFinanceNotification = "notification:finance"

type FinanceNotifyJob struct {
	StudentID     uint                       `json:"student_id"`
	Bill          *financedomain.StudentBill `json:"bill,omitempty"`
	Payment       *financedomain.Payment     `json:"payment,omitempty"`
	NotifType     string                     `json:"notif_type"`
	CustomMessage string                     `json:"custom_message,omitempty"`
	CustomReason  string                     `json:"custom_reason,omitempty"`
}

type FinanceNotificationService interface {
	Notify(job FinanceNotifyJob)
}

type financeNotificationService struct {
	db          *bun.DB
	stuRepo     academicrepo.StudentRepo
	userRepo    userauthrepo.UserRepo
	notiRepo    notificationrepo.NotificationRepo
	msg         utils.Messenger
	audit       auditusecase.AuditLogService
	asynqClient *asynq.Client
}

func NewFinanceNotificationService(db *bun.DB, stuRepo academicrepo.StudentRepo, userRepo userauthrepo.UserRepo, notiRepo notificationrepo.NotificationRepo, msg utils.Messenger, audit auditusecase.AuditLogService, asynqClient *asynq.Client) FinanceNotificationService {
	return &financeNotificationService{
		db:          db,
		stuRepo:     stuRepo,
		userRepo:    userRepo,
		notiRepo:    notiRepo,
		msg:         msg,
		audit:       audit,
		asynqClient: asynqClient,
	}
}

func (s *financeNotificationService) Notify(job FinanceNotifyJob) {
	payload, err := json.Marshal(job)
	if err != nil {
		fmt.Printf("[FinanceNotificationService] failed to marshal job: %v\n", err)
		return
	}
	task := asynq.NewTask(TaskFinanceNotification, payload, asynq.MaxRetry(3))
	if _, err := s.asynqClient.Enqueue(task); err != nil {
		fmt.Printf("[FinanceNotificationService] failed to enqueue task: %v\n", err)
	}
}
