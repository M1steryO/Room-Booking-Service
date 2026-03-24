package bookings

import (
	"context"
	"fmt"
)

func (r *BookingsRepository) CancelByOwner(ctx context.Context, bookingID string, userID string) error {
	const op = "repository.bookings.CancelByOwner"

	query := `
		UPDATE bookings
		SET status = 'cancelled'
		WHERE id = $1
		  AND status = 'active'
		  AND user_id = $2
	`

	_, err := r.pool.Exec(ctx, query, bookingID, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
