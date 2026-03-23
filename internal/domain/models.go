package domain

import "time"

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

func (r Role) Valid() bool {
	return r == RoleAdmin || r == RoleUser
}

type BookingStatus string

const (
	BookingStatusActive    BookingStatus = "active"
	BookingStatusCancelled BookingStatus = "cancelled"
)

type User struct {
	ID           string
	Email        string
	PasswordHash string
	Role         Role
	CreatedAt    time.Time
}

type Room struct {
	ID          string
	Name        string
	Description *string
	Capacity    *int
	CreatedAt   time.Time
}

type Schedule struct {
	ID         string
	RoomID     string
	DaysOfWeek []int
	StartTime  string
	EndTime    string
	CreatedAt  time.Time
}

type Slot struct {
	ID        string
	RoomID    string
	Start     time.Time
	End       time.Time
	CreatedAt time.Time
}

type Booking struct {
	ID             string
	SlotID         string
	UserID         string
	Status         BookingStatus
	ConferenceLink *string
	CreatedAt      time.Time

	SlotStart time.Time
	SlotEnd   time.Time
	RoomID    string
}
