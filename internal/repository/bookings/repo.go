package bookings

import (
	"github.com/jackc/pgx/v4/pgxpool"
)

type BookingsRepository struct {
	pool *pgxpool.Pool
}

func NewBookingsRepository(pool *pgxpool.Pool) *BookingsRepository {
	return &BookingsRepository{pool: pool}
}
