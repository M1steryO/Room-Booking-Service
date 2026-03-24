package schedules

import (
	"context"

	"github.com/avito-internships/test-backend-1-M1steryO/internal/domain"
	"github.com/avito-internships/test-backend-1-M1steryO/internal/usecase/schedules/helpers"
	"github.com/avito-internships/test-backend-1-M1steryO/pkg/identity"
)

func (u *SchedulesUsecase) Create(ctx context.Context, actorRole domain.Role, roomID string, days []int, startTime string, endTime string) (domain.Schedule, error) {
	if actorRole != domain.RoleAdmin {
		return domain.Schedule{}, domain.Forbidden("admin role required")
	}

	if _, err := u.roomsRepo.GetByID(ctx, roomID); err != nil {
		return domain.Schedule{}, err
	}

	if err := helpers.ValidateSchedule(days, startTime, endTime); err != nil {
		return domain.Schedule{}, domain.InvalidRequest(err.Error())
	}

	schedule := domain.Schedule{
		ID:         identity.New(),
		RoomID:     roomID,
		DaysOfWeek: days,
		StartTime:  startTime,
		EndTime:    endTime,
		CreatedAt:  u.clock.Now(),
	}

	created, err := u.schedulesRepo.Create(ctx, schedule)
	if err != nil {
		return domain.Schedule{}, err
	}

	window := helpers.DatesWindow(u.clock.Now(), u.slotHorizonDay)
	allSlots := make([]domain.Slot, 0)
	for _, date := range window {
		slots, genErr := helpers.GenerateSlotsForDate(roomID, created, date)
		if genErr != nil {
			return domain.Schedule{}, genErr
		}

		allSlots = append(allSlots, slots...)
	}

	if len(allSlots) > 0 {
		if err := u.slotsRepo.BulkUpsert(ctx, allSlots); err != nil {
			return domain.Schedule{}, err
		}
	}

	return created, nil
}
