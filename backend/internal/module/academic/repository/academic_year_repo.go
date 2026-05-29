package repository

import (
	"context"
	"fmt"
	"github.com/ahmadzakyarifin/schoolpay/internal/module/academic/domain"
	"github.com/uptrace/bun"
)

type AcademicYearRepo interface {
	Create(ctx context.Context, ay *domain.AcademicYear) error
	FindAll(ctx context.Context, page, limit int, search, status, sort string) ([]domain.AcademicYear, int, error)
	FindByID(ctx context.Context, id uint) (*domain.AcademicYear, error)
	FindByYear(ctx context.Context, year int) (*domain.AcademicYear, error)
	Update(ctx context.Context, ay *domain.AcademicYear) error
	Delete(ctx context.Context, id uint) error
	Restore(ctx context.Context, id uint) error
	BulkDelete(ctx context.Context, ids []uint) error
	BulkRestore(ctx context.Context, ids []uint) error
	Exists(ctx context.Context, year int, excludeID uint) (bool, error)
	AssignMajors(ctx context.Context, ayID uint, majorIDs []uint) error
	GetMajorsByYear(ctx context.Context, ayID uint) ([]domain.Major, error)
	AssignClasses(ctx context.Context, ayID uint, classIDs []uint) error
	GetClassesByYear(ctx context.Context, ayID uint) ([]domain.Class, error)
	FindActive(ctx context.Context) ([]domain.AcademicYear, error)
	CountStudentsByEntryYear(ctx context.Context, year int) (int, error)
}

type academicYearRepo struct {
	db *bun.DB
}

func NewAcademicYearRepo(db *bun.DB) AcademicYearRepo {
	return &academicYearRepo{db: db}
}

func (r *academicYearRepo) Create(ctx context.Context, ay *domain.AcademicYear) error {
	_, err := r.db.NewInsert().Model(ay).Exec(ctx)
	return err
}

func (r *academicYearRepo) FindAll(ctx context.Context, page, limit int, search, status, sort string) ([]domain.AcademicYear, int, error) {
	var list []domain.AcademicYear
	q := r.db.NewSelect().Model(&list)

	if status == "trash" {
		q.WhereAllWithDeleted().Where("ay.deleted_at IS NOT NULL")
	} else {
		q.Where("ay.deleted_at IS NULL")
		switch status {
		case "active":
			q.Where("ay.is_active = ?", true)
		case "inactive":
			q.Where("ay.is_active = ?", false)
		}
	}

	if search != "" {
		q.Where("CAST(ay.year AS CHAR) LIKE ?", "%"+search+"%")
	}

	switch sort {
	case "year_asc":
		q.Order("ay.year ASC", "ay.id ASC")
	case "created_desc":
		q.Order("ay.created_at DESC", "ay.id DESC")
	case "created_asc":
		q.Order("ay.created_at ASC", "ay.id ASC")
	default:
		q.Order("ay.year DESC", "ay.id DESC")
	}

	total, err := q.Limit(limit).Offset((page - 1) * limit).ScanAndCount(ctx)
	if err != nil {
		return nil, 0, err
	}

	// Post-processing to get counts and major names
	for i := range list {
		majorCount, _ := r.db.NewSelect().Model((*domain.AcademicYearMajor)(nil)).Where("academic_year_id = ?", list[i].ID).Count(ctx)
		classCount, _ := r.db.NewSelect().Model((*domain.AcademicYearClass)(nil)).Where("academic_year_id = ?", list[i].ID).Count(ctx)
		studentCount, _ := r.CountStudentsByEntryYear(ctx, list[i].Year)

		var majorNames []string
		var majorIDs []uint
		_ = r.db.NewSelect().
			Model((*domain.Major)(nil)).
			Column("m.name", "m.id").
			Join("JOIN academic_year_majors aym ON aym.major_id = m.id").
			Where("aym.academic_year_id = ?", list[i].ID).
			Scan(ctx, &majorNames, &majorIDs)

		var classNames []string
		var classIDs []uint
		_ = r.db.NewSelect().
			Model((*domain.Class)(nil)).
			Column("c.name", "c.id").
			Join("JOIN academic_year_classes ayc ON ayc.class_id = c.id").
			Where("ayc.academic_year_id = ?", list[i].ID).
			Scan(ctx, &classNames, &classIDs)

		list[i].MajorCount = int(majorCount)
		list[i].ClassCount = int(classCount)
		list[i].StudentCount = int(studentCount)
		list[i].MajorNames = majorNames
		list[i].ClassNames = classNames
		list[i].MajorIDs = majorIDs
		list[i].ClassIDs = classIDs
	}

	return list, total, nil
}

func (r *academicYearRepo) FindByID(ctx context.Context, id uint) (*domain.AcademicYear, error) {
	var ay domain.AcademicYear
	err := r.db.NewSelect().Model(&ay).Where("id = ?", id).Scan(ctx)
	return &ay, err
}

func (r *academicYearRepo) FindByYear(ctx context.Context, year int) (*domain.AcademicYear, error) {
	var ay domain.AcademicYear
	err := r.db.NewSelect().Model(&ay).Where("year = ?", year).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &ay, nil
}

func (r *academicYearRepo) Update(ctx context.Context, ay *domain.AcademicYear) error {
	_, err := r.db.NewUpdate().Model(ay).WherePK().Exec(ctx)
	return err
}

func (r *academicYearRepo) Delete(ctx context.Context, id uint) error {
	res, err := r.db.NewDelete().Model((*domain.AcademicYear)(nil)).Where("id = ? AND deleted_at IS NULL", id).Exec(ctx)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("gagal: Angkatan tidak ditemukan atau sudah terhapus")
	}
	return nil
}

func (r *academicYearRepo) Restore(ctx context.Context, id uint) error {
	res, err := r.db.NewUpdate().Model((*domain.AcademicYear)(nil)).WhereAllWithDeleted().Set("deleted_at = NULL").Where("id = ? AND deleted_at IS NOT NULL", id).Exec(ctx)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("gagal: Angkatan tidak ditemukan di riwayat penghapusan atau sudah aktif")
	}
	return nil
}

func (r *academicYearRepo) BulkDelete(ctx context.Context, ids []uint) error {
	res, err := r.db.NewDelete().Model((*domain.AcademicYear)(nil)).Where("id IN (?) AND deleted_at IS NULL", bun.In(ids)).Exec(ctx)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("gagal: Tidak ada data angkatan yang bisa dihapus (mungkin sudah terhapus)")
	}
	return nil
}

func (r *academicYearRepo) BulkRestore(ctx context.Context, ids []uint) error {
	res, err := r.db.NewUpdate().Model((*domain.AcademicYear)(nil)).WhereAllWithDeleted().Set("deleted_at = NULL").Where("id IN (?) AND deleted_at IS NOT NULL", bun.In(ids)).Exec(ctx)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("gagal: Tidak ada data angkatan yang bisa dipulihkan (mungkin sudah aktif)")
	}
	return nil
}

func (r *academicYearRepo) Exists(ctx context.Context, year int, excludeID uint) (bool, error) {
	q := r.db.NewSelect().Model((*domain.AcademicYear)(nil)).Where("year = ?", year)
	if excludeID > 0 {
		q.Where("id != ?", excludeID)
	}
	return q.Exists(ctx)
}

func (r *academicYearRepo) AssignMajors(ctx context.Context, ayID uint, majorIDs []uint) error {
	return r.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		// Remove existing
		_, err := tx.NewDelete().Model((*domain.AcademicYearMajor)(nil)).Where("academic_year_id = ?", ayID).Exec(ctx)
		if err != nil {
			return err
		}

		if len(majorIDs) == 0 {
			return nil
		}

		// Pengecekan apakah semua jurusan yang dipilih berstatus aktif dan belum dihapus
		activeCount, err := tx.NewSelect().Model((*domain.Major)(nil)).
			Where("id IN (?) AND is_active = ? AND deleted_at IS NULL", bun.In(majorIDs), true).
			Count(ctx)
		if err != nil {
			return err
		}
		if activeCount != len(majorIDs) {
			return fmt.Errorf("gagal: Terdapat jurusan yang dipilih yang tidak aktif atau tidak valid")
		}

		// Insert new
		var list []domain.AcademicYearMajor
		for _, mid := range majorIDs {
			list = append(list, domain.AcademicYearMajor{AcademicYearID: ayID, MajorID: mid})
		}
		_, err = tx.NewInsert().Model(&list).Exec(ctx)
		return err
	})
}

func (r *academicYearRepo) GetMajorsByYear(ctx context.Context, ayID uint) ([]domain.Major, error) {
	var list []domain.Major
	err := r.db.NewSelect().
		Model(&list).
		Join("JOIN academic_year_majors aym ON aym.major_id = m.id").
		Where("aym.academic_year_id = ?", ayID).
		Scan(ctx)
	return list, err
}

func (r *academicYearRepo) AssignClasses(ctx context.Context, ayID uint, classIDs []uint) error {
	return r.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		// Remove existing
		_, err := tx.NewDelete().Model((*domain.AcademicYearClass)(nil)).Where("academic_year_id = ?", ayID).Exec(ctx)
		if err != nil {
			return err
		}

		if len(classIDs) == 0 {
			return nil
		}

		// Pengecekan apakah semua kelas yang dipilih berstatus aktif dan belum dihapus
		activeCount, err := tx.NewSelect().Model((*domain.Class)(nil)).
			Where("id IN (?) AND is_active = ? AND deleted_at IS NULL", bun.In(classIDs), true).
			Count(ctx)
		if err != nil {
			return err
		}
		if activeCount != len(classIDs) {
			return fmt.Errorf("gagal: Terdapat kelas yang dipilih yang tidak aktif atau tidak valid")
		}

		// Insert new
		var list []domain.AcademicYearClass
		for _, cid := range classIDs {
			list = append(list, domain.AcademicYearClass{AcademicYearID: ayID, ClassID: cid})
		}
		_, err = tx.NewInsert().Model(&list).Exec(ctx)
		return err
	})
}

func (r *academicYearRepo) GetClassesByYear(ctx context.Context, ayID uint) ([]domain.Class, error) {
	var list []domain.Class
	err := r.db.NewSelect().
		Model(&list).
		Join("JOIN academic_year_classes ayc ON ayc.class_id = c.id").
		Where("ayc.academic_year_id = ?", ayID).
		Scan(ctx)
	return list, err
}

func (r *academicYearRepo) FindActive(ctx context.Context) ([]domain.AcademicYear, error) {
	var list []domain.AcademicYear
	err := r.db.NewSelect().Model(&list).Where("is_active = 1").Order("year DESC").Scan(ctx)
	return list, err
}

func (r *academicYearRepo) CountStudentsByEntryYear(ctx context.Context, year int) (int, error) {
	return r.db.NewSelect().
		Model((*domain.Student)(nil)).
		Where("entry_year = ? AND status = 'active' AND deleted_at IS NULL", year).
		Count(ctx)
}
