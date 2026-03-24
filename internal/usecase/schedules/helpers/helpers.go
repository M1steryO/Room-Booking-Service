package helpers

import (
	"fmt"
	"github.com/avito-internships/test-backend-1-M1steryO/internal/domain"
	"github.com/avito-internships/test-backend-1-M1steryO/pkg/identity"
	"slices"
	"time"
)

func ValidateSchedule(days []int, startTime string, endTime string) error {
	if len(days) == 0 {
		return fmt.Errorf("daysOfWeek must not be empty")
	}

	for _, day := range days {
		if day < 1 || day > 7 {
			return fmt.Errorf("daysOfWeek must be in range 1..7")
		}
	}

	start, err := time.Parse("15:04", startTime)
	if err != nil {
		return fmt.Errorf("invalid startTime: %w", err)
	}

	end, err := time.Parse("15:04", endTime)
	if err != nil {
		return fmt.Errorf("invalid endTime: %w", err)
	}

	if !end.After(start) {
		return fmt.Errorf("endTime must be after startTime")
	}

	if end.Sub(start)%(30*time.Minute) != 0 {
		return fmt.Errorf("time range must be divisible by 30 minutes")
	}

	return nil
}

func GenerateSlotsForDate(roomID string, schedule domain.Schedule, date time.Time) ([]domain.Slot, error) {
	date = date.UTC().Truncate(24 * time.Hour)
	weekday := int(date.Weekday())
	if weekday == 0 {
		weekday = 7
	}

	if !slices.Contains(schedule.DaysOfWeek, weekday) {
		return nil, nil
	}

	startClock, err := time.Parse("15:04", schedule.StartTime)
	if err != nil {
		return nil, err
	}

	endClock, err := time.Parse("15:04", schedule.EndTime)
	if err != nil {
		return nil, err
	}

	start := time.Date(date.Year(), date.Month(), date.Day(), startClock.Hour(), startClock.Minute(), 0, 0, time.UTC)
	end := time.Date(date.Year(), date.Month(), date.Day(), endClock.Hour(), endClock.Minute(), 0, 0, time.UTC)

	slots := make([]domain.Slot, 0, int(end.Sub(start)/(30*time.Minute)))
	for cursor := start; cursor.Before(end); cursor = cursor.Add(30 * time.Minute) {
		next := cursor.Add(30 * time.Minute)
		if next.After(end) {
			break
		}

		slots = append(slots, domain.Slot{
			ID:        identity.DeterministicSlotID(roomID, cursor.Format(time.RFC3339), next.Format(time.RFC3339)),
			RoomID:    roomID,
			Start:     cursor,
			End:       next,
			CreatedAt: time.Now().UTC(),
		})
	}

	return slots, nil
}

func DatesWindow(from time.Time, days int) []time.Time {
	out := make([]time.Time, 0, days)
	base := from.UTC().Truncate(24 * time.Hour)
	for i := 0; i < days; i++ {
		out = append(out, base.AddDate(0, 0, i))
	}

	return out
}
