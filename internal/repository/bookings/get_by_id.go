package bookings

import (
	"context"
	"errors"
	"fmt"

	"github.com/M1steryO/Room-Booking-Service/internal/domain"
	"github.com/jackc/pgx/v4"
)

func (r *BookingsRepository) GetByID(ctx context.Context, bookingID string) (domain.Booking, error) {
	const op = "repository.bookings.GetByID"

	query := `
		SELECT b.id, b.slot_id, b.user_id, b.status, b.conference_link, b.created_at,
		       s.room_id, s.start_at AS slot_start, s.end_at AS slot_end
		FROM bookings b
		JOIN slots s ON s.id = b.slot_id
		WHERE b.id = $1
	`

	var booking domain.Booking
	err := r.pool.QueryRow(ctx, query, bookingID).Scan(
		&booking.ID,
		&booking.SlotID,
		&booking.UserID,
		&booking.Status,
		&booking.ConferenceLink,
		&booking.CreatedAt,
		&booking.RoomID,
		&booking.SlotStart,
		&booking.SlotEnd,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Booking{}, domain.BookingNotFound()
		}
		return domain.Booking{}, fmt.Errorf("%s: %w", op, err)
	}

	return booking, nil
}
