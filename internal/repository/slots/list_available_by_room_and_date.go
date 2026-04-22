package slots

import (
	"context"
	"fmt"
	"time"

	"github.com/M1steryO/Room-Booking-Service/internal/domain"
	"github.com/georgysavva/scany/pgxscan"
)

func (r *SlotsRepository) ListAvailableByRoomAndDate(ctx context.Context, roomID string, date time.Time) ([]domain.Slot, error) {
	const op = "repository.slots.ListAvailableByRoomAndDate"

	start := date.UTC().Truncate(24 * time.Hour)
	end := start.Add(24 * time.Hour)

	query := `
		SELECT s.id, s.room_id, s.start_at as start, s.end_at as "end", s.created_at
		FROM slots s
		LEFT JOIN bookings b
		  ON b.slot_id = s.id
		 AND b.status = 'active'
		WHERE s.room_id = $1
		  AND s.start_at >= $2
		  AND s.start_at < $3
		  AND b.id IS NULL
		ORDER BY s.start_at
	`

	rows, err := r.pool.Query(ctx, query, roomID, start, end)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	slots := make([]domain.Slot, 0)
	err = pgxscan.ScanAll(&slots, rows)
	if err != nil {
		return nil, err
	}

	return slots, rows.Err()
}
