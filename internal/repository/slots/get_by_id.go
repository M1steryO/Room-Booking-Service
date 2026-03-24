package slots

import (
	"context"
	"errors"
	"fmt"

	"github.com/avito-internships/test-backend-1-M1steryO/internal/domain"
	"github.com/jackc/pgx/v4"
)

func (r *SlotsRepository) GetByID(ctx context.Context, slotID string) (domain.Slot, error) {
	const op = "repository.slots.GetByID"

	query := `
		SELECT id, room_id, start_at, end_at, created_at
		FROM slots
		WHERE id = $1
	`

	var slot domain.Slot
	err := r.pool.QueryRow(ctx, query, slotID).
		Scan(&slot.ID, &slot.RoomID, &slot.Start, &slot.End, &slot.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Slot{}, domain.SlotNotFound()
		}
		return domain.Slot{}, fmt.Errorf("%s: %w", op, err)
	}

	return slot, nil
}
