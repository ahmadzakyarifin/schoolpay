package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
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
	mux.HandleFunc(authusecase.TaskAuthEmail, HandleAuthEmailTask(msg))

	log.Println("[Asynq] Worker server starting...")
	if err := srv.Run(mux); err != nil {
		log.Fatalf("[Asynq] could not run server: %v", err)
	}
}

func HandleAuthEmailTask(msg utils.Messenger) asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {
		var job authusecase.AuthEmailJob
		if err := json.Unmarshal(t.Payload(), &job); err != nil {
			return fmt.Errorf("json.Unmarshal failed: %v", err)
		}

		err := msg.SendEmail(job.Email, job.Subject, job.Body)
		if err != nil {
			return fmt.Errorf("SendEmail failed: %v", err)
		}
		return nil
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
				subject = "Konfirmasi Pembayaran SchoolPay"
				depositLine := ""
				if job.Payment.DepositApplied > 0 {
					depositLine = fmt.Sprintf("🔹 *Saldo Deposit Dipakai:* %s\n🔹 *Dibayar via %s:* %s\n", utils.FormatCurrency(job.Payment.DepositApplied), job.Payment.Method, utils.FormatCurrency(job.Payment.Amount-job.Payment.DepositApplied))
				}
				message = fmt.Sprintf(
					"✅ *PEMBAYARAN BERHASIL*\n\nHalo Ayah/Bunda dari *%s*,\n\n"+
						"Terima kasih, pembayaran Anda telah kami terima:\n"+
						"🔹 *Total Dialokasikan:* %s\n"+
						"%s"+
						"🔹 *Metode:* %s\n"+
						"🔹 *Referensi:* %s\n"+
						"🔹 *Tanggal:* %s\n\n"+
						"Bukti pembayaran resmi dapat Anda unduh melalui portal orang tua. Terima kasih telah melakukan pembayaran tepat waktu! ✨",
					student.Name, utils.FormatCurrency(job.Payment.Amount), depositLine, job.Payment.Method, job.Payment.TransactionRef, time.Now().Format("02-01-2006 15:04"),
				)
			case "refund_deposit":
				subject = "Refund Saldo Deposit SchoolPay"
				message = job.CustomMessage
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

				switch job.NotifType {
				case "initial":
					subject = "Tagihan Baru SchoolPay"
					if job.CustomMessage != "" {
						message = job.CustomMessage
					} else {
						message = fmt.Sprintf(
							"📋 *TAGIHAN BARU - SCHOOLPAY*\n\nHalo Ayah/Bunda dari *%s*,\n\n"+
								"Tagihan baru telah diterbitkan:\n"+
								"🔹 *Jenis:* %s\n"+
								"🔹 *Periode:* %s\n"+
								"🔹 *Nominal:* %s\n"+
								"🔹 *Jatuh Tempo:* %s\n\n",
							student.Name, billTypeName, periodStr, amountStr, job.Bill.DueDate.Format("02-01-2006"),
						)
						if job.CustomReason != "" {
							message += fmt.Sprintf("*Keterangan Tambahan:*\n_%s_\n\n", job.CustomReason)
						}
						message += "Silakan lakukan pembayaran melalui portal orang tua. Terima kasih."
					}
				case "reminder":
					subject = "Pengingat Pembayaran SchoolPay"
					message = fmt.Sprintf(
						"🔔 *PENGINGAT PEMBAYARAN*\n\nHalo Ayah/Bunda dari *%s*,\n\n"+
							"Kami menginformasikan bahwa tagihan *%s* (%s) sebesar *%s* akan jatuh tempo dalam 3 hari lagi.\n\n"+
							"Abaikan pesan ini jika Anda sudah melakukan pembayaran. Terima kasih.",
						student.Name, billTypeName, periodStr, amountStr,
					)
				case "overdue":
					subject = "Peringatan Tunggakan SchoolPay"
					message = fmt.Sprintf(
						"⚠️ *PERINGATAN TUNGGAKAN*\n\nHalo Ayah/Bunda dari *%s*,\n\n"+
							"Tagihan *%s* (%s) sebesar *%s* telah MELEWATI batas jatuh tempo pada tanggal %s.\n\n"+
							"Pembayaran online untuk tagihan jatuh tempo ditutup. Mohon lakukan pembayaran langsung ke admin sekolah. Terima kasih.",
						student.Name, billTypeName, periodStr, amountStr, job.Bill.DueDate.Format("02-01-2006"),
					)
				case "adjustment":
					subject = "Penyesuaian Tarif Tagihan SchoolPay"
					if job.CustomMessage != "" {
						message = job.CustomMessage
					} else {
						message = fmt.Sprintf(
							"📢 *PEMBERITAHUAN RESMI SCHOOLPAY*\n*Penyesuaian Tarif Tagihan Sekolah*\n\n"+
								"Yth. Orang Tua / Wali dari *%s*,\n\n"+
								"Melalui pesan ini, kami menginformasikan adanya penyesuaian kebijakan tarif untuk tagihan *%s* (%s) dengan rincian selisih penyesuaian sebesar:\n"+
								"▫️ *Selisih Kekurangan:* *%s*\n\n",
							student.Name, billTypeName, periodStr, amountStr,
						)
						if job.CustomReason != "" {
							message += fmt.Sprintf("*Mengapa ada penyesuaian ini?*\n_%s_\n\n", job.CustomReason)
						}
						message += "Sistem kami telah menerbitkan tagihan penyesuaian pada akun SchoolPay Anda. Pembayaran dapat dilakukan melalui portal orang tua. Terima kasih atas perhatian dan kerja sama Bapak/Ibu."
					}
				}
			}

			if parent.PhoneNumber != "" {
				waID, waErr := msg.SendWhatsApp(parent.PhoneNumber, message)
				status := "sent"
				var deliveryErr *string
				if waErr != nil {
					status = "failed"
					deliveryErr = utils.StringPtr(waErr.Error())
				}

				notiObj := &notificationdomain.Notification{
					UserID:         parent.ID,
					Title:          subject,
					Message:        message,
					DeliveryStatus: status,
					DeliveryError:  deliveryErr,
					Type:           "finance",
				}
				if waID != "" {
					notiObj.WhatsappID = &waID
				}
				_ = notiRepo.Create(ctx, db, notiObj)

				if audit != nil {
					_ = audit.Log(ctx, db, 0, "System/Automation", "system", "SEND_NOTIFICATION", "notifications", notiObj.ID, nil, map[string]interface{}{
						"title":            subject,
						"message":          message,
						"channel":          "WhatsApp",
						"whatsapp_id":      waID,
						"status":           status,
						"error":            func() string { if deliveryErr != nil { return *deliveryErr }; return "" }(),
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
				}
				_ = notiRepo.Create(ctx, db, notiObj)

				if audit != nil {
					_ = audit.Log(ctx, db, 0, "System/Automation", "system", "SEND_NOTIFICATION", "notifications", notiObj.ID, nil, map[string]interface{}{
						"title":            subject,
						"message":          message,
						"channel":          "Email",
						"status":           status,
						"error":            func() string { if deliveryErr != nil { return *deliveryErr }; return "" }(),
						"target_user_id":   parent.ID,
						"target_user_name": parent.Name,
					}, "127.0.0.1", "System Worker")
				}
			}
		}
		return nil
	}
}
