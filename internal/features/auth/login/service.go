package login

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	sharedErrors "github.com/nanasuryana335/honda-leasing-api/internal/shared/errors"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo      *Repository
	jwtSecret string
}

func NewService(repo *Repository, jwtSecret string) *Service {
	return &Service{
		repo:      repo,
		jwtSecret: jwtSecret,
	}
}

func (s *Service) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {

	user, err := s.repo.FindUserByPhone(ctx, req.PhoneNumber)
	if err != nil {
		return nil, sharedErrors.ErrInvalidCredentials
	}

	if user.LockedUntil.After(time.Now()) {
		return nil, fmt.Errorf("account is locked until %s", user.LockedUntil.Format(time.RFC3339))
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PinKey), []byte(req.PIN)); err != nil {
		_ = s.repo.IncrementFailedAttempts(ctx, user.UserID)

		if user.FailedAttempts >= 4 { 
			// TODO REK
		}

		return nil, sharedErrors.ErrInvalidCredentials
	}

	_ = s.repo.ResetFailedAttempts(ctx, user.UserID)

	roles, err := s.repo.GetUserRoles(ctx, user.UserID)
	if err != nil {
		return nil, err
	}

	accessToken, err := s.generateAccessToken(user.UserID, roles)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateRefreshToken(user.UserID)
	if err != nil {
		return nil, err
	}

	_ = s.repo.UpdateLastLogin(ctx, user.UserID)

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: UserInfo{
			UserID:      user.UserID,
			FullName:    user.FullName,
			PhoneNumber: user.PhoneNumber,
			Email:       user.Email,
			Roles:       roles,
		},
	}, nil
}

func (s *Service) generateAccessToken(userID int64, roles []string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"roles":   roles,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *Service) generateRefreshToken(userID int64) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}
