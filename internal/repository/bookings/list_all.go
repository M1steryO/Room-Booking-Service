package bookings

import (
	"context"
	"fmt"

	"github.com/avito-internships/test-backend-1-M1steryO/internal/domain"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
)

func (r *BookingsRepository) ListAll(ctx context.Context, page int, pageSize int) ([]domain.Booking, int, error) {
	const op = "repository.bookings.ListAll"

	offset := (page - 1) * pageSize

	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.RepeatableRead,
		AccessMode: pgx.ReadOnly,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	var total int

	if err := tx.QueryRow(ctx, `SELECT count(*) FROM bookings`).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}

	query := `
		SELECT b.id, b.slot_id, b.user_id, b.status, b.conference_link, b.created_at,
		       s.room_id, s.start_at AS slot_start, s.end_at AS slot_end
		FROM bookings b
		JOIN slots s ON s.id = b.slot_id
		ORDER BY b.created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	bookings := make([]domain.Booking, 0, pageSize)
	err = pgxscan.ScanAll(&bookings, rows)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}

	return bookings, total, nil
}
