package schedules

import (
	"github.com/M1steryO/Room-Booking-Service/internal/repository"
	"github.com/M1steryO/Room-Booking-Service/pkg/clock"
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
