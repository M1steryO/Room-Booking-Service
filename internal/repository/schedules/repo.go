package schedules

import (
	"github.com/jackc/pgx/v4/pgxpool"
)

type SchedulesRepository struct {
	pool *pgxpool.Pool
}

func NewSchedulesRepository(pool *pgxpool.Pool) *SchedulesRepository {
	return &SchedulesRepository{pool: pool}
}
