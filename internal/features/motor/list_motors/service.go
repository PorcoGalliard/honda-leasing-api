package list_motors

import (
	"context"
	"fmt"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) ListMotors(ctx context.Context, req ListMotorsRequest) (*ListMotorsResponse, error) {
	// Set default values
	if req.Page == 0 {
		req.Page = 1
	}
	if req.Limit == 0 {
		req.Limit = 10
	}

	if req.GroupByType {
		return s.listMotorsGrouped(ctx, req)
	}

	return s.listMotorsFlat(ctx, req)
}

func (s *Service) listMotorsFlat(ctx context.Context, req ListMotorsRequest) (*ListMotorsResponse, error) {
	motors, totalCount, err := s.repo.FindMotors(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch motors: %w", err)
	}

	motorIDs := make([]int64, len(motors))
	for i, motor := range motors {
		motorIDs[i] = motor.MotorID
	}

	imagesMap, err := s.repo.GetMotorImages(ctx, motorIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch motor images: %w", err)
	}

	motorItems := make([]MotorItem, len(motors))
	for i, motor := range motors {
		motorTypeName, _ := s.repo.GetMotorTypeNameByType(ctx, motor.MotorType)

		images := imagesMap[motor.MotorID]
		if images == nil {
			images = []string{}
		}

		motorItems[i] = MotorItem{
			MotorID:       motor.MotorID,
			Merk:          motor.Merk,
			MotorType:     motor.MotorType,
			MotorTypeName: motorTypeName,
			Tahun:         motor.Tahun,
			Warna:         motor.Warna,
			NomorRangka:   motor.NomorRangka,
			NomorMesin:    motor.NomorMesin,
			CcMesin:       motor.CcMesin,
			NomorPolisi:   motor.NomorPolisi,
			StatusUnit:    motor.StatusUnit,
			HargaOtr:      motor.HargaOtr,
			Images:        images,
		}
	}

	totalPages := int(totalCount) / req.Limit
	if int(totalCount)%req.Limit > 0 {
		totalPages++
	}

	return &ListMotorsResponse{
		Motors: motorItems,
		Pagination: &Pagination{
			Page:       req.Page,
			Limit:      req.Limit,
			TotalRows:  totalCount,
			TotalPages: totalPages,
		},
		TotalMotors: len(motorItems),
	}, nil
}

func (s *Service) listMotorsGrouped(ctx context.Context, req ListMotorsRequest) (*ListMotorsResponse, error) {
	groupedMotors, err := s.repo.FindMotorsGroupedByType(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch grouped motors: %w", err)
	}

	var allMotorIDs []int64
	for _, motors := range groupedMotors {
		for _, motor := range motors {
			allMotorIDs = append(allMotorIDs, motor.MotorID)
		}
	}

	imagesMap, err := s.repo.GetMotorImages(ctx, allMotorIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch motor images: %w", err)
	}

	countsByType, err := s.repo.CountMotorsByType(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to count motors by type: %w", err)
	}

	var motorTypeGroups []MotorTypeGroup
	totalMotors := 0

	motorTypeOrder := []string{"Sport", "Matic", "Classic", "Maxi", "Bebek"}

	for _, motorType := range motorTypeOrder {
		motors, exists := groupedMotors[motorType]
		if !exists || len(motors) == 0 {
			continue
		}

		motorTypeName, _ := s.repo.GetMotorTypeNameByType(ctx, motorType)

		motorItems := make([]MotorItem, len(motors))
		for i, motor := range motors {
			images := imagesMap[motor.MotorID]
			if images == nil {
				images = []string{}
			}

			motorItems[i] = MotorItem{
				MotorID:       motor.MotorID,
				Merk:          motor.Merk,
				MotorType:     motor.MotorType,
				MotorTypeName: motorTypeName,
				Tahun:         motor.Tahun,
				Warna:         motor.Warna,
				NomorRangka:   motor.NomorRangka,
				NomorMesin:    motor.NomorMesin,
				CcMesin:       motor.CcMesin,
				NomorPolisi:   motor.NomorPolisi,
				StatusUnit:    motor.StatusUnit,
				HargaOtr:      motor.HargaOtr,
				Images:        images,
			}
		}

		if req.SortBy == "harga_otr" {
			motorItems = s.sortMotorItems(motorItems, req.OrderBy)
		}

		count := countsByType[motorType]
		motorTypeGroups = append(motorTypeGroups, MotorTypeGroup{
			MotorType:     motorType,
			MotorTypeName: motorTypeName,
			MotorCount:    int(count),
			Motors:        motorItems,
		})

		totalMotors += len(motorItems)
	}

	return &ListMotorsResponse{
		MotorsByType: motorTypeGroups,
		TotalMotors:  totalMotors,
	}, nil
}

func (s *Service) sortMotorItems(items []MotorItem, orderBy string) []MotorItem {
	if orderBy == "asc" {
		for i := 0; i < len(items)-1; i++ {
			for j := i + 1; j < len(items); j++ {
				if items[i].HargaOtr > items[j].HargaOtr {
					items[i], items[j] = items[j], items[i]
				}
			}
		}
	} else {
		for i := 0; i < len(items)-1; i++ {
			for j := i + 1; j < len(items); j++ {
				if items[i].HargaOtr < items[j].HargaOtr {
					items[i], items[j] = items[j], items[i]
				}
			}
		}
	}
	return items
}
