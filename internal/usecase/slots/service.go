package slots

import "github.com/avito-internships/test-backend-1-M1steryO/internal/repository"

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
