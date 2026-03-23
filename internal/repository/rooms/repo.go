package rooms

import (
	"github.com/jackc/pgx/v4/pgxpool"
)

type RoomsRepository struct {
	pool *pgxpool.Pool
}

func NewRoomsRepository(pool *pgxpool.Pool) *RoomsRepository {
	return &RoomsRepository{pool: pool}
}
