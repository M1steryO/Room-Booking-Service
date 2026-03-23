package bookings

import (
	"context"
	"github.com/avito-internships/test-backend-1-M1steryO/internal/domain"
	"github.com/avito-internships/test-backend-1-M1steryO/internal/platform/identity"
)

func (u *BookingsUsecase) Create(ctx context.Context, actorID string, actorRole domain.Role, slotID string, createConferenceLink bool) (domain.Booking, error) {
	if actorRole != domain.RoleUser {
		return domain.Booking{}, domain.Forbidden("booking available only for user role")
	}

	slot, err := u.slotsRepo.GetByID(ctx, slotID)
	if err != nil {
		return domain.Booking{}, err
	}

	if slot.Start.Before(u.clock.Now()) {
		return domain.Booking{}, domain.InvalidRequest("cannot create booking for slot in the past")
	}

	booking := domain.Booking{
		ID:        identity.New(),
		SlotID:    slot.ID,
		UserID:    actorID,
		Status:    domain.BookingStatusActive,
		CreatedAt: u.clock.Now(),
		RoomID:    slot.RoomID,
		SlotStart: slot.Start,
		SlotEnd:   slot.End,
	}

	if createConferenceLink {
		link, err := u.conference.CreateLink(ctx, booking.ID)
		if err != nil {
			return domain.Booking{}, err
		}

		booking.ConferenceLink = &link
	}

	return u.bookingsRepo.Create(ctx, booking)
}

func (u *BookingsUsecase) ListAll(ctx context.Context, actorRole domain.Role, page int, pageSize int) ([]domain.Booking, int, error) {
	if actorRole != domain.RoleAdmin {
		return nil, 0, domain.Forbidden("admin role required")
	}

	if page <= 0 || pageSize <= 0 || pageSize > 100 {
		return nil, 0, domain.InvalidRequest("invalid pagination parameters")
	}

	return u.bookingsRepo.ListAll(ctx, page, pageSize)
}

func (u *BookingsUsecase) ListMine(ctx context.Context, actorID string, actorRole domain.Role) ([]domain.Booking, error) {
	if actorRole != domain.RoleUser {
		return nil, domain.Forbidden("user role required")
	}

	return u.bookingsRepo.ListFutureByUser(ctx, actorID, u.clock.Now())
}

func (u *BookingsUsecase) Cancel(ctx context.Context, actorID string, actorRole domain.Role, bookingID string) (domain.Booking, error) {
	if actorRole != domain.RoleUser {
		return domain.Booking{}, domain.Forbidden("user role required")
	}

	return u.bookingsRepo.CancelByOwner(ctx, bookingID, actorID)
}
