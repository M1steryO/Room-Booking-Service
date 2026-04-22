package auth

import (
	"context"

	"github.com/M1steryO/Room-Booking-Service/internal/domain"
	"github.com/M1steryO/Room-Booking-Service/pkg/security"
)

func (u *AuthUsecase) Login(ctx context.Context, email string, password string) (string, error) {
	user, err := u.usersRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	if !security.CheckPassword(password, user.PasswordHash) {
		return "", domain.Unauthorized("invalid credentials")
	}

	return u.tokens.Issue(user.ID, user.Role, u.clock.Now())
}
