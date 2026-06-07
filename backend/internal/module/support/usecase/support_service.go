package usecase

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	auditusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/audit/usecase"
	notificationusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/notification/usecase"
	"github.com/ahmadzakyarifin/schoolpay/internal/module/support/domain"
	"github.com/ahmadzakyarifin/schoolpay/internal/module/support/repository"
	userdomain "github.com/ahmadzakyarifin/schoolpay/internal/module/user_auth/domain"
	userrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/user_auth/repository"
	"github.com/ahmadzakyarifin/schoolpay/internal/helper"
	"github.com/uptrace/bun"
)

type SupportService interface {
	RecordIncoming(ctx context.Context, phone string, parent *userdomain.User) (*domain.Conversation, error)
	List(ctx context.Context, status string, page, limit int) ([]domain.Conversation, int, error)
	Assign(ctx context.Context, conversationID, adminID uint) error
	Close(ctx context.Context, conversationID uint) error
	UpdateStatus(ctx context.Context, conversationID uint, status string) error
	HasActiveConversation(ctx context.Context, phone string) (bool, error)
}

type supportService struct {
	db    *bun.DB
	repo  repository.SupportRepo
	users userrepo.UserRepo
	wa    notificationusecase.WhatsAppService
	audit auditusecase.AuditLogService
}

func NewSupportService(db *bun.DB, repo repository.SupportRepo, users userrepo.UserRepo, wa notificationusecase.WhatsAppService, audit auditusecase.AuditLogService) SupportService {
	return &supportService{db: db, repo: repo, users: users, wa: wa, audit: audit}
}

func normalizePhone(phone string) string {
	phone = strings.TrimSpace(strings.Split(phone, "@")[0])
	phone = strings.TrimLeft(phone, "+")
	if strings.HasPrefix(phone, "0") {
		phone = "62" + phone[1:]
	}
	return phone
}

func buildWhatsAppWebURL(phone string) string {
	phone = normalizePhone(phone)
	text := url.QueryEscape("Halo Bapak/Ibu, kami dari Admin SchoolPay. Ada yang bisa kami bantu?")
	return fmt.Sprintf("https://web.whatsapp.com/send?phone=%s&text=%s", phone, text)
}

func (s *supportService) RecordIncoming(ctx context.Context, phone string, parent *userdomain.User) (*domain.Conversation, error) {
	phone = normalizePhone(phone)
	if phone == "" {
		return nil, errors.New("phone wajib diisi")
	}

	var result *domain.Conversation
	err := s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		conv, err := s.repo.FindOpenByPhone(ctx, phone)
		if err != nil || conv == nil {
			conv = &domain.Conversation{PhoneNumber: phone, Status: "open"}
			if parent != nil {
				conv.ParentID = &parent.ID
				conv.ParentName = &parent.Name
			}
			if err := s.repo.CreateConversation(ctx, tx, conv); err != nil {
				return err
			}
		} else if parent != nil && (conv.ParentName == nil || *conv.ParentName == "") {
			conv.ParentName = &parent.Name
			_, _ = tx.NewUpdate().Model(conv).Column("parent_name", "parent_id").WherePK().Exec(ctx)
		}
		if s.audit != nil {
			userID := uint(0)
			userName := "WhatsApp Parent"
			role := "parent"
			if parent != nil {
				userID = parent.ID
				userName = parent.Name
				role = parent.Role
			}
			_ = s.audit.Log(ctx, tx, userID, userName, role, "RECORD_INCOMING_SUPPORT_CHAT", "support_conversations", conv.ID, nil, map[string]interface{}{"phone": phone}, "whatsapp", "waha-webhook")
		}
		result = conv
		return nil
	})
	return result, err
}

func (s *supportService) List(ctx context.Context, status string, page, limit int) ([]domain.Conversation, int, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	list, total, err := s.repo.FindAll(ctx, status, page, limit)
	if err != nil {
		return nil, 0, err
	}

	for i := range list {
		if list[i].ParentID == nil {
			var parent userdomain.User
			err := s.db.NewSelect().
				Model(&parent).
				Where("phone_number = ?", list[i].PhoneNumber).
				Limit(1).
				Scan(ctx)
			if err == nil {
				list[i].ParentID = &parent.ID
				list[i].ParentName = &parent.Name
				_, _ = s.db.NewUpdate().
					Model(&list[i]).
					Column("parent_id", "parent_name").
					WherePK().
					Exec(ctx)
			}
		}

		if list[i].ParentID != nil {
			var studentNames []string
			_ = s.db.NewSelect().
				Table("students").
				Column("name").
				Where("parent_id = ?", *list[i].ParentID).
				Scan(ctx, &studentNames)
			if len(studentNames) > 0 {
				list[i].StudentNames = strings.Join(studentNames, ", ")
			}
		}
		list[i].WhatsAppWebURL = buildWhatsAppWebURL(list[i].PhoneNumber)
	}
	return list, total, nil
}

func (s *supportService) Assign(ctx context.Context, conversationID, adminID uint) error {
	if err := s.repo.Assign(ctx, s.db, conversationID, adminID); err != nil {
		return err
	}
	if s.audit != nil {
		userID, userName, role, ipAddress, userAgent := helper.GetAuditMeta(ctx)
		_ = s.audit.Log(ctx, s.db, userID, userName, role, "ASSIGN_SUPPORT_CHAT", "support_conversations", conversationID, nil, map[string]interface{}{"assigned_admin_id": adminID}, ipAddress, userAgent)
	}
	return nil
}

func (s *supportService) Close(ctx context.Context, conversationID uint) error {
	if err := s.repo.Close(ctx, s.db, conversationID); err != nil {
		return err
	}
	if s.audit != nil {
		userID, userName, role, ipAddress, userAgent := helper.GetAuditMeta(ctx)
		_ = s.audit.Log(ctx, s.db, userID, userName, role, "CLOSE_SUPPORT_CHAT", "support_conversations", conversationID, nil, map[string]interface{}{"status": "closed"}, ipAddress, userAgent)
	}
	return nil
}

func (s *supportService) UpdateStatus(ctx context.Context, conversationID uint, status string) error {
	status = strings.ToLower(strings.TrimSpace(status))
	switch status {
	case "open", "pending", "closed":
	default:
		return errors.New("status tiket CS tidak valid")
	}
	if err := s.repo.UpdateStatus(ctx, s.db, conversationID, status); err != nil {
		return err
	}
	if s.audit != nil {
		userID, userName, role, ipAddress, userAgent := helper.GetAuditMeta(ctx)
		_ = s.audit.Log(ctx, s.db, userID, userName, role, "UPDATE_SUPPORT_CHAT_STATUS", "support_conversations", conversationID, nil, map[string]interface{}{"status": status}, ipAddress, userAgent)
	}
	return nil
}

func (s *supportService) HasActiveConversation(ctx context.Context, phone string) (bool, error) {
	phone = normalizePhone(phone)
	conv, err := s.repo.FindOpenByPhone(ctx, phone)
	if err != nil {
		return false, err
	}
	return conv != nil, nil
}
