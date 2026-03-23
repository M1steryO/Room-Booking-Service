package schedules

import (
	"github.com/avito-internships/test-backend-1-M1steryO/internal/platform/clock"
	"github.com/avito-internships/test-backend-1-M1steryO/internal/repository"
)

type SchedulesUsecase struct {
	roomsRepo      repository.RoomRepository
	schedulesRepo  repository.ScheduleRepository
	slotsRepo      repository.SlotRepository
	clock          clock.Clock
	slotHorizonDay int
}

func NewSchedulesUsecase(
	roomsRepo repository.RoomRepository,
	schedulesRepo repository.ScheduleRepository,
	slotsRepo repository.SlotRepository,
	clk clock.Clock,
	slotHorizonDays int,
) *SchedulesUsecase {
	return &SchedulesUsecase{
		roomsRepo:      roomsRepo,
		schedulesRepo:  schedulesRepo,
		slotsRepo:      slotsRepo,
		clock:          clk,
		slotHorizonDay: slotHorizonDays,
	}
}
