package repository

import (
	"context"
	"fmt"
	"github.com/ahmadzakyarifin/schoolpay/internal/module/academic/domain"
	"github.com/uptrace/bun"
)

type MajorRepo interface {
	Create(ctx context.Context, j *domain.Major) error
	FindAll(ctx context.Context, page, limit int, search, status, sort string) ([]domain.Major, int, error)
	FindByID(ctx context.Context, id uint) (*domain.Major, error)
	FindByName(ctx context.Context, name string) (*domain.Major, error)
	Update(ctx context.Context, j *domain.Major) error
	Delete(ctx context.Context, id uint) error
	Restore(ctx context.Context, id uint) error
	ToggleStatus(ctx context.Context, id uint) error
	BulkDelete(ctx context.Context, ids []uint) error
	BulkRestore(ctx context.Context, ids []uint) error
	Exists(ctx context.Context, name string, excludeID uint) (bool, error)
	ExistsByCode(ctx context.Context, code string, excludeID uint) (bool, error)
	CountClasses(ctx context.Context, majorID uint) (int, error)
	CountStudents(ctx context.Context, majorID uint) (int, error)
	CountAcademicYears(ctx context.Context, majorID uint) (int, error)
}

type majorRepo struct {
	db *bun.DB
}

func NewMajorRepo(db *bun.DB) MajorRepo {
	return &majorRepo{db: db}
}

func (r *majorRepo) Create(ctx context.Context, j *domain.Major) error {
	_, err := r.db.NewInsert().Model(j).Exec(ctx)
	return err
}

func (r *majorRepo) FindAll(ctx context.Context, page, limit int, search, status, sort string) ([]domain.Major, int, error) {
	var list []domain.Major
	q := r.db.NewSelect().Model(&list)

	if status == "trash" {
		q.WhereAllWithDeleted().Where("deleted_at IS NOT NULL")
	} else {
		q.Where("deleted_at IS NULL")
		switch status {
		case "active":
			q.Where("is_active = ?", true)
		case "inactive":
			q.Where("is_active = ?", false)
		case "has_class":
			q.Where("EXISTS (SELECT 1 FROM classes WHERE major_id = major.id AND deleted_at IS NULL)")
		case "no_class":
			q.Where("NOT EXISTS (SELECT 1 FROM classes WHERE major_id = major.id AND deleted_at IS NULL)")
		}
	}

	if search != "" {
		s := "%" + search + "%"
		q.Where("name LIKE ? OR code LIKE ?", s, s)
	}
	switch sort {
	case "created_asc":
		q.Order("created_at ASC", "id ASC")
	case "created_desc":
		q.Order("created_at DESC", "id DESC")
	case "name_desc":
		q.Order("name DESC", "id DESC")
	case "code_asc":
		q.Order("code ASC", "name ASC")
	case "code_desc":
		q.Order("code DESC", "name ASC")
	default:
		q.Order("name ASC", "id ASC")
	}

	total, err := q.Limit(limit).Offset((page - 1) * limit).ScanAndCount(ctx)
	if err != nil {
		return list, total, err
	}

	if len(list) > 0 {
		var majorIDs []uint
		majorMap := make(map[uint]int)
		for i, m := range list {
			majorIDs = append(majorIDs, m.ID)
			majorMap[m.ID] = i
		}

		var relations []struct {
			MajorID        uint `bun:"major_id"`
			AcademicYearID uint `bun:"academic_year_id"`
		}
		_ = r.db.NewSelect().Table("academic_year_majors").
			Column("major_id", "academic_year_id").
			Where("major_id IN (?)", bun.In(majorIDs)).
			Scan(ctx, &relations)

		for _, rel := range relations {
			if idx, ok := majorMap[rel.MajorID]; ok {
				list[idx].YearIDs = append(list[idx].YearIDs, rel.AcademicYearID)
				list[idx].AcademicYearCount++
			}
		}

		// Batch fetch classes belonging to these majors
		var classes []struct {
			MajorID uint   `bun:"major_id"`
			Name    string `bun:"name"`
		}
		_ = r.db.NewSelect().Table("classes").
			Column("major_id", "name").
			Where("major_id IN (?) AND is_active = ? AND deleted_at IS NULL", bun.In(majorIDs), true).
			Order("grade ASC", "name ASC").
			Scan(ctx, &classes)

		for _, c := range classes {
			if idx, ok := majorMap[c.MajorID]; ok {
				list[idx].ClassNames = append(list[idx].ClassNames, c.Name)
				list[idx].ClassCount++
			}
		}

		// Batch fetch students belonging to these majors
		var studentCounts []struct {
			MajorID uint `bun:"major_id"`
			Count   int  `bun:"count"`
		}
		_ = r.db.NewSelect().TableExpr("students AS s").
			ColumnExpr("COALESCE(c.major_id, s.major_id) AS major_id").
			ColumnExpr("COUNT(s.id) AS count").
			Join("LEFT JOIN classes c ON s.class_id = c.id").
			Where("s.status = 'active' AND s.deleted_at IS NULL").
			Where("COALESCE(c.major_id, s.major_id) IN (?)", bun.In(majorIDs)).
			GroupExpr("COALESCE(c.major_id, s.major_id)").
			Scan(ctx, &studentCounts)

		for _, sc := range studentCounts {
			if idx, ok := majorMap[sc.MajorID]; ok {
				list[idx].StudentCount = sc.Count
			}
		}
	}

	return list, total, nil
}

func (r *majorRepo) FindByID(ctx context.Context, id uint) (*domain.Major, error) {
	var j domain.Major
	err := r.db.NewSelect().Model(&j).Where("id = ?", id).Scan(ctx)
	return &j, err
}

func (r *majorRepo) FindByName(ctx context.Context, name string) (*domain.Major, error) {
	var j domain.Major
	err := r.db.NewSelect().Model(&j).Where("name = ?", name).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &j, nil
}

func (r *majorRepo) Update(ctx context.Context, j *domain.Major) error {
	_, err := r.db.NewUpdate().Model(j).WherePK().Exec(ctx)
	return err
}

func (r *majorRepo) Delete(ctx context.Context, id uint) error {
	res, err := r.db.NewDelete().Model((*domain.Major)(nil)).Where("id = ? AND deleted_at IS NULL", id).Exec(ctx)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("jurusan tidak ditemukan atau sudah terhapus")
	}
	return nil
}

func (r *majorRepo) Restore(ctx context.Context, id uint) error {
	res, err := r.db.NewUpdate().
		Model((*domain.Major)(nil)).
		WhereAllWithDeleted().
		Where("id = ? AND deleted_at IS NOT NULL", id).
		Set("deleted_at = NULL").
		Exec(ctx)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("jurusan tidak ditemukan di riwayat penghapusan")
	}
	return nil
}

func (r *majorRepo) Exists(ctx context.Context, name string, excludeID uint) (bool, error) {
	q := r.db.NewSelect().Model((*domain.Major)(nil)).Where("name = ?", name)
	if excludeID > 0 {
		q.Where("id != ?", excludeID)
	}
	return q.Exists(ctx)
}

func (r *majorRepo) ExistsByCode(ctx context.Context, code string, excludeID uint) (bool, error) {
	q := r.db.NewSelect().Model((*domain.Major)(nil)).WhereAllWithDeleted().Where("LOWER(code) = LOWER(?)", code)
	if excludeID > 0 {
		q.Where("id != ?", excludeID)
	}
	return q.Exists(ctx)
}

func (r *majorRepo) ToggleStatus(ctx context.Context, id uint) error {
	_, err := r.db.NewUpdate().
		Model((*domain.Major)(nil)).
		Set("is_active = NOT is_active").
		Where("id = ?", id).
		Exec(ctx)
	return err
}

func (r *majorRepo) BulkDelete(ctx context.Context, ids []uint) error {
	res, err := r.db.NewDelete().Model((*domain.Major)(nil)).Where("id IN (?) AND deleted_at IS NULL", bun.In(ids)).Exec(ctx)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("tidak ada data yang dapat dihapus")
	}
	return nil
}

func (r *majorRepo) BulkRestore(ctx context.Context, ids []uint) error {
	res, err := r.db.NewUpdate().
		Model((*domain.Major)(nil)).
		WhereAllWithDeleted().
		Set("deleted_at = NULL").
		Where("id IN (?) AND deleted_at IS NOT NULL", bun.In(ids)).
		Exec(ctx)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("tidak ada data yang perlu dipulihkan")
	}
	return nil
}

func (r *majorRepo) CountClasses(ctx context.Context, majorID uint) (int, error) {
	return r.db.NewSelect().
		Model((*domain.Class)(nil)).
		Where("major_id = ? AND is_active = ? AND deleted_at IS NULL", majorID, true).
		Count(ctx)
}

func (r *majorRepo) CountStudents(ctx context.Context, majorID uint) (int, error) {
	return r.db.NewSelect().
		TableExpr("students AS s").
		Join("LEFT JOIN classes c ON s.class_id = c.id").
		Where("s.status = 'active' AND s.deleted_at IS NULL").
		Where("COALESCE(c.major_id, s.major_id) = ?", majorID).
		Count(ctx)
}

func (r *majorRepo) CountAcademicYears(ctx context.Context, majorID uint) (int, error) {
	return r.db.NewSelect().
		Model((*domain.AcademicYearMajor)(nil)).
		Join("JOIN academic_years ay ON ay.id = aym.academic_year_id").
		Where("aym.major_id = ? AND ay.deleted_at IS NULL", majorID).
		Count(ctx)
}
