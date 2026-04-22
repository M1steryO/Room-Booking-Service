package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/M1steryO/Room-Booking-Service/internal/domain"
	"github.com/M1steryO/Room-Booking-Service/internal/repository"
	"github.com/jackc/pgconn"
)

func (r *UsersRepository) Create(ctx context.Context, user domain.User) (domain.User, error) {
	const op = "repository.users.Create"

	query := `
		INSERT INTO users (id, email, password_hash, role, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, email, password_hash, role, created_at
	`

	var created domain.User
	err := r.pool.QueryRow(ctx, query, user.ID, user.Email, user.PasswordHash, user.Role, user.CreatedAt).
		Scan(&created.ID, &created.Email, &created.PasswordHash, &created.Role, &created.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == repository.UniqueViolationCode {
			return domain.User{}, domain.InvalidRequest("email already exists")
		}
		return domain.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return created, nil
}
