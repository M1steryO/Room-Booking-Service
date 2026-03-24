package slots

import (
	"context"
	"errors"
	"fmt"
	"github.com/avito-internships/test-backend-1-M1steryO/internal/domain"
	"github.com/georgysavva/scany/pgxscan"
	"strings"
	"time"

	"github.com/jackc/pgx/v4"
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

func (r *SlotsRepository) BulkUpsert(ctx context.Context, slots []domain.Slot) error {
	const op = "repository.slots.BulkUpsert"

	if len(slots) == 0 {
		return nil
	}

	const colsPerRow = 5

	var (
		sb   strings.Builder
		args = make([]any, 0, len(slots)*colsPerRow)
	)

	sb.WriteString(`
		INSERT INTO slots (id, room_id, start_at, end_at, created_at)
		VALUES
	`)

	for i, slot := range slots {
		if i > 0 {
			sb.WriteString(",")
		}

		base := i * colsPerRow
		
		_, _ = fmt.Fprintf(&sb, "($%d, $%d, $%d, $%d, $%d)",
			base+1, base+2, base+3, base+4, base+5,
		)

		args = append(args,
			slot.ID,
			slot.RoomID,
			slot.Start,
			slot.End,
			slot.CreatedAt,
		)
	}

	sb.WriteString(`
		ON CONFLICT (room_id, start_at, end_at) DO NOTHING
	`)

	_, err := r.pool.Exec(ctx, sb.String(), args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

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
