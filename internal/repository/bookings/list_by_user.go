package bookings

import (
	"context"
	"fmt"
	"time"

	"github.com/avito-internships/test-backend-1-M1steryO/internal/domain"
	"github.com/georgysavva/scany/pgxscan"
)

func (r *BookingsRepository) ListByUser(ctx context.Context, userID string, now time.Time) ([]domain.Booking, error) {
	const op = "repository.bookings.ListByUser"

	query := `
		SELECT b.id, b.slot_id, b.user_id, b.status, b.conference_link, b.created_at,
		       s.room_id, s.start_at AS slot_start, s.end_at AS slot_end
		FROM bookings b
		JOIN slots s ON s.id = b.slot_id
		WHERE b.user_id = $1
		  AND s.start_at >= $2
		ORDER BY s.start_at
	`

	rows, err := r.pool.Query(ctx, query, userID, now)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	bookings := make([]domain.Booking, 0)
	err = pgxscan.ScanAll(&bookings, rows)
	if err != nil {
		return nil, err
	}

	return bookings, rows.Err()
}
