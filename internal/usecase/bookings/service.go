package bookings

import (
	"github.com/avito-internships/test-backend-1-M1steryO/internal/client"
	"github.com/avito-internships/test-backend-1-M1steryO/internal/repository"
	"github.com/avito-internships/test-backend-1-M1steryO/pkg/clock"
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
