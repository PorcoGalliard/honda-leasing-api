package login

import (
	"context"

	"github.com/nanasuryana335/honda-leasing-api/internal/domain/query/account"
	"github.com/nanasuryana335/honda-leasing-api/internal/models"
	"gorm.io/gorm"
)

type Repository struct {
	db    *gorm.DB
	query *account.Query
}

func NewRepository(db *gorm.DB) *Repository {
	query := account.Use(db)
	return &Repository{
		db:    db,
		query: query,
	}
}

func (r *Repository) FindUserByPhone(ctx context.Context, phoneNumber string) (*models.User, error) {
	u := r.query.User
	user, err := u.WithContext(ctx).
		Where(u.PhoneNumber.Eq(phoneNumber)).
		Where(u.IsActive.Is(true)).
		First()

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *Repository) GetUserRoles(ctx context.Context, userID int64) ([]string, error) {
	var roles []string

	err := r.db.WithContext(ctx).
		Table("account.user_roles").
		Select("account.roles.role_name").
		Joins("JOIN account.roles ON account.user_roles.role_id = account.roles.role_id").
		Where("account.user_roles.user_id = ?", userID).
		Scan(&roles).Error

	return roles, err
}

func (r *Repository) UpdateLastLogin(ctx context.Context, userID int64) error {
	u := r.query.User
	_, err := u.WithContext(ctx).
		Where(u.UserID.Eq(userID)).
		Update(u.LastLogin, "NOW()")
	return err
}

func (r *Repository) IncrementFailedAttempts(ctx context.Context, userID int64) error {
	return r.db.WithContext(ctx).
		Exec("UPDATE account.users SET failed_attempts = failed_attempts + 1 WHERE user_id = ?", userID).
		Error
}

func (r *Repository) ResetFailedAttempts(ctx context.Context, userID int64) error {
	u := r.query.User
	_, err := u.WithContext(ctx).
		Where(u.UserID.Eq(userID)).
		Update(u.FailedAttempts, 0)
	return err
}
