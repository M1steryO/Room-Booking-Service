package bookings

import (
	"context"

	"github.com/avito-internships/test-backend-1-M1steryO/internal/domain"
	"github.com/google/uuid"
)

func (u *BookingsUsecase) Cancel(ctx context.Context, actorID string, actorRole domain.Role, bookingID string) (domain.Booking, error) {
	if actorRole != domain.RoleUser {
		return domain.Booking{}, domain.Forbidden("user role required")
	}
	if _, err := uuid.Parse(bookingID); err != nil {
		return domain.Booking{}, domain.InvalidRequest("invalid booking_id")
	}
	booking, err := u.bookingsRepo.GetByID(ctx, bookingID)
	if err != nil {
		return domain.Booking{}, err
	}

	if booking.UserID != actorID {
		return domain.Booking{}, domain.Forbidden("cannot cancel another user's booking")
	}

	err = u.bookingsRepo.CancelByOwner(ctx, bookingID, actorID)
	if err != nil {
		return domain.Booking{}, err
	}

	booking.Status = domain.BookingStatusCancelled
	return booking, nil
}
