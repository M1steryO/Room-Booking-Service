package auth_test

import (
	"context"
	"github.com/avito-internships/test-backend-1-M1steryO/pkg/security"
	"testing"
	"time"

	"github.com/avito-internships/test-backend-1-M1steryO/internal/domain"
	repmocks "github.com/avito-internships/test-backend-1-M1steryO/internal/repository/mocks"
	authuc "github.com/avito-internships/test-backend-1-M1steryO/internal/usecase/auth"
)

type fixedClock struct {
	now time.Time
}

func (f fixedClock) Now() time.Time { return f.now }

func TestEnsureDummyUsers(t *testing.T) {
	var got []domain.User
	repo := repmocks.NewUserRepositoryMock(t)
	repo.UpsertSystemUsersMock.Set(func(_ context.Context, users []domain.User) error {
		got = users
		return nil
	})
	repo.CreateMock.Optional()
	repo.GetByEmailMock.Optional()
	uc := authuc.NewAuthUsecase(repo, fixedClock{now: time.Now().UTC()}, security.NewJWTManager("secret", time.Hour))

	if err := uc.EnsureDummyUsers(context.Background()); err != nil {
		t.Fatalf("EnsureDummyUsers error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 dummy users, got %d", len(got))
	}
}

func TestDummyLoginInvalidRole(t *testing.T) {
	repo := repmocks.NewUserRepositoryMock(t)
	repo.CreateMock.Optional()
	repo.GetByEmailMock.Optional()
	repo.UpsertSystemUsersMock.Optional()
	uc := authuc.NewAuthUsecase(repo, fixedClock{now: time.Now().UTC()}, security.NewJWTManager("secret", time.Hour))

	_, err := uc.DummyLogin("invalid-role")
	if err == nil {
		t.Fatal("expected error for invalid role")
	}
	if domain.AsAppError(err).Code != "INVALID_REQUEST" {
		t.Fatalf("expected INVALID_REQUEST, got %s", domain.AsAppError(err).Code)
	}
}

func TestDummyLoginUsesFixedUserIDsByRole(t *testing.T) {
	repo := repmocks.NewUserRepositoryMock(t)
	repo.CreateMock.Optional()
	repo.GetByEmailMock.Optional()
	repo.UpsertSystemUsersMock.Optional()
	jwtManager := security.NewJWTManager("secret", time.Hour)
	uc := authuc.NewAuthUsecase(repo, fixedClock{now: time.Now().UTC()}, jwtManager)

	cases := []struct {
		name           string
		role           domain.Role
		expectedUserID string
	}{
		{
			name:           "user role maps to user dummy id",
			role:           domain.RoleUser,
			expectedUserID: "00000000-0000-0000-0000-000000000002",
		},
		{
			name:           "admin role maps to admin dummy id",
			role:           domain.RoleAdmin,
			expectedUserID: "00000000-0000-0000-0000-000000000001",
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			token, err := uc.DummyLogin(testCase.role)
			if err != nil {
				t.Fatalf("DummyLogin error: %v", err)
			}

			claims, err := jwtManager.Parse(token)
			if err != nil {
				t.Fatalf("Parse token error: %v", err)
			}

			if claims.UserID != testCase.expectedUserID {
				t.Fatalf("unexpected user_id in token: got=%s want=%s", claims.UserID, testCase.expectedUserID)
			}

			if claims.Role != testCase.role {
				t.Fatalf("unexpected role in token: got=%s want=%s", claims.Role, testCase.role)
			}
		})
	}
}

func TestRegisterAndLogin(t *testing.T) {
	var created domain.User
	password := "qwerty123"
	repo := repmocks.NewUserRepositoryMock(t)
	repo.CreateMock.Set(func(_ context.Context, user domain.User) (domain.User, error) {
		created = user
		return user, nil
	})
	repo.GetByEmailMock.Set(func(_ context.Context, _ string) (domain.User, error) {
		return created, nil
	})
	repo.UpsertSystemUsersMock.Optional()
	uc := authuc.NewAuthUsecase(repo, fixedClock{now: time.Now().UTC()}, security.NewJWTManager("secret", time.Hour))

	_, err := uc.Register(context.Background(), "u@test.dev", password, domain.RoleUser)
	if err != nil {
		t.Fatalf("Register error: %v", err)
	}
	if created.PasswordHash == "" || created.PasswordHash == password {
		t.Fatal("expected hashed password to be stored")
	}

	token, err := uc.Login(context.Background(), "u@test.dev", password)
	if err != nil {
		t.Fatalf("Login error: %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token")
	}
}
