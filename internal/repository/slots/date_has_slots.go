package slots

import (
	"context"
	"fmt"
	"time"
)

func (r *SlotsRepository) DateHasSlots(ctx context.Context, roomID string, date time.Time) (bool, error) {
	const op = "repository.slots.DateHasSlots"

	start := date.UTC().Truncate(24 * time.Hour)
	end := start.Add(24 * time.Hour)

	query := `
		SELECT EXISTS(
			SELECT 1
			FROM slots
			WHERE room_id = $1
			  AND start_at >= $2
			  AND start_at < $3
		)
	`

	var exists bool
	if err := r.pool.QueryRow(ctx, query, roomID, start, end).Scan(&exists); err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return exists, nil
}
