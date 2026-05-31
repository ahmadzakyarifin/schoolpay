package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	academicrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/academic/repository"
	auditusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/audit/usecase"
	financedomain "github.com/ahmadzakyarifin/schoolpay/internal/module/finance/domain"
	notificationdomain "github.com/ahmadzakyarifin/schoolpay/internal/module/notification/domain"
	notificationrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/notification/repository"
	notificationusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/notification/usecase"
	userauthrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/user_auth/repository"
	authusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/user_auth/usecase"
	"github.com/ahmadzakyarifin/schoolpay/pkg/utils"
	"github.com/hibiken/asynq"
	"github.com/uptrace/bun"
)

// StartAsynqWorker initializes and starts the Asynq worker server
func StartAsynqWorker(srv *asynq.Server, db *bun.DB, stuRepo academicrepo.StudentRepo, userRepo userauthrepo.UserRepo, notiRepo notificationrepo.NotificationRepo, msg utils.Messenger, audit auditusecase.AuditLogService) {
	mux := asynq.NewServeMux()

	// Register handlers
	mux.HandleFunc(notificationusecase.TaskFinanceNotification, HandleFinanceNotificationTask(db, stuRepo, userRepo, notiRepo, msg, audit))
	mux.HandleFunc(authusecase.TaskAuthEmail, HandleAuthEmailTask(db, notiRepo, msg))

	log.Println("[Asynq] Worker server starting...")
	if err := srv.Run(mux); err != nil {
		log.Fatalf("[Asynq] could not run server: %v", err)
	}
}

func HandleAuthEmailTask(db *bun.DB, notiRepo notificationrepo.NotificationRepo, msg utils.Messenger) asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {
		var job authusecase.AuthEmailJob
		if err := json.Unmarshal(t.Payload(), &job); err != nil {
			return fmt.Errorf("json.Unmarshal failed: %v", err)
		}

		err := msg.SendEmail(job.Email, job.Subject, job.Body)
		status := "sent"
		var deliveryErr *string
		if err != nil {
			status = "failed"
			deliveryErr = utils.StringPtr(err.Error())
		}
		if job.UserID > 0 {
			_ = notiRepo.Create(ctx, db, &notificationdomain.Notification{
				UserID:         job.UserID,
				Title:          job.Subject,
				Message:        job.Body,
				Type:           "auth",
				Channel:        "email",
				DeliveryStatus: status,
				DeliveryError:  deliveryErr,
			})
		}
		if err != nil {
			return fmt.Errorf("SendEmail failed: %v", err)
		}
		return nil
	}
}

type financeNotificationText struct {
	Subject string
	Message string
}

func formatFinancePaymentSuccess(studentName string, payment *financedomain.Payment) financeNotificationText {
	if payment == nil {
		return financeNotificationText{Subject: "Konfirmasi Pembayaran SchoolPay", Message: "Pembayaran berhasil diterima oleh SchoolPay."}
	}
	lines := []string{
		"✅ *PEMBAYARAN BERHASIL*",
		"",
		fmt.Sprintf("Halo Ayah/Bunda dari *%s*,", studentName),
		"",
		"Pembayaran sudah kami terima dengan rincian:",
		"",
		fmt.Sprintf("• Total dialokasikan: *%s*", utils.FormatCurrency(payment.Amount)),
	}
	if payment.DepositApplied > 0 {
		lines = append(lines,
			fmt.Sprintf("• Saldo deposit dipakai: *%s*", utils.FormatCurrency(payment.DepositApplied)),
			fmt.Sprintf("• Dibayar via %s: *%s*", payment.Method, utils.FormatCurrency(payment.Amount-payment.DepositApplied)),
		)
	}
	lines = append(lines,
		fmt.Sprintf("• Metode: *%s*", payment.Method),
		fmt.Sprintf("• Referensi: *%s*", payment.TransactionRef),
		fmt.Sprintf("• Tanggal: *%s*", time.Now().Format("02-01-2006 15:04")),
		"",
		"Bukti pembayaran resmi dapat diunduh melalui portal orang tua.",
		"Terima kasih.",
	)
	return financeNotificationText{Subject: "Konfirmasi Pembayaran SchoolPay", Message: strings.Join(lines, "\n")}
}

func formatFinanceBillNotification(notifType, studentName, billTypeName, period, amount, dueDate, customReason, customMessage string) financeNotificationText {
	if customMessage != "" && (notifType == "initial" || notifType == "adjustment") {
		return financeNotificationText{Subject: financeSubject(notifType), Message: strings.TrimSpace(customMessage)}
	}
	switch notifType {
	case "initial":
		lines := []string{
			"📋 *TAGIHAN BARU*",
			"",
			fmt.Sprintf("Halo Ayah/Bunda dari *%s*,", studentName),
			"",
			"Tagihan baru telah diterbitkan:",
			"",
			fmt.Sprintf("• Jenis: *%s*", billTypeName),
			fmt.Sprintf("• Periode: *%s*", period),
			fmt.Sprintf("• Nominal: *%s*", amount),
			fmt.Sprintf("• Jatuh tempo: *%s*", dueDate),
		}
		if customReason != "" {
			lines = append(lines, "", "Keterangan:", customReason)
		}
		lines = append(lines, "", "Silakan lakukan pembayaran melalui portal orang tua.", "Terima kasih.")
		return financeNotificationText{Subject: financeSubject(notifType), Message: strings.Join(lines, "\n")}
	case "reminder":
		return financeNotificationText{Subject: financeSubject(notifType), Message: strings.Join([]string{
			"🔔 *PENGINGAT PEMBAYARAN*",
			"",
			fmt.Sprintf("Halo Ayah/Bunda dari *%s*,", studentName),
			"",
			"Tagihan berikut akan jatuh tempo dalam 3 hari:",
			"",
			fmt.Sprintf("• Jenis: *%s*", billTypeName),
			fmt.Sprintf("• Periode: *%s*", period),
			fmt.Sprintf("• Nominal: *%s*", amount),
			fmt.Sprintf("• Jatuh tempo: *%s*", dueDate),
			"",
			"Abaikan pesan ini jika pembayaran sudah dilakukan.",
			"Terima kasih.",
		}, "\n")}
	case "overdue":
		return financeNotificationText{Subject: financeSubject(notifType), Message: strings.Join([]string{
			"⚠️ *PERINGATAN TUNGGAKAN*",
			"",
			fmt.Sprintf("Halo Ayah/Bunda dari *%s*,", studentName),
			"",
			"Tagihan berikut sudah melewati jatuh tempo:",
			"",
			fmt.Sprintf("• Jenis: *%s*", billTypeName),
			fmt.Sprintf("• Periode: *%s*", period),
			fmt.Sprintf("• Nominal: *%s*", amount),
			fmt.Sprintf("• Jatuh tempo: *%s*", dueDate),
			"",
			"Pembayaran online untuk tagihan jatuh tempo ditutup.",
			"Mohon lakukan pembayaran langsung ke admin sekolah.",
		}, "\n")}
	case "adjustment":
		lines := []string{
			"📢 *PENYESUAIAN TAGIHAN*",
			"",
			fmt.Sprintf("Yth. Orang Tua/Wali dari *%s*,", studentName),
			"",
			"Ada penyesuaian tagihan dengan rincian:",
			"",
			fmt.Sprintf("• Jenis: *%s*", billTypeName),
			fmt.Sprintf("• Periode: *%s*", period),
			fmt.Sprintf("• Selisih kekurangan: *%s*", amount),
		}
		if customReason != "" {
			lines = append(lines, "", "Alasan penyesuaian:", customReason)
		}
		lines = append(lines, "", "Tagihan penyesuaian sudah tersedia di portal orang tua.", "Terima kasih atas perhatian dan kerja sama Bapak/Ibu.")
		return financeNotificationText{Subject: financeSubject(notifType), Message: strings.Join(lines, "\n")}
	case "refund_deposit":
		lines := []string{
			"💳 *REFUND SALDO DEPOSIT*",
			"",
			fmt.Sprintf("Halo Ayah/Bunda dari *%s*,", studentName),
			"",
		}
		if customMessage != "" {
			lines = append(lines, "Informasi refund:", "", strings.TrimSpace(customMessage))
		} else {
			lines = append(lines,
				"Ada dana pembayaran yang dialihkan ke saldo deposit.",
				"Silakan cek portal orang tua untuk melihat saldo terbaru.",
			)
		}
		lines = append(lines,
			"",
			"Terima kasih.",
		)
		return financeNotificationText{Subject: financeSubject(notifType), Message: strings.Join(lines, "\n")}
	default:
		return financeNotificationText{Subject: "Notifikasi SchoolPay", Message: "Ada informasi baru dari SchoolPay. Silakan cek portal orang tua."}
	}
}

func financeSubject(notifType string) string {
	switch notifType {
	case "initial":
		return "Tagihan Baru SchoolPay"
	case "reminder":
		return "Pengingat Pembayaran SchoolPay"
	case "overdue":
		return "Peringatan Tunggakan SchoolPay"
	case "adjustment":
		return "Penyesuaian Tarif Tagihan SchoolPay"
	case "refund_deposit":
		return "Refund Saldo Deposit SchoolPay"
	default:
		return "Notifikasi SchoolPay"
	}
}

func HandleFinanceNotificationTask(db *bun.DB, stuRepo academicrepo.StudentRepo, userRepo userauthrepo.UserRepo, notiRepo notificationrepo.NotificationRepo, msg utils.Messenger, audit auditusecase.AuditLogService) asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {
		var job notificationusecase.FinanceNotifyJob
		if err := json.Unmarshal(t.Payload(), &job); err != nil {
			return fmt.Errorf("json.Unmarshal failed: %v", err)
		}

		parents, err := userRepo.FindParentsByStudentID(ctx, job.StudentID)
		if err != nil {
			return fmt.Errorf("FindParentsByStudentID failed: %v", err)
		}
		student, err := stuRepo.FindByID(ctx, job.StudentID)
		if err != nil || student == nil {
			return fmt.Errorf("FindByID (Student) failed or nil: %v", err)
		}

		for _, parent := range parents {
			var message string
			var subject string

			switch job.NotifType {
			case "payment_success":
				text := formatFinancePaymentSuccess(student.Name, job.Payment)
				subject = text.Subject
				message = text.Message
			case "refund_deposit":
				text := formatFinanceBillNotification(job.NotifType, student.Name, "", "", "", "", "", job.CustomMessage)
				subject = text.Subject
				message = text.Message
			default:
				amountStr := utils.FormatCurrency(job.Bill.Amount)
				periodStr := "-"
				if job.Bill.Period != nil {
					periodStr = *job.Bill.Period
				}

				var billTypeName string
				if job.Bill.BillTypeName != "" {
					billTypeName = job.Bill.BillTypeName
				} else {
					var bt financedomain.BillType
					if err := db.NewSelect().Model(&bt).Where("id = ?", job.Bill.BillTypeID).Scan(ctx); err == nil {
						billTypeName = bt.Name
					}
				}

				text := formatFinanceBillNotification(job.NotifType, student.Name, billTypeName, periodStr, amountStr, job.Bill.DueDate.Format("02-01-2006"), job.CustomReason, job.CustomMessage)
				subject = text.Subject
				message = text.Message
			}

			if parent.PhoneNumber != "" {
				waID, waErr := msg.SendWhatsApp(parent.PhoneNumber, message)
				status := "sent"
				var deliveryErr *string
				if waErr != nil {
					status = "failed"
					deliveryErr = utils.StringPtr(waErr.Error())
				}
				if waID == "" {
					waID = fmt.Sprintf("local-wa-%d-%d", parent.ID, time.Now().UnixNano())
				}

				notiObj := &notificationdomain.Notification{
					UserID:         parent.ID,
					Title:          subject,
					Message:        message,
					DeliveryStatus: status,
					DeliveryError:  deliveryErr,
					Type:           "finance",
					Channel:        "whatsapp",
				}
				notiObj.WhatsappID = &waID
				_ = notiRepo.Create(ctx, db, notiObj)

				if audit != nil {
					_ = audit.Log(ctx, db, 0, "System/Automation", "system", "SEND_NOTIFICATION", "notifications", notiObj.ID, nil, map[string]interface{}{
						"title":       subject,
						"message":     message,
						"channel":     "WhatsApp",
						"whatsapp_id": waID,
						"status":      status,
						"error": func() string {
							if deliveryErr != nil {
								return *deliveryErr
							}
							return ""
						}(),
						"target_user_id":   parent.ID,
						"target_user_name": parent.Name,
					}, "127.0.0.1", "System Worker")
				}
			}

			if parent.Email != "" {
				mailErr := msg.SendEmail(parent.Email, subject, message)
				status := "sent"
				var deliveryErr *string
				if mailErr != nil {
					status = "failed"
					deliveryErr = utils.StringPtr(mailErr.Error())
				}

				notiObj := &notificationdomain.Notification{
					UserID:         parent.ID,
					Title:          subject,
					Message:        message,
					DeliveryStatus: status,
					DeliveryError:  deliveryErr,
					Type:           "finance",
					Channel:        "email",
				}
				_ = notiRepo.Create(ctx, db, notiObj)

				if audit != nil {
					_ = audit.Log(ctx, db, 0, "System/Automation", "system", "SEND_NOTIFICATION", "notifications", notiObj.ID, nil, map[string]interface{}{
						"title":   subject,
						"message": message,
						"channel": "Email",
						"status":  status,
						"error": func() string {
							if deliveryErr != nil {
								return *deliveryErr
							}
							return ""
						}(),
						"target_user_id":   parent.ID,
						"target_user_name": parent.Name,
					}, "127.0.0.1", "System Worker")
				}
			}
		}
		return nil
	}
}
