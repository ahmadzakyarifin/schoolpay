package repository

import (
	"context"
	"fmt"
	"github.com/ahmadzakyarifin/schoolpay/internal/module/academic/domain"
	"github.com/uptrace/bun"
)

type ClassRepo interface {
	Create(ctx context.Context, c *domain.Class) error
	FindAll(ctx context.Context, page, limit int, search, status, majorID, ayID, sort string) ([]domain.Class, int, error)
	FindByID(ctx context.Context, id uint) (*domain.Class, error)
	FindByName(ctx context.Context, name string) (*domain.Class, error)
	Update(ctx context.Context, c *domain.Class) error
	Delete(ctx context.Context, id uint) error
	Restore(ctx context.Context, id uint) error
	ToggleStatus(ctx context.Context, id uint) error
	BulkDelete(ctx context.Context, ids []uint) error
	BulkRestore(ctx context.Context, ids []uint) error
	Exists(ctx context.Context, name string, ayID, majorID uint, excludeID uint) (bool, error)
	CountStudents(ctx context.Context, classID uint) (int, error)
	CountAcademicYears(ctx context.Context, classID uint) (int, error)
}

type classRepo struct {
	db *bun.DB
}

func NewClassRepo(db *bun.DB) ClassRepo {
	return &classRepo{db: db}
}

func (r *classRepo) Create(ctx context.Context, c *domain.Class) error {
	_, err := r.db.NewInsert().Model(c).Exec(ctx)
	return err
}

func (r *classRepo) FindAll(ctx context.Context, page, limit int, search, status, majorID, ayID, sort string) ([]domain.Class, int, error) {
	var list []domain.Class
	q := r.db.NewSelect().Model(&list).
		ColumnExpr("c.*").
		ColumnExpr("j.name as major_name").
		ColumnExpr("ay.year as academic_year_name").
		Join("LEFT JOIN majors j ON c.major_id = j.id").
		Join("LEFT JOIN academic_years ay ON c.academic_year_id = ay.id")

	if status == "trash" {
		q.WhereAllWithDeleted().Where("c.deleted_at IS NOT NULL")
	} else {
		q.Where("c.deleted_at IS NULL")
		switch status {
		case "active":
			q.Where("c.is_active = ?", true)
		case "inactive":
			q.Where("c.is_active = ?", false)
		}
	}

	if search != "" {
		s := "%" + search + "%"
		q.Where("c.name LIKE ? OR j.name LIKE ?", s, s)
	}

	if majorID != "" {
		q.Where("c.major_id = ?", majorID)
	}

	if ayID != "" {
		q.Where("c.academic_year_id = ?", ayID)
	}

	switch sort {
	case "created_asc":
		q.Order("c.created_at ASC", "c.id ASC")
	case "created_desc":
		q.Order("c.created_at DESC", "c.id DESC")
	case "name_asc":
		q.Order("c.name ASC", "c.id ASC")
	case "name_desc":
		q.Order("c.name DESC", "c.id DESC")
	case "grade_asc":
		q.Order("c.grade ASC", "c.name ASC")
	case "grade_desc":
		q.Order("c.grade DESC", "c.name ASC")
	default:
		q.Order("ay.year DESC", "c.grade ASC", "c.name ASC")
	}

	total, err := q.Limit(limit).Offset((page - 1) * limit).ScanAndCount(ctx)
	if err != nil {
		return list, total, err
	}

	if len(list) > 0 {
		var classIDs []uint
		classMap := make(map[uint]int)
		for i, c := range list {
			classIDs = append(classIDs, c.ID)
			classMap[c.ID] = i
			if c.AcademicYearID != nil && *c.AcademicYearID > 0 {
				list[i].AcademicYearIDs = append(list[i].AcademicYearIDs, *c.AcademicYearID)
			}
		}

		// 1. Batch fetch academic_year_classes
		var relations []struct {
			ClassID        uint `bun:"class_id"`
			AcademicYearID uint `bun:"academic_year_id"`
		}
		_ = r.db.NewSelect().Table("academic_year_classes").
			Column("class_id", "academic_year_id").
			Where("class_id IN (?)", bun.In(classIDs)).
			Scan(ctx, &relations)

		for _, rel := range relations {
			if idx, ok := classMap[rel.ClassID]; ok {
				found := false
				for _, existingID := range list[idx].AcademicYearIDs {
					if existingID == rel.AcademicYearID {
						found = true
						break
					}
				}
				if !found {
					list[idx].AcademicYearIDs = append(list[idx].AcademicYearIDs, rel.AcademicYearID)
				}
			}
		}

		for i := range list {
			list[i].AcademicYearCount = len(list[i].AcademicYearIDs)
		}

		// 2. Batch count active students
		var studentCounts []struct {
			ClassID uint `bun:"class_id"`
			Count   int  `bun:"cnt"`
		}
		_ = r.db.NewSelect().Table("students").
			Column("class_id").
			ColumnExpr("count(*) AS cnt").
			Where("class_id IN (?) AND status = 'active' AND deleted_at IS NULL", bun.In(classIDs)).
			Group("class_id").
			Scan(ctx, &studentCounts)

		for _, sc := range studentCounts {
			if idx, ok := classMap[sc.ClassID]; ok {
				list[idx].StudentCount = sc.Count
			}
		}
	}

	return list, total, nil
}

func (r *classRepo) FindByID(ctx context.Context, id uint) (*domain.Class, error) {
	var c domain.Class
	err := r.db.NewSelect().Model(&c).
		ColumnExpr("c.*").
		ColumnExpr("j.name as major_name").
		Join("LEFT JOIN majors j ON c.major_id = j.id").
		Where("c.id = ?", id).Scan(ctx)
	return &c, err
}

func (r *classRepo) FindByName(ctx context.Context, name string) (*domain.Class, error) {
	var c domain.Class
	err := r.db.NewSelect().Model(&c).Where("name = ?", name).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *classRepo) Update(ctx context.Context, c *domain.Class) error {
	_, err := r.db.NewUpdate().Model(c).WherePK().Exec(ctx)
	return err
}

func (r *classRepo) Delete(ctx context.Context, id uint) error {
	res, err := r.db.NewDelete().Model((*domain.Class)(nil)).Where("id = ? AND deleted_at IS NULL", id).Exec(ctx)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("kelas tidak ditemukan atau sudah terhapus")
	}
	return nil
}

func (r *classRepo) Restore(ctx context.Context, id uint) error {
	res, err := r.db.NewUpdate().
		Model((*domain.Class)(nil)).
		WhereAllWithDeleted().
		Where("id = ? AND deleted_at IS NOT NULL", id).
		Set("deleted_at = NULL").
		Exec(ctx)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("kelas tidak ditemukan di riwayat penghapusan")
	}
	return nil
}

func (r *classRepo) ToggleStatus(ctx context.Context, id uint) error {
	_, err := r.db.NewUpdate().
		Model((*domain.Class)(nil)).
		Set("is_active = NOT is_active").
		Where("id = ?", id).
		Exec(ctx)
	return err
}

func (r *classRepo) BulkDelete(ctx context.Context, ids []uint) error {
	res, err := r.db.NewDelete().Model((*domain.Class)(nil)).Where("id IN (?) AND deleted_at IS NULL", bun.In(ids)).Exec(ctx)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("tidak ada data yang dapat dihapus")
	}
	return nil
}

func (r *classRepo) BulkRestore(ctx context.Context, ids []uint) error {
	res, err := r.db.NewUpdate().
		Model((*domain.Class)(nil)).
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

func (r *classRepo) Exists(ctx context.Context, name string, ayID, majorID uint, excludeID uint) (bool, error) {
	q := r.db.NewSelect().Model((*domain.Class)(nil)).
		Where("name = ?", name)

	if ayID > 0 {
		q.Where("academic_year_id = ?", ayID)
	}
	if majorID > 0 {
		q.Where("major_id = ?", majorID)
	}

	if excludeID > 0 {
		q.Where("id != ?", excludeID)
	}
	return q.Exists(ctx)
}

func (r *classRepo) CountStudents(ctx context.Context, classID uint) (int, error) {
	return r.db.NewSelect().
		Model((*domain.Student)(nil)).
		Where("class_id = ? AND status = 'active' AND deleted_at IS NULL", classID).
		Count(ctx)
}

func (r *classRepo) CountAcademicYears(ctx context.Context, classID uint) (int, error) {
	return r.db.NewSelect().
		Model((*domain.AcademicYearClass)(nil)).
		Join("JOIN academic_years ay ON ay.id = ayc.academic_year_id").
		Where("ayc.class_id = ? AND ay.deleted_at IS NULL", classID).
		Count(ctx)
}
