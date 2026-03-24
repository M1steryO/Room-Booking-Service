package handlers

import (
	"encoding/json"
	"github.com/avito-internships/test-backend-1-M1steryO/internal/delivery/http/handlers/helpers"
	"net/http"
	"time"

	"github.com/avito-internships/test-backend-1-M1steryO/internal/delivery/http/middleware"
	"github.com/avito-internships/test-backend-1-M1steryO/internal/delivery/http/models"
	"github.com/avito-internships/test-backend-1-M1steryO/internal/domain"
	"github.com/go-chi/chi/v5"
)

// ListRooms godoc
// @Summary List rooms
// @Tags Rooms
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.RoomsEnvelope
// @Failure 401 {object} models.ErrorEnvelope
// @Router /rooms/list [get]
func (h *Handlers) ListRooms(w http.ResponseWriter, r *http.Request) {
	rooms, err := h.rooms.List(r.Context())
	if err != nil {
		helpers.WriteError(w, r, err)
		return
	}

	models.WriteJSON(w, http.StatusOK, models.RoomsEnvelope{Rooms: rooms})
}

// CreateRoom godoc
// @Summary Create room
// @Tags Rooms
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.CreateRoomRequest true "Create room payload"
// @Success 201 {object} models.RoomEnvelope
// @Failure 400 {object} models.ErrorEnvelope
// @Failure 403 {object} models.ErrorEnvelope
// @Router /rooms/create [post]
func (h *Handlers) CreateRoom(w http.ResponseWriter, r *http.Request) {
	var request models.CreateRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		helpers.WriteError(w, r, domain.InvalidRequest("invalid json body"))
		return
	}

	actor := middleware.ActorFromContext(r.Context())
	room, err := h.rooms.Create(r.Context(), actor.Role, request.Name, request.Description, request.Capacity)
	if err != nil {
		helpers.WriteError(w, r, err)
		return
	}

	models.WriteJSON(w, http.StatusCreated, models.RoomEnvelope{Room: room})
}

// CreateSchedule godoc
// @Summary Create room schedule
// @Tags Schedules
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param roomId path string true "Room ID"
// @Param request body models.CreateScheduleRequest true "Create schedule payload"
// @Success 201 {object} models.ScheduleEnvelope
// @Failure 400 {object} models.ErrorEnvelope
// @Failure 403 {object} models.ErrorEnvelope
// @Failure 404 {object} models.ErrorEnvelope
// @Failure 409 {object} models.ErrorEnvelope
// @Router /rooms/{roomId}/schedule/create [post]
func (h *Handlers) CreateSchedule(w http.ResponseWriter, r *http.Request) {
	var request models.CreateScheduleRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		helpers.WriteError(w, r, domain.InvalidRequest("invalid json body"))
		return
	}

	actor := middleware.ActorFromContext(r.Context())
	roomID := chi.URLParam(r, "roomId")

	schedule, err := h.schedules.Create(r.Context(), actor.Role, roomID, request.DaysOfWeek, request.StartTime, request.EndTime)
	if err != nil {
		helpers.WriteError(w, r, err)
		return
	}

	models.WriteJSON(w, http.StatusCreated, models.ScheduleEnvelope{Schedule: schedule})
}

// ListSlots godoc
// @Summary List available room slots by date
// @Tags Slots
// @Produce json
// @Security BearerAuth
// @Param roomId path string true "Room ID"
// @Param date query string true "Date YYYY-MM-DD"
// @Success 200 {object} models.SlotsEnvelope
// @Failure 400 {object} models.ErrorEnvelope
// @Failure 404 {object} models.ErrorEnvelope
// @Router /rooms/{roomId}/slots/list [get]
func (h *Handlers) ListSlots(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "roomId")
	dateRaw := r.URL.Query().Get("date")
	if dateRaw == "" {
		helpers.WriteError(w, r, domain.InvalidRequest("date query parameter is required"))
		return
	}

	date, err := time.Parse("2006-01-02", dateRaw)
	if err != nil {
		helpers.WriteError(w, r, domain.InvalidRequest("date must be in format YYYY-MM-DD"))
		return
	}

	slots, err := h.slots.ListAvailable(r.Context(), roomID, date)
	if err != nil {
		helpers.WriteError(w, r, err)
		return
	}

	models.WriteJSON(w, http.StatusOK, models.SlotsEnvelope{Slots: slots})
}
