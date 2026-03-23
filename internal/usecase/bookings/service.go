package bookings

import (
	"github.com/avito-internships/test-backend-1-M1steryO/internal/conference"
	"github.com/avito-internships/test-backend-1-M1steryO/internal/platform/clock"
	"github.com/avito-internships/test-backend-1-M1steryO/internal/repository"
)

type BookingsUsecase struct {
	slotsRepo      repository.SlotRepository
	bookingsRepo   repository.BookingRepository
	conference     conference.Service
	clock          clock.Clock
}

func NewBookingsUsecase(
	slotsRepo repository.SlotRepository,
	bookingsRepo repository.BookingRepository,
	conferenceService conference.Service,
	clk clock.Clock,
) *BookingsUsecase {
	return &BookingsUsecase{
		slotsRepo:    slotsRepo,
		bookingsRepo: bookingsRepo,
		conference:   conferenceService,
		clock:        clk,
	}
}
