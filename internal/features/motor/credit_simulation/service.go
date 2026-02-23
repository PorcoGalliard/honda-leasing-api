package credit_simulation

import (
	"context"
	"fmt"
	"math"
)

// Tenor options in months
var tenorOptions = []int{23, 29, 35}

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Simulate(ctx context.Context, req CreditSimulationRequest) (*CreditSimulationResponse, error) {
	motor, err := s.repo.FindMotorByID(ctx, req.MotorID)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data motor: %w", err)
	}
	if motor == nil {
		return nil, ErrMotorNotFound
	}

	if req.DP >= motor.HargaOtr {
		return nil, ErrDPExceedsPrice
	}

	pokok := motor.HargaOtr - req.DP

	simulasi := make([]TenorSimulasi, 0, len(tenorOptions))
	for _, tenor := range tenorOptions {
		angsuranPerBulan := pokok / float64(tenor)
		angsuranPerBulan = math.Ceil(angsuranPerBulan/1000) * 1000

		simulasi = append(simulasi, TenorSimulasi{
			Tenor:            tenor,
			TenorLabel:       fmt.Sprintf("%d Bulan", tenor),
			AngsuranPerBulan: angsuranPerBulan,
			TotalBayar:       motor.HargaOtr,
		})
	}

	namaMotor := fmt.Sprintf("%s %s %d", motor.Merk, motor.MotorType, motor.Tahun)

	return &CreditSimulationResponse{
		MotorID:   motor.MotorID,
		NamaMotor: namaMotor,
		HargaOtr:  motor.HargaOtr,
		DP:        req.DP,
		Pokok:     pokok,
		Simulasi:  simulasi,
	}, nil
}
