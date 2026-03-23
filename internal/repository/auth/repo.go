package auth

import (
	"github.com/jackc/pgx/v4/pgxpool"
)

type UsersRepository struct {
	pool *pgxpool.Pool
}

func NewUsersRepository(pool *pgxpool.Pool) *UsersRepository {
	return &UsersRepository{pool: pool}
}
