package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ahmadzakyarifin/schoolpay/config"
	academicdomain "github.com/ahmadzakyarifin/schoolpay/internal/module/academic/domain"
	academicrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/academic/repository"
	auditusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/audit/usecase"
	"github.com/ahmadzakyarifin/schoolpay/internal/module/finance/domain"
	"github.com/ahmadzakyarifin/schoolpay/internal/module/finance/repository"
	notificationusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/notification/usecase"
	"github.com/ahmadzakyarifin/schoolpay/internal/websocket"
	"github.com/ahmadzakyarifin/schoolpay/pkg/utils"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"github.com/uptrace/bun"
)

type PaymentService interface {
	Process(ctx context.Context, studentID uint, p *domain.Payment, billIDs []uint) error
	CreateIntent(ctx context.Context, studentID uint, amount float64, depositApplied float64, billIDs []uint, isBypassRule bool) (*domain.Payment, error)
	HandleWebhook(ctx context.Context, payload map[string]interface{}) error
	GetReceipt(ctx context.Context, paymentID uint) (*domain.Receipt, error)
	GetHistory(ctx context.Context, studentID uint) ([]domain.Payment, error)
}

type paymentService struct {
	db       *bun.DB
	repo     repository.PaymentRepo
	sbRepo   repository.StudentBillRepo
	stuRepo  academicrepo.StudentRepo
	notifSvc notificationusecase.FinanceNotificationService
	cfg      *config.Config
	hub      *websocket.Hub
	audit    auditusecase.AuditLogService
}

func NewPaymentService(db *bun.DB, repo repository.PaymentRepo, sbRepo repository.StudentBillRepo, stuRepo academicrepo.StudentRepo, notifSvc notificationusecase.FinanceNotificationService, cfg *config.Config, hub *websocket.Hub, audit auditusecase.AuditLogService) PaymentService {
	svc := &paymentService{db, repo, sbRepo, stuRepo, notifSvc, cfg, hub, audit}
	svc.startReconciler()
	return svc
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func (s *paymentService) CreateIntent(ctx context.Context, studentID uint, amount float64, depositApplied float64, billIDs []uint, isBypassRule bool) (*domain.Payment, error) {
	if isBypassRule {
		return nil, fmt.Errorf("gagal: bypass aturan tidak tersedia untuk pembayaran online")
	}
	if depositApplied < 0 {
		return nil, fmt.Errorf("gagal: saldo yang digunakan tidak valid")
	}
	if depositApplied >= amount {
		return nil, fmt.Errorf("gagal: pembayaran Midtrans membutuhkan sisa nominal setelah potongan saldo")
	}
	role := ""
	if s.audit != nil {
		userID, _, auditRole, _, _ := utils.GetAuditMeta(ctx)
		role = auditRole
		if role == "parent" {
			student, err := s.stuRepo.FindByID(ctx, studentID)
			if err != nil || student == nil {
				return nil, fmt.Errorf("siswa tidak ditemukan")
			}
			if student.ParentID != userID {
				return nil, fmt.Errorf("unauthorized: Anda tidak memiliki akses ke data siswa ini")
			}
		}
	}

	if len(billIDs) == 0 {
		unpaidBills, _ := s.sbRepo.FindUnpaidBillsByStudent(ctx, studentID)
		for _, ub := range unpaidBills {
			billIDs = append(billIDs, ub.ID)
		}
	}

	if len(billIDs) == 0 {
		return nil, fmt.Errorf("gagal: tidak ada tagihan yang bisa dibayar")
	}

	student, err := s.stuRepo.FindByID(ctx, studentID)
	if err != nil || student == nil {
		return nil, fmt.Errorf("siswa tidak ditemukan")
	}

	today := time.Now().Truncate(24 * time.Hour)
	targetOutstanding := 0.0
	for _, billID := range billIDs {
		bill, err := s.sbRepo.FindByID(ctx, billID)
		if err != nil {
			return nil, err
		}
		if bill.StudentID != studentID {
			return nil, fmt.Errorf("gagal: tagihan tidak sesuai dengan siswa")
		}
		if bill.Status == "voided" {
			return nil, fmt.Errorf("gagal: tagihan yang sudah ditarik tidak bisa dibayar lewat Midtrans")
		}
		if role == "parent" && bill.DueDate.Before(today) {
			return nil, fmt.Errorf("gagal: tagihan %s sudah jatuh tempo. Pembayaran online ditutup, silakan bayar langsung ke admin sekolah", bill.BillTypeName)
		}
		if bill.Status != "paid" {
			targetOutstanding += bill.Amount - bill.TotalPaid
		}
	}
	if targetOutstanding <= 0 {
		return nil, fmt.Errorf("gagal: tagihan yang dipilih sudah lunas")
	}
	if amount > targetOutstanding {
		return nil, fmt.Errorf("gagal: nominal pembayaran melebihi total sisa tagihan yang dipilih")
	}
	if depositApplied > targetOutstanding {
		return nil, fmt.Errorf("gagal: saldo yang digunakan melebihi total sisa tagihan yang dipilih")
	}
	if depositApplied > student.DepositBalance {
		return nil, fmt.Errorf("gagal: Saldo deposit siswa tidak mencukupi (Saldo: %s, Digunakan: %s)", utils.FormatCurrency(student.DepositBalance), utils.FormatCurrency(depositApplied))
	}
	if err := s.validateAllocationRules(ctx, amount, billIDs); err != nil {
		return nil, err
	}

	intentBillIDs, _ := json.Marshal(billIDs)
	transactionRef := fmt.Sprintf("PAY-%d-%d", studentID, time.Now().UnixNano())
	gatewayAmount := amount - depositApplied
	p := &domain.Payment{
		StudentID: studentID, Amount: amount, DepositApplied: depositApplied, Channel: "gateway", Method: "midtrans",
		Status: "pending", GatewayProvider: "midtrans",
		TransactionRef:  transactionRef,
		ExternalOrderID: utils.StringPtr(transactionRef),
		IdempotencyKey:  utils.StringPtr(transactionRef),
		IntentBillIDs:   utils.StringPtr(string(intentBillIDs)),
		IsBypassRule:    isBypassRule,
		CreatedBy:       "SYSTEM",
	}

	if err := s.repo.Create(ctx, s.db, p); err != nil {
		return nil, err
	}

	midtrans.ServerKey = s.cfg.MidtransServerKey
	midtrans.Environment = midtrans.Sandbox
	if !s.cfg.MidtransIsSandbox {
		midtrans.Environment = midtrans.Production
	}

	studentName := "Siswa"
	if student != nil {
		studentName = student.Name
	}

	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{OrderID: p.TransactionRef, GrossAmt: int64(gatewayAmount)},
		CustomerDetail:     &midtrans.CustomerDetails{FName: studentName},
		Items:              &[]midtrans.ItemDetails{{ID: "BILL-PAY", Name: fmt.Sprintf("Tagihan: %s", studentName), Price: int64(gatewayAmount), Qty: 1}},
	}

	snapResp, err := snap.CreateTransaction(req)
	if err != nil {
		return nil, err
	}
	p.PaymentURL = snapResp.RedirectURL
	p.SnapToken = snapResp.Token

	if s.audit != nil {
		userID, userName, role, ipAddress, userAgent := utils.GetAuditMeta(ctx)
		newVals := map[string]interface{}{"student_id": studentID, "amount": amount, "deposit_applied": depositApplied, "gateway_amount": gatewayAmount, "transaction_ref": p.TransactionRef, "status": p.Status}
		_ = s.audit.Log(ctx, s.db, userID, userName, role, "CREATE_PAYMENT_INTENT", "payments", p.ID, nil, newVals, ipAddress, userAgent)
	}

	return p, nil
}

func (s *paymentService) HandleWebhook(ctx context.Context, payload map[string]interface{}) error {
	orderID, _ := payload["order_id"].(string)
	status, _ := payload["transaction_status"].(string)
	paymentType, _ := payload["payment_type"].(string)

	if status == "expire" || status == "cancel" || status == "deny" {
		p, err := s.repo.FindByRef(ctx, orderID)
		if err != nil {
			return err
		}
		if p.Status == "success" || p.Status == "failed" {
			return nil
		}

		raw, _ := json.Marshal(payload)
		now := time.Now()
		p.Status = "failed"
		p.Note = utils.StringPtr(fmt.Sprintf("Payment %s", status))
		p.GatewayRawResponse = utils.StringPtr(string(raw))
		p.LastCheckedAt = &now
		_, _ = s.db.NewUpdate().Model(p).Column("status", "note", "gateway_raw_response", "last_checked_at", "updated_at").WherePK().Exec(ctx)
		return nil
	}

	if status == "capture" || status == "settlement" {
		p, err := s.repo.FindByRef(ctx, orderID)
		if err != nil {
			return err
		}
		if p.Status == "success" {
			return nil
		}

		// Webhook gross amount validation to prevent client-side pricing exploits
		var webhookAmount float64
		if amtStr, ok := payload["gross_amount"].(string); ok {
			webhookAmount, _ = strconv.ParseFloat(amtStr, 64)
		} else if amtNum, ok := payload["gross_amount"].(float64); ok {
			webhookAmount = amtNum
		}

		raw, _ := json.Marshal(payload)
		now := time.Now()
		p.GatewayRawResponse = utils.StringPtr(string(raw))
		p.LastCheckedAt = &now

		expectedGatewayAmount := p.Amount - p.DepositApplied
		if webhookAmount != expectedGatewayAmount {
			p.Status = "failed"
			p.Note = utils.StringPtr(fmt.Sprintf("DITOLAK: Nominal webhook (Rp %.2f) berbeda dengan nominal gateway lokal (Rp %.2f)", webhookAmount, expectedGatewayAmount))
			_, _ = s.db.NewUpdate().Model(p).Column("status", "note", "gateway_raw_response", "last_checked_at", "updated_at").WherePK().Exec(ctx)
			return fmt.Errorf("gagal: nominal webhook tidak cocok dengan data lokal")
		}

		p.Status = "success"
		p.Method = paymentType
		p.PaidAt = &now
		p.ReconciledAt = &now
		if err := s.Process(ctx, p.StudentID, p, nil); err != nil {
			if s.isRecoverableGatewayAllocationError(err) {
				return s.moveGatewayPaymentToDeposit(ctx, p, paymentType, string(raw), now, err)
			}
			return err
		}
		return nil
	}
	return nil
}

func (s *paymentService) validateAllocationRules(ctx context.Context, amount float64, billIDs []uint) error {
	remainingAmount := amount
	for _, billID := range billIDs {
		if remainingAmount <= 0 {
			break
		}
		bill, err := s.sbRepo.FindByID(ctx, billID)
		if err != nil {
			return err
		}
		if bill.Status == "paid" || bill.Status == "voided" {
			continue
		}
		needed := bill.Amount - bill.TotalPaid
		if needed <= 0 {
			continue
		}
		if remainingAmount >= needed {
			remainingAmount -= needed
			continue
		}
		if !bill.AllowInstallment {
			return fmt.Errorf("gagal: Tagihan '%s' tidak melayani pembayaran sebagian/cicilan. Wajib dibayar lunas sebesar %s", bill.BillTypeName, utils.FormatCurrency(needed))
		}
		if bill.MaxInstallment > 0 {
			existingInstallments, err := s.db.NewSelect().Model((*domain.PaymentDetail)(nil)).
				Where("student_bill_id = ?", bill.ID).
				Count(ctx)
			if err != nil {
				return err
			}
			if existingInstallments >= bill.MaxInstallment {
				return fmt.Errorf("gagal: Tagihan '%s' sudah mencapai batas maksimal %d kali cicilan", bill.BillTypeName, bill.MaxInstallment)
			}
		}
	}
	return nil
}

func (s *paymentService) isRecoverableGatewayAllocationError(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "tidak melayani pembayaran sebagian/cicilan") ||
		strings.Contains(msg, "sudah mencapai batas maksimal")
}

func (s *paymentService) moveGatewayPaymentToDeposit(ctx context.Context, p *domain.Payment, method string, rawResponse string, settledAt time.Time, allocationErr error) error {
	return s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		var locked domain.Payment
		if err := tx.NewSelect().Model(&locked).Where("id = ?", p.ID).For("UPDATE").Scan(ctx); err != nil {
			return err
		}
		if locked.Status == "success" {
			return nil
		}

		gatewayAmount := locked.Amount - locked.DepositApplied
		if gatewayAmount <= 0 {
			return allocationErr
		}

		if _, err := tx.NewUpdate().Model((*academicdomain.Student)(nil)).
			Set("deposit_balance = deposit_balance + ?", gatewayAmount).
			Where("id = ?", locked.StudentID).
			Exec(ctx); err != nil {
			return err
		}

		dm := &domain.DepositMovement{
			StudentID:   locked.StudentID,
			Type:        "IN",
			Amount:      gatewayAmount,
			Reason:      "GATEWAY_ALLOCATION_FAILED",
			ReferenceID: utils.StringPtr(fmt.Sprintf("%d", locked.ID)),
			CreatedBy:   "SYSTEM",
		}
		if _, err := tx.NewInsert().Model(dm).Exec(ctx); err != nil {
			return err
		}

		note := fmt.Sprintf("Dana Midtrans dialihkan ke saldo deposit karena alokasi tagihan gagal: %s", allocationErr.Error())
		locked.Amount = gatewayAmount
		locked.DepositApplied = 0
		locked.Status = "success"
		locked.Method = method
		locked.PaidAt = &settledAt
		locked.ReconciledAt = &settledAt
		locked.LastCheckedAt = &settledAt
		locked.Note = &note
		locked.GatewayRawResponse = &rawResponse
		if _, err := tx.NewUpdate().Model(&locked).
			Column("amount", "deposit_applied", "status", "method", "paid_at", "reconciled_at", "last_checked_at", "note", "gateway_raw_response", "updated_at").
			WherePK().
			Exec(ctx); err != nil {
			return err
		}

		if s.audit != nil {
			userID, userName, role, ipAddress, userAgent := utils.GetAuditMeta(ctx)
			newVals := map[string]interface{}{
				"student_id": locked.StudentID, "amount": gatewayAmount, "transaction_ref": locked.TransactionRef,
				"deposit_reason": "GATEWAY_ALLOCATION_FAILED", "allocation_error": allocationErr.Error(),
			}
			_ = s.audit.Log(ctx, tx, userID, userName, role, "GATEWAY_PAYMENT_TO_DEPOSIT", "payments", locked.ID, nil, newVals, ipAddress, userAgent)
		}

		if s.notifSvc != nil {
			paymentCopy := locked
			s.notifSvc.Notify(notificationusecase.FinanceNotifyJob{StudentID: locked.StudentID, Payment: &paymentCopy, NotifType: "payment_success"})
		}
		return nil
	})
}

func (s *paymentService) GetReceipt(ctx context.Context, paymentID uint) (*domain.Receipt, error) {
	p, err := s.repo.FindByID(ctx, paymentID)
	if err != nil {
		return nil, err
	}

	role := ""
	if s.audit != nil {
		userID, _, auditRole, _, _ := utils.GetAuditMeta(ctx)
		role = auditRole
		if role == "parent" {
			student, err := s.stuRepo.FindByID(ctx, p.StudentID)
			if err != nil || student == nil {
				return nil, fmt.Errorf("siswa tidak ditemukan")
			}
			if student.ParentID != userID {
				return nil, fmt.Errorf("unauthorized: Anda tidak memiliki akses ke data siswa ini")
			}
		}
	}

	details, _ := s.repo.FindDetailsByPaymentID(ctx, paymentID)
	student, _ := s.stuRepo.FindByID(ctx, p.StudentID)
	studentName := "Siswa"
	nis := "-"
	className := "-"
	if student != nil {
		studentName = student.Name
		if student.NIS != "" {
			nis = student.NIS
		}
		if student.ClassName != nil && *student.ClassName != "" {
			className = *student.ClassName
		}
	}

	receipt := &domain.Receipt{
		PaymentID: p.ID, ReceiptNumber: p.TransactionRef, StudentName: studentName,
		NIS: nis, ClassName: className, Amount: p.Amount, PaymentMethod: p.Method, PaidAt: time.Now(), Items: []domain.ReceiptItem{},
	}
	if p.PaidAt != nil {
		receipt.PaidAt = *p.PaidAt
	}
	for _, d := range details {
		receipt.Items = append(receipt.Items, domain.ReceiptItem{BillName: d.BillTypeName, Period: d.Period, Amount: d.Amount})
	}
	if len(receipt.Items) == 0 {
		receipt.Items = append(receipt.Items, domain.ReceiptItem{
			BillName: "Pembayaran Tagihan Keuangan",
			Period:   "Sekali Bayar / Alokasi Otomatis",
			Amount:   p.Amount,
		})
	}
	return receipt, nil
}

func (s *paymentService) Process(ctx context.Context, studentID uint, p *domain.Payment, billIDs []uint) error {
	return s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		adminID, adminName, role, _, _ := utils.GetAuditMeta(ctx)
		createdBy := "SYSTEM"
		if adminName != "" {
			createdBy = fmt.Sprintf("ADMIN-%d (%s)", adminID, adminName)
		}
		if p.CreatedBy == "" {
			p.CreatedBy = createdBy
		}

		if p.Method == "DEPOSIT" {
			p.DepositApplied = p.Amount
			p.Method = "Saldo Deposit"
			p.Channel = "deposit"
		}
		if p.DepositApplied < 0 || p.DepositApplied > p.Amount {
			return fmt.Errorf("gagal: saldo yang digunakan tidak valid")
		}
		if p.IsBypassRule && strings.TrimSpace(stringValue(p.BypassReason)) == "" {
			return fmt.Errorf("gagal: alasan bypass wajib diisi untuk pembayaran sebagian di luar aturan")
		}

		if role == "parent" {
			if p.IsBypassRule {
				return fmt.Errorf("unauthorized: bypass aturan pembayaran hanya boleh dilakukan admin")
			}
			if p.Channel != "gateway" && p.Method != "Saldo Deposit" {
				return fmt.Errorf("unauthorized: orang tua hanya boleh membayar lewat Midtrans atau memakai potongan saldo deposit")
			}
			student, err := s.stuRepo.FindByID(ctx, studentID)
			if err != nil || student == nil {
				return fmt.Errorf("siswa tidak ditemukan")
			}
			if student.ParentID != adminID {
				return fmt.Errorf("unauthorized: Anda tidak memiliki akses ke data siswa ini")
			}
		}

		isExisting := p.ID != 0
		if p.ID == 0 {
			p.StudentID = studentID
			p.Status = "success"
			if p.TransactionRef == "" {
				p.TransactionRef = fmt.Sprintf("CSH-%d-%d", studentID, time.Now().Unix())
			}
			if p.PaidAt == nil {
				now := time.Now()
				p.PaidAt = &now
			}
			if err := s.repo.Create(ctx, tx, p); err != nil {
				return err
			}
		} else {
			var existingP domain.Payment
			if err := tx.NewSelect().Model(&existingP).Where("id = ?", p.ID).For("UPDATE").Scan(ctx); err != nil {
				return err
			}
			if existingP.Status == "success" {
				return nil
			}

			p.Status = "success"
			if _, err := tx.NewUpdate().Model(p).Column("status", "method", "paid_at", "created_by", "note", "deposit_applied", "gateway_raw_response", "last_checked_at", "reconciled_at", "updated_at").WherePK().Exec(ctx); err != nil {
				return err
			}
		}

		if len(billIDs) == 0 && p.IntentBillIDs != nil && *p.IntentBillIDs != "" {
			_ = json.Unmarshal([]byte(*p.IntentBillIDs), &billIDs)
		}
		if len(billIDs) == 0 {
			details, _ := s.repo.FindDetailsByPaymentID(ctx, p.ID)
			for _, d := range details {
				billIDs = append(billIDs, d.StudentBillID)
			}
			if len(billIDs) == 0 {
				unpaidBills, _ := s.sbRepo.FindUnpaidBillsByStudent(ctx, studentID)
				for _, ub := range unpaidBills {
					billIDs = append(billIDs, ub.ID)
				}
			}
		}

		if p.DepositApplied > 0 {
			targetOutstanding := 0.0
			for _, bID := range billIDs {
				var b domain.StudentBill
				if err := tx.NewSelect().Model(&b).Where("id = ?", bID).Scan(ctx); err != nil {
					return err
				}
				if b.Status != "paid" && b.Status != "voided" {
					targetOutstanding += b.Amount - b.TotalPaid
				}
			}
			if targetOutstanding > 0 && p.DepositApplied > targetOutstanding {
				return fmt.Errorf("gagal: saldo yang digunakan melebihi total tagihan yang dipilih")
			}

			var student academicdomain.Student
			if err := tx.NewSelect().Model(&student).Where("id = ?", studentID).For("UPDATE").Scan(ctx); err != nil {
				return err
			}
			if student.DepositBalance < p.DepositApplied {
				return fmt.Errorf("gagal: Saldo deposit siswa tidak mencukupi (Saldo: %s, Digunakan: %s)", utils.FormatCurrency(student.DepositBalance), utils.FormatCurrency(p.DepositApplied))
			}

			if _, err := tx.NewUpdate().Model(&student).
				Set("deposit_balance = deposit_balance - ?", p.DepositApplied).
				Where("id = ?", studentID).
				Exec(ctx); err != nil {
				return err
			}

			dm := &domain.DepositMovement{
				StudentID:   studentID,
				Type:        "OUT",
				Amount:      p.DepositApplied,
				Reason:      "PAY_BILL",
				ReferenceID: utils.StringPtr(fmt.Sprintf("%d", p.ID)),
				CreatedBy:   p.CreatedBy,
			}
			if _, err := tx.NewInsert().Model(dm).Exec(ctx); err != nil {
				return err
			}
		}

		if isExisting {
			_ = s.repo.DeleteDetailsByPaymentID(ctx, tx, p.ID)
		}

		remainingAmount := p.Amount
		if len(billIDs) > 0 {
			for _, bID := range billIDs {
				if remainingAmount <= 0 {
					break
				}
				var b domain.StudentBill
				if err := tx.NewSelect().Model(&b).
					ColumnExpr("sb.*").
					ColumnExpr("bt.name as bill_type_name, br.allow_installment, br.max_installment").
					Join("JOIN bill_types bt ON sb.bill_type_id = bt.id").
					Join("LEFT JOIN billing_rules br ON sb.billing_rule_id = br.id").
					Where("sb.id = ?", bID).For("UPDATE").Scan(ctx); err != nil {
					return err
				}
				if role == "parent" {
					today := time.Now().Truncate(24 * time.Hour)
					if b.DueDate.Before(today) {
						return fmt.Errorf("gagal: Tagihan %s sudah jatuh tempo. Pembayaran online ditutup, silakan bayar langsung ke admin sekolah", b.BillTypeName)
					}
				}
				if b.Status == "voided" {
					return fmt.Errorf("gagal: Tagihan %s sudah ditarik/void dan tidak bisa dibayar", b.BillTypeName)
				}
				needed := b.Amount - b.TotalPaid
				if needed <= 0 {
					continue
				}
				paymentForThisBill := 0.0

				if remainingAmount >= needed {
					paymentForThisBill = needed
					b.TotalPaid = b.Amount
					b.Status = "paid"
					remainingAmount -= needed
				} else {
					if !p.IsBypassRule {
						if !b.AllowInstallment {
							return fmt.Errorf("gagal: Tagihan '%s' tidak melayani pembayaran sebagian/cicilan. Wajib dibayar lunas sebesar %s", b.BillTypeName, utils.FormatCurrency(needed))
						}
						if b.MaxInstallment > 0 {
							existingInstallments, err := tx.NewSelect().Model((*domain.PaymentDetail)(nil)).
								Where("student_bill_id = ?", b.ID).
								Count(ctx)
							if err != nil {
								return err
							}
							if existingInstallments >= b.MaxInstallment {
								return fmt.Errorf("gagal: Tagihan '%s' sudah mencapai batas maksimal %d kali cicilan", b.BillTypeName, b.MaxInstallment)
							}
						}
					}
					paymentForThisBill = remainingAmount
					b.TotalPaid += remainingAmount
					b.Status = "partial"
					remainingAmount = 0
				}
				if err := s.sbRepo.UpdateStatus(ctx, tx, b.ID, b.Status, b.TotalPaid); err != nil {
					return err
				}
				pd := &domain.PaymentDetail{PaymentID: p.ID, StudentBillID: b.ID, Amount: paymentForThisBill}
				_ = s.repo.CreateDetail(ctx, tx, pd)
			}
		}

		// Handle Overpayment: Kelebihan dana dialihkan ke saldo deposit siswa secara tertulis
		if remainingAmount > 0 {
			_, err := tx.NewUpdate().Model((*academicdomain.Student)(nil)).
				Set("deposit_balance = deposit_balance + ?", remainingAmount).
				Where("id = ?", studentID).
				Exec(ctx)
			if err != nil {
				return err
			}

			dm := &domain.DepositMovement{
				StudentID:   studentID,
				Type:        "IN",
				Amount:      remainingAmount,
				Reason:      "OVERPAYMENT",
				ReferenceID: utils.StringPtr(fmt.Sprintf("%d", p.ID)),
				CreatedBy:   p.CreatedBy,
			}
			if _, err := tx.NewInsert().Model(dm).Exec(ctx); err != nil {
				return err
			}
		}

		s.notifSvc.Notify(notificationusecase.FinanceNotifyJob{StudentID: p.StudentID, Payment: p, NotifType: "payment_success"})

		if s.audit != nil {
			userID, userName, role, ipAddress, userAgent := utils.GetAuditMeta(ctx)
			newVals := map[string]interface{}{
				"student_id": p.StudentID, "amount": p.Amount, "deposit_applied": p.DepositApplied, "cash_or_gateway_amount": p.Amount - p.DepositApplied, "transaction_ref": p.TransactionRef, "status": p.Status, "method": p.Method,
				"is_bypass_rule": p.IsBypassRule, "bypass_reason": stringValue(p.BypassReason), "note": stringValue(p.Note),
				"proof_attachment": stringValue(p.ProofAttachment), "created_by": p.CreatedBy,
			}
			_ = s.audit.Log(ctx, tx, userID, userName, role, "PROCESS_PAYMENT", "payments", p.ID, nil, newVals, ipAddress, userAgent)
		}

		if s.hub != nil {
			student, _ := s.stuRepo.FindByID(ctx, p.StudentID)
			studentName := "Siswa"
			if student != nil {
				studentName = student.Name
			}
			s.hub.Broadcast("NEW_PAYMENT", map[string]interface{}{"payment_id": p.ID, "student_name": studentName, "amount": p.Amount, "timestamp": time.Now().Format(time.RFC3339)})
		}
		return nil
	})
}

func (s *paymentService) GetHistory(ctx context.Context, studentID uint) ([]domain.Payment, error) {
	if s.audit != nil {
		userID, _, role, _, _ := utils.GetAuditMeta(ctx)
		if role == "parent" {
			student, err := s.stuRepo.FindByID(ctx, studentID)
			if err != nil || student == nil {
				return nil, fmt.Errorf("siswa tidak ditemukan")
			}
			if student.ParentID != userID {
				return nil, fmt.Errorf("unauthorized: Anda tidak memiliki akses ke data siswa ini")
			}
		}
	}
	return s.repo.FindByStudent(ctx, studentID)
}

func (s *paymentService) checkMidtransStatus(ctx context.Context, orderID string) (string, string, string, error) {
	baseURL := "https://api.sandbox.midtrans.com"
	if !s.cfg.MidtransIsSandbox {
		baseURL = "https://api.midtrans.com"
	}
	url := fmt.Sprintf("%s/v2/%s/status", baseURL, orderID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", "", "", err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(s.cfg.MidtransServerKey, "")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return "not_found", "", "", nil
	}
	if resp.StatusCode >= 400 {
		return "", "", "", fmt.Errorf("midtrans api returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", "", err
	}

	var result struct {
		TransactionStatus string `json:"transaction_status"`
		PaymentType       string `json:"payment_type"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", "", string(body), err
	}

	return result.TransactionStatus, result.PaymentType, string(body), nil
}

func (s *paymentService) startReconciler() {
	go func() {
		ticker := time.NewTicker(15 * time.Minute)
		ctx := context.Background()
		for range ticker.C {
			var pendingPayments []domain.Payment
			nowWindow := time.Now()
			err := s.db.NewSelect().Model(&pendingPayments).
				Where("status = ?", "pending").
				Where("gateway_provider = ?", "midtrans").
				Where("created_at >= ?", nowWindow.Add(-48*time.Hour)).
				Where("created_at <= ?", nowWindow.Add(-2*time.Minute)).
				Where("last_checked_at IS NULL OR last_checked_at <= ?", nowWindow.Add(-15*time.Minute)).
				Order("created_at ASC").
				Limit(50).
				Scan(ctx)
			if err != nil {
				continue
			}

			for _, p := range pendingPayments {
				pCopy := p
				now := time.Now()
				txStatus, payType, raw, err := s.checkMidtransStatus(ctx, pCopy.TransactionRef)
				pCopy.ReconcileAttempts++
				pCopy.LastCheckedAt = &now
				if raw != "" {
					pCopy.GatewayRawResponse = utils.StringPtr(raw)
				}
				if err != nil {
					pCopy.LastReconcileError = utils.StringPtr(err.Error())
					_, _ = s.db.NewUpdate().Model(&pCopy).Column("reconcile_attempts", "last_reconcile_error", "last_checked_at", "gateway_raw_response", "updated_at").WherePK().Exec(ctx)
					continue
				}

				switch txStatus {
				case "capture", "settlement":
					pCopy.Status = "success"
					pCopy.Method = payType
					pCopy.PaidAt = &now
					pCopy.ReconciledAt = &now
					if err := s.Process(ctx, pCopy.StudentID, &pCopy, nil); err != nil {
						if s.isRecoverableGatewayAllocationError(err) {
							_ = s.moveGatewayPaymentToDeposit(ctx, &pCopy, payType, raw, now, err)
							continue
						}
						pCopy.LastReconcileError = utils.StringPtr(err.Error())
						_, _ = s.db.NewUpdate().Model(&pCopy).Column("reconcile_attempts", "last_reconcile_error", "last_checked_at", "gateway_raw_response", "updated_at").WherePK().Exec(ctx)
					}
				case "expire", "cancel", "deny":
					pCopy.Status = "failed"
					pCopy.Note = utils.StringPtr(fmt.Sprintf("Payment %s", txStatus))
					_, _ = s.db.NewUpdate().Model(&pCopy).Column("status", "note", "reconcile_attempts", "last_checked_at", "gateway_raw_response", "updated_at").WherePK().Exec(ctx)
				case "pending":
					_, _ = s.db.NewUpdate().Model(&pCopy).Column("reconcile_attempts", "last_checked_at", "gateway_raw_response", "updated_at").WherePK().Exec(ctx)
				}
			}
		}
	}()
}
