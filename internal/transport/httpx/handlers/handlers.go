package handlers

import (
	"strconv"

	"github.com/avito-internships/test-backend-1-M1steryO/internal/usecase/auth"
	"github.com/avito-internships/test-backend-1-M1steryO/internal/usecase/bookings"
	"github.com/avito-internships/test-backend-1-M1steryO/internal/usecase/rooms"
	"github.com/avito-internships/test-backend-1-M1steryO/internal/usecase/schedules"
	"github.com/avito-internships/test-backend-1-M1steryO/internal/usecase/slots"
)

type Handlers struct {
	auth      *auth.AuthUsecase
	rooms     *rooms.RoomsUsecase
	schedules *schedules.SchedulesUsecase
	slots     *slots.SlotsUsecase
	bookings  *bookings.BookingsUsecase
}

func New(
	auth *auth.AuthUsecase,
	rooms *rooms.RoomsUsecase,
	schedules *schedules.SchedulesUsecase,
	slots *slots.SlotsUsecase,
	bookings *bookings.BookingsUsecase,
) *Handlers {
	return &Handlers{
		auth:      auth,
		rooms:     rooms,
		schedules: schedules,
		slots:     slots,
		bookings:  bookings,
	}
}

func parsePositiveInt(value string, fallback int) int {
	if value == "" {
		return fallback
	}

	number, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}

	return number
}
