package auth

import (
	"context"

	"github.com/M1steryO/Room-Booking-Service/internal/domain"
	"github.com/M1steryO/Room-Booking-Service/pkg/identity"
	"github.com/M1steryO/Room-Booking-Service/pkg/security"
)

func (u *AuthUsecase) Register(ctx context.Context, email string, password string, role domain.Role) (domain.User, error) {
	if email == "" || password == "" || !role.Valid() {
		return domain.User{}, domain.InvalidRequest("invalid registration payload")
	}

	hash, err := security.HashPassword(password)
	if err != nil {
		return domain.User{}, err
	}

	user := domain.User{
		ID:           identity.New(),
		Email:        email,
		PasswordHash: hash,
		Role:         role,
		CreatedAt:    u.clock.Now(),
	}

	return u.usersRepo.Create(ctx, user)
}
