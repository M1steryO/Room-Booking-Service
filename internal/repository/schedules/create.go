package schedules

import (
	"context"
	"errors"
	"fmt"

	"github.com/M1steryO/Room-Booking-Service/internal/domain"
	"github.com/M1steryO/Room-Booking-Service/internal/repository"
	"github.com/jackc/pgconn"
)

func (r *SchedulesRepository) Create(ctx context.Context, schedule domain.Schedule) (domain.Schedule, error) {
	const op = "repository.schedules.Create"

	query := `
		INSERT INTO schedules (id, room_id, days_of_week, start_time, end_time, created_at)
		VALUES ($1, $2, $3, $4::time, $5::time, $6)
		RETURNING id, room_id, days_of_week, to_char(start_time, 'HH24:MI'), to_char(end_time, 'HH24:MI'), created_at
	`

	var created domain.Schedule
	err := r.pool.QueryRow(ctx, query, schedule.ID, schedule.RoomID, schedule.DaysOfWeek, schedule.StartTime, schedule.EndTime, schedule.CreatedAt).
		Scan(&created.ID, &created.RoomID, &created.DaysOfWeek, &created.StartTime, &created.EndTime, &created.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == repository.UniqueViolationCode {
			return domain.Schedule{}, domain.ScheduleExists()
		}
		return domain.Schedule{}, fmt.Errorf("%s: %w", op, err)
	}

	return created, nil
}
