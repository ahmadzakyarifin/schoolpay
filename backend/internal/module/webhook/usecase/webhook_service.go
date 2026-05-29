package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/ahmadzakyarifin/schoolpay/config"
	academicrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/academic/repository"
	financerepo "github.com/ahmadzakyarifin/schoolpay/internal/module/finance/repository"
	notificationrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/notification/repository"
	notificationusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/notification/usecase"
	supportusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/support/usecase"
	userauthdomain "github.com/ahmadzakyarifin/schoolpay/internal/module/user_auth/domain"
	userauthrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/user_auth/repository"
	webhookrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/webhook/repository"
	"github.com/ahmadzakyarifin/schoolpay/internal/websocket"
	"github.com/ahmadzakyarifin/schoolpay/pkg/utils"
)

type WebhookService interface {
	HandleWAHAWebhook(ctx context.Context, payload json.RawMessage) error
}

type webhookService struct {
	repo     webhookrepo.WebhookRepo
	wa       notificationusecase.WhatsAppService
	notiRepo notificationrepo.NotificationRepo
	sbRepo   financerepo.StudentBillRepo
	payRepo  financerepo.PaymentRepo
	stuRepo  academicrepo.StudentRepo
	userRepo userauthrepo.UserRepo
	hub      *websocket.Hub
	support  supportusecase.SupportService
	cfg      *config.Config
}

func NewWebhookService(
	repo webhookrepo.WebhookRepo,
	wa notificationusecase.WhatsAppService,
	notiRepo notificationrepo.NotificationRepo,
	sbRepo financerepo.StudentBillRepo,
	payRepo financerepo.PaymentRepo,
	stuRepo academicrepo.StudentRepo,
	userRepo userauthrepo.UserRepo,
	hub *websocket.Hub,
	support supportusecase.SupportService,
	cfg *config.Config,
) WebhookService {
	return &webhookService{repo: repo, wa: wa, notiRepo: notiRepo, sbRepo: sbRepo, payRepo: payRepo, stuRepo: stuRepo, userRepo: userRepo, hub: hub, support: support, cfg: cfg}
}

func (s *webhookService) HandleWAHAWebhook(ctx context.Context, payload json.RawMessage) error {
	var data map[string]interface{}
	if err := json.Unmarshal(payload, &data); err != nil {
		return err
	}

	event, _ := data["event"].(string)
	if event == "session.status" {
		p, ok := data["payload"].(map[string]interface{})
		if ok {
			status, _ := p["status"].(string)

			// Broadcast status change to all clients
			if s.hub != nil {
				s.hub.Broadcast("WA_STATUS_CHANGED", map[string]interface{}{
					"status": status,
				})
			}

			if status == "STOPPED" {
				fmt.Println("[WA-WEBHOOK] Session STOPPED detected! Attempting auto-restart...")
				go s.wa.StartSession()
			}
		}
	}

	if event == "message" {
		p, ok := data["payload"].(map[string]interface{})
		if ok {
			s.handleIncomingMessage(ctx, p)
		}
	}

	if event == "message.ack" {
		p, ok := data["payload"].(map[string]interface{})
		if ok {
			var messageID string
			if idObj, isMap := p["id"].(map[string]interface{}); isMap {
				if serialized, hasSerial := idObj["_serialized"].(string); hasSerial {
					messageID = serialized
				}
			} else if idStr, isStr := p["id"].(string); isStr {
				messageID = idStr
			}

			if ackVal, ok := p["ack"].(float64); ok && messageID != "" {
				status := "SENT"
				switch int(ackVal) {
				case 2:
					status = "DELIVERED"
				case 3, 4:
					status = "READ"
				}
				fmt.Printf("[WA-WEBHOOK] Message ACK: %s -> %s\n", messageID, status)
				_ = s.notiRepo.UpdateStatusByWhatsappID(ctx, messageID, status)
			}
		}
	}

	return nil
}

func (s *webhookService) handleIncomingMessage(ctx context.Context, payload map[string]interface{}) {
	from, _ := payload["from"].(string)
	body, _ := payload["body"].(string)
	if from == "" || body == "" {
		return
	}

	cleanPhone := strings.Split(from, "@")[0]
	user, err := s.userRepo.FindByPhone(ctx, cleanPhone)
	if err != nil || user == nil {
		s.createSupportTicket(ctx, from, body, nil)
		return
	}

	cmd := strings.ToLower(strings.TrimSpace(body))
	switch {
	case cmd == "menu" || cmd == "halo" || cmd == "hi":
		s.sendMenu(from, user.Name)
	case strings.Contains(cmd, "tagihan"):
		s.handleCekTagihan(ctx, from, user)
	case strings.Contains(cmd, "tunggakan"):
		s.handleCekTunggakan(ctx, from, user)
	case strings.Contains(cmd, "pembayaran") || strings.Contains(cmd, "riwayat"):
		s.handleCekPembayaran(ctx, from, user)
	case strings.Contains(cmd, "cara bayar"):
		s.sendInstruction(from)
	case strings.Contains(cmd, "cs") || strings.Contains(cmd, "admin"):
		s.createSupportTicket(ctx, from, body, user)
	default:
		s.createSupportTicket(ctx, from, body, user)
	}
}

func (s *webhookService) sendMenu(to string, name string) {
	msg := fmt.Sprintf("Halo Bapak/Ibu *%s*,\nSelamat datang di *Layanan Bot SchoolPay* 🎓\n\nSilakan pilih menu:\n1. *Cek Tagihan*\n2. *Cek Tunggakan*\n3. *Cek Pembayaran*\n4. *Cara Bayar*\n5. *CS/Admin*", name)
	s.wa.SendChatMessage(to, msg)
}

func (s *webhookService) handleCekTagihan(ctx context.Context, to string, user *userauthdomain.User) {
	students, _ := s.stuRepo.GetStudentsByParentID(ctx, user.ID)
	if len(students) == 0 {
		s.wa.SendChatMessage(to, "Data siswa tidak ditemukan.")
		return
	}

	var total float64
	var msg strings.Builder
	msg.WriteString("📋 *DAFTAR TAGIHAN AKTIF*\n\n")

	for _, stu := range students {
		bills, _ := s.sbRepo.FindUnpaidBillsByStudent(ctx, stu.ID)
		if len(bills) > 0 {
			msg.WriteString(fmt.Sprintf("👤 *Siswa: %s*\n", stu.Name))
			for _, b := range bills {
				if b.DueDate.After(time.Now()) {
					period := "-"
					if b.Period != nil {
						period = *b.Period
					}
					msg.WriteString(fmt.Sprintf("- %s (%s): %s\n", b.BillTypeName, period, utils.FormatCurrency(b.Amount-b.TotalPaid)))
					total += (b.Amount - b.TotalPaid)
				}
			}
			msg.WriteString("\n")
		}
	}

	if total == 0 {
		s.wa.SendChatMessage(to, "Tidak ada tagihan aktif. Terima kasih! 😊")
		return
	}

	msg.WriteString(fmt.Sprintf("💰 *Total: %s*\n\n✅ Silakan bayar melalui portal parent.", utils.FormatCurrency(total)))
	s.wa.SendChatMessage(to, msg.String())
}

func (s *webhookService) handleCekTunggakan(ctx context.Context, to string, user *userauthdomain.User) {
	students, _ := s.stuRepo.GetStudentsByParentID(ctx, user.ID)
	var total float64
	var msg strings.Builder
	msg.WriteString("⚠️ *DAFTAR TUNGGAKAN*\n\n")

	for _, stu := range students {
		bills, _ := s.sbRepo.FindUnpaidBillsByStudent(ctx, stu.ID)
		hasOverdue := false
		for _, b := range bills {
			if b.DueDate.Before(time.Now()) {
				if !hasOverdue {
					msg.WriteString(fmt.Sprintf("👤 *Siswa: %s*\n", stu.Name))
					hasOverdue = true
				}
				period := "-"
				if b.Period != nil {
					period = *b.Period
				}
				msg.WriteString(fmt.Sprintf("- %s (%s): %s ❌\n", b.BillTypeName, period, utils.FormatCurrency(b.Amount-b.TotalPaid)))
				total += (b.Amount - b.TotalPaid)
			}
		}
		if hasOverdue {
			msg.WriteString("\n")
		}
	}

	if total == 0 {
		s.wa.SendChatMessage(to, "Tidak ada tunggakan. Terima kasih! ✨")
		return
	}

	msg.WriteString(fmt.Sprintf("💰 *Total: %s*\n\n🚫 *PENTING:* Pembayaran tunggakan melalui Admin Sekolah.", utils.FormatCurrency(total)))
	s.wa.SendChatMessage(to, msg.String())
}

func (s *webhookService) handleCekPembayaran(ctx context.Context, to string, user *userauthdomain.User) {
	students, _ := s.stuRepo.GetStudentsByParentID(ctx, user.ID)
	var msg strings.Builder
	msg.WriteString("🕒 *RIWAYAT PEMBAYARAN TERBARU*\n\n")

	for _, stu := range students {
		payments, _ := s.payRepo.FindByStudent(ctx, stu.ID)
		if len(payments) > 0 {
			msg.WriteString(fmt.Sprintf("👤 *Siswa: %s*\n", stu.Name))
			count := 0
			for _, p := range payments {
				if count >= 3 {
					break
				}
				date := p.CreatedAt.Format("02/01/2006")
				msg.WriteString(fmt.Sprintf("- %s: %s (%s) ✅\n", date, utils.FormatCurrency(p.Amount), p.Method))
				count++
			}
			msg.WriteString("\n")
		}
	}
	s.wa.SendChatMessage(to, msg.String())
}

func (s *webhookService) sendInstruction(to string) {
	msg := "💳 *INSTRUKSI PEMBAYARAN*\n1. *Portal Parent*\n2. *Transfer Bank ABC*\n3. *Scan QRIS*"
	s.wa.SendChatMessage(to, msg)
}

func (s *webhookService) createSupportTicket(ctx context.Context, to, body string, user *userauthdomain.User) {
	if s.support != nil {
		if conv, err := s.support.RecordIncoming(ctx, to, body, user); err == nil && s.hub != nil {
			s.hub.Broadcast("SUPPORT_CHAT_UPDATED", map[string]interface{}{"conversation_id": conv.ID, "phone": conv.PhoneNumber, "message": body})
		}
	}
	if user == nil {
		s.wa.SendChatMessage(to, "Nomor WhatsApp ini belum terdaftar sebagai wali siswa di SchoolPay. Pesan Anda tetap kami teruskan ke antrian CS/Admin agar dapat dibantu.")
		return
	}
	s.wa.SendChatMessage(to, "Pesan Anda sudah masuk ke antrian CS/Admin SchoolPay. Admin akan membalas dari dashboard sekolah, jadi Bapak/Ibu tidak perlu membuka WhatsApp Web.")
}
