package rooms

import (
	"context"

	"github.com/M1steryO/Room-Booking-Service/internal/domain"
)

func (u *RoomsUsecase) List(ctx context.Context) ([]domain.Room, error) {
	return u.roomsRepo.List(ctx)
}
