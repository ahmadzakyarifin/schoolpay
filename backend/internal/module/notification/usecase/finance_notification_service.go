package usecase

import (
	"context"
	"encoding/json"
	"fmt"

	academicrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/academic/repository"
	auditusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/audit/usecase"
	financedomain "github.com/ahmadzakyarifin/schoolpay/internal/module/finance/domain"
	notificationdomain "github.com/ahmadzakyarifin/schoolpay/internal/module/notification/domain"
	notificationrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/notification/repository"
	userauthrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/user_auth/repository"
	"github.com/ahmadzakyarifin/schoolpay/pkg/utils"
	"github.com/uptrace/bun"
)

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
	db       *bun.DB
	stuRepo  academicrepo.StudentRepo
	userRepo userauthrepo.UserRepo
	notiRepo notificationrepo.NotificationRepo
	msg      utils.Messenger
	audit    auditusecase.AuditLogService
}

func NewFinanceNotificationService(db *bun.DB, stuRepo academicrepo.StudentRepo, userRepo userauthrepo.UserRepo, notiRepo notificationrepo.NotificationRepo, msg utils.Messenger, audit auditusecase.AuditLogService) FinanceNotificationService {
	return &financeNotificationService{
		db:       db,
		stuRepo:  stuRepo,
		userRepo: userRepo,
		notiRepo: notiRepo,
		msg:      msg,
		audit:    audit,
	}
}

func (s *financeNotificationService) Notify(job FinanceNotifyJob) {
	payload, err := json.Marshal(job)
	if err != nil {
		fmt.Printf("[FinanceNotificationService] failed to marshal job: %v\n", err)
		return
	}

	bj := &notificationdomain.BackgroundJob{
		TaskName: "notification:finance",
		Payload:  string(payload),
		Status:   "pending",
	}

	if _, err := s.db.NewInsert().Model(bj).Exec(context.Background()); err != nil {
		fmt.Printf("[FinanceNotificationService] failed to save job: %v\n", err)
	}
}
