package list_motors

import (
	"context"

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

func (r *Repository) FindMotors(ctx context.Context, req ListMotorsRequest) ([]*models.Motor, int64, error) {
	m := r.query.Motor

	query := m.WithContext(ctx)

	// Apply filters
	if req.MotorType != "" {
		query = query.Where(m.MotorType.Eq(req.MotorType))
	}

	if req.Status != "" {
		query = query.Where(m.StatusUnit.Eq(req.Status))
	}

	if req.MinPrice > 0 {
		query = query.Where(m.HargaOtr.Gte(req.MinPrice))
	}

	if req.MaxPrice > 0 {
		query = query.Where(m.HargaOtr.Lte(req.MaxPrice))
	}

	count, err := query.Count()
	if err != nil {
		return nil, 0, err
	}

	// Apply ordering
	query = r.applyOrdering(query, req)

	// Apply pagination (kalau tidak grouping)
	if !req.GroupByType {
		if req.Page > 0 && req.Limit > 0 {
			offset := (req.Page - 1) * req.Limit
			query = query.Offset(offset).Limit(req.Limit)
		}
	}

	motors, err := query.Find()
	return motors, count, err
}

func (r *Repository) FindMotorsGroupedByType(ctx context.Context, req ListMotorsRequest) (map[string][]*models.Motor, error) {
	motors, _, err := r.FindMotors(ctx, req)
	if err != nil {
		return nil, err
	}

	grouped := make(map[string][]*models.Motor)
	for _, motor := range motors {
		grouped[motor.MotorType] = append(grouped[motor.MotorType], motor)
	}

	return grouped, nil
}

func (r *Repository) GetMotorTypeNames(ctx context.Context) (map[int64]string, error) {
	mt := r.query.MotorType
	motorTypes, err := mt.WithContext(ctx).Find()
	if err != nil {
		return nil, err
	}

	typeMap := make(map[int64]string)
	for _, t := range motorTypes {
		typeMap[t.MotyID] = t.MotyName
	}

	return typeMap, nil
}

func (r *Repository) GetMotorImages(ctx context.Context, motorIDs []int64) (map[int64][]string, error) {
	if len(motorIDs) == 0 {
		return make(map[int64][]string), nil
	}

	type Result struct {
		MotorID int64  `gorm:"column:moas_motor_id"`
		FileURL string `gorm:"column:file_url"`
	}

	var results []Result
	err := r.db.WithContext(ctx).
		Table("dealer.motor_assets").
		Select("moas_motor_id, file_url").
		Where("moas_motor_id IN ?", motorIDs).
		Where("file_type IN ?", []string{"png", "jpg", "jpeg"}).
		Order("created_at ASC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	imagesMap := make(map[int64][]string)
	for _, r := range results {
		imagesMap[r.MotorID] = append(imagesMap[r.MotorID], r.FileURL)
	}

	return imagesMap, nil
}

func (r *Repository) GetMotorTypeName(ctx context.Context, motyID int64) (string, error) {
	mt := r.query.MotorType
	motorType, err := mt.WithContext(ctx).
		Where(mt.MotyID.Eq(motyID)).
		First()

	if err != nil {
		return "", err
	}

	return motorType.MotyName, nil
}

func (r *Repository) GetMotorTypeNameByType(ctx context.Context, motorType string) (string, error) {

	mt := r.query.MotorType
	motorTypeModel, err := mt.WithContext(ctx).
		Where(mt.MotyName.Eq(motorType)).
		First()

	if err != nil {
		return motorType, nil
	}

	return motorTypeModel.MotyName, nil
}

func (r *Repository) CountMotorsByType(ctx context.Context, req ListMotorsRequest) (map[string]int64, error) {
	type CountResult struct {
		MotorType string `gorm:"column:motor_type"`
		Count     int64  `gorm:"column:count"`
	}

	query := r.db.WithContext(ctx).
		Table("dealer.motors").
		Select("motor_type, COUNT(*) as count").
		Group("motor_type")

	if req.Status != "" {
		query = query.Where("status_unit = ?", req.Status)
	}

	if req.MinPrice > 0 {
		query = query.Where("harga_otr >= ?", req.MinPrice)
	}

	if req.MaxPrice > 0 {
		query = query.Where("harga_otr <= ?", req.MaxPrice)
	}

	var results []CountResult
	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	countMap := make(map[string]int64)
	for _, r := range results {
		countMap[r.MotorType] = r.Count
	}

	return countMap, nil
}

func (r *Repository) applyOrdering(query dealer.IMotorDo, req ListMotorsRequest) dealer.IMotorDo {
	m := r.query.Motor

	// Default ordering
	sortBy := req.SortBy
	if sortBy == "" {
		sortBy = "created_at"
	}

	orderBy := req.OrderBy
	if orderBy == "" {
		orderBy = "desc"
	}

	// Apply ordering based on sortBy field
	switch sortBy {
	case "motor_type":
		if orderBy == "asc" {
			query = query.Order(m.MotorType)
		} else {
			query = query.Order(m.MotorType.Desc())
		}
	case "harga_otr":
		if orderBy == "asc" {
			query = query.Order(m.HargaOtr)
		} else {
			query = query.Order(m.HargaOtr.Desc())
		}
	case "created_at":
		if orderBy == "asc" {
			query = query.Order(m.CreatedAt)
		} else {
			query = query.Order(m.CreatedAt.Desc())
		}
	default:
		query = query.Order(m.CreatedAt.Desc())
	}

	return query
}
