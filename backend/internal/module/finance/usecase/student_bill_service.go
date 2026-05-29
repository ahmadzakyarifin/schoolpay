package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	academicdomain "github.com/ahmadzakyarifin/schoolpay/internal/module/academic/domain"
	academicrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/academic/repository"
	auditusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/audit/usecase"
	"github.com/ahmadzakyarifin/schoolpay/internal/module/finance/domain"
	"github.com/ahmadzakyarifin/schoolpay/internal/module/finance/repository"
	notificationusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/notification/usecase"
	"github.com/ahmadzakyarifin/schoolpay/pkg/utils"
	"github.com/uptrace/bun"
)

type StudentBillService interface {
	GenerateFromRule(ctx context.Context, ruleID uint, customReason, customMessage string, skipNotification bool) error
	BulkGenerateFromRules(ctx context.Context, ruleIDs []uint, customReason, customMessage string, skipNotification bool) error
	CancelGeneratedBills(ctx context.Context, ruleIDs []uint, customReason, customMessage string, skipNotification bool) error
	GetByStudent(ctx context.Context, studentID uint) ([]domain.StudentBill, error)
	GetByParent(ctx context.Context, parentID uint) ([]domain.StudentBill, error)
	GetAll(ctx context.Context, search, sort string) ([]domain.StudentBill, error)
	GetByID(ctx context.Context, id uint) (*domain.StudentBill, error)
	Create(ctx context.Context, sb *domain.StudentBill) error
	Update(ctx context.Context, sb *domain.StudentBill) error
	Delete(ctx context.Context, id uint) error
	SendManualReminder(ctx context.Context, id uint) error
	RunScheduler()
}

type studentBillService struct {
	db       *bun.DB
	repo     repository.StudentBillRepo
	ruleRepo repository.BillingRuleRepo
	stuRepo  academicrepo.StudentRepo
	notifSvc notificationusecase.FinanceNotificationService
	audit    auditusecase.AuditLogService
}

func NewStudentBillService(db *bun.DB, repo repository.StudentBillRepo, ruleRepo repository.BillingRuleRepo, stuRepo academicrepo.StudentRepo, notifSvc notificationusecase.FinanceNotificationService, audit auditusecase.AuditLogService) StudentBillService {
	return &studentBillService{db, repo, ruleRepo, stuRepo, notifSvc, audit}
}

type generatedBillPeriod struct {
	Period    string
	Month     *int
	Year      *int
	StartDate *time.Time
	EndDate   *time.Time
	DueDate   time.Time
}

func monthNameID(month time.Month) string {
	names := map[time.Month]string{
		time.January: "Januari", time.February: "Februari", time.March: "Maret",
		time.April: "April", time.May: "Mei", time.June: "Juni",
		time.July: "Juli", time.August: "Agustus", time.September: "September",
		time.October: "Oktober", time.November: "November", time.December: "Desember",
	}
	return names[month]
}

func dueDateInMonth(year int, month time.Month, dueDay int) time.Time {
	lastDay := time.Date(year, month+1, 0, 0, 0, 0, 0, time.Local).Day()
	day := dueDay
	if day > lastDay {
		day = lastDay
	}
	return time.Date(year, month, day, 0, 0, 0, 0, time.Local)
}

func buildBillName(billTypeName string, period generatedBillPeriod) string {
	if period.Month != nil && period.Year != nil {
		return fmt.Sprintf("%s %s %d", billTypeName, monthNameID(time.Month(*period.Month)), *period.Year)
	}
	return billTypeName
}

func buildCurrentMonthlyPeriod(rule *domain.BillingRule, now time.Time) (generatedBillPeriod, bool) {
	if rule.PeriodType == nil || *rule.PeriodType != "bulanan" {
		return generatedBillPeriod{}, false
	}
	if rule.StartDate != nil && now.Before(*rule.StartDate) {
		return generatedBillPeriod{}, false
	}
	if rule.EndDate != nil && now.After(*rule.EndDate) {
		return generatedBillPeriod{}, false
	}

	periodStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
	periodEnd := periodStart.AddDate(0, 1, -1)
	if rule.StartDate != nil && rule.StartDate.After(periodStart) {
		periodStart = time.Date(rule.StartDate.Year(), rule.StartDate.Month(), rule.StartDate.Day(), 0, 0, 0, 0, time.Local)
	}
	if rule.EndDate != nil && rule.EndDate.Before(periodEnd) {
		periodEnd = time.Date(rule.EndDate.Year(), rule.EndDate.Month(), rule.EndDate.Day(), 0, 0, 0, 0, time.Local)
	}
	if periodStart.After(periodEnd) {
		return generatedBillPeriod{}, false
	}

	month := int(now.Month())
	year := now.Year()
	dueDate := dueDateInMonth(year, now.Month(), rule.DueDay)
	if dueDate.Before(periodStart) {
		dueDate = periodStart
	}
	if dueDate.After(periodEnd) {
		dueDate = periodEnd
	}

	return generatedBillPeriod{
		Period:    now.Format("2006-01"),
		Month:     &month,
		Year:      &year,
		StartDate: &periodStart,
		EndDate:   &periodEnd,
		DueDate:   dueDate,
	}, true
}

func (s *studentBillService) buildPeriods(rule *domain.BillingRule) ([]generatedBillPeriod, error) {
	if rule.DueDay < 1 || rule.DueDay > 31 {
		return nil, errors.New("gagal: Tanggal jatuh tempo bulanan harus berada di antara tanggal 1 sampai 31")
	}

	periods := []generatedBillPeriod{{DueDate: time.Now()}}
	if rule.EndDate != nil {
		periods[0].DueDate = *rule.EndDate
		periods[0].EndDate = rule.EndDate
	}

	if rule.PeriodType != nil && *rule.PeriodType == "bulanan" && rule.StartDate != nil && rule.EndDate != nil {
		periods = []generatedBillPeriod{}
		start := time.Date(rule.StartDate.Year(), rule.StartDate.Month(), 1, 0, 0, 0, 0, time.Local)
		end := time.Date(rule.EndDate.Year(), rule.EndDate.Month(), 1, 0, 0, 0, 0, time.Local)
		for !start.After(end) {
			month := int(start.Month())
			year := start.Year()
			periodStart := start
			periodEnd := start.AddDate(0, 1, -1)
			dueDate := dueDateInMonth(year, start.Month(), rule.DueDay)
			periods = append(periods, generatedBillPeriod{
				Period:    start.Format("2006-01"),
				Month:     &month,
				Year:      &year,
				StartDate: &periodStart,
				EndDate:   &periodEnd,
				DueDate:   dueDate,
			})
			start = start.AddDate(0, 1, 0)
		}
	}

	if rule.PeriodType != nil && *rule.PeriodType == "tahunan" && rule.StartDate != nil && rule.EndDate != nil {
		periods = []generatedBillPeriod{}
		for year := rule.StartDate.Year(); year <= rule.EndDate.Year(); year++ {
			periodStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.Local)
			periodEnd := time.Date(year, 12, 31, 0, 0, 0, 0, time.Local)
			if year == rule.StartDate.Year() {
				periodStart = *rule.StartDate
			}
			if year == rule.EndDate.Year() {
				periodEnd = *rule.EndDate
			}
			dueMonth := periodStart.Month()
			dueDate := dueDateInMonth(year, dueMonth, rule.DueDay)
			if dueDate.Before(periodStart) {
				dueDate = periodStart
			}
			if dueDate.After(periodEnd) {
				dueDate = periodEnd
			}
			yearCopy := year
			periods = append(periods, generatedBillPeriod{
				Period:    fmt.Sprintf("%d", year),
				Year:      &yearCopy,
				StartDate: &periodStart,
				EndDate:   &periodEnd,
				DueDate:   dueDate,
			})
		}
	}
	return periods, nil
}

func (s *studentBillService) GenerateFromRule(ctx context.Context, ruleID uint, customReason, customMessage string, skipNotification bool) error {
	rule, err := s.ruleRepo.FindByID(ctx, ruleID)
	if err != nil {
		return errors.New("gagal: Aturan tagihan tidak ditemukan")
	}

	if !rule.IsActive {
		return errors.New("gagal: Aturan tagihan sedang nonaktif. Aktifkan aturan sebelum generate tagihan")
	}

	var billType domain.BillType
	if err := s.db.NewSelect().Model(&billType).Where("id = ? AND deleted_at IS NULL", rule.BillTypeID).Scan(ctx); err != nil {
		return errors.New("gagal: Jenis tagihan tidak ditemukan")
	}
	if !billType.IsActive {
		return errors.New("gagal: Jenis tagihan sedang nonaktif. Aktifkan jenis tagihan sebelum generate")
	}

	targetType := rule.TargetType
	targetID := rule.TargetID
	if rule.ClassID != nil && *rule.ClassID > 0 {
		targetType = "class"
		targetID = *rule.ClassID
	}

	studentIDs, err := s.stuRepo.FindIDsByTarget(ctx, targetType, targetID)
	if err != nil {
		return err
	}
	if len(studentIDs) == 0 {
		return errors.New("gagal: Tidak ada siswa aktif yang sesuai dengan target aturan tagihan")
	}

	periods, err := s.buildPeriods(rule)
	if err != nil {
		return err
	}

	for _, sID := range studentIDs {
		for _, period := range periods {
			periodCopy := period.Period
			billName := buildBillName(billType.Name, period)

			exists, err := s.repo.ExistsByPeriod(ctx, sID, rule.BillTypeID, periodCopy)
			if err != nil {
				return err
			}
			if exists {
				var existingBill domain.StudentBill
				if err := s.db.NewSelect().Model(&existingBill).
					Where("student_id = ? AND bill_type_id = ? AND COALESCE(period, '') = ? AND status != 'voided'", sID, rule.BillTypeID, periodCopy).
					Scan(ctx); err != nil {
					return err
				}
				existingBill.Name = &billName
				existingBill.DueDate = period.DueDate
				existingBill.EndDate = rule.EndDate
				existingBill.BillingRuleID = &rule.ID
				existingBill.PeriodMonth = period.Month
				existingBill.PeriodYear = period.Year
				existingBill.PeriodStartDate = period.StartDate
				existingBill.PeriodEndDate = period.EndDate

				oldAmount := existingBill.Amount
				oldPaid := existingBill.TotalPaid
				depositRefund := 0.0
				if existingBill.Status == "unpaid" {
					existingBill.Amount = rule.Amount
				} else if rule.Amount < oldAmount {
					existingBill.Amount = rule.Amount
					if existingBill.TotalPaid > rule.Amount {
						depositRefund = existingBill.TotalPaid - rule.Amount
						existingBill.TotalPaid = rule.Amount
					}
					if existingBill.TotalPaid >= existingBill.Amount {
						existingBill.Status = "paid"
					} else if existingBill.TotalPaid > 0 {
						existingBill.Status = "partial"
					} else {
						existingBill.Status = "unpaid"
					}
				}
				if _, err := s.db.NewUpdate().Model(&existingBill).WherePK().Exec(ctx); err != nil {
					return err
				}

				if depositRefund > 0 {
					_, err := s.db.NewUpdate().Model((*academicdomain.Student)(nil)).
						Set("deposit_balance = deposit_balance + ?", depositRefund).
						Where("id = ?", sID).
						Exec(ctx)
					if err != nil {
						return err
					}
					dm := &domain.DepositMovement{StudentID: sID, Type: "IN", Amount: depositRefund, Reason: "BILL_AMOUNT_REDUCED", ReferenceID: utils.StringPtr(fmt.Sprintf("%d", existingBill.ID)), CreatedBy: "SYSTEM"}
					if _, err := s.db.NewInsert().Model(dm).Exec(ctx); err != nil {
						return err
					}
					if s.audit != nil {
						userID, userName, role, ipAddress, userAgent := utils.GetAuditMeta(ctx)
						oldVals := map[string]interface{}{"amount": oldAmount, "total_paid": oldPaid}
						newVals := map[string]interface{}{"amount": existingBill.Amount, "total_paid": existingBill.TotalPaid, "deposit_refund": depositRefund}
						_ = s.audit.Log(ctx, s.db, userID, userName, role, "REFUND_BILL_REDUCTION", "student_bills", existingBill.ID, oldVals, newVals, ipAddress, userAgent)
					}
					if !skipNotification {
						message := customMessage
						if message == "" {
							message = fmt.Sprintf("Nominal tagihan %s disesuaikan turun. Selisih pembayaran sebesar Rp %.0f dialihkan ke Saldo Deposit siswa.", billName, depositRefund)
							if customReason != "" {
								message += " Alasan: " + customReason
							}
						}
						s.notifSvc.Notify(notificationusecase.FinanceNotifyJob{StudentID: sID, Bill: &existingBill, NotifType: "refund_deposit", CustomMessage: message, CustomReason: customReason})
					}
				}

				if existingBill.Status != "unpaid" && rule.Amount > oldAmount {
					adjustmentBase := periodCopy
					if adjustmentBase == "" {
						adjustmentBase = fmt.Sprintf("single-%d", rule.ID)
					}

					var adjustmentBills []domain.StudentBill
					if err := s.db.NewSelect().Model(&adjustmentBills).
						Where("student_id = ? AND bill_type_id = ? AND billing_rule_id = ? AND COALESCE(period, '') LIKE ? AND status != 'voided'", sID, rule.BillTypeID, rule.ID, adjustmentBase+"-ADJ-%").
						Scan(ctx); err != nil {
						return err
					}
					existingAdjustmentAmount := 0.0
					for _, adj := range adjustmentBills {
						existingAdjustmentAmount += adj.Amount
					}

					desiredAdjustment := rule.Amount - oldAmount
					deficit := desiredAdjustment - existingAdjustmentAmount
					if deficit <= 0 {
						continue
					}

					desc := fmt.Sprintf("Tagihan penyesuaian atas perubahan kebijakan tarif dari Rp %.0f menjadi Rp %.0f. Selisih kekurangan: Rp %.0f.", oldAmount, rule.Amount, deficit)
					if customReason != "" {
						desc += " Alasan penyesuaian: " + customReason
					}

					adjName := billName + " - Penyesuaian"
					adjPeriod := fmt.Sprintf("%s-ADJ-%d-%d", adjustmentBase, rule.ID, time.Now().UnixNano())
					adjBill := &domain.StudentBill{
						StudentID: sID, BillTypeID: rule.BillTypeID, BillingRuleID: &rule.ID,
						Name: &adjName, AcademicYear: existingBill.AcademicYear, Period: &adjPeriod,
						PeriodMonth: period.Month, PeriodYear: period.Year, PeriodStartDate: period.StartDate, PeriodEndDate: period.EndDate,
						Amount: deficit, Status: "unpaid", DueDate: period.DueDate, EndDate: rule.EndDate, Description: desc,
					}

					if err := s.repo.Create(ctx, s.db, adjBill); err != nil {
						return err
					}
					if s.audit != nil {
						userID, userName, role, ipAddress, userAgent := utils.GetAuditMeta(ctx)
						newVals := map[string]interface{}{"student_id": sID, "bill_type_id": rule.BillTypeID, "amount": deficit, "academic_year": existingBill.AcademicYear, "period": periodCopy, "description": desc}
						_ = s.audit.Log(ctx, s.db, userID, userName, role, "GENERATE_ADJUSTMENT_BILL", "student_bills", adjBill.ID, nil, newVals, ipAddress, userAgent)
					}
					if !skipNotification {
						s.notifSvc.Notify(notificationusecase.FinanceNotifyJob{StudentID: sID, Bill: adjBill, NotifType: "adjustment", CustomMessage: customMessage, CustomReason: customReason})
					}
				}
				continue
			}

			academicYear := ""
			pTime := time.Now()
			if period.StartDate != nil {
				pTime = *period.StartDate
			}
			if pTime.Month() >= 7 {
				academicYear = fmt.Sprintf("%d/%d", pTime.Year(), pTime.Year()+1)
			} else {
				academicYear = fmt.Sprintf("%d/%d", pTime.Year()-1, pTime.Year())
			}

			bill := &domain.StudentBill{
				StudentID: sID, BillTypeID: rule.BillTypeID, BillingRuleID: &rule.ID,
				Name: &billName, AcademicYear: academicYear, Period: &periodCopy,
				PeriodMonth: period.Month, PeriodYear: period.Year, PeriodStartDate: period.StartDate, PeriodEndDate: period.EndDate,
				Amount: rule.Amount, Status: "unpaid", DueDate: period.DueDate, EndDate: rule.EndDate,
			}

			if err := s.repo.Create(ctx, s.db, bill); err != nil {
				return err
			}
			if s.audit != nil {
				userID, userName, role, ipAddress, userAgent := utils.GetAuditMeta(ctx)
				newVals := map[string]interface{}{"student_id": sID, "bill_type_id": rule.BillTypeID, "amount": rule.Amount, "academic_year": academicYear, "period": periodCopy, "name": billName}
				_ = s.audit.Log(ctx, s.db, userID, userName, role, "GENERATE_BILL", "student_bills", bill.ID, nil, newVals, ipAddress, userAgent)
			}
			if !skipNotification {
				s.notifSvc.Notify(notificationusecase.FinanceNotifyJob{StudentID: sID, Bill: bill, NotifType: "initial", CustomMessage: customMessage, CustomReason: customReason})
			}
		}
	}
	return nil
}

func (s *studentBillService) BulkGenerateFromRules(ctx context.Context, ruleIDs []uint, customReason, customMessage string, skipNotification bool) error {
	for _, ruleID := range ruleIDs {
		if err := s.GenerateFromRule(ctx, ruleID, customReason, customMessage, skipNotification); err != nil {
			return err
		}
	}
	return nil
}

func (s *studentBillService) CancelGeneratedBills(ctx context.Context, ruleIDs []uint, customReason, customMessage string, skipNotification bool) error {
	return s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		var bills []domain.StudentBill
		if err := tx.NewSelect().Model(&bills).
			Where("billing_rule_id IN (?) AND status != 'voided'", bun.In(ruleIDs)).
			For("UPDATE").
			Scan(ctx); err != nil {
			return err
		}

		adminID, adminName, _, _, _ := utils.GetAuditMeta(ctx)
		createdBy := "SYSTEM"
		if adminName != "" {
			createdBy = fmt.Sprintf("ADMIN-%d (%s)", adminID, adminName)
		}
		now := time.Now()

		for _, b := range bills {
			refundAmount := b.TotalPaid
			oldStatus := b.Status
			b.Status = "voided"
			b.VoidedAt = &now
			voidReason := "BULK_CANCEL_GENERATED_BILL"
			if customReason != "" {
				voidReason = customReason
			}
			b.VoidReason = utils.StringPtr(voidReason)
			if _, err := tx.NewUpdate().Model(&b).Column("status", "voided_at", "void_reason").WherePK().Exec(ctx); err != nil {
				return err
			}

			if refundAmount > 0 {
				if _, err := tx.NewUpdate().Model((*academicdomain.Student)(nil)).
					Set("deposit_balance = deposit_balance + ?", refundAmount).
					Where("id = ?", b.StudentID).
					Exec(ctx); err != nil {
					return err
				}
				dm := &domain.DepositMovement{StudentID: b.StudentID, Type: "IN", Amount: refundAmount, Reason: "BILL_VOIDED", ReferenceID: utils.StringPtr(fmt.Sprintf("%d", b.ID)), CreatedBy: createdBy}
				if _, err := tx.NewInsert().Model(dm).Exec(ctx); err != nil {
					return err
				}

				if !skipNotification {
					message := customMessage
					if message == "" {
						message = fmt.Sprintf("Tagihan ditarik oleh sekolah. Dana yang sudah dibayarkan sebesar Rp %.0f dialihkan ke Saldo Deposit siswa.", refundAmount)
						if customReason != "" {
							message += " Alasan: " + customReason
						}
					}
					bCopy := b
					go s.notifSvc.Notify(notificationusecase.FinanceNotifyJob{
						StudentID:     b.StudentID,
						Bill:          &bCopy,
						NotifType:     "refund_deposit",
						CustomMessage: message,
						CustomReason:  customReason,
					})
				}
			}

			if s.audit != nil {
				userID, userName, role, ipAddress, userAgent := utils.GetAuditMeta(ctx)
				oldVals := map[string]interface{}{"status": oldStatus, "total_paid": refundAmount}
				newVals := map[string]interface{}{"status": "voided", "deposit_refund": refundAmount, "void_reason": voidReason}
				_ = s.audit.Log(ctx, tx, userID, userName, role, "BULK_CANCEL_GENERATED_BILL", "student_bills", b.ID, oldVals, newVals, ipAddress, userAgent)
			}
		}
		return nil
	})
}

func (s *studentBillService) GetByStudent(ctx context.Context, studentID uint) ([]domain.StudentBill, error) {
	return s.repo.FindByStudent(ctx, studentID)
}

func (s *studentBillService) GetByParent(ctx context.Context, parentID uint) ([]domain.StudentBill, error) {
	return s.repo.FindByParent(ctx, parentID)
}

func (s *studentBillService) GetAll(ctx context.Context, search, sort string) ([]domain.StudentBill, error) {
	return s.repo.FindAll(ctx, search, sort)
}

func (s *studentBillService) Create(ctx context.Context, sb *domain.StudentBill) error {
	period := ""
	if sb.Period != nil {
		period = *sb.Period
	}
	exists, err := s.repo.ExistsByPeriod(ctx, sb.StudentID, sb.BillTypeID, period)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("gagal: Tagihan untuk siswa dan periode ini sudah terdaftar")
	}

	err = s.repo.Create(ctx, s.db, sb)
	if err == nil && s.audit != nil {
		userID, userName, role, ipAddress, userAgent := utils.GetAuditMeta(ctx)
		newVals := map[string]interface{}{"student_id": sb.StudentID, "bill_type_id": sb.BillTypeID, "amount": sb.Amount, "academic_year": sb.AcademicYear, "period": period}
		_ = s.audit.Log(ctx, s.db, userID, userName, role, "CREATE", "student_bills", sb.ID, nil, newVals, ipAddress, userAgent)
	}
	return err
}

func (s *studentBillService) Update(ctx context.Context, sb *domain.StudentBill) error {
	if sb.Amount <= 0 {
		return errors.New("gagal: Nominal tagihan wajib lebih dari 0")
	}

	return s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		var existing domain.StudentBill
		if err := tx.NewSelect().Model(&existing).Where("id = ?", sb.ID).For("UPDATE").Scan(ctx); err != nil {
			return err
		}
		if existing.Status == "voided" {
			return errors.New("gagal: Tagihan yang sudah dibatalkan tidak dapat diubah")
		}

		// Merge partial update payload with the locked row so missing JSON fields do not zero existing bill data.
		if sb.StudentID == 0 {
			sb.StudentID = existing.StudentID
		}
		if sb.BillTypeID == 0 {
			sb.BillTypeID = existing.BillTypeID
		}
		if sb.BillingRuleID == nil {
			sb.BillingRuleID = existing.BillingRuleID
		}
		if sb.Name == nil {
			sb.Name = existing.Name
		}
		if sb.AcademicYear == "" {
			sb.AcademicYear = existing.AcademicYear
		}
		if sb.Period == nil {
			sb.Period = existing.Period
		}
		if sb.PeriodMonth == nil {
			sb.PeriodMonth = existing.PeriodMonth
		}
		if sb.PeriodYear == nil {
			sb.PeriodYear = existing.PeriodYear
		}
		if sb.PeriodStartDate == nil {
			sb.PeriodStartDate = existing.PeriodStartDate
		}
		if sb.PeriodEndDate == nil {
			sb.PeriodEndDate = existing.PeriodEndDate
		}
		if sb.DueDate.IsZero() {
			sb.DueDate = existing.DueDate
		}
		if sb.EndDate == nil {
			sb.EndDate = existing.EndDate
		}
		if sb.Description == "" {
			sb.Description = existing.Description
		}
		sb.CreatedAt = existing.CreatedAt
		sb.DeletedAt = existing.DeletedAt
		sb.VoidedAt = existing.VoidedAt
		sb.VoidReason = existing.VoidReason
		sb.LastNotifiedAt = existing.LastNotifiedAt
		sb.NextNotifiedAt = existing.NextNotifiedAt
		sb.UQKeyPeriod = existing.UQKeyPeriod

		period := ""
		if sb.Period != nil {
			period = *sb.Period
		}
		exists, err := s.repo.ExistsByPeriodExcludeID(ctx, sb.StudentID, sb.BillTypeID, period, sb.ID)
		if err != nil {
			return err
		}
		if exists {
			return errors.New("gagal: Data tagihan serupa sudah ada untuk siswa ini")
		}

		depositRefund := 0.0
		oldVals := map[string]interface{}{"student_id": existing.StudentID, "bill_type_id": existing.BillTypeID, "amount": existing.Amount, "total_paid": existing.TotalPaid, "status": existing.Status, "academic_year": existing.AcademicYear, "period": period}

		if existing.TotalPaid > sb.Amount {
			depositRefund = existing.TotalPaid - sb.Amount
			sb.TotalPaid = sb.Amount
		} else {
			sb.TotalPaid = existing.TotalPaid
		}

		if sb.TotalPaid >= sb.Amount {
			sb.Status = "paid"
		} else if sb.TotalPaid > 0 {
			sb.Status = "partial"
		} else {
			sb.Status = "unpaid"
		}

		if err := s.repo.Update(ctx, tx, sb); err != nil {
			return err
		}

		if depositRefund > 0 {
			adminID, adminName, _, _, _ := utils.GetAuditMeta(ctx)
			createdBy := "SYSTEM"
			if adminName != "" {
				createdBy = fmt.Sprintf("ADMIN-%d (%s)", adminID, adminName)
			}

			if _, err := tx.NewUpdate().Model((*academicdomain.Student)(nil)).
				Set("deposit_balance = deposit_balance + ?", depositRefund).
				Where("id = ?", sb.StudentID).
				Exec(ctx); err != nil {
				return err
			}

			dm := &domain.DepositMovement{
				StudentID:   sb.StudentID,
				Type:        "IN",
				Amount:      depositRefund,
				Reason:      "BILL_AMOUNT_REDUCED",
				ReferenceID: utils.StringPtr(fmt.Sprintf("%d", sb.ID)),
				CreatedBy:   createdBy,
			}
			if _, err := tx.NewInsert().Model(dm).Exec(ctx); err != nil {
				return err
			}
		}

		if s.audit != nil {
			userID, userName, role, ipAddress, userAgent := utils.GetAuditMeta(ctx)
			newVals := map[string]interface{}{"student_id": sb.StudentID, "bill_type_id": sb.BillTypeID, "amount": sb.Amount, "total_paid": sb.TotalPaid, "status": sb.Status, "academic_year": sb.AcademicYear, "period": period, "deposit_refund": depositRefund}
			_ = s.audit.Log(ctx, tx, userID, userName, role, "UPDATE", "student_bills", sb.ID, oldVals, newVals, ipAddress, userAgent)
		}
		return nil
	})
}

func (s *studentBillService) Delete(ctx context.Context, id uint) error {
	return s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		var b domain.StudentBill
		if err := tx.NewSelect().Model(&b).Where("id = ?", id).For("UPDATE").Scan(ctx); err != nil {
			return err
		}

		if b.Status == "voided" {
			return nil
		}

		adminID, adminName, _, _, _ := utils.GetAuditMeta(ctx)
		createdBy := "SYSTEM"
		if adminName != "" {
			createdBy = fmt.Sprintf("ADMIN-%d (%s)", adminID, adminName)
		}

		refundAmount := b.TotalPaid
		oldStatus := b.Status
		now := time.Now()
		b.Status = "voided"
		b.VoidedAt = &now
		b.VoidReason = utils.StringPtr("MANUAL_VOID_BILL")
		if _, err := tx.NewUpdate().Model(&b).Column("status", "voided_at", "void_reason").WherePK().Exec(ctx); err != nil {
			return err
		}

		if refundAmount > 0 {
			_, err := tx.NewUpdate().Model((*academicdomain.Student)(nil)).
				Set("deposit_balance = deposit_balance + ?", refundAmount).
				Where("id = ?", b.StudentID).
				Exec(ctx)
			if err != nil {
				return err
			}

			dm := &domain.DepositMovement{
				StudentID:   b.StudentID,
				Type:        "IN",
				Amount:      refundAmount,
				Reason:      "BILL_VOIDED",
				ReferenceID: utils.StringPtr(fmt.Sprintf("%d", b.ID)),
				CreatedBy:   createdBy,
			}
			if _, err := tx.NewInsert().Model(dm).Exec(ctx); err != nil {
				return err
			}

			go func(studentID uint, amount float64) {
				s.notifSvc.Notify(notificationusecase.FinanceNotifyJob{
					StudentID:     studentID,
					NotifType:     "refund_deposit",
					CustomMessage: fmt.Sprintf("Tagihan ditarik. Dana pembayaran sebesar Rp %.0f dialihkan ke Saldo Deposit anak Anda.", amount),
				})
			}(b.StudentID, refundAmount)
		}

		if s.audit != nil {
			userID, userName, role, ipAddress, userAgent := utils.GetAuditMeta(ctx)
			oldVals := map[string]interface{}{"status": oldStatus, "total_paid": refundAmount}
			newVals := map[string]interface{}{"status": "voided"}
			_ = s.audit.Log(ctx, tx, userID, userName, role, "VOID_BILL", "student_bills", id, oldVals, newVals, ipAddress, userAgent)
		}
		return nil
	})
}

func (s *studentBillService) RunScheduler() {
	ticker := time.NewTicker(24 * time.Hour)
	go func() {
		ctx := context.Background()
		// Run once on startup
		s.autoGenerateMonthlyBills(ctx)

		for range ticker.C {
			// 1. Auto generate monthly bills
			s.autoGenerateMonthlyBills(ctx)

			// 2. Send reminders
			billsH3, _ := s.repo.FindForReminder(ctx, 3, false)
			for _, b := range billsH3 {
				bCopy := b
				s.notifSvc.Notify(notificationusecase.FinanceNotifyJob{StudentID: b.StudentID, Bill: &bCopy, NotifType: "reminder"})
				now := time.Now()
				b.LastNotifiedAt = &now
				_, _ = s.db.NewUpdate().Model(&b).Column("last_notified_at").WherePK().Exec(ctx)
			}

			// 3. Send overdue notices
			overdueBills, _ := s.repo.FindForReminder(ctx, 0, true)
			for _, b := range overdueBills {
				bCopy := b
				s.notifSvc.Notify(notificationusecase.FinanceNotifyJob{StudentID: b.StudentID, Bill: &bCopy, NotifType: "overdue"})
				now := time.Now()
				b.LastNotifiedAt = &now
				_, _ = s.db.NewUpdate().Model(&b).Column("last_notified_at").WherePK().Exec(ctx)
			}
		}
	}()
}

func (s *studentBillService) autoGenerateMonthlyBills(ctx context.Context) {
	rules, err := s.ruleRepo.FindAll(ctx)
	if err != nil {
		return
	}

	now := time.Now()
	for _, rule := range rules {
		if !rule.IsActive {
			continue
		}

		period, ok := buildCurrentMonthlyPeriod(&rule, now)
		if !ok {
			continue
		}

		var billType domain.BillType
		if err := s.db.NewSelect().Model(&billType).Where("id = ? AND deleted_at IS NULL AND is_active = ?", rule.BillTypeID, true).Scan(ctx); err != nil {
			continue
		}

		targetType := rule.TargetType
		targetID := rule.TargetID
		if rule.ClassID != nil && *rule.ClassID > 0 {
			targetType = "class"
			targetID = *rule.ClassID
		}

		studentIDs, err := s.stuRepo.FindIDsByTarget(ctx, targetType, targetID)
		if err != nil || len(studentIDs) == 0 {
			continue
		}

		var existing []domain.StudentBill
		_ = s.db.NewSelect().Model(&existing).
			Column("student_id").
			Where("student_id IN (?)", bun.In(studentIDs)).
			Where("bill_type_id = ?", rule.BillTypeID).
			Where("COALESCE(period, '') = ?", period.Period).
			Where("status != ?", "voided").
			Scan(ctx)

		existingByStudent := map[uint]bool{}
		for _, bill := range existing {
			existingByStudent[bill.StudentID] = true
		}

		billName := buildBillName(billType.Name, period)
		academicYear := ""
		if period.StartDate != nil && period.StartDate.Month() >= 7 {
			academicYear = fmt.Sprintf("%d/%d", period.StartDate.Year(), period.StartDate.Year()+1)
		} else if period.StartDate != nil {
			academicYear = fmt.Sprintf("%d/%d", period.StartDate.Year()-1, period.StartDate.Year())
		}

		for _, studentID := range studentIDs {
			if existingByStudent[studentID] {
				continue
			}
			periodCopy := period.Period
			nameCopy := billName
			bill := &domain.StudentBill{
				StudentID: studentID, BillTypeID: rule.BillTypeID, BillingRuleID: &rule.ID,
				Name: &nameCopy, AcademicYear: academicYear, Period: &periodCopy,
				PeriodMonth: period.Month, PeriodYear: period.Year, PeriodStartDate: period.StartDate, PeriodEndDate: period.EndDate,
				Amount: rule.Amount, Status: "unpaid", DueDate: period.DueDate, EndDate: rule.EndDate,
			}
			if err := s.repo.Create(ctx, s.db, bill); err != nil {
				continue
			}
			if s.audit != nil {
				newVals := map[string]interface{}{"student_id": studentID, "bill_type_id": rule.BillTypeID, "amount": rule.Amount, "academic_year": academicYear, "period": periodCopy, "name": billName, "source": "scheduler_current_month"}
				_ = s.audit.Log(ctx, s.db, 0, "System/Automation", "system", "GENERATE_BILL", "student_bills", bill.ID, nil, newVals, "system", "scheduler")
			}
			s.notifSvc.Notify(notificationusecase.FinanceNotifyJob{StudentID: studentID, Bill: bill, NotifType: "initial"})
		}
	}
}

func (s *studentBillService) SendManualReminder(ctx context.Context, id uint) error {
	var b domain.StudentBill
	err := s.db.NewSelect().Model(&b).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return err
	}

	now := time.Now()
	if b.LastNotifiedAt != nil && now.Sub(*b.LastNotifiedAt) < 5*time.Minute {
		return errors.New("gagal: Reminder baru saja dikirim. Tunggu minimal 5 menit sebelum mengirim ulang")
	}

	notifType := "reminder"
	if b.DueDate.Before(now) {
		notifType = "overdue"
	}

	s.notifSvc.Notify(notificationusecase.FinanceNotifyJob{
		StudentID: b.StudentID,
		Bill:      &b,
		NotifType: notifType,
	})

	b.LastNotifiedAt = &now
	_, _ = s.db.NewUpdate().Model(&b).Column("last_notified_at").WherePK().Exec(ctx)

	if s.audit != nil {
		userID, userName, role, ipAddress, userAgent := utils.GetAuditMeta(ctx)
		newVals := map[string]interface{}{"student_id": b.StudentID, "bill_id": b.ID, "notif_type": notifType}
		_ = s.audit.Log(ctx, s.db, userID, userName, role, "SEND_MANUAL_BILL_REMINDER", "student_bills", b.ID, nil, newVals, ipAddress, userAgent)
	}

	return nil
}

func (s *studentBillService) GetByID(ctx context.Context, id uint) (*domain.StudentBill, error) {
	return s.repo.FindByID(ctx, id)
}
