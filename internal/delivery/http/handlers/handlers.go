package handlers

import (
	"github.com/M1steryO/Room-Booking-Service/internal/usecase/auth"
	"github.com/M1steryO/Room-Booking-Service/internal/usecase/bookings"
	"github.com/M1steryO/Room-Booking-Service/internal/usecase/rooms"
	"github.com/M1steryO/Room-Booking-Service/internal/usecase/schedules"
	"github.com/M1steryO/Room-Booking-Service/internal/usecase/slots"
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
