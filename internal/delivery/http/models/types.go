package models

import (
	"github.com/M1steryO/Room-Booking-Service/internal/domain"
	"time"
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

type UserResponse struct {
	ID        string      `json:"id"`
	Email     string      `json:"email"`
	Role      domain.Role `json:"role"`
	CreatedAt time.Time   `json:"createdAt"`
}

type UserEnvelope struct {
	User UserResponse `json:"user"`
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

func NewUserResponse(user domain.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}
}
