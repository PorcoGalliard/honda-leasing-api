package create_order

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/nanasuryana335/honda-leasing-api/internal/models"
	"gorm.io/gorm"
)

const (
	BiayaAdmin = 200_000.0
	Asuransi   = 250_000.0
	Fidusia    = 200_000.0
	Materai    = 10_000.0
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateOrder(ctx context.Context, userID int64, req CreateOrderRequest) (*CreateOrderResponse, error) {
	motor, err := s.repo.FindMotorByID(ctx, req.MotorID)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data motor: %w", err)
	}
	if motor == nil {
		return nil, ErrMotorNotFound
	}
	if motor.StatusUnit != "ready" {
		return nil, ErrMotorNotAvailable
	}

	if req.DP >= motor.HargaOtr {
		return nil, ErrDPExceedsPrice
	}

	existingCustomer, err := s.repo.FindCustomerByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil profil customer: %w", err)
	}
	if existingCustomer == nil && req.NIK == "" {
		return nil, ErrNIKRequired
	}

	product, err := s.repo.FindProductByTenor(ctx, req.Tenor)
	if err != nil {
		return nil, fmt.Errorf("gagal menemukan produk leasing: %w", err)
	}

	pokok := motor.HargaOtr - req.DP
	cicilanPerBulan := math.Ceil(pokok/float64(req.Tenor)/1000) * 1000
	subTotal := motor.HargaOtr - req.DP - req.PromoDiscount + BiayaAdmin + Asuransi + Fidusia + Materai

	reqDate, err := time.Parse("2006-01-02", req.RequestDate)
	if err != nil {
		return nil, fmt.Errorf("format tanggal tidak valid, gunakan YYYY-MM-DD")
	}

	contractNumber, err := s.repo.GenerateContractNumber(ctx)
	if err != nil {
		return nil, fmt.Errorf("gagal generate nomor kontrak: %w", err)
	}

	templates, err := s.repo.GetTemplateTasks(ctx)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil template task: %w", err)
	}

	var contract *models.LeasingContract
	var tasks []*models.LeasingTask

	err = s.repo.CreateOrderTx(ctx, func(tx *gorm.DB) error {
		var customerID int64
		if existingCustomer != nil {
			customerID = existingCustomer.CustomerID
		} else {
			locationID, err := CreateLocation(ctx, tx, req.Latitude, req.Longitude)
			if err != nil {
				return fmt.Errorf("gagal menyimpan lokasi: %w", err)
			}

			newCustomer, err := CreateCustomer(ctx, tx, userID, req.ContactName, req.PhoneNumber, req.NIK, locationID)
			if err != nil {
				return fmt.Errorf("gagal membuat profil customer: %w", err)
			}
			customerID = newCustomer.CustomerID
		}

		contract = &models.LeasingContract{
			ContractNumber:  contractNumber,
			RequestDate:     reqDate,
			TenorBulan:      req.Tenor,
			NilaiKendaraan:  motor.HargaOtr,
			DpDibayar:       req.DP,
			PokokPinjaman:   pokok,
			TotalPinjaman:   subTotal,
			CicilanPerBulan: cicilanPerBulan,
			Status:          "draft",
			CustomerID:      customerID,
			MotorID:         motor.MotorID,
			ProductID:       product.ProductID,
		}
		if err := CreateContract(ctx, tx, contract); err != nil {
			return fmt.Errorf("gagal membuat kontrak: %w", err)
		}

		if err := CreateLeasingTasksFromTemplates(ctx, tx, contract.ContractID, templates); err != nil {
			return fmt.Errorf("gagal membuat task order: %w", err)
		}

		if err := UpdateMotorStatus(ctx, tx, motor.MotorID, "booked"); err != nil {
			return fmt.Errorf("gagal update status motor: %w", err)
		}

		leasingTasks, err := FetchTasksByContractID(ctx, tx, contract.ContractID)
		if err != nil {
			return fmt.Errorf("gagal mengambil tasks: %w", err)
		}
		tasks = leasingTasks

		return nil
	})
	if err != nil {
		return nil, err
	}

	taskItems := make([]OrderTaskItem, 0, len(tasks))
	for _, t := range tasks {
		taskItems = append(taskItems, OrderTaskItem{
			TaskID:     t.TaskID,
			TaskName:   t.TaskName,
			SequenceNo: t.SequenceNo,
			Status:     t.Status,
		})
	}

	return &CreateOrderResponse{
		ContractID:     contract.ContractID,
		ContractNumber: contract.ContractNumber,
		Status:         contract.Status,
		Motor: OrderMotorInfo{
			MotorID:   motor.MotorID,
			Merk:      motor.Merk,
			MotorType: motor.MotorType,
			Tahun:     motor.Tahun,
			HargaOtr:  motor.HargaOtr,
		},
		Tenor: req.Tenor,
		TotalSummary: OrderSummary{
			HargaOtr:        motor.HargaOtr,
			DP:              req.DP,
			PromoDiscount:   req.PromoDiscount,
			BiayaAdmin:      BiayaAdmin,
			Asuransi:        Asuransi,
			Fidusia:         Fidusia,
			Materai:         Materai,
			PokokPinjaman:   pokok,
			CicilanPerBulan: cicilanPerBulan,
			SubTotal:        subTotal,
		},
		Tasks: taskItems,
	}, nil
}
