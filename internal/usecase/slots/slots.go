package slots

import (
	"context"
	"github.com/avito-internships/test-backend-1-M1steryO/internal/domain"
	"github.com/avito-internships/test-backend-1-M1steryO/internal/usecase/schedules/helpers"
	"time"
)

func (u *SlotsUsecase) ListAvailable(ctx context.Context, roomID string, date time.Time) ([]domain.Slot, error) {
	if _, err := u.roomsRepo.GetByID(ctx, roomID); err != nil {
		return nil, err
	}

	schedule, err := u.schedulesRepo.GetByRoomID(ctx, roomID)
	if err != nil {
		if domain.AsAppError(err).Code == "NOT_FOUND" {
			return []domain.Slot{}, nil
		}
		return nil, err
	}

	hasSlots, err := u.slotsRepo.DateHasSlots(ctx, roomID, date)
	if err != nil {
		return nil, err
	}

	if !hasSlots {
		generated, err := helpers.GenerateSlotsForDate(roomID, schedule, date)
		if err != nil {
			return nil, domain.InvalidRequest(err.Error())
		}

		if len(generated) > 0 {
			if err := u.slotsRepo.BulkUpsert(ctx, generated); err != nil {
				return nil, err
			}
		}
	}

	return u.slotsRepo.ListAvailableByRoomAndDate(ctx, roomID, date)
}
