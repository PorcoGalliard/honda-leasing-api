package register

import (
	"context"
	"fmt"
	"regexp"

	"github.com/nanasuryana335/honda-leasing-api/internal/models"
	sharedErrors "github.com/nanasuryana335/honda-leasing-api/internal/shared/errors"
	"golang.org/x/crypto/bcrypt"
)

const defaultRole = "CUSTOMER"

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Register(ctx context.Context, req RegisterRequest) (*RegisterResponse, error) {
	if err := s.validate(req); err != nil {
		return nil, err
	}

	roleName := defaultRole
	if req.Role != "" {
		if !AllowedRegisterRoles[req.Role] {
			return nil, fmt.Errorf("role '%s' tidak diizinkan untuk registrasi mandiri", req.Role)
		}
		roleName = req.Role
	}

	exists, err := s.repo.IsPhoneRegistered(ctx, req.PhoneNumber)
	if err != nil {
		return nil, fmt.Errorf("gagal melakukan pengecekan nomor HP: %w", err)
	}
	if exists {
		return nil, sharedErrors.ErrDuplicateEntry
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("gagal memproses password: %w", err)
	}

	hashedPIN, err := bcrypt.GenerateFromPassword([]byte(req.PIN), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("gagal memproses PIN: %w", err)
	}

	newUser := &models.User{
		FullName:    req.FullName,
		PhoneNumber: req.PhoneNumber,
		Email:       req.Email,
		Password:    string(hashedPassword),
		PinKey:      string(hashedPIN),
		IsActive:    true,
	}

	if err := s.repo.CreateUser(ctx, newUser); err != nil {
		return nil, fmt.Errorf("gagal menyimpan data user: %w", err)
	}

	role, err := s.repo.FindRoleByName(ctx, roleName)
	if err != nil {
		return nil, fmt.Errorf("gagal mendapatkan data role: %w", err)
	}

	userRole := &models.UserRole{
		UserID: newUser.UserID,
		RoleID: role.RoleID,
	}
	if err := s.repo.AssignRole(ctx, userRole); err != nil {
		return nil, fmt.Errorf("gagal assign role ke user: %w", err)
	}

	return &RegisterResponse{
		UserID:      newUser.UserID,
		FullName:    newUser.FullName,
		PhoneNumber: newUser.PhoneNumber,
		Email:       newUser.Email,
		Role:        roleName,
	}, nil
}

func (s *Service) validate(req RegisterRequest) error {

	phoneRegex := regexp.MustCompile(`^(\+62|62|0)[0-9]{9,12}$`)
	if !phoneRegex.MatchString(req.PhoneNumber) {
		return fmt.Errorf("format nomor HP tidak valid")
	}

	pinRegex := regexp.MustCompile(`^[0-9]{6}$`)
	if !pinRegex.MatchString(req.PIN) {
		return fmt.Errorf("PIN harus 6 digit angka")
	}

	return nil
}
