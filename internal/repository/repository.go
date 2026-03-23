package repository

import (
	"context"
	"github.com/avito-internships/test-backend-1-M1steryO/internal/domain"
	"time"
)

type UserRepository interface {
	Create(ctx context.Context, user domain.User) (domain.User, error)
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	UpsertSystemUsers(ctx context.Context, users []domain.User) error
}

type RoomRepository interface {
	Create(ctx context.Context, room domain.Room) (domain.Room, error)
	List(ctx context.Context) ([]domain.Room, error)
	GetByID(ctx context.Context, roomID string) (domain.Room, error)
}

type ScheduleRepository interface {
	Create(ctx context.Context, schedule domain.Schedule) (domain.Schedule, error)
	GetByRoomID(ctx context.Context, roomID string) (domain.Schedule, error)
}

type SlotRepository interface {
	DateHasSlots(ctx context.Context, roomID string, date time.Time) (bool, error)
	BulkUpsert(ctx context.Context, slots []domain.Slot) error
	GetByID(ctx context.Context, slotID string) (domain.Slot, error)
	ListAvailableByRoomAndDate(ctx context.Context, roomID string, date time.Time) ([]domain.Slot, error)
}

type BookingRepository interface {
	Create(ctx context.Context, booking domain.Booking) (domain.Booking, error)
	ListAll(ctx context.Context, page, pageSize int) ([]domain.Booking, int, error)
	ListFutureByUser(ctx context.Context, userID string, now time.Time) ([]domain.Booking, error)
	GetByID(ctx context.Context, bookingID string) (domain.Booking, error)
	CancelByOwner(ctx context.Context, bookingID, userID string) (domain.Booking, error)
}
