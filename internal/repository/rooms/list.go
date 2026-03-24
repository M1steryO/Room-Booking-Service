package rooms

import (
	"context"
	"fmt"

	"github.com/avito-internships/test-backend-1-M1steryO/internal/domain"
	"github.com/georgysavva/scany/pgxscan"
)

func (r *RoomsRepository) List(ctx context.Context) ([]domain.Room, error) {
	const op = "repository.rooms.List"

	query := `
		SELECT id, name, description, capacity, created_at
		FROM rooms
		ORDER BY created_at, name
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	rooms := make([]domain.Room, 0)
	err = pgxscan.ScanAll(&rooms, rows)
	if err != nil {
		return nil, err
	}

	return rooms, rows.Err()
}
