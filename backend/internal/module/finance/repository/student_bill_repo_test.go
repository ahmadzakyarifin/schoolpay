package repository

import (
	"context"
	"reflect"
	"testing"

	"github.com/ahmadzakyarifin/schoolpay/internal/module/finance/domain"
	"github.com/uptrace/bun"
)

func TestNewStudentBillRepo(t *testing.T) {
	type args struct {
		db *bun.DB
	}
	tests := []struct {
		name string
		args args
		want StudentBillRepo
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewStudentBillRepo(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStudentBillRepo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_studentBillRepo_ExistsByPeriod(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx        context.Context
		studentID  uint
		billTypeID uint
		period     string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &studentBillRepo{
				db: tt.fields.db,
			}
			got, err := r.ExistsByPeriod(tt.args.ctx, tt.args.studentID, tt.args.billTypeID, tt.args.period)
			if (err != nil) != tt.wantErr {
				t.Fatalf("studentBillRepo.ExistsByPeriod() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got != tt.want {
				t.Errorf("studentBillRepo.ExistsByPeriod() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_studentBillRepo_ExistsByPeriodExcludeID(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx        context.Context
		studentID  uint
		billTypeID uint
		period     string
		excludeID  uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &studentBillRepo{
				db: tt.fields.db,
			}
			got, err := r.ExistsByPeriodExcludeID(tt.args.ctx, tt.args.studentID, tt.args.billTypeID, tt.args.period, tt.args.excludeID)
			if (err != nil) != tt.wantErr {
				t.Fatalf("studentBillRepo.ExistsByPeriodExcludeID() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got != tt.want {
				t.Errorf("studentBillRepo.ExistsByPeriodExcludeID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_studentBillRepo_Create(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx context.Context
		db  bun.IDB
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
			r := &studentBillRepo{
				db: tt.fields.db,
			}
			if err := r.Create(tt.args.ctx, tt.args.db, tt.args.sb); (err != nil) != tt.wantErr {
				t.Errorf("studentBillRepo.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_studentBillRepo_FindByStudent(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx       context.Context
		studentID uint
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
			r := &studentBillRepo{
				db: tt.fields.db,
			}
			got, err := r.FindByStudent(tt.args.ctx, tt.args.studentID)
			if (err != nil) != tt.wantErr {
				t.Fatalf("studentBillRepo.FindByStudent() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("studentBillRepo.FindByStudent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_studentBillRepo_FindUnpaidBillsByStudent(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx       context.Context
		studentID uint
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
			r := &studentBillRepo{
				db: tt.fields.db,
			}
			got, err := r.FindUnpaidBillsByStudent(tt.args.ctx, tt.args.studentID)
			if (err != nil) != tt.wantErr {
				t.Fatalf("studentBillRepo.FindUnpaidBillsByStudent() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("studentBillRepo.FindUnpaidBillsByStudent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_studentBillRepo_FindByParent(t *testing.T) {
	type fields struct {
		db *bun.DB
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
			r := &studentBillRepo{
				db: tt.fields.db,
			}
			got, err := r.FindByParent(tt.args.ctx, tt.args.parentID)
			if (err != nil) != tt.wantErr {
				t.Fatalf("studentBillRepo.FindByParent() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("studentBillRepo.FindByParent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_studentBillRepo_FindAll(t *testing.T) {
	type fields struct {
		db *bun.DB
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
			r := &studentBillRepo{
				db: tt.fields.db,
			}
			got, err := r.FindAll(tt.args.ctx, tt.args.search, tt.args.sort)
			if (err != nil) != tt.wantErr {
				t.Fatalf("studentBillRepo.FindAll() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("studentBillRepo.FindAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_studentBillRepo_FindByID(t *testing.T) {
	type fields struct {
		db *bun.DB
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
			r := &studentBillRepo{
				db: tt.fields.db,
			}
			got, err := r.FindByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Fatalf("studentBillRepo.FindByID() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("studentBillRepo.FindByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_studentBillRepo_UpdateStatus(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx       context.Context
		db        bun.IDB
		id        uint
		status    string
		totalPaid float64
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
			r := &studentBillRepo{
				db: tt.fields.db,
			}
			if err := r.UpdateStatus(tt.args.ctx, tt.args.db, tt.args.id, tt.args.status, tt.args.totalPaid); (err != nil) != tt.wantErr {
				t.Errorf("studentBillRepo.UpdateStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_studentBillRepo_Update(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx context.Context
		db  bun.IDB
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
			r := &studentBillRepo{
				db: tt.fields.db,
			}
			if err := r.Update(tt.args.ctx, tt.args.db, tt.args.sb); (err != nil) != tt.wantErr {
				t.Errorf("studentBillRepo.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_studentBillRepo_Delete(t *testing.T) {
	type fields struct {
		db *bun.DB
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
			r := &studentBillRepo{
				db: tt.fields.db,
			}
			if err := r.Delete(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("studentBillRepo.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_studentBillRepo_FindForReminder(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx         context.Context
		dueInDays   int
		overdueOnly bool
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
			r := &studentBillRepo{
				db: tt.fields.db,
			}
			got, err := r.FindForReminder(tt.args.ctx, tt.args.dueInDays, tt.args.overdueOnly)
			if (err != nil) != tt.wantErr {
				t.Fatalf("studentBillRepo.FindForReminder() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("studentBillRepo.FindForReminder() = %v, want %v", got, tt.want)
			}
		})
	}
}
