package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/avito-internships/test-backend-1-M1steryO/internal/domain"
	"github.com/avito-internships/test-backend-1-M1steryO/internal/transport/httpx/middleware"
	"github.com/avito-internships/test-backend-1-M1steryO/internal/transport/httpx/models"
	"github.com/go-chi/chi/v5"
)

// CreateBooking godoc
// @Summary Create booking
// @Tags Bookings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateBookingRequest true "Create booking payload"
// @Success 201 {object} BookingEnvelope
// @Failure 400 {object} ErrorEnvelope
// @Failure 403 {object} ErrorEnvelope
// @Failure 404 {object} ErrorEnvelope
// @Failure 409 {object} ErrorEnvelope
// @Router /bookings/create [post]
func (h *Handlers) CreateBooking(w http.ResponseWriter, r *http.Request) {
	var request models.CreateBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		models.WriteError(w, domain.InvalidRequest("invalid json body"))
		return
	}

	actor := middleware.ActorFromContext(r.Context())
	
	booking, err := h.bookings.Create(r.Context(), actor.UserID, actor.Role, request.SlotID, request.CreateConferenceLink)
	if err != nil {
		models.WriteError(w, err)
		return
	}

	models.WriteJSON(w, http.StatusCreated, models.BookingEnvelope{Booking: booking})
}

// ListBookings godoc
// @Summary List all bookings
// @Tags Bookings
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page"
// @Param pageSize query int false "Page size"
// @Success 200 {object} BookingsListEnvelope
// @Failure 400 {object} ErrorEnvelope
// @Failure 403 {object} ErrorEnvelope
// @Router /bookings/list [get]
func (h *Handlers) ListBookings(w http.ResponseWriter, r *http.Request) {
	actor := middleware.ActorFromContext(r.Context())
	page := parsePositiveInt(r.URL.Query().Get("page"), 1)
	pageSize := parsePositiveInt(r.URL.Query().Get("pageSize"), 20)

	bookings, total, err := h.bookings.ListAll(r.Context(), actor.Role, page, pageSize)
	if err != nil {
		models.WriteError(w, err)
		return
	}

	models.WriteJSON(w, http.StatusOK, models.BookingsListEnvelope{
		Bookings: bookings,
		Pagination: models.Pagination{
			Page:     page,
			PageSize: pageSize,
			Total:    total,
		},
	})
}

// MyBookings godoc
// @Summary List current user bookings
// @Tags Bookings
// @Produce json
// @Security BearerAuth
// @Success 200 {object} BookingsEnvelope
// @Failure 403 {object} ErrorEnvelope
// @Router /bookings/my [get]
func (h *Handlers) MyBookings(w http.ResponseWriter, r *http.Request) {
	actor := middleware.ActorFromContext(r.Context())
	bookings, err := h.bookings.ListMine(r.Context(), actor.UserID, actor.Role)
	if err != nil {
		models.WriteError(w, err)
		return
	}

	models.WriteJSON(w, http.StatusOK, models.BookingsEnvelope{Bookings: bookings})
}

// CancelBooking godoc
// @Summary Cancel booking
// @Tags Bookings
// @Produce json
// @Security BearerAuth
// @Param bookingId path string true "Booking ID"
// @Success 200 {object} BookingEnvelope
// @Failure 403 {object} ErrorEnvelope
// @Failure 404 {object} ErrorEnvelope
// @Router /bookings/{bookingId}/cancel [post]
func (h *Handlers) CancelBooking(w http.ResponseWriter, r *http.Request) {
	actor := middleware.ActorFromContext(r.Context())
	bookingID := chi.URLParam(r, "bookingId")

	booking, err := h.bookings.Cancel(r.Context(), actor.UserID, actor.Role, bookingID)
	if err != nil {
		models.WriteError(w, err)
		return
	}

	models.WriteJSON(w, http.StatusOK, models.BookingEnvelope{Booking: booking})
}
