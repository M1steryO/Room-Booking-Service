package httpx

import (
	"github.com/avito-internships/test-backend-1-M1steryO/internal/transport/httpx/handlers"
	"github.com/avito-internships/test-backend-1-M1steryO/internal/usecase/auth"
	"github.com/avito-internships/test-backend-1-M1steryO/internal/usecase/bookings"
	"github.com/avito-internships/test-backend-1-M1steryO/internal/usecase/rooms"
	"github.com/avito-internships/test-backend-1-M1steryO/internal/usecase/schedules"
	"github.com/avito-internships/test-backend-1-M1steryO/internal/usecase/slots"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter(
	auth *auth.AuthUsecase,
	rooms *rooms.RoomsUsecase,
	schedules *schedules.SchedulesUsecase,
	slots *slots.SlotsUsecase,
	bookings *bookings.BookingsUsecase,
	authMiddleware func(next http.Handler) http.Handler,
) http.Handler {
	router := chi.NewRouter()

	handler := handlers.New(auth, rooms, schedules, slots, bookings)

	router.Get("/_info", handler.Info)

	router.Post("/dummyLogin", handler.DummyLogin)
	router.Post("/register", handler.Register)
	router.Post("/login", handler.Login)

	router.Group(func(private chi.Router) {
		private.Use(authMiddleware)

		private.Get("/rooms/list", handler.ListRooms)
		private.Post("/rooms/create", handler.CreateRoom)

		private.Post("/rooms/{roomId}/schedule/create", handler.CreateSchedule)
		private.Get("/rooms/{roomId}/slots/list", handler.ListSlots)

		private.Post("/bookings/create", handler.CreateBooking)
		private.Get("/bookings/list", handler.ListBookings)
		private.Get("/bookings/my", handler.MyBookings)
		private.Post("/bookings/{bookingId}/cancel", handler.CancelBooking)
	})

	return router
}
