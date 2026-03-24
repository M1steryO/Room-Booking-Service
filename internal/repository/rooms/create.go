package rooms

import (
	"context"
	"fmt"

	"github.com/avito-internships/test-backend-1-M1steryO/internal/domain"
)

func (r *RoomsRepository) Create(ctx context.Context, room domain.Room) (domain.Room, error) {
	const op = "repository.rooms.Create"

	query := `
		INSERT INTO rooms (id, name, description, capacity)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, description, capacity, created_at
	`

	var created domain.Room
	err := r.pool.QueryRow(ctx, query, room.ID, room.Name, room.Description, room.Capacity).
		Scan(&created.ID, &created.Name, &created.Description, &created.Capacity, &created.CreatedAt)
	if err != nil {
		return domain.Room{}, fmt.Errorf("%s: %w", op, err)
	}

	return created, nil
}
