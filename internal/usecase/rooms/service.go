package rooms

import "github.com/avito-internships/test-backend-1-M1steryO/internal/repository"

type RoomsUsecase struct {
	roomsRepo repository.RoomRepository
}

func NewRoomsUsecase(roomsRepo repository.RoomRepository) *RoomsUsecase {
	return &RoomsUsecase{roomsRepo: roomsRepo}
}
