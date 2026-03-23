package slots

import (
	"context"
	"errors"
	"fmt"
	"github.com/avito-internships/test-backend-1-M1steryO/internal/domain"
	"github.com/georgysavva/scany/pgxscan"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

func (r *SlotsRepository) DateHasSlots(ctx context.Context, roomID string, date time.Time) (bool, error) {
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
		return false, fmt.Errorf("check slots existence: %w", err)
	}

	return exists, nil
}

func (r *SlotsRepository) BulkUpsert(ctx context.Context, slots []domain.Slot) error {
	if len(slots) == 0 {
		return nil
	}

	query := `
		INSERT INTO slots (id, room_id, start_at, end_at, created_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (room_id, start_at, end_at) DO UPDATE
		SET id = EXCLUDED.id
	`

	batch := &pgx.Batch{}
	for _, slot := range slots {
		batch.Queue(query, slot.ID, slot.RoomID, slot.Start, slot.End, slot.CreatedAt)
	}

	results := r.pool.SendBatch(ctx, batch)
	defer func() {
		err := results.Close()
		if err != nil {
			return
		}
	}()

	for range slots {
		if _, err := results.Exec(); err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				return fmt.Errorf("bulk upsert slots (%s): %w", pgErr.Code, err)
			}
			return fmt.Errorf("bulk upsert slots: %w", err)
		}
	}

	return nil
}

func (r *SlotsRepository) GetByID(ctx context.Context, slotID string) (domain.Slot, error) {
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
		return domain.Slot{}, fmt.Errorf("get slot by id: %w", err)
	}

	return slot, nil
}

func (r *SlotsRepository) ListAvailableByRoomAndDate(ctx context.Context, roomID string, date time.Time) ([]domain.Slot, error) {
	start := date.UTC().Truncate(24 * time.Hour)
	end := start.Add(24 * time.Hour)

	query := `
		SELECT s.id, s.room_id, s.start_at, s.end_at, s.created_at
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
		return nil, fmt.Errorf("list available slots: %w", err)
	}
	defer rows.Close()

	slots := make([]domain.Slot, 0)
	err = pgxscan.ScanAll(&slots, rows)
	if err != nil {
		return nil, err
	}

	return slots, rows.Err()
}
