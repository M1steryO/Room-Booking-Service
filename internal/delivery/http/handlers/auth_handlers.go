package handlers

import (
	"encoding/json"
	"github.com/M1steryO/Room-Booking-Service/internal/delivery/http/handlers/helpers"
	"net/http"

	"github.com/M1steryO/Room-Booking-Service/internal/delivery/http/models"
	"github.com/M1steryO/Room-Booking-Service/internal/domain"
)

// Info godoc
// @Summary Healthcheck
// @Tags System
// @Produce json
// @Success 200 {object} map[string]string
// @Router /_info [get]
func (h *Handlers) Info(w http.ResponseWriter, _ *http.Request) {
	models.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// DummyLogin godoc
// @Summary Get test JWT by role
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.DummyLoginRequest true "Role"
// @Success 200 {object} models.TokenResponse
// @Failure 400 {object} models.ErrorEnvelope
// @Router /dummyLogin [post]
func (h *Handlers) DummyLogin(w http.ResponseWriter, r *http.Request) {
	var request models.DummyLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		helpers.WriteError(w, r, domain.InvalidRequest("invalid json body"))
		return
	}

	token, err := h.auth.DummyLogin(request.Role)
	if err != nil {
		helpers.WriteError(w, r, err)
		return
	}

	models.WriteJSON(w, http.StatusOK, models.TokenResponse{Token: token})
}

// Register godoc
// @Summary Register user
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "Register payload"
// @Success 201 {object} models.UserEnvelope
// @Failure 400 {object} models.ErrorEnvelope
// @Failure 500 {object} models.ErrorEnvelope
// @Router /register [post]
func (h *Handlers) Register(w http.ResponseWriter, r *http.Request) {
	var request models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		helpers.WriteError(w, r, domain.InvalidRequest("invalid json body"))
		return
	}

	user, err := h.auth.Register(r.Context(), request.Email, request.Password, request.Role)
	if err != nil {
		helpers.WriteError(w, r, err)
		return
	}

	models.WriteJSON(w, http.StatusCreated, models.UserEnvelope{User: models.NewUserResponse(user)})
}

// Login godoc
// @Summary Login by email and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login payload"
// @Success 200 {object} models.TokenResponse
// @Failure 401 {object} models.ErrorEnvelope
// @Failure 500 {object} models.ErrorEnvelope
// @Router /login [post]
func (h *Handlers) Login(w http.ResponseWriter, r *http.Request) {
	var request models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		helpers.WriteError(w, r, domain.InvalidRequest("invalid json body"))
		return
	}

	token, err := h.auth.Login(r.Context(), request.Email, request.Password)
	if err != nil {
		helpers.WriteError(w, r, err)
		return
	}

	models.WriteJSON(w, http.StatusOK, models.TokenResponse{Token: token})
}
