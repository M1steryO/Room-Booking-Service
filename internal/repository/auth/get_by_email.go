package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/avito-internships/test-backend-1-M1steryO/internal/domain"
	"github.com/jackc/pgx/v4"
)

func (r *UsersRepository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	const op = "repository.users.GetByEmail"

	query := `
		SELECT id, email, password_hash, role, created_at
		FROM users
		WHERE email = $1
	`

	var user domain.User
	err := r.pool.QueryRow(ctx, query, email).
		Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, domain.Unauthorized("invalid credentials")
		}
		return domain.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}
