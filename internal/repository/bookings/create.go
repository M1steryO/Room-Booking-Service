package bookings

import (
	"context"
	"errors"
	"fmt"

	"github.com/M1steryO/Room-Booking-Service/internal/domain"
	"github.com/M1steryO/Room-Booking-Service/internal/repository"
	"github.com/jackc/pgconn"
)

func (r *BookingsRepository) Create(ctx context.Context, booking domain.Booking) (domain.Booking, error) {
	const op = "repository.bookings.Create"

	query := `INSERT INTO bookings (id, slot_id, user_id, status, conference_link, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, slot_id, user_id, status, conference_link, created_at`

	var created domain.Booking
	err := r.pool.QueryRow(ctx, query, booking.ID, booking.SlotID, booking.UserID, booking.Status, booking.ConferenceLink, booking.CreatedAt).
		Scan(&created.ID, &created.SlotID, &created.UserID, &created.Status, &created.ConferenceLink, &created.CreatedAt)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == repository.UniqueViolationCode {
			return domain.Booking{}, domain.SlotAlreadyBooked()
		}
		return domain.Booking{}, fmt.Errorf("%s: %w", op, err)
	}

	created.RoomID = booking.RoomID
	created.SlotStart = booking.SlotStart
	created.SlotEnd = booking.SlotEnd

	return created, nil
}
