package rooms

import (
	"context"
	"errors"
	"fmt"

	"github.com/M1steryO/Room-Booking-Service/internal/domain"
	"github.com/jackc/pgx/v4"
)

func (r *RoomsRepository) GetByID(ctx context.Context, roomID string) (domain.Room, error) {
	const op = "repository.rooms.GetByID"

	query := `
		SELECT id, name, description, capacity, created_at
		FROM rooms
		WHERE id = $1
	`

	var room domain.Room
	err := r.pool.QueryRow(ctx, query, roomID).
		Scan(&room.ID, &room.Name, &room.Description, &room.Capacity, &room.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Room{}, domain.RoomNotFound()
		}
		return domain.Room{}, fmt.Errorf("%s: %w", op, err)
	}

	return room, nil
}
