package bookings

import (
	"context"

	"github.com/M1steryO/Room-Booking-Service/internal/domain"
)

func (u *BookingsUsecase) ListAll(ctx context.Context, actorRole domain.Role, page int, pageSize int) ([]domain.Booking, int, error) {
	if actorRole != domain.RoleAdmin {
		return nil, 0, domain.Forbidden("admin role required")
	}

	if page <= 0 || pageSize <= 0 || pageSize > 100 {
		return nil, 0, domain.InvalidRequest("invalid pagination parameters")
	}

	return u.bookingsRepo.ListAll(ctx, page, pageSize)
}
