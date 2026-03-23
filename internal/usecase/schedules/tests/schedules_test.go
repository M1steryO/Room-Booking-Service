package schedules_test

import (
	"context"
	"testing"
	"time"

	"github.com/avito-internships/test-backend-1-M1steryO/internal/domain"
	repmocks "github.com/avito-internships/test-backend-1-M1steryO/internal/repository/mocks"
	"github.com/avito-internships/test-backend-1-M1steryO/internal/usecase/schedules"
)

type fixedClock struct{ now time.Time }

func (f fixedClock) Now() time.Time { return f.now }

func TestCreateForbiddenForUser(t *testing.T) {
	roomsRepo := repmocks.NewRoomRepositoryMock(t)
	schedulesRepo := repmocks.NewScheduleRepositoryMock(t)
	slotsRepo := repmocks.NewSlotRepositoryMock(t)
	roomsRepo.GetByIDMock.Optional()
	roomsRepo.CreateMock.Optional()
	roomsRepo.ListMock.Optional()
	schedulesRepo.CreateMock.Optional()
	schedulesRepo.GetByRoomIDMock.Optional()
	slotsRepo.BulkUpsertMock.Optional()
	slotsRepo.DateHasSlotsMock.Optional()
	slotsRepo.GetByIDMock.Optional()
	slotsRepo.ListAvailableByRoomAndDateMock.Optional()
	uc := schedules.NewSchedulesUsecase(
		roomsRepo,
		schedulesRepo,
		slotsRepo,
		fixedClock{now: time.Now().UTC()},
		7,
	)

	_, err := uc.Create(context.Background(), domain.RoleUser, "r1", []int{1}, "09:00", "10:00")
	if err == nil || domain.AsAppError(err).Code != "FORBIDDEN" {
		t.Fatalf("expected FORBIDDEN, got %v", err)
	}
}

func TestCreateSuccessGeneratesSlots(t *testing.T) {
	var upserted int
	roomsRepo := repmocks.NewRoomRepositoryMock(t)
	schedulesRepo := repmocks.NewScheduleRepositoryMock(t)
	slotsRepo := repmocks.NewSlotRepositoryMock(t)
	roomsRepo.GetByIDMock.Set(func(_ context.Context, _ string) (domain.Room, error) { return domain.Room{ID: "r1"}, nil })
	schedulesRepo.CreateMock.Set(func(_ context.Context, s domain.Schedule) (domain.Schedule, error) { return s, nil })
	slotsRepo.BulkUpsertMock.Set(func(_ context.Context, slots []domain.Slot) error {
		upserted = len(slots)
		return nil
	})
	roomsRepo.CreateMock.Optional()
	roomsRepo.ListMock.Optional()
	schedulesRepo.GetByRoomIDMock.Optional()
	slotsRepo.DateHasSlotsMock.Optional()
	slotsRepo.GetByIDMock.Optional()
	slotsRepo.ListAvailableByRoomAndDateMock.Optional()
	uc := schedules.NewSchedulesUsecase(
		roomsRepo,
		schedulesRepo,
		slotsRepo,
		fixedClock{now: time.Date(2026, 3, 23, 8, 0, 0, 0, time.UTC)},
		2,
	)

	_, err := uc.Create(context.Background(), domain.RoleAdmin, "r1", []int{1}, "09:00", "10:00")
	if err != nil {
		t.Fatalf("Create error: %v", err)
	}
	if upserted == 0 {
		t.Fatal("expected generated slots to be upserted")
	}
}
