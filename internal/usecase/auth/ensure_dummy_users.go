package auth

import (
	"context"

	"github.com/avito-internships/test-backend-1-M1steryO/internal/domain"
)

func (u *AuthUsecase) EnsureDummyUsers(ctx context.Context) error {
	now := u.clock.Now()
	return u.usersRepo.UpsertSystemUsers(ctx, []domain.User{
		{
			ID:        dummyAdminID,
			Email:     "admin@dummy.local",
			Role:      domain.RoleAdmin,
			CreatedAt: now,
		},
		{
			ID:        dummyUserID,
			Email:     "user@dummy.local",
			Role:      domain.RoleUser,
			CreatedAt: now,
		},
	})
}
