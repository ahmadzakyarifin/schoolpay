package repository

import (
	"context"
	"fmt"
	"github.com/ahmadzakyarifin/schoolpay/internal/module/academic/domain"
	userauthdomain "github.com/ahmadzakyarifin/schoolpay/internal/module/user_auth/domain"
	"github.com/uptrace/bun"
	"time"
)

type StudentRepo interface {
	Create(ctx context.Context, db bun.IDB, s *domain.Student) error
	LinkParent(ctx context.Context, db bun.IDB, studentID uint, parentID uint, relation string) error
	FindAll(ctx context.Context) ([]domain.Student, error)
	FindAllPaginated(ctx context.Context, page, limit int, search, filter, status string, entryYear int, classID, majorID uint, sort string) ([]domain.Student, int, error)
	GetDistinctEntryYears(ctx context.Context) ([]int, error)
	FindByIdentifiers(ctx context.Context, id string) (*domain.Student, error)
	FindByEmail(ctx context.Context, email string) (*domain.Student, error)
	FindByPhone(ctx context.Context, phone string) (*domain.Student, error)
	FindByNISN(ctx context.Context, nisn string) (*domain.Student, error)
	FindByNIS(ctx context.Context, nis string) (*domain.Student, error)
	FindByNIK(ctx context.Context, nik string) (*domain.Student, error)
	GetParents(ctx context.Context, studentID uint) ([]userauthdomain.User, error)
	GetStudentsByParentID(ctx context.Context, parentID uint) ([]domain.Student, error)
	Update(ctx context.Context, db bun.IDB, s *domain.Student) error
	Delete(ctx context.Context, id uint) error
	ToggleStatus(ctx context.Context, id uint) error
	FindIDsByTarget(ctx context.Context, targetType string, targetID uint) ([]uint, error)
	CountActive(ctx context.Context) (int, error)
	CountActiveByPeriod(ctx context.Context, start, end *time.Time, academicYear int, classID, majorID uint) (int, error)
	FindByID(ctx context.Context, id uint) (*domain.Student, error)
	CountByGender(ctx context.Context, academicYear int, classID, majorID uint) (map[string]int, error)
	CountByStatus(ctx context.Context, academicYear int, classID, majorID uint) (map[string]int, error)
	CountByMajor(ctx context.Context, academicYear int, classID, majorID uint) (map[string]int, error)
	CountByClass(ctx context.Context, academicYear int, classID, majorID uint) (map[string]int, error)
	CountByYear(ctx context.Context) (map[int]int, error)
	GetGenderDemographicsByClass(ctx context.Context, academicYear int, classID, majorID uint) ([]map[string]interface{}, error)
	GetGenderDemographicsByMajor(ctx context.Context, academicYear int, classID, majorID uint) ([]map[string]interface{}, error)
	GetClassHistory(ctx context.Context, studentID uint) ([]domain.ClassHistory, error)
	AddClassHistory(ctx context.Context, db bun.IDB, studentID, classID uint) error
	UpdateActiveHistory(ctx context.Context, db bun.IDB, studentID, classID uint) error
	Restore(ctx context.Context, id uint) error
	BulkRestore(ctx context.Context, ids []uint) error
}

type studentRepo struct {
	db *bun.DB
}

func NewStudentRepo(db *bun.DB) StudentRepo {
	return &studentRepo{db: db}
}

func (r *studentRepo) Create(ctx context.Context, db bun.IDB, s *domain.Student) error {
	_, err := db.NewInsert().Model(s).Exec(ctx)
	return err
}

func (r *studentRepo) LinkParent(ctx context.Context, db bun.IDB, studentID uint, parentID uint, relation string) error {
	_, err := db.NewUpdate().
		Model((*domain.Student)(nil)).
		Set("parent_id = ?", parentID).
		Where("id = ?", studentID).
		Exec(ctx)
	return err
}

func (r *studentRepo) FindAll(ctx context.Context) ([]domain.Student, error) {
	var list []domain.Student
	err := r.db.NewSelect().Model(&list).Scan(ctx)
	return list, err
}

func (r *studentRepo) FindAllPaginated(ctx context.Context, page, limit int, search, filter, status string, entryYear int, classID, majorID uint, sort string) ([]domain.Student, int, error) {
	var list []domain.Student

	q := r.db.NewSelect().
		Model(&list).
		ColumnExpr("s.*").
		ColumnExpr("c.name as class_name, j.name as major_name").
		Join("LEFT JOIN classes c ON s.class_id = c.id").
		Join("LEFT JOIN majors j ON (c.major_id = j.id OR s.major_id = j.id)")

	if search != "" {
		searchQuery := "%" + search + "%"
		q.WhereGroup(" AND ", func(sq *bun.SelectQuery) *bun.SelectQuery {
			return sq.Where("s.name LIKE ?", searchQuery).
				WhereOr("s.nisn LIKE ?", searchQuery).
				WhereOr("s.nis LIKE ?", searchQuery).
				WhereOr("s.email LIKE ?", searchQuery)
		})
	}

	if filter == "no_parent" {
		q.Where("s.parent_id IS NULL")
	}

	if status == "trash" {
		q.WhereAllWithDeleted().Where("s.deleted_at IS NOT NULL")
	} else if status != "" {
		q.Where("s.status = ?", status)
	}

	if entryYear > 0 {
		q.Where("s.entry_year = ?", entryYear)
	}

	if classID > 0 {
		q.Where("s.class_id = ?", classID)
	}

	if majorID > 0 {
		q.Where("c.major_id = ? OR s.major_id = ?", majorID, majorID)
	}

	switch sort {
	case "created_desc":
		q.Order("s.created_at DESC", "s.id DESC")
	case "created_asc":
		q.Order("s.created_at ASC", "s.id ASC")
	case "name_desc":
		q.Order("s.name DESC", "s.id DESC")
	case "entry_year_desc":
		q.Order("s.entry_year DESC", "s.name ASC")
	case "entry_year_asc":
		q.Order("s.entry_year ASC", "s.name ASC")
	default:
		q.Order("s.name ASC", "s.id ASC")
	}

	total, err := q.
		Limit(limit).
		Offset((page - 1) * limit).
		ScanAndCount(ctx)

	return list, total, err
}

func (r *studentRepo) GetDistinctEntryYears(ctx context.Context) ([]int, error) {
	var years []int
	err := r.db.NewSelect().
		Model((*domain.Student)(nil)).
		ColumnExpr("DISTINCT entry_year").
		Order("entry_year DESC").
		Scan(ctx, &years)
	return years, err
}

func (r *studentRepo) ToggleStatus(ctx context.Context, id uint) error {
	_, err := r.db.NewUpdate().
		Model((*domain.Student)(nil)).
		Set("status = CASE WHEN status = 'active' THEN 'inactive' WHEN status = 'inactive' THEN 'active' ELSE status END").
		Where("id = ?", id).
		Exec(ctx)
	return err
}

func (r *studentRepo) FindByIdentifiers(ctx context.Context, id string) (*domain.Student, error) {
	var s domain.Student
	err := r.db.NewSelect().
		Model(&s).
		Where("nisn = ? OR nis = ? OR nik = ? OR email = ? OR phone_number = ?", id, id, id, id, id).
		Limit(1).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *studentRepo) FindByEmail(ctx context.Context, email string) (*domain.Student, error) {
	var s domain.Student
	err := r.db.NewSelect().Model(&s).WhereAllWithDeleted().Where("email = ?", email).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *studentRepo) FindByPhone(ctx context.Context, phone string) (*domain.Student, error) {
	var s domain.Student
	err := r.db.NewSelect().Model(&s).WhereAllWithDeleted().Where("phone_number = ?", phone).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *studentRepo) FindByNISN(ctx context.Context, nisn string) (*domain.Student, error) {
	var s domain.Student
	err := r.db.NewSelect().Model(&s).WhereAllWithDeleted().Where("nisn = ?", nisn).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *studentRepo) FindByNIS(ctx context.Context, nis string) (*domain.Student, error) {
	var s domain.Student
	err := r.db.NewSelect().Model(&s).WhereAllWithDeleted().Where("nis = ?", nis).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *studentRepo) FindByNIK(ctx context.Context, nik string) (*domain.Student, error) {
	var s domain.Student
	err := r.db.NewSelect().Model(&s).WhereAllWithDeleted().Where("nik = ?", nik).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *studentRepo) GetParents(ctx context.Context, studentID uint) ([]userauthdomain.User, error) {
	var parents []userauthdomain.User
	err := r.db.NewSelect().
		Model(&parents).
		ColumnExpr("u.*").
		Join("JOIN students s ON s.parent_id = u.id").
		Where("s.id = ?", studentID).
		Group("u.id").
		Scan(ctx)
	return parents, err
}

func (r *studentRepo) GetStudentsByParentID(ctx context.Context, parentID uint) ([]domain.Student, error) {
	var students []domain.Student
	err := r.db.NewSelect().
		Model(&students).
		ColumnExpr("s.*").
		ColumnExpr("c.name as class_name").
		ColumnExpr("m.name as major_name").
		Join("LEFT JOIN classes c ON s.class_id = c.id AND c.deleted_at IS NULL").
		Join("LEFT JOIN majors m ON m.id = COALESCE(c.major_id, s.major_id) AND m.deleted_at IS NULL").
		Where("s.parent_id = ?", parentID).
		Where("s.deleted_at IS NULL").
		Order("s.name ASC").
		Scan(ctx)
	return students, err
}

func (r *studentRepo) Update(ctx context.Context, db bun.IDB, s *domain.Student) error {
	_, err := db.NewUpdate().Model(s).WherePK().Exec(ctx)
	return err
}

func (r *studentRepo) Delete(ctx context.Context, id uint) error {
	res, err := r.db.NewDelete().Model((*domain.Student)(nil)).Where("id = ? AND deleted_at IS NULL", id).Exec(ctx)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("siswa sudah terhapus atau tidak ditemukan")
	}
	return nil
}

func (r *studentRepo) FindIDsByTarget(ctx context.Context, targetType string, targetID uint) ([]uint, error) {
	var ids []uint
	q := r.db.NewSelect().Model((*domain.Student)(nil)).
		ColumnExpr("s.id").
		Join("JOIN classes c ON s.class_id = c.id").
		Join("JOIN majors j ON (c.major_id = j.id OR s.major_id = j.id)").
		Join("JOIN academic_years ay ON s.entry_year = ay.year").
		Where("s.status = 'active' AND s.deleted_at IS NULL").
		Where("c.is_active = ? AND c.deleted_at IS NULL", true).
		Where("j.is_active = ? AND j.deleted_at IS NULL", true).
		Where("ay.is_active = ? AND ay.deleted_at IS NULL", true)

	switch targetType {
	case "class":
		q.Where("s.class_id = ?", targetID)
	case "major":
		q.Where("c.major_id = ? OR s.major_id = ?", targetID, targetID)
	case "academic_year":
		q.Where("s.entry_year = (SELECT year FROM academic_years WHERE id = ?)", targetID)
	case "all":
		// No extra filter
	default:
		return ids, nil
	}

	err := q.Scan(ctx, &ids)
	return ids, err
}

func (r *studentRepo) CountActive(ctx context.Context) (int, error) {
	return r.db.NewSelect().Model((*domain.Student)(nil)).Where("status = 'active'").Count(ctx)
}

func (r *studentRepo) CountActiveByPeriod(ctx context.Context, start, end *time.Time, academicYear int, classID, majorID uint) (int, error) {
	q := r.db.NewSelect().Model((*domain.Student)(nil)).Where("status = 'active'")

	if start != nil {
		q.Where("created_at >= ?", start)
	}
	if end != nil {
		q.Where("created_at <= ?", end)
	}
	if academicYear > 0 {
		if academicYear < 1000 {
			q.Where("entry_year = (SELECT year FROM academic_years WHERE id = ?)", academicYear)
		} else {
			q.Where("entry_year = ?", academicYear)
		}
	}
	if classID > 0 {
		q.Where("class_id = ?", classID)
	}
	if majorID > 0 {
		q.Where("major_id = ?", majorID)
	}

	return q.Count(ctx)
}

func (r *studentRepo) FindByID(ctx context.Context, id uint) (*domain.Student, error) {
	var s domain.Student
	err := r.db.NewSelect().
		Model(&s).
		ColumnExpr("s.*").
		ColumnExpr("c.name as class_name, j.name as major_name, u.name as parent_name").
		Join("LEFT JOIN classes c ON s.class_id = c.id").
		Join("LEFT JOIN majors j ON s.major_id = j.id").
		Join("LEFT JOIN users u ON s.parent_id = u.id").
		Where("s.id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &s, nil
}
func (r *studentRepo) CountByGender(ctx context.Context, academicYear int, classID, majorID uint) (map[string]int, error) {
	counts := make(map[string]int)
	var stats []struct {
		Gender string `bun:"gender"`
		Count  int    `bun:"count"`
	}

	q := r.db.NewSelect().
		Model((*domain.Student)(nil)).
		ColumnExpr("gender, COUNT(*) as count")

	if academicYear > 0 {
		if academicYear < 1000 {
			q.Where("entry_year = (SELECT year FROM academic_years WHERE id = ?)", academicYear)
		} else {
			q.Where("entry_year = ?", academicYear)
		}
	}
	if classID > 0 {
		q.Where("class_id = ?", classID)
	}
	if majorID > 0 {
		q.Where("major_id = ?", majorID)
	}

	err := q.Group("gender").Scan(ctx, &stats)
	if err != nil {
		return counts, err
	}

	for _, s := range stats {
		counts[s.Gender] = s.Count
	}
	return counts, nil
}

func (r *studentRepo) CountByStatus(ctx context.Context, academicYear int, classID, majorID uint) (map[string]int, error) {
	counts := make(map[string]int)
	var stats []struct {
		Status string `bun:"status"`
		Count  int    `bun:"count"`
	}

	q := r.db.NewSelect().
		Model((*domain.Student)(nil)).
		ColumnExpr("status, COUNT(*) as count")

	if academicYear > 0 {
		if academicYear < 1000 {
			q.Where("entry_year = (SELECT year FROM academic_years WHERE id = ?)", academicYear)
		} else {
			q.Where("entry_year = ?", academicYear)
		}
	}
	if classID > 0 {
		q.Where("class_id = ?", classID)
	}
	if majorID > 0 {
		q.Where("major_id = ?", majorID)
	}

	err := q.Group("status").Scan(ctx, &stats)
	if err != nil {
		return counts, err
	}

	for _, s := range stats {
		counts[s.Status] = s.Count
	}
	return counts, nil
}

func (r *studentRepo) CountByMajor(ctx context.Context, academicYear int, classID, majorID uint) (map[string]int, error) {
	counts := make(map[string]int)
	var stats []struct {
		Name  string `bun:"name"`
		Count int    `bun:"count"`
	}

	q := r.db.NewSelect().
		Model((*domain.Student)(nil)).
		ColumnExpr("j.name, COUNT(s.id) as count").
		Join("JOIN majors j ON s.major_id = j.id")

	if academicYear > 0 {
		if academicYear < 1000 {
			q.Where("s.entry_year = (SELECT year FROM academic_years WHERE id = ?)", academicYear)
		} else {
			q.Where("s.entry_year = ?", academicYear)
		}
	}
	if classID > 0 {
		q.Where("s.class_id = ?", classID)
	}
	if majorID > 0 {
		q.Where("s.major_id = ?", majorID)
	}

	err := q.Group("j.id", "j.name").Scan(ctx, &stats)
	if err != nil {
		return counts, err
	}

	for _, s := range stats {
		counts[s.Name] = s.Count
	}
	return counts, nil
}

func (r *studentRepo) CountByClass(ctx context.Context, academicYear int, classID, majorID uint) (map[string]int, error) {
	counts := make(map[string]int)
	var stats []struct {
		Name  string `bun:"name"`
		Count int    `bun:"count"`
	}

	q := r.db.NewSelect().
		Model((*domain.Student)(nil)).
		ColumnExpr("c.name, COUNT(s.id) as count").
		Join("JOIN classes c ON s.class_id = c.id")

	if academicYear > 0 {
		if academicYear < 1000 {
			q.Where("s.entry_year = (SELECT year FROM academic_years WHERE id = ?)", academicYear)
		} else {
			q.Where("s.entry_year = ?", academicYear)
		}
	}
	if classID > 0 {
		q.Where("s.class_id = ?", classID)
	}
	if majorID > 0 {
		q.Where("s.major_id = ?", majorID)
	}

	err := q.Group("c.id", "c.name").Scan(ctx, &stats)
	if err != nil {
		return counts, err
	}

	for _, s := range stats {
		counts[s.Name] = s.Count
	}
	return counts, nil
}

func (r *studentRepo) CountByYear(ctx context.Context) (map[int]int, error) {
	counts := make(map[int]int)
	var stats []struct {
		Year  int `bun:"entry_year"`
		Count int `bun:"count"`
	}

	err := r.db.NewSelect().
		Model((*domain.Student)(nil)).
		ColumnExpr("entry_year, COUNT(*) as count").
		Group("entry_year").
		Order("entry_year DESC").
		Scan(ctx, &stats)

	if err != nil {
		return counts, err
	}

	for _, s := range stats {
		counts[s.Year] = s.Count
	}
	return counts, nil
}
func (r *studentRepo) GetGenderDemographicsByClass(ctx context.Context, academicYear int, classID, majorID uint) ([]map[string]interface{}, error) {
	q := r.db.NewSelect().
		Model((*domain.Student)(nil)).
		ColumnExpr("c.name").
		ColumnExpr("SUM(CASE WHEN s.gender = 'L' THEN 1 ELSE 0 END) as male").
		ColumnExpr("SUM(CASE WHEN s.gender = 'P' THEN 1 ELSE 0 END) as female").
		Join("JOIN classes c ON s.class_id = c.id")

	if academicYear > 0 {
		if academicYear < 1000 {
			q.Where("s.entry_year = (SELECT year FROM academic_years WHERE id = ?)", academicYear)
		} else {
			q.Where("s.entry_year = ?", academicYear)
		}
	}
	if classID > 0 {
		q.Where("s.class_id = ?", classID)
	}
	if majorID > 0 {
		q.Where("s.major_id = ?", majorID)
	}

	var results []map[string]interface{}
	err := q.Group("c.id", "c.name").Scan(ctx, &results)
	return results, err
}

func (r *studentRepo) GetGenderDemographicsByMajor(ctx context.Context, academicYear int, classID, majorID uint) ([]map[string]interface{}, error) {
	q := r.db.NewSelect().
		Model((*domain.Student)(nil)).
		ColumnExpr("j.name").
		ColumnExpr("SUM(CASE WHEN s.gender = 'L' THEN 1 ELSE 0 END) as male").
		ColumnExpr("SUM(CASE WHEN s.gender = 'P' THEN 1 ELSE 0 END) as female").
		Join("JOIN majors j ON s.major_id = j.id")

	if academicYear > 0 {
		if academicYear < 1000 {
			q.Where("s.entry_year = (SELECT year FROM academic_years WHERE id = ?)", academicYear)
		} else {
			q.Where("s.entry_year = ?", academicYear)
		}
	}
	if classID > 0 {
		q.Where("s.class_id = ?", classID)
	}
	if majorID > 0 {
		q.Where("s.major_id = ?", majorID)
	}

	var results []map[string]interface{}
	err := q.Group("j.id", "j.name").Scan(ctx, &results)
	return results, err
}

func (r *studentRepo) GetClassHistory(ctx context.Context, studentID uint) ([]domain.ClassHistory, error) {
	var list []domain.ClassHistory
	err := r.db.NewSelect().
		Model(&list).
		ModelTableExpr("student_classes as sc").
		ColumnExpr("sc.id, sc.student_id, sc.class_id, sc.is_active, sc.created_at").
		ColumnExpr("c.name as class_name, c.grade").
		ColumnExpr(`CASE WHEN MONTH(sc.created_at) >= 7 THEN CONCAT(YEAR(sc.created_at), '/', YEAR(sc.created_at)+1) ELSE CONCAT(YEAR(sc.created_at)-1, '/', YEAR(sc.created_at)) END as academic_year`).
		Join("JOIN classes c ON sc.class_id = c.id").
		Where("sc.student_id = ?", studentID).
		Order("sc.created_at DESC").
		Scan(ctx)
	return list, err
}

func (r *studentRepo) AddClassHistory(ctx context.Context, db bun.IDB, studentID, classID uint) error {
	// 1. Deactivate current active history
	_, _ = db.NewUpdate().
		Model((*domain.StudentClass)(nil)).
		Set("is_active = 0").
		Where("student_id = ? AND is_active = 1", studentID).
		Exec(ctx)

	// 2. Insert new active history
	sc := &domain.StudentClass{
		StudentID: studentID,
		ClassID:   classID,
		IsActive:  true,
	}
	_, err := db.NewInsert().Model(sc).Exec(ctx)
	return err
}

func (r *studentRepo) UpdateActiveHistory(ctx context.Context, db bun.IDB, studentID, classID uint) error {
	_, err := db.NewUpdate().
		Model((*domain.StudentClass)(nil)).
		Set("class_id = ?", classID).
		Where("student_id = ? AND is_active = 1", studentID).
		Exec(ctx)
	return err
}

func (r *studentRepo) Restore(ctx context.Context, id uint) error {
	res, err := r.db.NewUpdate().
		Model((*domain.Student)(nil)).
		Set("deleted_at = NULL").
		Set("status = 'active'").
		WhereAllWithDeleted().
		Where("id = ? AND deleted_at IS NOT NULL", id).
		Exec(ctx)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("siswa tidak ditemukan di tempat sampah atau sudah aktif")
	}
	return nil
}

func (r *studentRepo) BulkRestore(ctx context.Context, ids []uint) error {
	res, err := r.db.NewUpdate().
		Model((*domain.Student)(nil)).
		Set("deleted_at = NULL").
		Set("status = 'active'").
		WhereAllWithDeleted().
		Where("id IN (?) AND deleted_at IS NOT NULL", bun.In(ids)).
		Exec(ctx)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("tidak ada data siswa yang perlu dipulihkan")
	}
	return nil
}
