package bookings

import (
	"context"
	"errors"
	"fmt"
	"github.com/avito-internships/test-backend-1-M1steryO/internal/domain"
	"github.com/jackc/pgx/v4"
	"time"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgconn"
)

func (r *BookingsRepository) Create(ctx context.Context, booking domain.Booking) (domain.Booking, error) {
	query := `INSERT INTO bookings (id, slot_id, user_id, status, conference_link, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, slot_id, user_id, status, conference_link, created_at`

	var created domain.Booking
	err := r.pool.QueryRow(ctx, query, booking.ID, booking.SlotID, booking.UserID, booking.Status, booking.ConferenceLink, booking.CreatedAt).
		Scan(&created.ID, &created.SlotID, &created.UserID, &created.Status, &created.ConferenceLink, &created.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.Booking{}, domain.SlotAlreadyBooked()
		}
		return domain.Booking{}, fmt.Errorf("create booking: %w", err)
	}

	created.RoomID = booking.RoomID
	created.SlotStart = booking.SlotStart
	created.SlotEnd = booking.SlotEnd

	return created, nil
}

func (r *BookingsRepository) ListAll(ctx context.Context, page int, pageSize int) ([]domain.Booking, int, error) {
	offset := (page - 1) * pageSize

	var total int

	if err := r.pool.QueryRow(ctx, `SELECT count(*) FROM bookings`).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count bookings: %w", err)
	}

	query := `
		SELECT b.id, b.slot_id, b.user_id, b.status, b.conference_link, b.created_at,
		       s.room_id, s.start_at, s.end_at
		FROM bookings b
		JOIN slots s ON s.id = b.slot_id
		ORDER BY b.created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.pool.Query(ctx, query, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list bookings: %w", err)
	}
	defer rows.Close()

	bookings := make([]domain.Booking, 0, pageSize)
	err = pgxscan.ScanAll(&bookings, rows)
	if err != nil {
		return nil, 0, err
	}

	return bookings, total, rows.Err()
}

func (r *BookingsRepository) ListFutureByUser(ctx context.Context, userID string, now time.Time) ([]domain.Booking, error) {
	query := `
		SELECT b.id, b.slot_id, b.user_id, b.status, b.conference_link, b.created_at,
		       s.room_id, s.start_at, s.end_at
		FROM bookings b
		JOIN slots s ON s.id = b.slot_id
		WHERE b.user_id = $1
		  AND s.start_at >= $2
		ORDER BY s.start_at
	`

	rows, err := r.pool.Query(ctx, query, userID, now)
	if err != nil {
		return nil, fmt.Errorf("list user future bookings: %w", err)
	}
	defer rows.Close()

	bookings := make([]domain.Booking, 0)
	err = pgxscan.ScanAll(&bookings, rows)
	if err != nil {
		return nil, err
	}

	return bookings, rows.Err()
}

func (r *BookingsRepository) GetByID(ctx context.Context, bookingID string) (domain.Booking, error) {
	query := `
		SELECT b.id, b.slot_id, b.user_id, b.status, b.conference_link, b.created_at,
		       s.room_id, s.start_at, s.end_at
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
		return domain.Booking{}, fmt.Errorf("get booking: %w", err)
	}

	return booking, nil
}

func (r *BookingsRepository) CancelByOwner(ctx context.Context, bookingID string, userID string) (domain.Booking, error) {
	booking, err := r.GetByID(ctx, bookingID)
	if err != nil {
		return domain.Booking{}, err
	}

	if booking.UserID != userID {
		return domain.Booking{}, domain.Forbidden("cannot cancel another user's booking")
	}

	query := `
		UPDATE bookings
		SET status = 'cancelled'
		WHERE id = $1
		  AND status = 'active'
	`

	if _, err := r.pool.Exec(ctx, query, bookingID); err != nil {
		return domain.Booking{}, fmt.Errorf("cancel booking: %w", err)
	}

	return r.GetByID(ctx, bookingID)
}
