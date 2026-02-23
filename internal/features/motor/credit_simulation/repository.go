package credit_simulation

import (
	"context"
	"errors"

	"github.com/nanasuryana335/honda-leasing-api/internal/domain/query/dealer"
	"github.com/nanasuryana335/honda-leasing-api/internal/models"
	"gorm.io/gorm"
)

type Repository struct {
	db    *gorm.DB
	query *dealer.Query
}

func NewRepository(db *gorm.DB) *Repository {
	query := dealer.Use(db)
	return &Repository{
		db:    db,
		query: query,
	}
}

// FindMotorByID - Get motor detail by ID
func (r *Repository) FindMotorByID(ctx context.Context, motorID int64) (*models.Motor, error) {
	m := r.query.Motor

	motor, err := m.WithContext(ctx).
		Where(m.MotorID.Eq(motorID)).
		First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return motor, nil
}
