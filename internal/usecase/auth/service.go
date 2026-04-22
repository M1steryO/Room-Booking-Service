package auth

import (
	"github.com/M1steryO/Room-Booking-Service/internal/repository"
	"github.com/M1steryO/Room-Booking-Service/pkg/clock"
	"github.com/M1steryO/Room-Booking-Service/pkg/security"
)

const (
	dummyAdminID = "00000000-0000-0000-0000-000000000001"
	dummyUserID  = "00000000-0000-0000-0000-000000000002"
)

type AuthUsecase struct {
	usersRepo repository.UserRepository
	clock     clock.Clock
	tokens    *security.JWTManager
}

func NewAuthUsecase(usersRepo repository.UserRepository, clk clock.Clock, tokens *security.JWTManager) *AuthUsecase {
	return &AuthUsecase{
		usersRepo: usersRepo,
		clock:     clk,
		tokens:    tokens,
	}
}
