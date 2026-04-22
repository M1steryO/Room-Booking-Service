package auth_test

import (
	"context"
	"github.com/M1steryO/Room-Booking-Service/pkg/security"
	"testing"
	"time"

	"github.com/M1steryO/Room-Booking-Service/internal/domain"
	repmocks "github.com/M1steryO/Room-Booking-Service/internal/repository/mocks"
	authuc "github.com/M1steryO/Room-Booking-Service/internal/usecase/auth"
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

	expectedByID := map[string]struct {
		email string
		role  domain.Role
	}{
		"00000000-0000-0000-0000-000000000001": {email: "admin@dummy.local", role: domain.RoleAdmin},
		"00000000-0000-0000-0000-000000000002": {email: "user@dummy.local", role: domain.RoleUser},
	}

	for _, user := range got {
		expected, ok := expectedByID[user.ID]
		if !ok {
			t.Fatalf("unexpected dummy user id: %s", user.ID)
		}
		if user.Email != expected.email {
			t.Fatalf("unexpected email for id=%s: got=%s want=%s", user.ID, user.Email, expected.email)
		}
		if user.Role != expected.role {
			t.Fatalf("unexpected role for id=%s: got=%s want=%s", user.ID, user.Role, expected.role)
		}
		if user.CreatedAt.IsZero() {
			t.Fatalf("expected non-zero created_at for id=%s", user.ID)
		}
		delete(expectedByID, user.ID)
	}

	if len(expectedByID) != 0 {
		t.Fatalf("missing expected dummy users: %+v", expectedByID)
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
	repo.GetByEmailMock.Set(func(_ context.Context, email string) (domain.User, error) {
		if email != "u@test.dev" {
			t.Fatalf("unexpected email: %s", email)
		}
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

func TestLoginTokenContainsExpectedClaims(t *testing.T) {
	now := time.Now().UTC()
	jwtManager := security.NewJWTManager("secret", time.Hour)
	hash, err := security.HashPassword("qwerty123")
	if err != nil {
		t.Fatalf("HashPassword error: %v", err)
	}

	repo := repmocks.NewUserRepositoryMock(t)
	repo.CreateMock.Optional()
	repo.UpsertSystemUsersMock.Optional()
	repo.GetByEmailMock.Set(func(_ context.Context, _ string) (domain.User, error) {
		return domain.User{
			ID:           "11111111-2222-3333-4444-555555555555",
			Email:        "u@test.dev",
			PasswordHash: hash,
			Role:         domain.RoleUser,
			CreatedAt:    now,
		}, nil
	})

	uc := authuc.NewAuthUsecase(repo, fixedClock{now: now}, jwtManager)

	token, err := uc.Login(context.Background(), "u@test.dev", "qwerty123")
	if err != nil {
		t.Fatalf("Login error: %v", err)
	}

	claims, err := jwtManager.Parse(token)
	if err != nil {
		t.Fatalf("Parse token error: %v", err)
	}

	if claims.UserID != "11111111-2222-3333-4444-555555555555" {
		t.Fatalf("unexpected user_id claim: got=%s", claims.UserID)
	}
	if claims.Role != domain.RoleUser {
		t.Fatalf("unexpected role claim: got=%s want=%s", claims.Role, domain.RoleUser)
	}
}

func TestRegisterInvalidRole(t *testing.T) {
	repo := repmocks.NewUserRepositoryMock(t)
	repo.CreateMock.Optional()
	repo.GetByEmailMock.Optional()
	repo.UpsertSystemUsersMock.Optional()
	uc := authuc.NewAuthUsecase(repo, fixedClock{now: time.Now().UTC()}, security.NewJWTManager("secret", time.Hour))

	_, err := uc.Register(context.Background(), "u@test.dev", "qwerty123", "invalid-role")
	if err == nil {
		t.Fatal("expected error for invalid role")
	}
	if domain.AsAppError(err).Code != "INVALID_REQUEST" {
		t.Fatalf("expected INVALID_REQUEST, got %s", domain.AsAppError(err).Code)
	}
}

func TestRegisterEmailAlreadyExists(t *testing.T) {
	repo := repmocks.NewUserRepositoryMock(t)
	repo.CreateMock.Set(func(_ context.Context, _ domain.User) (domain.User, error) {
		return domain.User{}, domain.InvalidRequest("email already exists")
	})
	repo.GetByEmailMock.Optional()
	repo.UpsertSystemUsersMock.Optional()
	uc := authuc.NewAuthUsecase(repo, fixedClock{now: time.Now().UTC()}, security.NewJWTManager("secret", time.Hour))

	_, err := uc.Register(context.Background(), "u@test.dev", "qwerty123", domain.RoleUser)
	if err == nil {
		t.Fatal("expected error")
	}
	if domain.AsAppError(err).Code != "INVALID_REQUEST" {
		t.Fatalf("expected INVALID_REQUEST, got %s", domain.AsAppError(err).Code)
	}
}

func TestLoginWrongPassword(t *testing.T) {
	hash, err := security.HashPassword("qwerty123")
	if err != nil {
		t.Fatalf("HashPassword error: %v", err)
	}

	repo := repmocks.NewUserRepositoryMock(t)
	repo.CreateMock.Optional()
	repo.UpsertSystemUsersMock.Optional()
	repo.GetByEmailMock.Set(func(_ context.Context, email string) (domain.User, error) {
		if email != "u@test.dev" {
			t.Fatalf("unexpected email: %s", email)
		}
		return domain.User{
			ID:           "11111111-2222-3333-4444-555555555555",
			Email:        "u@test.dev",
			PasswordHash: hash,
			Role:         domain.RoleUser,
			CreatedAt:    time.Now().UTC(),
		}, nil
	})
	uc := authuc.NewAuthUsecase(repo, fixedClock{now: time.Now().UTC()}, security.NewJWTManager("secret", time.Hour))

	_, err = uc.Login(context.Background(), "u@test.dev", "wrong-password")
	if err == nil {
		t.Fatal("expected error")
	}
	if domain.AsAppError(err).Code != "UNAUTHORIZED" {
		t.Fatalf("expected UNAUTHORIZED, got %s", domain.AsAppError(err).Code)
	}
}

func TestLoginUserNotFound(t *testing.T) {
	repo := repmocks.NewUserRepositoryMock(t)
	repo.CreateMock.Optional()
	repo.UpsertSystemUsersMock.Optional()
	repo.GetByEmailMock.Set(func(_ context.Context, email string) (domain.User, error) {
		if email != "u@test.dev" {
			t.Fatalf("unexpected email: %s", email)
		}
		return domain.User{}, domain.Unauthorized("invalid credentials")
	})
	uc := authuc.NewAuthUsecase(repo, fixedClock{now: time.Now().UTC()}, security.NewJWTManager("secret", time.Hour))

	_, err := uc.Login(context.Background(), "u@test.dev", "qwerty123")
	if err == nil {
		t.Fatal("expected error")
	}
	if domain.AsAppError(err).Code != "UNAUTHORIZED" {
		t.Fatalf("expected UNAUTHORIZED, got %s", domain.AsAppError(err).Code)
	}
}
