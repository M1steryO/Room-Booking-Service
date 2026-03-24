package bookings

import (
	"context"

	"github.com/avito-internships/test-backend-1-M1steryO/internal/domain"
)

func (u *BookingsUsecase) Cancel(ctx context.Context, actorID string, actorRole domain.Role, bookingID string) (domain.Booking, error) {
	if actorRole != domain.RoleUser {
		return domain.Booking{}, domain.Forbidden("user role required")
	}

	return u.bookingsRepo.CancelByOwner(ctx, bookingID, actorID)
}
