package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/avito-internships/test-backend-1-M1steryO/internal/domain"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgconn"
)

func (r *UsersRepository) Create(ctx context.Context, user domain.User) (domain.User, error) {
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
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.User{}, domain.InvalidRequest("email already exists")
		}
		return domain.User{}, fmt.Errorf("create user: %w", err)
	}

	return created, nil
}

func (r *UsersRepository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
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
		return domain.User{}, fmt.Errorf("get user by email: %w", err)
	}

	return user, nil
}

func (r *UsersRepository) UpsertSystemUsers(ctx context.Context, users []domain.User) error {
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
			return fmt.Errorf("upsert system users: %w", err)
		}
	}

	return nil
}
