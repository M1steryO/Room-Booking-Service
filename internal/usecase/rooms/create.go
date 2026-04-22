package rooms

import (
	"context"

	"github.com/M1steryO/Room-Booking-Service/internal/domain"
	"github.com/M1steryO/Room-Booking-Service/pkg/identity"
)

func (u *RoomsUsecase) Create(ctx context.Context, actorRole domain.Role, name string, description *string, capacity *int) (domain.Room, error) {
	if actorRole != domain.RoleAdmin {
		return domain.Room{}, domain.Forbidden("admin role required")
	}

	if name == "" {
		return domain.Room{}, domain.InvalidRequest("name is required")
	}

	room := domain.Room{
		ID:          identity.New(),
		Name:        name,
		Description: description,
		Capacity:    capacity,
	}

	return u.roomsRepo.Create(ctx, room)
}
