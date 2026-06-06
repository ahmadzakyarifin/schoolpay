package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/ahmadzakyarifin/schoolpay/internal/module/user_auth/domain"
	"github.com/ahmadzakyarifin/schoolpay/internal/module/user_auth/repository/model"
	"github.com/uptrace/bun"
)

type UserRepo interface {
	Create(ctx context.Context, db bun.IDB, user *domain.User) error
	FindAll(ctx context.Context, role string) ([]domain.User, error)
	FindByID(ctx context.Context, id uint) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindByPhone(ctx context.Context, phone string) (*domain.User, error)
	FindByNIK(ctx context.Context, nik string) (*domain.User, error)
	Update(ctx context.Context, db bun.IDB, user *domain.User) error
	Delete(ctx context.Context, id uint) error
	UpdatePassword(ctx context.Context, id uint, passwordHash string) error
	Activate(ctx context.Context, id uint) error
	ToggleStatus(ctx context.Context, id uint) error
	FindPaginated(ctx context.Context, page, limit int, search, role, filter, status, sort string) ([]domain.User, int, error)
	CountActive(ctx context.Context) (int, error)
	CountActiveByPeriod(ctx context.Context, start, end *time.Time) (int, error)
	FindParentsByStudentID(ctx context.Context, studentID uint) ([]domain.User, error)
	BulkDelete(ctx context.Context, ids []uint) error
	Restore(ctx context.Context, id uint) error
	BulkRestore(ctx context.Context, ids []uint) error
}

type userRepo struct {
	db *bun.DB
}

func NewUserRepo(db *bun.DB) UserRepo {
	return &userRepo{db: db}
}

func (r *userRepo) Create(ctx context.Context, db bun.IDB, u *domain.User) error {
	m := model.FromDomain(u)
	_, err := db.NewInsert().Model(m).Exec(ctx)
	if err == nil {
		u.ID = m.ID
		u.CreatedAt = m.CreatedAt
		u.UpdatedAt = m.UpdatedAt
	}
	return err
}

func (r *userRepo) FindAll(ctx context.Context, role string) ([]domain.User, error) {
	var models []model.UserModel
	q := r.db.NewSelect().Model(&models)
	if role != "" {
		q.Where("role = ?", role)
	}
	err := q.Scan(ctx)
	if err != nil {
		return nil, err
	}
	users := make([]domain.User, len(models))
	for i, m := range models {
		users[i] = *m.ToDomain()
	}
	return users, nil
}

func (r *userRepo) FindByID(ctx context.Context, id uint) (*domain.User, error) {
	var m model.UserModel
	err := r.db.NewSelect().
		Model(&m).
		ColumnExpr("u.*").
		ColumnExpr("(SELECT COUNT(*) FROM students s WHERE s.parent_id = u.id AND s.deleted_at IS NULL) as student_count").
		ColumnExpr("(SELECT GROUP_CONCAT(CONCAT(s.id, '::', s.name) ORDER BY s.name ASC SEPARATOR '||') FROM students s WHERE s.parent_id = u.id AND s.deleted_at IS NULL) as student_names").
		Where("u.id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return m.ToDomain(), nil
}

func (r *userRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var m model.UserModel
	err := r.db.NewSelect().Model(&m).WhereAllWithDeleted().Where("email = ?", email).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return m.ToDomain(), nil
}

func (r *userRepo) FindByPhone(ctx context.Context, phone string) (*domain.User, error) {
	var m model.UserModel
	err := r.db.NewSelect().Model(&m).Where("phone_number = ?", phone).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return m.ToDomain(), nil
}

func (r *userRepo) FindByNIK(ctx context.Context, nik string) (*domain.User, error) {
	var m model.UserModel
	err := r.db.NewSelect().Model(&m).WhereAllWithDeleted().Where("nik = ?", nik).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return m.ToDomain(), nil
}

func (r *userRepo) Update(ctx context.Context, db bun.IDB, u *domain.User) error {
	m := model.FromDomain(u)
	_, err := db.NewUpdate().Model(m).WherePK().Exec(ctx)
	return err
}

func (r *userRepo) Delete(ctx context.Context, id uint) error {
	res, err := r.db.NewDelete().Model((*model.UserModel)(nil)).Where("id = ? AND deleted_at IS NULL", id).Exec(ctx)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("pengguna sudah terhapus atau tidak ditemukan")
	}
	return nil
}

func (r *userRepo) BulkDelete(ctx context.Context, ids []uint) error {
	res, err := r.db.NewDelete().Model((*model.UserModel)(nil)).Where("id IN (?) AND deleted_at IS NULL", bun.In(ids)).Exec(ctx)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("tidak ada data yang perlu dihapus")
	}
	return nil
}

func (r *userRepo) Restore(ctx context.Context, id uint) error {
	res, err := r.db.NewUpdate().
		Model((*model.UserModel)(nil)).
		Set("deleted_at = NULL").
		WhereAllWithDeleted().
		Where("id = ? AND deleted_at IS NOT NULL", id).
		Exec(ctx)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("data tidak ditemukan di tempat sampah atau sudah aktif")
	}
	return nil
}

func (r *userRepo) BulkRestore(ctx context.Context, ids []uint) error {
	res, err := r.db.NewUpdate().
		Model((*model.UserModel)(nil)).
		Set("deleted_at = NULL").
		WhereAllWithDeleted().
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

func (r *userRepo) UpdatePassword(ctx context.Context, id uint, hash string) error {
	res, err := r.db.NewUpdate().
		Model((*model.UserModel)(nil)).
		Set("password_hash = ?", hash).
		Set("is_active = 1").
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("pengguna tidak ditemukan")
	}
	return nil
}

func (r *userRepo) Activate(ctx context.Context, id uint) error {
	_, err := r.db.NewUpdate().
		Model((*model.UserModel)(nil)).
		Set("is_active = 1").
		Where("id = ?", id).
		Exec(ctx)
	return err
}

func (r *userRepo) ToggleStatus(ctx context.Context, id uint) error {
	_, err := r.db.NewUpdate().
		Model((*model.UserModel)(nil)).
		Set("is_active = NOT is_active").
		Where("id = ?", id).
		Exec(ctx)
	return err
}

func (r *userRepo) FindPaginated(ctx context.Context, page, limit int, search, role, filter, status, sort string) ([]domain.User, int, error) {
	var models []model.UserModel

	q := r.db.NewSelect().
		Model(&models).
		ColumnExpr("u.*").
		ColumnExpr("(SELECT COUNT(*) FROM students st WHERE st.parent_id = u.id AND st.deleted_at IS NULL) as student_count").
		ColumnExpr("(SELECT GROUP_CONCAT(CONCAT(st.id, '::', st.name) ORDER BY st.name ASC SEPARATOR '||') FROM students st WHERE st.parent_id = u.id AND st.deleted_at IS NULL) as student_names")

	switch status {
	case "trash":
		q.WhereAllWithDeleted().Where("u.deleted_at IS NOT NULL")
	case "active":
		q.Where("u.is_active = 1")
	case "inactive":
		q.Where("u.is_active = 0")
	}

	if search != "" {
		searchQuery := "%" + search + "%"
		q.WhereGroup(" AND ", func(sq *bun.SelectQuery) *bun.SelectQuery {
			return sq.Where("u.name LIKE ?", searchQuery).
				WhereOr("u.email LIKE ?", searchQuery).
				WhereOr("u.phone_number LIKE ?", searchQuery).
				WhereOr("EXISTS (SELECT 1 FROM students st WHERE st.parent_id = u.id AND st.name LIKE ? AND st.deleted_at IS NULL)", searchQuery)
		})
	}

	if role != "" {
		q.Where("u.role = ?", role)
	}

	switch filter {
	case "no_child":
		q.Where("u.role = 'parent'").Where("NOT EXISTS (SELECT 1 FROM students st WHERE st.parent_id = u.id AND st.deleted_at IS NULL)")
	case "has_child":
		q.Where("u.role = 'parent'").Where("EXISTS (SELECT 1 FROM students st WHERE st.parent_id = u.id AND st.deleted_at IS NULL)")
	}

	switch sort {
	case "name_asc":
		q.Order("u.name ASC")
	case "name_desc":
		q.Order("u.name DESC")
	case "created_asc":
		q.Order("u.created_at ASC")
	case "created_desc":
		q.Order("u.created_at DESC")
	default:
		q.Order("u.created_at DESC")
	}

	count, err := q.Limit(limit).
		Offset((page - 1) * limit).
		ScanAndCount(ctx)

	if err != nil {
		return nil, 0, err
	}

	users := make([]domain.User, len(models))
	for i, m := range models {
		users[i] = *m.ToDomain()
	}

	return users, count, nil
}

func (r *userRepo) CountActive(ctx context.Context) (int, error) {
	return r.db.NewSelect().Model((*model.UserModel)(nil)).Where("is_active = 1").Count(ctx)
}

func (r *userRepo) CountActiveByPeriod(ctx context.Context, start, end *time.Time) (int, error) {
	q := r.db.NewSelect().Model((*model.UserModel)(nil))
	if start != nil {
		q.Where("created_at >= ?", start)
	}
	if end != nil {
		q.Where("created_at <= ?", end)
	}
	return q.Count(ctx)
}

func (r *userRepo) FindParentsByStudentID(ctx context.Context, studentID uint) ([]domain.User, error) {
	var parents []model.UserModel
	err := r.db.NewSelect().
		Model(&parents).
		Join("JOIN students s ON u.id = s.parent_id").
		Where("s.id = ?", studentID).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	users := make([]domain.User, len(parents))
	for i, m := range parents {
		users[i] = *m.ToDomain()
	}
	return users, nil
}
