package models

import (
	"github.com/avito-internships/test-backend-1-M1steryO/internal/domain"
)

type DummyLoginRequest struct {
	Role domain.Role `json:"role"`
}

type RegisterRequest struct {
	Email    string      `json:"email"`
	Password string      `json:"password"`
	Role     domain.Role `json:"role"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateRoomRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Capacity    *int    `json:"capacity"`
}

type CreateScheduleRequest struct {
	DaysOfWeek []int  `json:"daysOfWeek"`
	StartTime  string `json:"startTime"`
	EndTime    string `json:"endTime"`
}

type CreateBookingRequest struct {
	SlotID               string `json:"slotId"`
	CreateConferenceLink bool   `json:"createConferenceLink"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

type UserEnvelope struct {
	User domain.User `json:"user"`
}

type RoomEnvelope struct {
	Room domain.Room `json:"room"`
}

type RoomsEnvelope struct {
	Rooms []domain.Room `json:"rooms"`
}

type ScheduleEnvelope struct {
	Schedule domain.Schedule `json:"schedule"`
}

type SlotsEnvelope struct {
	Slots []domain.Slot `json:"slots"`
}

type BookingEnvelope struct {
	Booking domain.Booking `json:"booking"`
}

type BookingsEnvelope struct {
	Bookings []domain.Booking `json:"bookings"`
}

type BookingsListEnvelope struct {
	Bookings   []domain.Booking `json:"bookings"`
	Pagination Pagination       `json:"pagination"`
}

type Pagination struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
	Total    int `json:"total"`
}
