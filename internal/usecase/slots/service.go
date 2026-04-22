package slots

import "github.com/M1steryO/Room-Booking-Service/internal/repository"

type SlotsUsecase struct {
	roomsRepo     repository.RoomRepository
	schedulesRepo repository.ScheduleRepository
	slotsRepo     repository.SlotRepository
}

func NewSlotsUsecase(roomsRepo repository.RoomRepository, schedulesRepo repository.ScheduleRepository, slotsRepo repository.SlotRepository) *SlotsUsecase {
	return &SlotsUsecase{
		roomsRepo:     roomsRepo,
		schedulesRepo: schedulesRepo,
		slotsRepo:     slotsRepo,
	}
}
