package rooms

import (
	"context"
	"github.com/avito-internships/test-backend-1-M1steryO/internal/domain"
	"github.com/avito-internships/test-backend-1-M1steryO/internal/platform/identity"
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

func (u *RoomsUsecase) List(ctx context.Context) ([]domain.Room, error) {
	return u.roomsRepo.List(ctx)
}
