package slots

import (
	"context"
	"fmt"
	"strings"

	"github.com/avito-internships/test-backend-1-M1steryO/internal/domain"
)

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
