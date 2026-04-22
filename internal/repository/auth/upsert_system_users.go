package auth

import (
	"context"
	"fmt"

	"github.com/M1steryO/Room-Booking-Service/internal/domain"
	"github.com/jackc/pgx/v4"
)

func (r *UsersRepository) UpsertSystemUsers(ctx context.Context, users []domain.User) error {
	const op = "repository.users.UpsertSystemUsers"

	query := `
		INSERT INTO users (id, email, password_hash, role, created_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO UPDATE
		SET email = EXCLUDED.email,
		    role = EXCLUDED.role
	`

	batch := &pgx.Batch{}
	for _, user := range users {
		batch.Queue(query, user.ID, user.Email, user.PasswordHash, user.Role, user.CreatedAt)
	}

	results := r.pool.SendBatch(ctx, batch)
	defer func(results pgx.BatchResults) {
		err := results.Close()
		if err != nil {
			return
		}
	}(results)

	for range users {
		if _, err := results.Exec(); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil
}
