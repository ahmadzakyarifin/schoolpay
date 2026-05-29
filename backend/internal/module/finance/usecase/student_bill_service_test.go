package usecase

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/ahmadzakyarifin/schoolpay/internal/mocks"
	academicrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/academic/repository"
	auditusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/audit/usecase"
	"github.com/ahmadzakyarifin/schoolpay/internal/module/finance/domain"
	"github.com/ahmadzakyarifin/schoolpay/internal/module/finance/repository"
	notificationusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/notification/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"
)

func TestNewStudentBillService(t *testing.T) {
	type args struct {
		db       *bun.DB
		repo     repository.StudentBillRepo
		ruleRepo repository.BillingRuleRepo
		stuRepo  academicrepo.StudentRepo
		notifSvc notificationusecase.FinanceNotificationService
		audit    auditusecase.AuditLogService
	}
	tests := []struct {
		name string
		args args
		want StudentBillService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewStudentBillService(tt.args.db, tt.args.repo, tt.args.ruleRepo, tt.args.stuRepo, tt.args.notifSvc, tt.args.audit); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStudentBillService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_monthNameID(t *testing.T) {
	type args struct {
		month time.Month
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := monthNameID(tt.args.month); got != tt.want {
				t.Errorf("monthNameID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_buildBillName(t *testing.T) {
	type args struct {
		billTypeName string
		period       generatedBillPeriod
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildBillName(tt.args.billTypeName, tt.args.period); got != tt.want {
				t.Errorf("buildBillName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_studentBillService_buildPeriods(t *testing.T) {
	type fields struct {
		db       *bun.DB
		repo     repository.StudentBillRepo
		ruleRepo repository.BillingRuleRepo
		stuRepo  academicrepo.StudentRepo
		notifSvc notificationusecase.FinanceNotificationService
		audit    auditusecase.AuditLogService
	}
	type args struct {
		rule *domain.BillingRule
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []generatedBillPeriod
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &studentBillService{
				db:       tt.fields.db,
				repo:     tt.fields.repo,
				ruleRepo: tt.fields.ruleRepo,
				stuRepo:  tt.fields.stuRepo,
				notifSvc: tt.fields.notifSvc,
				audit:    tt.fields.audit,
			}
			got, err := s.buildPeriods(tt.args.rule)
			if (err != nil) != tt.wantErr {
				t.Fatalf("studentBillService.buildPeriods() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("studentBillService.buildPeriods() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_studentBillService_GenerateFromRule(t *testing.T) {
	type fields struct {
		db       *bun.DB
		repo     repository.StudentBillRepo
		ruleRepo repository.BillingRuleRepo
		stuRepo  academicrepo.StudentRepo
		notifSvc notificationusecase.FinanceNotificationService
		audit    auditusecase.AuditLogService
	}
	type args struct {
		ctx              context.Context
		ruleID           uint
		customReason     string
		customMessage    string
		skipNotification bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &studentBillService{
				db:       tt.fields.db,
				repo:     tt.fields.repo,
				ruleRepo: tt.fields.ruleRepo,
				stuRepo:  tt.fields.stuRepo,
				notifSvc: tt.fields.notifSvc,
				audit:    tt.fields.audit,
			}
			if err := s.GenerateFromRule(tt.args.ctx, tt.args.ruleID, tt.args.customReason, tt.args.customMessage, tt.args.skipNotification); (err != nil) != tt.wantErr {
				t.Errorf("studentBillService.GenerateFromRule() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_studentBillService_BulkGenerateFromRules(t *testing.T) {
	type fields struct {
		db       *bun.DB
		repo     repository.StudentBillRepo
		ruleRepo repository.BillingRuleRepo
		stuRepo  academicrepo.StudentRepo
		notifSvc notificationusecase.FinanceNotificationService
		audit    auditusecase.AuditLogService
	}
	type args struct {
		ctx              context.Context
		ruleIDs          []uint
		customReason     string
		customMessage    string
		skipNotification bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &studentBillService{
				db:       tt.fields.db,
				repo:     tt.fields.repo,
				ruleRepo: tt.fields.ruleRepo,
				stuRepo:  tt.fields.stuRepo,
				notifSvc: tt.fields.notifSvc,
				audit:    tt.fields.audit,
			}
			if err := s.BulkGenerateFromRules(tt.args.ctx, tt.args.ruleIDs, tt.args.customReason, tt.args.customMessage, tt.args.skipNotification); (err != nil) != tt.wantErr {
				t.Errorf("studentBillService.BulkGenerateFromRules() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_studentBillService_CancelGeneratedBills(t *testing.T) {
	type fields struct {
		db       *bun.DB
		repo     repository.StudentBillRepo
		ruleRepo repository.BillingRuleRepo
		stuRepo  academicrepo.StudentRepo
		notifSvc notificationusecase.FinanceNotificationService
		audit    auditusecase.AuditLogService
	}
	type args struct {
		ctx     context.Context
		ruleIDs []uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &studentBillService{
				db:       tt.fields.db,
				repo:     tt.fields.repo,
				ruleRepo: tt.fields.ruleRepo,
				stuRepo:  tt.fields.stuRepo,
				notifSvc: tt.fields.notifSvc,
				audit:    tt.fields.audit,
			}
			if err := s.CancelGeneratedBills(tt.args.ctx, tt.args.ruleIDs, "", "", false); (err != nil) != tt.wantErr {
				t.Errorf("studentBillService.CancelGeneratedBills() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_studentBillService_GetByStudent(t *testing.T) {
	ctx := context.Background()

	// Create Mocks
	repoMock := mocks.NewStudentBillRepo(t)

	s := &studentBillService{
		repo: repoMock,
	}

	expectedBills := []domain.StudentBill{
		{ID: 1, StudentID: 5, Amount: 100000},
	}

	repoMock.On("FindByStudent", ctx, uint(5)).Return(expectedBills, nil)

	got, err := s.GetByStudent(ctx, 5)

	assert.NoError(t, err)
	assert.Len(t, got, 1)
	assert.Equal(t, uint(1), got[0].ID)
	repoMock.AssertExpectations(t)
}

func Test_studentBillService_GetByParent(t *testing.T) {
	type fields struct {
		db       *bun.DB
		repo     repository.StudentBillRepo
		ruleRepo repository.BillingRuleRepo
		stuRepo  academicrepo.StudentRepo
		notifSvc notificationusecase.FinanceNotificationService
		audit    auditusecase.AuditLogService
	}
	type args struct {
		ctx      context.Context
		parentID uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []domain.StudentBill
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &studentBillService{
				db:       tt.fields.db,
				repo:     tt.fields.repo,
				ruleRepo: tt.fields.ruleRepo,
				stuRepo:  tt.fields.stuRepo,
				notifSvc: tt.fields.notifSvc,
				audit:    tt.fields.audit,
			}
			got, err := s.GetByParent(tt.args.ctx, tt.args.parentID)
			if (err != nil) != tt.wantErr {
				t.Fatalf("studentBillService.GetByParent() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("studentBillService.GetByParent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_studentBillService_GetAll(t *testing.T) {
	type fields struct {
		db       *bun.DB
		repo     repository.StudentBillRepo
		ruleRepo repository.BillingRuleRepo
		stuRepo  academicrepo.StudentRepo
		notifSvc notificationusecase.FinanceNotificationService
		audit    auditusecase.AuditLogService
	}
	type args struct {
		ctx    context.Context
		search string
		sort   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []domain.StudentBill
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &studentBillService{
				db:       tt.fields.db,
				repo:     tt.fields.repo,
				ruleRepo: tt.fields.ruleRepo,
				stuRepo:  tt.fields.stuRepo,
				notifSvc: tt.fields.notifSvc,
				audit:    tt.fields.audit,
			}
			got, err := s.GetAll(tt.args.ctx, tt.args.search, tt.args.sort)
			if (err != nil) != tt.wantErr {
				t.Fatalf("studentBillService.GetAll() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("studentBillService.GetAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_studentBillService_Create(t *testing.T) {
	type fields struct {
		db       *bun.DB
		repo     repository.StudentBillRepo
		ruleRepo repository.BillingRuleRepo
		stuRepo  academicrepo.StudentRepo
		notifSvc notificationusecase.FinanceNotificationService
		audit    auditusecase.AuditLogService
	}
	type args struct {
		ctx context.Context
		sb  *domain.StudentBill
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &studentBillService{
				db:       tt.fields.db,
				repo:     tt.fields.repo,
				ruleRepo: tt.fields.ruleRepo,
				stuRepo:  tt.fields.stuRepo,
				notifSvc: tt.fields.notifSvc,
				audit:    tt.fields.audit,
			}
			if err := s.Create(tt.args.ctx, tt.args.sb); (err != nil) != tt.wantErr {
				t.Errorf("studentBillService.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_studentBillService_Update(t *testing.T) {
	type fields struct {
		db       *bun.DB
		repo     repository.StudentBillRepo
		ruleRepo repository.BillingRuleRepo
		stuRepo  academicrepo.StudentRepo
		notifSvc notificationusecase.FinanceNotificationService
		audit    auditusecase.AuditLogService
	}
	type args struct {
		ctx context.Context
		sb  *domain.StudentBill
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &studentBillService{
				db:       tt.fields.db,
				repo:     tt.fields.repo,
				ruleRepo: tt.fields.ruleRepo,
				stuRepo:  tt.fields.stuRepo,
				notifSvc: tt.fields.notifSvc,
				audit:    tt.fields.audit,
			}
			if err := s.Update(tt.args.ctx, tt.args.sb); (err != nil) != tt.wantErr {
				t.Errorf("studentBillService.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_studentBillService_Delete(t *testing.T) {
	type fields struct {
		db       *bun.DB
		repo     repository.StudentBillRepo
		ruleRepo repository.BillingRuleRepo
		stuRepo  academicrepo.StudentRepo
		notifSvc notificationusecase.FinanceNotificationService
		audit    auditusecase.AuditLogService
	}
	type args struct {
		ctx context.Context
		id  uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &studentBillService{
				db:       tt.fields.db,
				repo:     tt.fields.repo,
				ruleRepo: tt.fields.ruleRepo,
				stuRepo:  tt.fields.stuRepo,
				notifSvc: tt.fields.notifSvc,
				audit:    tt.fields.audit,
			}
			if err := s.Delete(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("studentBillService.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_studentBillService_RunScheduler(t *testing.T) {
	type fields struct {
		db       *bun.DB
		repo     repository.StudentBillRepo
		ruleRepo repository.BillingRuleRepo
		stuRepo  academicrepo.StudentRepo
		notifSvc notificationusecase.FinanceNotificationService
		audit    auditusecase.AuditLogService
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &studentBillService{
				db:       tt.fields.db,
				repo:     tt.fields.repo,
				ruleRepo: tt.fields.ruleRepo,
				stuRepo:  tt.fields.stuRepo,
				notifSvc: tt.fields.notifSvc,
				audit:    tt.fields.audit,
			}
			s.RunScheduler()
		})
	}
}

func Test_studentBillService_autoGenerateMonthlyBills(t *testing.T) {
	type fields struct {
		db       *bun.DB
		repo     repository.StudentBillRepo
		ruleRepo repository.BillingRuleRepo
		stuRepo  academicrepo.StudentRepo
		notifSvc notificationusecase.FinanceNotificationService
		audit    auditusecase.AuditLogService
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &studentBillService{
				db:       tt.fields.db,
				repo:     tt.fields.repo,
				ruleRepo: tt.fields.ruleRepo,
				stuRepo:  tt.fields.stuRepo,
				notifSvc: tt.fields.notifSvc,
				audit:    tt.fields.audit,
			}
			s.autoGenerateMonthlyBills(tt.args.ctx)
		})
	}
}

func Test_studentBillService_SendManualReminder(t *testing.T) {
	type fields struct {
		db       *bun.DB
		repo     repository.StudentBillRepo
		ruleRepo repository.BillingRuleRepo
		stuRepo  academicrepo.StudentRepo
		notifSvc notificationusecase.FinanceNotificationService
		audit    auditusecase.AuditLogService
	}
	type args struct {
		ctx context.Context
		id  uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &studentBillService{
				db:       tt.fields.db,
				repo:     tt.fields.repo,
				ruleRepo: tt.fields.ruleRepo,
				stuRepo:  tt.fields.stuRepo,
				notifSvc: tt.fields.notifSvc,
				audit:    tt.fields.audit,
			}
			if err := s.SendManualReminder(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("studentBillService.SendManualReminder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_studentBillService_GetByID(t *testing.T) {
	type fields struct {
		db       *bun.DB
		repo     repository.StudentBillRepo
		ruleRepo repository.BillingRuleRepo
		stuRepo  academicrepo.StudentRepo
		notifSvc notificationusecase.FinanceNotificationService
		audit    auditusecase.AuditLogService
	}
	type args struct {
		ctx context.Context
		id  uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.StudentBill
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &studentBillService{
				db:       tt.fields.db,
				repo:     tt.fields.repo,
				ruleRepo: tt.fields.ruleRepo,
				stuRepo:  tt.fields.stuRepo,
				notifSvc: tt.fields.notifSvc,
				audit:    tt.fields.audit,
			}
			got, err := s.GetByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Fatalf("studentBillService.GetByID() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("studentBillService.GetByID() = %v, want %v", got, tt.want)
			}
		})
	}
}
