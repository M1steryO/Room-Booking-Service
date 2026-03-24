package bookings

import (
	"context"

	"github.com/google/uuid"

	"github.com/avito-internships/test-backend-1-M1steryO/internal/domain"
	"github.com/avito-internships/test-backend-1-M1steryO/pkg/identity"
)

func (u *BookingsUsecase) Create(ctx context.Context, actorID string, actorRole domain.Role, slotID string, createConferenceLink bool) (domain.Booking, error) {
	if actorRole != domain.RoleUser {
		return domain.Booking{}, domain.Forbidden("booking available only for user role")
	}
	if _, err := uuid.Parse(slotID); err != nil {
		return domain.Booking{}, domain.InvalidRequest("invalid slot_id")
	}

	slot, err := u.slotsRepo.GetByID(ctx, slotID)
	if err != nil {
		return domain.Booking{}, err
	}

	if slot.Start.Before(u.clock.Now()) {
		return domain.Booking{}, domain.InvalidRequest("can't create booking for slot in the past")
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
