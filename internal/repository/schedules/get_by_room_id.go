package schedules

import (
	"context"
	"errors"
	"fmt"

	"github.com/avito-internships/test-backend-1-M1steryO/internal/domain"
	"github.com/jackc/pgx/v4"
)

func (r *SchedulesRepository) GetByRoomID(ctx context.Context, roomID string) (domain.Schedule, error) {
	const op = "repository.schedules.GetByRoomID"

	query := `
		SELECT id, room_id, days_of_week, to_char(start_time, 'HH24:MI'), to_char(end_time, 'HH24:MI'), created_at
		FROM schedules
		WHERE room_id = $1
	`

	var schedule domain.Schedule
	err := r.pool.QueryRow(ctx, query, roomID).
		Scan(&schedule.ID, &schedule.RoomID, &schedule.DaysOfWeek, &schedule.StartTime, &schedule.EndTime, &schedule.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Schedule{}, domain.NotFound("schedule not found")
		}
		return domain.Schedule{}, fmt.Errorf("%s: %w", op, err)
	}

	return schedule, nil
}
