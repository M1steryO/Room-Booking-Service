package bookings

import (
	"context"

	"github.com/M1steryO/Room-Booking-Service/internal/domain"
)

func (u *BookingsUsecase) ListMine(ctx context.Context, actorID string, actorRole domain.Role) ([]domain.Booking, error) {
	if actorRole != domain.RoleUser {
		return nil, domain.Forbidden("user role required")
	}

	return u.bookingsRepo.ListByUser(ctx, actorID, u.clock.Now())
}
