package auth

import (
	"context"
	"github.com/avito-internships/test-backend-1-M1steryO/internal/domain"
	"github.com/avito-internships/test-backend-1-M1steryO/internal/platform/identity"
	"github.com/avito-internships/test-backend-1-M1steryO/internal/platform/security"
)

const (
	DummyAdminID = "00000000-0000-0000-0000-000000000001"
	DummyUserID  = "00000000-0000-0000-0000-000000000002"
)

func (u *AuthUsecase) EnsureDummyUsers(ctx context.Context) error {
	now := u.clock.Now()
	return u.usersRepo.UpsertSystemUsers(ctx, []domain.User{
		{
			ID:        DummyAdminID,
			Email:     "admin@dummy.local",
			Role:      domain.RoleAdmin,
			CreatedAt: now,
		},
		{
			ID:        DummyUserID,
			Email:     "user@dummy.local",
			Role:      domain.RoleUser,
			CreatedAt: now,
		},
	})
}

func (u *AuthUsecase) DummyLogin(role domain.Role) (string, error) {
	if !role.Valid() {
		return "", domain.InvalidRequest("role must be admin or user")
	}

	userID := DummyUserID
	if role == domain.RoleAdmin {
		userID = DummyAdminID
	}

	return u.tokens.Issue(userID, role, u.clock.Now())
}

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
