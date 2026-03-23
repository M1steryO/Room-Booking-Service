package slots

import (
	"github.com/jackc/pgx/v4/pgxpool"
)

type SlotsRepository struct {
	pool *pgxpool.Pool
}

func NewSlotsRepository(pool *pgxpool.Pool) *SlotsRepository {
	return &SlotsRepository{pool: pool}
}
