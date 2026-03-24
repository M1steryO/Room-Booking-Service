package bookings

import (
	"context"

	"github.com/avito-internships/test-backend-1-M1steryO/internal/domain"
)

func (u *BookingsUsecase) ListMine(ctx context.Context, actorID string, actorRole domain.Role) ([]domain.Booking, error) {
	if actorRole != domain.RoleUser {
		return nil, domain.Forbidden("user role required")
	}

	return u.bookingsRepo.ListByUser(ctx, actorID, u.clock.Now())
}
