package rooms

import (
	"context"
	"errors"
	"fmt"
	"github.com/avito-internships/test-backend-1-M1steryO/internal/domain"
	"github.com/georgysavva/scany/pgxscan"

	"github.com/jackc/pgx/v4"
)

func (r *RoomsRepository) Create(ctx context.Context, room domain.Room) (domain.Room, error) {
	query := `
		INSERT INTO rooms (id, name, description, capacity)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, description, capacity, created_at
	`

	var created domain.Room
	err := r.pool.QueryRow(ctx, query, room.ID, room.Name, room.Description, room.Capacity).
		Scan(&created.ID, &created.Name, &created.Description, &created.Capacity, &created.CreatedAt)
	if err != nil {
		return domain.Room{}, fmt.Errorf("create room: %w", err)
	}

	return created, nil
}

func (r *RoomsRepository) List(ctx context.Context) ([]domain.Room, error) {
	query := `
		SELECT id, name, description, capacity, created_at
		FROM rooms
		ORDER BY created_at, name
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list rooms: %w", err)
	}
	defer rows.Close()

	rooms := make([]domain.Room, 0)
	err = pgxscan.ScanAll(&rooms, rows)
	if err != nil {
		return nil, err
	}

	return rooms, rows.Err()
}

func (r *RoomsRepository) GetByID(ctx context.Context, roomID string) (domain.Room, error) {
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
		return domain.Room{}, fmt.Errorf("get room by id: %w", err)
	}

	return room, nil
}
