package get_order_progress

import (
	"context"
	"errors"

	"github.com/nanasuryana335/honda-leasing-api/internal/domain/query/dealer"
	"github.com/nanasuryana335/honda-leasing-api/internal/domain/query/leasing"
	"github.com/nanasuryana335/honda-leasing-api/internal/models"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindContractByID(ctx context.Context, contractID int64) (*models.LeasingContract, error) {
	lq := leasing.Use(r.db)
	lc := lq.LeasingContract
	contract, err := lc.WithContext(ctx).
		Where(lc.ContractID.Eq(contractID)).
		First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return contract, nil
}

func (r *Repository) FindContractByIDAndUserID(ctx context.Context, contractID, userID int64) (*models.LeasingContract, error) {
	dq := dealer.Use(r.db)
	customer, err := dq.Customer.WithContext(ctx).
		Where(dq.Customer.UserID.Eq(userID)).
		First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	lq := leasing.Use(r.db)
	lc := lq.LeasingContract
	contract, err := lc.WithContext(ctx).
		Where(lc.ContractID.Eq(contractID)).
		Where(lc.CustomerID.Eq(customer.CustomerID)).
		First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return contract, nil
}

func (r *Repository) FindTasksByContractID(ctx context.Context, contractID int64) ([]*models.LeasingTask, error) {
	lq := leasing.Use(r.db)
	return lq.LeasingTask.WithContext(ctx).
		Where(lq.LeasingTask.ContractID.Eq(contractID)).
		Order(lq.LeasingTask.SequenceNo).
		Find()
}
