package rooms

import (
	"context"

	"github.com/avito-internships/test-backend-1-M1steryO/internal/domain"
)

func (u *RoomsUsecase) List(ctx context.Context) ([]domain.Room, error) {
	return u.roomsRepo.List(ctx)
}
