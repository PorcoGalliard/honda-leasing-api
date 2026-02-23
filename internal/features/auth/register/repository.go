package register

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
	return &Repository{
		db:    db,
		query: account.Use(db),
	}
}

func (r *Repository) IsPhoneRegistered(ctx context.Context, phoneNumber string) (bool, error) {
	u := r.query.User
	count, err := u.WithContext(ctx).
		Where(u.PhoneNumber.Eq(phoneNumber)).
		Count()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *Repository) CreateUser(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).
		Table("account.users").
		Omit("username").
		Create(user).Error
}

func (r *Repository) FindRoleByName(ctx context.Context, roleName string) (*models.Role, error) {
	ro := r.query.Role
	return ro.WithContext(ctx).
		Where(ro.RoleName.Eq(roleName)).
		First()
}

func (r *Repository) AssignRole(ctx context.Context, userRole *models.UserRole) error {
	return r.db.WithContext(ctx).
		Table("account.user_roles").
		Omit("assigned_by").
		Create(userRole).Error
}
