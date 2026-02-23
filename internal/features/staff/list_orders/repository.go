package list_orders

import (
	"context"
	"math"

	"github.com/nanasuryana335/honda-leasing-api/internal/models"
	"gorm.io/gorm"
)

var fullAccessRoles = map[string]bool{
	"SUPER_ADMIN":  true,
	"ADMIN_CABANG": true,
}

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetRoleIDByName(ctx context.Context, roleName string) (int64, error) {
	var role models.Role
	err := r.db.WithContext(ctx).
		Table("account.roles").
		Where("role_name = ?", roleName).
		First(&role).Error
	if err != nil {
		return 0, err
	}
	return role.RoleID, nil
}

type orderRow struct {
	ContractID      int64   `gorm:"column:contract_id"`
	ContractNumber  string  `gorm:"column:contract_number"`
	RequestDate     string  `gorm:"column:request_date"`
	Status          string  `gorm:"column:status"`
	CustomerName    string  `gorm:"column:customer_name"`
	CustomerPhone   string  `gorm:"column:customer_phone"`
	MotorMerk       string  `gorm:"column:motor_merk"`
	MotorType       string  `gorm:"column:motor_type"`
	NilaiKendaraan  float64 `gorm:"column:nilai_kendaraan"`
	DpDibayar       float64 `gorm:"column:dp_dibayar"`
	TenorBulan      int16   `gorm:"column:tenor_bulan"`
	CicilanPerBulan float64 `gorm:"column:cicilan_per_bulan"`
}

func (r *Repository) FindOrders(ctx context.Context, req ListOrdersRequest, roleName string) ([]orderRow, int64, error) {
	// default pagination
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}

	baseQuery := r.db.WithContext(ctx).
		Table("leasing.leasing_contract lc").
		Select(`lc.contract_id, lc.contract_number,
			TO_CHAR(lc.request_date, 'YYYY-MM-DD') AS request_date,
			lc.status,
			c.nama_lengkap AS customer_name,
			c.no_hp       AS customer_phone,
			m.merk        AS motor_merk,
			m.motor_type,
			lc.nilai_kendaraan, lc.dp_dibayar,
			lc.tenor_bulan, lc.cicilan_per_bulan`).
		Joins("JOIN dealer.customers c ON c.customer_id = lc.customer_id").
		Joins("JOIN dealer.motors m ON m.motor_id = lc.motor_id")

	if req.Status != "" {
		baseQuery = baseQuery.Where("lc.status = ?", req.Status)
	}

	if !fullAccessRoles[roleName] {
		roleID, err := r.GetRoleIDByName(ctx, roleName)
		if err != nil {
			return nil, 0, err
		}
		baseQuery = baseQuery.
			Joins("JOIN leasing.leasing_tasks lt ON lt.contract_id = lc.contract_id AND lt.role_id = ?", roleID).
			Group("lc.contract_id, lc.contract_number, lc.request_date, lc.status, c.nama_lengkap, c.no_hp, m.merk, m.motor_type, lc.nilai_kendaraan, lc.dp_dibayar, lc.tenor_bulan, lc.cicilan_per_bulan")
	}

	var total int64
	if err := r.db.WithContext(ctx).
		Table("(?) AS sub", baseQuery).
		Count(&total).Error; err != nil {
		total = 0
	}

	var rows []orderRow
	offset := (req.Page - 1) * req.Limit
	err := baseQuery.
		Order("lc.created_at DESC").
		Limit(req.Limit).
		Offset(offset).
		Scan(&rows).Error
	if err != nil {
		return nil, 0, err
	}

	if total == 0 && len(rows) > 0 {
		total = int64(len(rows))
	}

	return rows, total, nil
}

func (r *Repository) CountOrders(ctx context.Context, req ListOrdersRequest, roleName string) (int64, error) {
	q := r.db.WithContext(ctx).
		Table("leasing.leasing_contract lc").
		Joins("JOIN dealer.customers c ON c.customer_id = lc.customer_id")

	if req.Status != "" {
		q = q.Where("lc.status = ?", req.Status)
	}

	if !fullAccessRoles[roleName] {
		roleID, err := r.GetRoleIDByName(ctx, roleName)
		if err != nil {
			return 0, err
		}
		q = q.Joins("JOIN leasing.leasing_tasks lt ON lt.contract_id = lc.contract_id AND lt.role_id = ?", roleID).
			Group("lc.contract_id")
	}

	var total int64
	err := r.db.WithContext(ctx).
		Table("(?) AS sub", q.Select("lc.contract_id")).
		Count(&total).Error
	return total, err
}

func calcTotalPages(total int64, limit int) int {
	return int(math.Ceil(float64(total) / float64(limit)))
}
