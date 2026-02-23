package create_order

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/nanasuryana335/honda-leasing-api/internal/domain/query/dealer"
	"github.com/nanasuryana335/honda-leasing-api/internal/domain/query/leasing"
	"github.com/nanasuryana335/honda-leasing-api/internal/domain/query/mst"
	"github.com/nanasuryana335/honda-leasing-api/internal/models"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindCustomerByUserID(ctx context.Context, userID int64) (*models.Customer, error) {
	q := dealer.Use(r.db)
	c := q.Customer
	customer, err := c.WithContext(ctx).
		Where(c.UserID.Eq(userID)).
		First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return customer, nil
}

func (r *Repository) FindMotorByID(ctx context.Context, motorID int64) (*models.Motor, error) {
	q := dealer.Use(r.db)
	m := q.Motor
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

func (r *Repository) FindProductByTenor(ctx context.Context, tenor int16) (*models.LeasingProduct, error) {
	var product models.LeasingProduct
	err := r.db.WithContext(ctx).
		Table("leasing.leasing_product").
		Where("tenor_bulan >= ?", tenor).
		Order("tenor_bulan ASC").
		First(&product).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = r.db.WithContext(ctx).
				Table("leasing.leasing_product").
				Order("tenor_bulan DESC").
				First(&product).Error
		}
		if err != nil {
			return nil, err
		}
	}
	return &product, nil
}

func (r *Repository) GetTemplateTasks(ctx context.Context) ([]*models.TemplateTask, error) {
	q := mst.Use(r.db)
	return q.TemplateTask.WithContext(ctx).
		Order(q.TemplateTask.SequenceNo).
		Find()
}

func (r *Repository) GenerateContractNumber(ctx context.Context) (string, error) {
	var count int64
	year := time.Now().Year()
	err := r.db.WithContext(ctx).
		Table("leasing.leasing_contract").
		Where("EXTRACT(YEAR FROM created_at) = ?", year).
		Count(&count).Error
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("KTR-%d-%04d", year, count+1), nil
}

func (r *Repository) CreateOrderTx(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return r.db.WithContext(ctx).Transaction(fn)
}

func CreateLocation(ctx context.Context, tx *gorm.DB, lat, lng float64) (int64, error) {
	loc := &models.Location{
		Latitude:  lat,
		Longitude: lng,
	}
	err := tx.WithContext(ctx).
		Table("mst.locations").
		Omit("kel_id").
		Create(loc).Error
	if err != nil {
		return 0, err
	}
	return loc.LocationID, nil
}

func CreateCustomer(ctx context.Context, tx *gorm.DB, userID int64, name, phone, nik string, locationID int64) (*models.Customer, error) {
	customer := &models.Customer{
		Nik:         nik,
		NamaLengkap: name,
		NoHp:        phone,
		LocationID:  locationID,
		UserID:      userID,
	}
	err := tx.WithContext(ctx).
		Table("dealer.customers").
		Omit("tanggal_lahir").
		Create(customer).Error
	if err != nil {
		return nil, err
	}
	return customer, nil
}

func FetchTasksByContractID(ctx context.Context, tx *gorm.DB, contractID int64) ([]*models.LeasingTask, error) {
	lq := leasing.Use(tx)
	return lq.LeasingTask.WithContext(ctx).
		Where(lq.LeasingTask.ContractID.Eq(contractID)).
		Order(lq.LeasingTask.SequenceNo).
		Find()
}

func CreateContract(ctx context.Context, tx *gorm.DB, contract *models.LeasingContract) error {
	return tx.WithContext(ctx).
		Table("leasing.leasing_contract").
		Omit("tanggal_akad", "tanggal_mulai_cicil").
		Create(contract).Error
}

func CreateLeasingTasksFromTemplates(ctx context.Context, tx *gorm.DB, contractID int64, templates []*models.TemplateTask) error {
	tasks := make([]*models.LeasingTask, 0, len(templates))
	for _, t := range templates {
		tasks = append(tasks, &models.LeasingTask{
			TaskName:   t.TetaName,
			Status:     "inprogress",
			ContractID: contractID,
			RoleID:     t.TetaRoleID,
			SequenceNo: t.SequenceNo,
		})
	}

	return tx.WithContext(ctx).
		Table("leasing.leasing_tasks").
		Create(&tasks).Error
}

func UpdateMotorStatus(ctx context.Context, tx *gorm.DB, motorID int64, status string) error {
	dq := dealer.Use(tx)
	m := dq.Motor
	_, err := m.WithContext(ctx).
		Where(m.MotorID.Eq(motorID)).
		Update(m.StatusUnit, status)
	return err
}
