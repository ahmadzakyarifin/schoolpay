package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"

	auditusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/audit/usecase"
	notificationusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/notification/usecase"
	"github.com/ahmadzakyarifin/schoolpay/internal/module/support/domain"
	"github.com/ahmadzakyarifin/schoolpay/internal/module/support/repository"
	userdomain "github.com/ahmadzakyarifin/schoolpay/internal/module/user_auth/domain"
	userrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/user_auth/repository"
	"github.com/ahmadzakyarifin/schoolpay/pkg/utils"
	"github.com/uptrace/bun"
)

type SupportService interface {
	RecordIncoming(ctx context.Context, phone, message string, parent *userdomain.User) (*domain.Conversation, error)
	List(ctx context.Context, status string, page, limit int) ([]domain.Conversation, int, error)
	Messages(ctx context.Context, conversationID uint) ([]domain.Message, error)
	Reply(ctx context.Context, conversationID, adminID uint, message string) error
	Assign(ctx context.Context, conversationID, adminID uint) error
	Close(ctx context.Context, conversationID uint) error
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

func (s *supportService) RecordIncoming(ctx context.Context, phone, message string, parent *userdomain.User) (*domain.Conversation, error) {
	phone = normalizePhone(phone)
	message = strings.TrimSpace(message)
	if phone == "" || message == "" {
		return nil, errors.New("phone dan message wajib diisi")
	}

	var result *domain.Conversation
	err := s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		conv, err := s.repo.FindOpenByPhone(ctx, phone)
		if err != nil || conv == nil {
			conv = &domain.Conversation{PhoneNumber: phone, Status: "open", Subject: utils.StringPtr("Pertanyaan Orang Tua")}
			if parent != nil {
				conv.ParentID = &parent.ID
				conv.ParentName = &parent.Name
			}
			if err := s.repo.CreateConversation(ctx, tx, conv); err != nil {
				return err
			}
		}
		msg := &domain.Message{ConversationID: conv.ID, SenderType: "parent", Message: message, DeliveryStatus: "received"}
		if parent != nil {
			msg.SenderID = &parent.ID
		}
		if err := s.repo.CreateMessage(ctx, tx, msg); err != nil {
			return err
		}
		if err := s.repo.UpdateConversationPreview(ctx, tx, conv.ID, message, 1); err != nil {
			return err
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
	return s.repo.FindAll(ctx, status, page, limit)
}

func (s *supportService) Messages(ctx context.Context, conversationID uint) ([]domain.Message, error) {
	_ = s.repo.MarkRead(ctx, s.db, conversationID)
	return s.repo.FindMessages(ctx, conversationID)
}

func (s *supportService) Reply(ctx context.Context, conversationID, adminID uint, message string) error {
	message = strings.TrimSpace(message)
	if message == "" {
		return errors.New("message wajib diisi")
	}

	conv, err := s.repo.FindByID(ctx, conversationID)
	if err != nil {
		return err
	}

	var msgID uint
	if err := s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		msg := &domain.Message{ConversationID: conversationID, SenderType: "admin", SenderID: &adminID, Message: message, DeliveryStatus: "pending"}
		if err := s.repo.CreateMessage(ctx, tx, msg); err != nil {
			return err
		}
		msgID = msg.ID
		if err := s.repo.UpdateConversationPreview(ctx, tx, conversationID, message, 0); err != nil {
			return err
		}
		if conv.AssignedAdminID == nil {
			_ = s.repo.Assign(ctx, tx, conversationID, adminID)
		}
		if s.audit != nil {
			userID, userName, role, ipAddress, userAgent := utils.GetAuditMeta(ctx)
			_ = s.audit.Log(ctx, tx, userID, userName, role, "REPLY_SUPPORT_CHAT", "support_conversations", conversationID, nil, map[string]interface{}{"message": message, "phone": conv.PhoneNumber, "delivery_status": "pending"}, ipAddress, userAgent)
		}
		return nil
	}); err != nil {
		return err
	}

	deliveryStatus := "sent"
	if err := s.wa.SendChatMessage(conv.PhoneNumber, message); err != nil {
		deliveryStatus = "failed"
		_, _ = s.db.NewUpdate().Model((*domain.Message)(nil)).Set("delivery_status = ?", deliveryStatus).Where("id = ?", msgID).Exec(ctx)
		return err
	}
	_, _ = s.db.NewUpdate().Model((*domain.Message)(nil)).Set("delivery_status = ?", deliveryStatus).Where("id = ?", msgID).Exec(ctx)
	return nil
}

func (s *supportService) Assign(ctx context.Context, conversationID, adminID uint) error {
	if adminID == 0 {
		return fmt.Errorf("admin_id tidak valid")
	}
	if err := s.repo.Assign(ctx, s.db, conversationID, adminID); err != nil {
		return err
	}
	if s.audit != nil {
		userID, userName, role, ipAddress, userAgent := utils.GetAuditMeta(ctx)
		_ = s.audit.Log(ctx, s.db, userID, userName, role, "ASSIGN_SUPPORT_CHAT", "support_conversations", conversationID, nil, map[string]interface{}{"assigned_admin_id": adminID}, ipAddress, userAgent)
	}
	return nil
}

func (s *supportService) Close(ctx context.Context, conversationID uint) error {
	if err := s.repo.Close(ctx, s.db, conversationID); err != nil {
		return err
	}
	if s.audit != nil {
		userID, userName, role, ipAddress, userAgent := utils.GetAuditMeta(ctx)
		_ = s.audit.Log(ctx, s.db, userID, userName, role, "CLOSE_SUPPORT_CHAT", "support_conversations", conversationID, nil, map[string]interface{}{"status": "closed"}, ipAddress, userAgent)
	}
	return nil
}
