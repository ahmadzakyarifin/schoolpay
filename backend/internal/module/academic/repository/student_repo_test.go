package repository

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/ahmadzakyarifin/schoolpay/internal/module/academic/domain"
	userauthdomain "github.com/ahmadzakyarifin/schoolpay/internal/module/user_auth/domain"
	"github.com/uptrace/bun"
)

func TestNewStudentRepo(t *testing.T) {
	type args struct {
		db *bun.DB
	}
	tests := []struct {
		name string
		args args
		want StudentRepo
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewStudentRepo(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStudentRepo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_studentRepo_Create(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx context.Context
		db  bun.IDB
		s   *domain.Student
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
			r := &studentRepo{
				db: tt.fields.db,
			}
			if err := r.Create(tt.args.ctx, tt.args.db, tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("studentRepo.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_studentRepo_LinkParent(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx       context.Context
		db        bun.IDB
		studentID uint
		parentID  uint
		relation  string
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
			r := &studentRepo{
				db: tt.fields.db,
			}
			if err := r.LinkParent(tt.args.ctx, tt.args.db, tt.args.studentID, tt.args.parentID, tt.args.relation); (err != nil) != tt.wantErr {
				t.Errorf("studentRepo.LinkParent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_studentRepo_FindAll(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []domain.Student
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &studentRepo{
				db: tt.fields.db,
			}
			got, err := r.FindAll(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Fatalf("studentRepo.FindAll() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("studentRepo.FindAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_studentRepo_FindAllPaginated(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx       context.Context
		page      int
		limit     int
		search    string
		filter    string
		status    string
		entryYear int
		classID   uint
		majorID   uint
		sort      string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []domain.Student
		want1   int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &studentRepo{
				db: tt.fields.db,
			}
			got, got1, err := r.FindAllPaginated(tt.args.ctx, tt.args.page, tt.args.limit, tt.args.search, tt.args.filter, tt.args.status, tt.args.entryYear, tt.args.classID, tt.args.majorID, tt.args.sort)
			if (err != nil) != tt.wantErr {
				t.Fatalf("studentRepo.FindAllPaginated() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("studentRepo.FindAllPaginated() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("studentRepo.FindAllPaginated() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_studentRepo_GetDistinctEntryYears(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &studentRepo{
				db: tt.fields.db,
			}
			got, err := r.GetDistinctEntryYears(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Fatalf("studentRepo.GetDistinctEntryYears() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("studentRepo.GetDistinctEntryYears() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_studentRepo_ToggleStatus(t *testing.T) {
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
			r := &studentRepo{
				db: tt.fields.db,
			}
			if err := r.ToggleStatus(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("studentRepo.ToggleStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_studentRepo_FindByIdentifiers(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.Student
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &studentRepo{
				db: tt.fields.db,
			}
			got, err := r.FindByIdentifiers(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Fatalf("studentRepo.FindByIdentifiers() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("studentRepo.FindByIdentifiers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_studentRepo_FindByEmail(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx   context.Context
		email string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.Student
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &studentRepo{
				db: tt.fields.db,
			}
			got, err := r.FindByEmail(tt.args.ctx, tt.args.email)
			if (err != nil) != tt.wantErr {
				t.Fatalf("studentRepo.FindByEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("studentRepo.FindByEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_studentRepo_FindByPhone(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx   context.Context
		phone string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.Student
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &studentRepo{
				db: tt.fields.db,
			}
			got, err := r.FindByPhone(tt.args.ctx, tt.args.phone)
			if (err != nil) != tt.wantErr {
				t.Fatalf("studentRepo.FindByPhone() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("studentRepo.FindByPhone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_studentRepo_FindByNISN(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx  context.Context
		nisn string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.Student
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &studentRepo{
				db: tt.fields.db,
			}
			got, err := r.FindByNISN(tt.args.ctx, tt.args.nisn)
			if (err != nil) != tt.wantErr {
				t.Fatalf("studentRepo.FindByNISN() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("studentRepo.FindByNISN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_studentRepo_FindByNIS(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx context.Context
		nis string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.Student
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &studentRepo{
				db: tt.fields.db,
			}
			got, err := r.FindByNIS(tt.args.ctx, tt.args.nis)
			if (err != nil) != tt.wantErr {
				t.Fatalf("studentRepo.FindByNIS() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("studentRepo.FindByNIS() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_studentRepo_FindByNIK(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx context.Context
		nik string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.Student
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &studentRepo{
				db: tt.fields.db,
			}
			got, err := r.FindByNIK(tt.args.ctx, tt.args.nik)
			if (err != nil) != tt.wantErr {
				t.Fatalf("studentRepo.FindByNIK() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("studentRepo.FindByNIK() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_studentRepo_GetParents(t *testing.T) {
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
		want    []userauthdomain.User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &studentRepo{
				db: tt.fields.db,
			}
			got, err := r.GetParents(tt.args.ctx, tt.args.studentID)
			if (err != nil) != tt.wantErr {
				t.Fatalf("studentRepo.GetParents() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("studentRepo.GetParents() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_studentRepo_GetStudentsByParentID(t *testing.T) {
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
		want    []domain.Student
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &studentRepo{
				db: tt.fields.db,
			}
			got, err := r.GetStudentsByParentID(tt.args.ctx, tt.args.parentID)
			if (err != nil) != tt.wantErr {
				t.Fatalf("studentRepo.GetStudentsByParentID() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("studentRepo.GetStudentsByParentID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_studentRepo_Update(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx context.Context
		db  bun.IDB
		s   *domain.Student
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
			r := &studentRepo{
				db: tt.fields.db,
			}
			if err := r.Update(tt.args.ctx, tt.args.db, tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("studentRepo.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_studentRepo_Delete(t *testing.T) {
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
			r := &studentRepo{
				db: tt.fields.db,
			}
			if err := r.Delete(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("studentRepo.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_studentRepo_FindIDsByTarget(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx        context.Context
		targetType string
		targetID   uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []uint
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &studentRepo{
				db: tt.fields.db,
			}
			got, err := r.FindIDsByTarget(tt.args.ctx, tt.args.targetType, tt.args.targetID)
			if (err != nil) != tt.wantErr {
				t.Fatalf("studentRepo.FindIDsByTarget() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("studentRepo.FindIDsByTarget() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_studentRepo_CountActive(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &studentRepo{
				db: tt.fields.db,
			}
			got, err := r.CountActive(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Fatalf("studentRepo.CountActive() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got != tt.want {
				t.Errorf("studentRepo.CountActive() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_studentRepo_CountActiveByPeriod(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx          context.Context
		start        *time.Time
		end          *time.Time
		academicYear int
		classID      uint
		majorID      uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &studentRepo{
				db: tt.fields.db,
			}
			got, err := r.CountActiveByPeriod(tt.args.ctx, tt.args.start, tt.args.end, tt.args.academicYear, tt.args.classID, tt.args.majorID)
			if (err != nil) != tt.wantErr {
				t.Fatalf("studentRepo.CountActiveByPeriod() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got != tt.want {
				t.Errorf("studentRepo.CountActiveByPeriod() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_studentRepo_FindByID(t *testing.T) {
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
		want    *domain.Student
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &studentRepo{
				db: tt.fields.db,
			}
			got, err := r.FindByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Fatalf("studentRepo.FindByID() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("studentRepo.FindByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_studentRepo_CountByGender(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx          context.Context
		academicYear int
		classID      uint
		majorID      uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &studentRepo{
				db: tt.fields.db,
			}
			got, err := r.CountByGender(tt.args.ctx, tt.args.academicYear, tt.args.classID, tt.args.majorID)
			if (err != nil) != tt.wantErr {
				t.Fatalf("studentRepo.CountByGender() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("studentRepo.CountByGender() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_studentRepo_CountByStatus(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx          context.Context
		academicYear int
		classID      uint
		majorID      uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &studentRepo{
				db: tt.fields.db,
			}
			got, err := r.CountByStatus(tt.args.ctx, tt.args.academicYear, tt.args.classID, tt.args.majorID)
			if (err != nil) != tt.wantErr {
				t.Fatalf("studentRepo.CountByStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("studentRepo.CountByStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_studentRepo_CountByMajor(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx          context.Context
		academicYear int
		classID      uint
		majorID      uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &studentRepo{
				db: tt.fields.db,
			}
			got, err := r.CountByMajor(tt.args.ctx, tt.args.academicYear, tt.args.classID, tt.args.majorID)
			if (err != nil) != tt.wantErr {
				t.Fatalf("studentRepo.CountByMajor() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("studentRepo.CountByMajor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_studentRepo_CountByClass(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx          context.Context
		academicYear int
		classID      uint
		majorID      uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &studentRepo{
				db: tt.fields.db,
			}
			got, err := r.CountByClass(tt.args.ctx, tt.args.academicYear, tt.args.classID, tt.args.majorID)
			if (err != nil) != tt.wantErr {
				t.Fatalf("studentRepo.CountByClass() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("studentRepo.CountByClass() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_studentRepo_CountByYear(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[int]int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &studentRepo{
				db: tt.fields.db,
			}
			got, err := r.CountByYear(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Fatalf("studentRepo.CountByYear() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("studentRepo.CountByYear() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_studentRepo_GetGenderDemographicsByClass(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx          context.Context
		academicYear int
		classID      uint
		majorID      uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []map[string]interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &studentRepo{
				db: tt.fields.db,
			}
			got, err := r.GetGenderDemographicsByClass(tt.args.ctx, tt.args.academicYear, tt.args.classID, tt.args.majorID)
			if (err != nil) != tt.wantErr {
				t.Fatalf("studentRepo.GetGenderDemographicsByClass() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("studentRepo.GetGenderDemographicsByClass() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_studentRepo_GetGenderDemographicsByMajor(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx          context.Context
		academicYear int
		classID      uint
		majorID      uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []map[string]interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &studentRepo{
				db: tt.fields.db,
			}
			got, err := r.GetGenderDemographicsByMajor(tt.args.ctx, tt.args.academicYear, tt.args.classID, tt.args.majorID)
			if (err != nil) != tt.wantErr {
				t.Fatalf("studentRepo.GetGenderDemographicsByMajor() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("studentRepo.GetGenderDemographicsByMajor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_studentRepo_GetClassHistory(t *testing.T) {
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
		want    []domain.ClassHistory
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &studentRepo{
				db: tt.fields.db,
			}
			got, err := r.GetClassHistory(tt.args.ctx, tt.args.studentID)
			if (err != nil) != tt.wantErr {
				t.Fatalf("studentRepo.GetClassHistory() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("studentRepo.GetClassHistory() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_studentRepo_AddClassHistory(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx       context.Context
		db        bun.IDB
		studentID uint
		classID   uint
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
			r := &studentRepo{
				db: tt.fields.db,
			}
			if err := r.AddClassHistory(tt.args.ctx, tt.args.db, tt.args.studentID, tt.args.classID); (err != nil) != tt.wantErr {
				t.Errorf("studentRepo.AddClassHistory() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_studentRepo_UpdateActiveHistory(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx       context.Context
		db        bun.IDB
		studentID uint
		classID   uint
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
			r := &studentRepo{
				db: tt.fields.db,
			}
			if err := r.UpdateActiveHistory(tt.args.ctx, tt.args.db, tt.args.studentID, tt.args.classID); (err != nil) != tt.wantErr {
				t.Errorf("studentRepo.UpdateActiveHistory() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_studentRepo_Restore(t *testing.T) {
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
			r := &studentRepo{
				db: tt.fields.db,
			}
			if err := r.Restore(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("studentRepo.Restore() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_studentRepo_BulkRestore(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx context.Context
		ids []uint
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
			r := &studentRepo{
				db: tt.fields.db,
			}
			if err := r.BulkRestore(tt.args.ctx, tt.args.ids); (err != nil) != tt.wantErr {
				t.Errorf("studentRepo.BulkRestore() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
