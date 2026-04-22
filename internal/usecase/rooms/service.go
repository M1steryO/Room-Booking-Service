package rooms

import "github.com/M1steryO/Room-Booking-Service/internal/repository"

type RoomsUsecase struct {
	roomsRepo repository.RoomRepository
}

func NewRoomsUsecase(roomsRepo repository.RoomRepository) *RoomsUsecase {
	return &RoomsUsecase{roomsRepo: roomsRepo}
}
