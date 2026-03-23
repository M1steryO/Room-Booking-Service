package auth

import (
	"github.com/avito-internships/test-backend-1-M1steryO/internal/platform/clock"
	"github.com/avito-internships/test-backend-1-M1steryO/internal/platform/security"
	"github.com/avito-internships/test-backend-1-M1steryO/internal/repository"
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
