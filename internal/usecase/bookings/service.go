package bookings

import (
	"github.com/M1steryO/Room-Booking-Service/internal/client"
	"github.com/M1steryO/Room-Booking-Service/internal/repository"
	"github.com/M1steryO/Room-Booking-Service/pkg/clock"
)

type BookingsUsecase struct {
	slotsRepo    repository.SlotRepository
	bookingsRepo repository.BookingRepository
	conference   client.ConferenceService
	clock        clock.Clock
}

func NewBookingsUsecase(
	slotsRepo repository.SlotRepository,
	bookingsRepo repository.BookingRepository,
	conferenceService client.ConferenceService,
	clk clock.Clock,
) *BookingsUsecase {
	return &BookingsUsecase{
		slotsRepo:    slotsRepo,
		bookingsRepo: bookingsRepo,
		conference:   conferenceService,
		clock:        clk,
	}
}
